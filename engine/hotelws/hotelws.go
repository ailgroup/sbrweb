/* Package hotelws (Hotel Web Services) implements Sabre SOAP hotel booking through availability, property and rate descriptions, passenger details (PNR), reservation, and transaction Web Services. It handles exclusively the XML-based Simple Object Access Protocol endpoints the Sabre supports. Look elsewhere for support of Sabre rest services.

# Typical Workflow definition from Sabre Web Services
	Book Hotel Reservation
	The following workflow allows you to search and book a hotel room.
	Steps

	--Step 1: Retrieve hotel availability using OTA_HotelAvailLLSRQ.
	--Step 2: Retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
	--Step 3: Retrieve hotel rules and policies using HotelRateDescriptionLLSRQ.*
	--Step 4: Add any additional (required) information to create the passenger name record (PNR) using PassengerDetailsRQ.**
	--Step 5: Book a room for the selected hotel using OTA_HotelResLLSRQ.
	--Step 6: End the transaction of the passenger name record using EndTransactionLLSRQ.
	Note

	* Mandatory only if selected option in response of HotelPropertyDescriptionLLSRQ contains HRD_RequiredForSell="true".
	** Ensure Agency address is added within call to PassengerDetails, so as the OTA_HotelResLLSRQ call is not rejected.

One may implement Sabre hotel searching through building various criteria functions with proper criterion types. Many criterion exist that are not yet implemented: (Award, ContactNumbers, CommissionProgram, HotelAmenity, PointOfInterest, RefPoint, RoomAmenity, HotelFeaturesCriterion). To add more criterion create a criterion type (e.g, XCriterion) as well as its accompanying function to handle the data params (e.g., XSearch); see functions in hotel_search_criteria.go and for types look in this file (e.g., HotelSearchCriteria, Criterion, ...).
*/
package hotelws

import (
	b64 "encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

const (
	hotelRQVersion        = "2.3.0"
	TimeFormatMD          = "01-02"
	TimeFormatMDTHM       = "01-02T15:04"
	TimeFormatMDHM        = "01-02 15:04"
	StreetQueryField      = "street_qf"
	CityQueryField        = "city_qf"
	PostalQueryField      = "postal_qf"
	CountryCodeQueryField = "countryCode_qf"
	LatlngQueryField      = "latlng_qf"
	HotelidQueryField     = "hotelID_qf"
	TrackEncDelimiter     = "|"
	TrackEncIndex         = "idx"
	TrackEncRPH           = "rph"
	TrackEncIATAChar      = "iatachar"
	TrackEncTotal         = "total"
	returnHostCommand     = true
	ESA                   = "\u0087" //UNICODE: End of Selected Area
	CrossLorraine         = "\u2628" //UNICODE Cross of Lorraine
)

var hostCommandReplacer = strings.NewReplacer("\\", "", "/", "", ESA, "")

// B64Enc  base64 encode a string
func B64Enc(str string) string {
	return b64.URLEncoding.EncodeToString([]byte(str))
}

// B64Dec decode a base64 string
func B64Dec(b64str string) (string, error) {
	uDec, err := b64.URLEncoding.DecodeString(b64str)
	return string(uDec), err
}

// TimeSpanFormatter parse string data value into time value.
func TimeSpanFormatter(arrive, depart, formIn, formOut string) TimeSpan {
	a, _ := time.Parse(formIn, arrive)
	d, _ := time.Parse(formIn, depart)
	return TimeSpan{
		Depart: d.Format(formOut),
		Arrive: a.Format(formOut),
	}
}

// sanatize cleans up filtered string terms for file names. removes whitespace and slashes as these either get in the
func sanatize(str string) string {
	// use a NewReplacer to clean this up and make it easy to use with multiple replacers
	//trim := strings.Trim(str, " ")
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, hostCommandReplacer.Replace(strings.ToLower(str)))
}

func (s SystemResults) Translate() string {
	clean := sanatize(s.Message)
	switch {
	case strings.Contains(clean, "ckdate"):
		return fmt.Sprintf("%s=%s", "Check Date Parameters", s.Message)
	default:
		return fmt.Sprintf("%s=%s", s.Message, "No Translation")
	}
}

