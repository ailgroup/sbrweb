package srvc

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
	cycleEvery = time.Minute * 3
)

func TestSessionPoolNewCounters(t *testing.T) {
	poolSize := 3
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	if p.ServiceURL != serverCreateRQ.URL {
		t.Errorf("SessionPool.ServiceURL expect: %s, got: %s", serverCreateRQ.URL, p.ServiceURL)
	}
	if p.ConfigPoolSize != poolSize {
		t.Errorf("ConfigPoolSize should be same as user defined value. expect: %d, got: %d", poolSize, p.ConfigPoolSize)
	}
	if p.ConfigPoolSize == p.PoolSizeCounter {
		t.Errorf("ConfigPoolSize: %d should equal PoolSizeCounter: %d", p.ConfigPoolSize, p.PoolSizeCounter)
	}
	if len(p.Sessions) == poolSize {
		t.Errorf("Sessions should not equal user defined value on NewPool(). expect: %d, got: %d", 0, len(p.Sessions))
	}
	if (p.PoolSizeCounter - len(p.Sessions)) != 0 {
		t.Errorf("Open Sessions expect: %d, got: %d", 0, (p.PoolSizeCounter - len(p.Sessions)))
	}
}

func TestSessionPoolZeroPoolSize(t *testing.T) {
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, 0)
	err := p.Populate()
	if err == nil {
		t.Error("0 ConfigPoolSize should return error:", err)
	}
	if len(p.Sessions) != 0 {
		t.Error("0 ConfigPoolSize should have 0 Sessions")
	}
	if p.PoolSizeCounter != 0 {
		t.Error("0 ConfigPoolSize should have 0 PoolSizeCounter")
	}
	sess := p.Pick()
	if sess.ID != "" {
		t.Errorf("Pick session from closed queue should not have id; expect: %s, got: %s", "", sess.ID)
	}
}

func TestSessionPoolPopluateServerDown(t *testing.T) {
	poolSize := 3
	sampleSessionConf.ServiceURL = serverDown.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	err := p.Populate()
	if err != nil {
		t.Error("Bad Populate should not return error:", err)
	}
	if len(p.NetworkErrors) == 0 {
		t.Error("Bad Populate should have network errors:", p.NetworkErrors)
	}
	if len(p.Sessions) != poolSize {
		t.Errorf("Expect %d sessions when server down, got (len.Sessions)=%d", 0, len(p.Sessions))
	}
	if NumberOfBadSessions != poolSize {
		t.Errorf("NumberOfBadSessions expect: %d, got: %d", 0, NumberOfBadSessions)
	}
	if len(p.NetworkErrors) <= 0 {
		t.Errorf("Expect NetworkErrors, got: %v", p.NetworkErrors)
	}
	if p.PoolSizeCounter != poolSize {
		t.Error("PoolSizeCounter should == poolSize since we expect them to heal:", p.PoolSizeCounter)
	}
	if len(p.Sessions) != poolSize {
		t.Error("SessionPool size should == poolSize since we expect them to heal:", len(p.Sessions))
	}

	sess := p.Pick()
	if (p.PoolSizeCounter - len(p.Sessions)) != 1 {
		t.Errorf("Sessions should be missing 1. expect: %d, got: %d", 1, (p.PoolSizeCounter - len(p.Sessions)))
	}
	if sess.ID == "" {
		t.Errorf("Pick bad session from bad queue should still have id; expect: %s, got: %s", sess.ID, "")
	}
	if sess.OK {
		t.Error("Session should not be OK when server not available")
	}
	if sess.BinSecTokCached != "" {
		t.Errorf("BinSecTokCached should be: %s for bad session, instead: %s", "", sess.BinSecTokCached)
	}
}

