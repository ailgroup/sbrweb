package sbrhotel

import (
	"encoding/xml"
	"time"

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
	XMLName             xml.Name  `xml:"AvailRequestSegment"`
	Customer            *Customer //nil pointer ignored if empty
	GuestCounts         GuestCounts
	HotelSearchCriteria HotelSearchCriteria
	ArriveDepart        TimeSpan `xml:"TimeSpan"`
}
type TimeSpan struct {
	XMLName xml.Name `xml:"TimeSpan"`
	Depart  string   `xml:"End"`
	Arrive  string   `xml:"Start"`
}

type HotelSearchCriteria struct {
	XMLName   xml.Name `xml:"HotelSearchCriteria"`
	Criterion Criterion
}

type Criterion struct {
	XMLName  xml.Name `xml:"Criterion"`
	HotelRef []HotelRef
	Address  Address
}

type HotelRef struct {
	XMLName       xml.Name `xml:"HotelRef,omitempty"`
	HotelCityCode string   `xml:",attr,omitempty"`
	HotelCode     string   `xml:",attr,omitempty"`
	//HotelName     string `xml:",attr,omitempty"`
	Latitude  string `xml:",attr,omitempty"`
	Longitude string `xml:",attr,omitempty"`
}

type GuestCounts struct {
	XMLName xml.Name `xml:"GuestCounts"`
	Count   int      `xml:",attr"`
}

type Customer struct {
	XMLName    xml.Name    `xml:"Customer,omitempty"`
	Corporate  *Corporate  //nil pointer ignored if empty
	CustomerID *CustomerID //nil pointer ignored if empty
}

type CustomerID struct {
	XMLName xml.Name `xml:"ID,omitempty"`
	Number  string   `xml:"Number,omitempty"`
}

type Corporate struct {
	XMLName xml.Name `xml:"Corporate,omitempty"`
	ID      string   `xml:"ID,omitempty"`
}

func (a *OTA_HotelAvailRQ) addCorporateID(cID string) {
	a.Avail.Customer = &Customer{
		Corporate: &Corporate{
			ID: cID,
		},
	}
}
func (a *OTA_HotelAvailRQ) addCustomerID(cID string) {
	a.Avail.Customer = &Customer{
		CustomerID: &CustomerID{
			Number: cID,
		},
	}
}

func BuildHotelAvailRq(guestCount int, query HotelSearchCriteria, arrive, depart time.Time) OTA_HotelAvailRQ {
	return OTA_HotelAvailRQ{
		Version:           hotelAvailVersion,
		XMLNS:             sbrweb.BaseWebServicesNS,
		XMLNSXs:           sbrweb.BaseXSDNameSpace,
		XMLNSXsi:          sbrweb.BaseXSINamespace,
		ReturnHostCommand: returnHostCommand,
		Avail: AvailRequestSegment{
			GuestCounts:         GuestCounts{Count: guestCount},
			HotelSearchCriteria: query,
			ArriveDepart: TimeSpan{
				Depart: depart.Format(timeSpanFormatter),
				Arrive: arrive.Format(timeSpanFormatter),
			},
		},
	}
}
