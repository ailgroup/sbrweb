package itin

import "encoding/xml"

type ProductDetails struct {
	XMLName            xml.Name `xml:"ProductDetails"`
	VendorCode         string   `xml:"vendorCode,attr"`
	StatusCode         string   `xml:"statusCode,attr"`
	PreviousStatusCode string   `xml:"previousStatusCode,attr"`
	StartDateTime      string   `xml:"startDateTime,attr"`
	EndDateTime        string   `xml:"endDateTime,attr"`
	ProductName        struct {
		PType string `xml:"type,attr"`
	} `xml:"ProductName"`
	Hotel HotelSegmentElem
	//Vehicle VehicleSegmentElem
}
type ProudctSegmentElem struct {
	XMLName     xml.Name `xml:"Product"`
	ProductBase struct {
		XMLName          xml.Name `xml:"ProductBase"`
		SegmentReference string   `xml:"SegmentReference"`
	} `xml:"ProductBase"`
	ProductDetails ProductDetails
}
type ReservationVehicle struct {
	//XMLName xml.Name `xml:"Reservation"`
}
type VehicleSegmentElem struct {
	//XMLName     xml.Name           `xml:"Vehicle"`
	ID          string             `xml:"id,attr"`
	Sequence    string             `xml:"sequence,attr"`
	IsPast      bool               `xml:"isPast,attr"`
	Reservation ReservationVehicle //`xml:",omitempty"`
}
type RoomTypeRes struct {
	XMLName       xml.Name `xml:"RoomType"`
	RoomTypeCode  string   `xml:"RoomTypeCode"`
	NumberOfUnits string   `xml:"NumberOfUnits"`
	ShortText     string   `xml:"ShortText"`
}
type RoomRatesRes struct {
	XMLName         xml.Name `xml:"RoomRates"`
	AmountBeforeTax string   `xml:"AmountBeforeTax"`
	CurrencyCode    string   `xml:"CurrencyCode"`
}
type GuestCountsRes struct {
	XMLName         xml.Name `xml:"GuestCounts"`
	GuestCount      string   `xml:"GuestCount"`
	ExtraGuestCount string   `xml:"ExtraGuestCount"`
	RollAwayCount   string   `xml:"RollAwayCount"`
	CribCount       string   `xml:"CribCount"`
}
type GuaranteeRes struct {
	XMLName xml.Name `xml:"Guarantee"`
	Text    string   `xml:"Text"`
}
type TaxRes struct {
	XMLName xml.Name `xml:"Tax"`
	ID      string   `xml:"Id,attr"`
	Val     string   `xml:",chardata"`
}
type TotalTaxRes struct {
	XMLName xml.Name `xml:"TotalTax"`
	Amount  string   `xml:"Amount,attr"`
	Tax     TaxRes
}
type ApproximateTotalRes struct {
	XMLName           xml.Name `xml:"ApproximateTotal"`
	AmountAndCurrency string   `xml:"AmountAndCurrency,attr"`
}
type DisclaimerRes struct {
	XMLName xml.Name `xml:"Disclaimer"`
	ID      string   `xml:"Id,attr"`
	Val     string   `xml:",chardata"`
}
type HotelTotalPricingRes struct {
	XMLName          xml.Name `xml:"HotelTotalPricing"`
	TotalTax         TotalTaxRes
	ApproximateTotal ApproximateTotalRes
	Disclaimer       DisclaimerRes
}
type ReservatioHotel struct {
	XMLName           xml.Name `xml:"Reservation"`
	DayOfWeekInd      string   `xml:"DayOfWeekInd,attr"`
	NumberInParty     string   `xml:"NumberInParty,attr"`
	LineNumber        string   `xml:"LineNumber"`
	LineType          string   `xml:"LineType"`
	LineStatus        string   `xml:"LineStatus"`
	POSRequestor      string   `xml:"POSRequestorID"`
	RoomType          RoomTypeRes
	RoomRates         RoomRatesRes
	GuestCounts       GuestCountsRes
	TimeSpanStart     string `xml:"TimeSpanStart"`
	TimeSpanDuration  string `xml:"TimeSpanDuration"`
	TimeSpanEnd       string `xml:"TimeSpanEnd"`
	Guarantee         GuaranteeRes
	ChainCode         string `xml:"ChainCode"`
	HotelCode         string `xml:"HotelCode"`
	HotelCityCode     string `xml:"HotelCityCode"`
	HotelName         string `xml:"HotelName"`
	HotelTotalPricing HotelTotalPricingRes
}
type AddressAdditional struct {
	XMLName     xml.Name `xml:"Address"`
	AdressLine  []string `xml:"AddressLine"`
	CountryCode string   `xml:"CountryCode"`
	City        string   `xml:"City"`
	ZipCode     string   `xml:"ZipCode"`
}
type ContactNumbersAdditional struct {
	XMLName     xml.Name `xml:"ContactNumbers"`
	PhoneNumber string   `xml:"PhoneNumber"`
	FaxNumber   string   `xml:"FaxNumber"`
}
type CommissionAdditional struct {
	XMLName   xml.Name `xml:"Commission"`
	Indicator string   `xml:"Indicator"`
	Text      string   `xml:"Text"`
}
type AdditionalInformation struct {
	XMLName                 xml.Name `xml:"AdditionalInformation"`
	Address                 AddressAdditional
	ContactNumbers          ContactNumbersAdditional
	CancelPenaltyPolicyCode string `xml:"CancelPenaltyPolicyCode"`
	Commission              CommissionAdditional
}
type RateDescriptionSegment struct {
	XMLName  xml.Name `xml:"RateDescription"`
	TextLine []string `xml:"TextLine"`
}
type HotelPolicySegment struct {
	XMLName            xml.Name `xml:"HotelPolicy"`
	GuaranteePolicy    string   `xml:"GuaranteePolicy"`
	CancellationPolicy string   `xml:"CancellationPolicy"`
}
type HotelSegmentElem struct {
	XMLName               xml.Name `xml:"Hotel"`
	ID                    string   `xml:"id,attr"`
	Sequence              string   `xml:"sequence,attr"`
	IsPast                bool     `xml:"isPast,attr"`
	Reservation           ReservatioHotel
	AdditionalInformation AdditionalInformation
	SegmentText           string                 `xml:"SegmentText"`
	RateDescription       RateDescriptionSegment //only on product
	HotelPolicy           HotelPolicySegment     //only on product
}
type SegmentReservation struct {
	XMLName  xml.Name `xml:"Segment"`
	Sequence string   `xml:"sequence,attr"`
	ID       string   `xml:"id,attr"`
	Hotel    HotelSegmentElem
	//Vehicle  VehicleSegmentElem
	Product ProudctSegmentElem
}
