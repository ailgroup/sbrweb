package sbrhotel

import (
	"fmt"
	"strings"
)

/*
	propertyQueryField = "property"
	locationQueryField = "location"
	amenityQueryField  = "amenity"
	airportQueryField  = "airport"
	creditQueryField   = "creditcard"
	diningTypeField    = "Dining"
	alertTypeField     = "Alerts"

	sabreHotelContentVersion = "1.0.0"

	// LatLng represents a location on the Earth.
	type LatLng struct {
		Lat float64
		Lng float64
	}
*/

const (
	hotelAvailVersion = "2.3.0"
	timeSpanFormatter = "01-02"

	streetQueryField      = "street_qf"
	cityQueryField        = "city_qf"
	postalQueryField      = "postal_qf"
	countryCodeQueryField = "countryCode_qf"
	latlngQueryField      = "latlng_qf"
	hotelidQueryField     = "hotelID_qf"
	returnHostCommand     = true
)

// Address represents typical building addresses
type Address struct {
	City        string `xml:"CityName,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
	Street      string `xml:"StreetNumber,omitempty"`
}

// QueryParams is a typed function to support optional query params on creation of new search criterion
type QueryParams func(*HotelSearchCriteria) error

/*
	Many criterion exist:
		Award
		ContactNumbers
		CommissionProgram
		HotelAmenity
		Package
		PointOfInterest
		PropertyType
		RefPoint
		RoomAmenity
	only implementing these for now
*/
//type HotelFeaturesCriterion map[string][]string
type HotelRefCriterion map[string][]string
type AddressCriterion map[string]string

func NewHotelSearchCriteria(queryParams ...QueryParams) (HotelSearchCriteria, error) {
	criteria := &HotelSearchCriteria{}
	for _, qm := range queryParams {
		err := qm(criteria)
		if err != nil {
			return HotelSearchCriteria{}, err
		}
	}
	return *criteria, nil
}

// AddressOption todo...
func AddressSearch(params AddressCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		a := Address{}
		if len(params) < 1 {
			return fmt.Errorf("AddressSearch params cannot be empty: %v", params)
		}
		for k, v := range params {
			switch k {
			case streetQueryField:
				a.Street = v
			case cityQueryField:
				a.City = v
			case postalQueryField:
				a.Postal = v
			case countryCodeQueryField:
				a.CountryCode = v
			}
		}
		q.Criterion.Address = a
		return nil
	}
}

// HotelRefSearch accepts HotelRef criterion and returns a function for hotel search critera.
// Supports CityCode and HotelCode for now... later support for HotelName, Latitude, Longitude.
func HotelRefSearch(params HotelRefCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("HotelRefCriterion params cannot be empty: %v", params)
		}
		for k, v := range params {
			switch k {
			case cityQueryField:
				for _, city := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCityCode: city})
				}
			case hotelidQueryField:
				for _, code := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCode: code})
				}
			case latlngQueryField:
				for _, l := range v {
					latlng := strings.Split(l, ",")
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{Latitude: latlng[0], Longitude: latlng[1]})
				}
			}
		}
		return nil
	}
}
