package itin

import (
	"encoding/xml"
	"testing"

	"github.com/ailgroup/sbrweb/sbrerr"
)

func TestEndTransactionXML(t *testing.T) {
	et := BuildEndTransactionRequest(sampleConf, samplebinsectoken)
	_, err := xml.Marshal(et)
	if err != nil {
		t.Error("Error marshal build end transaction", err)
	}
	// if string(b) != string(sampleEndTReq) {
	// 	t.Errorf("Expect end transaction \n given: %s \n built: %s", sampleEndTReq, b)
	// }
}

func TestEndTCallBadBody(t *testing.T) {
	req := BuildEndTransactionRequest(sampleConf, samplebinsectoken)
	resp, err := CallEndTransaction(serverBadBody.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverBadBody")
	}
	if err.Error() != resp.ErrorSabreXML.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreXML.Code != sbrerr.BadParse {
		t.Errorf("Expect %d got %d", sbrerr.BadParse, resp.ErrorSabreXML.Code)
	}
	if resp.ErrorSabreXML.AppMessage != sbrerr.ErrCallEndTransaction {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallEndTransaction, resp.ErrorSabreXML.AppMessage)
	}
}

func TestEndTCallBusLogic(t *testing.T) {
	req := BuildEndTransactionRequest(sampleConf, samplebinsectoken)
	resp, err := CallEndTransaction(serverEndTBizLogic.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverBadBody")
	}
	if !resp.Body.Fault.Ok() {
		t.Error("Soap Fault be Ok() since errors was nil")
	}
	appRes := resp.Body.EndTransactionRS.AppResults
	if appRes.Ok() {
		t.Error("Application Results should not be Ok()")
	}
	if len(appRes.Errors) == 0 {
		t.Errorf("Application Results should have errors, want: %d, got %d", len(appRes.Errors), 0)
	}
}

func TestEndTCallResponseBody(t *testing.T) {
	req := BuildEndTransactionRequest(sampleConf, samplebinsectoken)
	resp, err := CallEndTransaction(serverEndT.URL, req)
	if err != nil {
		t.Errorf("Error should be nil: %s", err)
	}
	if resp.Body.EndTransactionRS.ItineraryRef.ID != "NMOXQF" {
		t.Errorf("Expected %s, got %s", "NMOXQF", resp.Body.EndTransactionRS.ItineraryRef.ID)
	}
}
