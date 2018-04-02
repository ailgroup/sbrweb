package sbrhotel

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/ailgroup/sbrweb"
)

// OTAHotelAvailRQ retrieve sabre hotel content using various query criteria, see SearchCriteria
type OTAHotelAvailRQ struct {
	XMLName           xml.Name `xml:"OTA_HotelAvailRQ"`
	Version           string   `xml:"version,attr"`
	XMLNS             string   `xml:"xmlns,attr"`
	XMLNSXs           string   `xml:"xmlns:xs,attr"`
	XMLNSXsi          string   `xml:"xmlns:xsi,attr"`
	ReturnHostCommand bool
	Avail             AvailRequestSegment
}

// AvailAvailRequestSegment holds basic hotel availability params: customer ids, guest count, hotel search criteria and arrival departure
type AvailRequestSegment struct {
	XMLName             xml.Name  `xml:"AvailRequestSegment"`
	Customer            *Customer //nil pointer ignored if empty
	GuestCounts         GuestCounts
	HotelSearchCriteria HotelSearchCriteria
	ArriveDepart        TimeSpan `xml:"TimeSpan"`
}

// Timepsan for arrival and departure params
type TimeSpan struct {
	XMLName xml.Name `xml:"TimeSpan"`
	Depart  string   `xml:"End"`
	Arrive  string   `xml:"Start"`
}

// HotelSearchCriteria top level element for criterion
type HotelSearchCriteria struct {
	XMLName   xml.Name `xml:"HotelSearchCriteria"`
	Criterion Criterion
}

// Criterion holds various serach criteria
type Criterion struct {
	XMLName  xml.Name `xml:"Criterion"`
	HotelRef []HotelRef
	Address  Address
}

// HotelRef contains any number of serach criteria under the HotelRef element.
type HotelRef struct {
	XMLName       xml.Name `xml:"HotelRef,omitempty"`
	HotelCityCode string   `xml:",attr,omitempty"`
	HotelCode     string   `xml:",attr,omitempty"`
	//HotelName     string `xml:",attr,omitempty"`
	Latitude  string `xml:",attr,omitempty"`
	Longitude string `xml:",attr,omitempty"`
}

// GuestCounts how many guests per night-room. TODO: check on Sabre validation limits (think it is 4)
type GuestCounts struct {
	XMLName xml.Name `xml:"GuestCounts"`
	Count   int      `xml:",attr"`
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
	// NEXT on the list:
	tccepype HotelFeaturesCriterion map[string][]string
*/

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

type AddressCriterion map[string]string

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

type HotelRefCriterion map[string][]string

// HotelRefSearch accepts HotelRef criterion and returns a function for hotel search critera.
// Supports CityCode, HotelCode, Latitude, and Longitude for now... later support for HotelName, Latitude, Longitude.
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

// addCorporateID to the existing avail struct for a corporate customer
func (a *OTAHotelAvailRQ) addCorporateID(cID string) {
	a.Avail.Customer = &Customer{
		Corporate: &Corporate{
			ID: cID,
		},
	}
}

// addCorpoaddCustomerID rateID to the existing avail struct for a corporate customer
func (a *OTAHotelAvailRQ) addCustomerID(cID string) {
	a.Avail.Customer = &Customer{
		CustomerID: &CustomerID{
			Number: cID,
		},
	}
}

// SetHotelAvailRqStruct hotel availability request using input parameters
func SetHotelAvailRqStruct(guestCount int, query HotelSearchCriteria, arrive, depart time.Time) OTAHotelAvailRQ {
	return OTAHotelAvailRQ{
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

// AvailCreateRequest is wrapper with namespace prefix definitions for payload
type HotelAvailCreateRequest struct {
	sbrweb.Envelope
	Header sbrweb.SessionHeader
	Body   OTAHotelAvailRQ
}

// BuildHotelAvailRequest to make hotel availability request.
func BuildHotelAvailRequest(from, pcc, binsectoken, convid, mid, time string, otaHotelAvail OTAHotelAvailRQ) HotelAvailCreateRequest {
	return HotelAvailCreateRequest{
		Envelope: sbrweb.CreateEnvelope(),
		Header: sbrweb.SessionHeader{
			MessageHeader: sbrweb.MessageHeader{
				MustUnderstand: sbrweb.SabreMustUnderstand,
				EbVersion:      sbrweb.SabreEBVersion,
				From: sbrweb.FromElem{
					PartyID: sbrweb.CreatePartyID(from, sbrweb.PartyIDTypeURN),
				},
				To: sbrweb.ToElem{
					PartyID: sbrweb.CreatePartyID(sbrweb.SabreToBase, sbrweb.PartyIDTypeURN),
				},
				CPAID:          pcc,
				ConversationID: convid,
				Service:        sbrweb.ServiceElem{Value: "OTA_HotelAvailRQ", Type: "OTA"},
				Action:         "OTA_HotelAvailRQ",
				MessageData: sbrweb.MessageDataElem{
					MessageID: mid,
					Timestamp: time,
				},
			},
			Security: sbrweb.Security{
				XMLNSWsseBase:       sbrweb.BaseWsse,
				XMLNSWsu:            sbrweb.BaseWsuNameSpace,
				BinarySecurityToken: binsectoken,
			},
		},
		Body: otaHotelAvail,
	}
}

/*
// CallSessionCreate to sabre web services.
func CallSessionCreate(serviceURL string, req SessionCreateRequest) (SessionCreateResponse, error) {
	sessionResponse := SessionCreateResponse{}
	//construct payload

	byteReq, _ := xml.Marshal(req)
	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		return sessionResponse, fmt.Errorf("CallSessionCreate http.Post(). %v", err)
	}

	//parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	//defer func() { resp.Body.Close() }()
	resp.Body.Close()

	//marshal byte body sabre response body into session envelope response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &sessionResponse)
	if err != nil {
		return sessionResponse, fmt.Errorf("CallSessionCreate Unmarshal(,&sessionResponse). %v", err)
	}
	return sessionResponse, nil
}
*/
