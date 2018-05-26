package itin

import (
	"encoding/xml"
	"testing"
)

var (
	samplefrom        = "z.com"
	samplepcc         = "ABCD1"
	samplebinsectoken = string([]byte(`Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!`))
	sampleconvid      = "fds8789h|dev@z.com"
	samplemid         = "mid:20180216-07:18:42.3|14oUa"
	sampletime        = "2018-05-25T19:29:20Z"
	sampleFirstName   = "Charles"
	sampleLastName    = "Babbage"
	samplePhone       = "123-456-7890"
	samplePsngrReq    = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">PassengerDetailsRQ</eb:Service><eb:Action>PassengerDetailsRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-05-25T19:29:20Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><PassengerDetailsRQ xmlns="http://services.sabre.com/sp/pd/v3_3" version="3.3.0" IgnoreOnError="false" HaltOnError="false"><PostProcessing IgnoreAfter="false" RedisplayReservation="true" UnmaskCreditCard="false"></PostProcessing><PreProcessing IgnoreBefore="true"></PreProcessing><TravelItineraryAddInfoRQ><CustomerInfo><ContactNumbers><ContactNumber NameNumber="1.1" Phone="123-456-7890" PhoneUseType="H"></ContactNumber></ContactNumbers><PersonName NameNumber="1.1" NameReference="ABC123" PassengerType="ADT"><GivenName>Charles</GivenName><Surname>Babbage</Surname></PersonName></CustomerInfo></TravelItineraryAddInfoRQ></PassengerDetailsRQ></soap-env:Body></soap-env:Envelope>`)
)

func TestSetPsngr(t *testing.T) {
	s := SetPsngrDetailsRequestStruct(samplePhone, sampleFirstName, sampleLastName)
	s.AddSpecialDetails()
	s.AddUniqueID("1234ABCD")

	if s.PassengerDetailsRQ.PreProcess.UniqueID.ID != "1234ABCD" {
		t.Errorf("s.PassengerDetailsRQ.PreProcess.UniqueID.ID given %v, built %v", "1234ABCD", s.PassengerDetailsRQ.PreProcess.UniqueID.ID)
	}

	spd := &SpecialReqDetails{}
	if s.PassengerDetailsRQ.SpecialReq.SpecialServiceRQ.SpecialServiceInfo.AdvancedPassenger.VendorPrefs.Airline.Hosted != spd.SpecialServiceRQ.SpecialServiceInfo.AdvancedPassenger.VendorPrefs.Airline.Hosted {
		t.Errorf("AddSpecialDetails \ngiven: %v \nbuilt: %v", spd, s.PassengerDetailsRQ.SpecialReq)
	}

	pn := s.PassengerDetailsRQ.TravelItinInfo.Customer.PersonName
	if pn.First.Val != sampleFirstName {
		t.Errorf("TravelItinInfo.Customer.PersonName.First expect: %s, got %s", sampleFirstName, pn.First.Val)
	}
	if pn.Last.Val != sampleLastName {
		t.Errorf("TravelItinInfo.Customer.PersonName.Last expect: %s, got %s", sampleLastName, pn.Last.Val)
	}
}

func TestBuildPsngr(t *testing.T) {
	body := SetPsngrDetailsRequestStruct(samplePhone, sampleFirstName, sampleLastName)
	req := BuildPsngrDetailsRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, body)

	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling passenger details request", err)
	}
	if string(b) != string(samplePsngrReq) {
		t.Errorf("Expected marshal passenger details request \n given: %s \n built: %s", string(samplePsngrReq), string(b))
	}
}