/* OTHER ERRORS

--see hotel_res_direct_connect.xml, credit card???
	<stl:Message>      ** - DIRECT CONNECT NOT PROCESSED - **</stl:Message>

--see hotel_res_rq_not_proc_format.xml, too many options??
	<stl:Message code="0">FORMAT</stl:Message>

--see...
	<stl:Message>INVALID CARD NUMBER</stl:Message> card number/code no match, other invalid card number reasons

*/
func (result ApplicationResults) ErrFormat() sbrerr.ErrorSabreResult {
	return sbrerr.ErrorSabreResult{
		Code: sbrerr.SabreEngineStatusCode(result.Status),
		AppMessage: fmt.Sprintf(
			"%s because %s. %s. HostCommand[LNIATA: %s Cryptic: %s]",
			result.Status,
			result.Error.Type,
			result.Error.System.Translate(),
			result.Error.System.HostCommand.LNIATA,
			result.Error.System.HostCommand.Cryptic,
		),
	}
}

func (result ApplicationResults) Ok() bool {
	switch result.Status {
	case sbrerr.StatusNotProcess(): //queries
		return false
	case sbrerr.StatusApproved(): //sessions
		return true
	case sbrerr.StatusComplete(): //queries, pnr
		return true
	default:
		return true
	}
}

// HotelSearchCriteria top level element for criterion
type HotelSearchCriteria struct {
	XMLName   xml.Name `xml:"HotelSearchCriteria"`
	Criterion Criterion
}

// QuerySearchParams is a typed function to support optional query params on creation of new search criterion
type QuerySearchParams func(*HotelSearchCriteria) error

// HotelRefCriterion map of hotel ref criteria
type HotelRefCriterion map[string][]string

// AddressCriterion map of address search criteria
type AddressCriterion map[string]string

// PropertyTypeCriterion slice of property type strings (APTS, LUXRY)
type PropertyTypeCriterion []string

// PackageCriterion slice of property type strings (GF, HM, BB)
type PackageCriterion []string

// HotelRef contains any number of search criteria under the HotelRef element.
type HotelRef struct {
	XMLName       xml.Name `xml:"HotelRef,omitempty"`
	HotelCityCode string   `xml:"HotelCityCode,attr,omitempty"`
	HotelCode     string   `xml:"HotelCode,attr,omitempty"`
	HotelName     string   `xml:"HotelName,attr,omitempty"`
	Latitude      string   `xml:"Latitude,attr,omitempty"`
	Longitude     string   `xml:"Longitude,attr,omitempty"`
}

// PropertyType container for searhing types of properties (APTS, LUXRY...)
type PropertyType struct {
	Val string `xml:",chardata"`
}

// Package container for searching types of packages (GF, HM, BB...)
type Package struct {
	Val string `xml:",chardata"`
}

