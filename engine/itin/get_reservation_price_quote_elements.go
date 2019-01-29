package itin

import "encoding/xml"

type HomeLocation struct {
	//XMLName xml.Name `xml:"HomeLocation"`
	Val string `xml:",chardata"`
}
type WorkLocation struct {
	//XMLName xml.Name `xml:"WorkLocation"`
	Val string `xml:",chardata"`
}
type AgentInfo struct {
	//XMLName       xml.Name `xml:"AgentInfo"`
	Duty         string `xml:"duty,attr"`
	Sine         string `xml:"sine,attr"`
	HomeLocation HomeLocation
	WorkLocation WorkLocation
}
type CreateDateTime struct {
	//XMLName xml.Name `xml:"CreateDateTime"`
	Val string `xml:",chardata"`
}
type LocalCreateDateTime struct {
	//XMLName xml.Name `xml:"LocalCreateDateTime"`
	Val string `xml:",chardata"`
}
type ExpiryDateTime struct {
	//XMLName xml.Name `xml:"ExpiryDateTime"`
	Val string `xml:",chardata"`
}
type InputEntry struct {
	//XMLName xml.Name `xml:"InputEntry"`
	Val string `xml:",chardata"`
}
type TransactionInfo struct {
	//XMLName       xml.Name `xml:"TransactionInfo"`
	CreateDateTime      CreateDateTime //string `xml:"CreateDateTime"`
	LocalCreateDateTime LocalCreateDateTime
	ExpiryDateTime      ExpiryDateTime
	InputEntry          InputEntry
}
type NameAssociationInfo struct {
	//XMLName    xml.Name `xml:"NameAssociationInfo"`
	FirstName  string `xml:"firstName,attr"`
	LastName   string `xml:"lastName,attr"`
	NameID     string `xml:"nameId,attr"`
	NameNumber string `xml:"nameNumber,attr"`
}
type MarketingFlight struct {
	//XMLName    xml.Name `xml:"MarketingFlight"`
	Number string `xml:"number,attr"`
	Val    string `xml:",chardata"`
}
type ClassOfService struct {
	//XMLName    xml.Name `xml:"ClassOfService"`
	Val string `xml:",chardata"`
}
type Flight struct {
	//XMLName    xml.Name `xml:"Flight"`
	ConnectionIndicator string `xml:"connectionIndicator,attr"`
	MarketingFlight     MarketingFlight
	ClassOfService      ClassOfService
	Departure           Departure
	Arrival             Arrival
}
type FareBasis struct {
	//XMLName    xml.Name `xml:"FareBasis"`
	Val string `xml:",chardata"`
}
type NotValidAfter struct {
	//XMLName    xml.Name `xml:"NotValidAfter"`
	Val string `xml:",chardata"`
}
type Baggage struct {
	//XMLName    xml.Name `xml:"Baggage"`
	Allowance string `xml:"allowance,attr"`
	Btype     string `xml:"type,attr"`
}
type SegmentInfo struct {
	//XMLName    xml.Name `xml:"SegmentInfo"`
	Number        string `xml:"number,attr"` //string easier deal w/ response and 0-th
	SegmentStatus string `xml:"segmentStatus,attr"`
	Flight        Flight //this is needed even for hotel/car
	FareBasis     FareBasis
	NotValidAfter NotValidAfter
	Baggage       Baggage
}
type FareIndicator struct {
	//XMLName       xml.Name `xml:"FareIndicator"`
}
type BaseFare struct {
	//XMLName       xml.Name `xml:"BaseFare"`
	CurrencyCode string `xml:"currencyCode,attr"`
	Val          string `xml:",chardata"`
}
type TotalTax struct {
	//XMLName       xml.Name `xml:"TotalTax"`
	CurrencyCode string `xml:"currencyCode,attr"`
	Val          string `xml:",chardata"`
}
type TotalFare struct {
	//XMLName       xml.Name `xml:"TotalFare"`
	CurrencyCode string `xml:"currencyCode,attr"`
	Val          string `xml:",chardata"`
}
type CombinedTax struct {
	//XMLName       xml.Name `xml:"CombinedTax"`
	Code   string `xml:"code,attr"`
	Amount Amount
}
type Tax struct {
	//XMLName       xml.Name `xml:"CombinedTax"`
	Code   string `xml:"code,attr"`
	Amount Amount
}
type TaxInfo struct {
	//XMLName       xml.Name `xml:"TaxInfo"`
	CombinedTaxes []CombinedTax `xml:"CombinedTax"`
	Taxes         []Tax         `xml:"Tax"`
}
type FareCalculation struct {
	//XMLName       xml.Name `xml:"FareCalculation"`
	Val string `xml:",chardata"`
}
type SegmentNumber struct {
	//XMLName xml.Name `xml:"SegmentNumber"`
	Val string `xml:",chardata"`
}
type FlightSegmentNumbers struct {
	//XMLName       xml.Name `xml:"FlightSegmentNumbers"`
	SegmentNumbers []SegmentNumber `xml:"SegmentNumber"`
}
type FareDirectionality struct {
	//XMLName       xml.Name `xml:"FareDirectionality"`
	OneWay bool `xml:"oneWay,attr"`
}
type GoverningCarrier struct {
	//XMLName       xml.Name `xml:"GoverningCarrier"`
	Val string `xml:",chardata"`
}
type FareComponent struct {
	//XMLName       xml.Name `xml:"FareComponent"`
	FareBasisCode        string `xml:"fareBasisCode,attr"`
	Number               string `xml:"number,attr"`
	FlightSegmentNumbers FlightSegmentNumbers
	FareDirectionality   FareDirectionality
	Departure            Departure
	Arrival              Arrival
	Amount               Amount
	GoverningCarrier     GoverningCarrier
}
type FareInfo struct {
	//XMLName       xml.Name `xml:"FareInfo"`
	FareIndicators  FareIndicator
	BaseFare        BaseFare
	TotalTax        TotalTax
	TotalFare       TotalFare
	TaxInfo         TaxInfo
	FareCalculation FareCalculation
	FareComponent   FareComponent
}
type OBFee struct {
	//XMLName       xml.Name `xml:"OBFee"`
	Code        string `xml:"code"`
	OBFType     string `xml:"type"`
	Amount      Amount
	Description struct {
		Val string `xml:",chardata"`
	} `xml:"Description"`
}
type FeeInfo struct {
	//XMLName       xml.Name `xml:"FeeInfo"`
	OBFee OBFee
}
type MiscellaneousInfo struct {
	//XMLName       xml.Name `xml:"MiscellaneousInfo"`
	ValidatingCarrier ValidatingCarrier
	ItineraryType     ItineraryType
}
type ReservationMessage struct {
	//XMLName xml.Name `xml:"Message"`
	//Code   string `xml:"code,attr"`
	Number string `xml:"number,attr"`
	Mtype  string `xml:"type,attr"`
	Val    string `xml:",chardata"`
}
type MessageInfo struct {
	//XMLName       xml.Name `xml:"MessageInfo"`
	Messages []ReservationMessage `xml:"Message"`
}
type DetailsPriceQuoteElem struct {
	//XMLName       xml.Name `xml:"Details"`
	Number              string `xml:"number,attr"`
	PassengerType       string `xml:"passengerType,attr"`
	PricingType         string `xml:"pricingType,attr"`
	Status              string `xml:"status,attr"`
	Dtype               string `xml:"type,attr"`
	AgentInfo           AgentInfo
	TransactionInfo     TransactionInfo
	NameAssociationInfo NameAssociationInfo
	SegmentInfo         SegmentInfo
	FareInfo            FareInfo
	FeeInfo             FeeInfo
	MiscellaneousInfo   MiscellaneousInfo
	MessageInfo         MessageInfo
}