func TestSessionPoolPopluateBadBody(t *testing.T) {
	poolSize := 2
	sampleSessionConf.ServiceURL = serverBadBody.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	err := p.Populate()
	if err != nil {
		t.Error("Bad Populate should not return error:", err)
	}
	if len(p.NetworkErrors) == 0 {
		t.Error("Bad Populate should have network errors:", p.NetworkErrors)
	}
	if p.PoolSizeCounter != poolSize {
		t.Error("PoolSizeCounter should == poolSize since we expect them to heal")
	}
}
func TestSessionPoolCloseOnDownServer(t *testing.T) {
	poolSize := 2
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	_ = p.Populate()
	//reroute to unavailable server
	p.ServiceURL = serverDown.URL
	p.Close()
	if len(p.NetworkErrors) <= 0 {
		t.Fail()
	}
}

func TestSessionPoolPopluateUnauth(t *testing.T) {
	poolSize := 2
	sampleSessionConf.ServiceURL = serverCreateRSUnauth.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	_ = p.Populate()
	if len(p.NetworkErrors) != 0 {
		t.Error("Network errors should be 0:", p.NetworkErrors)
	}
	if p.PoolSizeCounter != poolSize {
		t.Error("PoolSizeCounter should be more than zero even with SOAP Fault")
	}
	sess := p.Pick()
	if sess.OK {
		t.Error("Session should not be OK when server not available")
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
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	_ = p.Populate()
	if len(p.NetworkErrors) != 0 {
		t.Error("Network errors should be 0:", p.NetworkErrors)
	}
	if p.PoolSizeCounter != poolSize {
		t.Error("PoolSizeCounter should be more than zero even with SOAP Fault")
	}
	//reroute serivce to server with close invalid response...
	p.ServiceURL = serverCloseRSInvalid.URL
	p.Close()
	if len(p.FaultErrors) != poolSize {
		t.Error("FaultErrors should exist:", p.FaultErrors)
	}
	if len(p.NetworkErrors) != 0 {
		t.Error("NetworkErrors should not exist")
	}
	if p.PoolSizeCounter != 0 {
		t.Error("PoolSizeCounter should be more than zero even with SOAP Fault")
	}
	if len(p.FaultErrors) != poolSize {
		t.Error("SessionPool fault errors shoudl exist")
	}
}

func TestSessionPoolPopluatePickPut(t *testing.T) {
	poolSize := 5
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	_ = p.Populate()
	if poolSize != p.PoolSizeCounter {
		t.Errorf("given poolSize: %d should equal PoolSizeCounter on Populate: %d", poolSize, p.PoolSizeCounter)
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
		if p.PoolSizeCounter != poolSize {
			t.Errorf("AoowPoolSize after Pick() expect: %d, got: %d", poolSize, p.PoolSizeCounter)
		}
		if (p.PoolSizeCounter - len(p.Sessions)) != i {
			t.Errorf("Busy Sessions after Pick() expect: %d, got: %d", i, (p.PoolSizeCounter - len(p.Sessions)))
		}
	}
	for num, sess := range sessions {
		if len(p.Sessions) != num {
			t.Errorf("Open Sessions for Put(sess) expect: %d, got: %d", num, len(p.Sessions))
		}
		if p.PoolSizeCounter != poolSize {
			t.Errorf("PoolSizeCounter for Put(sess) expect: %d, got: %d", poolSize, p.PoolSizeCounter)
		}
		expectBusy := len(sessions) - num
		if (p.PoolSizeCounter - len(p.Sessions)) != expectBusy {
			t.Errorf("Busy Sessions for Put(sess) expect: %d, got: %d", expectBusy, (p.PoolSizeCounter - len(p.Sessions)))
		}
		p.Put(sess)
	}
}
func TestSessionPoolBlocking(t *testing.T) {
	poolSize := 5
	blockingSize := 3
	sampleSessionConf.ServiceURL = serverCreateRQ.URL
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	_ = p.Populate()

	sessions := []Session{}
	for i := 1; i <= poolSize; i++ {
		sessions = append(sessions, p.Pick())
		expectOpen := (poolSize - i)
		if len(p.Sessions) != expectOpen {
			t.Errorf("Open Sessions after Pick() expect: %d, got: %d", expectOpen, (p.PoolSizeCounter - len(p.Sessions)))
		}
		if p.PoolSizeCounter != poolSize {
			t.Errorf("PoolSizeCounter after Put(sess) expect: %d, got: %d", poolSize, p.PoolSizeCounter)
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
	if (p.PoolSizeCounter - len(p.Sessions)) != poolSize {
		t.Errorf("Busy Sessions after all Pick() and we have blocking requests expect: %d, got: %d", poolSize, (p.PoolSizeCounter - len(p.Sessions)))
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
	if (p.PoolSizeCounter - len(p.Sessions)) != blockingSize {
		t.Errorf("Busy Sessions after busy Put() back and blocking requests handled: %d, got: %d", blockingSize, (p.PoolSizeCounter - len(p.Sessions)))
	}
	if p.PoolSizeCounter != poolSize {
		t.Errorf("PoolSizeCounter after handling blocked requests expect: %d, got: %d", poolSize, p.PoolSizeCounter)
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
	if (p.PoolSizeCounter - len(p.Sessions)) != 0 {
		t.Errorf("Busy Sessions after all requests handled: %d, got: %d", 0, (p.PoolSizeCounter - len(p.Sessions)))
	}
	if p.PoolSizeCounter != poolSize {
		t.Errorf("PoolSizeCounter after all requests handled: %d, got: %d", poolSize, p.PoolSizeCounter)
	}
}

func TestSessionPoolSafeBlocking(t *testing.T) {
	poolSize := 3
	p := NewPool(sampleExpireScheme, sampleSessionConf, cycleEvery, poolSize)
	if p.ConfigPoolSize == p.PoolSizeCounter {
		t.Errorf("ConfigPoolSize: %d not equal PoolSizeCounter: %d", p.ConfigPoolSize, p.PoolSizeCounter)
	}

	p.ServiceURL = serverDown.URL
	err := p.Populate()
	if err != nil {
		t.Error("Populate pool with sessions from down server should not return error", err)
	}
	if len(p.NetworkErrors) == 0 {
		t.Error("Network errors should exist:", p.NetworkErrors)
	}
	if p.ConfigPoolSize != p.PoolSizeCounter {
		t.Errorf("ConfigPoolSize: %d not equal PoolSizeCounter: %d AFTER server goes down", p.ConfigPoolSize, p.PoolSizeCounter)
	}
	if NumberOfBadSessions != poolSize {
		t.Errorf("NumberOfBadSessions expect: %d, got: %d", 0, NumberOfBadSessions)
	}

	sess := p.Pick()
	if NumberOfBadSessions != poolSize {
		t.Errorf("NumberOfBadSessions expect: %d, got: %d", 0, NumberOfBadSessions)
	}
	// if no safe blocking this will never be reached
	// fact that you get here is a success for safe blocking!
	s := Session{}
	if sess == s {
		t.Error("Bad Session should not be empty:", sess)
	}
	if (p.PoolSizeCounter - len(p.Sessions)) != 1 {
		t.Errorf("Sessions should be missing 1. expect: %d, got: %d", 1, (p.PoolSizeCounter - len(p.Sessions)))
	}
	if sess.ID == "" {
		t.Errorf("Pick bad session from bad queue should still have id; expect: %s, got: %s", sess.ID, "")
	}
	if sess.OK {
		t.Error("Session should not be OK when server not available")
	}
	if sess.BinSecTokCached != "" {
		t.Errorf("BinSecTokCached should be: %s for bad session, instead: %s", "", sess.BinSecTokCached)
	}

	p.Put(sess)
	if (p.PoolSizeCounter - len(p.Sessions)) != 0 {
		t.Errorf("PoolSizeCounter and Sessions should be same. expect: %d, got: %d", 0, (p.PoolSizeCounter - len(p.Sessions)))
	}
	if NumberOfBadSessions != poolSize {
		t.Errorf("NumberOfBadSessions expect: %d, got: %d", 0, NumberOfBadSessions)
	}
}