// Address represents typical building addresses
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNumber,omitempty"`
	City          string `xml:"CityName,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr"`
	} `xml:"StateCountyProv,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
}

// Criterion holds various serach criteria
type Criterion struct {
	XMLName       xml.Name `xml:"Criterion"`
	HotelRefs     []*HotelRef
	Address       *Address
	PropertyTypes []*PropertyType
	Packages      []*Package
}

// GuestCounts how many guests per night-room. TODO: check on Sabre validation limits (think it is 4)
type GuestCounts struct {
	XMLName xml.Name `xml:"GuestCounts"`
	Count   int      `xml:"Count,attr"`
}

// Customer for corporate or typical sabre customer ids
type Customer struct {
	XMLName    xml.Name    `xml:"Customer,omitempty"`
	NameNumber string      `xml:"NameNumber,attr,omitempty"`
	Corporate  *Corporate  //nil pointer ignored if empty
	CustomerID *CustomerID //nil pointer ignored if empty
}

// CustomerID number
type CustomerID struct {
	XMLName xml.Name `xml:"ID"`
	Number  string   `xml:"Number,omitempty"`
}

// Corporate customer id
type Corporate struct {
	XMLName xml.Name `xml:"Corporate"`
	ID      string   `xml:"ID,omitempty"`
}

// Timepsan for arrival and departure params
type TimeSpan struct {
	XMLName  xml.Name `xml:"TimeSpan" json:"-"`
	Duration int      `xml:"Duration,omitempty"`
	Depart   string   `xml:"End,attr,omitempty"`
	Arrive   string   `xml:"Start,attr,omitempty"`
}

type RatePlan struct {
	XMLName         xml.Name `xml:"RatePlanCandidate" json:"-"`
	CurrencyCode    string   `xml:"CurrencyCode,attr,omitempty"`
	DCA_ProductCode string   `xml:"DCA_ProductCode,attr,omitempty"`
	DecodeAll       string   `xml:"DecodeAll,attr,omitempty"`
	RateCode        string   `xml:"RateCode,attr,omitempty"`
	RPH             string   `xml:"RPH,attr,omitempty"`
}

// RatePlanCandidates determines types of rates queried
type RatePlanCandidates struct {
	XMLName   xml.Name `xml:"RatePlanCandidates"`
	RatePlans []*RatePlan
}

// AvailAvailRequestSegment holds basic hotel availability params: customer ids, guest count, hotel search criteria and arrival departure. nil pointers ignored if empty
type AvailRequestSegment struct {
	XMLName             xml.Name `xml:"AvailRequestSegment"`
	Customer            *Customer
	GuestCounts         *GuestCounts
	HotelSearchCriteria *HotelSearchCriteria
	RatePlanCandidates  *RatePlanCandidates
	TimeSpan            *TimeSpan
}

// RoomStay contains all info relevant to the property's available rooms. It is the root-level element after service element for hotel_rate_desc and hotel_property_desc.
type RoomStay struct {
	XMLName           xml.Name `xml:"RoomStay" json:"-"`
	BasicPropertyInfo BasicPropertyInfo
	Guarantee         Guarantee
	RoomRates         []RoomRate `xml:"RoomRates>RoomRate"`
	TimeSpan          TimeSpan
}

// Guarantee shows forms of payment accepted by property
type Guarantee struct {
	XMLName            xml.Name `xml:"Guarantee" json:"-"`
	GuaranteesAccepted GuaranteesAccepted
	DepositsAccepted   DepositsAccepted
}
type GuaranteesAccepted struct {
	XMLName      xml.Name      `xml:"GuaranteesAccepted" json:"-"`
	PaymentCards []PaymentCard `xml:"PaymentCard"`
}
type DepositsAccepted struct {
	XMLName      xml.Name      `xml:"DepositsAccepted" json:"-"`
	PaymentCards []PaymentCard `xml:"PaymentCard"`
}
type PaymentCard struct {
	Code       string `xml:"Code,attr,omitempty"`
	Type       string `xml:"Type,attr,omitempty"`
	ExpireDate string `xml:"ExpireDate,attr,omitempty"`
	Number     string `xml:"Number,attr,omitempty"`
}

type RoomRate struct {
	XMLName            xml.Name `xml:"RoomRate" json:"-"`
	DirectConnect      string   `xml:"RDirectConnect,attr"`
	GuaranteeSurcharge string   `xml:"GuaranteeSurchargeRequired,attr"`
	GuaranteedRate     string   `xml:"GuaranteedRateProgram,attr"`
	IATA_Character     string   `xml:"IATA_CharacteristicIdentification,attr"`
	IATA_Product       string   `xml:"IATA_ProductIdentification,attr"`
	LowInventory       string   `xml:"LowInventoryThreshold,attr"`
	RateLevelCode      string   `xml:"RateLevelCode,attr"`
	RPH                string   `xml:"RPH,attr"`
	RateChangeInd      string   `xml:"RateChangeInd,attr"`
	RateConversionInd  string   `xml:"RateConversionInd,attr"`
	SpecialOffer       string   `xml:"SpecialOffer,attr"`
	Rates              []Rate   `xml:"Rates>Rate"`
	AdditionalInfo     AdditionalInfo
	HotelRateCode      string `xml:"HotelRateCode"`
	TrackedEncoding    string `json:"tracked_encoding"`
}

func (r *RoomRate) DecodeTrackedEncoding() ([]string, error) {
	res := []string{}
	bytEnc, err := B64Dec(r.TrackedEncoding)
	if err != nil {
		return res, err
	}
	res = strings.Split(string(bytEnc), TrackEncDelimiter)
	return res, nil
}

type AdditionalInfo struct {
	XMLName    xml.Name `xml:"AdditionalInfo" json:"-"`
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

// BasicPropertyInfo contains all info relevant to property. It is the root-level element after service element for hotel_avail; and is embedded in the RoomStay element.
type BasicPropertyInfo struct {
	XMLName         xml.Name `xml:"BasicPropertyInfo" json:"-"`
	AreadID         string   `xml:"AreadID,attr"`
	ChainCode       string   `xml:"ChainCode,attr"`
	Distance        string   `xml:"Distance,attr"`
	GEO_ConfAvail   string   `xml:"GEO_ConfidenceLevel,attr"` //hotel avail
	GeoConfPropDesc string   `xml:"GeoConfidenceLevel,attr"`  //property description
	HotelCityCode   string   `xml:"HotelCityCode,attr"`
	HotelCode       string   `xml:"HotelCode,attr"`
	HotelName       string   `xml:"HotelName,attr"`
	CancelPenalty   struct {
		PolicyCode string `xml:"PolicyCode"`
	} `xml:"CancelPenalty"`
	Latitude  string `xml:"Latitude,attr"`
	Longitude string `xml:"Longitude,attr"`
	Address   struct {
		Line []string `xml:"AddressLine"`
	} `xml:"Address"`
	Awards struct {
		AwardProvider string `xml:"AwardProvider"`
	} `xml:"Awards"`
	CheckIn        string `xml:"CheckInTime"`
	CheckOut       string `xml:"CheckOutTime"`
	ContactNumbers struct {
		Number ContactNumber `xml:"ContactNumber"`
	} `xml:"ContactNumbers"`
	DirectCon DirectConnect `xml:"DirectConnect"`
	LocDesc   struct {
		Text string `xml:"Text"`
	} `xml:"LocationDescription"`
	IndexDatum []IndexD `xml:"IndexData>Index"`
	Prop       struct {
		Rating string `xml:"Rating,attr"`
		Text   string `xml:"Text"`
	} `xml:"Property"`
	PropertyOptionInfo PropertyOptionInfo `xml:"PropertyOptionInfo"`
	PropertyTypeInfo   PropertyTypeInfo   `xml:"PropertyTypeInfo"`
	SpecialOffers      struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"SpecialOffers"`
	Taxes struct {
		Text []string `xml:"Text"`
	} `xml:"Taxes"`
	VendorMessages VendorMessages
	RateRange      struct {
		CurrencyCode string `xml:"CurrencyCode,attr"`
		Max          string `xml:"Max,attr"`
		Min          string `xml:"Min,attr"`
	} `xml:"RateRange"`
	RoomRateAvail RoomRate //hotel avail
}
type Charge struct {
	//XMLName     xml.Name `xml:"Charges"`
	Crib        string `xml:"Crib,attr"`
	ExtraPerson string `xml:"ExtraPerson,attr"`
}