type Indicators struct {
	//XMLName    xml.Name `xml:"Indicators"`
	ItineraryChange bool `xml:"itineraryChange,attr"`
}

type ItineraryType struct {
	//XMLName    xml.Name `xml:"ItineraryType"`
	Val string `xml:",chardata"`
}
type Amount struct {
	//XMLName    xml.Name `xml:"Amount"`
	CurrencyCode string `xml:"currencyCode,attr"`
	DecimalPlace string `xml:"decimalPlace,attr"`
	Val          string `xml:",chardata"`
}
type Waiver struct {
	//XMLName    xml.Name `xml:"Waiver"`
	Val string `xml:",chardata"`
}
type Fee struct {
	//XMLName    xml.Name `xml:"Fee"`
	Code   string `xml:"code,attr"`
	ItemID string `xml:"itemId,attr"`
	Ftype  string `xml:"type,attr"`
	Amount Amount
	Waiver Waiver
}
type ValidatingCarrier struct {
	//XMLName    xml.Name `xml:"ValidatingCarrier"`
	Val string `xml:",chardata"`
}
type Total struct {
	//XMLName    xml.Name `xml:"Total"`
	CurrencyCode string `xml:"currencyCode,attr"`
	Val          string `xml:",chardata"`
}
type Amounts struct {
	//XMLName    xml.Name `xml:"Amounts"`
	Total Total
}
type PassengerPriceQuoteElem struct {
	XMLName            xml.Name `xml:"Passenger"`
	PassengerTypeCount string   `xml:"passengerTypeCount,attr"`
	RequestedType      string   `xml:"requestedType,attr"`
	PType              string   `xml:"type,attr"`
}
type PriceQuoteNameAssocElem struct {
	//XMLName    xml.Name `xml:"PriceQuote"`
	Number              string `xml:"number,attr"`
	PricingType         string `xml:"pricingType,attr"`
	Status              string `xml:"status,attr"`
	PQNAType            string `xml:"type,attr"`
	Indicators          Indicators
	Passenger           PassengerPriceQuoteElem
	ItineraryType       ItineraryType
	Fee                 Fee
	ValidatingCarrier   ValidatingCarrier
	Amounts             Amounts
	LocalCreateDateTime string `xml:"LocalCreateDateTime"`
}
type NameAssociation struct {
	//XMLName    xml.Name `xml:"NameAssociation"`
	FirstName  string                  `xml:"firstName,attr"`
	LastName   string                  `xml:"lastName,attr"`
	NameID     string                  `xml:"nameId,attr"`
	NameNumber string                  `xml:"nameNumber,attr"`
	PriceQuote PriceQuoteNameAssocElem `xml:"PriceQuote"`
}
type SummaryPriceQuoteElem struct {
	//XMLName          xml.Name          `xml:"Summary"`
	NameAssociations []NameAssociation `xml:"NameAssociation"`
}
type ReservationPriceQuoteElem struct {
	//XMLName     xml.Name `xml:"Reservation"`
	Val         string `xml:",chardata"`
	UpdateToken string `xml:"updateToken,attr"`
}
type PriceQuoteInfo struct {
	//XMLName     xml.Name `xml:"PriceQuoteInfo"`
	XMLNS       string `xml:"xmlns,attr"` //"http://www.sabre.com/ns/Ticketing/pqs/1.0"
	Reservation ReservationPriceQuoteElem
	Summary     SummaryPriceQuoteElem
	Details     []DetailsPriceQuoteElem
}
type PriceQuote struct {
	//XMLName        xml.Name `xml:"PriceQuote"`
	PriceQuoteInfo PriceQuoteInfo
}
