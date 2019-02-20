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
	OTA_CancelLLSRQ is used to cancel itinerary segments contained within a passenger name record (PNR).
	Please note that TravelItineraryReadLLSRQ must be executed prior to calling OTA_CancelLLSRQ.
*/

// CancelSegmentBody holds namespaced body
type CancelSegmentBody struct {
	XMLName         xml.Name `xml:"soap-env:Body"`
	CancelSegmentRQ CancelSegmentRQ
}

// CancelSeqmentRequest container for soap envelope, header, body
type CancelSeqmentRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   CancelSegmentBody
}

// SegmentToCancel is the PNR segment. Signals the type of segment to delete.
// Typ specifies type to cancel. Options: "air"(XIA), "vehicle"(XIC), "hotel"(XIH), "other"(XIO), or "entire"(XI).
// EndNumber is used to specify a range of segments to cancel.
// Number is used to specify a particular segment number to cancel.
// Typ cannot combine with EndNumber or Number.
type SegmentToCancel struct {
	XMLName   xml.Name `xml:"Segment"`
	Typ       string   `xml:"Type,attr,omitempty"`
	EndNumber string   `xml:"EndNumber,attr,omitempty"`
	Number    string   `xml:"Number,attr,omitempty"`
}

// CancelSegmentRQ root element
type CancelSegmentRQ struct {
	XMLName xml.Name `xml:"OTA_CancelRQ"`
	Version string   `xml:"Version,attr"`
	Segment SegmentToCancel
}

func BuildCancelSegmentRequest(c *srvc.SessionConf, typ, endnumber, number string) CancelSeqmentRequest {
	seg := SegmentToCancel{
		Typ:       typ,
		EndNumber: endnumber,
		Number:    number,
	}
	return CancelSeqmentRequest{
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
				Service:        srvc.ServiceElem{Value: "OTA_CancelRQ", Type: "sabreXML"},
				Action:         "OTA_CancelRQ", //OTA_CancelLLSRQ
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
		Body: CancelSegmentBody{
			CancelSegmentRQ: CancelSegmentRQ{
				Version: "2.0.2",
				Segment: seg,
			},
		},
	}
}

// GetReservationRS response schema for get reservations endpoint. Currently only concerned with Reservation payload.
type CancelSegmentRS struct {
	XMLName    xml.Name           `xml:"OTA_CancelRS"`
	Errors     PNRErrors          //is this available or needed?
	AppResults ApplicationResults //not sure service has this element
}

// GetReservationResponse container for soap envelope, header, body, and other errors.
type CancelSegmentResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		CancelSegmentRS CancelSegmentRS
		Fault           srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// Ok check for errors on get reservations requests.
func (r *CancelSegmentRS) Ok() bool {
	return len(r.Errors.Error) > 0
}

// CallGetReservation to execute GetReservationRequest, which must be done in order to finish the booking transaction.
func CallCancelSegment(serviceURL string, req GetReservationRequest) (CancelSegmentResponse, error) {
	cSeg := CancelSegmentResponse{}
	byteReq, _ := xml.Marshal(req)

	//-----------------------------------
	fmt.Printf("\n\n CallCancelSegment RAW REQUEST: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		cSeg.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadService,
		)
		return cSeg, cSeg.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// note ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, err = io.Copy(bodyBuffer, resp.Body)
	//close body no defer
	resp.Body.Close()
	//handle and return error if bad body
	if err != nil {
		cSeg.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadParse,
		)
		return cSeg, cSeg.ErrorSabreService
	}

	//-----------------------------------
	fmt.Printf("\n\n CallCancelSegment RAW RESPONSE: %s\n\n", bodyBuffer.Bytes())

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &cSeg)
	if err != nil {
		cSeg.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallGetReservation,
			sbrerr.BadParse,
		)
		return cSeg, cSeg.ErrorSabreXML
	}
	if !cSeg.Body.Fault.Ok() {
		return cSeg, sbrerr.NewErrorSoapFault(cSeg.Body.Fault.Format().ErrMessage)
	}

	// does this even return AppResults ??
	if !cSeg.Body.CancelSegmentRS.AppResults.Ok() {
		return cSeg, cSeg.Body.CancelSegmentRS.AppResults.ErrFormat()
	}

	if !cSeg.Body.CancelSegmentRS.Ok() {
		return cSeg, cSeg.Body.CancelSegmentRS.Errors.Format()
	}

	return cSeg, nil
}
