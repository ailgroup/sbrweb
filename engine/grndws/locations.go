/*
The Car Locations API allows you to search for vehicles using more than just an airport code. The Car Locations API searches for and returns locations, including vehicle rates and availability, using the following the search criteria:

    Latitude and longitude
    City names
    Addresses
    Zip codes
    Points of interest passed as brief text descriptions (for example, “Walt Disney World” or "Eiffel Tower")
    Hotel location codes
    Hotel segment in a booked Passenger Name Record (PNR)
    Rail station name

You can even further specify your search by including additional shop attributes, such as:

    Corporate discount numbers
    Customer ID numbers
    Vehicle type
    Individual or a subset of car suppliers
    One-way rentals
    Car Extras
*/
package grndws

import (
	"encoding/xml"

	"github.com/ailgroup/sbrweb/engine/srvc"
)

// VehicleLocationRequest for soap package on VehLocationFinderRQ service
type VehicleLocationRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   VehicleLocationBody
}
type VehicleLocationBody struct {
	XMLName             xml.Name `xml:"soap-env:Body"`
	VehLocationFinderRQ VehLocationFinderRQ
}
type VehLocationFinderRQ struct {
	XMLName  xml.Name `xml:"VehLocationFinderRQ" json:"-"`
	Version  string   `xml:"Version,attr"`
	XMLNS    string   `xml:"xmlns,attr"`     //srvc.BaseWebServicesNS
	XMLNSXs  string   `xml:"xmlns:xs,attr"`  //srvc.BaseXSDNameSpace
	XMLNSXsi string   `xml:"xmlns:xsi,attr"` //srvc.BaseXSINamespace
}

// VehAvailRQCore wrapper for location details and rental times
type VehAvailRQCore struct {
	XMLName xml.Name `xml:"VehAvailCore" json:"-"`
	Pickup  string   `xml:"PickUpDateTime,attr"` //"12-22T09:00"
	Return  string   `xml:"ReturnDateTime,attr"` //"12-29T11:00"
}

// LocationDetails generates different location params
type LocationDetails struct {
	XMLName xml.Name `xml:"LocationDetails" json:"-"`
	Dropoff bool     `xml:"DropOff,attr,omitempty"`
	Address Address
}

// Address represents typical building addresses
type Address struct {
	//AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNmbr,omitempty"`
	City          string `xml:"CityName,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr"`
	} `xml:"StateCountyProv,omitempty"`
	CountryCode     string `xml:"CountryCode,omitempty"`
	Postal          string `xml:"PostalCode,omitempty"`
	Latitude        string `xml:"Latitude,attr,omitempty"`
	Longitude       string `xml:"Longitude,attr,omitempty"`
	CounterLocation string `xml:"X,attr,omitempty"`               //example: "N"
	Direction       string `xml:"CounterLocation,attr,omitempty"` //example: "SE"
	Distance        string `xml:"Distance,attr,omitempty"`        //example: "9.3"
	LocationCode    string `xml:"LocationCode,attr,omitempty"`    //example: "DFW"
	LocationName    string `xml:"LocationName,attr,omitempty"`    //example: "DALLAS FORT WORTH"
	LocationOwner   string `xml:"LocationOwner,attr,omitempty"`   //example: "C"
	UnitOfMeasure   string `xml:"UnitOfMeasure,attr,omitempty"`   //example: "MI"
}
