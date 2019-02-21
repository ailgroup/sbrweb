package srvc

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

// Session holds sabre session data and other fields for handling in the SessionPool
type Session struct {
	ID            string
	FaultError    error
	TimeValidated time.Time
	TimeStarted   time.Time
	ExpireTime    time.Time
	Sabre         SessionCreateResponse
}

// ExpireScheme for when to expire sessions
type ExpireScheme struct {
	Max int
	Min int
}

// SessionPool container for pool of sessions with specs on size, cycles, counters, errors, timers, configuration, etc...
type SessionPool struct {
	NumberOfCycles  int
	ConfigPoolSize  int
	AllowPoolSize   int
	Expire          ExpireScheme
	CycleEvery      time.Duration
	ServiceURL      string
	NetworkErrors   []error
	FaultErrors     []error
	InitializedTime time.Time
	Sessions        chan Session
	Conf            *SessionConf
}

// NewPool initializes a new session pool of given size, for ServiceURL, with a buffered channel of
// type Session, and with specification for cycled keepalive checks, expiration and
// timer for an initialized time.
func NewPool(expire ExpireScheme, cred *SessionConf, cycle time.Duration, size int) *SessionPool {
	return &SessionPool{
		ConfigPoolSize:  size,
		Sessions:        make(chan Session, size), //buffered channel blocks!
		InitializedTime: time.Now(),
		ServiceURL:      cred.ServiceURL,
		CycleEvery:      cycle,
		Expire:          expire,
		Conf:            cred,
	}
}

// Deamonize initializes and populates new session pool, accepts signal handler to
// gracefully manage valid shutdown of Sabre sessions, session pool, and keepalive.
// For example: pool.Deamonize(os.Interrupt)
func (p *SessionPool) Deamonize(sig os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, sig)
	go func() {
		err := p.Populate()
		if err != nil {
			fmt.Printf("Error popluating session pool %v\n", err)
			os.Exit(1)
		}
		p.Keepalive(done)
		fmt.Printf("\nGot '%s' SIGNAL. Shutting down keepalive and session pool; exiting program...\n", sig)
		p.Close()
		os.Exit(0)
	}()
}

// GenerateSessionID for small easy to find ids in logs; returns format 'xca123'
func GenerateSessionID() string {
	randStr := randStringBytesMaskImprSrc(3)
	nowtime := time.Now().Format(".999")
	return randStr + nowtime
}

// RandomInt select random integer within range min/max.
// This is used to better randomize expiration times of sessions
// fitting within a min and max time range.
func RandomInt(min, max int) int {
	//make random index, with Permute [0,n] for size min-max+1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := max - min + 1
	p := r.Perm(idx)
	//make key from position idx-3 of permuted slice
	key := p[idx-3]
	a := make([]int, idx)
	for i := range a {
		a[i] = min + i
	}
	//key to select position value out of range [min...max]
	return a[key]
}

// newSession initializes a new valid Sabre session for AAA workspace. It checks new sessions
// can be added to workspace, checks for any faults on creation, creates a new Session struct
// with metadata, logs it, and retursn that Session for placement into the pool.
func (p *SessionPool) newSession() (Session, error) {
	createRQ := BuildSessionCreateRequest(p.Conf)
	createRS, err := CallSessionCreate(p.ServiceURL, createRQ)
	if err != nil {
		p.NetworkErrors = append(p.NetworkErrors, err)
		// create is special, we don't want to put crappy sessions into the buffer
		return Session{}, err
	}
	now := time.Now()
	var faultErr error
	fc := createRS.Body.Fault.Code
	if fc != "" {
		st := createRS.Body.Fault.Detail.StackTrace
		fs := createRS.Body.Fault.String
		faultErr = fmt.Errorf("%s-%s: %s", fs, fc, st)
		p.FaultErrors = append(p.FaultErrors, faultErr)
	}
	//still want to fill buffer even if we get a fault from Sabre
	// there may sceanrios where this is legit (sabre flushes all sessions weekly,
	// or we are over threshold... don't want to block the queue or try to keep
	// filling it with bad sessions. Instead, accept bad sessions, and let the
	// vleanup will clear out faulted sessions later).
	sess := Session{
		ID:            GenerateSessionID(),
		Sabre:         createRS,
		TimeStarted:   now,
		TimeValidated: now,
		ExpireTime:    now.Add(time.Minute * time.Duration(RandomInt(p.Expire.Min, p.Expire.Max))),
		FaultError:    faultErr,
	}
	logSession.Printf(
		"Status='%s' created session ID=%s with Expirey='%s' for token=%s", sess.Sabre.Body.SessionCreateRS.Status,
		sess.ID,
		sess.ExpireTime,
		SabreTokenParse(sess.Sabre.Header.Security.BinarySecurityToken.Value),
	)
	return sess, nil
}

