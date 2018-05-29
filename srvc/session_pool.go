package srvc

import (
	"fmt"
	"math/rand"
	"time"
	//"github.com/pkg/profile"
)

// Session holds sabre session data and other fields for handling in the SessionPool
type Session struct {
	ID            string
	Sabre         SessionCreateResponse
	TimeValidated time.Time
	TimeStarted   time.Time
	ExpireTime    time.Time
	FaultError    error
}

// SessionConfig holds info to create and manage session
type SessionConfig struct {
	from      string
	pcc       string
	convid    string
	mid       string
	timeStamp string
	username  string
	password  string
}

// ExpireScheme for when to expire sessions
type ExpireScheme struct {
	Max int
	Min int
}

// SessionPool holds sessions
type SessionPool struct {
	Cycles          int
	ConfigPoolSize  int
	AllowPoolSize   int
	Sessions        chan Session
	InitializedTime time.Time
	ServiceURL      string
	NetworkErrors   []error
	FaultErrors     []error
	Expire          ExpireScheme
	Conf            SessionConfig
}

// NewPool initializes a new pool with the given tasks and at the given
// concurrency.
func NewPool(expire ExpireScheme, size int, serviceURL, from, pcc, convid, mid, timeStamp, username, password string) *SessionPool {
	return &SessionPool{
		ConfigPoolSize:  size,
		Sessions:        make(chan Session, size), //buffered channel blocks!
		InitializedTime: time.Now(),
		ServiceURL:      serviceURL,
		Expire:          expire,
		Conf: SessionConfig{
			from:      from,
			pcc:       pcc,
			convid:    convid,
			mid:       mid,
			timeStamp: timeStamp,
			username:  username,
			password:  password,
		},
	}
}

// GenerateSessionID returns 'xca123'
func GenerateSessionID() string {
	randStr := randStringBytesMaskImprSrc(3)
	nowtime := time.Now().Format(".999")
	return randStr + nowtime
}

//RandomInt select random integer within range min/max
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
	//use key select position value out of range [min...max]
	return a[key]
}

func (p *SessionPool) newSession() (Session, error) {
	createRQ := BuildSessionCreateRequest(p.Conf.from, p.Conf.pcc, p.Conf.convid, p.Conf.mid, p.Conf.timeStamp, p.Conf.username, p.Conf.password)
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
		//ExpireTime:    now.Add(time.Minute * time.Duration(RandomInt(3, 14))),
		ExpireTime: now.Add(time.Minute * time.Duration(1)),
		FaultError: faultErr,
	}
	logSession.Printf(
		"Status='%s' created session ID=%s with Expirey='%s' for token=%s", sess.Sabre.Body.SessionCreateRS.Status,
		sess.ID,
		sess.ExpireTime,
		SabreTokenParse(sess.Sabre.Header.Security.BinarySecurityToken.Value),
	)
	return sess, nil
}

