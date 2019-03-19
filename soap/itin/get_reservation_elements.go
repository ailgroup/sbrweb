package itin

import "encoding/xml"

type BookingDetails struct {
	XMLName                 xml.Name `xml:"BookingDetails"`
	RecordLocator           string   `xml:"RecordLocator"`
	CreationTimestamp       string   `xml:"CreationTimestamp"`
	SystemCreationTimestamp string   `xml:"SystemCreationTimestamp"`
	CreationAgentID         string   `xml:"CreationAgentID"`
	UpdateTimestamp         string   `xml:"UpdateTimestamp"`
	PNRSequence             string   `xml:"PNRSequence"` //string easier w/ 0-th index
	DivideSplitDetails      string   `xml:"DivideSplitDetails"`
	EstimatedPurgeTimestamp string   `xml:"EstimatedPurgeTimestamp"`
	UpdateToken             string   `xml:"UpdateToken"`
}
type SourceReservationElem struct {
	XMLName        xml.Name `xml:"Source"`
	BookingSource  string   `xml:"BookingSource,attr"`
	AgentSine      string   `xml:"AgentSine,attr"`
	PseudoCityCode string   `xml:"PseudoCityCode,attr"`
	ISOCountry     string   `xml:"ISOCountry,attr"`
	AgentDutyCode  string   `xml:"AgentDutyCode,attr"`
	//AirlineVendorID   bool     `xml:"AirlineVendorID,attr"`
	HomePseudoCityCode string `xml:"HomePseudoCityCode,attr"`
}
type POSReservationElem struct {
	XMLName     xml.Name `xml:"POS"`
	AirExtras   bool     `xml:"AirExtras,attr"`
	InhibitCode string   `xml:"InhibitCode,attr"`
	Source      SourceReservationElem
}
type Passenger struct {
	XMLName         xml.Name `xml:"Passenger"`
	ID              string   `xml:"id,attr"`
	NameType        string   `xml:"nameType,attr"`
	PassengerType   string   `xml:"passengerType,attr"`
	ReferenceNumber string   `xml:"referenceNumber,attr"`
	NameID          string   `xml:"nameId,attr"`
	NameAssocID     string   `xml:"nameAssocId,attr"`
	ElementID       string   `xml:"elementId,attr"`
	LastName        string   `xml:"LastName"`
	FirstName       string   `xml:"FirstName"`
	//Seats              string `xml:"Seats"`
}

type PassengerReservation struct {
	XMLName          xml.Name           `xml:"PassengerReservation"`
	Passengers       []Passenger        `xml:"Passengers>Passenger"`
	Segments         SegmentReservation `xml:"Segments>Segment"`
	TicketingInfo    string             `xml:"TicketingInfo"`
	ItineraryPricing string             `xml:"ItineraryPricing"`
}
type ReceivedFrom struct {
	XMLName xml.Name `xml:"ReceivedFrom"`
	Name    string   `xml:"Name"`
}
type AddressLineResElem struct {
	XMLName xml.Name `xml:"AddressLine"`
	ID      string   `xml:"id,attr"`
	Atype   string   `xml:"type,attr"`
	Text    string   `xml:"Text"`
}
type AddressReservationElem struct {
	XMLName      xml.Name           `xml:"Address"`
	AddressLines AddressLineResElem `xml:"AddressLines>AddressLine"`
}
type PhoneNumberReservationElem struct {
	XMLName   xml.Name `xml:"PhoneNumber"`
	ID        string   `xml:"id,attr"`
	Index     string   `xml:"index,attr"`
	ElementID string   `xml:"elementId,attr"`
	CityCode  string   `xml:"CityCode"`
	Number    string   `xml:"Number"`
}
type EmailReservationElem struct {
	XMLName   xml.Name `xml:"Email"`
	ID        string   `xml:"id,attr"`
	Index     string   `xml:"index,attr"`
	ElementID string   `xml:"elementId,attr"`
	Text      string   `xml:"Text"`
}
type GenericSpecialRequests struct {
	XMLName     xml.Name `xml:"GenericSpecialRequests"`
	ID          string   `xml:"id,attr"`
	GType       string   `xml:"type,attr"`
	MsgType     string   `xml:"msgType,attr"`
	FreeText    string   `xml:"FreeText"`
	AirlineCode string   `xml:"AirlineCode"`
	FullText    string   `xml:"FullText"`
}
type AssociationMatrix struct {
	XMLName xml.Name `xml:"AssociationMatrix"`
	Name    string   `xml:"Name"`
	Parent  struct {
		Ref string `xml:"ref,attr"`
	} `xml:"Parent"`
	Child struct {
		Ref string `xml:"ref,attr"`
	} `xml:"Child"`
}
type ServiceRequestOpenRes struct {
	XMLName     xml.Name `xml:"ServiceRequest"`
	AirlineCode string   `xml:"airlineCode,attr"`
	ServiceType string   `xml:"serviceType,attr"`
	SsrType     string   `xml:"ssrType,attr"`
	FreeText    string   `xml:"FreeText"`
	FullText    string   `xml:"FullText"`
}
type OpenReservationElement struct {
	XMLName        xml.Name              `xml:"OpenReservationElement"`
	ID             string                `xml:"id,attr"`
	SType          string                `xml:"type,attr"`
	ElementID      string                `xml:"elementId,attr"`
	ServiceRequest ServiceRequestOpenRes `xml:"ServiceRequest"`
}
type Reservation struct {
	XMLName                 xml.Name `xml:"Reservation"`
	NumberInParty           string   `xml:"numberInParty,attr"`
	NumberOfInfants         string   `xml:"numberOfInfants,attr"`
	NumberInSegment         string   `xml:"numberInSegment,attr"`
	BookingDetails          BookingDetails
	POS                     POSReservationElem
	PassengerReservation    PassengerReservation
	ReceivedFrom            ReceivedFrom
	Addresses               []AddressReservationElem     `xml:"Addresses>Address"`
	PhoneNumbers            []PhoneNumberReservationElem `xml:"PhoneNumbers>PhoneNumber"`
	EmailAddresses          []EmailReservationElem       `xml:"EmailAddresses>Email"`
	GenericSpecialRequests  []GenericSpecialRequests
	AssociationMatrices     []AssociationMatrix      `xml:"AssociationMatrices>AssociationMatrix"`
	OpenReservationElements []OpenReservationElement `xml:"OpenReservationElements>OpenReservationElement"`
}
