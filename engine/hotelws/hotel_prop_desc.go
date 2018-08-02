package hotelws

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
	"github.com/ailgroup/sbrweb/engine/srvc"
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

func (r *HotelPropDescResponse) SetRoomMetaData() {
	for i, rate := range r.Body.HotelDesc.RoomStay.RoomRates {
		strslc := []string{}
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaRPHKey, rate.RPH))
		strslc = append(strslc, fmt.Sprintf("%s:%s", RoomMetaIATACharKey, rate.IATA_Character))
		for ri, rr := range rate.Rates {
			strslc = append(strslc, fmt.Sprintf("%s:%d-%s:%s-%s:%s", RoomMetaRatesIdxKey, ri, RoomMetaTotalKey, rr.HotelPricing.Amount, RoomMetaRateNextKey, rr.HRD_RequiredForSell))
		}
		r.Body.HotelDesc.RoomStay.RoomRates[i].B64RoomMetaData = B64Enc(strings.Join(strslc, RoomMetaDelimiter))
	}
}

// CallHotelPropDesc to sabre web services retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
func CallHotelPropDesc(serviceURL string, req HotelPropDescRequest) (HotelPropDescResponse, error) {
	propResp := HotelPropDescResponse{}
	byteReq, _ := xml.Marshal(req)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		propResp.ErrorSabreService = sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallHotelPropDesc, sbrerr.BadService)
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