// Populate writes to buffered channel the max pool size and returns the current pool, boolean value if current pool size matches max pool size, along with a slice of any overall network errors. Note that if the SOAP service fails that will be located in the SOAP Fault field, not as a general network error. We want to fill buffer/queue even if we get a fault from Sabre: there may scenarios where this is legit (sabre flushes all sessions weekly, sabre sessions service goes down, or we are over session limit.... Under these conditions we don't want to block the queue or reQUEUESIZEpeatedly attempt to fill it. Instead, accept a bad session and let the keepalive cleanup faulted sessions later.
func (p *SessionPool) Populate() error {
	var err error
	ok := (p.AllowPoolSize == len(p.Sessions))
	for i := 0; i < p.ConfigPoolSize; i++ {
		sess, err := p.newSession()
		if err != nil {
			logSession.Printf("Network ERROR for attempt=%d, continue", i)
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

// Pick session from buffered queue
func (p *SessionPool) Pick() Session {
	sess := <-p.Sessions
	p.logReport("Pick-" + sess.ID)
	//logSession.Printf("PICK ID=%s, Alow=%d, Busy=%d, Queue=%d", sess.ID, p.AllowPoolSize, (p.AllowPoolSize - len(p.Sessions)), len(p.Sessions))
	return sess
}

// Put session back onto the buffered queue
func (p *SessionPool) Put(sess Session) {
	p.Sessions <- sess
	p.logReport("Put-" + sess.ID)
	//logSession.Printf("PUT ID=%s, Alow=%d, Busy=%d, Queue=%d", sess.ID, p.AllowPoolSize, (p.AllowPoolSize - len(p.Sessions)), len(p.Sessions))
}

// logReport helper to log info about session pool
func (p *SessionPool) logReport(ctx string) {
	open := len(p.Sessions)
	allow := p.AllowPoolSize
	busy := (allow - open)
	logSession.Printf("[%s] ALLOW=%d, OPEN=%d, BUSY=%d, CYCLES=%d, OK=%v", ctx, allow, open, busy, p.Cycles, (allow == (open + busy)))
}

//RangeKeepalive sessions to validate.
// Range over sessions, if expire is over current time, valdiate, reset expire
// and place back into the queue, otherwise place it back on the queue and move on.
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
				SabreTimeFormat(),
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
					newSess.ExpireTime = time.Now().Add(time.Second * 1)
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

func generateKeepAliveID() string {
	randStr := randStringBytesMaskImprSrc(5)
	nowtime := time.Now().Format("0102-15:04")
	return "kid:" + randStr + "|" + nowtime
}

//Keepalive sessions in the session pool validated
func Keepalive(p *SessionPool, repeatEvery, endAfter time.Duration) {
	//defer profile.Start(profile.MemProfile).Stop()
	doneChan := time.NewTimer(endAfter).C
	started := time.Now()
	keepAliveID := generateKeepAliveID()
	logSession.Println("Starting KEEPALIVE...", keepAliveID)

	//keep validating until exhaust endAfter
	for {
		select {
		case <-time.After(repeatEvery):
			p.Cycles++
			p.RangeKeepalive(keepAliveID)
			logSession.Printf("KEEPALIVE run(InMin=%.2f, InHour=%.2f)", time.Since(started).Minutes(), time.Since(started).Hours())
			p.logReport(keepAliveID + "-KeepAlive")
		case <-doneChan:
			logSession.Println("KEEPALIVE killed, total time:", time.Since(started))
			return
		}
	}
}

// Close down all sessions in a nice manner
func (p *SessionPool) Close() {
	networkErrors := []error{}
	faultErrors := []error{}
	for createRS := range p.Sessions {
		closeRQ := BuildSessionCloseRequest(
			createRS.Sabre.Header.MessageHeader.To.PartyID.Value,
			createRS.Sabre.Header.MessageHeader.CPAID,
			createRS.Sabre.Header.Security.BinarySecurityToken.Value,
			createRS.Sabre.Header.MessageHeader.ConversationID,
			createRS.Sabre.Header.MessageHeader.MessageData.RefToMessageID,
			SabreTimeFormat(),
		)
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

	logSession.Printf("Close report... AllowPoolSize=%d, Busy=%d, Queuesize=%d, NetworkErrors=%v,  FaultErrors=%v", p.AllowPoolSize, (p.AllowPoolSize - len(p.Sessions)), len(p.Sessions), networkErrors, faultErrors)

	logSession.Println("Close SessionPool complete")
}

/*
func HotelAvail(p *SessionPool) HotelData {
	sess := p.Pick()
	defer func(s *SessionPool) {
		p.Put(sess)
	}(sess)
	HotelData := CallAvail(sess.DATA, other.DATA)
	return HotelData
}



 //TODO figure this one out
//SessionTable print data about all sessions
func (p *SessionPool) SessionTable() map[int][]string {
	var sessionData map[int][]string
	sessionData = make(map[int][]string)
	for s := range p.Sessions {
		sessionData[s.ID] = []string{fmt.Sprintf(
			"Started: %s, Validated: %s, SabreToken: %s",
			s.TimeStarted,
			s.TimeValidated,
			s.Sabre.Header.Security.BinarySecurityToken.Value),
			s.Sabre.Body.SessionCreateRS.Status,
		}
	}
	return sessionData
}

// Report on state of session pool
func (p *SessionPool) Report() string {
	return fmt.Sprintf(
		"=> PoolReport => POOLSIZE: %d, BUFFERSIZE: %d, BUSY: %d, NOTBUSY: %d, CYCLES: %d",
		p.CurrentPoolSize,
		len(p.Sessions),
		p.Busy,
		p.NotBusy,
		p.Cycles,
	)
}

//ReportRun to see info about the session pool and validators
func ReportRun(p *SessionPool, repeatEvery, endAfter, runReport time.Duration) {
	doneChan := time.NewTimer(endAfter).C
	started := time.Now()
	log.Println("Starting report runner...")
	//keep reporting until exhaust endAfter
	for {
		select {
		case <-time.After(repeatEvery):
			log.Println(p.Report())
		case <-doneChan:
			log.Println("ReportRun killed, total time:", time.Since(started))
			return
		}
	}
}

// Run works
func (p *SessionPool) Run() {
	p.pickCounter()
	sess := <-p.Sessions
	log.Printf("PICK ID=%d, Busy=%d, NotBusy=%d, Queue=%d", sess.ID, p.Busy, p.NotBusy, len(p.Sessions))
	defer func() {
		p.Sessions <- sess
	}()
}
*/
