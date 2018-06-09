package sbrerr

import (
	"errors"
)

type AppStatus int

// List statuses for common error handling and parsing. Generally, the lower number are more serious.
const (
	Unknown      AppStatus = iota + 1 //1
	BadService                        //2
	BadParse                          //3
	SoapFault                         //4
	NotProcessed                      //5
	Approved                          //6
	Complete                          //7
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
)

var (
	ErrPropDescCityCode  = errors.New("HotelCityCode not allowed in HotelPropDesc")
	ErrPropDescLatLng    = errors.New("Latitude or Longitude not allowed in HotelPropDesc")
	ErrPropDescHotelRefs = errors.New("Criterion.HotelRef cannot be greater than 1, can only search using one criterion")

	appStatuses = [...]string{
		"0",
		"Unknown",
		"BadService",
		"BadParse",
		"SoapFault",
		"NotProcessed",
		"Approved",
		"Complete",
	}
)

func StatusNotProcess() string {
	return appStatuses[NotProcessed]
}
func StatusApproved() string {
	return appStatuses[Approved]
}
func StatusComplete() string {
	return appStatuses[Complete]
}
func (code AppStatus) String() string {
	if code < Unknown || code > Complete {
		return "Unknown"
	}
	return appStatuses[code]
}
func AppStatusCode(input string) AppStatus {
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
	ErrMessage string    `json:"err_message_sabre_service,omitempty"`
	AppMessage string    `json:"app_message_sabre_service,omitempty"`
	Code       AppStatus `json:"app_status"`
}

// NewErrorSabreService for http or sabre web services networking problems
func NewErrorSabreService(errIn, appIn string, code AppStatus) ErrorSabreService {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreService{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// Error for ErrorSabreService implements std lib error interface
func (e ErrorSabreService) Error() string {
	return e.ErrMessage
}

// ErrorSabreXML container for xml issues
type ErrorSabreXML struct {
	ErrMessage string    `json:"err_message_sabre_xml,omitempty"`
	AppMessage string    `json:"app_message_sabre_xml,omitempty"`
	Code       AppStatus `json:"app_status"`
}

// ErrorSabreXML for parsing web services responses
func NewErrorSabreXML(errIn, appIn string, code AppStatus) ErrorSabreXML {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreXML{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// Error for ErrorSabreXML implements std lib error interface
func (e ErrorSabreXML) Error() string {
	return e.ErrMessage
}

// ErrorSabreResult for results issues
type ErrorSabreResult struct {
	AppMessage string    `json:"app_message_sabre_results,omitempty"`
	Code       AppStatus `json:"app_status"`
}

// NewErrorSabreResult for response results with errors(bad dates, credit card, etc...)
func NewErrorSabreResult(appIn string, code AppStatus) ErrorSabreResult {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreResult{AppMessage: appIn, Code: code}
}

// Error for ErrorSabreResult implements std lib error interface
func (e ErrorSabreResult) Error() string {
	return e.AppMessage
}

// ErrorSoapFault for results issues
type ErrorSoapFault struct {
	ErrMessage string    `json:"err_message_soap_fault_string,omitempty"`
	FaultCode  string    `json:"soap_fault_code,omitempty"`
	StackTrace string    `json:"soap_fault_stacktrace,omitempty"`
	Code       AppStatus `json:"app_status"`
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