type AdditionalGuestAmount struct {
	XMLName        xml.Name `xml:"AdditionalGuestAmount" json:"-"`
	MaxExtraPerson int      `xml:"MaxExtraPersonsAllowed,attr"`
	NumCribs       int      `xml:"NumCribs,attr"`
	Charges        []Charge `xml:"Charges"`
}

type TotalSurcharges struct {
	XMLName xml.Name `xml:"TotalSurcharges" json:"-"`
	Amount  string   `xml:"Amount,attr"`
}
type TotalTaxes struct {
	XMLName     xml.Name `xml:"TotalTaxes" json:"-"`
	Amount      string   `xml:"Amount,attr"`
	TaxFieldOne string   `xml:"TaxFieldOne"`
	TaxFieldTwo string   `xml:"TaxFieldTwo"`
	Text        []string `xml:"Text"`
}

type HotelPricing struct {
	XMLName         xml.Name `xml:"HotelTotalPricing" json:"-"`
	Amount          string   `xml:"Amount,attr"`
	Disclaimer      string   `xml:"Disclaimer"`
	TotalSurcharges TotalSurcharges
	TotalTaxes      TotalTaxes
}

type Rate struct {
	XMLName                xml.Name                `xml:"Rate" json:"-"`
	Amount                 string                  `xml:"Amount,attr"`
	ChangeIndicator        string                  `xml:"ChangeIndicator,attr"`
	CurrencyCode           string                  `xml:"CurrencyCode,attr"`
	HRD_RequiredForSell    string                  `xml:"HRD_RequiredForSell,attr"`
	PackageIndicator       string                  `xml:"PackageIndicator,attr"`
	RateConversionInd      string                  `xml:"RateConversionInd,attr"`
	ReturnOfRateInd        string                  `xml:"ReturnOfRateInd,attr"`
	RoomOnRequest          string                  `xml:"RoomOnRequest,attr"`
	AdditionalGuestAmounts []AdditionalGuestAmount `xml:"AdditionalGuestAmounts>AdditionalGuestAmount"`
	HotelPricing           HotelPricing
}

