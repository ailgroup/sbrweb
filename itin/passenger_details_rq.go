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
	SpecialRequest *SpecialRequestDetails
}

type PostProcessing struct{}
type PreProcessing struct{}
type SpecialRequestDetails struct {
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
}
type Document struct {
	IssueCountry struct {
	} `xml:"IssueCountry,omitempty"`
	NationalityCountry struct {
	} `xml:"NationalityCountry,omitempty"`
}
type AgencyInfo struct {
	Address addr.Address
}
type CustomerInfo struct{}
type TravelItineraryAddInfoRQ struct {
}
