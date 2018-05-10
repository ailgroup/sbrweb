package hotelws

import (
	"encoding/xml"
)

type BasicPropertyInfo struct {
	XMLName         xml.Name `xml:"BasicPropertyInfo"`
	AreadID         string   `xml:"AreadID,attr"`
	ChainCode       string   `xml:"ChainCode,attr"`
	Distance        string   `xml:"Distance,attr"`
	GEO_ConfAvail   string   `xml:"GEO_ConfidenceLevel,attr"` //hotel avail
	GeoConfPropDesc string   `xml:"GeoConfidenceLevel,attr"`  //property description
	HotelCityCode   string   `xml:"HotelCityCode,attr"`
	HotelCode       string   `xml:"HotelCode,attr"`
	HotelName       string   `xml:"HotelName,attr"`
	Latitude        string   `xml:"Latitude,attr"`
	Longitude       string   `xml:"Longitude,attr"`
	Address         struct {
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
	XMLName        xml.Name `xml:"AdditionalGuestAmount"`
	MaxExtraPerson int      `xml:"MaxExtraPersonsAllowed,attr"`
	NumCribs       int      `xml:"NumCribs,attr"`
	Charges        []Charge `xml:"Charges"`
}

type TotalSurcharges struct {
	XMLName xml.Name `xml:"TotalSurcharges"`
	Amount  string   `xml:"Amount,attr"`
}
type TotalTaxes struct {
	XMLName xml.Name `xml:"TotalTaxes"`
	Amount  string   `xml:"Amount,attr"`
}

type HotelPricing struct {
	XMLName         xml.Name `xml:"HotelTotalPricing"`
	Amount          string   `xml:"Amount,attr"`
	Disclaimer      string   `xml:"Disclaimer"`
	TotalSurcharges TotalSurcharges
	TotalTaxes      TotalTaxes
}

type Rate struct {
	XMLName                xml.Name                `xml:"Rate"`
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
	AdditionalInfo     struct {
		CancelPolicy struct {
			Numeric int    `xml:"Numeric,attr"` //string? 001 versus 1
			Option  string `xml:"Option,attr"`
		} `xml:"CancelPolicy"`
		Text []string `xml:"Text"`
	} `xml:"AdditionalInfo"`
	HotelRateCode string `xml:"HotelRateCode"`
}

type VendorMessages struct {
	XMLName              xml.Name             `xml:"VendorMessages"`
	Attractions          Attractions          `xml:"Attractions"`
	Awards               Awards               `xml:"Awards"`
	Cancellation         Cancellation         `xml:"Cancellation"`
	Deposit              Deposit              `xml:"Deposit"`
	Description          Description          `xnl:"Description"`
	Dining               Dining               `xml:"Dining"`
	Directions           Directions           `xml:"Directions"`
	Facilities           Facilities           `xml:"Facilities"`
	Guarantee            Guarantee            `xml:"Guarantee"`
	Location             Location             `xml:"Location"`
	MarketingInformation MarketingInformation `xml:"MarketingInformation"`
	MiscServices         MiscServices         `xml:"MiscServices"`
	Policies             Policies             `xml:"Policies"`
	Recreation           Recreation           `xml:"Recreation"`
	Rooms                Rooms                `xml:"Rooms"`
	Safety               Safety               `xml:"Safety"`
	Services             Services             `xml:"Services"`
	Transportation       Transportation       `xml:"Transportation"`
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
type Guarantee struct {
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

type IndexD struct {
	XMLName            xml.Name `xml:"Index"`
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
	XMLName xml.Name   `xml:"ApplicationResults"`
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
