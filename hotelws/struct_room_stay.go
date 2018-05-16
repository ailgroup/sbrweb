package hotelws

import "encoding/xml"

// RoomStay contains all info relevant to the property's available rooms. It is the root-level element after service element for hotel_rate_desc and hotel_property_desc.
type RoomStay struct {
	XMLName           xml.Name `xml:"RoomStay"`
	BasicPropertyInfo BasicPropertyInfo
	Guarantee         Guarantee
	RoomRates         []RoomRate `xml:"RoomRates>RoomRate"`
	TimeSpan          struct {
		Duration string `xml:"Duration,attr"` //string 0001 or int 1?
		End      string `xml:"End,attr"`
		Start    string `xml:"Start,attr"`
	} `xml:"TimeSpan"`
}

// Guarantee shows forms of payment accepted by property
type Guarantee struct {
	XMLName            xml.Name `xml:"Guarantee"`
	GuaranteesAccepted GuaranteesAccepted
	DepositsAccepted   DepositsAccepted
}
type GuaranteesAccepted struct {
	XMLName      xml.Name      `xml:"GuaranteesAccepted"`
	PaymentCards []PaymentCard `xml:"PaymentCard"`
}
type DepositsAccepted struct {
	XMLName      xml.Name      `xml:"DepositsAccepted"`
	PaymentCards []PaymentCard `xml:"PaymentCard"`
}
type PaymentCard struct {
	Code string `xml:"Code,attr"`
	Type string `xml:"Type,attr"`
}

type RoomRate struct {
	XMLName            xml.Name `xml:"RoomRate"`
	DirectConnect      string   `xml:"RDirectConnect,attr"`
	GuaranteeSurcharge string   `xml:"GuaranteeSurchargeRequired,attr"`
	GuaranteedRate     string   `xml:"GuaranteedRateProgram,attr"`
	IATA_Character     string   `xml:"IATA_CharacteristicIdentification,attr"`
	IATA_Product       string   `xml:"IATA_ProductIdentification,attr"`
	LowInventory       string   `xml:"LowInventoryThreshold,attr"`
	RateLevelCode      string   `xml:"RateLevelCode,attr"`
	RPH                int      `xml:"RPH,attr"`
	RateChangeInd      string   `xml:"RateChangeInd,attr"`
	RateConversionInd  string   `xml:"RateConversionInd,attr"`
	SpecialOffer       string   `xml:"SpecialOffer,attr"`
	Rates              []Rate   `xml:"Rates>Rate"`
	AdditionalInfo     AdditionalInfo
	HotelRateCode      string `xml:"HotelRateCode"`
}
type AdditionalInfo struct {
	XMLName    xml.Name `xml:"AdditionalInfo"`
	Commission struct {
		NonCommission string `xml:"NonCommission,attr"`
		Val           string `xml:",char"`
	} `xml:"Commission"`
	DCACancellation struct {
		Text []string `xml:"Text"`
	} `xml:"DCA_Cancellation"`
	DCAGuarantee struct {
		Text []string `xml:"Text"`
	} `xml:"DCA_Guarantee"`
	DCAOther struct {
		Text []string `xml:"Text"`
	} `xml:"DCA_Other"`
	PaymentCards []PaymentCard `xml:"PaymentCard"`
	Taxes        string        `xml:"Taxes"`
	CancelPolicy struct {
		Numeric int    `xml:"Numeric,attr"` //string? 001 versus 1
		Option  string `xml:"Option,attr"`
	} `xml:"CancelPolicy"`
	Text []string `xml:"Text"`
}
