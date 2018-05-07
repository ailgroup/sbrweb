package hotel

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// QueryParams is a typed function to support optional query params on creation of new search criterion
type QueryParams func(*HotelSearchCriteria) error

// HotelRefCriterion map of hotel ref criteria
type HotelRefCriterion map[string][]string

// AddressCriterion map of address search criteria
type AddressCriterion map[string]string

// PropertyTypeCriterion slice of property type strings (APTS, LUXRY)
type PropertyTypeCriterion []string

// PackageCriterion slice of property type strings (GF, HM, BB)
type PackageCriterion []string

// Timepsan for arrival and departure params
type TimeSpan struct {
	XMLName xml.Name `xml:"TimeSpan"`
	Depart  string   `xml:"End,attr"`
	Arrive  string   `xml:"Start,attr"`
}

// HotelSearchCriteria top level element for criterion
type HotelSearchCriteria struct {
	XMLName   xml.Name `xml:"HotelSearchCriteria"`
	Criterion Criterion
}

// HotelRef contains any number of search criteria under the HotelRef element.
type HotelRef struct {
	XMLName       xml.Name `xml:"HotelRef,omitempty"`
	HotelCityCode string   `xml:"HotelCityCode,attr,omitempty"`
	HotelCode     string   `xml:"HotelCode,attr,omitempty"`
	HotelName     string   `xml:"HotelName,attr,omitempty"`
	Latitude      string   `xml:"Latitude,attr,omitempty"`
	Longitude     string   `xml:"Longitude,attr,omitempty"`
}

// Address represents typical building addresses
type Address struct {
	City        string `xml:"CityName,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
	Street      string `xml:"StreetNumber,omitempty"`
}

// PropertyType container for searhing types of properties (APTS, LUXRY...)
type PropertyType struct {
	Val string `xml:",chardata"`
}

// Package container for searching types of packages (GF, HM, BB...)
type Package struct {
	Val string `xml:",chardata"`
}

// Criterion holds various serach criteria
type Criterion struct {
	XMLName       xml.Name `xml:"Criterion"`
	HotelRefs     []*HotelRef
	Address       *Address
	PropertyTypes []*PropertyType
	Packages      []*Package
}

// GuestCounts how many guests per night-room. TODO: check on Sabre validation limits (think it is 4)
type GuestCounts struct {
	XMLName xml.Name `xml:"GuestCounts"`
	Count   int      `xml:"Count,attr"`
}

// Customer for corporate or typical sabre customer ids
type Customer struct {
	XMLName    xml.Name    `xml:"Customer,omitempty"`
	Corporate  *Corporate  //nil pointer ignored if empty
	CustomerID *CustomerID //nil pointer ignored if empty
}

// CustomerID number
type CustomerID struct {
	XMLName xml.Name `xml:"ID,omitempty"`
	Number  string   `xml:"Number,omitempty"`
}

// Corporate customer id
type Corporate struct {
	XMLName xml.Name `xml:"Corporate,omitempty"`
	ID      string   `xml:"ID,omitempty"`
}

// NewHotelSearchCriteria accepts set of QueryParams functions, executes over hotel search criteria and returns modified criteria
func NewHotelSearchCriteria(queryParams ...QueryParams) (HotelSearchCriteria, error) {
	criteria := &HotelSearchCriteria{}
	for _, qm := range queryParams {
		err := qm(criteria)
		if err != nil {
			return *criteria, err
		}
	}
	return *criteria, nil
}

// PackageSearch ...
func PackageSearch(params PackageCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		for _, p := range params {
			q.Criterion.Packages = append(q.Criterion.Packages, &Package{Val: p})
		}
		return nil
	}
}

// PropertyTypeSearch ...
func PropertyTypeSearch(params PropertyTypeCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		for _, p := range params {
			q.Criterion.PropertyTypes = append(q.Criterion.PropertyTypes, &PropertyType{Val: p})
		}
		return nil
	}
}

// AddressSearch parse incoming params, build Address, and put it on the serach criterion
func AddressSearch(params AddressCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		a := &Address{}
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
// Supports CityCode, HotelCode, Latitude, and Longitude for now... later support for HotelName.
func HotelRefSearch(params HotelRefCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("HotelRefCriterion params cannot be empty: %v", params)
		}
		for k, v := range params {
			switch k {
			case cityQueryField:
				for _, city := range v {
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{HotelCityCode: city})
				}
			case hotelidQueryField:
				for _, code := range v {
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{HotelCode: code})
				}
			case latlngQueryField:
				for _, l := range v {
					latlng := strings.Split(l, ",")
					q.Criterion.HotelRefs = append(q.Criterion.HotelRefs, &HotelRef{Latitude: latlng[0], Longitude: latlng[1]})
				}
			}
		}
		return nil
	}
}
