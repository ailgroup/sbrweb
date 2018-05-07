package hotel

import "encoding/xml"

type BasicPropertyInfo struct {
	XMLName             xml.Name `xml:"BasicPropertyInfo"`
	AreadID             string   `xml:"AreadID,attr"`
	ChainCode           string   `xml:"ChainCode,attr"`
	Distance            string   `xml:"Distance,attr"`
	GEO_ConfidenceLevel string   `xml:"GEO_ConfidenceLevel,attr"`
	HotelCityCode       string   `xml:"HotelCityCode,attr"`
	HotelCode           string   `xml:"HotelCode,attr"`
	HotelName           string   `xml:"HotelName,attr"`
	Latitude            string   `xml:"Latitude,attr"`
	Longitude           string   `xml:"Longitude,attr"`
	Address             struct {
		Line []string `xml:"AddressLine"`
	} `xml:"Address"`
	ContactNumbers struct {
		Number ContactNumber `xml:"ContactNumber"`
	} `xml:"ContactNumbers"`
	DirectCon DirectConnect `xml:"DirectConnect"`
	LocDesc   struct {
		Text string `xml:"Text"`
	} `xml:"LocationDescription"`
	Prop struct {
		Rating string `xml:"Rating,attr"`
		Text   string `xml:"Text"`
	} `xml:"Property"`
	PropOptInfo PropertyOptionInfo `xml:"PropertyOptionInfo"`
	RateRange   struct {
		CurrencyCode string `xml:"CurrencyCode,attr"`
		Max          string `xml:"Max,attr"`
		Min          string `xml:"Min,attr"`
	} `xml:"RateRange"`
	RoomRate      RoomRate
	SpecialOffers struct {
		Ind bool `xml:"Ind,attr"`
	} `xml:"SpecialOffers"`
}

type RoomRate struct {
	XMLName        xml.Name `xml:"RoomRate"`
	RateLevelCode  string   `xml:"RateLevelCode,attr"`
	AdditionalInfo struct {
		CancelPolicy struct {
			Numeric int `xml:"Numeric,attr"` //string? 001 versus 1
		} `xml:"CancelPolicy"`
	} `xml:"AdditionalInfo"`
	HotelRateCode string `xml:"HotelRateCode"`
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
