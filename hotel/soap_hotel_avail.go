package sbrhotel

import (
	"encoding/xml"
	"fmt"

	"github.com/ailgroup/sbrweb"
)

// OTA_HotelAvailRQ retrieve sabre hotel content using various query criteria, see SearchCriteria
type OTA_HotelAvailRQ struct {
	XMLName           xml.Name `xml:"OTA_HotelAvailRQ"`
	Version           string   `xml:"version,attr"`
	XMLNS             string   `xml:"xmlns,attr"`
	XMLNSXs           string   `xml:"xmlns:xs,attr"`
	XMLNSXsi          string   `xml:"xmlns:xsi,attr"`
	ReturnHostCommand bool
	Avail             AvailRequestSegment
}

type AvailRequestSegment struct {
	XMLName xml.Name `xml:"AvailRequestSegment"`
	Customer
	GuestCounts         GuestCounts
	HotelSearchCriteria HotelSearchCriteria
}

type GuestCounts struct {
	Count int `xml:",attr"`
}
type Customer struct {
	Corporate struct {
		ID string `xml:",omitempty"`
	} `xml:",omitempty"`
}

type HotelSearchCriteria struct {
	Criterion CriterionElem
}

type CriterionElem struct {
	HotelRef []HotelRef
}
type HotelRef struct {
	HotelCityCode string `xml:",attr,omitempty"`
	HotelCode     string `xml:",attr,omitempty"`
	//HotelName     string `xml:",attr,omitempty"`
	//Latitude     string `xml:",attr,omitempty"`
	//Longitude     string `xml:",attr,omitempty"`
}

type HotelQueryParams map[string][]string

func buildHotelSearch(params HotelQueryParams) (HotelSearchCriteria, error) {
	q := HotelSearchCriteria{}
	if len(params) > 1 {
		return q, fmt.Errorf("Cannot use more than 1 search type criteria, have %v", params)
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
		}
	}

	return q, nil
}

func BuildHotelAvailRq(corpID string, guestCount int, query HotelSearchCriteria) OTA_HotelAvailRQ {
	rq := OTA_HotelAvailRQ{
		Version:           hotelAvailVersion,
		XMLNS:             sbrweb.BaseWebServicesNS,
		XMLNSXs:           sbrweb.BaseXSDNameSpace,
		XMLNSXsi:          sbrweb.BaseXSINamespace,
		ReturnHostCommand: returnHostCommand,
		Avail: AvailRequestSegment{
			GuestCounts:         GuestCounts{Count: guestCount},
			HotelSearchCriteria: query,
		},
	}

	if len(corpID) > 0 {
		rq.Avail.Customer = Customer{
			Corporate: struct {
				ID string `xml:",omitempty"`
			}{
				ID: corpID,
			},
		}
	}

	return rq
}
