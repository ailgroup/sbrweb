package sbrweb

import (
	"testing"
	"time"
)

/*
	TODO mock SessionCreateResponse objects and stuff them into the SessionPool
	Can test according to that without needing mock network requests
*/

var (
	sampleExpireScheme = ExpireScheme{
		Min: 3,
		Max: 14,
	}
)

func TestSessionPoolEmpty(t *testing.T) {
	poolSize := 3
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	if p.ServiceURL != serverCreateRQ.URL {
		t.Errorf("SessionPool.ServiceURL expect: %s, got: %s", serverCreateRQ.URL, p.ServiceURL)
	}
	if p.ConfigPoolSize != poolSize {
		t.Errorf("ConfigPoolSize expect: %d, got: %d", poolSize, p.ConfigPoolSize)
	}
	if p.ConfigPoolSize == p.AllowPoolSize {
		t.Errorf("ConfigPoolSize: %d should not equal AllowPoolSize: %d", p.ConfigPoolSize, p.AllowPoolSize)
	}
	if len(p.Sessions) != 0 {
		t.Errorf("Sessions expect: %d, got: %d", 0, len(p.Sessions))
	}
	if (p.AllowPoolSize - len(p.Sessions)) != 0 {
		t.Errorf("Open Sessions expect: %d, got: %d", 0, (p.AllowPoolSize - len(p.Sessions)))
	}
}