// Populate puts sessions into a buffered channel of SessionPool if the current
// pool size will allow it; it logs and collects any errors (see newSession())
// and attempts to fill during network errors but still honoring the pool size.
// Note that if the SOAP service fails that will be located in the SOAP Fault field, not as
// a general network error. We want to populate even if we get a fault from Sabre:
// many scenarios where this is legit (sabre flushes all sessions weekly, sabre sessions
// service goes down, or we are over session limit...).
// Under these conditions we don't want to block or repeatedly attempt to populate;
// instead, accept a bad session and let the keepalive cleanup bad sessions later.
func (p *SessionPool) Populate() error {
	var err error
	ok := (p.AllowPoolSize == len(p.Sessions))
	for i := 0; i < p.ConfigPoolSize; i++ {
		sess, err := p.newSession()
		if err != nil {
			logSession.Printf("ERROR %v for attempt=%d, continuing...", err, i)
			continue
		}
		p.Sessions <- sess
		p.AllowPoolSize++
	}
	//close blocking channel, message, return error
	if p.AllowPoolSize == 0 {
		err = fmt.Errorf("Cannot create valid sessions closing pool")
		p.NetworkErrors = append(p.NetworkErrors, err)
		close(p.Sessions)
		ok = (p.AllowPoolSize == len(p.Sessions))
	}
	logSession.Printf("Create AllowPoolSize=%d, Create OK=%v. NetworkErrors=%v, FaultErrors=%v", p.AllowPoolSize, ok, p.NetworkErrors, p.FaultErrors)
	return err
}

// Pick session from buffered queue, returns the Session.
func (p *SessionPool) Pick() Session {
	sess := <-p.Sessions
	p.logReport("Pick-" + sess.ID)
	return sess
}

// Put session back onto the buffered queue.
func (p *SessionPool) Put(sess Session) {
	p.Sessions <- sess
	p.logReport("Put-" + sess.ID)
}

// logReport helper to log info about session pool
func (p *SessionPool) logReport(ctx string) {
	open := len(p.Sessions)
	allow := p.AllowPoolSize
	busy := (allow - open)
	logSession.Printf("[%s] ALLOW=%d, OPEN=%d, BUSY=%d, CYCLES=%d, OK=%v", ctx, allow, open, busy, p.NumberOfCycles, (allow == (open + busy)))
}

// RangeKeepalive pulls session out of the pool, checks if expire time is over current time,
// and if so it validates the session against Sabre (which forces Sabre to extend the lifetime)
// and we reset the expire time, placing session back into the pool. Otherwise
// we place session back into the pool leaving the expire time untouched.
func (p *SessionPool) RangeKeepalive(keepaliveID string) {
	breaker := len(p.Sessions)
	counter := 0
	for sess := range p.Sessions {
		//counter==breaker: no looping in the range indefinitely. 1 pass of buffer size is enough
		if counter == breaker {
			p.Sessions <- sess
			return
		}
		counter++
		if time.Now().After(sess.ExpireTime) {
			validateRQ := BuildSessionValidateRequest(
				sess.Sabre.Header.MessageHeader.From.PartyID.Value,
				sess.Sabre.Header.MessageHeader.CPAID, //don't need, send anyway
				sess.Sabre.Header.Security.BinarySecurityToken.Value,
				sess.Sabre.Header.MessageHeader.ConversationID,
				sess.Sabre.Header.MessageHeader.MessageData.RefToMessageID,
				SabreTimeNowFmt(),
			)
			validateRS, err := CallSessionValidate(p.ServiceURL, validateRQ)
			if err != nil {
				//if network error, log and continue. We'll update the queue item with a new expire and allow it to cycle through again. The session may still be valid and useable even if the session validate endpoint is down. Even if it is no longer valid, we don't want to dequeue the pool becuase if sabre is totally down we will end up with an empty queue that will block forever. If Sabre is down they are down, a nothing we can do, so we just go forward as usual and self-repair as Sabre services come back online.
				logSession.Print(err)
			}
			if validateRS.Header.MessageHeader.Action == StatusErrorRS {
				msg := fmt.Sprintf(
					"%s %s, %s, %s",
					StatusErrorRS,
					validateRS.Body.Fault.String,
					validateRS.Body.Fault.Code,
					validateRS.Body.Fault.Detail.StackTrace,
				)
				logSession.Printf("FAULT='%s', %s\n", validateRS.Header.MessageHeader.Action, msg)
				newSess, err := p.newSession()
				if err != nil {
					logSession.Printf("Network ERROR for ID=%s, expire and retry", newSess.ID)
					newSess.ExpireTime = time.Now().Add(time.Second * 5)
				}
				logSession.Printf("NewSession-%s ID=%s token=%s\n",
					keepaliveID,
					newSess.ID,
					SabreTokenParse(sess.Sabre.Header.Security.BinarySecurityToken.Value),
				)
				//kill sess::Session by GC, already pulled off queue
				p.Sessions <- newSess
				continue
			}
			//reset expire, validated time, binary token(shouldn't change but set to Sabre return)
			sess.ExpireTime = time.Now().Add(time.Minute * time.Duration(RandomInt(p.Expire.Min, p.Expire.Max)))
			sess.TimeValidated = time.Now()
			//sess.Sabre.Header.Security.BinarySecurityToken = validateRS.Header.Security.BinarySecurityToken
			logSession.Printf(
				"UPDATED-%s-%s token=%s ID=%s AliveFor=%.2f(mins) Next ExpireIn=%.2f(mins)\n",
				keepaliveID,
				validateRS.Header.MessageHeader.Action,
				SabreTokenParse(sess.Sabre.Header.Security.BinarySecurityToken.Value),
				sess.ID,
				time.Since(sess.TimeStarted).Minutes(),
				time.Until(sess.ExpireTime).Minutes(),
			)
			//put session back on queue
			p.Sessions <- sess
		} else {
			//put session back on queue
			p.Sessions <- sess
			logSession.Printf("VALID-%s token=%s ID=%s ExpiresIn=%.2f(mins)\n",
				keepaliveID,
				SabreTokenParse(sess.Sabre.Header.Security.BinarySecurityToken.Value),
				sess.ID,
				time.Until(sess.ExpireTime).Minutes(),
			)
		}
	}
}

