package htlsp

/*
	https://developer.sabre.com/docs/read/soap_apis/hotel/search/hotel_property_description

	The Hotel Property Description API (HOD) provides details on available room rates by room and rate type for a single property. Responses are based on real time requests to hotel suppliers with actual rates and rooms available at the time of request. The API allows the user to provide rate codes and qualifiers to shop for the applicable rates, and robust property descriptive content is provided with each rate.
*/

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

// HotelPropDescRequest for soap package on HotelPropertyDescriptionRQ service
type HotelPropDescRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
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

// SetHotelPropDescBody hotel availability request using input parameters
func SetHotelPropDescBody(guestCount int, query *HotelSearchCriteria, arrive, depart string) (HotelPropDescBody, error) {
	err := query.validatePropertyRequest()
	if err != nil {
		return HotelPropDescBody{}, err
	}
	ts := TimeSpanFormatter(arrive, depart, TimeFormatMD, TimeFormatMD)
	return HotelPropDescBody{
		HotelPropDescRQ: HotelPropDescRQ{
			Version:           hotelRQVersion,
			XMLNS:             srvc.BaseWebServicesNS,
			XMLNSXs:           srvc.BaseXSDNameSpace,
			XMLNSXsi:          srvc.BaseXSINamespace,
			ReturnHostCommand: true,
			Avail: AvailRequestSegment{
				GuestCounts:         &GuestCounts{Count: guestCount},
				HotelSearchCriteria: query,
				TimeSpan:            &ts,
			},
		},
	}, nil
}

// BuildHotelPropDescRequest to make hotel property description request, which will have rate availability information on the response.
func BuildHotelPropDescRequest(c *srvc.SessionConf, propDesc HotelPropDescBody) HotelPropDescRequest {
	return HotelPropDescRequest{
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
				Service:        srvc.ServiceElem{Value: "HotelPropertyDescription", Type: "sabreXML"},
				Action:         "HotelPropertyDescriptionLLSRQ",
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
		Body: propDesc,
	}
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
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		HotelDesc HotelPropertyDescriptionRS
		Fault     srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// SetRoomMetaData builds a b64 encoded string cache of rate request for later retrieval. See NewParsedRoomMeta for this data is parsed.
func (r *HotelPropDescResponse) SetRoomMetaData(guest int, arrive, depart, hotelid string) {
	for i, roomrate := range r.Body.HotelDesc.RoomStay.RoomRates {
		strslc := []string{}
		rrates := ""
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaArvKey, arrive))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaDptKey, depart))
		strslc = append(strslc, fmt.Sprintf("%s:%d", RoomMetaGstKey, guest))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaHcKey, r.Body.HotelDesc.Result.Success.System.HostCommand.Cryptic))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaHidKey, hotelid))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaRphKey, roomrate.RPH))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaRmtKey, roomrate.IATA_Character))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaGuarenteeKey, roomrate.GuaranteeSurcharge))

		for ri, rate := range roomrate.Rates {
			rrates += fmt.Sprintf("%s:%s-%s:%s-%s:%s", RrateMetaCurKey, rate.CurrencyCode, RrateMetaRqsKey, rate.HRD_RequiredForSell, RrateMetaAmtKey, rate.HotelPricing.Amount)
			if ((len(roomrate.Rates) - 1) - ri) != 0 {
				rrates += SColDelim
			}
		}
		strslc = append(strslc, fmt.Sprintf("%s%s%s", RBrackDelim, rrates, LBrackDelim))
		r.Body.HotelDesc.RoomStay.RoomRates[i].RoomToBook.B64RoomMetaData = B64Enc(strings.Join(strslc, PipeDelim))
	}
}

// CallHotelPropDesc to sabre web services retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
func CallHotelPropDesc(serviceURL string, req HotelPropDescRequest) (HotelPropDescResponse, error) {
	propResp := HotelPropDescResponse{}
	byteReq, _ := xml.Marshal(req)
	srvc.LogSoap.Printf("CallHotelPropDesc-REQUEST %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		propResp.ErrorSabreService = sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallHotelPropDesc, sbrerr.BadService)
		return propResp, propResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, _ = io.Copy(bodyBuffer, resp.Body)
	srvc.LogSoap.Printf("CallHotelPropDesc-RESPONSE %s\n\n", bodyBuffer)
	resp.Body.Close()

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &propResp)
	if err != nil {
		propResp.ErrorSabreXML = sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallHotelPropDesc, sbrerr.BadParse)
		return propResp, propResp.ErrorSabreXML
	}
	if !propResp.Body.Fault.Ok() {
		return propResp, propResp.Body.Fault.Format()
	}
	if !propResp.Body.HotelDesc.Result.Ok() {
		return propResp, propResp.Body.HotelDesc.Result.ErrFormat()
	}
	//propResp.SetTrackedEncode()
	return propResp, nil
}
