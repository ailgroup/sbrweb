/*
Package hotel (Hotel SOAP) implements Sabre SOAP hotel booking through availability, property and rate descriptions, passenger details (PNR), reservation, and transaction Web Services. It handles exclusively the XML-based Simple Object Access Protocol endpoints the Sabre supports. Look elsewhere for support of Sabre rest services.

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

One may implement Sabre hotel searching through building various criteria functions with proper criterion types.
Many criterion exist that are not yet implemented:
	* Award,
	* ContactNumbers,
	* CommissionProgram,
	* HotelAmenity,
	* RefPoint,
	* RoomAmenity,
	* HotelFeaturesCriterion

To add more criterion create a criterion type and function to handle the data params; see HotelSearchCriteria, Criterion and others.
*/
package htlsp

import (
	b64 "encoding/base64"
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

const (
	hotelRQVersion  = "2.3.0"
	TimeFormatMD    = "01-02"
	TimeFormatMDTHM = "01-02T15:04"
	TimeFormatMDHM  = "01-02 15:04"

	StreetQueryField      = "street_qf"
	CityQueryField        = "city_qf"
	StateCodeQueryField   = "stateCode_qf"
	CountryCodeQueryField = "countryCode_qf"
	POIQueryField         = "pOInterest_qf"
	HotelidQueryField     = "hotelID_qf"
	LatlngQueryField      = "latlng_qf"
	PostalQueryField      = "postal_qf"

	ColDelim    = ":"
	DashDelim   = "-"
	LBrackDelim = "]"
	PipeDelim   = "|"
	RBrackDelim = "["
	SColDelim   = ";"

	RoomMetaArvKey       = "arv"  //arrival
	RoomMetaDptKey       = "dpt"  //depart
	RoomMetaGstKey       = "gst"  //guest count
	RoomMetaHcKey        = "hc"   //host command
	RoomMetaHidKey       = "hid"  //hotel id
	RoomMetaRmtKey       = "rmt"  //room type
	RoomMetaRphKey       = "rph"  //reference place holder
	RoomMetaGuarenteeKey = "guar" //guarantee surchage
	RrateMetaAmtKey      = "amt"  //total
	RrateMetaCurKey      = "cur"  //total
	RrateMetaRqsKey      = "rqs"  //next

	returnHostCommand = true
)

var (
	hostCommandReplacer = strings.NewReplacer("\\", "", "/", "", srvc.ESA, "")
	ratesMetaMatch      = regexp.MustCompile(`^\[.*\]$`)
)

// TimeSpanFormatter parse string data value into time value.
func TimeSpanFormatter(arrive, depart, formIn, formOut string) TimeSpan {
	a, _ := time.Parse(formIn, arrive)
	d, _ := time.Parse(formIn, depart)
	return TimeSpan{
		Depart: d.Format(formOut),
		Arrive: a.Format(formOut),
	}
}

// sanatize cleans up filtered string terms for file names. removes whitespace and slashes
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

// Translate sabre SystemResults messages into human readable.
func (s SystemResults) Translate() string {
	clean := sanatize(s.Message)
	switch {
	case strings.Contains(clean, "ckdate"):
		return fmt.Sprintf("%s=%s", "Check Date Parameters", s.Message)
	case strings.Contains(clean, "noavail"):
		return fmt.Sprintf("%s=%s", "No Hotel Availability", s.Message)
	case strings.Contains(clean, "nomoredata"):
		return fmt.Sprintf("%s=%s", "No More Hotel Availability Data", s.Message)
	default:
		return fmt.Sprintf("%s=%s", s.Message, "No Translation")
	}
}

/*
ErrFormat formatter on ApplicationResults for printing error string from sabre soap calls.

	OTHER ERRORS TO BE FORMATTED
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

// Ok for ApplicationResults returns boolean check for sbrerr on SOAP requests
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

// B64Enc  base64 encode a string
func B64Enc(str string) string {
	return b64.URLEncoding.EncodeToString([]byte(str))
}

// B64Dec decode a base64 string
func B64Dec(b64str string) (string, error) {
	uDec, err := b64.URLEncoding.DecodeString(b64str)
	return string(uDec), err
}

// parseB64DecodeRates parses cached room stay rates
func (p *ParsedRoomMeta) parseB64DecodeRates() {
	for _, r := range p.StayRatesCache {
		mr := parsedStayRateCache{}
		dash := strings.Split(r, DashDelim)
		for _, d := range dash {
			rs := strings.Split(d, ColDelim)
			if len(rs) < 2 {
				continue
			}
			switch rs[0] {
			case RrateMetaAmtKey:
				mr.Amt = rs[1]
			case RrateMetaCurKey:
				mr.Cur = rs[1]
			case RrateMetaRqsKey:
				if rs[1] == "true" {
					mr.Rqs = true
				}
			}
		}
		p.ParsedStayRatesCache = append(p.ParsedStayRatesCache, mr)
	}
}

// NewParsedRoomMeta builds a struct from parsing b64 encoded string cache of previous rate request. See SetRoomMetaData for how this data is constucted.
func NewParsedRoomMeta(b64Str string) (ParsedRoomMeta, error) {
	rmp := ParsedRoomMeta{}
	b64Str, err := B64Dec(b64Str)
	if err != nil {
		return rmp, err
	}
	for _, p := range strings.Split(b64Str, PipeDelim) {
		if ratesMetaMatch.MatchString(p) {
			b := strings.TrimPrefix(p, RBrackDelim)
			b = strings.TrimSuffix(b, LBrackDelim)
			rmp.StayRatesCache = strings.Split(b, SColDelim)
		} else {
			c := strings.Split(p, ColDelim)
			if len(c) < 2 {
				continue
			}
			switch c[0] {
			case RoomMetaArvKey:
				rmp.Arrive = c[1]
			case RoomMetaDptKey:
				rmp.Depart = c[1]
			case RoomMetaGstKey:
				rmp.Guest = c[1]
			case RoomMetaHcKey:
				rmp.Hc = c[1]
			case RoomMetaHidKey:
				rmp.HotelID = c[1]
			case RoomMetaRphKey:
				rmp.Rph = c[1]
			case RoomMetaRmtKey:
				rmp.Rmt = c[1]
			case RoomMetaGuarenteeKey:
				rmp.GuaranteeSurcharge = c[1]
			}
		}
	}
	rmp.parseB64DecodeRates()
	return rmp, nil
}

type parsedStayRateCache struct {
	Amt string
	Cur string
	Rqs bool
}

type ParsedRoomMeta struct {
	Arrive               string   // arrival
	Depart               string   // departure
	Guest                string   // guest count
	Hc                   string   // host command
	HotelID              string   // hotel id
	Rmt                  string   // room type;;iata characteristic
	Rph                  string   // reference place holder
	GuaranteeSurcharge   string   //guarantee surcharge
	StayRatesCache       []string // room_stay.room_rates
	ParsedStayRatesCache []parsedStayRateCache
}

// HotelSearchCriteria top level element for criterion
type HotelSearchCriteria struct {
	XMLName       xml.Name `xml:"HotelSearchCriteria"`
	NumProperties int      `xml:"NumProperties,attr,omitempty"`
	Criterion     Criterion
}

func (h *HotelSearchCriteria) SetNumberOfHotels(n int) {
	h.NumProperties = n
}

// QuerySearchParams is a typed function to support optional query params on creation of new search criterion
type QuerySearchParams func(*HotelSearchCriteria) error

// HotelRefCriterion map of hotel ref criteria
type HotelRefCriterion map[string][]string

// PointOfInterestCriterion map of geographical area search
type PointOfInterestCriterion map[string]string

// AddressCriterion map of address search criteria; must be used with other criterion
// not recommended...
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

// PointOfInterest contains a number of specific search criteria for geographically named points.
// It supports City names, tourist attractions, or other well known landmarks.
type PointOfInterest struct {
	XMLName          xml.Name `xml:"PointOfInterest"`
	CountryStateCode string   `xml:"CountryStateCode,attr"`
	NonUS            bool     `xml:"NonUS,attr"`
	RPH              string   `xml:"RHP,attr,omitempty"`
	Val              string   `xml:",chardata"`
}

// PropertyType container for searhing types of properties (APTS, LUXRY...)
type PropertyType struct {
	Val string `xml:",chardata"`
}

// Package container for searching types of packages (GF, HM, BB...)
type Package struct {
	Val string `xml:",chardata"`
}

// Address represents typical building addresses; state province nil pointer ignored if empty.
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNumber,omitempty"`
	City          string `xml:"CityName,omitempty"`
	CountryCode   string `xml:"CountryCode,omitempty"`
	Postal        string `xml:"PostalCode,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr"`
	} `xml:"StateCountyProv"`
}

// AddressSearchStruct speical container for Criterion searches with address.
type AddressSearchStruct struct {
	XMLName     xml.Name `xml:"Address"`
	CityName    string   `xml:"CityName,omitempty"`
	CountryCode string   `xml:"CountryCode,omitempty"`
	PostalCode  string   `xml:"PostalCode,omitempty"`
	StreetNmbr  string   `xml:"StreetNmbr,omitempty"`
}

// Criterion holds various serach criteria
type Criterion struct {
	XMLName         xml.Name `xml:"Criterion"`
	HotelRefs       []*HotelRef
	PointOfInterest *PointOfInterest
	AddressSearch   *AddressSearchStruct
	PropertyTypes   []*PropertyType
	Packages        []*Package
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

// RatePlan container for attriubtes on rates
type RatePlan struct {
	CurrencyCode    string `xml:"CurrencyCode,attr,omitempty"`
	DCA_ProductCode string `xml:"DCA_ProductCode,attr,omitempty"`
	DecodeAll       string `xml:"DecodeAll,attr,omitempty"`
	RateCode        string `xml:"RateCode,attr,omitempty"`
	RPH             string `xml:"RPH,attr,omitempty"`
}

// ContractNegotiatedRateCode specify a contract or negotiated rate code (Client ID)
// Combined amount of "ContractNegotiatedRateCode" and .../RatePlanCandidates/RateAccessCode" elements cannot exceed 8
type ContractNegotiatedRateCode struct {
	Val string `xml:",chardata"`
}

// RateAccessCode used to specify ID associated with RAC Code.
// Equivalent Sabre host command: HOTFSG/21JUN-24JUN2/RC-X‡AAA,X‡AOM
type RateAccessCode struct {
	Code string `xml:"Code,attr"`
	Val  string `xml:",chardata"`
}

// RatePlanCode specify a Rate Plan Code (rate category)
// Acceptable values are: V, C, D, I, F, GOV, M, P, S, TVL, W, R, N, X, or ALL.
type RatePlanCode struct {
	Val string `xml:",chardata"`
}

// RateRange is used to specify minimum/maximum rates to be returned, as well as currency.
type RateRange struct {
	CurrencyCode string `xml:"CurrencyCode,attr"`
	Max          string `xml:",omitempty"`
	Min          string `xml:",omitempty"`
}

// RatePlanCandidates determines types of rates.
type RatePlanCandidates struct {
	XMLName                     xml.Name                      `xml:"RatePlanCandidates"`
	PromotionalSpot             string                        `xml:"PromotionalSpot,attr,omitempty"`
	RateAssured                 bool                          `xml:"RateAssured,attr,omitempty"`
	SuppressRackRate            bool                          `xml:"SuppressRackRate,attr,omitempty"`
	RatePlans                   []*RatePlan                   `xml:"RatePlanCandidate"`
	ContractNegotiatedRateCodes []*ContractNegotiatedRateCode `xml:"ContractNegotiatedRateCode"`
	RateAccessCodes             []*RateAccessCode             `xml:"RateAccessCode"`
	RatePlanCodes               []*RatePlanCode               `xml:"RatePlanCode"`
	RateRange                   *RateRange                    `xml:"RateRange"`
}

// SetRatePlans helper to create a slice of rate plans to append to an Avail Segement for search or description services.
func (rpc *RatePlanCandidates) SetRatePlans(ratePlans []RatePlan) {
	for _, plan := range ratePlans {
		rpc.RatePlans = append(rpc.RatePlans, &plan)
	}
}

// SetContractNegotiatedRates helper for setting the negotiated rates on availability requests.
func (a *AvailRequestSegment) SetContractNegotiatedRates(rateCodes []string) {
	rpc := &RatePlanCandidates{}
	for _, code := range rateCodes {
		rpc.ContractNegotiatedRateCodes = append(rpc.ContractNegotiatedRateCodes, &ContractNegotiatedRateCode{code})
	}
	a.RatePlanCandidates = rpc
}

// SetContractNegotiatedRates helper for setting the negotiated rates on availability requests.
// RC-G,S,C "\u2628" ALL
func (rpc *RatePlanCandidates) SetContractNegotiatedRates(rateCodes []string) {
	for _, code := range rateCodes {
		rpc.ContractNegotiatedRateCodes = append(rpc.ContractNegotiatedRateCodes, &ContractNegotiatedRateCode{code})
	}
}

// SetRatePlanCodes convencience for setting the negotiated rates on availability requests.
func (rpc *RatePlanCandidates) SetRatePlanCodes(rateCodes []string) {
	for _, code := range rateCodes {
		rpc.RatePlanCodes = append(rpc.RatePlanCodes, &RatePlanCode{code})
	}
}

type AdditionalAvail struct {
	XMLName xml.Name `xml:"AdditionalAvail"`
	Ind     bool     `xml:"Ind,attr"`
}

// AvailAvailRequestSegment holds basic hotel availability params: customer ids, guest count, hotel search criteria and arrival departure. nil pointers ignored if empty
type AvailRequestSegment struct {
	XMLName             xml.Name `xml:"AvailRequestSegment"`
	AdditionalAvail     *AdditionalAvail
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

type RoomToBook struct {
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	NumRooms  string `form:"num_rooms"`
	CCPhone   string `form:"cc_phone"`
	CCCode    string `form:"cc_code"`
	CCExpire  string `form:"cc_expire"`
	CCNumber  string `form:"cc_number"`
	RoomMeta  string `form:"room_meta"`

	B64RoomMetaData string            //passing room meta in form
	ParsedRoomMeta  ParsedRoomMeta    //parsing the room_meta from from
	FormErrors      map[string]string //collecting errors on validate
}

//Validate BookRoomParams checks form fields for submitting booking
//TODO add more validations...
func (b *RoomToBook) ValidateAndSetParsedRoomMeta() bool {
	b.FormErrors = make(map[string]string)
	meta, err := NewParsedRoomMeta(b.RoomMeta)
	if err != nil {
		b.FormErrors["RoomMeta"] = "Room Metadata is malformed"
	} else {
		b.ParsedRoomMeta = meta
	}
	if strings.TrimSpace(b.FirstName) == "" {
		b.FormErrors["FirstName"] = "First Name cannot be empty"
	}
	if strings.TrimSpace(b.LastName) == "" {
		b.FormErrors["LastName"] = "Last Name cannot be empty"
	}
	return len(b.FormErrors) == 0
}

type RoomRate struct {
	XMLName            xml.Name `xml:"RoomRate" json:"-"`
	ClientID           string   `xml:"ClientID,attr"`
	DirectConnect      string   `xml:"DirectConnect,attr"`
	GuaranteeSurcharge string   `xml:"GuaranteeSurchargeRequired,attr"`
	GuaranteedRate     string   `xml:"GuaranteedRateProgram,attr"`
	IATA_Character     string   `xml:"IATA_CharacteristicIdentification,attr"`
	IATA_Product       string   `xml:"IATA_ProductIdentification,attr"`
	LowInventory       string   `xml:"LowInventoryThreshold,attr"`
	RateAccessCode     string   `xml:"RateAccessCode,attr"`
	RateCategory       string   `xml:"RateCategory,attr"`
	RateLevelCode      string   `xml:"RateLevelCode,attr"`
	RPH                string   `xml:"RPH,attr"`
	RateChangeInd      string   `xml:"RateChangeInd,attr"`
	RateConversionInd  string   `xml:"RateConversionInd,attr"`
	RoomLocationCode   string   `xml:"RoomLocationCode,attr"`
	SpecialOffer       string   `xml:"SpecialOffer,attr"`
	Rates              []Rate   `xml:"Rates>Rate"`
	AdditionalInfo     AdditionalInfo
	HotelRateCode      string `xml:"HotelRateCode"`
	//B64RoomMetaData    string
	RoomToBook RoomToBook
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
	XMLName            xml.Name `xml:"BasicPropertyInfo" json:"-"`
	AreaID             string   `xml:"AreaID,attr"` //TODO: Area ID... this is wrong
	ChainCode          string   `xml:"ChainCode,attr"`
	Distance           string   `xml:"Distance,attr"`
	GEO_ConfAvail      string   `xml:"GEO_ConfidenceLevel,attr"` //hotel avail
	GeoConfPropDesc    string   `xml:"GeoConfidenceLevel,attr"`  //property description
	HotelCityCode      string   `xml:"HotelCityCode,attr"`
	HotelCode          string   `xml:"HotelCode,attr"`
	HotelName          string   `xml:"HotelName,attr"`
	ConfirmationNumber struct {
		Val string `xml:",chardata"`
	} `xml:"ConfirmationNumber"`
	CancelPenalty struct {
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
	XMLName     xml.Name `xml:"Charges"`
	Crib        string   `xml:"Crib,attr"`
	ExtraPerson string   `xml:"ExtraPerson,attr"`
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
