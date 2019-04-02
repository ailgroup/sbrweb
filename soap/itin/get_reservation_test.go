package itin

import (
	"encoding/xml"
	"testing"
)

var (
	sampleLocator          = "IJKZUQ"
	sampleGetReservationRQ = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">GetReservationRQ</eb:Service><eb:Action>GetReservationRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-05-25T19:29:20Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><GetReservationRQ xmlns="http://webservices.sabre.com/pnrbuilder/v1_19" Version="1.19.0"><Locator>IJKZUQ</Locator><RequestType>Stateful</RequestType></GetReservationRQ></soap-env:Body></soap-env:Envelope>`)
)

func TestBuildGetReservationMarshal(t *testing.T) {
	req := BuildGetReservationRequest(sampleConf, samplebinsectoken, sampleLocator)
	_, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling pnr read request", err)
	}
}
