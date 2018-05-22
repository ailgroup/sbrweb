package itin

import (
	"encoding/xml"

	"github.com/ailgroup/sbrweb/srvc"
)

type PsngrDetRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   PassengerDetailBody
}

type PassengerDetailBody struct {
	XMLName            xml.Name `xml:"soap-env:Body"`
	PassengerDetailsRQ PassengerDetailsRQ
}
type PassengerDetailsRQ struct {
	XMLName        xml.Name `xml:"PassengerDetailsRQ"`
	XMLNS          string   `xml:"xmlns,attr"`
	Version        string   `xml:"version,attr"`
	IgnoreOnError  string   `xml:"IgnoreOnError,attr"`
	HaltOnError    string   `xml:"HaltOnError,attr"`
	PostProcess    PostProcessing
	PreProcess     PreProcessing
	SpecialReq     *SpecialReqDetails
	TravelItinInfo TravelItineraryAddInfoRQ
}

type PostProcessing struct{}
type PreProcessing struct{}
type SpecialReqDetails struct {
	XMLName          xml.Name `xml:"SpecialRequestDetails"`
	SpecialServiceRQ SpecialServiceRQ
}
type SpecialServiceRQ struct {
	XMLName            xml.Name `xml:"SpecialServiceRQ"`
	SpecialServiceInfo SpecialServiceInfo
}
type SpecialServiceInfo struct {
	XMLName           xml.Name `xml:"SpecialServiceInfo"`
	AdvancedPassenger AdvancedPassenger
}
type AdvancedPassenger struct {
	XMLName       xml.Name `xml:"AdvancePassenger"`
	SegmentNumber string   `xml:"SegmentNumber,attr"`
	Document      Document
	PersonName    PersonName
	VendorPrefs   VendorPrefs
}
type Document struct {
	IssueCountry struct {
	} `xml:"IssueCountry,omitempty"`
	NationalityCountry struct {
	} `xml:"NationalityCountry,omitempty"`
}
type GivenName struct {
	XMLName xml.Name `xml:"GivenName"`
	Val     string   `xml:",chardata"`
}
type MiddleName struct {
	XMLName xml.Name `xml:"MiddleName"`
	Val     string   `xml:",chardata"`
}
type LastName struct {
	XMLName xml.Name `xml:"LastName"`
	Val     string   `xml:",chardata"`
}
type PersonName struct {
	XMLName       xml.Name `xml:"PersonName"`
	NameNumber    string   `xml:"NameNumber,attr"`    //1.1
	NameReference string   `xml:"NameReference,attr"` //ABC123
	PassengerType string   `xml:"PassengerType,attr"` //ADT
	Given         GivenName
	Middle        *MiddleName
	Last          LastName
}
type VendorPrefs struct {
	XMLName xml.Name `xml:"VendorPrefs"`
	Airline struct {
		Hosted bool `xml:"Hosted,attr"`
	} `xml:"Airline"`
}

// Address represents typical building addresses
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNumber,omitempty"`
	City          string `xml:"CityName,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr"`
	} `xml:"StateCountyProv,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
}
type AgencyInfo struct {
	Address     Address
	VendorPrefs VendorPrefs
}
type ContactNumber struct {
	XMLName      xml.Name `xml:"ContactNumber"`
	NameNumber   string   `xml:"NameNumber,attr"`   //1.1
	Phone        string   `xml:"Phone,attr"`        //123-456-7890
	PhoneUseType string   `xml:"PhoneUseType,attr"` //H|M
}
type CustomerInfo struct {
	ContactNumbers []ContactNumber `xml:"ContactNumbers>ContactNumber"`
	PersonName     PersonName
}
type TravelItineraryAddInfoRQ struct {
	Agency   AgencyInfo
	Customer CustomerInfo
}
