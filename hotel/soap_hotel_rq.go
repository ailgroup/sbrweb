/* Pakcage hotel implements Sabre hotel searching SOAP payloads through various criterion. Many criterion exist are not yet implemented: (Award, ContactNumbers, CommissionProgram, HotelAmenity, Package, PointOfInterest, PropertyType, RefPoint, RoomAmenity, HotelFeaturesCriterion,)
 */
package hotel

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ailgroup/sbrweb"
)

// QueryParams is a typed function to support optional query params on creation of new search criterion
type QueryParams func(*HotelSearchCriteria) error

// HotelRefCriterion map of hotel ref criteria
type HotelRefCriterion map[string][]string

// AddressCriterion map of address search criteria
type AddressCriterion map[string]string

// HotelAvailRequest for soap package on OTA_HotelAvailRQ service
type HotelAvailRequest struct {
	sbrweb.Envelope
	Header sbrweb.SessionHeader
	Body   HotelAvailBody
}

// HotelAvailBody constructs soap body element
type HotelAvailBody struct {
	XMLName         xml.Name `xml:"soap-env:Body"`
	OTAHotelAvailRQ OTAHotelAvailRQ
}

// OTAHotelAvailRQ retrieve sabre hotel content using various query criteria, see SearchCriteria
type OTAHotelAvailRQ struct {
	XMLName           xml.Name `xml:"OTA_HotelAvailRQ"`
	Version           string   `xml:"Version,attr"`
	XMLNS             string   `xml:"xmlns,attr"`
	XMLNSXs           string   `xml:"xmlns:xs,attr"`
	XMLNSXsi          string   `xml:"xmlns:xsi,attr"`
	ReturnHostCommand bool     `xml:"ReturnHostCommand,attr"`
	Avail             AvailRequestSegment
}

// HotelPropDescRequest for soap package on HotelPropertyDescriptionRQ service
type HotelPropDescRequest struct {
	sbrweb.Envelope
	Header sbrweb.SessionHeader
	Body   HotelPropDescBody
}

// HotelPropDescBody constructs soap body element
type HotelPropDescBody struct {
	XMLName         xml.Name `xml:"soap-env:Body"`
	HotelPropDescRQ HotelPropDescRQ
}

// HotelPropDescRQ retrieve sabre hotel content using various query criteria, see SearchCriteria
type HotelPropDescRQ struct {
	XMLName           xml.Name `xml:"HotelPropertyDescriptionRQ"`
	Version           string   `xml:"Version,attr"`
	XMLNS             string   `xml:"xmlns,attr"`
	XMLNSXs           string   `xml:"xmlns:xs,attr"`
	XMLNSXsi          string   `xml:"xmlns:xsi,attr"`
	ReturnHostCommand bool     `xml:"ReturnHostCommand,attr"`
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
	Depart  string   `xml:"End,attr"`
	Arrive  string   `xml:"Start,attr"`
}

// HotelSearchCriteria top level element for criterion
type HotelSearchCriteria struct {
	XMLName   xml.Name `xml:"HotelSearchCriteria"`
	Criterion Criterion
}

// Criterion holds various serach criteria
type Criterion struct {
	XMLName  xml.Name `xml:"Criterion"`
	HotelRef []*HotelRef
	Address  *Address
}

