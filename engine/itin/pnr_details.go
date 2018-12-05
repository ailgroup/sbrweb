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

/* PNRDetailsRequest root level struct for dealing with an PNR. Taken from Sabre docs:

According to your specification in the request, does one of the following:
    * Displays the entire PNR. Returns the record locator when all processing of the service is completed.
    * You can use this service to create a PNR by adding traveler information for a maximum of 99 travelers, or you can add remarks and SSRs to an existing PNR and travelers.
    * A group can also be added.

In either case, at least one segment must be sold with content present in the Sabre work area (the AAA). The segments can be of the following types: air, hotel, vehicle, rail, or cruise. OTA_AirBookLLSRQ, Enhanced_AirBookRQ, or Enhanced_AirBookWithTaXRQ can be utilized to sell air segments. OTA_VehResLLSRQ can be used to sell car segments. OTA_HotelResLLSRQ can be used to sell hotel segments.

A successful transaction creates a new PNR or updates an existing PNR, saving the content you pass in the Sabre system. The system assigns a record locator for a new PNR, and returns the record locator of an existing PNR. When the processing of the service is complete, the content remains in the Sabre work area.
*/
type PNRDetailsRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   PassengerDetailBody
}

// PassengerDetailBody holds namespaced body
type PassengerDetailBody struct {
	XMLName            xml.Name `xml:"soap-env:Body"`
	PassengerDetailsRQ PassengerDetailsRQ
}

// PassengerDetailsRQ main element
type PassengerDetailsRQ struct {
	XMLName        xml.Name `xml:"PassengerDetailsRQ"`
	XMLNS          string   `xml:"xmlns,attr"`
	Version        string   `xml:"version,attr"`
	IgnoreOnError  bool     `xml:"IgnoreOnError,attr"`
	HaltOnError    bool     `xml:"HaltOnError,attr"`
	PostProcess    PostProcessing
	PreProcess     PreProcessing
	SpecialReq     *SpecialReqDetails
	TravelItinInfo TravelItineraryAddInfoRQ
}
type PostProcessing struct {
	XMLName              xml.Name `xml:"PostProcessing"`
	IgnoreAfter          bool     `xml:"IgnoreAfter,attr"`
	RedisplayReservation bool     `xml:"RedisplayReservation,attr"`
	UnmaskCreditCard     bool     `xml:"UnmaskCreditCard,attr"`
}
type PreProcessing struct {
	XMLName      xml.Name  `xml:"PreProcessing"`
	IgnoreBefore bool      `xml:"IgnoreBefore,attr"`
	UniqueID     *UniqueID // existing PNR id goes here...
}
type UniqueID struct {
	XMLName xml.Name `xml:"UniqueID"`
	ID      string   `xml:"ID,attr"`
}
type SpecialReqDetails struct {
	XMLName          xml.Name `xml:"SpecialRequestDetails"`
	SpecialServiceRQ SpecialServiceRQ
}
type SpecialServiceRQ struct {
	XMLName            xml.Name `xml:"SpecialServiceRQ"`
	SpecialServiceInfo SpecialServiceInfo
}
type SpecialServiceInfo struct {
	XMLName           xml.Name `xml:"SpecialServiceInfo"`
	AdvancedPassenger AdvancedPassenger
}
type AdvancedPassenger struct {
	XMLName       xml.Name `xml:"AdvancePassenger"`
	SegmentNumber string   `xml:"SegmentNumber,attr"` //A
	Document      Document
	PersonName    PersonName
	VendorPrefs   VendorPrefs
}
type Document struct {
	IssueCountry struct {
	} `xml:"IssueCountry,omitempty"`
	NationalityCountry struct {
	} `xml:"NationalityCountry,omitempty"`
}
type PersonName struct {
	XMLName       xml.Name `xml:"PersonName"`
	WithInfant    bool     `xml:"WithInfant,attr,omitempty"`
	NameNumber    string   `xml:"NameNumber,attr,omitempty"`    //1.1
	NameReference string   `xml:"NameReference,attr,omitempty"` //ABC123
	PassengerType string   `xml:"PassengerType,attr,omitempty"` //ADT
	RPH           int      `xml:"RPH,attr,omitempty"`           //1 OR 001
	First         *GivenName
	Middle        *MiddleName
	Last          Surname
}
type GivenName struct {
	XMLName xml.Name `xml:"GivenName"`
	Val     string   `xml:",chardata"`
}
type MiddleName struct {
	XMLName xml.Name `xml:"MiddleName"`
	Val     string   `xml:",chardata"`
}
type Surname struct {
	XMLName xml.Name `xml:"Surname"`
	Val     string   `xml:",chardata"`
}

