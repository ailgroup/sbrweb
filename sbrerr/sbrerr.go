package sbrerr

import (
	"errors"
)

type AppStatus int

const (
	Unknown AppStatus = 1 << (2 * iota)
	BadService
	BadParse
)

const (
	StatusComplete       = "Complete"
	ErrCallHotelAvail    = "Error CallHotelAvail::OTA_HotelAvailLLSRQ"
	ErrCallHotelPropDesc = "Error CallHotelPropDesc::HotelPropertyDescriptionLLSRQ"
	ErrCallHotelRateDesc = "Error CallHotelRateDesc::HotelRateDescriptionLLSRQ"
	ErrCallHotelRes      = "Error CallHotelRes::OTA_HotelResLLSRQ"
	ErrCallPNRDetails    = "Error CallPNRDetails::PassengerDetailsRQ" //TODO in its own package
)

var (
	ErrPropDescCityCode  = errors.New("HotelCityCode not allowed in HotelPropDesc")
	ErrPropDescLatLng    = errors.New("Latitude or Longitude not allowed in HotelPropDesc")
	ErrPropDescHotelRefs = errors.New("Criterion.HotelRef cannot be greater than 1, can only search using one criterion")
)

func (code AppStatus) String() string {
	statuses := [...]string{
		"Unknown",
		"BadService",
		"BadParse",
	}
	if code < BadParse || code < Unknown {
		return "Unknown"
	}
	return statuses[code]
}

type ErrorSabreService struct {
	ErrMessage string    `json:"err_message_sabre_service,omitempty"`
	AppMessage string    `json:"app_message_sabre_service,omitempty"`
	Code       AppStatus `json:"app_status"`
}

func NewErrorSabreService(errIn, appIn string, code AppStatus) ErrorSabreService {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreService{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// ErrorSabreService implements std lib error interface
func (e ErrorSabreService) Error() string {
	return e.ErrMessage
}

type ErrorSabreXML struct {
	ErrMessage string    `json:"err_message_sabre_xml,omitempty"`
	AppMessage string    `json:"app_message_sabre_xml,omitempty"`
	Code       AppStatus `json:"app_status"`
}

func NewErrorSabreXML(errIn, appIn string, code AppStatus) ErrorSabreXML {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreXML{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// ErrorSabreXML implements std lib error interface
func (e ErrorSabreXML) Error() string {
	return e.ErrMessage
}
