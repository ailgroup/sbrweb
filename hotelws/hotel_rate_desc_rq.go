package hotelws

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb"
)

// HotelRateDescRequest for soap package on HotelRateertyDescriptionRQ service
type HotelRateDescRequest struct {
	sbrweb.Envelope
	Header sbrweb.SessionHeader
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

// addCorporateID to the existing avail struct for a corporate customer
func (a *HotelRateDescRQ) addCorporateID(cID string) {
	a.Avail.Customer = &Customer{
		Corporate: &Corporate{
			ID: cID,
		},
	}
}

// addCustomerID rateID to the existing avail struct for a corporate customer
func (a *HotelRateDescRQ) addCustomerID(cID string) {
	a.Avail.Customer = &Customer{
		CustomerID: &CustomerID{
			Number: cID,
		},
	}
}

// SetHotelRateDescRqStruct hotel availability request using input parameters
func SetHotelRateDescRqStruct(guestCount int, query HotelSearchCriteria, arrive, depart string) (HotelRateDescBody, error) {
	err := query.validatePropertyRequest()
	if err != nil {
		return HotelRateDescBody{}, err
	}
	a, d := arriveDepartParser(arrive, depart)
	return HotelRateDescBody{
		HotelRateDescRQ: HotelRateDescRQ{
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

// BuildHotelRateDescRequest to make hotel property description request, which will have rate availability information on the response.
func BuildHotelRateDescRequest(from, pcc, binsectoken, convid, mid, time string, propDesc HotelRateDescBody) HotelRateDescRequest {
	return HotelRateDescRequest{
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
				Service:        sbrweb.ServiceElem{Value: "HotelRateertyDescription", Type: "sabreXML"},
				Action:         "HotelRateertyDescriptionLLSRQ",
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

// HotelRateDescriptionRS parse sabre hotel availability
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

// HotelAvailResponse is wrapper with namespace prefix definitions for payload
type HotelRateDescResponse struct {
	Envelope sbrweb.EnvelopeUnMarsh
	Header   sbrweb.SessionHeaderUnmarsh
	Body     struct {
		HotelDesc HotelRateDescriptionRS
		Fault     sbrweb.SOAPFault
	}
	ErrorSabreService ErrorSabreService
	ErrorSabreXML     ErrorSabreXML
}

// CallHotelRate to sabre web services retrieve hotel rates using HotelRateertyDescriptionLLSRQ.
func CallHotelRate(serviceURL string, req HotelRateDescRequest) (HotelRateDescResponse, error) {
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
