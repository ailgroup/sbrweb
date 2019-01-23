package sbrerr

import (
	"errors"
)

type SabreStatus int

// List statuses for common error handling and parsing. Generally, the lower number are more serious.
const (
	Unknown      SabreStatus = iota + 1 //1
	BadService                          //2
	BadParse                            //3
	SoapFault                           //4
	NotProcessed                        //5
	Approved                            //6
	Complete                            //7
)

const (
	ErrCallSessionCreate   = "Error CallSessionCreate::SessionCreateRQ"
	ErrCallSessionClose    = "Error CallSessionClose::SessionCloseRQ"
	ErrCallSessionValidate = "Error CallSessionValidate::SessionValidateRQ"
	ErrCallHotelAvail      = "Error CallHotelAvail::OTA_HotelAvailLLSRQ"
	ErrCallHotelPropDesc   = "Error CallHotelPropDesc::HotelPropertyDescriptionLLSRQ"
	ErrCallHotelRateDesc   = "Error CallHotelRateDesc::HotelRateDescriptionLLSRQ"
	ErrCallHotelRes        = "Error CallHotelRes::OTA_HotelResLLSRQ"
	ErrCallPNRDetails      = "Error CallPNRDetails::PassengerDetailsRQ"
	ErrCallEndTransaction  = "Error CallEndTransaction::EndTransactionRQ"
)

var (
	ErrPropDescCityCode  = errors.New("HotelCityCode not allowed in HotelPropDesc")
	ErrPropDescLatLng    = errors.New("Latitude or Longitude not allowed in HotelPropDesc")
	ErrPropDescHotelRefs = errors.New("Criterion.HotelRef cannot be greater than 1, can only search using one criterion")

	// sabreEngineStatuses strings to map to consts....
	// TODO come back and refactor to something less fragile
	sabreEngineStatuses = [...]string{
		"0", //place holder so we can get index values starting 1 :)
		"Unknown",
		"BadService",
		"BadParse",
		"SoapFault",
		"NotProcessed",
		"Approved",
		"Complete",
	}
)

// TODO come back and refactor statuses to something less fragile
// StatusNotProcess returns enum value for sabre string
func StatusNotProcess() string {
	return sabreEngineStatuses[NotProcessed]
}

// StatusApproved returns enum value for sabre string
func StatusApproved() string {
	return sabreEngineStatuses[Approved]
}

// StatusComplete returns enum value for sabre string
func StatusComplete() string {
	return sabreEngineStatuses[Complete]
}

// String reverse lookup for statuses, get string value for code
func (code SabreStatus) String() string {
	if code < Unknown || code > Complete {
		return "Unknown"
	}
	return sabreEngineStatuses[code]
}

// SabreEngineStatusCode input is xml parsed string value from sabre that we match with String() function to get the error code.
func SabreEngineStatusCode(input string) SabreStatus {
	if input == "0" {
		return Unknown
	}
	switch input {
	case NotProcessed.String():
		return NotProcessed
	case Approved.String():
		return Approved
	case Complete.String():
		return Complete
	case BadParse.String():
		return BadParse
	case BadService.String():
		return BadService
	case SoapFault.String():
		return SoapFault
	default:
		return Unknown
	}
}

// ErrorSabreService container for network issues
type ErrorSabreService struct {
	ErrMessage string `json:",omitempty"`
	AppMessage string `json:",omitempty"`
	Code       SabreStatus
}

// NewErrorSabreService for http or sabre web services networking problems
func NewErrorSabreService(errIn, appIn string, code SabreStatus) ErrorSabreService {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreService{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// Error for ErrorSabreService implements std lib error interface
func (e ErrorSabreService) Error() string {
	return e.ErrMessage
}

// ErrorSabreXML container for xml issues
type ErrorSabreXML struct {
	ErrMessage string `json:",omitempty"`
	AppMessage string `json:",omitempty"`
	Code       SabreStatus
}

// ErrorSabreXML for parsing web services responses
func NewErrorSabreXML(errIn, appIn string, code SabreStatus) ErrorSabreXML {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreXML{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// Error for ErrorSabreXML implements std lib error interface
func (e ErrorSabreXML) Error() string {
	return e.ErrMessage
}

// ErrorSabreResult for results issues
type ErrorSabreResult struct {
	AppMessage string `json:",omitempty"`
	Code       SabreStatus
}

// NewErrorSabreResult for response results with errors(bad dates, credit card, etc...)
func NewErrorSabreResult(appIn string, code SabreStatus) ErrorSabreResult {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreResult{AppMessage: appIn, Code: code}
}

// Error for ErrorSabreResult implements std lib error interface
func (e ErrorSabreResult) Error() string {
	return e.AppMessage
}

// ErrorSoapFault for results issues
type ErrorSoapFault struct {
	ErrMessage string `json:",omitempty"`
	FaultCode  string `json:",omitempty"`
	StackTrace string `json:",omitempty"`
	Code       SabreStatus
}

// NewErrorSoapFault for response results with errors(bad dates, credit card, etc...)
func NewErrorSoapFault(errIn string) ErrorSoapFault {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSoapFault{ErrMessage: errIn, Code: SoapFault}
}

// Error for ErrorSoapFault implements std lib error interface
func (e ErrorSoapFault) Error() string {
	return e.ErrMessage
}
