package hotelws

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb"
)

var (
	ErrPropDescCityCode  = errors.New("HotelCityCode not allowed in HotelPropertyDescription")
	ErrPropDescLatLng    = errors.New("Latitude or Longitude not allowed in HotelPropertyDescription")
	ErrPropDescHotelRefs = errors.New("Criterion.HotelRef cannot be greater than 1, can only search using one criterion")
)

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

// validatePropertyRequest ensures property description requests are well-formed
func (c *HotelSearchCriteria) validatePropertyRequest() error {
	for _, criterion := range c.Criterion.HotelRefs {
		if len(criterion.HotelCityCode) > 0 {
			return ErrPropDescCityCode
		}
		if (len(criterion.Latitude) > 0) || (len(criterion.Longitude) > 0) {
			return ErrPropDescLatLng
		}

		if len(c.Criterion.HotelRefs) > 1 {
			return ErrPropDescHotelRefs
		}
	}
	return nil
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

type RoomStay struct {
	XMLName           xml.Name `xml:"RoomStay"`
	BasicPropertyInfo BasicPropertyInfo
	RoomRates         []RoomRate `xml:"RoomRates>RoomRate"`
	TimeSpan          struct {
		Duration int    `xml:"Duration,attr"` //string 0001 or int 1?
		End      string `xml:"End,attr"`
		Start    string `xml:"Start,attr"`
	} `xml:"TimeSpan"`
}

// OTAHotelAvailRS parse sabre hotel availability
type HotelPropertyDescriptionRS struct {
	XMLName  xml.Name `xml:"HotelPropertyDescriptionRS"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSXs  string   `xml:"xs,attr"`
	XMLNSXsi string   `xml:"xsi,attr"`
	XMLNSStl string   `xml:"stl,attr"`
	Version  string   `xml:"Version,attr"`
	Result   ApplicationResults
	RoomStay RoomStay
}

// HotelAvailResponse is wrapper with namespace prefix definitions for payload
type HotelPropDescResponse struct {
	Envelope sbrweb.EnvelopeUnMarsh
	Header   sbrweb.SessionHeaderUnmarsh
	Body     struct {
		HotelDesc HotelPropertyDescriptionRS
		Fault     sbrweb.SOAPFault
	}
}

// CallHotelProperty to sabre web services
func CallHotelProperty(serviceURL string, req HotelPropDescRequest) error {
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
