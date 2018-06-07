package hotelws

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/srvc"
)

/*
More options here: http://webservices.sabre.com/drc/servicedoc/OTA_HotelResLLSRQ_v2.2.0_Design.xml

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
	GuestCounts      GuestCounts
	RoomType         RoomType
	SpecialPrefs     *SpecialPrefs
	TimeSpan         TimeSpan
}

// BasicPropertyRes is the BasicPropertyInfo element specifically for executing hotel reservations. Easier to duplicate this simple case than omit all the struct fields in the BasicPropertyInfo type.
type BasicPropertyRes struct {
	XMLName     xml.Name `xml:"BasicPropertyInfo"`
	ChainCode   string   `xml:"ChainCode,attr,omitempty"`
	HotelCode   string   `xml:"HotelCode,attr,omitempty"`
	InsertAfter string   `xml:"InsertAfter,attr,omitempty"`
	RPH         int      `xml:"RPH,attr"`
}

// CCInfo for passing credit card
type CCInfo struct {
	XMLName     xml.Name `xml:"CC_Info"`
	PaymentCard PaymentCard
	PersonName  PersonNameRes
}

type PersonNameRes struct {
	XMLName xml.Name `xml:"PersonName"`
	Surname Surname
}

type Surname struct {
	XMLName xml.Name `xml:"Surname"`
	Val     string   `xml:",chardata"`
}

// GuaranteeReservation is a gurantee type specifically for executing hotel reservations
// Type can be "G", "GAGT", "GDPST", "GC", "GCR", "GH", "GDPSTH", "GT", or "GDPSTT", or "D"
type GuaranteeReservation struct {
	XMLName xml.Name `xml:"Guarantee"`
	Type    string   `xml:"Type,attr"`
	CCInfo  CCInfo
}

type RoomType struct {
	XMLName       xml.Name `xml:"RoomType"`
	NumberOfUnits int      `xml:"NumberOfUnits,attr"`
	RoomTypeCode  string   `xml:"RoomTypeCode,attr"`
}

type WrittenConfirmation struct {
	XMLName xml.Name `xml:"WrittenConfirmation"`
	Ind     bool     `xml:"Ind,attr"`
}

type SpecRefText struct {
	Val string `xml:",chardata"`
}

// SpecialPrefs allows adding extra customer information
type SpecialPrefs struct {
	XMLName             xml.Name `xml:"SpecialPrefs"`
	WrittenConfirmation WrittenConfirmation
	Text                []SpecRefText //`xml:"Text"`
}

// CreateSpecialPrefs creates the value. Must be done before adding special prefs
func (h *HotelRsrvBody) AddSpecialPrefs(p *SpecialPrefs) {
	h.OTAHotelResRQ.Hotel.SpecialPrefs = p
}

// AddSpecialPrefs to the existing hotel reservation body. See CreateSpecialPrefs.
func (s *SpecialPrefs) AddSpecPrefWritConf(opt bool) {
	s.WrittenConfirmation = WrittenConfirmation{Ind: opt}
}

// AddSpecialPrefs to the existing hotel reservation body. See CreateSpecialPrefs.
func (s *SpecialPrefs) AddSpecPrefText(vals []string) {
	s.Text = []SpecRefText{}
	for _, v := range vals {
		s.Text = append(s.Text, SpecRefText{Val: v})
	}
}

// AddRoomType to the existing hotel reservation body
func (h *HotelRsrvBody) AddRoomType(units int, roomCode string) {
	h.OTAHotelResRQ.Hotel.RoomType = RoomType{
		NumberOfUnits: units,
		RoomTypeCode:  roomCode,
	}
}

// NewGuaranteeRes builds and sets guarantee and credit card info on hotel res
func (h *HotelRsrvBody) NewGuaranteeRes(lastName, gtype, ccCode, ccExpire, ccNumber string) {
	h.OTAHotelResRQ.Hotel.Guarantee = GuaranteeReservation{
		Type: gtype,
		CCInfo: CCInfo{
			PaymentCard: PaymentCard{
				Code:       ccCode,
				ExpireDate: ccExpire,
				Number:     ccNumber,
			},
			PersonName: PersonNameRes{
				Surname: Surname{Val: lastName},
			},
		},
	}
}

func (h *HotelRsrvBody) NewPropertyRes(rph int, chain, hotel string) {
	h.OTAHotelResRQ.Hotel.BasicPropertyRes = BasicPropertyRes{
		ChainCode: chain,
		HotelCode: hotel,
		RPH:       rph,
	}
}

func SetHotelResBody(guestCount int, timesp TimeSpan) HotelRsrvBody {
	return HotelRsrvBody{
		OTAHotelResRQ: OTAHotelResRQ{
			XMLNS:    srvc.BaseWebServicesNS,
			XMLNSXs:  srvc.BaseXSDNameSpace,
			XMLNSXsi: srvc.BaseXSINamespace,
			Version:  "2.2.0",
			Hotel: HotelRequest{
				GuestCounts: GuestCounts{Count: guestCount},
				TimeSpan:    timesp,
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
				Service:        srvc.ServiceElem{Value: "OTA_HotelRes", Type: "sabreXML"},
				Action:         "OTA_HotelResLLSRQ",
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
		return resResp, sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadService)
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
		fmt.Println("1")
		return resResp, sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallHotelAvail, sbrerr.BadParse)
	}
	if !resResp.Body.Fault.Ok() {
		fmt.Println("2")
		return resResp, sbrerr.NewErrorSoapFault(resResp.Body.Fault.String)
	}
	if !resResp.Body.HotelRes.Result.Ok() {
		fmt.Println("3")
		return resResp, resResp.Body.HotelRes.Result.ErrFormat()
	}
	fmt.Println("4")
	return resResp, nil
}
