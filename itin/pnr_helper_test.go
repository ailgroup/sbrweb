package itin

import (
	"net/http"
	"net/http/httptest"
)

var (
	//serverDown unavailable service
	serverDown = &httptest.Server{}
	//serverBadBody mocks a server that returns malformed body
	serverBadBody = &httptest.Server{}
	//serverPNRDetails responds with successful and valid post for passenger details
	serverPNRDetails = &httptest.Server{}
	//serverBizLogic mocks a server that returns warning business logic repsonse (due to missing first name)
	serverBizLogic = &httptest.Server{}
)

//Initialize Mock Sabre Web Servers and test data
func init() {
	serverDown = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write([]byte(`hello`))
			},
		),
	)
	defer func() { serverDown.Close() }()

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

	// init test servers...
	serverPNRDetails = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(samplePNRRes)
			},
		),
	)
	//defer func() { serverPNRDetails.Close() }()

	serverBizLogic = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(samplePNRResWarnBizLogic)
			},
		),
	)
	//defer func() { serverBizLogic.Close() }()
}

var (
	samplefrom        = "z.com"
	samplepcc         = "ABCD1"
	samplebinsectoken = string([]byte(`Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!`))
	sampleconvid      = "fds8789h|dev@z.com"
	samplemid         = "mid:20180216-07:18:42.3|14oUa"
	sampletime        = "2018-05-25T19:29:20Z"
	sampletimeOffset  = "2018-05-25T20:29:21.213-05:00"
	sampleFirstName   = "Charles"
	sampleLastName    = "Babbage"
	samplePhoneRes    = "123-456-7890-H.1.1"
	samplePhoneReq    = "123-456-7890"
)

var (
	samplePNRReq             = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">PassengerDetailsRQ</eb:Service><eb:Action>PassengerDetailsRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-05-25T19:29:20Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><PassengerDetailsRQ xmlns="http://services.sabre.com/sp/pd/v3_3" version="3.3.0" IgnoreOnError="false" HaltOnError="false"><PostProcessing IgnoreAfter="false" RedisplayReservation="true" UnmaskCreditCard="false"></PostProcessing><PreProcessing IgnoreBefore="true"></PreProcessing><TravelItineraryAddInfoRQ><CustomerInfo><ContactNumbers><ContactNumber NameNumber="1.1" Phone="123-456-7890" PhoneUseType="H"></ContactNumber></ContactNumbers><PersonName NameNumber="1.1" NameReference="ABC123" PassengerType="ADT"><GivenName>Charles</GivenName><Surname>Babbage</Surname></PersonName></CustomerInfo></TravelItineraryAddInfoRQ></PassengerDetailsRQ></soap-env:Body></soap-env:Envelope>`)
	samplePNRResWarnBizLogic = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">z.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">PassengerDetailsRQ</eb:Service><eb:Action>PassengerDetailsRS</eb:Action><eb:MessageData><eb:MessageId>1bb067aj2</eb:MessageId><eb:Timestamp>2018-05-26T02:12:47</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><PassengerDetailsRS xmlns="http://services.sabre.com/sp/pd/v3_3"><ApplicationResults xmlns="http://services.sabre.com/STL_Payload/v02_01" status="Complete"><Success timeStamp="2018-05-26T14:26:47.962-05:00"/><Warning type="BusinessLogic" timeStamp="2018-05-26T14:26:47.873-05:00"><SystemSpecificResults><Message code="WARN.SWS.HOST.ERROR_IN_RESPONSE">.CQT.NBR.FIRST NAMES.NOT ENT BGNG WITH</Message></SystemSpecificResults></Warning><Warning type="Application" timeStamp="2018-05-26T14:26:47.960-05:00"><SystemSpecificResults><Message code="WARN.SP.PROVIDER_ERROR">No PNR in AAA, caused by [No PNR in AAA, code: 500306, severity: WARNING]</Message><Message code="700408">No PNR in AAA, caused by [No PNR in AAA, code: 500306, severity: WARNING]</Message></SystemSpecificResults></Warning></ApplicationResults></PassengerDetailsRS></soap-env:Body></soap-env:Envelope>`)
	samplePNRRes = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">z.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">PassengerDetailsRQ</eb:Service><eb:Action>PassengerDetailsRS</eb:Action><eb:MessageData><eb:MessageId>1auug50cv</eb:MessageId><eb:Timestamp>2018-05-26T02:12:47</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><PassengerDetailsRS xmlns="http://services.sabre.com/sp/pd/v3_3"><ApplicationResults xmlns="http://services.sabre.com/STL_Payload/v02_01" status="Complete"><Success timeStamp="2018-05-25T20:29:21.213-05:00"/></ApplicationResults><TravelItineraryReadRS><TravelItinerary><CustomerInfo><ContactNumbers><ContactNumber LocationCode="SLC" Phone="123-456-7890-H.1.1" RPH="001"/></ContactNumbers><PersonName WithInfant="false" NameNumber="01.01" NameReference="ABC123" PassengerType="ADT" RPH="1"><GivenName>CHARLES</GivenName><Surname>BABBAGE</Surname></PersonName></CustomerInfo><ItineraryInfo><ReservationItems/></ItineraryInfo><ItineraryRef AirExtras="false" InhibitCode="U" PartitionID="AA" PrimeHostID="1S"><Source PseudoCityCode="ABCD1"/></ItineraryRef></TravelItinerary></TravelItineraryReadRS></PassengerDetailsRS></soap-env:Body></soap-env:Envelope>`)
)