// NOTE: this helps test case where we could not get any valid sessions, so we close the buffered channel to prevent indefinite blocking.
func TestSessionPoolPopluateServerDown(t *testing.T) {
	poolSize := 3
	p := NewPool(sampleExpireScheme, poolSize, serverDown.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()

	if len(p.Sessions) != 0 {
		t.Errorf("Expect %d sessions when server down, got (len.Sessions)=%d", 0, len(p.Sessions))
	}
	if len(p.NetworkErrors) <= 0 {
		t.Errorf("Expect NetworkErrors, got: %v", p.NetworkErrors)
	}
	if p.AllowPoolSize != 0 {
		t.Error("AllowPoolSize should be zero with all failed sessions", p.AllowPoolSize)
	}
	if len(p.Sessions) != 0 {
		t.Error("QUEUE size should be zero since it has been closed due to block safety", len(p.Sessions))
	}
	sess := p.Pick()
	if sess.ID != "" {
		t.Errorf("Pick session from closed queue should have id==0; expect: %s, got: %s", "", sess.ID)
	}
	if (p.AllowPoolSize - len(p.Sessions)) != 0 {
		t.Errorf("Pick session from closed queue should NotBusy == -1; expect: %d, got: %d", 0, (p.AllowPoolSize - len(p.Sessions)))
	}
}

func TestSessionPoolPopluateBadBody(t *testing.T) {
	poolSize := 2
	p := NewPool(sampleExpireScheme, poolSize, serverBadBody.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()
	if len(p.NetworkErrors) <= 0 {
		t.Fail()
	}
	if p.AllowPoolSize != 0 {
		t.Error("AllowPoolSize should be zero with all failed sessions")
	}
}
func TestSessionPoolCloseOnDownServer(t *testing.T) {
	poolSize := 2
	//must be a good server or no sessions and buffer blocks...
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()

	//reroute to unavailable server
	p.ServiceURL = serverDown.URL
	//p.ServiceURL = serverBadBody.URL

	p.Close()
	if len(p.NetworkErrors) <= 0 {
		t.Fail()
	}
}

func TestSessionPoolPopluateUnauth(t *testing.T) {
	poolSize := 2
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRSUnauth.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()
	if len(p.NetworkErrors) != 0 {
		t.Fail()
	}
	if p.AllowPoolSize != poolSize {
		t.Error("AllowPoolSize should be more than zero even with SOAP Fault")
	}
	sess := p.Pick()
	if sess.Sabre.Body.Fault.Detail.StackTrace != sampleSessionNoAuthStackTrace {
		t.Errorf("Session Fault Invalid token expect: %s, got: %s", sampleSessionNoAuthStackTrace, sess.Sabre.Body.Fault.String)
	}
	if sess.FaultError.Error() != sampleSessionPoolMsgNoAuth {
		t.Errorf("Session Soap Error should be nasty string. expect: %s, got: %s", sampleSessionPoolMsgNoAuth, sess.FaultError)
	}
	if len(p.FaultErrors) != poolSize {
		t.Error("SessionPool fault errors shoudl exist")
	}
}

func TestSessionPoolCloseInvalidToken(t *testing.T) {
	poolSize := 2
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()
	if len(p.NetworkErrors) != 0 {
		t.Fail()
	}
	if p.AllowPoolSize != poolSize {
		t.Error("AllowPoolSize should be more than zero even with SOAP Fault")
	}

	//reroute serivce to server with close invalid response...
	p.ServiceURL = serverCloseRSInvalid.URL
	p.Close()

	if len(p.NetworkErrors) != 0 {
		t.Error("Network errors should not exist")
	}
	if p.AllowPoolSize != 0 {
		t.Error("AllowPoolSize should be more than zero even with SOAP Fault")
	}
	if len(p.FaultErrors) != poolSize {
		t.Error("SessionPool fault errors shoudl exist")
	}
}

func TestSessionPoolPopluatePickPut(t *testing.T) {
	poolSize := 5
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()
	if poolSize != p.AllowPoolSize {
		t.Errorf("given poolSize: %d should equal AllowPoolSize on Populate: %d", poolSize, p.AllowPoolSize)
	}
	if len(p.Sessions) != poolSize {
		t.Errorf("Sessions expect: %d, got: %d", poolSize, len(p.Sessions))
	}

	sessions := []Session{}
	for i := 1; i <= poolSize; i++ {
		sessions = append(sessions, p.Pick())
		expectOpen := (poolSize - i)
		if len(p.Sessions) != expectOpen {
			t.Errorf("Open Sessions after Pick() expect: %d, got: %d", expectOpen, len(p.Sessions))
		}
		if p.AllowPoolSize != poolSize {
			t.Errorf("AoowPoolSize after Pick() expect: %d, got: %d", poolSize, p.AllowPoolSize)
		}
		if (p.AllowPoolSize - len(p.Sessions)) != i {
			t.Errorf("Busy Sessions after Pick() expect: %d, got: %d", i, (p.AllowPoolSize - len(p.Sessions)))
		}
	}
	for num, sess := range sessions {
		if len(p.Sessions) != num {
			t.Errorf("Open Sessions for Put(sess) expect: %d, got: %d", num, len(p.Sessions))
		}
		if p.AllowPoolSize != poolSize {
			t.Errorf("AllowPoolSize for Put(sess) expect: %d, got: %d", poolSize, p.AllowPoolSize)
		}
		expectBusy := len(sessions) - num
		if (p.AllowPoolSize - len(p.Sessions)) != expectBusy {
			t.Errorf("Busy Sessions for Put(sess) expect: %d, got: %d", expectBusy, (p.AllowPoolSize - len(p.Sessions)))
		}
		p.Put(sess)
	}
}
func TestSessionPoolBlocking(t *testing.T) {
	poolSize := 5
	blockingSize := 3
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	p.Populate()

	sessions := []Session{}
	for i := 1; i <= poolSize; i++ {
		sessions = append(sessions, p.Pick())
		expectOpen := (poolSize - i)
		if len(p.Sessions) != expectOpen {
			t.Errorf("Open Sessions after Pick() expect: %d, got: %d", expectOpen, (p.AllowPoolSize - len(p.Sessions)))
		}
		if p.AllowPoolSize != poolSize {
			t.Errorf("AllowPoolSize after Put(sess) expect: %d, got: %d", poolSize, p.AllowPoolSize)
		}
	}

	// setup blocking requests to the session pool
	bg := []Session{}
	go func(b []Session, pool *SessionPool) {
		for i := 1; i <= blockingSize; i++ {
			bg = append(bg, pool.Pick())
		}
	}(bg, p)

	//verify session pool
	if len(p.Sessions) != 0 {
		t.Errorf("Open Sessions after all Pick() and we have blocking requests expect: %d, got: %d", 0, len(p.Sessions))
	}
	if (p.AllowPoolSize - len(p.Sessions)) != poolSize {
		t.Errorf("Busy Sessions after all Pick() and we have blocking requests expect: %d, got: %d", poolSize, (p.AllowPoolSize - len(p.Sessions)))
	}

	//put back busy sessions with requests waiting
	// give a few milliseconds for blocking requests
	// to pick from queue
	for _, sess := range sessions {
		p.Put(sess)
		time.Sleep(7 * time.Millisecond)
	}
	//poolsize - blocking size
	// blocking requests have been handled, which means they are now busy
	// leaving remaining sessions open...
	if len(p.Sessions) != 2 {
		t.Errorf("Open Sessions after busy Put() back and blocking requests handled: %d, got: %d", 2, len(p.Sessions))
	}
	// expect all to still be busy
	if (p.AllowPoolSize - len(p.Sessions)) != blockingSize {
		t.Errorf("Busy Sessions after busy Put() back and blocking requests handled: %d, got: %d", blockingSize, (p.AllowPoolSize - len(p.Sessions)))
	}
	if p.AllowPoolSize != poolSize {
		t.Errorf("AllowPoolSize after handling blocked requests expect: %d, got: %d", poolSize, p.AllowPoolSize)
	}
	if len(bg) != blockingSize {
		t.Errorf("Waiting requests have now been handled, should be same number as blockingSize expect: %d, got: %d", blockingSize, len(bg))
	}

	//put back busy backgrounded sessions
	for _, bgsess := range bg {
		p.Put(bgsess)
	}
	if len(p.Sessions) != poolSize {
		t.Errorf("Open Sessions after all requests handled: %d, got: %d", poolSize, len(p.Sessions))
	}
	// expect all to still be busy
	if (p.AllowPoolSize - len(p.Sessions)) != 0 {
		t.Errorf("Busy Sessions after all requests handled: %d, got: %d", 0, (p.AllowPoolSize - len(p.Sessions)))
	}
	if p.AllowPoolSize != poolSize {
		t.Errorf("AllowPoolSize after all requests handled: %d, got: %d", poolSize, p.AllowPoolSize)
	}
}

func TestSessionPoolSafeBlocking(t *testing.T) {
	poolSize := 3
	p := NewPool(sampleExpireScheme, poolSize, serverCreateRQ.URL, samplefrom, samplepcc, sampleconvid, samplemid, sampletime, sampleusername, samplepassword)
	if p.ConfigPoolSize == p.AllowPoolSize {
		t.Errorf("ConfigPoolSize: %d not equal AllowPoolSize: %d", p.ConfigPoolSize, p.AllowPoolSize)
	}

	p.ServiceURL = serverDown.URL
	err := p.Populate()
	if p.ConfigPoolSize == p.AllowPoolSize {
		t.Errorf("ConfigPoolSize: %d not equal AllowPoolSize AFTER server goes down: %d", p.ConfigPoolSize, p.AllowPoolSize)
	}
	if err == nil {
		t.Error("Should have returned error on forever-blocking session pool", err)
	}
	sess := p.Pick()
	//if no safe blocking this will never be reached
	// fact that you get here is a success for safe blocking!
	s := Session{}
	if sess != s {
		t.Error("Should get to this code and return a empty session", sess)
	}
	// p.Put(sess) --> this will explode: SEND ON CLOSED CHANNEL
	//fmt.Printf("%+v\n", p.Sessions)
}
