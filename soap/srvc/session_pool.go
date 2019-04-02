package srvc

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	//these counts separate from session pool to prevent data races
	NumberSessionPoolCycles = 0
	NumberOfBadSessions     = 0
)

// Session holds sabre session data and other fields for handling in the SessionPool
type Session struct {
	ID               string
	FaultError       error
	OK               bool
	PromisePutAfter  time.Duration  //time elapsed put back in to poo
	PromisePut       chan bool      //put back in pool
	PromiseSigListen chan os.Signal //end promise on signal
	TimeValidated    time.Time
	TimeStarted      time.Time
	ExpireTime       time.Time
	BinSecTokCached  string
}

// ExpireScheme for when to expire sessions
type ExpireScheme struct {
	Max int
	Min int
}

// SessionPool container for pool of sessions with specs on size, cycles, counters, errors, timers, configuration, etc...
type SessionPool struct {
	ConfigPoolSize  int
	PoolSizeCounter int
	refreshMod      int
	Expire          ExpireScheme
	CycleEvery      time.Duration
	ServiceURL      string
	NetworkErrors   []error
	FaultErrors     []error
	InitializedTime time.Time
	Sessions        chan Session
	ShutDown        chan os.Signal
	Signals         []os.Signal
	Conf            *SessionConf
}

func findMod(total int) int {
	if total <= 3 {
		return 3
	}

	perc := (20 * total) / 100
	if perc < 4 {
		return 5
	}
	return perc
}

// NewPool initializes a new session pool of given size, for ServiceURL, with a buffered channel of
// type Session, and with specification for cycled keepalive checks, expiration and
// timer for an initialized time. Do not initialize sessions until ready to populate, this keeps it nil until the time we actually need it.
func NewPool(expire ExpireScheme, cred *SessionConf, cycle time.Duration, size int, sig ...os.Signal) *SessionPool {
	return &SessionPool{
		ConfigPoolSize:  size,
		refreshMod:      findMod(size),
		InitializedTime: time.Now(),
		ServiceURL:      cred.ServiceURL,
		CycleEvery:      cycle,
		Expire:          expire,
		Conf:            cred,
		Signals:         sig,
	}
}

// GenerateSessionID for small easy to find ids in logs; returns format 'PGC.346'
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
// with metadata, logs it, and returns that Session for placement into the pool.
func (p *SessionPool) newSession() (Session, error) {
	var err error
	var ok bool = true
	createRQ := BuildSessionCreateRequest(p.Conf)
	createRS, err := CallSessionCreate(p.ServiceURL, createRQ)
	if err != nil {
		// create is special, we still want to put crappy sessions into the buffer because RangeKeepAlive will eventually heal them
		ok = false
		p.NetworkErrors = append(p.NetworkErrors, err)
	}
	now := time.Now()
	var faultErr error
	fc := createRS.Body.Fault.Code
	if fc != "" {
		ok = false
		st := createRS.Body.Fault.Detail.StackTrace
		fs := createRS.Body.Fault.String
		faultErr = fmt.Errorf("%s-%s: %s", fs, fc, st)
		p.FaultErrors = append(p.FaultErrors, faultErr)
	}

	// still want to fill buffer even if we get a fault from Sabre
	// there may sceanrios where this is legit (sabre flushes all sessions weekly,
	// or we are over threshold... don't want to block the queue or try to keep
	// filling it with bad sessions. Instead, accept bad sessions, and let the
	// cleanup will clear out faulted sessions later).
	sess := Session{
		ID:               GenerateSessionID(),
		BinSecTokCached:  createRS.Header.Security.BinarySecurityToken.Value,
		TimeStarted:      now,
		TimeValidated:    now,
		ExpireTime:       now.Add(time.Minute * time.Duration(RandomInt(p.Expire.Min, p.Expire.Max))),
		FaultError:       faultErr,
		OK:               ok,
		PromisePutAfter:  time.Duration(5 * time.Minute),
		PromisePut:       make(chan bool),
		PromiseSigListen: make(chan os.Signal, 1),
	}
	signal.Notify(sess.PromiseSigListen, p.Signals...)
	var status string
	if createRS.Body.SessionCreateRS.Status == "" {
		status = "NO CREATE"
		sess.OK = false
	} else {
		status = createRS.Body.SessionCreateRS.Status
	}
	logSession.Printf(
		"ID-%s OK=%v Create Status='%s' Expirey='%s' for token=%s",
		sess.ID,
		sess.OK,
		status,
		sess.ExpireTime,
		SabreTokenParse(sess.BinSecTokCached),
	)
	countBadSessions(sess.OK, sess.ID, p.ConfigPoolSize)
	return sess, err
}