type VendorMessages struct {
	XMLName               xml.Name              `xml:"VendorMessages" json:"-"`
	Attractions           Attractions           `xml:"Attractions"`
	AdditionalAttractions AdditionalAttractions `xml:"AdditionalAttractions"`
	Awards                Awards                `xml:"Awards"`
	Cancellation          Cancellation          `xml:"Cancellation"`
	Deposit               Deposit               `xml:"Deposit"`
	Description           Description           `xnl:"Description"`
	Dining                Dining                `xml:"Dining"`
	Directions            Directions            `xml:"Directions"`
	Facilities            Facilities            `xml:"Facilities"`
	Guarantee             VendorGuarantee       `xml:"Guarantee"`
	Location              Location              `xml:"Location"`
	MarketingInformation  MarketingInformation  `xml:"MarketingInformation"`
	MiscServices          MiscServices          `xml:"MiscServices"`
	Policies              Policies              `xml:"Policies"`
	Recreation            Recreation            `xml:"Recreation"`
	Rooms                 Rooms                 `xml:"Rooms"`
	Safety                Safety                `xml:"Safety"`
	Services              Services              `xml:"Services"`
	Transportation        Transportation        `xml:"Transportation"`
}

type Transportation struct {
	Text []string `xml:"Text"`
}
type Services struct {
	Text []string `xml:"Text"`
}
type Safety struct {
	Text []string `xml:"Text"`
}
type Rooms struct {
	Text []string `xml:"Text"`
}
type Recreation struct {
	Text []string `xml:"Text"`
}
type Policies struct {
	Text []string `xml:"Text"`
}
type MiscServices struct {
	Text []string `xml:"Text"`
}
type MarketingInformation struct {
	Text []string `xml:"Text"`
}
type Location struct {
	Text []string `xml:"Text"`
}
type VendorGuarantee struct {
	//this is Guarantee under BasciPropertyInfo.VendorMessages
	//not Guarantee element in RoomStay
	Text []string `xml:"Text"`
}
type Facilities struct {
	Text []string `xml:"Text"`
}
type Directions struct {
	Text []string `xml:"Text"`
}
type Dining struct {
	Text []string `xml:"Text"`
}
type Description struct {
	Text []string `xml:"Text"`
}
type Deposit struct {
	Text []string `xml:"Text"`
}
type Cancellation struct {
	Text []string `xml:"Text"`
}
type Awards struct {
	Text []string `xml:"Text"`
}
type Attractions struct {
	Text []string `xml:"Text"`
}
type AdditionalAttractions struct {
	Text []string `xml:"Text"`
}

type IndexD struct {
	XMLName            xml.Name `xml:"Index" json:"-"`
	CountryState       string   `xml:"CountryState,attr"`
	DistanceDirection  string   `xml:"DistanceDirection,attr"`
	LocationCode       string   `xml:"LocationCode,attr"`
	Point              string   `xml:"Point,attr"`
	TransportationCode string   `xml:"TransportationCode,attr"`
}

