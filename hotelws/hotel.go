/* Package hotelws implements Sabre hotel searching SOAP payloads through various criterion for hotel availability as well as hotel property descriptions. Many criterion exist that are not yet implemented: (Award, ContactNumbers, CommissionProgram, HotelAmenity, Package, PointOfInterest, PropertyType, RefPoint, RoomAmenity, HotelFeaturesCriterion,). To add more criterion create a criterion type (e.g, XCriterion) as well as its accompanying function to handle the data parms (e.g., XSearch).
 */
package hotelws

import "time"

const (
	hotelRQVersion    = "2.3.0"
	timeSpanFormatter = "01-02"

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
