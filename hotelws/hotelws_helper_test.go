package hotelws

import (
	"net/http"
	"net/http/httptest"
)

/*
TESTING NOTES:
	- all data variables use for mocking tests are downcase and start with sample*
	- functions used in a test have their tests come first
	- Benchmarks come after tests using testable functionality
	- Benchmarks have same name as test that is benchmarked (sans Test/Benchmark prefix)
*/

var (
	//serverHotelDown server mocks unavailable service
	serverHotelDown = &httptest.Server{}
	//serverBadBody mocks a server that returns malformed body
	serverBadBody = &httptest.Server{}
	//serverHotelAvailability server to retrieve hotel availability using OTA_HotelAvailLLSRQ.
	serverHotelAvailability = &httptest.Server{}
	//serverHotelPropertyDesc server to retrieve hotel rates using HotelPropertyDescriptionLLSRQ.
	serverHotelPropertyDesc = &httptest.Server{}
	//serverHotelRateDesc server to retrieve rules and policies using HotelRateDescriptionLLSRQ.
	serverHotelRateDesc = &httptest.Server{}
)

//Initialize Mock Sabre Web Servers and test data
func init() {
	// init data chunks...
	hqcity[cityQueryField] = sampleHotelCityCode
	hqids[hotelidQueryField] = sampleHotelCode
	hqltln[latlngQueryField] = sampleLatLang

	addr[streetQueryField] = sampleStreet
	addr[cityQueryField] = sampleCity
	addr[postalQueryField] = samplePostal
	addr[countryCodeQueryField] = sampleCountryCode

	// init test servers...
	serverHotelDown = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleHotelAvailRSgood)
			},
		),
	)
	defer func() { serverHotelDown.Close() }()

	serverBadBody = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				//rs.Header()
				//rs.WriteHeader(500)
				//rs.Write(sampleBadBody)
				rs.Write([]byte(`<!# SOME BAD--XML_/__.*__\\fhji(*&^%^%<Boo<HA/>/>$%^&Y*(J)OPKL:/>`))
			},
		),
	)

	serverHotelAvailability = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleHotelAvailRSgood)
			},
		),
	)
	//defer func() { serverHotelAvailability.Close() }()

	serverHotelPropertyDesc = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleHotelPropDescRSgood)
			},
		),
	)
	//defer func() { serverHotelPropertyDesc.Close() }()

	serverHotelRateDesc = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleHotelRateDescRSgood)
			},
		),
	)
	//defer func() { serverHotelRateDesc.Close() }()
}

// data chunks for testing
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
	samplesite          = "www.z.com"
	samplepcc           = "7TZA"
	samplebinsectoken   = string([]byte(`Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0`))
	sampleconvid        = "fds8789h|dev@z.com"
	samplemid           = "mid:20180207-20:19:07.25|QVbg0"
	sampletime          = "2018-02-16T07:18:42Z"
)

var iataCharSample = []string{"P1KRAC", "D1KRAC", "L1KRAC", "P1KBRF", "T1KRAC", "P2TRAC", "K1KRAC", "L2TRAC", "E2DRAC", "E1KRAC", "D2DRAC", "C1KRAC", "U1QRAC", "A2TRAC", "N1KRAC", "N1QRAC"}

// table test for room rates on property description
var rateSamples = []struct {
	direct       string
	surcharge    string
	guarrate     string
	iatachar     string
	iataprod     string
	lowinventory string
	ratecode     string
	rph          string
	ratechange   string
	rateconv     string
	specialoff   string
	rates        []Rate
}{
	{
		direct:       "",
		surcharge:    "G",
		guarrate:     "false",
		iatachar:     "P1KRAC",
		iataprod:     "FULLY FLEXIBLE-",
		lowinventory: "false",
		ratecode:     "",
		rph:          "01",
		ratechange:   "false",
		rateconv:     "false",
		specialoff:   "false",
		rates: []Rate{
			Rate{
				Amount:       "285.00",
				CurrencyCode: "SGD",
				AdditionalGuestAmounts: []AdditionalGuestAmount{
					AdditionalGuestAmount{
						Charges: []Charge{
							Charge{
								ExtraPerson: "80.00",
							},
						},
					},
				},
				HotelPricing: HotelPricing{
					Amount: "335.45",
					TotalSurcharges: TotalSurcharges{
						Amount: "28.50",
					},
					TotalTaxes: TotalTaxes{
						Amount: "21.95",
					},
				},
			},
		},
	},
}

