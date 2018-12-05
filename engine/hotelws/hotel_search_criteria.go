package hotelws

/*
This file contains functions related to hotel search criteria. See hotelws.go file for the struct and type definitions.
*/

import (
	"fmt"
	"strings"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

// SetRateParams helper to create a slice of rate plans to append on a an Avail Segement
// for search or description services
func SetRateParams(ratePlans []RatePlan) *RatePlanCandidates {
	rpc := &RatePlanCandidates{}
	for _, plan := range ratePlans {
		rpc.RatePlans = append(rpc.RatePlans, &plan)
	}
	return rpc
}

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

// AddressSearch parse incoming params, build Address, and put it on the serach criterion
// search on small towns..
func AddressSearch(params AddressCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		a := &Address{}
		if len(params) < 1 {
			return fmt.Errorf("AddressSearch params cannot be empty: %v", params)
		}
		for k, v := range params {
			switch k {
			case StreetQueryField:
				a.Street = v
			case CityQueryField:
				a.City = v
			case PostalQueryField:
				a.Postal = v
			case CountryCodeQueryField:
				a.CountryCode = v
			}
		}
		q.Criterion.Address = a
		return nil
	}
}

// HotelRefSearch accepts HotelRef criterion and returns a function for hotel search critera.
// Supports CityCode, HotelCode, Latitude, and Longitude for now... later support for HotelName.
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
