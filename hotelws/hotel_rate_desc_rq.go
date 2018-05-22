package hotelws

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/srvc"
)

// HotelRateDescRequest for soap package on HotelRateDescRequest service
type HotelRateDescRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   HotelRateDescBody
}

// HotelRateDescBody constructs soap body element
type HotelRateDescBody struct {
	XMLName         xml.Name `xml:"soap-env:Body"`
	HotelRateDescRQ HotelRateDescRQ
}

// HotelRateDescRQ retrieve sabre hotel content using various query criteria, see SearchCriteria
type HotelRateDescRQ struct {
	XMLName           xml.Name `xml:"HotelRateDescriptionRQ"`
	Version           string   `xml:"Version,attr"`
	XMLNS             string   `xml:"xmlns,attr"`
	XMLNSXs           string   `xml:"xmlns:xs,attr"`
	XMLNSXsi          string   `xml:"xmlns:xsi,attr"`
	ReturnHostCommand bool     `xml:"ReturnHostCommand,attr"`
	Avail             AvailRequestSegment
}

// SetHotelRateDescRqStruct hotel rate description request using input parameters
func SetHotelRateDescRqStruct(rpc *RatePlanCandidates) (HotelRateDescBody, error) {
	return HotelRateDescBody{
		HotelRateDescRQ: HotelRateDescRQ{
			Version:           hotelRQVersion,
			XMLNS:             srvc.BaseWebServicesNS,
			XMLNSXs:           srvc.BaseXSDNameSpace,
			XMLNSXsi:          srvc.BaseXSINamespace,
			ReturnHostCommand: true,
			Avail: AvailRequestSegment{
				RatePlanCandidates: rpc,
			},
		},
	}, nil
}

// BuildHotelRateDescRequest to make hotel property description request, done after hotel property description iff HRD_RequiredForSell=true.
func BuildHotelRateDescRequest(from, pcc, binsectoken, convid, mid, time string, propDesc HotelRateDescBody) HotelRateDescRequest {
	return HotelRateDescRequest{
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
				Service:        srvc.ServiceElem{Value: "HotelRateDescriptionLLSRQ", Type: "sabreXML"},
				Action:         "HotelRateDescriptionLLSRQ",
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
		Body: propDesc,
	}
}

// HotelRateDescriptionRS parse sabre hotel rate description
type HotelRateDescriptionRS struct {
	XMLName  xml.Name `xml:"HotelRateDescriptionRS"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSXs  string   `xml:"xs,attr"`
	XMLNSXsi string   `xml:"xsi,attr"`
	XMLNSStl string   `xml:"stl,attr"`
	Version  string   `xml:"Version,attr"`
	Result   ApplicationResults
	RoomStay RoomStay
}

// HotelRateDescResponse is wrapper with namespace prefix definitions for payload
type HotelRateDescResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		HotelDesc HotelRateDescriptionRS
		Fault     srvc.SOAPFault
	}
	ErrorSabreService ErrorSabreService
	ErrorSabreXML     ErrorSabreXML
}

// CallHotelRateDesc to sabre web services retrieve hotel rates using HotelRateDescriptionLLSRQ. This call only supports requests that contain an RPH from a previous hotel_property_desc call, see BuildHotelRateDescRequest.
func CallHotelRateDesc(serviceURL string, req HotelRateDescRequest) (HotelRateDescResponse, error) {
	propResp := HotelRateDescResponse{}
	byteReq, _ := xml.Marshal(req)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		propResp.ErrorSabreService = NewErrorSabreService(err.Error(), ErrCallHotelRateDesc, BadService)
		return propResp, propResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &propResp)
	if err != nil {
		propResp.ErrorSabreXML = NewErrorErrorSabreXML(err.Error(), ErrCallHotelRateDesc, BadParse)
		return propResp, propResp.ErrorSabreXML
	}
	return propResp, nil
}