type Airline struct {
	XMLName xml.Name `xml:"Airline"`
	Hosted  bool     `xml:"Hosted,attr"`
}
type VendorPrefs struct {
	XMLName xml.Name `xml:"VendorPrefs"`
	Airline Airline
}

type StateProvince struct {
	XMLName   xml.Name `xml:"StateCountyProv,omitempty"`
	StateCode string   `xml:"StateCode,attr,omitempty"`
}

// Address PNR specific struct for addresses
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	City          string `xml:"CityName,omitempty"`
	CountryCode   string `xml:"CountryCode,omitempty"`
	Postal        string `xml:"PostalCode,omitempty"`
	StateProvince StateProvince
	Street        string `xml:"StreetNmbr,omitempty"`
	VendorPrefs   VendorPrefs
}

type CustomerInfo struct {
	XMLName        xml.Name        `xml:"CustomerInfo"`
	ContactNumbers []ContactNumber `xml:"ContactNumbers>ContactNumber"`
	PersonName     PersonName
}
type ContactNumber struct {
	XMLName      xml.Name `xml:"ContactNumber"`
	RPH          int      `xml:"RPH,attr,omitempty"` //1 OR 001
	LocationCode string   `xml:"LocationCode,attr,omitempty"`
	NameNumber   string   `xml:"NameNumber,attr,omitempty"`   //1.1
	Phone        string   `xml:"Phone,attr,omitempty"`        //123-456-7890 OR 123-456-7890-H.1.1
	PhoneUseType string   `xml:"PhoneUseType,attr,omitempty"` //H|M
}

//TravelItineraryAddInfoRQ basic information for agency and customer
type TravelItineraryAddInfoRQ struct {
	XMLName  xml.Name `xml:"TravelItineraryAddInfoRQ"`
	Agency   *AgencyInfo
	Customer CustomerInfo
}
type AgencyInfo struct {
	XMLName xml.Name `xml:"AgencyInfo"`
	Address Address
}

// AddAgencyInfo required to complete booking. Helper function allows it to be more fleixible to build up travel itinerary PNR.
func (ti *TravelItineraryAddInfoRQ) AddAgencyInfo(addr Address) {
	ti.Agency = &AgencyInfo{
		Address: addr,
	}
}

// AddSpecialDetails optionally include special details in special requests
func (p *PassengerDetailBody) AddSpecialDetails() {
	p.PassengerDetailsRQ.SpecialReq = &SpecialReqDetails{}
}

// AddUniqueID optionally include a pre processing unique ID
func (p *PassengerDetailBody) AddUniqueID(id string) {
	p.PassengerDetailsRQ.PreProcess.UniqueID = &UniqueID{ID: id}
}

// CreatePersonName standalone function for ease of use
func CreatePersonName(firstName, lastName string) PersonName {
	return PersonName{
		NameNumber:    "1.1",    // sabre example
		NameReference: "ABC123", // sabre example
		PassengerType: "ADT",    // sabre example
		First:         &GivenName{Val: firstName},
		Last:          Surname{Val: lastName},
	}
}

/*
SetHotelRateDescRqStruct hotel rate description request using input parameters
	IgnoreOnError: false,
	HaltOnError:   true,
	PostProcess: PostProcessing{
		IgnoreAfter:          false,
		RedisplayReservation: true,
		UnmaskCreditCard:     false,
	},
	PreProcess: PreProcessing{
		IgnoreBefore: true,

*/
func SetPNRDetailBody(phone string, person PersonName) PassengerDetailBody {
	return PassengerDetailBody{
		PassengerDetailsRQ: PassengerDetailsRQ{
			XMLNS:         "http://services.sabre.com/sp/pd/v3_3",
			Version:       "3.3.0",
			IgnoreOnError: false,
			HaltOnError:   true,
			PostProcess: PostProcessing{
				IgnoreAfter:          false,
				RedisplayReservation: true,
				UnmaskCreditCard:     false,
			},
			PreProcess: PreProcessing{
				IgnoreBefore: true,
			},
			TravelItinInfo: TravelItineraryAddInfoRQ{
				Customer: CustomerInfo{
					ContactNumbers: []ContactNumber{
						ContactNumber{
							NameNumber:   "1.1",
							Phone:        phone,
							PhoneUseType: "H",
						},
					},
					PersonName: person,
				},
			},
		},
	}
}

// BuildPNRDetailsRequest passenger details for booking
func BuildPNRDetailsRequest(c *srvc.SessionConf, body PassengerDetailBody) PNRDetailsRequest {
	return PNRDetailsRequest{
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
				Service:        srvc.ServiceElem{Value: "PassengerDetailsRQ", Type: "sabreXML"},
				Action:         "PassengerDetailsRQ",
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
		Body: body,
	}
}

