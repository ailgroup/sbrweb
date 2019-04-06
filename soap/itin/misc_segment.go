package itin

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

/*
The Sell Miscellaneous Segment (MiscSegmentSellLLSRQ) API is used to sell OTH, MCO, PTA, or INS miscellaneous segment types.

Segment types include the following:
    MCO – Miscellaneous Charge Order
    OTH – Other Product Segment
    PTA – Prepaid Ticket Advice
    INS – Insurance Segment

Using this API, you can:
    Create miscellaneous segment with a carrier code.
    Create miscellaneous segment with a vendor code.
    Create a miscellaneous segment using free text instead of three alpha character city/airport code.
    Create a miscellaneous hotel segment for a location without a city code.

Please note that .../Text can contain up to a maximum of 215 letters, numbers, commas and spaces.

When updating an existing record with a miscellaneous segment please note that TravelItineraryReadLLSRQ/GetItinerary must be executed prior to calling MiscSegmentSellLLSRQ.

https://beta.developer.sabre.com/docs/soap_apis/air/book/Sell_Miscellaneous_Segments
*/

// MiscSegmentBody holds namespaced body
type MiscSegmentBody struct {
	XMLName       xml.Name `xml:"soap-env:Body"`
	MiscSegmentRQ MiscSegmentRQ
}

// MiscSegmentRequest wrapper for soap payload.
type MiscSegmentRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   MiscSegmentBody
}

// MiscSegmentRQ root element
type MiscSegmentRQ struct {
	XMLName xml.Name `xml:"MiscSegmentSellRQ"`
	XMLNS   string   `xml:"xmlns,attr"`
	XMLXSI  string   `xml:"xmlns:xsi,attr"`
	//XSISchema string   `xml:"xsi:schemaLocation,attr"`
	ReturnHostCommand bool `xml:"ReturnHostCommand,attr"`
	//Timestamp         string `xml:"Timestamp,attr"`
	Version     string `xml:"Version,attr"`
	MiscSegment MiscSegment
}

type MiscSegment struct {
	XMLName           xml.Name `xml:"MiscSegment"`
	DepartureDateTime string   `xml:"DepartureDateTime,attr"` // MM-DD
	NumberInParty     int32    `xml:"NumberInParty,attr"`
	Status            string   `xml:"Status,attr"`
	Typ               string   `xml:"Type,attr"`
	//InsertAfter       int32    `xml:"InsertAfter,attr"`
	OriginLocation OriginLocation
	MiscSegText    MiscSegText
	VendorPrefs    VendorPrefs
}
type OriginLocation struct {
	XMLName      xml.Name `xml:"OriginLocation"`
	LocationCode string   `xml:"LocationCode,attr,omitempty"`
}
type MiscSegText struct {
	XMLName xml.Name `xml:"Text"`
	Val     string   `xml:",chardata"`
}

// BuildMiscSegmentRequest construct payload for request to copy a profile into a PNR. The FilterPath is parameter since it may need to be constructed in various ways depending on need; for example BuildFilterPathForProfileOnly buils a simple FilterPath Profile
func BuildMiscSegmentRequest(c *srvc.SessionConf, binsec string, seg MiscSegment) MiscSegmentRequest {
	return MiscSegmentRequest{
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
				Service:        srvc.ServiceElem{Value: "MiscSegmentSellLLSRQ", Type: "sabreXML"},
				Action:         "MiscSegmentSellLLSRQ",
				MessageData: srvc.MessageDataElem{
					MessageID: srvc.GenerateMessageID(),
					Timestamp: srvc.SabreTimeNowFmt(),
				},
			},
			Security: srvc.Security{
				XMLNSWsseBase:       srvc.BaseWsse,
				XMLNSWsu:            srvc.BaseWsuNameSpace,
				BinarySecurityToken: binsec,
			},
		},
		Body: MiscSegmentBody{
			MiscSegmentRQ: MiscSegmentRQ{
				XMLNS:  srvc.BaseWebServicesNS,
				XMLXSI: srvc.BaseTPFCSchema,
				//XSISchema: srvc.BaseMisSegSchema,
				ReturnHostCommand: true,
				//Timestamp:         srvc.SabreTimeNowFmt(),
				Version:     "2.0.0",
				MiscSegment: seg,
			},
		},
	}
}

type MiscSegmentRS struct {
	XMLName    xml.Name `xml:"MiscSegmentSellRS"`
	AppResults ApplicationResults
	Text       MiscSegText
}

type MiscSegmentResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		MiscSegmentRS MiscSegmentRS
		Fault         srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallMiscSegment to execute MiscSegmentRequest, which is done in order to add more segments to existing PNR.
func CallMiscSegment(serviceURL string, req MiscSegmentRequest) (MiscSegmentResponse, error) {
	miscS := MiscSegmentResponse{}
	byteReq, _ := xml.Marshal(req)
	srvc.LogSoap.Printf("CallMiscSegment-REQUEST %s \n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		miscS.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallMiscSegment,
			sbrerr.BadService,
		)
		return miscS, miscS.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// note ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, err = io.Copy(bodyBuffer, resp.Body)
	srvc.LogSoap.Printf("CallMiscSegment-RESPONSE %s \n\n", bodyBuffer)
	//close body no defer
	resp.Body.Close()
	//handle and return error if bad body
	if err != nil {
		miscS.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallMiscSegment,
			sbrerr.BadParse,
		)
		srvc.LogSoap.Printf("CallMiscSegment-Unmarshal %v \n\n", miscS.ErrorSabreService)
		return miscS, miscS.ErrorSabreService
	}

	//marshal bytes sabre response body into miscS response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &miscS)
	if err != nil {
		miscS.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallMiscSegment,
			sbrerr.BadParse,
		)
		srvc.LogSoap.Printf("CallMiscSegment-Unmarshal %v \n\n", miscS.ErrorSabreXML)
		return miscS, miscS.ErrorSabreXML
	}
	if !miscS.Body.Fault.Ok() {
		srvc.LogSoap.Printf("CallMiscSegment-Fault %v \n\n", sbrerr.NewErrorSoapFault(miscS.Body.Fault.Format().ErrMessage))
		return miscS, sbrerr.NewErrorSoapFault(miscS.Body.Fault.Format().ErrMessage)
	}

	if !miscS.Body.MiscSegmentRS.AppResults.Ok() {
		srvc.LogSoap.Printf("CallMiscSegment-AppResults %v \n\n", miscS.Body.MiscSegmentRS.AppResults)
		return miscS, miscS.Body.MiscSegmentRS.AppResults.ErrFormat()
	}
	return miscS, nil
}
