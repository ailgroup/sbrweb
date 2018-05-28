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
	"time"
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
)

// arriveDepartParser parse string data value into time value.
func arriveDepartParser(arrive, depart string) (time.Time, time.Time) {
	a, _ := time.Parse(timeSpanFormatter, arrive)
	d, _ := time.Parse(timeSpanFormatter, depart)
	return a, d
}