func (p *SessionPool) refreshSession(sess Session) {
	logSession.Printf("RefreshSession ID-%s ", sess.ID)
	//if session was created while network was down its not going to have this and won't exist on sabre side... no use closing what does not exist, just try re-creating
	if sess.BinSecTokCached != "" {
		logSession.Printf("BinSecToken valid for close ID-%s ", sess.ID)
		closeRQ := BuildSessionCloseRequest(p.Conf, sess.BinSecTokCached)
		_, err := CallSessionClose(p.ServiceURL, closeRQ)
		if err != nil {
			fmt.Println(err)
		}
	}
	s, _ := p.newSession()
	p.Sessions <- s
}

func countBadSessions(sessOK bool, id string, configuredPoolSize int) {
	if !sessOK {
		if NumberOfBadSessions >= configuredPoolSize {
			//don't count any higher
			return
		}
		NumberOfBadSessions++
	} else if NumberOfBadSessions >= 1 { //if count 1 likely THIS is last/only bad session
		NumberOfBadSessions--
	}
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
	var ok bool
	p.Sessions = make(chan Session, p.ConfigPoolSize) //buffered channel blocks!
	for i := 0; i < p.ConfigPoolSize; i++ {
		sess, err := p.newSession()
		if err != nil {
			logSession.Printf("ERROR %v for attempt=%d, adding bad session, KeepAlive will heal it...", err, i)
		}
		p.Sessions <- sess
		p.PoolSizeCounter++
	}
	ok = ((p.PoolSizeCounter == len(p.Sessions)) && (p.ConfigPoolSize == len(p.Sessions)) && (len(p.NetworkErrors) != len(p.Sessions)))
	//close blocking channel, message, return error
	if p.ConfigPoolSize == 0 {
		err = fmt.Errorf("You have not allowed any sessions to be created, check PoolSizeCounter on SessionPool; closing SessionPool for now.")
		// p.NetworkErrors = append(p.NetworkErrors, err)
		//this closes it so the app can be shutdown(don't want to leave an empty buffered channel open because it will forever block). It does not ever allow the pool to be populated again... Since ConfigPoolSize is user defined this may be the best way
		close(p.Sessions)
		//Is it really OK? it's not blocking and that is good, but ...
		ok = false
	}
	logSession.Printf("Create PoolSizeCounter=%d, Create OK=%v. NetworkErrors=%v, FaultErrors=%v", p.PoolSizeCounter, ok, p.NetworkErrors, p.FaultErrors)
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
func (p SessionPool) logReport(ctx string) {
	configured := p.ConfigPoolSize
	count := p.PoolSizeCounter
	open := len(p.Sessions) - NumberOfBadSessions
	notOpen := (count - open)
	logSession.Printf("[%s] CONFIGURED=%d, COUNTED=%d, BAD=%d, OPEN=%d, NOT_OPEN=%d, CYCLES=%d, STABLE=%v", ctx, configured, count, NumberOfBadSessions, open, notOpen, NumberSessionPoolCycles, (count == (open + notOpen)))
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
		// keeping the pool "fresh": semi-randomly pick session, close, create new, put in pool
		if counter%p.refreshMod == 0 {
			if time.Since(sess.TimeValidated)%3 == 0 {
				logSession.Printf("ID-%s select for refresh...\n", sess.ID)
				p.refreshSession(sess)
				continue
			}
		}
		//time to expire and/or try to recover from bad state
		if time.Now().After(sess.ExpireTime) || !sess.OK {
			validateRQ := BuildSessionValidateRequest(p.Conf, sess.BinSecTokCached)
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
					newSess.ExpireTime = time.Now().Add(time.Second * 30)
				}
				logSession.Printf("ID-%s OK=%v NewSession-%s token=%s\n",
					newSess.ID,
					newSess.OK,
					keepaliveID,
					SabreTokenParse(newSess.BinSecTokCached),
				)
				//kill sess::Session  already pulled off queue, GC will pick it up...
				countBadSessions(newSess.OK, newSess.ID, p.ConfigPoolSize)
				p.Sessions <- newSess
				continue
			}
			//reset expire, validated time, binary token (these shouldn't change but update anyway)
			sess.ExpireTime = time.Now().Add(time.Minute * time.Duration(RandomInt(p.Expire.Min, p.Expire.Max)))
			sess.TimeValidated = time.Now()
			sess.BinSecTokCached = validateRS.Header.Security.BinarySecurityToken.Value
			logSession.Printf(
				"ID-%s UPDATED-%s-%s token=%s AliveFor=%.2f(mins) Next ExpireIn=%.2f(mins)\n",
				sess.ID,
				keepaliveID,
				validateRS.Header.MessageHeader.Action,
				SabreTokenParse(sess.BinSecTokCached),
				time.Since(sess.TimeStarted).Minutes(),
				time.Until(sess.ExpireTime).Minutes(),
			)
			//put session back on queue
			countBadSessions(sess.OK, sess.ID, p.ConfigPoolSize)
			p.Sessions <- sess
		} else {
			//put session back on queue
			p.Sessions <- sess
			logSession.Printf("ID-%s OK=%v VALIDATE-%s token=%s ExpiresIn=%.2f(mins)\n",
				sess.ID,
				sess.OK,
				keepaliveID,
				SabreTokenParse(sess.BinSecTokCached),
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

/*
Daemonize initializes and populates new session pool.
Accepts a waitgroup and variadic args signal handler for
graceful shutdown sessions in a valid Sabre transaction.
	Example:
		func up(pool *srvc.SessionPool) {
			var wg sync.WaitGroup
			wg.Add(1)
			go pool.Deamonize(&wg, os.Interrupt, syscall.SIGTERM)
			wg.Wait()
		}
		s := &http.Server{
			Addr:    port,
			Handler: router,
		}
		// pool uses Signal to capture sigint/sigterm to shutdown sabre sessions
		// because of this we need to shutdown server in background
		go func() {
			up(pool)
			_ = s.Shutdown(context.Background())
		}()
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			panic(fmt.Errorf("FATAL HTTP ERROR: %s", err))
		}
*/
func (p *SessionPool) Daemonize(wg *sync.WaitGroup) {
	//initialize the shutdown channel and notify for signals we care about
	p.ShutDown = make(chan os.Signal, 1)
	signal.Notify(p.ShutDown, p.Signals...)
	//	signal.Notify(p.ShutDown, os.Interrupt, syscall.SIGTERM)

	err := p.Populate()
	if err != nil {
		// let this play out... session pool should eventually self-heal
		fmt.Printf("Error popluating session pool %v\n", err)
	}
	//begin keepalive
	p.Keepalive()
	for {
		<-p.ShutDown
		wg.Done()
		return
	}
}

/*
Promise allows us to puts a session into a number of channels that guarantee it will be closed. This is used in cases where its better to keep making subsequent requests to Sabre using the same session but accross different handlers, api endpoints, or html pages. There are three promise struct fields that allow us manage the Promise: PromisePut is boolean and allows the user to force put the session back into the pool, effectively ending the promise; the other two are automated with one acting as a timeout (e.g., user leaves browser window open, no more requests coming in on the api), while the other listens for shutdown signals on the pool and makes sure the channel is put back into the pool before Close() on the pool is called.; the last two are run in a goroutine
NOTE: Promise is an advanced useage feature created specficially to run HOT* (hotel display), which must use the exact session meta data AND must guarantee no other session uses that meta-data. So we essentially lock out the channel associated with that meta-data for subsequent use. The danger here is that if you do not manage the Promise well you can exhaust you session pool, blocking all other requests until the Promise times out; do not do this.
These are set on new session creation
	sess.PromisePutAfter = time.Duration(5 * time.Minute)
	sess.PromisePut = make(chan bool)
	sess.PromiseSigListen = chan os.Signal
*/
func (sess *Session) Promise(p *SessionPool) {
	msg := fmt.Sprintf(
		"Promise-%s in (%.2f minutes) at '%s'",
		sess.ID,
		sess.PromisePutAfter.Minutes(),
		time.Now().Add(sess.PromisePutAfter),
	)
	p.logReport(msg)

	select {
	case <-sess.PromisePut:
		p.Sessions <- *sess
		p.logReport("PromisePut-" + sess.ID)
	//shutdown signal means session goes back for Close()
	case <-sess.PromiseSigListen:
		p.Sessions <- *sess
		p.logReport("PromisePutSignal-" + sess.ID)
		return
	//time elapsed put back in to pool
	case <-time.After(sess.PromisePutAfter):
		p.Sessions <- *sess
		p.logReport("PromisePutAfter-" + sess.ID)
	}
}

// Keepalive cycles through the pool periodically, running RangeKeepalive.
// This means they are guaranteed to be valid and the pool to be correct size.
// Listens for p.Shutdown, which is an os.Signal, on match will Close SessionPool and close(p.Shutdown) channel for program end.
func (p *SessionPool) Keepalive() {
	started := time.Now()
	keepAliveID := generateKeepAliveID()
	logSession.Printf("Starting KEEPALIVE...%v refresh modulo: %d for total: %d", keepAliveID, p.refreshMod, len(p.Sessions))
	p.logReport(keepAliveID + "-KeepAlive")
	for {
		select {
		case <-time.After(p.CycleEvery):
			NumberSessionPoolCycles++
			p.RangeKeepalive(keepAliveID)
			logSession.Printf("KEEPALIVE run(InMin=%.2f, InHour=%.2f)", time.Since(started).Minutes(), time.Since(started).Hours())
			p.logReport(keepAliveID + "-KeepAlive")
		case <-p.ShutDown:
			logSession.Println("KEEPALIVE shutdown, total lifetime:", time.Since(started))
			fmt.Println("\n-> Deamonize -> Close")
			p.Close() //close sabre sessions
			fmt.Println("Close -> close(Shutdown)")
			close(p.ShutDown) //shutdown session pool
			return
		}
	}
}

// Close down all sessions gracefull and valid on Sabre.
func (p *SessionPool) Close() {
	p.NetworkErrors, p.FaultErrors = p.loopOverPool([]error{}, []error{})

	logSession.Printf("Closing report... PoolSizeCounter=%d, Busy=%d, Queuesize=%d, NetworkErrors=%v,  FaultErrors=%v", p.PoolSizeCounter, (p.PoolSizeCounter - len(p.Sessions)), len(p.Sessions), p.NetworkErrors, p.FaultErrors)

	logSession.Println("Close SessionPool complete")
}

// TODO refactor this so its easy to just close one session so we can recreate a new one...
//loopHole iterates through all sessions in the pool and initializing a correct close session request to Sabre. If sessions are not properly closed on Sabre side they remain open and invalidate the workspace for up to an hour, which means you cannot open new sessions.
func (p *SessionPool) loopOverPool(networkErrors, faultErrors []error) ([]error, []error) {
	//only close down if we have sessions
	if len(p.Sessions) > 0 {
		for sessChan := range p.Sessions {
			//jsut make sure noting is holding on to a promise
			closeRQ := BuildSessionCloseRequest(p.Conf, sessChan.BinSecTokCached)
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
			p.PoolSizeCounter--
			logSession.Printf("ID-%s Close Status='%s' for token='%s'", sessChan.ID, closeRS.Body.SessionCloseRS.Status, SabreTokenParse(closeRS.Header.Security.BinarySecurityToken.Value))

			//only after we close the actual number of sessions allocated
			if p.PoolSizeCounter == 0 {
				close(p.Sessions)
			}
		}
	}
	return networkErrors, faultErrors
}
