package itin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
	"github.com/ailgroup/sbrweb/engine/srvc"
)

/*
	GetReservationRQ Retrieve Itinerary API is used to retrieve and display a passenger name record (PNR) and data that is related to PNR..

	Once a PNR has been created on the Sabre Host, this Web Service offers capabilities allowing Airline or Agency to retrieve PNR data using PNR Locator as a search criterion. It also enables retrieving PNR from AAA session. Request payload can be further specified by using "ReturnOptions" which determine response message content.

	For Read Only Access use the Trip option as the PNR is not unpacked into the user AAA Session. The PNR Locator must always be specified in the request.

	For Update Access use the Stateful option as this will unpack the PNR into the user AAA session and be available for follow up Sabre entries. If a Locator is specified in the request the service checks the Locator in the AAA and if they match retrieves the current data in the AAA, if they do not match the service will unpack the PNR into the AAA session as long as the current session is available and there are no updates outstanding.
*/

// GetReservationBody holds namespaced body
type GetReservationBody struct {
	XMLName          xml.Name `xml:"soap-env:Body"`
	GetReservationRQ GetReservationRQ
}

// GetReservationRequest container for soap envelope, header, body
type GetReservationRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   GetReservationBody
}

// Locator PNR search criterion retrieve current data into AAA workspace
type Locator struct {
	XMLName xml.Name `xml:"Locator"`
	Val     string   `xml:",chardata"`
}

// RequestType signals hwo to interact with workspace session: options are Stateful (updates), Stateless (read only), Trip (read only, no data loaded into workspace).
type RequestType struct {
	XMLName xml.Name `xml:"RequestType"`
	Val     string   `xml:",chardata"`
}

/*
	// SubjectArea determines type of response, options include PRICE_QUOTE
	type SubjectArea struct {
		XMLName xml.Name `xml:"SubjectArea"`
		Val     string   `xml:",chardata"`
	}

	// ViewName
	type ViewName struct {
		XMLName xml.Name `xml:"ViewName"`
		Val     string   `xml:",chardata"`
	}

	// ResponseFormat
	type ResponseFormat struct {
		XMLName xml.Name `xml:"ResponseFormat"`
		Val     string   `xml:",chardata"`
	}

	// ReturnOptions allows client to specify various ways of reading out PNR segments
	type ReturnOptions struct {
		XMLName        xml.Name      `xml:"ReturnOptions"`
		PQVersion      string        `xml:"PriceQuoteServiceVersion,attr"`
		SubjectAreas   []SubjectArea `xml:"SubjectAreas>SubjectArea"`
		ViewName       ViewName
		ResponseFormat ResponseFormat
	}
*/

// GetReservationRQ root element
type GetReservationRQ struct {
	XMLName     xml.Name `xml:"GetReservationRQ"`
	XMLNS       string   `xml:"xmlns,attr"` //"http://webservices.sabre.com/pnrbuilder/v1_19"
	Version     string   `xml:"Version,attr"`
	Locator     Locator
	RequestType RequestType
	//ReturnOptions ReturnOptions //not necessary but leaving here for documentation
}

func BuildGetReservationRequest(c *srvc.SessionConf, locator string) GetReservationRequest {
	return GetReservationRequest{
		Envelope: srvc.CreateEnvelope(),
		Header: srvc.SessionHeader{
			MessageHeader: srvc.MessageHeader{
				MustUnderstand: srvc.SabreMustUnderstand,
				EbVersion:      srvc.SabreEBVersion,
				From: srvc.FromElem{
					PartyID: srvc.CreatePartyID(c.From, srvc.PartyIDTypeURN),
				},
				To: srvc.ToElem{
					PartyID: srvc.CreatePartyID(srvc.SabreToBase, srvc.PartyIDTypeURN),
				},
				CPAID:          c.PCC,
				ConversationID: c.Convid,
				Service:        srvc.ServiceElem{Value: "GetReservationRQ", Type: "sabreXML"},
				Action:         "GetReservationRQ",
				MessageData: srvc.MessageDataElem{
					MessageID: c.Msgid,
					Timestamp: c.Timestr,
				},
			},
			Security: srvc.Security{
				XMLNSWsseBase:       srvc.BaseWsse,
				XMLNSWsu:            srvc.BaseWsuNameSpace,
				BinarySecurityToken: c.Binsectok,
			},
		},
		Body: GetReservationBody{
			GetReservationRQ: GetReservationRQ{
				XMLNS:   "http://webservices.sabre.com/pnrbuilder/v1_19",
				Version: "1.19.0",
				Locator: Locator{
					Val: locator,
				},
				RequestType: RequestType{
					Val: "Stateful",
					//Val: "Stateless",
					//Val: "Trip",
				},
				/*
					ReturnOptions: ReturnOptions{
						PQVersion:      "3.2.0",
						SubjectAreas:   []SubjectArea{SubjectArea{Val: subject}},
						ViewName:       ViewName{Val: "Simple"},
						ResponseFormat: ResponseFormat{Val: "OTA"},
					},
				*/
			},
		},
	}
}