// HotelRef contains any number of serach criteria under the HotelRef element.
type HotelRef struct {
	XMLName       xml.Name `xml:"HotelRef,omitempty"`
	HotelCityCode string   `xml:"HotelCityCode,attr,omitempty"`
	HotelCode     string   `xml:"HotelCode,attr,omitempty"`
	//HotelName     string `xml:",attr,omitempty"`
	Latitude  string `xml:"Latitude,attr,omitempty"`
	Longitude string `xml:"Longitude,attr,omitempty"`
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

// Address represents typical building addresses
type Address struct {
	City        string `xml:"CityName,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
	Street      string `xml:"StreetNumber,omitempty"`
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

// AddressOption todo...
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
// Supports CityCode, HotelCode, Latitude, and Longitude for now... later support for HotelName, Latitude, Longitude.
func HotelRefSearch(params HotelRefCriterion) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		if len(params) < 1 {
			return fmt.Errorf("HotelRefCriterion params cannot be empty: %v", params)
		}
		/* limit property description requests to just hotel ids: FIND A BETTER WAY
		if property == true {
			for _, v := range params {
				for _, code := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, &HotelRef{HotelCode: code})
				}
			}
			return nil
		}
		*/
		for k, v := range params {
			switch k {
			case cityQueryField:
				for _, city := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, &HotelRef{HotelCityCode: city})
				}
			case hotelidQueryField:
				for _, code := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, &HotelRef{HotelCode: code})
				}
			case latlngQueryField:
				for _, l := range v {
					latlng := strings.Split(l, ",")
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, &HotelRef{Latitude: latlng[0], Longitude: latlng[1]})
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

// addCorporateID to the existing avail struct for a corporate customer
func (a *HotelPropDescRQ) addCorporateID(cID string) {
	a.Avail.Customer = &Customer{
		Corporate: &Corporate{
			ID: cID,
		},
	}
}

// addCustomerID rateID to the existing avail struct for a corporate customer
func (a *HotelPropDescRQ) addCustomerID(cID string) {
	a.Avail.Customer = &Customer{
		CustomerID: &CustomerID{
			Number: cID,
		},
	}
}

// arriveDepartParser parse string data value into time value.
func arriveDepartParser(arrive, depart string) (time.Time, time.Time) {
	a, _ := time.Parse(timeSpanFormatter, arrive)
	d, _ := time.Parse(timeSpanFormatter, depart)
	return a, d
}

// SetHotelAvailRqStruct hotel availability request using input parameters
func SetHotelAvailRqStruct(guestCount int, query HotelSearchCriteria, arrive, depart string) HotelAvailBody {
	a, d := arriveDepartParser(arrive, depart)
	return HotelAvailBody{
		OTAHotelAvailRQ: OTAHotelAvailRQ{
			Version:           hotelRQVersion,
			XMLNS:             sbrweb.BaseWebServicesNS,
			XMLNSXs:           sbrweb.BaseXSDNameSpace,
			XMLNSXsi:          sbrweb.BaseXSINamespace,
			ReturnHostCommand: returnHostCommand,
			Avail: AvailRequestSegment{
				GuestCounts:         GuestCounts{Count: guestCount},
				HotelSearchCriteria: query,
				ArriveDepart: TimeSpan{
					Depart: d.Format(timeSpanFormatter),
					Arrive: a.Format(timeSpanFormatter),
				},
			},
		},
	}
}

// BuildHotelAvailRequest to make hotel availability request.
func BuildHotelAvailRequest(from, pcc, binsectoken, convid, mid, time string, otaHotelAvail HotelAvailBody) HotelAvailRequest {
	return HotelAvailRequest{
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
				Service:        sbrweb.ServiceElem{Value: "OTA_HotelAvailRQ", Type: "sabreXML"},
				Action:         "OTA_HotelAvailLLSRQ",
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

type HotelAvailResponse struct {
}

// CallSessionValidate to sabre web services
func CallHotelAvail(serviceURL string, req HotelAvailRequest) (HotelAvailResponse, error) {
	availResp := HotelAvailResponse{}
	//construct payload
	byteReq, _ := xml.Marshal(req)
	fmt.Printf("\n\nREQUEST: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		return availResp, fmt.Errorf("CallHotelAvail http.Post(). %v", err)
	}

	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	fmt.Printf("\n\nBODYbuffer: %v\n\n", bodyBuffer)
	resp.Body.Close()

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &availResp)
	if err != nil {
		return availResp, fmt.Errorf("CallHotelAvail Unmarshal(bytes, &availResp): %v", err)
	}
	return availResp, nil
}

func (c *HotelSearchCriteria) validatePropertyRequest() error {
	for _, criterion := range c.Criterion.HotelRef {
		if len(criterion.HotelCityCode) > 0 {
			return fmt.Errorf("HotelCityCode not allowed in HotelPropertyDescription have %d for %v", len(criterion.HotelCityCode), criterion.HotelCityCode)
		}
		if len(criterion.Latitude) > 0 {
			return fmt.Errorf("Latitude not allowed in HotelPropertyDescription have %d for %v", len(criterion.Latitude), criterion.Latitude)
		}
		if len(criterion.Longitude) > 0 {
			return fmt.Errorf("Latitude not allowed in HotelPropertyDescription have %d for %v", len(criterion.Longitude), criterion.Longitude)
		}

		if len(c.Criterion.HotelRef) > 1 {
			return fmt.Errorf("Criterion.HotelRef cannot be greater than 1 for HotelPropertyDescription have %d for %v", len(c.Criterion.HotelRef), c.Criterion.HotelRef)
		}
	}
	return nil
}

// SetHotelPropDescRqStruct hotel availability request using input parameters
func SetHotelPropDescRqStruct(guestCount int, query HotelSearchCriteria, arrive, depart string) (HotelPropDescBody, error) {
	err := query.validatePropertyRequest()
	if err != nil {
		return HotelPropDescBody{}, err
	}
	a, d := arriveDepartParser(arrive, depart)
	return HotelPropDescBody{
		HotelPropDescRQ: HotelPropDescRQ{
			Version:           hotelRQVersion,
			XMLNS:             sbrweb.BaseWebServicesNS,
			XMLNSXs:           sbrweb.BaseXSDNameSpace,
			XMLNSXsi:          sbrweb.BaseXSINamespace,
			ReturnHostCommand: true,
			Avail: AvailRequestSegment{
				GuestCounts:         GuestCounts{Count: guestCount},
				HotelSearchCriteria: query,
				ArriveDepart: TimeSpan{
					Depart: d.Format(timeSpanFormatter),
					Arrive: a.Format(timeSpanFormatter),
				},
			},
		},
	}, nil
}

// BuildHotelPropDescRequest to make hotel property description request, which will have rate availability information on the response.
func BuildHotelPropDescRequest(from, pcc, binsectoken, convid, mid, time string, propDesc HotelPropDescBody) HotelPropDescRequest {
	return HotelPropDescRequest{
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
				Service:        sbrweb.ServiceElem{Value: "HotelPropertyDescription", Type: "sabreXML"},
				Action:         "HotelPropertyDescriptionLLSRQ",
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
		Body: propDesc,
	}
}

// CallHotelProp to sabre web services
func CallHotelProp(serviceURL string, req HotelPropDescRequest) error {

	byteReq, _ := xml.Marshal(req)
	fmt.Printf("\n\nREQUEST: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		return fmt.Errorf("CallHotelAvail http.Post(). %v", err)
	}
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	fmt.Printf("\n\nRESPONSE: %+v\n\n", resp)
	fmt.Printf("\n\nTLS: %+v\n\n", resp.Body)
	fmt.Printf("\n\nBODYbuffer: %v\n\n", bodyBuffer)
	return nil
}
