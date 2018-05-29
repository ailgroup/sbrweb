package hotelws

import (
	"encoding/xml"

	"github.com/ailgroup/sbrweb/srvc"
)

// HotelResRequest for soap package on OTA_HotelResRQ service for making reservations
type HotelResRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   HotelResBody
}

// HotelResBody implements Hotel element for SOAP
type HotelResBody struct {
	XMLName xml.Name `xml:"OTA_HotelResRQ"`
	Hotel   HotelRes
}

// HotelRes holds hotel information specific to making a hote reservation
type HotelRes struct {
	XMLName          xml.Name `xml:"Hotel"`
	BasicPropertyRes BasicPropertyRes
	Guarantee        GuaranteeReservation
}

// BasicPropertyRes is the BasicPropertyInfo element specifically for executing hotel reservations. Easier to duplicate this simple case than omit all the struct fields in the BasicPropertyInfo type.
type BasicPropertyRes struct {
	XMLName xml.Name `xml:"BasicPropertyInfo"`
	RPH     int      `xml:"RPH,attr"`
}

// GuaranteeReservation is a gurantee type specifically for executing hotel reservations
type GuaranteeReservation struct {
	XMLName xml.Name `xml:"Guarantee"`
	CCInfo  CCInfo
}

// CCInfo for passing credit card
type CCInfo struct {
	XMLName     xml.Name `xml:"CC_Info"`
	PaymentCard PaymentCard
}

func SetHotelResBody() HotelResBody {
	return HotelResBody{}
}

// BuildHotelResRequest build request body for SOAP reservation service
func BuildHotelResRequest(from, pcc, binsectoken, convid, mid, time string, body HotelResBody) HotelResRequest {
	return HotelResRequest{
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