// PNRError capture and format errors related to get reservation, pnr, and segments.
type PNRError struct {
	XMLName xml.Name `xml:"Error"`
	Code    struct {
		V string `xml:",chardata"`
	} `xml:"Code"`
	Message struct {
		V string `xml:",chardata"`
	} `xml:"Message"`
	Severity struct {
		V string `xml:",chardata"`
	} `xml:"Severity"`
}

// PNRErrors format for list of PNRError elements.
type PNRErrors struct {
	XMLName xml.Name `xml:"Errors"`
	Error   []PNRError
}

// GetReservationRS response schema for get reservations endpoint. Currently only concerned with Reservation payload.
type GetReservationRS struct {
	XMLName xml.Name `xml:"GetReservationRS"`
	Errors  PNRErrors
	//AppResults  ApplicationResults //not sure service has this element
	Reservation Reservation
	//PriceQuote  PriceQuote //need for response format request SubjectAreas=PRICE_QUOTE
}

// GetReservationResponse container for soap envelope, header, body, and other errors.
type GetReservationResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		GetReservationRS GetReservationRS
		Fault            srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// Ok check for errors on get reservations requests.
func (r *GetReservationRS) Ok() bool {
	return len(r.Errors.Error) > 0
}

// Format for PNRErrors returns standard formatted erorrs for api.
func (p *PNRErrors) Format() sbrerr.ErrorSabreService {
	var appmsg string
	var errmsg string
	for _, e := range p.Error {
		errmsg += e.Message.V
		appmsg += fmt.Sprintf("%s %s", e.Code.V, e.Severity.V)
	}
	return sbrerr.ErrorSabreService{
		Code:       sbrerr.SabreEngineStatusCode(sbrerr.StatusNotProcess()),
		ErrMessage: errmsg,
		AppMessage: appmsg,
	}
}

// CallGetReservation to execute GetReservationRequest, which must be done in order to finish the booking transaction.
func CallGetReservation(serviceURL string, req GetReservationRequest) (GetReservationResponse, error) {
	getRes := GetReservationResponse{}
	byteReq, _ := xml.Marshal(req)

	//-----------------------------------
	fmt.Printf("\n\nCallGetReservation RAW REQUEST: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		getRes.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadService,
		)
		return getRes, getRes.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// note ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, err = io.Copy(bodyBuffer, resp.Body)
	//close body no defer
	resp.Body.Close()
	//handle and return error if bad body
	if err != nil {
		getRes.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadParse,
		)
		return getRes, getRes.ErrorSabreService
	}

	//-----------------------------------
	fmt.Printf("\n\nCallGetReservation RAW RESPONSE: %s\n\n", bodyBuffer.Bytes())

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &getRes)
	if err != nil {
		getRes.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadParse,
		)
		return getRes, getRes.ErrorSabreXML
	}
	if !getRes.Body.Fault.Ok() {
		return getRes, sbrerr.NewErrorSoapFault(getRes.Body.Fault.Format().ErrMessage)
	}

	// does this even return AppResults ??
	//if !getRes.Body.GetReservationRS.AppResults.Ok() {
	//	return getRes, getRes.Body.GetReservationRS.AppResults.ErrFormat()
	//}

	if !getRes.Body.GetReservationRS.Ok() {
		return getRes, getRes.Body.GetReservationRS.Errors.Format()
	}

	return getRes, nil
}
