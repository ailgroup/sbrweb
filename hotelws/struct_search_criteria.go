package hotelws

import "encoding/xml"

/*
Implement Sabre hotel searching through various criteria. Many criterion exist that are not yet implemented: (Award, ContactNumbers, CommissionProgram, HotelAmenity, PointOfInterest, RefPoint, RoomAmenity, HotelFeaturesCriterion). To add more criterion create a criterion type (e.g, XCriterion) as well as its accompanying function to handle the data params (e.g., XSearch); see examples in hotel_search_criteria.go.
*/

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

// AvailAvailRequestSegment holds basic hotel availability params: customer ids, guest count, hotel search criteria and arrival departure
type AvailRequestSegment struct {
	XMLName             xml.Name  `xml:"AvailRequestSegment"`
	Customer            *Customer //nil pointer ignored if empty
	GuestCounts         GuestCounts
	HotelSearchCriteria HotelSearchCriteria
	ArriveDepart        TimeSpan `xml:"TimeSpan"`
}