type DirectConnect struct {
	AltAvail struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Alt_Avail"`
	DCAvailPart struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"DC_AvailParticipant"`
	DCSellPart struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"DC_SellParticipant"`
	RatesExceedMax struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RatesExceedMax"`
	UnAvail struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"UnAvail"`
}

type ContactNumber struct {
	Fax   string `xml:"Fax,attr"`
	Phone string `xml:"Phone,attr"`
}

type ApplicationResults struct {
	XMLName xml.Name   `xml:"ApplicationResults" json:"-"`
	Status  string     `xml:"status,attr"`
	Success ReqSuccess `xml:"Success"`
	Error   ReqError   `xml:"Error"`
}

type ReqError struct {
	Type      string        `xml:"type,attr"`
	Timestamp string        `xml:"timeStamp,attr"`
	System    SystemResults `xml:"SystemSpecificResults"`
}

type ReqSuccess struct {
	Timestamp string        `xml:"timeStamp,attr"`
	System    SystemResults `xml:"SystemSpecificResults"`
}

type SystemResults struct {
	HostCommand HostCommand `xml:"HostCommand"`
	Message     string      `xml:"Message,omitempty"`
}

type HostCommand struct {
	LNIATA  string `xml:"LNIATA,attr,omitempty"`
	Cryptic string `xml:",chardata"`
}

type PropertyTypeInfo struct {
	AllInclusive struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"AllInclusive"`
	Apartments struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Apartments"`
	BedBreakfast struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"BedBreakfast"`
	Castle struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Castle"`
	Conventions struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Conventions"`
	Economy struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Economy"`
	ExtendedStay struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"ExtendedStay"`
	Farm struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Farm"`
	Luxury struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Luxury"`
	Moderate struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Moderate"`
	Motel struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Motel"`
	Resort struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Resort"`
	Suites struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Suites"`
}

type PropertyOptionInfo struct {
	ADA_Accessible struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"ADA_Accessible"`
	AdultsOnly struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"AdultsOnlyA"`
	BeachFront struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"BeachFront"`
	Breakfast struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Breakfast"`
	BusinessCenter struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"BusinessCenter"`
	BusinessReady struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"BusinessReadyA"`
	Conventions struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Conventions"`
	Dataport struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Dataport"`
	Dining struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Dining"`
	DryClean struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"DryClean"`
	EcoCertified struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"EcoCertified"`
	ExecutiveFloors struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"ExecutiveFloorsA"`
	FitnessCenter struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FitnessCenter"`
	FreeLocalCalls struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeLocalCalls"`
	FreeParking struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeParking"`
	FreeShuttle struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeShuttle"`
	FreeWifiInMeetingRooms struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeWifiInMeetingRooms"`
	FreeWifiInPublicSpaces struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeWifiInPublicSpaces"`
	FreeWifiInRooms struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"FreeWifiInRooms"`
	GameFacilities struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"GameFacilities"`
	Golf struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Golf"`
	HighSpeedInternet struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"HighSpeedInternet"`
	HypoallergenicRooms struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"HypoallergenicRooms"`
	IndoorPool struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"IndoorPool"`
	InRoomCoffeeTea struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"InRoomCoffeeTea"`
	InRoomMiniBar struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"InRoomMiniBar"`
	InRoomRefrigerator struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"InRoomRefrigerator"`
	InRoomSafe struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"InRoomSafe"`
	InteriorDoorways struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"InteriorDoorways"`
	Jacuzzi struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Jacuzzi"`
	KidsFacilities struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"KidsFacilities"`
	KitchenFacilities struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"KitchenFacilities"`
	MealService struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"MealService"`
	MeetingFacilities struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"MeetingFacilities"`
	NoAdultTV struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"NoAdultTV"`
	NonSmoking struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"NonSmoking"`
	OutdoorPool struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"OutdoorPool"`
	Pets struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Pets"`
	Pool struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Pool"`
	PublicTransportationAdjacent struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"PublicTransportationAdjacent"`
	RateAssured struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RateAssured"`
	Recreation struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Recreation"`
	RestrictedRoomAccess struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RestrictedRoomAccess"`
	RoomService struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RoomService"`
	RoomService24Hours struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RoomService24Hours"`
	RoomsWithBalcony struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"RoomsWithBalcony"`
	SkiInOutProperty struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"SkiInOutProperty"`
	SmokeFree struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"SmokeFree"`
	SmokingRoomsAvail struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"SmokingRoomsAvail"`
	Tennis struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Tennis"`
	WaterPurificationSystem struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"WaterPurificationSystem"`
	Wheelchair struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"Wheelchair"`
}
