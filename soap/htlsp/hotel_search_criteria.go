package htlsp

/*
	This file contains functions related to hotel search criteria. See hotel.go file for the struct and type definitions.

	See also hotel_avail.go and
	https://developer.sabre.com/docs/read/soap_apis/hotel/search/hotel_availability
*/

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ailgroup/sbrweb/sbrerr"
)

// NewHotelSearchCriteria accepts set of QueryParams functions, executes over hotel search criteria and returns modified criteria
func NewHotelSearchCriteria(queryParams ...QuerySearchParams) (*HotelSearchCriteria, error) {
	criteria := &HotelSearchCriteria{}
	for _, qm := range queryParams {
		err := qm(criteria)
		if err != nil {
			return criteria, err
		}
	}
	return criteria, nil
}

// validatePropertyRequest ensures property description requests are well-formed
func (c *HotelSearchCriteria) validatePropertyRequest() error {
	for _, criterion := range c.Criterion.HotelRefs {
		if len(criterion.HotelCityCode) > 0 {
			return sbrerr.ErrPropDescCityCode
		}
		if (len(criterion.Latitude) > 0) || (len(criterion.Longitude) > 0) {
			return sbrerr.ErrPropDescLatLng
		}

		if len(c.Criterion.HotelRefs) > 1 {
			return sbrerr.ErrPropDescHotelRefs
		}
	}
	return nil
}

// PackageSearch ... TODO: create validation around the packages that can be used.
func PackageSearch(params PackageCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		for _, p := range params {
			q.Criterion.Packages = append(q.Criterion.Packages, &Package{Val: p})
		}
		return nil
	}
}

// PropertyTypeSearch ... TODO: create validation around the types that can be used.
func PropertyTypeSearch(params PropertyTypeCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		for _, p := range params {
			q.Criterion.PropertyTypes = append(q.Criterion.PropertyTypes, &PropertyType{Val: p})
		}
		return nil
	}
}

// validate for AddressSearchStruct based on what sabre allows
func (a AddressSearchStruct) validate() bool {
	// need postal or country
	if (a.PostalCode == "") && (a.CountryCode == "") {
		return false
	}
	// need city or postal or street
	if (a.PostalCode == "") && (a.CityName == "") && (a.StreetNmbr == "") {
		return false
	}
	return true
}

/*
AddressSearch builds AddressSearch and put it on the serach criterion.
NOTE: Must have country code and/or postal code state province not an acceptable criterion.
AddressSearch requires other criteria such as city code, package, property type.
This criteria is not recommended.
*/
func AddressSearch(params AddressCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("AddressSearch params cannot be empty: %v", params)
		}
		a := &AddressSearchStruct{}
		for k, v := range params {
			switch k {
			case StreetQueryField:
				a.StreetNmbr = v
			case CityQueryField:
				a.CityName = v
			case PostalQueryField:
				a.PostalCode = v
			case CountryCodeQueryField:
				a.CountryCode = v
			}
		}
		if !a.validate() {
			return errors.New("ERROR AddressSearch: Missing postal or country; OR need city, postal, street")
		}
		q.Criterion.AddressSearch = a
		return nil
	}
}

// HotelRefSearch accepts HotelRef criterion and returns a function for hotel search critera.
// Supports CityCode, HotelCode, Latitude-Longitude for now... later support for HotelName.
func HotelRefSearch(params HotelRefCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("HotelRefCriterion params cannot be empty: %v", params)
		}
		for k, v := range params {
			switch k {
			case CityQueryField:
				for _, city := range v {
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{HotelCityCode: city})
				}
			case HotelidQueryField:
				for _, code := range v {
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{HotelCode: code})
				}
			case LatlngQueryField:
				for _, l := range v {
					latlng := strings.Split(l, ",")
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{Latitude: latlng[0], Longitude: latlng[1]})
				}
			}
		}
		return nil
	}
}

// PointOfInterestSearch for hotel availability searches based on named points of interest.
// Supports CountryStateCode (for country or state), Val for PointOfInterest (city or landmark).
// NOTE: this may only include properties that are located in the city centre or surrounding areas.
// Other fields sit on the PointOfInterest struct, but NonUS we want default(false), and RPH not needed but there for correctness.
// Cryptic: HOTUT-ST GEORGE/21APR-24APR2 [HOT(StateCode)-(CityName)/(Arrive)-(Depart)GuestNumber]
func PointOfInterestSearch(params PointOfInterestCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("PointOfInterestSearch params cannot be empty: %v", params)
		}
		p := &PointOfInterest{}
		for k, v := range params {
			switch k {
			case StateCodeQueryField:
				p.CountryStateCode = v
			case POIQueryField:
				p.Val = v
			}
		}

		q.Criterion.PointOfInterest = p
		return nil
	}
}
