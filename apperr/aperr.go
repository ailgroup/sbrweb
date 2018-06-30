package apperr

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
	case Invalid.String():
		return Invalid
	case BadInput.String():
		return BadInput
	default:
		return Unknown
	}
}

// ErrorInvalid container for validation errors
type ErrorInvalid struct {
	ErrMessage string    `json:"err_invalid,omitempty"`
	AppMessage string    `json:"app_invalid,omitempty"`
	Code       AppStatus `json:"app_status"`
	HTTPStatus uint      `json:"http_status"`
}

// NewErrorInvalid for validation errors
func NewErrorInvalid(errIn, appIn string, code AppStatus, httpstatus uint) ErrorInvalid {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorInvalid{ErrMessage: errIn, AppMessage: appIn, Code: code, HTTPStatus: httpstatus}
}

// Error for ErrorInvalid implements std lib error interface
func (e ErrorInvalid) Error() string {
	return e.ErrMessage
}

// ErrorBadInput container for bad input or json parsing
type ErrorBadInput struct {
	ErrMessage string    `json:"err_badinput,omitempty"`
	AppMessage string    `json:"app_badinput,omitempty"`
	Code       AppStatus `json:"app_status"`
	HTTPStatus uint      `json:"http_status"`
}

// NewErrorBadInput for json parsing or other bad input
func NewErrorBadInput(errIn, appIn string, code AppStatus, httpstatus uint) ErrorBadInput {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorBadInput{ErrMessage: errIn, AppMessage: appIn, Code: code, HTTPStatus: httpstatus}
}

// Error for ErrorBadInput implements std lib error interface
func (e ErrorBadInput) Error() string {
	return e.ErrMessage
}
