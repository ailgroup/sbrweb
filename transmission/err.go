package transmission

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type AppStatus int

// List statuses for common error handling and parsing.
const (
	Unknown  AppStatus = iota + 1 //1
	BadInput                      //2
	Invalid                       //3
)

var (
	appStatuses = [...]string{
		"0",
		"Unknown",
		"BadInput",
		"Invalid",
	}
)

func StatusUnknown() string {
	return appStatuses[Unknown]
}
func StatusInvalid() string {
	return appStatuses[Invalid]
}
func StatusBadInput() string {
	return appStatuses[BadInput]
}
func (code AppStatus) String() string {
	if code < Unknown || code > Invalid {
		return "Unknown"
	}
	return appStatuses[code]
}
func AppStatusCode(input string) AppStatus {
	if input == "0" {
		return Unknown
	}
	switch input {
	case BadInput.String():
		return BadInput
	case Invalid.String():
		return Invalid
	default:
		return Unknown
	}
}

// ErrorUnknown container for validation errors
type ErrorUnknown struct {
	ErrMessage string `json:",omitempty"`
	AppMessage string `json:",omitempty"`
	Code       AppStatus
	HTTPStatus uint
	ServerTime time.Time
}

// NewErrorUnknown for validation errors
func NewErrorUnknown(errIn, appIn string, code AppStatus, hts uint) ErrorUnknown {
	return ErrorUnknown{ErrMessage: errIn, AppMessage: appIn, Code: code, HTTPStatus: hts, ServerTime: time.Now()}
}

// Error for ErrorUnknown implements std lib error interface
func (e ErrorUnknown) Error() string {
	return e.ErrMessage
}

// DecodeInvalid builds an ErrorInvalid response given a set of url.Values
func DecodeUnknown(handlerMsg string, val url.Values, err error, httpstatus uint) []byte {
	b, _ := json.Marshal(
		NewErrorUnknown(
			err.Error(),
			fmt.Sprintf("%s. Unknown for Query: %v", handlerMsg, val),
			Unknown,
			httpstatus,
		),
	)
	return b
}

// ErrorInvalid container for validation errors
type ErrorInvalid struct {
	ErrMessage string `json:",omitempty"`
	AppMessage string `json:",omitempty"`
	Code       AppStatus
	HTTPStatus uint
	ServerTime time.Time
}

// NewErrorInvalid for validation errors
func NewErrorInvalid(errIn, appIn string, code AppStatus, hts uint) ErrorInvalid {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorInvalid{ErrMessage: errIn, AppMessage: appIn, Code: code, HTTPStatus: hts, ServerTime: time.Now()}
}

// Error for ErrorInvalid implements std lib error interface
func (e ErrorInvalid) Error() string {
	return e.ErrMessage
}

// DecodeInvalid builds an ErrorInvalid response given a set of url.Values
func DecodeInvalid(handlerMsg string, err error, httpstatus uint) []byte {
	b, _ := json.Marshal(
		NewErrorInvalid(
			err.Error(),
			fmt.Sprintf("%s: %s.", "Invalid", handlerMsg),
			Invalid,
			httpstatus,
		),
	)
	return b
}

// ErrorBadInput container for bad input or json parsing
type ErrorBadInput struct {
	ErrMessage string `json:",omitempty"`
	AppMessage string `json:",omitempty"`
	Code       AppStatus
	HTTPStatus uint
	ServerTime time.Time
}

// NewErrorBadInput for json parsing or other bad input
func NewErrorBadInput(errIn, appIn string, code AppStatus, hts uint) ErrorBadInput {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorBadInput{ErrMessage: errIn, AppMessage: appIn, Code: code, HTTPStatus: hts, ServerTime: time.Now()}
}

// Error for ErrorBadInput implements std lib error interface
func (e ErrorBadInput) Error() string {
	return e.ErrMessage
}

// DecodeBadInput builds an ErrorBadInput response given a set of url.Values
func DecodeBadInput(handlerName string, val url.Values, err error, httpstatus uint) []byte {
	b, _ := json.Marshal(
		NewErrorBadInput(
			err.Error(),
			fmt.Sprintf("%s. Cannot Decode Query: %v", handlerName, val),
			BadInput,
			httpstatus,
		),
	)
	return b
}
