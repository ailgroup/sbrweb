package hotelws

import "encoding/xml"

type PaymentCard struct {
	XMLName xml.Name `xml:"PaymentCard"`
	Code    string   `xml:"Code,attr"`
	Type    string   `xml:"Type,attr"`
}

type DepositsAccepted struct {
	XMLName      xml.Name `xml:"DepositsAccepted"`
	PaymentCards []PaymentCard
}

type GuaranteesAccepted struct {
	XMLName      xml.Name `xml:"GuaranteesAccepted"`
	PaymentCards []PaymentCard
}

type GuaranteeAccepted struct {
	XMLName    xml.Name `xml:"Guarantee"`
	Guarantees GuaranteesAccepted
	Deposits   DepositsAccepted
}

type RoomStay struct {
	XMLName           xml.Name `xml:"RoomStay"`
	BasicPropertyInfo BasicPropertyInfo
	GuaranteeAccepted GuaranteeAccepted
	RoomRates         []RoomRate `xml:"RoomRates>RoomRate"`
	TimeSpan          struct {
		Duration string `xml:"Duration,attr"` //string 0001 or int 1?
		End      string `xml:"End,attr"`
		Start    string `xml:"Start,attr"`
	} `xml:"TimeSpan"`
}
