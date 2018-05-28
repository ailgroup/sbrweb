package sbrerr

import (
	"testing"
)

var (
	errInput = "error input string"
	appInput = "app input string"
)

var sampleStatuses = []struct {
	code   AppStatus
	expect string
}{
	{1, "Unknown"},
	{2, "BadService"},
	{3, "BadParse"},
	{-1, "Unknown"},
	{-2, "Unknown"},
	{-3, "Unknown"},
	{5, "Unknown"},
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
