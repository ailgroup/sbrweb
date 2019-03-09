package itin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

/*
EndTransactionLLSRQ is used to complete and store changes made to a Passenger Name Record (PNR). For additional information please refer to the Format Finder Help System reference: pnrfo005.

    If finalizing or completing the PNR is the only action desired to occur please only include the .../EndTransaction node in the request.
    If receiving the PNR is the only action desired please send the .../UpdatedBy node in the request and set .../EndTransactionInd="false".

When updating an existing record please note that TravelItineraryReadRQ must be executed prior to calling EndTransactionLLSRQ.
*/

// EndTransactionBody holds namespaced body
type EndTransactionBody struct {
	XMLName          xml.Name `xml:"soap-env:Body"`
	EndTransactionRQ EndTransactionRQ
}

// EndTransactionRequest
type EndTransactionRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   EndTransactionBody
}

type EndTransaction struct {
	XMLName xml.Name `xml:"EndTransaction"`
	Ind     bool     `xml:"Ind,attr"`
}
type Source struct {
	XMLName      xml.Name `xml:"Source"`
	ReceivedFrom string   `xml:"ReceivedFrom,attr"`
}

// EndTransactionRQ root element
type EndTransactionRQ struct {
	XMLName        xml.Name `xml:"EndTransactionRQ"`
	XMLNS          string   `xml:"xmlns,attr"`
	XMLXS          string   `xml:"xmlns:xs,attr"`
	XMLXSI         string   `xml:"xmlns:xsi,attr"`
	Version        string   `xml:"Version,attr"`
	EndTransaction EndTransaction
	Source         Source
}

func BuildEndTransactionRequest(c *srvc.SessionConf) EndTransactionRequest {
	return EndTransactionRequest{
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
				Service:        srvc.ServiceElem{Value: "EndTransactionRQ", Type: "sabreXML"},
				Action:         "EndTransactionLLSRQ",
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
		Body: EndTransactionBody{
			EndTransactionRQ: EndTransactionRQ{
				XMLNS:   srvc.BaseWebServicesNS,
				XMLXS:   srvc.BaseXSDNameSpace,
				XMLXSI:  srvc.BaseXSINamespace,
				Version: "2.0.9",
				EndTransaction: EndTransaction{
					Ind: true,
				},
				Source: Source{
					ReceivedFrom: "IBE", //c.From,
				},
			},
		},
	}

}

type EndTransactionRS struct {
	XMLName      xml.Name `xml:"EndTransactionRS"`
	AppResults   ApplicationResults
	ItineraryRef ItineraryRef
}
type EndTransactionResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		EndTransactionRS EndTransactionRS
		Fault            srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallEndTransaction to execute EndTransactionRequest, which must be done in order to finish the booking transaction.
func CallEndTransaction(serviceURL string, req EndTransactionRequest) (EndTransactionResponse, error) {
	endT := EndTransactionResponse{}
	byteReq, _ := xml.Marshal(req)
	fmt.Printf("CallEndTransaction-REQUEST %s \n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		endT.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallEndTransaction,
			sbrerr.BadService,
		)
		return endT, endT.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// note ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, err = io.Copy(bodyBuffer, resp.Body)
	fmt.Printf("CallEndTransaction-RESPONSE %s \n\n", bodyBuffer)
	//close body no defer
	resp.Body.Close()
	//handle and return error if bad body
	if err != nil {
		endT.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallEndTransaction,
			sbrerr.BadParse,
		)
		return endT, endT.ErrorSabreService
	}

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &endT)
	if err != nil {
		endT.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallEndTransaction,
			sbrerr.BadParse,
		)
		return endT, endT.ErrorSabreXML
	}
	if !endT.Body.Fault.Ok() {
		return endT, sbrerr.NewErrorSoapFault(endT.Body.Fault.Format().ErrMessage)
	}

	if !endT.Body.EndTransactionRS.AppResults.Ok() {
		return endT, endT.Body.EndTransactionRS.AppResults.ErrFormat()
	}
	return endT, nil
}
