package hotel

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/ailgroup/sbrweb"
)

var (
	hqbad               = make(HotelRefCriterion)
	hqcity              = make(HotelRefCriterion)
	hqids               = make(HotelRefCriterion)
	hqltln              = make(HotelRefCriterion)
	addr                = make(AddressCriterion)
	sampleCID           = "12345"
	sampleLatLang       = []string{"32.78,-96.81", "54.87,-102.96"}
	sampleHotelCode     = []string{"0012", "19876", "1109", "445098", "000034"}
	sampleHotelCityCode = []string{"DFW", "CHC", "LA"}
	samplePackages      = []string{"GF", "HM", "BB"}
	samplePropertyTypes = []string{"APTS", "LUXRY"}
	sampleGuestCount    = 2
	sampleStreet        = "2031 N. 100 W"
	sampleCity          = "Nowhere"
	samplePostal        = "999908"
	sampleCountryCode   = "US"
	sampleArrive        = "04-02"
	sampleDepart        = "04-05"

	samplesite        = "www.z.com"
	samplepcc         = "7TZA"
	samplebinsectoken = string([]byte(`Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0`))
	sampleconvid      = "fds8789h|dev@z.com"
	samplemid         = "mid:20180207-20:19:07.25|QVbg0"
	sampletime        = "2018-02-16T07:18:42Z"

	sampleAvailRQHotelIDSCoprID = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><Customer><Corporate><ID>12345</ID></Corporate></Customer><GuestCounts Count="4"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="0012"></HotelRef><HotelRef HotelCode="19876"></HotelRef><HotelRef HotelCode="1109"></HotelRef><HotelRef HotelCode="445098"></HotelRef><HotelRef HotelCode="000034"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQCitiesCustNumber = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><Customer><ID><Number>12345</Number></ID></Customer><GuestCounts Count="3"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCityCode="DFW"></HotelRef><HotelRef HotelCityCode="CHC"></HotelRef><HotelRef HotelCityCode="LA"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQLatLng = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef Latitude="32.78" Longitude="-96.81"></HotelRef><HotelRef Latitude="54.87" Longitude="-102.96"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQPropPackages = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><PropertyTypes>APTS</PropertyTypes><PropertyTypes>LUXRY</PropertyTypes><Packages>GF</Packages><Packages>HM</Packages><Packages>BB</Packages></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQHotelIDS = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">OTA_HotelAvailRQ</eb:Service><eb:Action>OTA_HotelAvailLLSRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180207-20:19:07.25|QVbg0</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="0012"></HotelRef><HotelRef HotelCode="19876"></HotelRef><HotelRef HotelCode="1109"></HotelRef><HotelRef HotelCode="445098"></HotelRef><HotelRef HotelCode="000034"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ></soap-env:Body></soap-env:Envelope>`)

	sampleHotelAvailRSgood = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">OTA_HotelAvailRQ</eb:Service><eb:Action>OTA_HotelAvailLLSRS</eb:Action><eb:MessageData><eb:MessageId>1374478129129220211</eb:MessageId><eb:Timestamp>2018-04-03T03:35:13</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESF!ICESMSLB\/RES.LB!-3161638152750045809!1191725!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><OTA_HotelAvailRS xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:stl="http://services.sabre.com/STL/v01" Version="2.3.0">
		 <stl:ApplicationResults status="Complete">
		  <stl:Success timeStamp="2018-04-02T22:35:13-05:00">
		   <stl:SystemSpecificResults>
			<stl:HostCommand LNIATA="222222">HOTPRPL-1/02APR-05APR2</stl:HostCommand>
		   </stl:SystemSpecificResults>
		  </stl:Success>
		 </stl:ApplicationResults>
		 <AdditionalAvail Ind="false"/>
		 <AvailabilityOptions>
		  <AvailabilityOption RPH="001">
		   <BasicPropertyInfo AreaID="000" ChainCode="HY" Distance="M" GEO_ConfidenceLevel="1" HotelCityCode="TUL" HotelCode="0000001" HotelName="HYATT REGENCY TULSA" Latitude="36.154800" Longitude="-95.990356">
			<Address>
			 <AddressLine>100 E SECOND STREET</AddressLine>
			 <AddressLine>TULSA OK 74103</AddressLine>
			</Address>
			<ContactNumbers>
			 <ContactNumber Fax="1-918-560-2232" Phone="1-918-234-1234"/>
			</ContactNumbers>
			<DirectConnect>
			 <Alt_Avail Ind="false"/>
			 <DC_AvailParticipant Ind="true"/>
			 <DC_SellParticipant Ind="true"/>
			 <RatesExceedMax Ind="false"/>
			 <UnAvail Ind="false"/>
			</DirectConnect>
			<LocationDescription>
			 <Text>TULSA OK</Text>
			</LocationDescription>
			<Property Rating="NTM">
			 <Text>4 CROWN</Text>
			</Property>
			<PropertyOptionInfo>
			 <ADA_Accessible Ind="true"/>
			 <AdultsOnly Ind="false"/>
			 <BeachFront Ind="false"/>
			 <Breakfast Ind="false"/>
			 <BusinessCenter Ind="true"/>
			 <BusinessReady Ind="false"/>
			 <Conventions Ind="true"/>
			 <Dataport Ind="true"/>
			 <Dining Ind="true"/>
			 <DryClean Ind="false"/>
			 <EcoCertified Ind="false"/>
			 <ExecutiveFloors Ind="true"/>
			 <FitnessCenter Ind="true"/>
			 <FreeLocalCalls Ind="false"/>
			 <FreeParking Ind="false"/>
			 <FreeShuttle Ind="true"/>
			 <FreeWifiInMeetingRooms Ind="false"/>
			 <FreeWifiInPublicSpaces Ind="false"/>
			 <FreeWifiInRooms Ind="true"/>
			 <FullServiceSpa Ind="true"/>
			 <GameFacilities Ind="false"/>
			 <Golf Ind="false"/>
			 <HighSpeedInternet Ind="true"/>
			 <HypoallergenicRooms Ind="true"/>
			 <IndoorPool Ind="true"/>
			 <InRoomCoffeeTea Ind="true"/>
			 <InRoomMiniBar Ind="false"/>
			 <InRoomRefrigerator Ind="false"/>
			 <InRoomSafe Ind="true"/>
			 <InteriorDoorways Ind="false"/>
			 <Jacuzzi Ind="false"/>
			 <KidsFacilities Ind="false"/>
			 <KitchenFacilities Ind="false"/>
			 <MealService Ind="false"/>
			 <MeetingFacilities Ind="true"/>
			 <NoAdultTV Ind="false"/>
			 <NonSmoking Ind="true"/>
			 <OutdoorPool Ind="true"/>
			 <Pets Ind="true"/>
			 <Pool Ind="true"/>
			 <PublicTransportationAdjacent Ind="false"/>
			 <RateAssured Ind="true"/>
			 <Recreation Ind="false"/>
			 <RestrictedRoomAccess Ind="true"/>
			 <RoomService Ind="true"/>
			 <RoomService24Hours Ind="false"/>
			 <RoomsWithBalcony Ind="false"/>
			 <SkiInOutProperty Ind="false"/>
			 <SmokeFree Ind="true"/>
			 <SmokingRoomsAvail Ind="false"/>
			 <Tennis Ind="false"/>
			 <WaterPurificationSystem Ind="false"/>
			 <Wheelchair Ind="true"/>
			</PropertyOptionInfo>
			<RateRange CurrencyCode="USD" Max="289.00" Min="134.00"/>
			<RoomRate RateLevelCode="RAC">
			 <AdditionalInfo>
			  <CancelPolicy Numeric="00"/>
			 </AdditionalInfo>
			 <HotelRateCode>RAC</HotelRateCode>
			</RoomRate>
			<SpecialOffers Ind="false"/>
		   </BasicPropertyInfo>
		  </AvailabilityOption>
		 </AvailabilityOptions>
		</OTA_HotelAvailRS></soap-env:Body></soap-env:Envelope>`)
)

func init() {
	hqcity[cityQueryField] = sampleHotelCityCode
	hqids[hotelidQueryField] = sampleHotelCode
	hqltln[latlngQueryField] = sampleLatLang

	addr[streetQueryField] = sampleStreet
	addr[cityQueryField] = sampleCity
	addr[postalQueryField] = samplePostal
	addr[countryCodeQueryField] = sampleCountryCode
}

func TestAddressSearchReturnError(t *testing.T) {
	_, err := NewHotelSearchCriteria(
		AddressSearch(AddressCriterion{}),
	)
	if err == nil {
		t.Errorf("AddressSearch empty params should return error: '%v'", err)
	}
}

func TestAddressSearchCriteria(t *testing.T) {
	a, err := NewHotelSearchCriteria(
		AddressSearch(addr),
	)

	if err != nil {
		t.Errorf("NewHotelSearchCriteria with AddressSearch error %v", err)
	}
	if a.Criterion.Address.Street != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if a.Criterion.Address.City != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if a.Criterion.Address.Postal != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if a.Criterion.Address.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
	}
}

func TestPackageSearchCriteria(t *testing.T) {
	p, err := NewHotelSearchCriteria(
		PackageSearch(samplePackages),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with PackageSearch error %v", err)
	}
	if len(p.Criterion.Packages) != len(samplePackages) {
		t.Errorf("PackageSearch wrong number of results, expected %d got %d", len(samplePackages), len(p.Criterion.Packages))
	}
}

func TestPropertyTypeSearchCriteria(t *testing.T) {
	p, err := NewHotelSearchCriteria(
		PropertyTypeSearch(samplePropertyTypes),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with PropertyTypeSearch error %v", err)
	}
	if len(p.Criterion.PropertyTypes) != len(samplePropertyTypes) {
		t.Errorf("PackageSearch wrong number of results, expected %d got %d", len(samplePropertyTypes), len(p.Criterion.PropertyTypes))
	}
}

func TestHotelRefSearchReturnError(t *testing.T) {
	_, err := NewHotelSearchCriteria(
		HotelRefSearch(hqbad),
	)
	if err == nil {
		t.Errorf("HotelRefSearch empty params should return error: '%v'", err)
	}
}

func TestHotelRefSearchCityCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleHotelCityCode {
		if r.Criterion.HotelRefs[i].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", i, code, r.Criterion.HotelRefs[i].HotelCityCode)
		}

	}
}

func TestHotelRefSearchHotelCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleHotelCode {
		if r.Criterion.HotelRefs[i].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", i, code, r.Criterion.HotelRefs[i].HotelCode)
		}

	}
}

func TestHotelRefSearchLatLngCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleLatLang {
		ll := strings.Split(code, ",")
		if r.Criterion.HotelRefs[i].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", i, ll[0], r.Criterion.HotelRefs[i].Latitude)
		}
		if r.Criterion.HotelRefs[i].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", i, ll[1], r.Criterion.HotelRefs[i].Longitude)
		}
	}
}

func TestMultipleHotelCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
		HotelRefSearch(hqcity),
		AddressSearch(addr),
		HotelRefSearch(hqltln),
		PackageSearch(samplePackages),
		PropertyTypeSearch(samplePropertyTypes),
	)

	if err != nil {
		t.Errorf("NewHotelSearchCriteria with basic criteria error %v", err)
	}

	counter := 0
	for _, code := range sampleHotelCode {
		if r.Criterion.HotelRefs[counter].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", counter, code, r.Criterion.HotelRefs[counter].HotelCode)
		}
		counter++
	}
	for _, code := range sampleHotelCityCode {
		if r.Criterion.HotelRefs[counter].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", counter, code, r.Criterion.HotelRefs[counter].HotelCityCode)
		}
		counter++
	}

	if r.Criterion.Address.Street != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if r.Criterion.Address.City != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if r.Criterion.Address.Postal != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if r.Criterion.Address.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
	}

	for _, code := range sampleLatLang {
		ll := strings.Split(code, ",")
		if r.Criterion.HotelRefs[counter].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", counter, ll[0], r.Criterion.HotelRefs[counter].Latitude)
		}
		if r.Criterion.HotelRefs[counter].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", counter, ll[1], r.Criterion.HotelRefs[counter].Longitude)
		}
		counter++
	}
}

func TestSetHotelAvailRqStructMarshal(t *testing.T) {
	availBody := SetHotelAvailRqStruct(sampleGuestCount, HotelSearchCriteria{}, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)

	if avail.XMLNSXsi != sbrweb.BaseXSINamespace {
		t.Errorf("SetHotelAvailRqStruct XMLNSXsi expect: %s, got %s", sbrweb.BaseXSINamespace, avail.XMLNSXsi)
	}
	if avail.Version != hotelRQVersion {
		t.Errorf("SetHotelAvailRqStruct Version expect: %s, got %s", hotelRQVersion, avail.Version)
	}
	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("SetHotelAvailRqStruct GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("SetHotelAvailRqStruct Customer.Corporate.ID expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestSetHotelAvailRqStructCorpID(t *testing.T) {
	availBody := SetHotelAvailRqStruct(sampleGuestCount, HotelSearchCriteria{}, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)

	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("SetHotelAvailRqStruct Customer.Corporate.ID  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

	avail.addCustomerID(sampleCID)
	if avail.Avail.Customer.CustomerID.Number != sampleCID {
		t.Errorf("SetHotelAvailRqStruct CustomerID.Number  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

}

func TestAvailIdsMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	gcount := 4
	availBody := SetHotelAvailRqStruct(gcount, q, sampleArrive, sampleDepart)

	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("SetHotelAvailRqStruct GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqids[hotelidQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqids[hotelidQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQHotelIDSCoprID) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDSCoprID), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestAvailCitiesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	gcount := 3
	availBody := SetHotelAvailRqStruct(gcount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCustomerID(sampleCID)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqcity[cityQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqcity[cityQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQCitiesCustNumber) {
		t.Errorf("Expected marshal hotel avail for hotel cities \n sample: %s \n result: %s", string(sampleAvailRQCitiesCustNumber), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestAvailLatLngMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	availBody := SetHotelAvailRqStruct(sampleGuestCount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ

	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqltln[latlngQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqltln[latlngQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQLatLng) {
		t.Errorf("Expected marshal set hotel avail for hotel lat/lng \n sample: %s \n result: %s", string(sampleAvailRQLatLng), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestAvailPropertyTypesPackagesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		PackageSearch(samplePackages),
		PropertyTypeSearch(samplePropertyTypes),
	)
	availBody := SetHotelAvailRqStruct(sampleGuestCount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQPropPackages) {
		t.Errorf("Expected marshal set hotel avail for hotel packages and property types \n sample: %s \n result: %s", string(sampleAvailRQLatLng), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestBuildHotelAvailRequestMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	avail := SetHotelAvailRqStruct(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelAvailRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, avail)

	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}

	if string(b) != string(sampleAvailRQHotelIDS) {
		t.Errorf("Expected marshal SOAP hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDS), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestHotelAvailUnmarshal(t *testing.T) {
	avail := HotelAvailResponse{}
	err := xml.Unmarshal(sampleHotelAvailRSgood, &avail)
	if err != nil {
		t.Errorf("Error unmarshaling hotel avail %s \nERROR: %v", sampleHotelAvailRSgood, err)
	}
	reqError := avail.Body.HotelAvail.Result.Error
	if reqError.Type != "" {
		t.Errorf("Request error %v should not have message %s", reqError, reqError.System.Message)
	}
	success := avail.Body.HotelAvail.Result.Success
	if success.System.HostCommand.LNIATA != "222222" {
		t.Errorf("System.HostCommand.LNIATA for success expect: %v, got: %v", "222222", success.System.HostCommand.LNIATA)
	}

	//fmt.Printf("SAMPLE: %s\n", sampleEnvelope)
	//fmt.Printf("CURRENT: %+v\n", success)
}
