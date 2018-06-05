package hotelws

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/itin"
	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/srvc"
)

/*
Using this API, you can book a hotel and specify:

    Confirmation number
    Passenger name number association
    Frequent flyer number
    Corporate ID number
    An ID number
    Number of cribs
    Number of extra guests
    Number of rollaways
    The booking source
    Special request-related information
    Plus, several other parameters
*/

// HotelRsrvRequest for soap package on OTA_HotelResRQ service for making reservations
type HotelRsrvRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   HotelRsrvBody
}

// HotelRsrvBody implements Hotel element for SOAP
type HotelRsrvBody struct {
	XMLName       xml.Name `xml:"soap-env:Body"`
	OTAHotelResRQ OTAHotelResRQ
}

// OTAHotelResRQ holds hotel information specific to making a hote reservation
type OTAHotelResRQ struct {
	XMLName  xml.Name `xml:"OTA_HotelResRQ"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSXs  string   `xml:"xmlns:xs,attr"`
	XMLNSXsi string   `xml:"xmlns:xsi,attr"`
	Version  string   `xml:"Version,attr"`
	Hotel    HotelRequest
}

// HotelForRQ Hotle segment for reservations
type HotelRequest struct {
	XMLName          xml.Name `xml:"Hotel"`
	BasicPropertyRes BasicPropertyRes
	Guarantee        GuaranteeReservation
	RoomType         *RoomType
	SpecialPrefs     *SpecialPrefs
}

// BasicPropertyRes is the BasicPropertyInfo element specifically for executing hotel reservations. Easier to duplicate this simple case than omit all the struct fields in the BasicPropertyInfo type.
type BasicPropertyRes struct {
	XMLName xml.Name `xml:"BasicPropertyInfo"`
	RPH     int      `xml:"RPH,attr"`
}

// CCInfo for passing credit card
type CCInfo struct {
	XMLName     xml.Name `xml:"CC_Info"`
	PaymentCard PaymentCard
	PersonName  itin.PersonName
}

// GuaranteeReservation is a gurantee type specifically for executing hotel reservations
type GuaranteeReservation struct {
	XMLName xml.Name `xml:"Guarantee"`
	Type    string   `xml:"Type,attr"`
	CCInfo  CCInfo
}

type RoomType struct {
	XMLName       xml.Name `xml:"RoomType"`
	NumberOfUnits int      `xml:"NumberOfUnits,attr"`
}

type WrittenConfirmation struct {
	XMLName xml.Name `xml:"WrittenConfirmation"`
	Ind     bool     `xml:"Ind,attr"`
}

type SpecialPrefs struct {
	XMLName             xml.Name `xml:"SpecialPrefs"`
	WrittenConfirmation *WrittenConfirmation
}

// addRoomTypeUnits to the existing hotel reservation body
func (h *HotelRsrvBody) addRoomTypeUnits(units int) {
	h.OTAHotelResRQ.Hotel.RoomType = &RoomType{
		NumberOfUnits: units,
	}
}

// QuerySearchParams is a typed function to support optional query params on creation of new search criterion
type SpecialPrefOptions func(*SpecialPrefs) error

// NewSpecialPrefs accepts set of SpecialPrefOptions functions and returns modified options
func NewSpecialPrefs(options ...SpecialPrefOptions) (*SpecialPrefs, error) {
	prefs := &SpecialPrefs{}
	for _, qm := range options {
		err := qm(prefs)
		if err != nil {
			return prefs, err
		}
	}
	return prefs, nil
}

// WrittenConf function sets written confirmation for special preferences
// Example: NewSpecialPrefs(WrittenConf(true))
func WrittenConf(opt bool) func(q *SpecialPrefs) error {
	return func(s *SpecialPrefs) error {
		s.WrittenConfirmation = &WrittenConfirmation{Ind: opt}
		return nil
	}
}

// addSpecialPrefs to the existing hotel reservation body
func (h *HotelRsrvBody) addSpecialPrefs(p *SpecialPrefs) {
	h.OTAHotelResRQ.Hotel.SpecialPrefs = p
}

func SetHotelResBody(rph int, gtype, ccCode, ccExpire, ccNumber, pnrLast string) HotelRsrvBody {
	return HotelRsrvBody{
		OTAHotelResRQ: OTAHotelResRQ{
			XMLNS:    srvc.BaseWebServicesNS,
			XMLNSXs:  srvc.BaseXSDNameSpace,
			XMLNSXsi: srvc.BaseXSINamespace,
			Version:  "2.2.0",
			Hotel: HotelRequest{
				BasicPropertyRes: BasicPropertyRes{RPH: rph},
				Guarantee: GuaranteeReservation{
					Type: gtype,
					CCInfo: CCInfo{
						PaymentCard: PaymentCard{
							Code:       ccCode,
							ExpireDate: ccExpire,
							Number:     ccNumber,
						},
						PersonName: itin.PersonName{
							Last: itin.Surname{Val: pnrLast},
						},
					},
				},
			},
		},
	}
}

// BuildHotelResRequest build request body for SOAP reservation service
func BuildHotelResRequest(from, pcc, binsectoken, convid, mid, time string, body HotelRsrvBody) HotelRsrvRequest {
	return HotelRsrvRequest{
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
				Service:        srvc.ServiceElem{Value: "OTA_HotelResRQ", Type: "sabreXML"},
				Action:         "OTA_HotelResRQ",
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
		Body: body,
	}
}

// OTAHotelResRS parse sabre hotel availability
type OTAHotelResRS struct {
	XMLName  xml.Name `xml:"OTA_HotelResRS"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSXs  string   `xml:"xs,attr"`
	XMLNSXsi string   `xml:"xsi,attr"`
	XMLNSStl string   `xml:"stl,attr"`
	Version  string   `xml:"Version,attr"`
	Result   ApplicationResults
	Hotel    HotelResponse
}

// HotelReserveRQ holds hotel information specific to making a hote reservation
type HotelResponse struct {
	XMLName       xml.Name `xml:"Hotel"`
	BasicProperty BasicPropertyInfo
	Guarantee     GuaranteeReservation
	RoomRates     []RoomRate `xml:"RoomRates>RoomRate"`
	SpecialPrefs  SpecialPrefs
	Text          []string `xml:"Text"`
	TimeSpan      TimeSpan
}

// HotelRsrvResponse for parsing hotel reservation request
type HotelRsrvResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		HotelRes OTAHotelResRS
		Fault    srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallHotelAvail to sabre web services
func CallHotelRes(serviceURL string, req HotelRsrvRequest) (HotelRsrvResponse, error) {
	resResp := HotelRsrvResponse{}
	//construct payload
	byteReq, _ := xml.Marshal(req)
	fmt.Printf("\n\nCallHotelResPAYLOAD: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		resResp.ErrorSabreService = sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadService)
		return resResp, resResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	fmt.Printf("\n\nCallHotelResRESPONSE: %s\n\n", bodyBuffer)
	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &resResp)
	if err != nil {
		resResp.ErrorSabreXML = sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadParse)
		return resResp, resResp.ErrorSabreXML
	}
	return resResp, nil
}