// generateKeepAliveID is an ID an easy to find and parse in the logs,
// attached to every new call on Keepalive; returns format 'kid:tXury|0220-17:12'
func generateKeepAliveID() string {
	randStr := randStringBytesMaskImprSrc(5)
	nowtime := time.Now().Format("0102-15:04")
	return "kid:" + randStr + "|" + nowtime
}

// Keepalive cycles through the pool periodically, running RangeKeepalive.
// This means they are guaranteed to be valid and the pool to be correct size.
// It accepts channel of type os.Signal which waits to shutdown the process.
func (p *SessionPool) Keepalive(doneChan chan os.Signal) {
	started := time.Now()
	keepAliveID := generateKeepAliveID()
	logSession.Println("Starting KEEPALIVE...", keepAliveID)

	for {
		select {
		case <-time.After(p.CycleEvery):
			p.NumberOfCycles++
			p.RangeKeepalive(keepAliveID)
			logSession.Printf("KEEPALIVE run(InMin=%.2f, InHour=%.2f)", time.Since(started).Minutes(), time.Since(started).Hours())
			p.logReport(keepAliveID + "-KeepAlive")
		case <-doneChan:
			logSession.Println("KEEPALIVE done, total lifetime:", time.Since(started))
			return
		}
	}
}

// Close down all sessions gracefull and valid on Sabre. It does this by iterating through all sessions in the pool and initializing a correct close session request to Sabre.
// If sessions are not properly closed on Sabre side they remain open and invalidate the workspace for up to an hour, which means you cannot open new sessions.
func (p *SessionPool) Close() {
	networkErrors := []error{}
	faultErrors := []error{}
	for createRS := range p.Sessions {
		p.Conf.SetBinSec(createRS.Sabre)
		closeRQ := BuildSessionCloseRequest(p.Conf)
		closeRS, err := CallSessionClose(p.ServiceURL, closeRQ)

		if err != nil {
			networkErrors = append(networkErrors, err)
		}

		fc := closeRS.Body.Fault.Code
		if fc != "" {
			st := closeRS.Body.Fault.Detail.StackTrace
			fs := closeRS.Body.Fault.String
			faultErrors = append(faultErrors, fmt.Errorf("%s-%s: %s", fs, fc, st))
		}
		p.AllowPoolSize--
		logSession.Printf("Status='%s' closed session with token='%s'", closeRS.Body.SessionCloseRS.Status, SabreTokenParse(closeRS.Header.Security.BinarySecurityToken.Value))

		if p.AllowPoolSize == 0 {
			close(p.Sessions)
		}
	}
	p.NetworkErrors = networkErrors
	p.FaultErrors = faultErrors

	logSession.Printf("Closing report... AllowPoolSize=%d, Busy=%d, Queuesize=%d, NetworkErrors=%v,  FaultErrors=%v", p.AllowPoolSize, (p.AllowPoolSize - len(p.Sessions)), len(p.Sessions), networkErrors, faultErrors)

	logSession.Println("Close SessionPool complete")
}
