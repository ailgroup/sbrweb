/* Package hotel implements Sabre hotel searching SOAP payloads through various criterion. Many criterion exist that are not yet implemented: (Award, ContactNumbers, CommissionProgram, HotelAmenity, Package, PointOfInterest, PropertyType, RefPoint, RoomAmenity, HotelFeaturesCriterion,). To add more criterion create a criterion type (e.g, XCriterion) as well as its accomponying function to handle the data parms (e.g., XSearch).
 */
package hotel

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb"
)

type AvailabilityOptions struct {
	XMLName          xml.Name             `xml:"AvailabilityOptions"`
	AvailableOptions []AvailabilityOption `xml:"AvailabilityOption"`
}

type AvailabilityOption struct {
	RPH          int `xml:"RPH,attr"` //string? 001 versus 1
	PropertyInfo BasicPropertyInfo
}

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

// OTAHotelAvailRQ retrieve sabre hotel availability using various query criteria, see SearchCriteria
type OTAHotelAvailRQ struct {
	XMLName           xml.Name `xml:"OTA_HotelAvailRQ"`
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

// addCorporateID to the existing avail struct for a corporate customer
func (a *OTAHotelAvailRQ) addCorporateID(cID string) {
	a.Avail.Customer = &Customer{
		Corporate: &Corporate{
			ID: cID,
		},
	}
}

// addCustomerID rateID to the existing avail struct for a corporate customer
func (a *OTAHotelAvailRQ) addCustomerID(cID string) {
	a.Avail.Customer = &Customer{
		CustomerID: &CustomerID{
			Number: cID,
		},
	}
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

// OTAHotelAvailRS parse sabre hotel availability
type OTAHotelAvailRS struct {
	XMLName         xml.Name `xml:"OTA_HotelAvailRS"`
	XMLNS           string   `xml:"xmlns,attr"`
	XMLNSXs         string   `xml:"xs,attr"`
	XMLNSXsi        string   `xml:"xsi,attr"`
	XMLNSStl        string   `xml:"stl,attr"`
	Version         string   `xml:"Version,attr"`
	Result          ApplicationResults
	AdditionalAvail struct {
		Ind bool `xml:",attr"`
	} `xml:"AdditionalAvail,attr"`
	AvailOpts AvailabilityOptions
}

// HotelAvailResponse is wrapper with namespace prefix definitions for payload
type HotelAvailResponse struct {
	Envelope sbrweb.EnvelopeUnMarsh
	Header   sbrweb.SessionHeaderUnmarsh
	Body     struct {
		HotelAvail OTAHotelAvailRS
		Fault      sbrweb.SOAPFault
	}
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
