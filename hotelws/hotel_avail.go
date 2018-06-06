package hotelws

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/srvc"
)

// HotelAvailRequest for soap package on OTA_HotelAvailRQ service
type HotelAvailRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
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

// SetHotelAvailBody hotel availability request using input parameters
func SetHotelAvailBody(guestCount int, query *HotelSearchCriteria, arrive, depart string) HotelAvailBody {
	a, d := arriveDepartParser(arrive, depart)
	return HotelAvailBody{
		OTAHotelAvailRQ: OTAHotelAvailRQ{
			Version:           hotelRQVersion,
			XMLNS:             srvc.BaseWebServicesNS,
			XMLNSXs:           srvc.BaseXSDNameSpace,
			XMLNSXsi:          srvc.BaseXSINamespace,
			ReturnHostCommand: returnHostCommand,
			Avail: AvailRequestSegment{
				GuestCounts:         &GuestCounts{Count: guestCount},
				HotelSearchCriteria: query,
				TimeSpan: &TimeSpan{
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
		Envelope: srvc.CreateEnvelope(),
		Header: srvc.SessionHeader{
			MessageHeader: srvc.MessageHeader{
				MustUnderstand: srvc.SabreMustUnderstand,
				EbVersion:      srvc.SabreEBVersion,
				From: srvc.FromElem{
					PartyID: srvc.CreatePartyID(from, srvc.PartyIDTypeURN),
				},
				To: srvc.ToElem{
					PartyID: srvc.CreatePartyID(srvc.SabreToBase, srvc.PartyIDTypeURN),
				},
				CPAID:          pcc,
				ConversationID: convid,
				Service:        srvc.ServiceElem{Value: "OTA_HotelAvailRQ", Type: "sabreXML"},
				Action:         "OTA_HotelAvailLLSRQ",
				MessageData: srvc.MessageDataElem{
					MessageID: mid,
					Timestamp: time,
				},
			},
			Security: srvc.Security{
				XMLNSWsseBase:       srvc.BaseWsse,
				XMLNSWsu:            srvc.BaseWsuNameSpace,
				BinarySecurityToken: binsectoken,
			},
		},
		Body: otaHotelAvail,
	}
}

type AvailabilityOptions struct {
	XMLName          xml.Name             `xml:"AvailabilityOptions"`
	AvailableOptions []AvailabilityOption `xml:"AvailabilityOption"`
}

type AvailabilityOption struct {
	RPH          int `xml:"RPH,attr"` //string? 001 versus 1
	PropertyInfo BasicPropertyInfo
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

// HotelAvailResponse for parsing hote availability response
type HotelAvailResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		HotelAvail OTAHotelAvailRS
		Fault      srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallHotelAvail to sabre web services
func CallHotelAvail(serviceURL string, req HotelAvailRequest) (HotelAvailResponse, error) {
	availResp := HotelAvailResponse{}
	//construct payload
	byteReq, _ := xml.Marshal(req)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		availResp.ErrorSabreService = sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadService)
		return availResp, availResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &availResp)
	if err != nil {
		availResp.ErrorSabreXML = sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadParse)
		return availResp, availResp.ErrorSabreXML
	}
	if !availResp.Body.Fault.Ok() {
		return availResp, availResp.Body.Fault.Format()
	}
	if !availResp.Body.HotelAvail.Result.Ok() {
		return availResp, availResp.Body.HotelAvail.Result.ErrFormat()
	}
	return availResp, nil
}