type Message struct {
	Code string `xml:"code"`
	Val  string `xml:",chardata"`
}
type SystemResult struct {
	Messages []Message `xml:"Message"`
}
type Warning struct {
	Type          string         `xml:"type,attr"`
	Timestamp     string         `xml:"timeStamp,attr"`
	SystemResults []SystemResult `xml:"SystemSpecificResults"`
}
type ApplicationResults struct {
	XMLName xml.Name `xml:"ApplicationResults"`
	Status  string   `xml:"status,attr"`
	Success struct {
		Timestamp string `xml:"timeStamp,attr"`
	} `xml:"Success"`
	Warnings []Warning `xml:"Warning"`
}

func (result ApplicationResults) Ok() bool {
	switch result.Status {
	case sbrerr.StatusNotProcess(): //queries
		return false
	case sbrerr.StatusComplete(): //queries, pnr
		if len(result.Warnings) > 0 {
			return false
		}
		return true
	default:
		return false
	}
}

func (result ApplicationResults) ErrFormat() sbrerr.ErrorSabreResult {
	var wmsg string
	for i, w := range result.Warnings {
		var msg string
		for is, s := range w.SystemResults {
			for ms, m := range s.Messages {
				msg += fmt.Sprintf("SystemResult-%d:Message-%d:Code-%s %s. ", is, ms, m.Code, m.Val)
			}
		}
		wmsg += fmt.Sprintf("Warning-%d:Type-%s Results: %s", i, w.Type, msg)
	}
	return sbrerr.ErrorSabreResult{
		Code:       sbrerr.SabreEngineStatusCode(result.Status),
		AppMessage: wmsg,
	}
}

type ReservationItem struct {
}
type ItineraryInfo struct {
	XMLName          xml.Name          `xml:"ItineraryInfo"`
	ReservationItems []ReservationItem `xml:"ReservationItems"`
}
type ItineraryRef struct {
	XMLName     xml.Name `xml:"ItineraryRef"`
	AirExtras   bool     `xml:"AirExtras,attr"`
	InhibitCode string   `xml:"InhibitCode,attr"`
	PartitionID string   `xml:"PartitionID,attr"`
	PrimeHostID string   `xml:"PrimeHostID,attr"`
	Source      struct {
		PseudoCityCode string `xml:"PseudoCityCode,attr"`
	} `xml:"Source"`
}
type TravelItinerary struct {
	XMLName       xml.Name `xml:"TravelItinerary"`
	Customer      CustomerInfo
	ItineraryInfo ItineraryInfo
	ItineraryRef  ItineraryRef
}
type TravelItineraryReadRS struct {
	XMLName         xml.Name `xml:"TravelItineraryReadRS"`
	TravelItinerary TravelItinerary
}
type PassengerDetailsRS struct {
	XMLName               xml.Name `xml:"PassengerDetailsRS"`
	AppResults            ApplicationResults
	TravelItineraryReadRS TravelItineraryReadRS
}
type PNRDetailsResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		PassengerDetailsRS PassengerDetailsRS
		Fault              srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallPNRDetailsRequest creates a new PNR or updates an existing PNR, saving the content you pass in the Sabre system. The system assigns a record locator for a new PNR, and returns the record locator of an existing PNR. When the processing of the service is complete, the content remains in the Sabre work area. Previous calls required are hotel_property_desc OR hotel_rate_desc call, see BuildPNRDetailsRequest.
func CallPNRDetail(serviceURL string, req PNRDetailsRequest) (PNRDetailsResponse, error) {
	pnrResp := PNRDetailsResponse{}
	byteReq, _ := xml.Marshal(req)
	fmt.Printf("\n\n %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		pnrResp.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallPNRDetails,
			sbrerr.BadService,
		)
		return pnrResp, pnrResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//fmt.Printf("\n\n:CallPNRDetail Response: %s\n\n", bodyBuffer)
	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &pnrResp)
	if err != nil {
		pnrResp.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallPNRDetails,
			sbrerr.BadParse,
		)
		return pnrResp, pnrResp.ErrorSabreXML
	}
	if !pnrResp.Body.Fault.Ok() {
		return pnrResp, sbrerr.NewErrorSoapFault(pnrResp.Body.Fault.String)
	}
	if !pnrResp.Body.PassengerDetailsRS.AppResults.Ok() {
		return pnrResp, pnrResp.Body.PassengerDetailsRS.AppResults.ErrFormat()
	}
	return pnrResp, nil
}
