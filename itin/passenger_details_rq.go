package itin

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/srvc"
)

/* PsngrDetailsRequest root level struct for dealing with an PNR. Taken from Sabre docs:

According to your specification in the request, does one of the following:
    * Displays the entire PNR. Returns the record locator when all processing of the service is completed.
    * You can use this service to create a PNR by adding traveler information for a maximum of 99 travelers, or you can add remarks and SSRs to an existing PNR and travelers.
    * A group can also be added.

In either case, at least one segment must be sold with content present in the Sabre work area (the AAA). The segments can be of the following types: air, hotel, vehicle, rail, or cruise. OTA_AirBookLLSRQ, Enhanced_AirBookRQ, or Enhanced_AirBookWithTaXRQ can be utilized to sell air segments. OTA_VehResLLSRQ can be used to sell car segments. OTA_HotelResLLSRQ can be used to sell hotel segments.

A successful transaction creates a new PNR or updates an existing PNR, saving the content you pass in the Sabre system. The system assigns a record locator for a new PNR, and returns the record locator of an existing PNR. When the processing of the service is complete, the content remains in the Sabre work area.
*/
type PsngrDetailsRequest struct {
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
	NameNumber    string   `xml:"NameNumber,attr"`    //1.1
	NameReference string   `xml:"NameReference,attr"` //ABC123
	PassengerType string   `xml:"PassengerType,attr"` //ADT
	Given         GivenName
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
type VendorPrefs struct {
	XMLName xml.Name `xml:"VendorPrefs"`
	Airline struct {
		Hosted bool `xml:"Hosted,attr"`
	} `xml:"Airline"`
}

// Address PNR specific struct for addresses
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNumber,omitempty"`
	City          string `xml:"CityName,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr,omitempty"`
	} `xml:"StateCountyProv,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
}
type AgencyInfo struct {
	Address     Address
	VendorPrefs VendorPrefs
}
type CustomerInfo struct {
	XMLName        xml.Name        `xml:"CustomerInfo"`
	ContactNumbers []ContactNumber `xml:"ContactNumbers>ContactNumber"`
	PersonName     PersonName
}
type ContactNumber struct {
	XMLName      xml.Name `xml:"ContactNumber"`
	NameNumber   string   `xml:"NameNumber,attr"`   //1.1
	Phone        string   `xml:"Phone,attr"`        //123-456-7890
	PhoneUseType string   `xml:"PhoneUseType,attr"` //H|M
}

//TravelItineraryAddInfoRQ basic information for agency and customer
type TravelItineraryAddInfoRQ struct {
	XMLName  xml.Name `xml:"TravelItineraryAddInfoRQ"`
	Agency   *AgencyInfo
	Customer CustomerInfo
}

func (p *PassengerDetailBody) AddSpecialDetails() {
	p.PassengerDetailsRQ.SpecialReq = &SpecialReqDetails{}
}

// SetHotelRateDescRqStruct hotel rate description request using input parameters
func SetPsngrDetailsRequestStruct(phone, firstName, lastName string) PassengerDetailBody {
	return PassengerDetailBody{
		PassengerDetailsRQ: PassengerDetailsRQ{
			XMLNS:         "http://services.sabre.com/sp/pd/v3_3",
			Version:       "3.3.0",
			IgnoreOnError: false,
			HaltOnError:   false,
			PostProcess: PostProcessing{
				IgnoreAfter:          false,
				RedisplayReservation: true,
				UnmaskCreditCard:     false,
			},
			PreProcess: PreProcessing{
				IgnoreBefore: true,
				//UniqueID:     UniqueID{ID: lastName + srvc.GenerateSessionID()},
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
					PersonName: PersonName{
						NameNumber:    "1.1",
						NameReference: "ABC123",
						PassengerType: "ADT",
						Given:         GivenName{Val: firstName},
						Last:          Surname{Val: lastName},
					},
				},
			},
		},
	}
}

// BuildPsngrDetailsRequest passenger details for booking
func BuildPsngrDetailsRequest(from, pcc, binsectoken, convid, mid, time string, body PassengerDetailBody) PsngrDetailsRequest {
	return PsngrDetailsRequest{
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
				Service:        srvc.ServiceElem{Value: "PassengerDetailsRQ", Type: "sabreXML"},
				Action:         "PassengerDetailsRQ",
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

// CallPsngrDetailsRequest creates a new PNR or updates an existing PNR, saving the content you pass in the Sabre system. The system assigns a record locator for a new PNR, and returns the record locator of an existing PNR. When the processing of the service is complete, the content remains in the Sabre work area. Previous calls required are hotel_property_desc OR hotel_rate_desc call, see BuildPsngrDetailsRequest.
func CallPsngrDetailsRequest(serviceURL string, req PsngrDetailsRequest) { //(PnsgrDetailsResponse, error) {
	//psngrResp := HotelRateDescResponse{}
	byteReq, _ := xml.Marshal(req)
	fmt.Printf("REQ: %s\n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	fmt.Printf("RESP: %s\n\n", resp.Status)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		//psngrResp.ErrorSabreService = NewErrorSabreService(err.Error(), ErrCallHotelRateDesc, BadService)
		//return psngrResp, psngrResp.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	fmt.Printf("MARSH-RESP: %s\n\n", bodyBuffer.Bytes())
	/*
		//marshal bytes sabre response body into availResp response struct
		err = xml.Unmarshal(bodyBuffer.Bytes(), &psngrResp)
		if err != nil {
			psngrResp.ErrorSabreXML = NewErrorErrorSabreXML(err.Error(), ErrCallHotelRateDesc, BadParse)
			return psngrResp, psngrResp.ErrorSabreXML
		}
		return psngrResp, nil
	*/
}
