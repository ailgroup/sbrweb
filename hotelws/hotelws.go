/* Package hotelws (Hotel Web Services) implements Sabre SOAP hotel booking through availability, property and rate descriptions, passenger details (PNR), reservation, and transaction Web Services. It handles exclusively the XML-based Simple Object Access Protocol endpoints the Sabre supports. Look elsewhere for support of Sabre rest services.

# Typical Workflow definition from Sabre Web Services
	Book Hotel Reservation
	The following workflow allows you to search and book a hotel room.
	Steps

	--Step 1: Retrieve hotel availability using OTA_HotelAvailLLSRQ.
	--Step 2: Retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
	--Step 3: Retrieve hotel rules and policies using HotelRateDescriptionLLSRQ.*
	--Step 4: Add any additional (required) information to create the passenger name record (PNR) using PassengerDetailsRQ.**
	--Step 5: Book a room for the selected hotel using OTA_HotelResLLSRQ.
	--Step 6: End the transaction of the passenger name record using EndTransactionLLSRQ.
	Note

	* Mandatory only if selected option in response of HotelPropertyDescriptionLLSRQ contains HRD_RequiredForSell="true".
	** Ensure Agency address is added within call to PassengerDetails, so as the OTA_HotelResLLSRQ call is not rejected.
*/
package hotelws

import (
	"errors"
	"time"
)

type AppStatus int

const (
	Unknown AppStatus = 1 << (2 * iota)
	BadService
	BadParse
)

const (
	hotelRQVersion        = "2.3.0"
	timeSpanFormatter     = "01-02"
	streetQueryField      = "street_qf"
	cityQueryField        = "city_qf"
	postalQueryField      = "postal_qf"
	countryCodeQueryField = "countryCode_qf"
	latlngQueryField      = "latlng_qf"
	hotelidQueryField     = "hotelID_qf"
	returnHostCommand     = true
	StatusComplete        = "Complete"
	ErrCallHotelAvail     = "Error CallHotelAvail::OTA_HotelAvailLLSRQ"
	ErrCallHotelPropDesc  = "Error CallHotelPropDesc::HotelPropertyDescriptionLLSRQ"
	ErrCallHotelRateDesc  = "Error CallHotelRateDesc::HotelRateDescriptionLLSRQ"
	ErrCallHotelRes       = "Error CallHotelRes::OTA_HotelResLLSRQ"
	//ErrCallPNRDetails     = "Error CallPNRDetails::PassengerDetailsRQ" TODO in its own package
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
	AppMessage string    `json:"app_message_sabre_service,omitempty"`
	ErrMessage string    `json:"err_message_sabre_service,omitempty"`
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
	AppMessage string    `json:"app_message_sabre_xml,omitempty"`
	ErrMessage string    `json:"err_message_sabre_xml,omitempty"`
	Code       AppStatus `json:"app_status"`
}

func NewErrorErrorSabreXML(errIn, appIn string, code AppStatus) ErrorSabreXML {
	//err = strings.Replace(err, "\n", "", -1)
	return ErrorSabreXML{ErrMessage: errIn, AppMessage: appIn, Code: code}
}

// ErrorSabreXML implements std lib error interface
func (e ErrorSabreXML) Error() string {
	return e.ErrMessage
}

// arriveDepartParser parse string data value into time value.
func arriveDepartParser(arrive, depart string) (time.Time, time.Time) {
	a, _ := time.Parse(timeSpanFormatter, arrive)
	d, _ := time.Parse(timeSpanFormatter, depart)
	return a, d
}
