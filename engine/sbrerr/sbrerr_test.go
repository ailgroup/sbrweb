package sbrerr

import (
	"testing"
)

var (
	errInput = "error input string"
	appInput = "app input string"
)

func TestStatusHelperFuncs(t *testing.T) {
	if StatusNotProcess() != "NotProcessed" {
		t.Errorf("StatusNotProcess() expect: %s, got: %s", "NotProcessed", StatusNotProcess())
	}
	if StatusApproved() != "Approved" {
		t.Errorf("StatusApproved() expect: %s, got: %s", "Approved", StatusApproved())
	}
	if StatusComplete() != "Complete" {
		t.Errorf("StatusComplete() expect: %s, got: %s", "Complete", StatusComplete())
	}

}

var sampleStatuses = []struct {
	code   SabreStatus
	expect string
}{
	{-100, "Unknown"},
	{-3, "Unknown"},
	{-2, "Unknown"},
	{-1, "Unknown"},
	{-0, "Unknown"},
	{0, "Unknown"},
	{1, "Unknown"},
	{2, "BadService"},
	{3, "BadParse"},
	{4, "SoapFault"},
	{5, "NotProcessed"},
	{6, "Approved"},
	{7, "Complete"},
	{100, "Unknown"},
	{223, "Unknown"},
}

func TestAppStatus(t *testing.T) {
	for i, stat := range sampleStatuses {
		if stat.code.String() != stat.expect {
			t.Errorf("sampleStatus: %d for Code: %d => expect: %s, got: %s", i, stat.code, stat.expect, stat.code.String())
		}
	}
}

var sampleCodes = []struct {
	expect SabreStatus
	input  string
}{
	{1, "Unknown"},
	{1, ""},
	{1, "0"},
	{1, "gobbleygook"},
	{1, "HelloBad"},
	{1, "Bad"},
	{1, "NoService"},
	{1, "Unknown"},
	{2, "BadService"},
	{3, "BadParse"},
	{4, "SoapFault"},
	{5, "NotProcessed"},
	{6, "Approved"},
	{7, "Complete"},
}

func TestGetStatus(t *testing.T) {
	for i, c := range sampleCodes {
		if SabreEngineStatusCode(c.input) != c.expect {
			t.Errorf("sampleCodes: %d for input %s expect %d got %d", i, c.input, c.expect, SabreEngineStatusCode(c.input))
		}
	}
}

func TestService(t *testing.T) {
	e := NewErrorSabreService(
		errInput,
		appInput,
		BadService,
	)
	if e.ErrMessage != errInput {
		t.Error("ErrMessage not correct")
	}
	if e.AppMessage != appInput {
		t.Error("AppMessage not correct")
	}
	if e.Code != BadService {
		t.Error("AppStatus Code not correct")
	}
	if e.Error() != errInput {
		t.Error("Error() method not correct")
	}
}
func TestXML(t *testing.T) {
	e := NewErrorSabreXML(
		errInput,
		appInput,
		BadParse,
	)
	if e.ErrMessage != errInput {
		t.Error("ErrMessage not correct")
	}
	if e.AppMessage != appInput {
		t.Error("AppMessage not correct")
	}
	if e.Code != BadParse {
		t.Error("AppStatus Code not correct")
	}
	if e.Error() != errInput {
		t.Error("Error() method not correct")
	}
}

func TestResult(t *testing.T) {
	e := NewErrorSabreResult(
		appInput,
		NotProcessed,
	)
	if e.AppMessage != appInput {
		t.Error("AppMessage not correct")
	}
	if e.Code != NotProcessed {
		t.Error("AppStatus Code not correct")
	}
	if e.Error() != appInput {
		t.Error("Error() method not correct")
	}
}

func TestSoapFault(t *testing.T) {
	e := ErrorSoapFault{
		ErrMessage: errInput,
		FaultCode:  appInput,
		StackTrace: "long stacky:stack here ugh.",
		Code:       SoapFault,
	}
	if e.ErrMessage != errInput {
		t.Error("ErrMessage not correct")
	}
	if e.FaultCode != appInput {
		t.Error("AppMessage not correct")
	}
	if e.StackTrace != "long stacky:stack here ugh." {
		t.Error("AppStatus Code not correct")
	}
	if e.Code != SoapFault {
		t.Error("AppStatus Code not correct")
	}
	if e.Error() != errInput {
		t.Error("Error() method not correct")
	}
}