// XML blobs for testing
var (
	samplePropRQIDs = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">HotelPropertyDescription</eb:Service><eb:Action>HotelPropertyDescriptionLLSRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180207-20:19:07.25|QVbg0</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><HotelPropertyDescriptionRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="10"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></HotelPropertyDescriptionRQ></soap-env:Body></soap-env:Envelope>`)

	sampleAvailRQHotelIDSCoprIDRatePlans = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><Customer><Corporate><ID>12345</ID></Corporate></Customer><GuestCounts Count="4"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="0012"></HotelRef><HotelRef HotelCode="19876"></HotelRef><HotelRef HotelCode="1109"></HotelRef><HotelRef HotelCode="445098"></HotelRef><HotelRef HotelCode="000034"></HotelRef></Criterion></HotelSearchCriteria><RatePlanCandidates><RatePlanCandidate CurrencyCode="USD" DCA_ProductCode="I7A"></RatePlanCandidate></RatePlanCandidates><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQCitiesCustNumber = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><Customer><ID><Number>12345</Number></ID></Customer><GuestCounts Count="3"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCityCode="DFW"></HotelRef><HotelRef HotelCityCode="CHC"></HotelRef><HotelRef HotelCityCode="LA"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQLatLng = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef Latitude="32.78" Longitude="-96.81"></HotelRef><HotelRef Latitude="54.87" Longitude="-102.96"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQPropPackages = []byte(`<OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><PropertyTypes>APTS</PropertyTypes><PropertyTypes>LUXRY</PropertyTypes><Packages>GF</Packages><Packages>HM</Packages><Packages>BB</Packages></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQHotelIDS = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">OTA_HotelAvailRQ</eb:Service><eb:Action>OTA_HotelAvailLLSRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180207-20:19:07.25|QVbg0</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><OTA_HotelAvailRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="0012"></HotelRef><HotelRef HotelCode="19876"></HotelRef><HotelRef HotelCode="1109"></HotelRef><HotelRef HotelCode="445098"></HotelRef><HotelRef HotelCode="000034"></HotelRef></Criterion></HotelSearchCriteria><TimeSpan End="04-05" Start="04-02"></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ></soap-env:Body></soap-env:Envelope>`)

	sampleHotelRateDescRQRPH = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">HotelRateDescriptionLLSRQ</eb:Service><eb:Action>HotelRateDescriptionLLSRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180207-20:19:07.25|QVbg0</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><HotelRateDescriptionRQ Version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" ReturnHostCommand="true"><AvailRequestSegment><RatePlanCandidates><RatePlanCandidate RPH="12"></RatePlanCandidate></RatePlanCandidates></AvailRequestSegment></HotelRateDescriptionRQ></soap-env:Body></soap-env:Envelope>`)

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

	sampleHotelPropDescRSgood = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">HotelPropertyDescription</eb:Service><eb:Action>HotelPropertyDescriptionLLSRS</eb:Action><eb:MessageData><eb:MessageId>775733075202330295</eb:MessageId><eb:Timestamp>2018-05-08T02:05:24</eb:Timestamp><eb:RefToMessageId>mid:20180207-20:19:07.25|QVbg0</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><HotelPropertyDescriptionRS xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:stl="http://services.sabre.com/STL/v01" Version="2.3.0">
			 <stl:ApplicationResults status="Complete">
			  <stl:Success timeStamp="2018-02-16T07:18:42Z">
			   <stl:SystemSpecificResults>
				<stl:HostCommand LNIATA="222222">HOD4/11MAY-12MAY2</stl:HostCommand>
			   </stl:SystemSpecificResults>
			  </stl:Success>
			 </stl:ApplicationResults>
			 <RoomStay>
			  <BasicPropertyInfo ChainCode="SL" GeoConfidenceLevel="0" HotelCityCode="SIN" HotelCode="0000004" HotelName="SWISSOTEL THE STAMFORD" Latitude="1.292960" Longitude="103.853058" NumFloors="72" RPH="001">
			   <Address>
				<AddressLine>2 STAMFORD ROAD</AddressLine>
				<AddressLine>SINGAPORE SG 178882</AddressLine>
				<CountryCode>SG</CountryCode>
			   </Address>
			   <Awards>
				<AwardProvider>NTM4  CROWN</AwardProvider>
			   </Awards>
			   <CheckInTime>15:00</CheckInTime>
			   <CheckOutTime>12:00</CheckOutTime>
			   <ContactNumbers>
				<ContactNumber Fax="65-6338-2862" Phone="65-6338-8585"/>
			   </ContactNumbers>
			   <DirectConnect Ind="true">
				<AdditionalData Ind="true"/>
				<CurrencyConverted Ind="false"/>
				<DC_AvailParticipant Ind="true"/>
				<DC_SellParticipant Ind="true"/>
				<RequestFail Ind="false"/>
				<UnAvail Ind="false"/>
			   </DirectConnect>
			   <IndexData>
				<Index CountryState="C/SG" DistanceDirection="0N" LocationCode="A" Point="PARLIAMENT HOUSE" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="2W" LocationCode="A" Point="SHENTON WAY" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="4S" LocationCode="A" Point="SENTOSA ISLAND" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="0N" LocationCode="A" Point="WAR MEMORIA" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="1N" LocationCode="A" Point="MARINA SQUARE" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="0N" LocationCode="A" Point="MERLION PARK" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="2N" LocationCode="A" Point="NATL MUSEUM ART GALLERY" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="2NE" LocationCode="A" Point="ORCHARD ROAD" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="0N" LocationCode="A" Point="SINGAPORE" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="0N" LocationCode="A" Point="SINGAPORE CITY HALL" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="10N" LocationCode="A" Point="SINGAPORE ZOO" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="1E" LocationCode="A" Point="SUNTEC CITY" TransportationCode="O"/>
				<Index CountryState="C/SG" DistanceDirection="2W" LocationCode="A" Point="BOAT QUAY" TransportationCode="O"/>
				<Index DistanceDirection="14E" LocationCode="C" Point="SIN" TransportationCode="O"/>
			   </IndexData>
			   <PropertyOptionInfo>
				<ADA_Accessible Ind="false"/>
				<AdultsOnly Ind="false"/>
				<AirportShuttle Ind="true"/>
				<BeachFront Ind="false"/>
				<Breakfast Ind="false"/>
				<BusinessCenter Ind="true"/>
				<BusinessReady Ind="false"/>
				<CarRentalCounter>CAN BE BOOKED V</CarRentalCounter>
				<Conventions Ind="true"/>
				<Dataport Ind="true"/>
				<Dining Ind="true"/>
				<DryClean Ind="true"/>
				<EcoCertified Ind="false"/>
				<ExecutiveFloors Ind="true"/>
				<FamilyPlan Ind="false"/>
				<FitnessCenter Ind="true"/>
				<FreeLocalCalls Ind="false"/>
				<FreeParking Ind="false"/>
				<FreeShuttle Ind="false"/>
				<FreeWifiInMeetingRooms Ind="false"/>
				<FreeWifiInPublicSpaces Ind="false"/>
				<FreeWifiInRooms Ind="false"/>
				<FullServiceSpa Ind="true"/>
				<GameFacilities Ind="false"/>
				<Golf Ind="false"/>
				<GovtSafetyFire Ind="true"/>
				<HighSpeedInternet Ind="true"/>
				<HypoallergenicRooms Ind="false"/>
				<IndoorPool Ind="false"/>
				<IndPetRestriction Ind="false"/>
				<InRoomCoffeeTea Ind="true"/>
				<InRoomMiniBar Ind="true"/>
				<InRoomRefrigerator Ind="true"/>
				<InRoomSafe Ind="true"/>
				<InteriorDoorways Ind="true"/>
				<Jacuzzi Ind="false"/>
				<KidsFacilities Ind="false"/>
				<KitchenFacilities Ind="false"/>
				<MealService Ind="false"/>
				<MeetingFacilities Ind="true"/>
				<NoAdultTV Ind="true"/>
				<NonSmoking Ind="true"/>
				<OutdoorPool Ind="true"/>
				<Parking>N</Parking>
				<Pets Ind="false"/>
				<Pool Ind="true"/>
				<PublicTransportationAdjacent Ind="true"/>
				<Recreation Ind="true"/>
				<RestrictedRoomAccess Ind="true"/>
				<RoomService Ind="true"/>
				<RoomService24Hours Ind="true"/>
				<RoomsWithBalcony Ind="true"/>
				<SkiInOutProperty Ind="false"/>
				<SmokeFree Ind="true"/>
				<SmokingRoomsAvail Ind="false"/>
				<Tennis Ind="true"/>
				<WaterPurificationSystem Ind="false"/>
				<Wheelchair Ind="true"/>
			   </PropertyOptionInfo>
			   <PropertyTypeInfo>
				<AllInclusive Ind="false"/>
				<Apartments Ind="false"/>
				<BedBreakfast Ind="false"/>
				<Castle Ind="false"/>
				<Conventions Ind="false"/>
				<Economy Ind="false"/>
				<ExtendedStay Ind="false"/>
				<Farm Ind="false"/>
				<First Ind="true"/>
				<Luxury Ind="false"/>
				<Moderate Ind="false"/>
				<Motel Ind="false"/>
				<Resort Ind="false"/>
				<Suites Ind="false"/>
			   </PropertyTypeInfo>
			   <SpecialOffers Ind="false"/>
			   <Taxes>
				<Text>17.7PCT TT</Text>
				<Text>.</Text>
			   </Taxes>
			   <VendorMessages>
				<Attractions>
				 <Text>CIVILIAN WAR MEMORIAL 0.1 MI</Text>
				 <Text>ESPLANADE THEATRES ON THE BAY 0.2 MI</Text>
				 <Text>FORMULA ONE AT MARINA BAY GARDENS BY THE BAY 1.2 MI</Text>
				 <Text>JURONG BIRD PARK 12.4 MI</Text>
				 <Text>LUNAR AND CHINESE NEW YEAR 0.8 MI</Text>
				 <Text>MARINA BAY SANDS 0.7 MI</Text>
				 <Text>NIGHT SAFARI 12.4 MI</Text>
				 <Text>RAFFLES PLACE CENTRAL BUSINESS DISTRICT RESORT WORLD SENTOSA</Text>
				 <Text>SENTOSA ISLAND 6.2 MI</Text>
				 <Text>SINGAPORE FLYER 0.5 MI</Text>
				 <Text>SINGAPORE ZOO 12.4 MI</Text>
				 <Text>THE MERLION 0.3 MI</Text>
				 <Text>CITY HALL MRT SPA-WILLOW STREAM SPA</Text>
				</Attractions>
				<Awards>
				 <Text>AAADIAMOND 5</Text>
				 <Text>AASTARRATING 5</Text>
				 <Text>MOBIL 5</Text>
				 <Text>OHG MODERATE DELUXE</Text>
				 <Text>STAR 5</Text>
				</Awards>
				<Cancellation>
				 <Text>CXL 48HRS PRIOR TO ARRIVAL TO AVOID PENALTY.  CANCELLATION</Text>
				 <Text>POLICY MAY VARY DEPENDING ON THE OFFER - CHECK RATE DISPLAY FOR</Text>
				 <Text>POLICIES **EXCEPTION -  SEPTEMBER 13TH TO SEPTEMBER 16TH 2018**</Text>
				 <Text>- FULL DEPOSIT TAKEN AT TIME OF BOOKING - NON-REFUNDABLE</Text>
				</Cancellation>
				<Deposit>
				 <Text>DEPOSIT POLICY MAY VARY DEPENDING ON THE OFFER - CHECK RATE</Text>
				 <Text>DISPLAY FOR POLICIES ****EXCEPTION -  SEPTEMBER 13TH TO</Text>
				 <Text>SEPTEMBER 16TH 2018** - FULL DEPOSIT TAKEN AT TIME OF BOOKING -</Text>
				 <Text>NON-REFUNDABLE</Text>
				</Deposit>
				<Description>
				 <Text>STEP INTO A WORLD OF EASE AND COMFORT AND EXPERIENCE THE FINEST</Text>
				 <Text>IN SWISS HOSPITALITY AT SWISSOTEL THE STAMFORD, THE TALLEST</Text>
				 <Text>HOTEL IN SOUTHEAST ASIA. STRATEGICALLY LOCATED IN THE HEART OF</Text>
				 <Text>SINGAPORES SHOPPING, DINING AND ENTERTAINMENT DISTRICTS WITH</Text>
				 <Text>THE CITY HALL MASS RAPID TRANSIT, MRT, TRAIN STATION AND OTHER</Text>
				 <Text>MAJOR TRANSPORTATION NODES AT ITS DOORSTEP, THIS FIVE-STAR</Text>
				 <Text>DELUXE HOTEL IS THE GATEWAY TO EXPLORE SINGAPORES LANDSCAPES AT</Text>
				 <Text>YOUR CONVENIENCE. STANDING TALL AT 226 METRES OVER 73 STOREYS,</Text>
				 <Text>THE IMPRESSIVE I.M. PEI ARCHITECTURE OFFERS THE MOST</Text>
				 <Text>MAGNIFICENT VIEWS.</Text>
				 <Text>GUESTROOMS</Text>
				 <Text>SWISSOTEL THE STAMFORD, ASIAS LEADING LUXURY CITY HOTEL,</Text>
				 <Text>SINGAPORES LEADING BUSINESS HOTEL AND THE TALLEST HOTEL IN</Text>
				 <Text>SOUTHEAST ASIA, ENCHANTS EVERY GUEST WITH UNRIVALLED COMFORT IN</Text>
				 <Text>1,252 GUESTROOMS AND LUXURIOUS SUITES WITH PRIVATE BALCONIES</Text>
				 <Text>PROVIDING BREATHTAKING PANORAMA OF THE COUNTRYS BUSTLING</Text>
				 <Text>LANDSCAPE AS WELL AS SCENIC VIEWS OF NEARBY ISLANDS OF MALAYSIA</Text>
				 <Text>AND INDONESIA. IF BUSINESS IS ON THE MIND, THE HIGH-SPEED</Text>
				 <Text>INTERNET ACCESS IN THE ROOM, AS WELL AS WIRELESS ACCESS IN THE</Text>
				 <Text>OTHER PUBLIC AREAS OF THE HOTEL, ALLOW GUESTS TO CATCH UP ON</Text>
				 <Text>E-MAILS AT THEIR CONVENIENCE. IT IS TAILORED TO MEET THE NEEDS</Text>
				 <Text>OF DISCERNING BUSINESS AND LEISURE TRAVELLERS.</Text>
				 <Text>THE STAMFORD CREST</Text>
				 <Text>FOR THE TRAVELLING EXECUTIVE, THE SWISS EXECUTIVE CLUB OFFERS</Text>
				 <Text>COMPREHENSIVE SERVICES TO MAKE WORKING AWAY FROM HOME SO MUCH</Text>
				 <Text>EASIER. GAIN EXCLUSIVE ACCESS TO THE CLUB LOUNGE WHERE GUESTS</Text>
				 <Text>CAN ENJOY DAILY BREAKFAST AND COMPLIMENTARY COCKTAILS EVERY</Text>
				 <Text>EVENING. FOR THOSE WITH LAVISH TASTES, LUXURIATE IN ONE OF THE</Text>
				 <Text>29 GRAND STAMFORD CREST SUITES, PERCHED ON THE TOP FLOORS OF</Text>
				 <Text>THE HOTEL. THESE ART-ADORNED SUITES OFFER THE BEST IN LIFESTYLE</Text>
				 <Text>EXPERIENCES, IPOD DOCKING STATION, BOSE HI-FI SYSTEM, DVD</Text>
				 <Text>PLAYER, AND AN ENSUITE BATHROOM WITH BUILT-IN TELEVISION. TRULY</Text>
				 <Text>A TOP-CLASS CHOICE FOR ACCOMMODATION, GUESTS CAN EXPECT</Text>
				 <Text>PERSONAL BUTLER, BUSINESS AND PRINTING SERVICES, ACCESS TO</Text>
				 <Text>MEETING ROOM FACILITIES AND ENJOY COMPLIMENTARY BREAKFAST AND</Text>
				 <Text>EVEN</Text>
				</Description>
				<Dining>
				 <Text>EMBARK ON A VOYAGE OF EPICUREAN AND ENTERTAINMENT PLEASURES AT</Text>
				 <Text>SWISSOTEL THE STAMFORD 16 RESTAURANTS AND BARS. INDULGE IN THE</Text>
				 <Text>FINER LIFE AT EQUINOX COMPLEX WITH 5 TIMELESS RESTAURANTS AND</Text>
				 <Text>BARS, AND 4 PRIVATE DINING ROOMS  WITH A CAPACITY OF 900</Text>
				 <Text>GUESTS. BE IT FRENCH NOUVELLE CUISINE, INTERNATIONAL FEASTS,</Text>
				 <Text>LOCAL DELIGHTS OR SAVORY CREPES, WE HAVE IT ALL. END THE NIGHT</Text>
				 <Text>ON A HIGH NOTE WITH COCKTAILS AT CITY SPACE AND INTROBAR OR</Text>
				 <Text>PARTY THE NIGHT AWAY AT THE ULTRA SEXY NEW ASIA.</Text>
				 <Text>ON PROPERTY:</Text>
				 <Text>CAFE SWISS</Text>
				 <Text>FILLED WITH CHARMING SWISS HOSPITALITY, CAF SWISS IS A</Text>
				 <Text>210-SEATER DINING RESTAURANT THAT IS DESIGNED IN A STYLISH</Text>
				 <Text>MODERN ARCHITECTURE. ILLUMINATED BY AN OVERHEAD NATURAL</Text>
				 <Text>SKYLIGHT, CAF SWISS EMANATES AN INVITING AURA OF WARMTH AND</Text>
				 <Text>ELEGANCE FOR A TRANQUIL RESPITE. OUR EXTENSIVE BUFFET LUNCH AND</Text>
				 <Text>DINNER OFFERS A VARIETY OF EUROPEAN FARE AND IS IMMENSELY</Text>
				 <Text>POPULAR.</Text>
				 <Text>KOPI TIAM</Text>
				 <Text>EXPLORE THE MULTI-CULTURAL FLAVOURS OF SINGAPORE AT KOPI TIAM,</Text>
				 <Text>UNDOUBTEDLY THE PERFECT PLACE FOR LOCALS AND TOURISTS. SAVOUR</Text>
				 <Text>SPECIALITIES SUCH AS CHILLI CRAB AND HAINANESE CHICKEN RICE OR</Text>
				 <Text>ALL-TIME FAVOURITES LIKE THE FISH HEAD CURRY AND LAKSA WHICH</Text>
				 <Text>CONTINUE TO DELIGHT THOSE CRAVING FOR SINGAPORES SIGNATURE</Text>
				 <Text>DISHES.</Text>
				 <Text>KOPI TIAM ALSO OFFERS A COLLECTION OF SPECIALTY CURRIES SUCH AS</Text>
				 <Text>DEVILS CURRY, MALAY CURRY, VINDALOO AND CHINESE CURRY AMONGST</Text>
				 <Text>OTHERS. SET IN A NOSTALGIC AND AIR-CONDITIONED AMBIENCE, KOPI</Text>
				 <Text>TIAM PROMOTES AUTHENTIC AND APPETISING LOCAL FARE WITH THE</Text>
				 <Text>FINEST QUALITY AND VALUE.</Text>
				 <Text>INTRO BAR</Text>
				 <Text>ELEGANT AND COZY, INTROBAR IS THE IDEAL ENTERTAINMENT AND</Text>
				 <Text>MEETING ALCOVE WITH A MELANGE OF MARTINIS, COCKTAILS, WINES AND</Text>
				 <Text>INTERNATIONAL SPIRITS. THE PERFECT PRELUDE TO THE EXHILARATING</Text>
				 <Text>EQUINOX EXPERIENCE, GUESTS MAY LAVISH ON THEIR FAVORITE</Text>
				 <Text>APERITIFS AT THIS PLUSH LOUNGE BEFORE EMBARKING ON A CULINARY</Text>
				 <Text>TREAT AT THE AWARD-WINNING EQUINOX RESTAURANT OR JAAN.</Text>
				 <Text>THE STAMFORD BRASSERIE</Text>
				 <Text>JAAN</Text>
				 <Text>DERIVED FROM THE ANCIENT SANSKRIT WORD FOR BOWL, JAAN IS AN</Text>
				 <Text>INTIMATE, 40-SEAT RESTAURANT DEDICATED TO SHOWCASING THE FINE</Text>
				</Dining>
				<Directions>
				 <Text>CHANGI INTL AIRPORT</Text>
				 <Text>EMBARKING FROM THE AIRPORT, TAKE THE EAST COAST PARKWAY</Text>
				 <Text>EXPRESSWAY ECP AND TAKE THE EXIT AT ROCHOR ROAD EXIT. TURN LEFT</Text>
				 <Text>AT THE 1ST JUNCTION INTO BEACH RD AND TURN RIGHT INTO STAMFORD</Text>
				 <Text>ROAD. THE HOTEL IS LOCATED ON THE JUNCTION OF BEACH ROAD AND</Text>
				 <Text>STAMFORD ROAD.</Text>
				</Directions>
				<Facilities>
				 <Text>WE HAVE THE BUSINESS CENTER, ASIA LARGEST SPA WILLOW STREAM</Text>
				 <Text>SPA, FITNESS CENTRE OUTDOOR POOL,  ON CALL VALET, 24HR ROOM</Text>
				 <Text>SERVICE, CONCIERGE, SHOPPING ARCADES, TRAIN STATION MRT, BEAUTY</Text>
				 <Text>SALON, BANK, FOREIGN EXCHANGE, CONVENTION CENTER, MEDICAL</Text>
				 <Text>SERVICE 24HR, MAJOR CREDIT CARDS ACCEPTED. PARKING FACILITIES</Text>
				 <Text>ARE AVAILABLE AT RAFFLES CITY SHOPPING CENTERS BASEMENT AREAS</Text>
				 <Text>AT PREVAILING CAR PARK FEES. VALET PARKING MAY ALSO BE AVAILED</Text>
				 <Text>AT SGD40.00 FOR 24-HOURS.</Text>
				 <Text>1 NIGHT CLUB S                15 RESTAURANT S</Text>
				 <Text>4 LOUNGE S                    AIR CONDITIONING</Text>
				 <Text>CASH MACHINE ATM ONSITE       COFFEE SHOP</Text>
				 <Text>CONVIENENCE STORE             ELEVATOR</Text>
				 <Text>EXECUTIVE FLOOR               FRONT DESK HOURS 0000-2359</Text>
				 <Text>GIFT SHOP                     GUEST LAUNDRY FACILITY</Text>
				 <Text>HANDICAP PARKING              HOTEL SAFE DEPOSIT BOX</Text>
				 <Text>HOUSEKEEPING                  PARKING CONTROLLED ACCESS</Text>
				 <Text>SMOKE FREE PROPERTY           STORAGE SPACE AVAILABLE</Text>
				 <Text>YEAR BUILT 1986</Text>
				 <Text>BUILDING MEETS LOCAL, STATE AND COUNTRY BUILDING CODES</Text>
				 <Text>PHYSICALLY CHALLENGED PUBLIC AREAS</Text>
				 <Text>YEAR PUBLIC AREAS REFURBISHED LAST 2006</Text>
				</Facilities>
				<Guarantee>
				 <Text>ALL RESERVATIONS MUST BE GUARANTEED WITH A VALID CREDIT CARD.</Text>
				</Guarantee>
				<Location>
				 <Text>STRATEGICALLY LOCATED IN THE HEART OF THE CITY, THE HOTEL IS</Text>
				 <Text>MINUTES AWAY FROM THE BUSTLING COMMERCIAL AND BANKING DISTRICT</Text>
				 <Text>OF SHENTON WAY. SWISSOTEL THE STAMFORD IS ALSO SURROUNDED BY</Text>
				 <Text>THEATERS, ART GALLERIES, MUSEUMS AND HISTORICAL SIGHTS. JUST 20</Text>
				 <Text>MINUTES AWAY FROM SINGAPORE CHANGI INTERNATIONAL AIRPORT AND</Text>
				 <Text>SINGAPORE EXPO, THE HOTEL IS EASILY ACCESSIBLE WITH A MAJOR</Text>
				 <Text>MASS RAPID TRANSIT MRT SUBWAY STATION LOCATED UNDERNEATH THE</Text>
				 <Text>HOTEL COMPLEX. IT IS ALSO A STONES THROW AWAY FROM SUNTEC CITY.</Text>
				</Location>
				<MarketingInformation>
				 <Text>THANK YOU FOR BOOKING WITH US</Text>
				</MarketingInformation>
				<MiscServices>
				 <Text>SWISSTEL THE STAMFORD HAS ALWAYS BEEN COMMITTED TO DELIVERING</Text>
				 <Text>THE BEST EXPERIENCES FOR OUR GUESTS. TO THAT END, WE ARE</Text>
				 <Text>DELIGHTED TO SHARE THAT THE HOTEL WILL BE EMBARKING ON A GUEST</Text>
				 <Text>ROOM IMPROVEMENT PROGRAMME THAT WILL BE SET TO ENHANCE THE</Text>
				 <Text>OVERALL LUXURY AND COMFORT OF SWISSTEL THE STAMFORD. FROM APRIL</Text>
				 <Text>TO OCTOBER 2017, WE WILL BE COMMENCING PHASE 1 OF THE</Text>
				 <Text>IMPROVEMENT WORKS AND UPGRADING GUEST ROOMS ON LEVELS 7 TO 28.</Text>
				 <Text>WORKS WILL BE CARRIED OUT BETWEEN 10.00 AM TO 6.00 PM DAILY AND</Text>
				 <Text>COMPLETED IN PHASES, WITH CONCERTED CONSIDERATIONS MADE TO</Text>
				 <Text>ENSURE GUEST COMFORT AND CONVENIENCE THROUGHOUT.</Text>
				 <Text>SUBSEQUENTLY, PHASE 2 AND 3 OF THE REFURBISHMENT WILL TAKE</Text>
				 <Text>PLACE BETWEEN JANUARY TO JULY 2018 FOR GUEST ROOMS ON LEVELS 30</Text>
				 <Text>TO 48  AND JULY TO DECEMBER 2018 FOR GUEST ROOMS ON LEVELS 50</Text>
				 <Text>TO 66, RESPECTIVELY.</Text>
				 <Text>THESE PROJECTS COLLECTIVELY REINFORCE SWISSTEL THE STAMFORDS</Text>
				 <Text>DEDICATION TO CREATE EVEN BETTER CUSTOMER-FOCUSED GUEST</Text>
				 <Text>EXPERIENCES FOR FUTURE VISITS. THE REFRESHED ROOMS ARE DESIGNED</Text>
				 <Text>BY AWARD-WINNING DESIGN COMPANY, WILSON ASSOCIATES  DRAWING</Text>
				 <Text>INSPIRATION FROM ICONIC SWISS CHARACTERISTICS THAT ARE</Text>
				 <Text>PRACTICAL AND ELEGANT. THESE DEVELOPMENTS WILL REVITALISE THE</Text>
				 <Text>HOTELS STYLE AND COMFORT IN CONTEMPORARY LIVING, WHILE SETTING</Text>
				 <Text>NEW BENCHMARKS IN HOSPITALITY FOR BEING TECHNOLOGICALLY ALIGNED</Text>
				 <Text>AND THOUGHTFULLY RELEVANT TO THE MODERN DAY TRAVELER.</Text>
				</MiscServices>
				<Policies>
				 <Text>CHECK-IN TIME IS AT 1500H AND CHECK-OUT TIME IS 1200H.</Text>
				 <Text>EARLY CHECK OUT POLICY</Text>
				 <Text>NO ADDITIONAL CHARGE</Text>
				 <Text>LATE CHECK OUT POLICY</Text>
				 <Text>CHECK OUT TIME IS AT 1200H.  THEREAFTER HALF DAY CHARGE APPLIES</Text>
				 <Text>FOR CHECKOUT BEFORE 1800H.  FOR LATE CHECKOUT REQUESTS AFTER</Text>
				 <Text>1800H FULL DAY RATE APPLIES.</Text>
				 <Text>PET POLICY</Text>
				 <Text>NO PETS ALLOWED ON THE HOTEL PREMISES. IF A GUEST IS TRAVELLING</Text>
				 <Text>WITH A SERVICE ANIMAL PLEASE ADVISE THE GUEST TO CONTACT THE</Text>
				 <Text>HOTEL OR THEIR TRAVEL PROFESSIONAL. NO CHARGES FOR SERVICE</Text>
				 <Text>ANIMALS.</Text>
				 <Text>FAMILY CHILDREN POLICY</Text>
				 <Text>MAXIMUM OCCUPANCY: EITHER 3 ADULTS OR 2 ADULTS PLUS 2 CHILDREN</Text>
				 <Text>IN 1 ROOM. IF A GUEST IS TRAVELING WITH A SERVICE ANIMAL PLEASE</Text>
				 <Text>ADVISE THE GUEST TO CONTACT THE HOTEL OR THEIR TRAVEL</Text>
				 <Text>PROFESSIONAL.</Text>
				 <Text>GROUP CONDITIONS</Text>
				 <Text>COMMISSION POLICY</Text>
				 <Text>-COMMISSION PERCENT - 10</Text>
				 <Text>COMMISSION IS ONLY PAYABLE TO AUTHORISED TRAVEL AGENTS ON ROOM</Text>
				 <Text>ONLY RATE.</Text>
				</Policies>
				<Recreation>
				 <Text>SWISSOTEL THE STAMFORD OFFERS NUMEROUS OPTIONS TO UNWIND</Text>
				 <Text>INCLUDING WILLOW STREAM SPA, ONE OF ASIAS LARGEST LUXURY SPA</Text>
				 <Text>AND FITNESS FACILITY-WITH 35 TREATMENT ROOMS COMPRISING OF 3</Text>
				 <Text>VIP COUPLE SUITES, A HYDROTHERAPY ROOM, A RELAXATION LOUNGE AND</Text>
				 <Text>MEDITATION ALCOVES. ADDITIONAL LEISURE FACILITIES INCLUDE A</Text>
				 <Text>FULLY-EQUIPPED FITNESS CENTRE, TWO OUTDOOR SWIMMING POOLS AND</Text>
				 <Text>SIX TENNIS COURTS.</Text>
				 <Text>RECREATIONS ON SITE</Text>
				 <Text>CARDIO VASCULAR EQUIPMENT     EXTENSIVE HEALTH CLUB</Text>
				 <Text>JACUZZI                       MASSAGE</Text>
				 <Text>OUTDOOR POOL                  SAUNA</Text>
				 <Text>SPA                           SUN BED</Text>
				 <Text>TENNIS COURTS                 WEIGHT LIFTING EQUIPMENT</Text>
				 <Text>WHIRLPOOL</Text>
				 <Text>RECREATIONS OFF SITE</Text>
				 <Text>BEACH                         BICYCLING</Text>
				 <Text>BOATING                       BOWLING</Text>
				 <Text>CAMPING                       CASINO</Text>
				 <Text>CHILD ACTIVITIES              GOLF</Text>
				 <Text>JOGGING TRACK                 MOUNTAIN BIKING</Text>
				 <Text>SAFARI</Text>
				</Recreation>
				<Rooms>
				 <Text>REVEL IN THE WARM AND COZY EMBRACE OF 1252 CONTEMPORARY</Text>
				 <Text>GUESTROOMS AND SUITES, WITH PRIVATE BALCONIES TO TAKE IN</Text>
				 <Text>BREATHTAKING PANORAMIC VIEWS.  ALL GUESTROOMS RELISH COPIOUS</Text>
				 <Text>AMENITIES INCLUDING INTERNET ACCESS, SAFE, MINI BAR,</Text>
				 <Text>IRON-IRONING BOARD, HAIRDRYER AND IN-ROOM ENTERTAINMENT, FOR</Text>
				 <Text>THE HIGHEST DEGREE OF PAMPERING.</Text>
				 <Text>MAGNIFICENT VIEWS OF SINGAPORE AWAIT ALL GUESTS STAYING AT</Text>
				 <Text>SWISSOTEL THE STAMFORD. EACH ROOM IS TASTEFULLY FURNISHED AND</Text>
				 <Text>DECORATED, ENSURING THE PURE COMFORT FOR ITS GUESTS. A</Text>
				 <Text>COMPREHENSIVE RANGE OF AMENITIES AND SERVICES ARE ALSO</Text>
				 <Text>INCLUDED, PROVIDING ALL THE CONVENIENCES ONE MAY REQUIRE.</Text>
				 <Text>STANDARD AMENITIES IN ALL ROOMS</Text>
				 <Text>ALARM CLOCK                   BALCONY</Text>
				 <Text>BATHROBE                      BATHROOM AMENITIES</Text>
				 <Text>BATHROOM PHONE                BATHTUB AND SHOWER</Text>
				 <Text>BLACK OUT CURTAINS            CABLE TV</Text>
				 <Text>COFFEE TEA                    COLOR TV</Text>
				 <Text>DUVET                         ELECTRICAL ADAPTERS</Text>
				 <Text>ELECTRICAL OUTLET DESK        ERGONOMIC CHAIR</Text>
				 <Text>FULL SIZE MIRROR              HAIR DRYER</Text>
				 <Text>HSPD                          IRON BOARD</Text>
				 <Text>LAUDRY BASKET                 LINEN THREAD COUNT 300</Text>
				 <Text>MINI FRIDGE                   MINIBAR</Text>
				 <Text>MODEM DATAPORT                MOVIES</Text>
				 <Text>NEWS                          NEWSPAPER</Text>
				 <Text>NO ADULT CHANNELS             NUMBER OF CLOSETS 1</Text>
				 <Text>NUMBER OF PHONES 2            OUTLET VOLTAGE 220</Text>
				 <Text>PHONE TWO LINES               PILLOW NONALLERGENIC</Text>
				 <Text>PILLOW TYPE FE                POWER CONVERTERS</Text>
				 <Text>PRIVATE BATH                  RADIO</Text>
				 <Text>REMOTE TV                     SAFE CHARGE 0.00</Text>
				 <Text>SAFE FOR LAPTOP               SAFE</Text>
				 <Text>SELF CONTROLLED HEATING       SEPARATE TOILET</Text>
				 <Text>SHOWER ONLY                   SPARE ELECTRIC OUTLET</Text>
				 <Text>STEREO</Text>
				</Rooms>
				<Safety>
				 <Text>24 HOUR SECURITY              ACCESSIBLE ELEVATORS</Text>
				 <Text>ALARMS CONTINUOUSLY MONITORED AUDIBLE ALARMS</Text>
				 <Text>AUDIBLE SMOKE ALARM HARDWIRED AUDIBLE SMOKE ALARMS IN ROOM</Text>
				 <Text>AUTO RECALL ELEVATORS         AUTOLINK TO FIRE DEPARTMENT</Text>
				 <Text>AUTOMATIC FIRE DOORS          DEADBOLTS ON ROOM DOORS</Text>
				 <Text>ELECTRONIC KEY CARDS          EMERGENCY BACKUP GENERATORS</Text>
				 <Text>EMERGENCY EVACUATION PLAN     EMERGENCY INFO IN ROOMS</Text>
				 <Text>EMERGENCY LIGHTING            EXIT MAPS IN ROOM</Text>
				 <Text>EXIT SIGNS LIT                EXTERIOR DOORS LOCK</Text>
				 <Text>FIRE DETECTORS IN HALLWAYS    FIRE DETECTORS IN ROOM</Text>
				 <Text>FIRST AID                     GUEST ROOM DOORS SELF CLOSING</Text>
				 <Text>OVERNIGHT SECURITY            PARKING AREAS WELL LIT</Text>
				 <Text>PARKING ATTENDANT             PATROLLED PARKING AREAS</Text>
				 <Text>PUBLIC ADDRESS SYSTEM         ROOM WINDOWS OPEN</Text>
				 <Text>SECONDARY LOCKS ON WINDOW     SECURED FLOORS</Text>
				 <Text>SECURITY ESCORT AVAILABLE     SECURITY PERSONNEL ON SITE</Text>
				 <Text>SPRINKLERS IN HALL            SPRINKLERS IN PUBLIC AREAS</Text>
				 <Text>SPRINKLERS IN ROOMS           STAFF CPR TRAINED</Text>
				 <Text>STAFF RED CROSS CERTIFIED CPR STAFF TRAINED DUPLICATE KEYS</Text>
				 <Text>STAFF TRAINED IN FIRST AID    STAFF TRAINED ON AED</Text>
				 <Text>SWINGBOLT LOCK                UNIFORM SECURITY ON PREMISES</Text>
				 <Text>VENTILATED STAIR WELLS        VIDEO SURVEILLANCE ENTRANCE</Text>
				 <Text>VIEW PORT                     WELL LIT WALKWAYS</Text>
				 <Text>AUDIBLE SMOKE ALARM IN PUBLIC AREAS</Text>
				 <Text>AUDIBLE SMOKE ALARMS IN HALLWAY</Text>
				 <Text>AUTOMATED EXTERNAL DEFIBRILLATOR ON SITE</Text>
				 <Text>BASIC MEDICAL EQUIPMENT ONSITE</Text>
				 <Text>EMERGENCY EVACUATION DRILL FREQUENCY 2</Text>
				 <Text>EMERGENCY SERVICE RESPONSE TIME 3 MINUTES</Text>
				 <Text>FIRE DETECTORS IN PUBLIC AREAS</Text>
				 <Text>FIRE EXTINGUISHERS IN HALLWAYS</Text>
				 <Text>GUEST ROOM DOORS HAVE SECOND LOCK</Text>
				 <Text>MULTIPLE FIRE EXITS EACH FLOOR</Text>
				 <Text>PROPERTY MEETS REQUIREMENT FOR FIRE SAFETY</Text>
				 <Text>SAFETY CHAIN ON GUEST ROOM DOORS</Text>
				 <Text>SECONDARY LOCKS ON SLIDING GLASS DOORS</Text>
				 <Text>VIDEO SURVEILLANCE MONITORED 24H</Text>
				 <Text>VIDEO SURVEILLANCE PARKING AREAS</Text>
				 <Text>VIDEO SURVEILLANCE PUBLIC AREAS</Text>
				 <Text>VIDEO SURVEILLANCE R</Text>
				</Safety>
				<Services>
				 <Text>DISCOVER THE EPITOME OF COMFORT IN THE TRANQUIL SANCTUARY THAT</Text>
				 <Text>IS UNDENIABLY SWISSOTEL THE STAMFORD. OFFERING UNPARALLELED</Text>
				 <Text>PANORAMIC VIEWS FROM THE BALCONIES OF 1261 GUESTROOMS AND</Text>
				 <Text>SUITES, INCLUDING 28 STAMFORD CREST SUITES, COUPLED WITH TOP OF</Text>
				 <Text>THE LINE FURNISHINGS AND AMENITIES, YOU WILL FIND ALL THE</Text>
				 <Text>ESSENTIALS FOR AN ULTIMATE PAMPERING AND A SEAMLESS WORK</Text>
				 <Text>AMBIANCE EXPERIENCE. FOR THE TRAVELING EXECUTIVE, THE SWISS</Text>
				 <Text>EXECUTIVE CLUB PUTS THE WORLD AT YOUR FINGERTIPS WITH</Text>
				 <Text>PRIVILEGES INCLUDING ACCESS TO THE LOUNGE WITH BREAKFAST AND</Text>
				 <Text>EVENING COCKTAILS AND BUSINESS CENTER FACILITIES. FOR THE BEST</Text>
				 <Text>IN LIFESTYLE EXPERIENCES, LUXURIATE IN 28 STAMFORD CREST SUITES</Text>
				 <Text>WITH EXCLUSIVE CHECK-IN AND CHECK OUT- ACCESS TO THE LIVING</Text>
				 <Text>ROOM AND PRIVATE GYM, USAGE OF MEETING ROOM, COMPLIMENTARY</Text>
				 <Text>BREAKFAST AND COCKTAILS AND MANY MORE. WITH 845 CLASSIC ROOMS,</Text>
				 <Text>255 CLASSIC HARBOR VIEW ROOMS AND 44 GRAND ROOMS, RELISH</Text>
				 <Text>COPIOUS AMENITIES INCLUDING INTERNET ACCESS, SAFE, MINI BAR,</Text>
				 <Text>IRON-IRONING BOARD, HAIRDRYER AND IN-ROOM ENTERTAINMENT. EMBARK</Text>
				 <Text>ON A VOYAGE OF CULINARY PLEASURES AT OUR COLLECTION OF 16 WORLD</Text>
				 <Text>CLASS RESTAURANTS AND BARS WITH SPECIALTIES RANGING FROM</Text>
				 <Text>SOUTHERN FRENCH, INTERNATIONAL AND SWISS TO LOCAL DELICACIES.</Text>
				 <Text>SOARING HIGH AT LEVEL 70, EQUINOX COMPLEX IS SINGAPORES MOST</Text>
				 <Text>EXCITING DINING AND ENTERTAINMENT HUB THAT REDEFINES LUXURY</Text>
				 <Text>LIFESTYLE. AS ASIA PACIFIC BEST BUSINESS HOTEL, SWISSOTEL THE</Text>
				 <Text>STAMFORD HOUSES THE CUTTING EDGE 70000 SQ FEET RAFFLES CITY</Text>
				 <Text>CONVENTION CENTRE, CAPABLE OF HOSTING EVENTS OF ALL SIZES FROM</Text>
				 <Text>INTIMATE PRIVATE FUNCTIONS TO LARGE SCALE CONVENTIONS. AFTER</Text>
				 <Text>ALL THE EXCITEMENT, RECHARGE AT THE 50000 SQ FEET WILLOW STREAM</Text>
				 <Text>SPA OR IN OUR 24-HOUR STATE-OF-THE-ART FITNESS FACILITY, TWO</Text>
				 <Text>SWIMMING POOLS AND SIX TENNIS COURTS. AT SWISSOTEL THE STAMFORD</Text>
				 <Text>YOU ARE ASSURED OF THE HIGHEST STANDARDS OF QUALITY,</Text>
				 <Text>RELIABILITY AND GENUINE SWISS HOSPITALITY.</Text>
				 <Text>AIRLINE DESK                  AIRPORT 1 SHUTTLE</Text>
				 <Text>AV EQUIPMENT                  BEAUTY SHOP</Text>
				 <Text>BELLMAN</Text>
				</Services>
				<Transportation>
				 <Text>CHANGI INTL AIRPORT</Text>
				 <Text>BUS</Text>
				 <Text>ONE WAY 1.90</Text>
				 <Text>ROUND TRIP 3.80</Text>
				 <Text>BUS STOPS AT CAPITOL BUILDING WHICH IS OUTSIDE THE HOTEL</Text>
				 <Text>CAR RENTAL</Text>
				 <Text>ONE WAY 120.00</Text>
				 <Text>DEPENDS ON CAR RENTAL COMPANY CHARGES WHICH VARIES</Text>
				 <Text>LIMO</Text>
				 <Text>ONE WAY 160.50</Text>
				 <Text>ROUND TRIP 321.00</Text>
				 <Text>TRIPS FROM  7AM-9.59PM. MIDNIGHT SURCHARGE SGD10.70 FOR TRIPS</Text>
				 <Text>FROM 10PM-6.59AM. INCL GST.</Text>
				 <Text>METROSUBWAY</Text>
				 <Text>ONE WAY 1.40</Text>
				 <Text>ROUND TRIP 2.80</Text>
				 <Text>METRO-SUBWAY MRT STATION BENEATH OUR HOTEL IS CITY HALL STATION</Text>
				 <Text>TAXI</Text>
				 <Text>ONE WAY 25.00</Text>
				 <Text>MIDNIGHT CHARGE-ADDITIONAL 50 PCT ATOP METERED FARE</Text>
				</Transportation>
			   </VendorMessages>
			  </BasicPropertyInfo>
			  <Guarantee>
			   <DepositsAccepted>
				<PaymentCard Code="AX" Type="AMERICAN EXPRESS"/>
				<PaymentCard Code="CA" Type="MASTERCARD"/>
				<PaymentCard Code="DC" Type="DINERS CLUB CARD"/>
				<PaymentCard Code="IK" Type="MASTER CARD"/>
				<PaymentCard Code="JC" Type="JCB CREDIT CARD"/>
				<PaymentCard Code="MC" Type="MASTER CARD"/>
				<PaymentCard Code="VI" Type="VISA"/>
				<PaymentCard Code="VS" Type="VISA"/>
			   </DepositsAccepted>
			   <GuaranteesAccepted>
				<PaymentCard Code="AX" Type="AMERICAN EXPRESS"/>
				<PaymentCard Code="CA" Type="MASTERCARD"/>
				<PaymentCard Code="DC" Type="DINERS CLUB CARD"/>
				<PaymentCard Code="IK" Type="MASTER CARD"/>
				<PaymentCard Code="JC" Type="JCB CREDIT CARD"/>
				<PaymentCard Code="MC" Type="MASTER CARD"/>
				<PaymentCard Code="VI" Type="VISA"/>
				<PaymentCard Code="VS" Type="VISA"/>
				<Text>/GAGT             GUARANTEE TO AGENCY TIDS/IATA NUMBER</Text>
				<Text>/GDPST            DEPOSIT WILL BE SENT</Text>
				<Text>/GDPST...         IMMEDIATE DEPOSIT TYPE OR FORM</Text>
			   </GuaranteesAccepted>
			  </Guarantee>
			  <RoomRates>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="P1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="001" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER KING, NEWLY RENOVATED, 40SQM</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAIN SHOWER</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="285.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="335.45">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="28.50"/>
				   <TotalTaxes Amount="21.95"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="D1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="002" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>CLASSIC KING, 40SQM/430SQF</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="285.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="335.45">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="28.50"/>
				   <TotalTaxes Amount="21.95"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="L1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="003" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER HARBOUR VIEW KING, NEWLY RENOVATED</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAINSHOWER, 40SQM</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="370.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="435.49">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="37.00"/>
				   <TotalTaxes Amount="28.49"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="P1KBRF" IATA_ProductIdentification="BREAKFAST FLEXIBLE" LowInventoryThreshold="false" RPH="004" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER KING, NEWLY RENOVATED, 40SQM</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAIN SHOWER</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="315.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="370.75">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="31.50"/>
				   <TotalTaxes Amount="24.25"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="T1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="005" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>SIGNATURE KING, 50 SQM/538SQF, FREE WIFI</Text>
				 <Text>NEWLY RENOVATED, PRIVATE BALCONY, RAIN SHOWER</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="325.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="382.53">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="32.50"/>
				   <TotalTaxes Amount="25.03"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="P2TRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="006" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER DOUBLE DOUBLE, NEWLY RENOVATED, 40SQM</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAIN SHOWER</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="285.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="335.45">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="28.50"/>
				   <TotalTaxes Amount="21.95"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="K1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="007" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>CREST SUITE, 1 KING BED, EXEC LOUNGE ACCESS</Text>
				 <Text>FREE WIFI, BREAKFAST AND COCKTAILS IN LOUNGE</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="585.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="688.55">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="58.50"/>
				   <TotalTaxes Amount="45.05"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="L2TRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="008" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER HARBOUR VIEW DOUBLE DOUBLE, RENOVATED</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAINSHOWER, 40SQM</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="370.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="435.49">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="37.00"/>
				   <TotalTaxes Amount="28.49"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="E2DRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="009" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>CLASSIC DOUBLE DOUBLE, HARBOUR VIEW</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, 40SQM/430SQF</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="295.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="347.22">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="29.50"/>
				   <TotalTaxes Amount="22.72"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="E1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="010" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>CLASSIC KING, HARBOUR VIEW, 40SQM/430SQF</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="295.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="347.22">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="29.50"/>
				   <TotalTaxes Amount="22.72"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="D2DRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="011" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>CLASSIC DOUBLE DOUBLE, 40SQM/430SQF</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="285.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="335.45">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="28.50"/>
				   <TotalTaxes Amount="21.95"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="C1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="012" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER KING PLUS, NEWLY RENOVATED, 40SQM</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAIN SHOWER</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="285.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="335.45">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="28.50"/>
				   <TotalTaxes Amount="21.95"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="U1QRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="013" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>EXECUTIVE KING, LOUNGE ACCESS BENEFITS</Text>
				 <Text>FREE WIFI, BREAKFAST AND COCKTAILS IN LOUNGE</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="390.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="459.03">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="39.00"/>
				   <TotalTaxes Amount="30.03"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="A2TRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="014" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>EXECUTIVE HARBOUR VIEW  DBL</Text>
				 <Text>EXEC LOUNGE WITH FREE BREAKFAST AND COCKTAILS</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="400.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="470.80">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="40.00"/>
				   <TotalTaxes Amount="30.80"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="N1KRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="015" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>PREMIER HARBOUR VIEW KING PLUS, NEWLY RENOVATED</Text>
				 <Text>PRIVATE BALCONY, FREE WIFI, RAINSHOWER, 40SQM</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="370.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="435.49">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="37.00"/>
				   <TotalTaxes Amount="28.49"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			   <RoomRate DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" IATA_CharacteristicIdentification="N1QRAC" IATA_ProductIdentification="FULLY FLEXIBLE-" LowInventoryThreshold="false" RPH="016" RateChangeInd="false" RateConversionInd="false" SpecialOffer="false">
				<AdditionalInfo>
				 <CancelPolicy Numeric="02" Option="D"/>
				 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
				 <Text>EXECUTIVE HARBOUR VIEW KING, EXEC LOUNGE</Text>
				 <Text>FREE WIFI, BREAKFAST AND COCKTAILS IN LOUNGE</Text>
				</AdditionalInfo>
				<Rates>
				 <Rate Amount="400.00" ChangeIndicator="false" CurrencyCode="SGD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
				  <AdditionalGuestAmounts>
				   <AdditionalGuestAmount MaxExtraPersonsAllowed="0" NumCribs="0">
					<Charges Crib="0" ExtraPerson="80.00"/>
				   </AdditionalGuestAmount>
				  </AdditionalGuestAmounts>
				  <HotelTotalPricing Amount="470.80">
				   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
				   <TotalSurcharges Amount="40.00"/>
				   <TotalTaxes Amount="30.80"/>
				  </HotelTotalPricing>
				 </Rate>
				</Rates>
			   </RoomRate>
			  </RoomRates>
			  <TimeSpan Duration="0005" End="2018-05-12" Start="2018-05-11"/>
			 </RoomStay>
			</HotelPropertyDescriptionRS></soap-env:Body></soap-env:Envelope>`)

	sampleHotelRateDescRSgood = []byte(` <?xml version="1.0" encoding="UTF-8"?>
		<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">HotelRateDescriptionLLSRQ</eb:Service><eb:Action>HotelRateDescriptionLLSRS</eb:Action><eb:MessageData><eb:MessageId>3319224857945470281</eb:MessageId><eb:Timestamp>2018-05-14T23:49:54</eb:Timestamp><eb:RefToMessageId>mid:20180207-20:19:07.25|QVbg0</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><HotelRateDescriptionRS xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:stl="http://services.sabre.com/STL/v01" Version="2.3.0">
		 <stl:ApplicationResults status="Complete">
		  <stl:Success timeStamp="2018-05-14T18:49:54-05:00">
		   <stl:SystemSpecificResults>
			<stl:HostCommand LNIATA="222222">HRD*12</stl:HostCommand>
		   </stl:SystemSpecificResults>
		  </stl:Success>
		 </stl:ApplicationResults>
		 <RoomStay>
		  <BasicPropertyInfo ChainCode="WI" GeoConfidenceLevel="0" HotelCityCode="TPA" HotelCode="0000232" HotelName="WESTIN TAMPA WATERSIDE" Latitude="27.937466" Longitude="-82.453541" NumFloors="12" RPH="001">
		   <Address>
			<AddressLine>725 SOUTH HARBOR ISLAND BLVD</AddressLine>
			<AddressLine>TAMPA FL 33602</AddressLine>
			<CountryCode>US</CountryCode>
		   </Address>
		   <Awards>
			<AwardProvider>NTM4  CROWN</AwardProvider>
		   </Awards>
		   <CheckInTime>15:00</CheckInTime>
		   <CheckOutTime>12:00</CheckOutTime>
		   <ContactNumbers>
			<ContactNumber Fax="1-813-229-5322" Phone="1-813-229-5000"/>
		   </ContactNumbers>
		   <DirectConnect Ind="true">
			<AdditionalData Ind="false"/>
			<CurrencyConverted Ind="false"/>
			<DC_AvailParticipant Ind="true"/>
			<DC_SellParticipant Ind="true"/>
			<RequestFail Ind="false"/>
			<UnAvail Ind="false"/>
		   </DirectConnect>
		   <IndexData>
			<Index CountryState="FL" DistanceDirection="1SE" LocationCode="A" Point="UNIVERSITY OF TAMPA" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="35SE" LocationCode="C" Point="CLEARWATER BEACH" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="9SW" LocationCode="A" Point="BUSCH GARDENS" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="2S" LocationCode="A" Point="PORT TAMPA" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="5SE" LocationCode="A" Point="RAYMOND JAMES STADIUM" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="21SW" LocationCode="C" Point="TROPICANA FIELD" TransportationCode="O"/>
			<Index CountryState="FL" DistanceDirection="1S" LocationCode="A" Point="TAMPA CONV CTR" TransportationCode="O"/>
			<Index DistanceDirection="7SE" LocationCode="A" Point="TPA" TransportationCode="O"/>
		   </IndexData>
		   <PropertyOptionInfo>
			<ADA_Accessible Ind="true"/>
			<AdultsOnly Ind="false"/>
			<AirportShuttle Ind="false"/>
			<BeachFront Ind="false"/>
			<Breakfast Ind="false"/>
			<BusinessCenter Ind="true"/>
			<BusinessReady Ind="false"/>
			<CarRentalCounter>N</CarRentalCounter>
			<Conventions Ind="true"/>
			<Dataport Ind="false"/>
			<Dining Ind="true"/>
			<DryClean Ind="false"/>
			<EcoCertified Ind="false"/>
			<ExecutiveFloors Ind="false"/>
			<FamilyPlan Ind="false"/>
			<FitnessCenter Ind="true"/>
			<FreeLocalCalls Ind="false"/>
			<FreeParking Ind="false"/>
			<FreeShuttle Ind="false"/>
			<FreeWifiInMeetingRooms Ind="false"/>
			<FreeWifiInPublicSpaces Ind="true"/>
			<FreeWifiInRooms Ind="false"/>
			<FullServiceSpa Ind="false"/>
			<GameFacilities Ind="false"/>
			<Golf Ind="false"/>
			<GovtSafetyFire Ind="true"/>
			<HighSpeedInternet Ind="false"/>
			<HypoallergenicRooms Ind="false"/>
			<IndoorPool Ind="false"/>
			<IndPetRestriction Ind="true"/>
			<InRoomCoffeeTea Ind="true"/>
			<InRoomMiniBar Ind="true"/>
			<InRoomRefrigerator Ind="false"/>
			<InRoomSafe Ind="true"/>
			<InteriorDoorways Ind="true"/>
			<Jacuzzi Ind="false"/>
			<KidsFacilities Ind="false"/>
			<KitchenFacilities Ind="false"/>
			<MealService Ind="false"/>
			<MeetingFacilities Ind="true"/>
			<NoAdultTV Ind="false"/>
			<NonSmoking Ind="true"/>
			<OutdoorPool Ind="true"/>
			<Parking>Y</Parking>
			<Pets Ind="true"/>
			<Pool Ind="true"/>
			<PublicTransportationAdjacent Ind="false"/>
			<Recreation Ind="true"/>
			<RestrictedRoomAccess Ind="false"/>
			<RoomService Ind="false"/>
			<RoomService24Hours Ind="false"/>
			<RoomsWithBalcony Ind="false"/>
			<SkiInOutProperty Ind="false"/>
			<SmokeFree Ind="true"/>
			<SmokingRoomsAvail Ind="false"/>
			<Tennis Ind="false"/>
			<WaterPurificationSystem Ind="false"/>
			<Wheelchair Ind="true"/>
		   </PropertyOptionInfo>
		   <PropertyTypeInfo>
			<AllInclusive Ind="false"/>
			<Apartments Ind="false"/>
			<BedBreakfast Ind="false"/>
			<Castle Ind="false"/>
			<Conventions Ind="true"/>
			<Economy Ind="false"/>
			<ExtendedStay Ind="false"/>
			<Farm Ind="false"/>
			<First Ind="true"/>
			<Luxury Ind="false"/>
			<Moderate Ind="false"/>
			<Motel Ind="false"/>
			<Resort Ind="false"/>
			<Suites Ind="false"/>
		   </PropertyTypeInfo>
		   <SpecialOffers Ind="false"/>
		   <Taxes>
			<Text>5.00 PCT</Text>
			<Text>7.00 PCT</Text>
		   </Taxes>
		   <VendorMessages>
			<AdditionalAttractions>
			 <Text>CORPORATE LOCATIONS -</Text>
			 <Text>WELLS FARGO BANK               1.00 MI</Text>
			 <Text>FEDEX OFFICE                   1.00 MI</Text>
			 <Text>SYKES ENTERPRISES              1.00 MI</Text>
			 <Text>ERNST   YOUNG                  1.00 MI</Text>
			 <Text>BANK OF AMERICA                1.00 MI</Text>
			 <Text>TECO ENERGY, INC.              1.00 MI</Text>
			 <Text>AACSB                          0.06 MI</Text>
			 <Text>VERIZON                        1.00 MI</Text>
			 <Text>CITIGROUP                      10.99 MI</Text>
			 <Text>MEDIA GENERAL BROADCASTING GRO 1.00 MI</Text>
			 <Text>TAMPA PORT AUTHORITY           2.00 MI</Text>
			 <Text>KPMG                           1.00 MI</Text>
			 <Text>CP SHIPS                       1.00 MI</Text>
			 <Text>PRICEWATERHOUSECOOPERS         1.00 MI</Text>
			</AdditionalAttractions>
			<Attractions>
			 <Text>ATTRACTIONS -</Text>
			 <Text>AMTRAK                       1.00 MI</Text>
			 <Text>SALVIDOR DALI MUSEUM         24.98 MI</Text>
			 <Text>FORT DE SOTO PARK            33.98 MI</Text>
			 <Text>UNIVERSITY OF TAMPA          1.00 MI</Text>
			 <Text>BAYSHORE BOULEVARD LINEAR PA 1.00 MI</Text>
			 <Text>UNIVERSITY OF SOUTH FLORIDA  8.00 MI</Text>
			 <Text>LEGEND S FIELD               5.00 MI</Text>
			 <Text>TAMPA BAY PERFORMING ARTS CE 1.00 MI</Text>
			 <Text>LOWRY PARK ZOO               8.00 MI</Text>
			 <Text>WESTSHORE PLAZA              6.00 MI</Text>
			 <Text>TAMPA CONVENTION CENTER      0.12 MI</Text>
			 <Text>OLD HYDE PARK VILLAGE        2.00 MI</Text>
			 <Text>AMALIE ARENA  FORMERLY ST. P 0.19 MI</Text>
			 <Text>TAMPA MUSEUM OF ART          1.00 MI</Text>
			 <Text>BUSCH GARDENS                8.99 MI</Text>
			 <Text>HB PLANT MUSEUM              1.00 MI</Text>
			 <Text>YACHT STARSHIP               1.00 MI</Text>
			 <Text>TAMPA THEATRE                1.00 MI</Text>
			 <Text>CHANNELSIDE                  1.00 MI</Text>
			 <Text>TAMPA GENERAL HOSPITAL       2.00 MI</Text>
			 <Text>INTERNATIONAL PLAZA   BAY ST 6.00 MI</Text>
			 <Text>PRICEWATERHOUSECOOPERS       1.00 MI</Text>
			 <Text>GERALD R. FORD AMPHITHEATRE  8.99 MI</Text>
			 <Text>RAYMOND JAMES STADIUM  HOME  5.00 MI</Text>
			 <Text>FLORIDA AQUARIUM             1.00 MI</Text>
			 <Text>PORT AUTHORITY               2.00 MI</Text>
			</Attractions>
			<Awards>
			 <Text>AWARDS -</Text>
			 <Text>GLOBAL BUSINESS TRAVEL ASSOCIATION  GBTA  - PROJECT ICARUS</Text>
			 <Text>SUPPLIER GOLD MEDAL  2016</Text>
			</Awards>
			<Cancellation>
			 <Text>01JAN14 - 31DEC99-</Text>
			 <Text>CANCEL BY 2 DAYS PRIOR TO ARRIVAL</Text>
			 <Text>TO AVOID A 100.00PCT CANCELLATION PENALTY</Text>
			 <Text>CANCELLATION POLICY TEXT -</Text>
			 <Text>THE CANCELLATION POLICY WILL VARY BASED ON THE RATE PLAN</Text>
			 <Text>AND/OR BOOKING DATE S . PLEASE SEE  RATE AND POLICY</Text>
			 <Text>INFORMATION  WHEN CHECKING AVAILABILITY.</Text>
			</Cancellation>
			<Deposit>
			 <Text>ACCEPTED FORMS OF DEPOSIT- 01JAN14 - 31DEC99</Text>
			 <Text>CREDIT CARD</Text>
			</Deposit>
			<Description>
			 <Text>YEAR BUILT - 1985           YEAR REMODELED - 2017</Text>
			 <Text>ADDITIONAL PROPERTY DESCRIPTION -</Text>
			 <Text>ESCAPE TO THE WESTIN TAMPA WATERSIDE, SITUATED ON A UNIQUE</Text>
			 <Text>LANDSCAPED ISLAND IN DOWNTOWN TAMPA. ADJACENT TO THE TAMPA</Text>
			 <Text>CONVENTION CENTER, OUR LOCATION IS CONVENIENT FOR BUSINESS</Text>
			 <Text>AND LEISURE TRAVELERS ALIKE. JUST TWO BLOCKS FROM AMALIE</Text>
			 <Text>ARENA  FORMERLY ST. PETE TIMES FORUM , WE HAVE ACCESS TO ALL</Text>
			 <Text>THAT TAMPA HAS TO OFFER.OUR WATERFRONT ATMOSPHERE LEAVES</Text>
			 <Text>NOTHING TO BE DESIRED, WITH AN OUTDOOR POOL AND A</Text>
			 <Text>FULL-SERVICE BUSINESS CENTER. HOST AN EVENT IN ONE OF OUR 13</Text>
			 <Text>FLEXIBLE MEETING ROOMS, SOME OFFERING STUNNING HARBOR VIEWS.</Text>
			 <Text>ENJOY DELICIOUS CONTINENTAL FARE AT OUR ELEGANT RESTAURANT,</Text>
			 <Text>BLUE HARBOUR EATERY  AMP  BAR. OUR OVERSIZED GUEST ROOMS AND</Text>
			 <Text>SUITES AT THE WESTIN TAMPA WATERSIDE ARE DESIGNED TO ENHANCE</Text>
			 <Text>YOUR RELAXATION. RELAX IN OUR SIGNATURE HEAVENLY  BED FOR A</Text>
			 <Text>PEACEFUL NIGHTS SLEEP. REFRESH THE NEXT MORNING IN ONE OF OUR</Text>
			 <Text>SPACIOUS BATHROOMS, EQUIPPED WITH THE HEAVENLY  SHOWER</Text>
			 <Text>FEATURING DUAL SHOWERHEADS.</Text>
			</Description>
			<Dining>
			 <Text>RESTAURANTS -</Text>
			 <Text>MARKET PLACE</Text>
			 <Text>RESTAURANT DESCRIPTION</Text>
			 <Text>STOP BY THE MARKET PLACE FOR GRAB-AND-GO OPTIONS LIKE</Text>
			 <Text>SNACKS, SODAS, AND COLD-BREW COFFEE.</Text>
			 <Text>BLUE HARBOUR EATERY   BAR</Text>
			 <Text>CUISINE - AMERICAN</Text>
			 <Text>MEALS SERVED - BREAKFAST - LUNCH - DINNER</Text>
			 <Text>RESTAURANT DESCRIPTION</Text>
			 <Text>THE RELAXED YET REFINED BLUE HARBOUR EATERY   BAR HAS BEEN</Text>
			 <Text>INSPIRED BY THE SEA. MADE FRESH FROM THE DAY S CATCH,</Text>
			 <Text>SIMPLE DISHES ARE FLAWLESSLY EXECUTED, AND A TOP-NOTCH</Text>
			 <Text>DRINK LIST INCLUDES PREMIUM COCKTAILS.</Text>
			</Dining>
			<Directions>
			 <Text>PRIMARY AIRPORT -</Text>
			 <Text>TPA - TAMPA INTERNATIONAL AIRPORT - 9.00 MI</Text>
			 <Text>OTHER AIRPORTS -</Text>
			 <Text>PIE - ST. PETERSBURG INTERNATIONAL AIRPORT - 25.00 MI</Text>
			 <Text>MCO - ORLANDO INTERNATIONAL AIRPORT - 86.00 MI</Text>
			 <Text>DIRECTIONS TO THE PROPERTY FROM THE NORTH -</Text>
			 <Text>TAKE INTERSTATE 75 TAMPA CROSSTOWN EXPRESSWAY, WESTBOUND.</Text>
			 <Text>TAKE THE CROSSTOWN EXPRESSWAY TO THE MORGAN STREET EXIT. AT</Text>
			 <Text>THE BASE OF THE RAMP, MERGE TO THE LEFT LANE. TURN LEFT ON</Text>
			 <Text>FRANKLIN STREET.</Text>
			 <Text>DIRECTIONS TO THE PROPERTY FROM THE SOUTH -</Text>
			 <Text>TAKE INTERSTATE 275 NORTH TO EXIT 44  ASHLEY STREET -</Text>
			 <Text>DOWNTOWN EAST/WEST . TAKE THE DOWNTOWN WEST EXIT. TAKE THE</Text>
			 <Text>TAMPA STREET SOUTH EXIT. TURN RIGHT ONTO FRANKLIN STREET</Text>
			 <Text>WHICH BECOMES HARBOUR ISLAND BOULEVARD. CONTINUE TO THE</Text>
			 <Text>HOTEL.</Text>
			 <Text>DIRECTIONS TO THE PROPERTY FROM THE WEST -</Text>
			 <Text>TAKE INTERSTATE 4 TO INTERSTATE 275 SOUTH. PROCEED TO EXIT 44</Text>
			 <Text>ASHLEY STREET - DOWNTOWN EAST/WEST . TAKE TAMPA STREET</Text>
			 <Text>SOUTH. TURN RIGHT ONTO FRANKLIN STREET WHICH BECOMES HARBOUR</Text>
			 <Text>ISLAND BOULEVARD. CONTINUE TO THE HOTEL.</Text>
			</Directions>
			<Facilities>
			 <Text>ON-SITE FACILITIES DESCRIPTION -</Text>
			 <Text>PARKING</Text>
			 <Text>ON-SITE FACILITIES -</Text>
			 <Text>ACCESSIBLE FACILITIES         BUSINESS CENTER</Text>
			 <Text>CONCIERGE LOUNGE              CONNECTING ROOMS</Text>
			 <Text>EXERCISE GYM                  HEALTH CLUB</Text>
			 <Text>MEETING ROOMS                 NON-SMOKING ROOMS  GENERIC</Text>
			 <Text>ONSITE LAUNDRY                OUTDOOR POOL</Text>
			 <Text>PARKING                       PARKING LOT</Text>
			 <Text>POOL                          RESTAURANT</Text>
			 <Text>WIRELESS INTERNET CONNECTION</Text>
			 <Text>MEETING AND CONVENTION FACILITIES -</Text>
			 <Text>TOTAL NBR OF MEETING ROOMS - 13</Text>
			 <Text>MAXIMUM SEATING CAPACITY OF LARGEST ROOM  - 700</Text>
			 <Text>MAXIMUM SEATING CAPACITY OF SMALLEST ROOM - 10</Text>
			 <Text>AREA OF LARGEST ROOM  - 442</Text>
			 <Text>AREA OF SMALLEST ROOM - 29</Text>
			 <Text>TOTAL MEASUREMENT, ALL ROOMS - 1672</Text>
			 <Text>MEETING ROOM SIZE AND SEATING FORMAT OPTIONS -</Text>
			 <Text>MEETING ROOM NAME - CORAL REEF</Text>
			 <Text>LENGTH: 76.00  WIDTH: 20.00  AREA: 1509.00 SQ FT</Text>
			 <Text>BANQUET - 100                THEATRE - 100</Text>
			 <Text>RECEPTION - 150</Text>
			 <Text>MEETING ROOM NAME - CHANNELSIDE</Text>
			 <Text>LENGTH: 144.00  WIDTH: 26.00  AREA: 3744.00 SQ FT</Text>
			 <Text>BANQUET - 200                THEATRE - 180</Text>
			 <Text>RECEPTION - 400</Text>
			 <Text>MEETING ROOM NAME - PALM</Text>
			 <Text>LENGTH: 18.00  WIDTH: 22.00  AREA: 396.00 SQ FT</Text>
			 <Text>CONFERENCE  - 14</Text>
			 <Text>MEETING ROOM NAME - REEF</Text>
			 <Text>LENGTH: 27.00  WIDTH: 20.00  AREA: 540.00 SQ FT</Text>
			 <Text>BANQUET - 40                 THEATRE - 50</Text>
			 <Text>RECEPTION - 60</Text>
			 <Text>MEETING ROOM NAME - CORAL</Text>
			 <Text>LENGTH: 27.00  WIDTH: 20.00  AREA: 540.00 SQ FT</Text>
			 <Text>BANQUET - 40                 THEATRE - 50</Text>
			 <Text>RECEPTION - 60</Text>
			 <Text>MEETING ROOM NAME - HIBISCUS</Text>
			 <Text>LENGTH: 25.00  WIDTH: 17.00  AREA: 425.00 SQ FT</Text>
			 <Text>BANQUET - 24                 THEATRE - 50</Text>
			 <Text>RECEPTION - 50</Text>
			 <Text>MEETING ROOM NAME - LAGOON</Text>
			 <Text>LENGTH: 20.00  WIDTH: 16.00  AREA: 320.00 SQ FT</Text>
			 <Text>BANQUET - 12                 THEATRE - 40</Text>
			 <Text>RECEPTION - 30</Text>
			 <Text>MEETING ROOM NAME - SUNSET</Text>
			 <Text>LENGTH: 35.00  WIDTH: 25.00  AREA: 875.00 SQ FT</Text>
			 <Text>BANQUET - 70                 THEATRE - 80</Text>
			 <Text>RECEPTION</Text>
			</Facilities>
			<Guarantee>
			 <Text>-01JAN14 - 31DEC99 MON-SUN</Text>
			 <Text>CREDIT CARD GUARANTEE MAY BE REQUIRED ON RESERVATIONS.</Text>
			 <Text>ACCEPTED FORMS OF GUARANTEE- 01JAN14 - 31DEC99</Text>
			 <Text>AGENCY IATA/ARC               CREDIT CARD</Text>
			 <Text>DEPOSIT</Text>
			</Guarantee>
			<Location>
			 <Text>PRIMARY PROPERTY LOCATION - CITY</Text>
			 <Text>DISTRICT -</Text>
			 <Text>AREA     -</Text>
			 <Text>PRIMARY AIRPORT -</Text>
			 <Text>TPA          9.00 MILE</Text>
			 <Text>ADDITIONAL AIRPORTS -</Text>
			 <Text>PIE          25.00 MILE</Text>
			 <Text>MCO          86.00 MILE</Text>
			 <Text>PRIMARY CITY -</Text>
			 <Text>TPA - CITY CENTER 0.72 MILE</Text>
			 <Text>EXPLORE AND DISCOVER OUR SURROUNDINGS.</Text>
			</Location>
			<MarketingInformation>
			 <Text>WESTIN HOTELS FOR A BETTER YOU</Text>
			 <Text>MEMBER OF STARWOOD HOTELS</Text>
			</MarketingInformation>
			<MiscServices>
			 <Text>EMAIL - RESERVATIONS WESTINTAMPAWATERSIDE.COM</Text>
			 <Text>TIME ZONE - EDT</Text>
			 <Text>TAXES AND SURCHARGES -</Text>
			 <Text>5.00 PCT OCCUPANCY TAX 11MAY18 - 31DEC99</Text>
			 <Text>7.00 PCT CITY TAX 11MAY18 - 31DEC99</Text>
			</MiscServices>
			<Policies>
			 <Text>EXTRA CHILD - 0 USD</Text>
			 <Text>CHILDREN STAY FREE</Text>
			 <Text>ROOM RATES INCLUDE THE ACCOMMODATION OF CHILDREN  17 YEARS</Text>
			 <Text>OLD OR YOUNGER  WHO SLEEP IN THE EXISTING BEDDING OF A GUEST</Text>
			 <Text>ROOM. ROLLAWAY BEDS AND CRIBS MAY INCUR EXTRA CHARGES.</Text>
			 <Text>FAMILY PLAN</Text>
			 <Text>ROOM RATES INCLUDE THE ACCOMMODATION OF CHILDREN  17 YEARS</Text>
			 <Text>OLD OR YOUNGER  WHO SLEEP IN THE EXISTING BEDDING OF A GUEST</Text>
			 <Text>ROOM. ROLLAWAY BEDS AND CRIBS MAY INCUR EXTRA CHARGES.</Text>
			 <Text>PETS ALLOWED</Text>
			 <Text>DOGS UP TO 40 POUNDS ARE ALLOWED. NO OTHER PETS ARE</Text>
			 <Text>PERMITTED. A 50 USD NON-REFUNDABLE CLEANING FEE WILL BE</Text>
			 <Text>CHARGED. PETS ARE NOT ALLOWED IN SUITES AND ARE ONLY</Text>
			 <Text>PERMITTED IN TRADITIONAL KING AND DOUBLE ROOM TYPES.  OWNERS</Text>
			 <Text>MUST SIGN A WAIVER AT CHECK-IN AND ARE RESPONSIBLE FOR ANY</Text>
			 <Text>DAMAGE OR ADDITIONAL CLEANING REQUIRED.  DOGS MUST BE</Text>
			 <Text>LEASHED AT ALL TIMES AND ARE NOT PERMITTED IN DINING AND</Text>
			 <Text>RECREATION AREAS.  DOGS MAY NOT BE LEFT UNATTENDED IN THE</Text>
			 <Text>GUESTROOM.  GUESTS MUST BE IN THEIR GUESTROOM IN ORDER TO</Text>
			 <Text>RECEIVE HOUSEKEEPING SERVICES.</Text>
			 <Text>PROPERTY OFFERS COMMISSION FOR SOME RATES</Text>
			 <Text>PROPERTY PARTICIPATES IN CENTRALIZED COMMISSION PROCESSING</Text>
			 <Text>WPS</Text>
			 <Text>ADDITIONAL COMMISSION INFORMATION</Text>
			 <Text>NOT ALL RATES COMMISSIONABLE. CHECK RATE PLAN DETAILS FOR</Text>
			 <Text>RATE COMMISSION INFORMATION.</Text>
			</Policies>
			<Recreation>
			 <Text>ON-SITE RECREATION -</Text>
			 <Text>FITNESS CENTER ON-SITE        POOL</Text>
			</Recreation>
			<Rooms>
			 <Text>GENERAL ROOM DESCRIPTION</Text>
			 <Text>ESCAPE TO A SUPERBLY COMFORTABLE SPACE. REFRESHING TOUCHES</Text>
			 <Text>WILL AWAKEN YOUR SENSES.</Text>
			 <Text>ROOM AMENITIES FOR ALL ROOMS</Text>
			 <Text>ADJOINING ROOMS                CABLE TELEVISION</Text>
			 <Text>COFFEE/TEA MAKER               CONNECTING ROOMS</Text>
			 <Text>CRIBS                          DESK</Text>
			 <Text>ERGONOMIC CHAIR                HAIRDRYER</Text>
			 <Text>IRON                           NON-SMOKING</Text>
			 <Text>SAFE                           SMOKE DETECTORS</Text>
			</Rooms>
			<Safety>
			 <Text>PROPERTY IS FIRE SAFETY COMPLIANT</Text>
			</Safety>
			<Services>
			 <Text>ON-SITE SERVICES -</Text>
			 <Text>BAGGAGE HOLD                  PETS ALLOWED</Text>
			 <Text>SAFE DEPOSIT BOX</Text>
			</Services>
			<Transportation>
			 <Text>TRANSPORTATION FROM PRIMARY AIRPORT TPA -</Text>
			 <Text>SHUTTLE - USD 13.0</Text>
			 <Text>FEE IS PER WAY, PER PERSON. THE HOTEL DOES NOT OFFER</Text>
			 <Text>COMPLIMENTARY TRANSPORTATION SERVICES.</Text>
			 <Text>TAXI - USD 25.0</Text>
			 <Text>TRANSPORTATION FROM AIRPORTS PIE -</Text>
			 <Text>TAXI - USD 40.0</Text>
			 <Text>LIMOUSINE - USD 75.0</Text>
			 <Text>TRANSPORTATION FROM AIRPORTS MCO -</Text>
			 <Text>TAXI - USD 200.0</Text>
			</Transportation>
		   </VendorMessages>
		  </BasicPropertyInfo>
		  <Guarantee>
		   <DepositsAccepted>
			<PaymentCard Code="AX" Type="AMERICAN EXPRESS"/>
			<PaymentCard Code="CA" Type="MASTERCARD"/>
			<PaymentCard Code="CB" Type="CARTE BLANCHE"/>
			<PaymentCard Code="DC" Type="DINERS CLUB CARD"/>
			<PaymentCard Code="DS" Type="DISCOVER CARD"/>
			<PaymentCard Code="JC" Type="JCB CREDIT CARD"/>
			<PaymentCard Code="VI" Type="VISA"/>
		   </DepositsAccepted>
		   <GuaranteesAccepted>
			<PaymentCard Code="AX" Type="AMERICAN EXPRESS"/>
			<PaymentCard Code="CA" Type="MASTERCARD"/>
			<PaymentCard Code="CB" Type="CARTE BLANCHE"/>
			<PaymentCard Code="DC" Type="DINERS CLUB CARD"/>
			<PaymentCard Code="DS" Type="DISCOVER CARD"/>
			<PaymentCard Code="JC" Type="JCB CREDIT CARD"/>
			<PaymentCard Code="VI" Type="VISA"/>
			<Text>/GAGT             GUARANTEE TO AGENCY TIDS/IATA NUMBER</Text>
			<Text>/GDPST            DEPOSIT WILL BE SENT</Text>
			<Text>/GDPST...         IMMEDIATE DEPOSIT TYPE OR FORM</Text>
		   </GuaranteesAccepted>
		  </Guarantee>
		  <RoomRates>
		   <RoomRate ClientID="¤¤¤" DirectConnect="false" GuaranteeSurchargeRequired="G" GuaranteedRateProgram="false" HRD_RequiredForSell="false" IATA_CharacteristicIdentification="J1KA16" IATA_ProductIdentification="1KING:AAA/CAA RATE" LowInventoryThreshold="false" RateAccessCode="¤¤¤" RateCategory="¤" RateChangeInd="false" RateConversionInd="false" RoomLocationCode="J1" SpecialOffer="false">
			<AdditionalInfo>
			 <CancelPolicy Numeric="02" Option="D"/>
			 <Commission NonCommission="false">10.00 PERCENT COMMISSION</Commission>
			 <DCA_Cancellation>
			  <Text>2 DAYS-PRIOR 1 NTS PENALTY</Text>
			 </DCA_Cancellation>
			 <DCA_Guarantee>
			  <Text>GUARANTEE REQRD- MAJOR CREDIT CARDS.</Text>
			 </DCA_Guarantee>
			 <PaymentCard Code="DS"/>
			 <PaymentCard Code="CA"/>
			 <PaymentCard Code="MC"/>
			 <PaymentCard Code="CB"/>
			 <PaymentCard Code="VI"/>
			 <PaymentCard Code="VS"/>
			 <PaymentCard Code="AX"/>
			 <PaymentCard Code="JC"/>
			 <PaymentCard Code="DC"/>
			 <Taxes>TAXES NOT INCLUDED IN ROOM RAT</Taxes>
			 <Text>AAA OR CAA MEMBERSHIP ID REQUIRED AT CHECK-IN.</Text>
			 <Text>NON SMOKING WATERVIEW: HIGH FLOOR ROOM: LED</Text>
			 <Text>SMART TV: FRIDGE / COMP BOTTLED WATER:</Text>
			 <Text>HEAVENLY</Text>
			 <Text>AAA OR CAA MEMBERSHIP REQUIRED AT BE SHOWN AT</Text>
			 <Text>CHECK-IN.</Text>
			</AdditionalInfo>
			<Rates>
			 <Rate Amount="274.55" ChangeIndicator="false" CurrencyCode="USD" HRD_RequiredForSell="false" PackageIndicator="false" RateConversionInd="false" ReturnOfRateInd="false" RoomOnRequest="false">
			  <AdditionalGuestAmounts>
			   <AdditionalGuestAmount MaxExtraPersonsAllowed="1" NumAdults="1" NumCribs="1">
				<Charges AdultRollAway="0.00" Crib="0.00" ExtraPerson="20.00"/>
			   </AdditionalGuestAmount>
			  </AdditionalGuestAmounts>
			  <HotelTotalPricing Amount="307.50">
			   <Disclaimer>INCLUDES TAXES AND SURCHARGES</Disclaimer>
			   <TotalTaxes Amount="32.95">
				<TaxFieldOne>19.22</TaxFieldOne>
				<TaxFieldTwo>13.73</TaxFieldTwo>
				<Text>CITY TAX</Text>
				<Text>OCCUPANCY TAX</Text>
			   </TotalTaxes>
			  </HotelTotalPricing>
			 </Rate>
			</Rates>
		   </RoomRate>
		  </RoomRates>
		  <TimeSpan Duration="0005" End="2018-05-16" Start="2018-05-15"/>
		 </RoomStay>
		</HotelRateDescriptionRS></soap-env:Body></soap-env:Envelope>`)
)
