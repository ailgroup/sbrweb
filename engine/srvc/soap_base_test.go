package srvc

/*
TESTING NOTES:
	- all data variables use for mocking tests are downcase and start with sample*
	- any table test strutcts come directly before their test (or first test using table)
	- functions used in a test have their tests come first
	- Benchmarks come after tests using testable functionality
	- Benchmarks have same name as test that is benchmarked (sans Test/Benchmark prefix)
*/

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

var (
	//serverDown mocks an unreachable service
	serverDown = &httptest.Server{}
	//serverBadBody mocks a server that returns malformed body
	serverBadBody = &httptest.Server{}
	//serverCreateRSUnauth for testing session create not authorized
	serverCreateRSUnauth = &httptest.Server{}
	//serverCloseRSInvalid for testing session create not authorized
	serverCloseRSInvalid = &httptest.Server{}
	//serverCreateRQ for testing session create
	serverCreateRQ = &httptest.Server{}
	//serverCloseRQ for testing session close
	serverCloseRQ = &httptest.Server{}
	//serverValidateRQ for testing session validate
	serverValidateRQ    = &httptest.Server{}
	samplerandStr       = regexp.MustCompile(`\w*`)
	samplerfc333pString = "2017-11-27T09:58:31Z"
	samplerfc333pReg    = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	samplemidReg        = regexp.MustCompile(`mid:\d{8}-\d{2}:\d{2}:\d{2}\.\d{1,2}\|\w{5}`)
	sampleconvReg       = regexp.MustCompile(`cid:\w{8}\|.{5,}`)
	samplemid           = "mid:20180216-07:18:42.3|14oUa"
	sampletime          = "2018-02-16T07:18:42Z"
	samplepcc           = "7TZA"
	sampleorg           = "7TZA"
	sampleconvid        = "fds8789h|dev@z.com"
	samplefrom          = "www.z.com"
	sampleservice       = "SessionCreateRQ"
	sampleaction        = "SessionCreateRQ"
	sampleusername      = "773400"
	samplepassword      = "PASSWORD_GOES_HER"
	sampledomain        = "DEFAULT"
	sampleSessionConf   = &SessionConf{
		From:      samplefrom,
		PCC:       samplepcc,
		Convid:    sampleconvid,
		Msgid:     samplemid,
		Timestr:   sampletime,
		Username:  sampleusername,
		Password:  samplepassword,
		Binsectok: samplebinsectoken,
	}
	samplebinsectoken              = string([]byte(`Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0`))
	samplebintokensplit            = "-3177016070087638144!110012!0"
	sampleSessionNoAuthFaultCode   = "soap-env:Client.AuthenticationFailed"
	sampleSessionNoAuthFaultString = " Authentication failed "
	sampleSessionNoAuthStackTrace  = "com.sabre.universalservices.base.security.AuthenticationException: errors.authentication.USG_AUTHENTICATION_FAILED"

	sampleSessionInvalidTokenFaultCode  = "soap-env:Client.InvalidSecurityToken"
	sampleSessionInvalidTokenStackTrace = "com.sabre.universalservices.base.session.SessionException: errors.session.USG_INVALID_SECURITY_TOKEN"

	sampleEnvelope = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"></soap-env:Envelope>`)

	sampleMessageHeader = []byte(`<eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCreateRQ</eb:Service><eb:Action>SessionCreateRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader>`)

	sampleSecurityRequest = []byte(`<wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:UsernameToken><wsse:Username>773400</wsse:Username><wsse:Password>PASSWORD_GOES_HER</wsse:Password><Organization>7TZA</Organization><Domain>DEFAULT</Domain></wsse:UsernameToken></wsse:Security>`)

	sampleSessionCreateRQ = []byte(`<ns:SessionCreateRQ xmlns:ns="http://www.opentravel.org/OTA/2002/11"><POS><Source PseudoCityCode="7TZA"></Source></POS></ns:SessionCreateRQ>`)

	sampleManifest = []byte(`<eb:Manifest soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:Reference xlink:href="cid:rootelement" xlink:type="simple"></eb:Reference></eb:Manifest>`)

	sampleSecurityResponse = []byte(`<wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security>`)

	sampleSessionRQHeader = []byte(`<soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>www.z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCreateRQ</eb:Service><eb:Action>SessionCreateRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:UsernameToken><wsse:Username>773400</wsse:Username><wsse:Password>PASSWORD_GOES_HER</wsse:Password><Organization>7TZA</Organization><Domain>DEFAULT</Domain></wsse:UsernameToken></wsse:Security></soap-env:Header>`)

	sampleSessionCreateRQBody = []byte(`<soap-env:Body><eb:Manifest soap-env:mustUnderstand="" eb:version=""><eb:Reference xlink:href="" xlink:type=""></eb:Reference></eb:Manifest><ns:SessionCreateRQ xmlns:ns=""><POS><Source PseudoCityCode=""></Source></POS></ns:SessionCreateRQ></soap-env:Body>`)

	sampleSessionEnvelope = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>www.z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCreateRQ</eb:Service><eb:Action>SessionCreateRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:UsernameToken><wsse:Username>773400</wsse:Username><wsse:Password>PASSWORD_GOES_HER</wsse:Password><Organization>7TZA</Organization><Domain>DEFAULT</Domain></wsse:UsernameToken></wsse:Security></soap-env:Header><soap-env:Body><eb:Manifest soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:Reference xlink:href="cid:rootelement" xlink:type="simple"></eb:Reference></eb:Manifest><ns:SessionCreateRQ xmlns:ns="http://www.opentravel.org/OTA/2002/11"><POS><Source PseudoCityCode="7TZA"></Source></POS></ns:SessionCreateRQ></soap-env:Body></soap-env:Envelope>`)

	sampleSessionEnvelopeWithValues = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCreateRQ</eb:Service><eb:Action>SessionCreateRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:UsernameToken><wsse:Username>773400</wsse:Username><wsse:Password>PASSWORD_GOES_HER</wsse:Password><Organization>7TZA</Organization><Domain>DEFAULT</Domain></wsse:UsernameToken></wsse:Security></soap-env:Header><soap-env:Body><eb:Manifest soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:Reference xlink:href="cid:rootelement" xlink:type="simple"></eb:Reference></eb:Manifest><ns:SessionCreateRQ xmlns:ns="http://www.opentravel.org/OTA/2002/11"><POS><Source PseudoCityCode="7TZA"></Source></POS></ns:SessionCreateRQ></soap-env:Body></soap-env:Envelope>`)

	sampleSessionSuccessResponse = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">Session</eb:Service><eb:Action>SessionCreateRS</eb:Action><eb:MessageData><eb:MessageId>4379957601383660213</eb:MessageId><eb:Timestamp>2018-02-18T16:42:18</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><SessionCreateRS xmlns="http://www.opentravel.org/OTA/2002/11" version="1" status="Approved">	<ConversationId>fds8789h|dev@z.com</ConversationId></SessionCreateRS></soap-env:Body></soap-env:Envelope>`)

	sampleSessionUnAuth = []byte(`<?xml version="1.0" encoding="UTF-8"?><soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service>Session</eb:Service><eb:Action>ErrorRS</eb:Action><eb:MessageData><eb:MessageId>4341295557539920551</eb:MessageId><eb:Timestamp>2018-02-18T15:29:14</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" /></soap-env:Header><soap-env:Body><soap-env:Fault><faultcode>soap-env:Client.AuthenticationFailed</faultcode><faultstring>Authentication failed</faultstring><detail><StackTrace>com.sabre.universalservices.base.security.AuthenticationException: errors.authentication.USG_AUTHENTICATION_FAILED</StackTrace></detail></soap-env:Fault></soap-env:Body></soap-env:Envelope>`)

	sampleSessionPoolMsgNoAuth = string([]byte(`Authentication failed-soap-env:Client.AuthenticationFailed: com.sabre.universalservices.base.security.AuthenticationException: errors.authentication.USG_AUTHENTICATION_FAILED`))

	sampleSessionCloseRQ = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCloseRQ</eb:Service><eb:Action>SessionCloseRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><SessionCloseRQ><POS><Source PseudoCityCode="7TZA"></Source></POS></SessionCloseRQ></soap-env:Body></soap-env:Envelope>`)

	sampleSessionCloseRespSuccess = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">Session</eb:Service><eb:Action>SessionCloseRS</eb:Action><eb:MessageData><eb:MessageId>5777927785778750193</eb:MessageId><eb:Timestamp>2018-02-18T21:49:37</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESC!ICESMSLB\/RES.LB!-3176941104583370867!441490!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><SessionCloseRS xmlns="http://www.opentravel.org/OTA/2002/11" version="1" status="Approved"/></soap-env:Body></soap-env:Envelope>`)

	sampleSessionInvalidTokenString = string([]byte(`Invalid or Expired binary security token: Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0`))

	sampleSessionCloseRespNoValidToken = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionCloseRQ</eb:Service><eb:Action>ErrorRS</eb:Action><eb:MessageData><eb:MessageId>6238327787953800281</eb:MessageId><eb:Timestamp>2018-02-18T21:53:15</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><soap-env:Fault><faultcode>soap-env:Client.InvalidSecurityToken</faultcode><faultstring>Invalid or Expired binary security token: Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</faultstring><detail><StackTrace>com.sabre.universalservices.base.session.SessionException: errors.session.USG_INVALID_SECURITY_TOKEN</StackTrace></detail></soap-env:Fault></soap-env:Body></soap-env:Envelope>`)

	sampleSessionValidateRQ = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">www.z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionValidateRQ</eb:Service><eb:Action>SessionValidateRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><SessionValidateRQ><POS><Source PseudoCityCode="7TZA"></Source></POS></SessionValidateRQ></soap-env:Body></soap-env:Envelope>`)

	sampleSessionValidateRespSuccess = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">Session</eb:Service><eb:Action>SessionValidateRS</eb:Action><eb:MessageData><eb:MessageId>931442098348670191</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken><wsse:UsernameToken><wsse:Username>7971</wsse:Username><Organization>7TZA</Organization><Domain>DEFAULT</Domain></wsse:UsernameToken></wsse:Security></soap-env:Header><soap-env:Body><SessionValidateRS/></soap-env:Body></soap-env:Envelope>`)

	sampleSessionValidateRSInvalidTokenRS = []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"><soap-env:Header><eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="2.0.0" soap-env:mustUnderstand="1"><eb:From><eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId></eb:From><eb:To><eb:PartyId eb:type="URI">www.z.com</eb:PartyId></eb:To><eb:CPAId>7TZA</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="OTA">SessionValidateRQ</eb:Service><eb:Action>ErrorRS</eb:Action><eb:MessageData><eb:MessageId>4215837468837900553</eb:MessageId><eb:Timestamp>2018-02-16T07:18:42Z</eb:Timestamp><eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"><wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><soap-env:Fault><faultcode>soap-env:Client.InvalidSecurityToken</faultcode><faultstring>Invalid or Expired binary security token: Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3177016070087638144!110012!0</faultstring><detail><StackTrace>com.sabre.universalservices.base.session.SessionException: errors.session.USG_INVALID_SECURITY_TOKEN</StackTrace></detail></soap-env:Fault></soap-env:Body></soap-env:Envelope>`)
)

//Initialize Mock Sabre Web Servers
func init() {
	serverDown = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				//rs.WriteHeader(500)
			},
		))
	serverDown.Close()

	serverBadBody = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				//rs.Header()
				//rs.WriteHeader(500)
				//rs.Write(sampleBadBody)
				rs.Write([]byte(`!#BAD_/_BODY_.*__\\fhji(*&^%^%$%^&Y*(J)OPKL:`))
			},
		),
	)
	//defer func() { serverBadBody.Close() }()

	serverCreateRQ = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleSessionSuccessResponse)
			},
		),
	)
	//defer func() { serverCreateRQ.Close() }()

	serverCreateRSUnauth = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleSessionUnAuth)
			},
		),
	)
	//defer func() { serverCreateRSInvalid.Close() }()

	serverCloseRQ = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleSessionCloseRespSuccess)
			},
		),
	)
	//defer func() { serverCloseRQ.Close() }()

	serverCloseRSInvalid = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleSessionCloseRespNoValidToken)
			},
		),
	)
	//defer func() { serverCloseRSInvalid.Close() }()

	serverValidateRQ = httptest.NewServer(
		http.HandlerFunc(
			func(rs http.ResponseWriter, rq *http.Request) {
				rs.Write(sampleSessionValidateRespSuccess)
			},
		),
	)
	//defer func() { serverValidateRQ.Close() }()

}

func TestLogSetup(t *testing.T) {
	setUpLogging()
	if _, err := os.Stat("sabre_web_soap.log"); os.IsNotExist(err) {
		t.Fatal("no log file for soap logging")
	}
	if _, err := os.Stat("sabre_web_session.log"); os.IsNotExist(err) {
		t.Fatal("no log file for session logging")
	}
}

func TestEnvelopeBaseMarshal(t *testing.T) {
	envelope := Envelope{
		XMLNSbase:  BaseNS,
		XMLNSeb:    BaseEBNameSpace,
		XMLNSxlink: baseXlinkNameSpace,
		XMLNSxsd:   BaseXSDNameSpace,
	}
	envelope2 := CreateEnvelope()

	for _, env := range []Envelope{envelope, envelope2} {
		//b, err := xml.MarshalIndent(env, "", "    ")
		b, err := xml.Marshal(env)
		if err != nil {
			t.Error("Error marshaling envelope", err)
		}
		if string(b) != string(sampleEnvelope) {
			t.Errorf("Expected marshal envelope \n sample: %s \n result: %s", string(sampleEnvelope), string(b))
		}
	}
}
func BenchmarkCreateEnvelope(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CreateEnvelope()
	}
}
func BenchmarkEnvelopeMarshal(b *testing.B) {
	envelope := CreateEnvelope()
	for n := 0; n < b.N; n++ {
		xml.Marshal(envelope)
	}
}

func TestEnvelopeBaseUnmarshal(t *testing.T) {
	env := EnvelopeUnMarsh{}
	err := xml.Unmarshal(sampleEnvelope, &env)
	if err != nil {
		t.Errorf("Error unmarshaling sample envelope %s \nERROR: %v", sampleEnvelope, err)
	}

	if env.XMLName.Local != "Envelope" {
		t.Errorf("Envelope xml Local wrong: expected: %s, got: %s", "Envelope", env.XMLName.Local)
	}
	if env.XMLName.Space != BaseNS {
		t.Errorf("Envelope xml Space wrong: expected: %s, got: %s", BaseNS, env.XMLName.Space)
	}
	if env.XMLNSbase != BaseNS {
		t.Errorf("Envelope XMLNSbase expected: %s, got: %s", BaseNS, env.XMLNSbase)
	}
	//fmt.Printf("SAMPLE: %s\n", sampleEnvelope)
	//fmt.Printf("CURRENT: %+v\n", env)
}
func BenchmarkEnvelopeUnmarshal(b *testing.B) {
	env := EnvelopeUnMarsh{}
	for n := 0; n < b.N; n++ {
		xml.Unmarshal(sampleEnvelope, &env)
	}
}

func TestTimeFormat(t *testing.T) {
	format := SabreTimeFormat()
	if !samplerfc333pReg.MatchString(format) {
		t.Errorf("Timestamp formt needs to be '%s' but got '%s'", samplerfc333pString, format)
	}
}
func BenchmarkTimeFormat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SabreTimeFormat()
	}
}

func TestSessionConfSetTime(t *testing.T) {
	conf := &SessionConf{
		Timestr: sampletime,
	}
	if conf.SetTime().Timestr != SabreTimeFormat() {
		t.Error("SessionConf SetTime() should be SabreTimeFormat()")
	}
}

var randstrtests = []struct {
	size int
}{
	{0},
	{1},
	{2},
	{3},
	{4},
	{10},
	{100},
}

func TestRandomStringGen(t *testing.T) {
	for _, randstr := range randstrtests {
		rstr1 := randStringBytesMaskImprSrc(randstr.size)
		if !samplerandStr.MatchString(rstr1) {
			t.Errorf("Random string does not match alpha-numeric pattern: %s", rstr1)
		}
		if len(rstr1) != randstr.size {
			t.Errorf("Random string should be %d characters long: got %d %s", randstr.size, len(rstr1), rstr1)
		}

	}

}
func BenchmarkRandomStringGen(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randStringBytesMaskImprSrc(20)
	}
}

func TestGenerateMessageID(t *testing.T) {
	mid := GenerateMessageID()
	if !samplemidReg.MatchString(mid) {
		t.Errorf("MessageID format wrong. Example: '%s', got '%s'", samplemid, mid)
	}
}
func BenchmarkGenerateMessageID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateMessageID()
	}
}

func TestGenerateConversationID(t *testing.T) {
	conv := GenerateConversationID(samplefrom)
	if !sampleconvReg.MatchString(conv) {
		t.Errorf("ConversaionID format wrong. Example: '%s', got '%s'", sampleconvid, conv)
	}
}
func BenchmarkGenerateConversationID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateConversationID(samplefrom)
	}
}

func TestSabreTokenParse(t *testing.T) {
	tok := SabreTokenParse(samplebinsectoken)
	if tok != samplebintokensplit {
		t.Errorf("BinaryTokenSplit epxect '%s', got '%s'", samplebintokensplit, tok)
	}
}
func BenchmarkSabreTokenParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SabreTokenParse(samplebinsectoken)
	}
}

func TestMessageHeaderBaseMarshal(t *testing.T) {
	mh := MessageHeader{
		MustUnderstand: "1",
		EbVersion:      SabreEBVersion,
		From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
		To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
		CPAID:          samplepcc,
		ConversationID: sampleconvid,
		Service:        ServiceElem{sampleservice, "OTA"},
		Action:         sampleaction,
		MessageData: MessageDataElem{
			MessageID: samplemid,  //not generating this so easy to test
			Timestamp: sampletime, //not genereating this so easy to test
		},
	}
	mh2 := MessageHeader{
		MustUnderstand: "1",
		EbVersion:      SabreEBVersion,
		From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
		To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
		CPAID:          samplepcc,
		ConversationID: sampleconvid,
		Service:        ServiceElem{sampleservice, "OTA"},
		Action:         sampleaction,
		MessageData: MessageDataElem{
			MessageID: samplemid,  //not generating this so easy to test
			Timestamp: sampletime, //not genereating this so easy to test
		},
	}

	for _, header := range []MessageHeader{mh, mh2} {
		//b, err := xml.MarshalIndent(header, "", "    ") //don't indent it fails
		b, err := xml.Marshal(header)
		if err != nil {
			t.Error("Error marshaling message header", err)
		}
		if string(b) != string(sampleMessageHeader) {
			t.Errorf("Expected marshal message header \n sample: %s \n result: %s", string(sampleMessageHeader), string(b))
		}
	}
}
func BenchmarkMessageHeaderBaseMarshal(b *testing.B) {
	mh := MessageHeader{
		MustUnderstand: "1",
		EbVersion:      SabreEBVersion,
		From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
		To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
		CPAID:          samplepcc,
		ConversationID: sampleconvid,
		Service:        ServiceElem{sampleservice, "OTA"},
		Action:         sampleaction,
		MessageData: MessageDataElem{
			MessageID: samplemid,  //not generating this so easy to test
			Timestamp: sampletime, //not genereating this so easy to test
		},
	}
	for n := 0; n < b.N; n++ {
		xml.Marshal(mh)
	}
}

func TestMessageHeaderBaseUnmarshal(t *testing.T) {
	mh := MessageHeaderUnmarsh{}
	err := xml.Unmarshal(sampleMessageHeader, &mh)
	if err != nil {
		t.Errorf("Error unmarshaling sample envelope %s \nERROR: %v", sampleMessageHeader, err)
	}
	if mh.XMLName.Local != "MessageHeader" {
		t.Errorf("MessageHeader xml Local wrong: expected: %s, got: %s", "MessageHeader", mh.XMLName.Local)
	}
	if mh.XMLName.Space != "eb" {
		t.Errorf("MessageHeader xml Space wrong: expected: %s, got: %s", "eb", mh.XMLName.Space)
	}
	if mh.MustUnderstand != "1" {
		t.Error("MustUnderstand shoudl be 1")
	}
	if mh.EbVersion != SabreEBVersion {
		t.Error("EbVersion should be ", SabreEBVersion)
	}
	if mh.From.XMLName.Space != "eb" {
		t.Errorf("From xml Space wrong: expected: %s, got: %s", "eb", mh.From.XMLName.Space)
	}
	if mh.To.XMLName.Space != "eb" {
		t.Errorf("To xml Space wrong: expected: %s, got: %s", "eb", mh.To.XMLName.Space)
	}
	//fmt.Printf("SAMPLE: %s\n", sampleMessageHeader)
	//fmt.Printf("CURRENT: %+v\n", mh)
}
func BenchmarkMessageHeaderBaseUnmarshal(b *testing.B) {
	mh := MessageHeaderUnmarsh{}
	for n := 0; n < b.N; n++ {
		xml.Unmarshal(sampleMessageHeader, &mh)
	}
}

func TestSecurityBaseMarshal(t *testing.T) {
	sec := Security{
		XMLNSWsseBase: BaseWsse,
		XMLNSWsu:      BaseWsuNameSpace,
		UserNameToken: &UsernameTokenElem{
			Username:     sampleusername,
			Password:     samplepassword,
			Organization: sampleorg,
			Domain:       sampledomain,
		},
	}
	sec2 := Security{
		XMLNSWsseBase: BaseWsse,
		XMLNSWsu:      BaseWsuNameSpace,
		UserNameToken: &UsernameTokenElem{
			Username:     sampleusername,
			Password:     samplepassword,
			Organization: sampleorg,
			Domain:       sampledomain,
		},
	}
	for _, s := range []Security{sec, sec2} {
		//b, err := xml.MarshalIndent(s, "", "    ")
		b, err := xml.Marshal(s)
		if err != nil {
			t.Error("Error marshaling security", err)
		}
		if string(b) != string(sampleSecurityRequest) {
			t.Errorf("Expected marshal security \n sample: %s \n result: %s", string(sampleSecurityRequest), string(b))
		}
	}
}
func BenchmarkSecurityBaseMarshal(b *testing.B) {
	sec := Security{
		XMLNSWsseBase: BaseWsse,
		XMLNSWsu:      BaseWsuNameSpace,
		UserNameToken: &UsernameTokenElem{
			Username:     sampleusername,
			Password:     samplepassword,
			Organization: sampleorg,
			Domain:       sampledomain,
		},
	}
	for n := 0; n < b.N; n++ {
		xml.Marshal(sec)
	}
}

func TestSecurityBaseUnmarshal(t *testing.T) {
	sec := SecurityUnmarsh{}
	err := xml.Unmarshal(sampleSecurityResponse, &sec)
	if err != nil {
		t.Errorf("Error unmarshaling sample security response %s \nERROR: %v", sampleSecurityResponse, err)
	}

	if sec.XMLName.Local != "Security" {
		t.Errorf("Security xml Local wrong: expected: %s, got: %s", "Security", sec.XMLName.Local)
	}
	if sec.XMLName.Space != BaseWsse {
		t.Errorf("Security xml Space wrong: expected: %s, got: %s", BaseWsse, sec.XMLName.Space)
	}
	if sec.XMLNSWsseBase != BaseWsse {
		t.Errorf("XMLNSWsseBase expected: %s, got: %s", BaseWsse, sec.XMLNSWsseBase)
	}
	if sec.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("BinarySecurityToken.Value expected: %s, got: %s", samplebinsectoken, sec.BinarySecurityToken)
	}
	if sec.BinarySecurityToken.ValueType != "String" {
		t.Errorf("BinarySecurityToken.ValueType expected: %s, got: %s", "String", sec.BinarySecurityToken.ValueType)
	}
	if sec.BinarySecurityToken.EncodingType != "wsse:Base64Binary" {
		t.Errorf("BinarySecurityToken.EncodingType expected: %s, got: %s", "wsse:Base64Binary", sec.BinarySecurityToken.EncodingType)
	}
	//fmt.Printf("SAMPLE: %s\n", sampleSecurityResponse)
	//fmt.Printf("CURRENT: %+v\n", sec)
}
func BenchmarkSecurityBaseUnmarshal(b *testing.B) {
	s := SecurityUnmarsh{}
	for n := 0; n < b.N; n++ {
		xml.Unmarshal(sampleSecurityResponse, &s)
	}
}

func TestManifestMarshal(t *testing.T) {
	mnf := Manifest{
		MustUnderstand: "1",
		EbVersion:      SabreEBVersion,
		Reference: ReferenceElem{
			Href: "cid:rootelement",
			Type: "simple",
		},
	}

	mnf2 := CreateManifest()

	for _, m := range []Manifest{mnf, mnf2} {
		b, err := xml.Marshal(m)
		if err != nil {
			t.Error("Error marshaling manifest", err)
		}
		if string(b) != string(sampleManifest) {
			t.Errorf("Expected marshal manifest \n sample: %s \n result: %s", string(sampleManifest), string(b))
		}
	}
	//fmt.Printf("SAMPLE: %s\n", sampleManifest)
	//fmt.Printf("CURRENT: %+v\n", mnf)
}
func BenchmarkManifestMarshal(b *testing.B) {
	mnf := CreateManifest()
	for n := 0; n < b.N; n++ {
		xml.Marshal(mnf)
	}
}

func TestSessionCreateRQMarshal(t *testing.T) {
	sess := SessionCreateRQ{
		XMLNS: baseOTANameSpace,
		POS: POSElem{
			Source: SourceElem{
				PseudoCityCode: sampleorg,
			},
		},
	}
	b, err := xml.Marshal(sess)
	if err != nil {
		t.Error("Error marshaling session create rq", err)
	}
	if string(b) != string(sampleSessionCreateRQ) {
		t.Errorf("Expected marshal session create rq \n sample: %s \n result: %s", string(sampleSessionCreateRQ), string(b))
	}

	//fmt.Printf("SAMPLE: %s\n", sampleSessionCreateRQ)
	//fmt.Printf("CURRENT: %+v\n", create)
}
func BenchmarkSessionCreateRQMarshal(b *testing.B) {
	s := SessionCreateRQ{
		XMLNS: baseOTANameSpace,
		POS: POSElem{
			Source: SourceElem{
				PseudoCityCode: sampleorg,
			},
		},
	}
	for n := 0; n < b.N; n++ {
		xml.Marshal(s)
	}
}

func TestSessionHeaderMarshal(t *testing.T) {
	shdr := SessionHeader{
		MessageHeader: MessageHeader{
			MustUnderstand: SabreMustUnderstand,
			EbVersion:      SabreEBVersion,
			From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
			To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
			CPAID:          samplepcc,
			ConversationID: samplefrom,
			Service:        ServiceElem{"SessionCreateRQ", "OTA"},
			Action:         "SessionCreateRQ",
			MessageData: MessageDataElem{
				MessageID: samplemid,
				Timestamp: sampletime,
			},
		},
		Security: Security{
			XMLNSWsseBase: BaseWsse,
			XMLNSWsu:      BaseWsuNameSpace,
			UserNameToken: &UsernameTokenElem{
				Username:     sampleusername,
				Password:     samplepassword,
				Organization: samplepcc,
				Domain:       sabreDefaultDomain,
			},
		},
	}

	b, err := xml.Marshal(&shdr)
	if err != nil {
		t.Error("Error marshalling session soap header", err)
	}
	if string(b) != string(sampleSessionRQHeader) {
		t.Errorf("Expected marshal session soap header \n sample: %s \n result: %s", string(sampleSessionRQHeader), string(b))
	}
}
func BenchmarkSessionHeaderMarshal(b *testing.B) {
	shdr := SessionHeader{
		MessageHeader: MessageHeader{
			MustUnderstand: SabreMustUnderstand,
			EbVersion:      SabreEBVersion,
			From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
			To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
			CPAID:          samplepcc,
			ConversationID: samplefrom,
			Service:        ServiceElem{"SessionCreateRQ", "OTA"},
			Action:         "SessionCreateRQ",
			MessageData: MessageDataElem{
				MessageID: samplemid,
				Timestamp: sampletime,
			},
		},
		Security: Security{
			XMLNSWsseBase: BaseWsse,
			XMLNSWsu:      BaseWsuNameSpace,
			UserNameToken: &UsernameTokenElem{
				Username:     sampleusername,
				Password:     samplepassword,
				Organization: samplepcc,
				Domain:       sabreDefaultDomain,
			},
		},
	}
	for n := 0; n < b.N; n++ {
		xml.Marshal(shdr)
	}
}
func BenchmarkSessionHeaderUnmarshal(b *testing.B) {
	s := SessionHeader{}
	for n := 0; n < b.N; n++ {
		xml.Unmarshal(sampleSessionRQHeader, &s)
	}
}

func TestSessionCreateRQBodyMarshal(t *testing.T) {
	sbdy := SessionCreateRQBody{}

	b, err := xml.Marshal(&sbdy)
	if err != nil {
		t.Error("Error marshalling session body", err)
	}
	if string(b) != string(sampleSessionCreateRQBody) {
		t.Errorf("Expected marshal session body \n sample: %s \n result: %s", string(sampleSessionCreateRQBody), string(b))
	}
}

func TestSessionCreateRequest(t *testing.T) {
	sessionRequest := SessionCreateRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          samplepcc,
				ConversationID: samplefrom,
				Service:        ServiceElem{"SessionCreateRQ", "OTA"},
				Action:         "SessionCreateRQ",
				MessageData: MessageDataElem{
					MessageID: samplemid,
					Timestamp: sampletime,
				},
			},
			Security: Security{
				XMLNSWsseBase: BaseWsse,
				XMLNSWsu:      BaseWsuNameSpace,
				UserNameToken: &UsernameTokenElem{
					Username:     sampleusername,
					Password:     samplepassword,
					Organization: samplepcc,
					Domain:       sabreDefaultDomain,
				},
			},
		},
		Body: SessionCreateRQBody{
			Manifest: CreateManifest(),
			SessionCreateRQ: SessionCreateRQ{
				XMLNS: baseOTANameSpace,
				POS: POSElem{
					Source: SourceElem{
						PseudoCityCode: samplepcc,
					},
				},
			},
		},
	}
	b, err := xml.Marshal(&sessionRequest)
	if err != nil {
		t.Error("Error marshalling session body", err)
	}
	if string(b) != string(sampleSessionEnvelope) {
		t.Errorf("Expected marshal session body \n sample: %s \n result: %s", string(sampleSessionEnvelope), string(b))
	}
}
func BenchmarkSessionCreateRequestMarshal(b *testing.B) {
	sessionRequest := SessionCreateRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          samplepcc,
				ConversationID: samplefrom,
				Service:        ServiceElem{"SessionCreateRQ", "OTA"},
				Action:         "SessionCreateRQ",
				MessageData: MessageDataElem{
					MessageID: samplemid,
					Timestamp: sampletime,
				},
			},
			Security: Security{
				XMLNSWsseBase: BaseWsse,
				XMLNSWsu:      BaseWsuNameSpace,
				UserNameToken: &UsernameTokenElem{
					Username:     sampleusername,
					Password:     samplepassword,
					Organization: samplepcc,
					Domain:       sabreDefaultDomain,
				},
			},
		},
		Body: SessionCreateRQBody{
			Manifest: CreateManifest(),
			SessionCreateRQ: SessionCreateRQ{
				XMLNS: baseOTANameSpace,
				POS: POSElem{
					Source: SourceElem{
						PseudoCityCode: samplepcc,
					},
				},
			},
		},
	}
	for n := 0; n < b.N; n++ {
		xml.Marshal(sessionRequest)
	}
}
func BenchmarkSessionCreateRequestUnmarshal(b *testing.B) {
	s := SessionCreateRequest{}
	for n := 0; n < b.N; n++ {
		xml.Unmarshal(sampleSessionEnvelope, &s)
	}
}

func TestBuildSessionCreateRequest(t *testing.T) {
	sess := BuildSessionCreateRequest(sampleSessionConf)

	if sess.Header.MessageHeader.From.PartyID.Value != samplefrom {
		t.Errorf("Header.MessageHeader.From.PartyID.Value expect: %s, got %s", samplefrom, sess.Header.MessageHeader.From.PartyID.Value)
	}
	if sess.Header.MessageHeader.MessageData.MessageID != samplemid {
		t.Errorf("Header.MessageHeader.MessageData.MessageID expect %s, got %s", samplemid, sess.Header.MessageHeader.MessageData.MessageID)
	}
	if sess.Header.MessageHeader.MessageData.Timestamp != sampletime {
		t.Errorf("Header.MessageHeader.MessageData.MessageID expect %s, got %s", sampletime, sess.Header.MessageHeader.MessageData.Timestamp)
	}
	if sess.Body.SessionCreateRQ.POS.Source.PseudoCityCode != sampleorg {
		t.Errorf("Body.SessionCreateRQ.POS.Source.PseudoCityCode expect %s, got %s", sampleorg, sess.Body.SessionCreateRQ.POS.Source.PseudoCityCode)
	}
}
func BenchmarkBuildSessionCreateRequest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BuildSessionCreateRequest(sampleSessionConf)
	}
}
func TestBuildSessionCreateRequestMarshal(t *testing.T) {
	sess := BuildSessionCreateRequest(sampleSessionConf)
	b, err := xml.Marshal(&sess)
	if err != nil {
		t.Error("Error marshalling session envelope", err)
	}
	if string(b) != string(sampleSessionEnvelopeWithValues) {
		t.Error("Session envelope with values does not match test sample")
	}
}
func BenchmarkBuildSessionCreateRequestMarshal(b *testing.B) {
	s := BuildSessionCreateRequest(sampleSessionConf)
	for n := 0; n < b.N; n++ {
		xml.Marshal(&s)
	}
}

func TestSessionCreateResponse(t *testing.T) {
	resp := SessionCreateResponse{}
	err := xml.Unmarshal(sampleSessionSuccessResponse, &resp)

	if err != nil {
		t.Errorf("Error unmarshaling session create response %s \nERROR: %v", sampleSessionSuccessResponse, err)
	}

	if resp.Header.MessageHeader.To.PartyID.Value != samplefrom {
		t.Errorf("SessionRSHeader.MessageHeader.To.PartyID.Value expected: %s, got: %s", samplefrom, resp.Header.MessageHeader.To.PartyID.Value)
	}
	if resp.Header.MessageHeader.CPAID != samplepcc {
		t.Errorf("SessionRSHeader.MessageHeader.CPAID expect %s, got %s", samplepcc, resp.Header.MessageHeader.CPAID)
	}
	if resp.Header.MessageHeader.ConversationID != sampleconvid {
		t.Errorf("SessionRSHeader.MessageHeader.ConversationID expect %s, got %s", sampleconvid, resp.Header.MessageHeader.ConversationID)
	}
	if resp.Header.MessageHeader.MessageData.RefToMessageID != samplemid {
		t.Errorf("SessionRSHeader.MessageHeader.MessageData.RefToMessageID expect %s, got %s", samplemid, resp.Header.MessageHeader.MessageData.RefToMessageID)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("SessionRSHeader.Security.BinarySecurityToken.Value \nexpected: %s, \nrecieved: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Body.SessionCreateRS.ConversationID != sampleconvid {
		t.Errorf("SessionRSBody.SessionCreateRS.ConversationID expect %s, got %s", sampleconvid, resp.Body.SessionCreateRS.ConversationID)
	}
	if resp.Body.SessionCreateRS.Status != "Approved" {
		t.Errorf("resp.SessionRSBody.SessionCreateRS.Status expect %s, got %s", "Approved", resp.Body.SessionCreateRS.Status)
	}
}

func TestSessionCreateResponseUnAuth(t *testing.T) {
	resp := SessionCreateResponse{}
	err := xml.Unmarshal(sampleSessionUnAuth, &resp)

	if err != nil {
		t.Errorf("Error unmarshaling session create response %s \nERROR: %v", sampleSessionSuccessResponse, err)
	}

	if resp.Header.MessageHeader.Action != StatusErrorRS {
		t.Errorf("Header.MessageHeader.Action expect: %s, got: %s", StatusErrorRS, resp.Header.MessageHeader.Action)
	}
	if resp.Header.Security.BinarySecurityToken.Value != "" {
		t.Errorf("SessionRSHeader.Security.BinarySecurityToken.Value \nexpected: '%s', \nrecieved: %s", "", resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Body.SessionCreateRS.Status != "" {
		t.Errorf("SessionRSBody.SessionCreateRS.Status expect %s, got %s", "", resp.Body.SessionCreateRS.Status)
	}
	if resp.Body.Fault.String != "Authentication failed" {
		t.Errorf("SessionRSBody.Fault.String expect %s, got %s", "Authentication failed", resp.Body.Fault.String)
	}
	if resp.Body.Fault.Code != sampleSessionNoAuthFaultCode {
		t.Errorf("Body.Fault.String expect: %s, got: %s", sampleSessionNoAuthFaultCode, resp.Body.Fault.String)
	}
	if resp.Body.Fault.Detail.StackTrace != sampleSessionNoAuthStackTrace {
		t.Errorf("SessionRSBody.Fault.Detail.StackTrace expect %s, got %s", sampleSessionNoAuthStackTrace, resp.Body.Fault.Detail.StackTrace)
	}
}

func TestSOAPFaultFormat(t *testing.T) {
	resp := SessionCreateResponse{}
	xml.Unmarshal(sampleSessionUnAuth, &resp)
	if resp.Body.Fault.Ok() {
		t.Error("SOAPFault should exist")
	}

	format := resp.Body.Fault.Format()
	if fmt.Sprintf("%T", format) != fmt.Sprintf("%T", sbrerr.ErrorSoapFault{}) {
		t.Errorf("SOAPFault Format() should be <T> ErrorSoapFault, got: %T", format)
	}
	if format.Code != sbrerr.SoapFault {
		t.Errorf("SOAPFault Format.Code expect: %d, got: %d", sbrerr.SoapFault, format.Code)
	}
	if format.ErrMessage != sampleSessionNoAuthFaultString {
		t.Errorf("SOAPFault Format.ErrMessage expect: '%s', got: '%s'", sampleSessionNoAuthFaultString, format.ErrMessage)
	}
	if format.StackTrace != sampleSessionNoAuthStackTrace {
		t.Errorf("SOAPFault Format.StackTrace expect: %s, got: %s", sampleSessionNoAuthStackTrace, format.StackTrace)
	}
	if format.FaultCode != sampleSessionNoAuthFaultCode {
		t.Errorf("SOAPFault Format.FaultCode expect: %s, got: %s", sampleSessionNoAuthFaultCode, format.FaultCode)
	}
	if format.Error() != sampleSessionNoAuthFaultString {
		t.Errorf("SOAPFault Format.Error() expect: %s, got: %s", sampleSessionNoAuthFaultString, format.Error())
	}

}

func TestSessionCloseRequest(t *testing.T) {
	close := SessionCloseRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(samplefrom, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          samplepcc,
				ConversationID: sampleconvid,
				Service:        ServiceElem{"SessionCloseRQ", "OTA"},
				Action:         "SessionCloseRQ",
				MessageData: MessageDataElem{
					MessageID: samplemid,
					Timestamp: sampletime,
				},
			},
			Security: Security{
				XMLNSWsseBase:       BaseWsse,
				XMLNSWsu:            BaseWsuNameSpace,
				BinarySecurityToken: samplebinsectoken,
			},
		},
		Body: SessionCloseRQBody{
			SessionCloseRQ: SessionCloseRQ{
				POS: POSElem{
					Source: SourceElem{
						PseudoCityCode: samplepcc,
					},
				},
			},
		},
	}
	b, err := xml.Marshal(close)
	if err != nil {
		t.Error("Error marshaling session close rq", err)
	}
	if string(b) != string(sampleSessionCloseRQ) {
		t.Errorf("Expected marshal session close body \n sample: %s \n result: %s", string(sampleSessionCloseRQ), string(b))
	}
}

func TestBuildSessionCloseRequestMarshal(t *testing.T) {
	close := BuildSessionCloseRequest(sampleSessionConf)
	b, err := xml.Marshal(&close)
	if err != nil {
		t.Error("Error marshalling session envelope", err)
	}
	if string(b) != string(sampleSessionCloseRQ) {
		t.Errorf("Close request marshal \n sample: %s \n result: %s", string(sampleSessionCloseRQ), string(b))
	}
}
func BenchmarkBuildSessionCloseRequest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BuildSessionCloseRequest(sampleSessionConf)
	}
}
func BenchmarkBuildSessionCloseRequestMarshal(b *testing.B) {
	close := BuildSessionCloseRequest(sampleSessionConf)
	for n := 0; n < b.N; n++ {
		xml.Marshal(&close)
	}
}

func TestSessionCloseResponse(t *testing.T) {
	resp := SessionCloseResponse{}
	err := xml.Unmarshal(sampleSessionCloseRespSuccess, &resp)
	if err != nil {
		t.Errorf("Error unmarshaling session close response %s \nERROR: %v", sampleSessionCloseRespSuccess, err)
	}

	if resp.Body.SessionCloseRS.Status != "Approved" {
		t.Errorf("Body.SessionCloseRS.Status expect: %s, got: %s", "Approved", resp.Body.SessionCloseRS.Status)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
}

func TestSessionCloseResponseInvalidToken(t *testing.T) {
	resp := SessionCloseResponse{}
	err := xml.Unmarshal(sampleSessionCloseRespNoValidToken, &resp)
	if err != nil {
		t.Errorf("Error unmarshaling session close response %s \nERROR: %v", sampleSessionCloseRespNoValidToken, err)
	}

	if resp.Body.SessionCloseRS.Status != "" {
		t.Errorf("Body.SessionCloseRS.Status expect: %s, got: %s", "", resp.Body.SessionCloseRS.Status)
	}
	if resp.Header.MessageHeader.Action != StatusErrorRS {
		t.Errorf("Header.MessageHeader.Action expect: %s, got: %s", StatusErrorRS, resp.Header.MessageHeader.Action)
	}
	if resp.Body.Fault.String != sampleSessionInvalidTokenString {
		t.Errorf("Body.Fault.String expect: %s, got: %s", sampleSessionInvalidTokenString, resp.Body.Fault.String)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("Header.Security.BinarySecurityToken.Value expect: %s, got: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
}

func TestSessionValidateResponse(t *testing.T) {
	resp := SessionValidateResponse{}
	err := xml.Unmarshal(sampleSessionValidateRespSuccess, &resp)
	if err != nil {
		t.Errorf("Error unmarshaling session close response %s \nERROR: %v", sampleSessionCloseRespSuccess, err)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("Header.Security.BinarySecurityToken.Value expect: %s, got: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Header.MessageHeader.ConversationID != sampleconvid {
		t.Errorf("Header.MessageHeader.ConversationID exptc: %s, got: %s", sampleconvid, resp.Header.MessageHeader.ConversationID)
	}
	if resp.Header.MessageHeader.MessageData.RefToMessageID != samplemid {
		t.Errorf("Header.MessageHeader.MessageData.RefToMessageID expect: %s, got: %s", samplemid, resp.Header.MessageHeader.MessageData.RefToMessageID)
	}
	if resp.Header.MessageHeader.MessageData.Timestamp != sampletime {
		t.Errorf("Header.MessageHeader.MessageData.Timestamp expect: %s, got: %s", sampletime, resp.Header.MessageHeader.MessageData.Timestamp)
	}
}

func TestSessionValidateResponseInvalidToken(t *testing.T) {
	resp := SessionValidateResponse{}
	err := xml.Unmarshal(sampleSessionValidateRSInvalidTokenRS, &resp)
	if err != nil {
		t.Errorf("Error unmarshaling session close response %s \nERROR: %v", sampleSessionValidateRSInvalidTokenRS, err)
	}
	if resp.Header.MessageHeader.Action != StatusErrorRS {
		t.Errorf("Header.MessageHeader.Action expect: %s, got: %s", StatusErrorRS, resp.Header.MessageHeader.Action)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("Header.Security.BinarySecurityToken.Value expect: %s, got: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Body.Fault.String != sampleSessionInvalidTokenString {
		t.Errorf("Body.Fault.String expect: %s, got: %s", sampleSessionInvalidTokenString, resp.Body.Fault.String)
	}
	if resp.Body.Fault.Code != sampleSessionInvalidTokenFaultCode {
		t.Errorf("Body.Fault.String expect: %s, got: %s", sampleSessionInvalidTokenFaultCode, resp.Body.Fault.String)
	}
	if resp.Body.Fault.Detail.StackTrace != sampleSessionInvalidTokenStackTrace {
		t.Errorf("Body.Fault.Detail.StackTrace expect: %s, got: %s", sampleSessionInvalidTokenStackTrace, resp.Body.Fault.Detail.StackTrace)
	}
}
func TestBuildSessionValidateRequestMarshal(t *testing.T) {
	close := BuildSessionValidateRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime)

	b, err := xml.Marshal(&close)
	if err != nil {
		t.Error("Error marshalling session envelope", err)
	}
	if string(b) != string(sampleSessionValidateRQ) {
		t.Error("Session envelope with values does not match test sample")
	}
}
func BenchmarkBuildSessionValidateRequest(b *testing.B) {
	for n := 0; n < b.N; n++ {
		BuildSessionValidateRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime)
	}
}
func BenchmarkBuildSessionValidateRequestMarshal(b *testing.B) {
	close := BuildSessionValidateRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime)
	for n := 0; n < b.N; n++ {
		xml.Marshal(&close)
	}
}

func TestCallSessionCreateServiceNoExist(t *testing.T) {
	resp, err := CallSessionCreate(serverDown.URL, SessionCreateRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m, _ := regexp.MatchString("Post http://127.0.0.1.*", err.Error())
	if !m {
		t.Error("Expect error to match string with url", err)
	}
	r := SessionCreateResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}
func TestCallSessionCloseServiceNoExist(t *testing.T) {
	resp, err := CallSessionClose(serverDown.URL, SessionCloseRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m, _ := regexp.MatchString("Post http://127.0.0.1.*", err.Error())
	if !m {
		t.Error("Expect error to match string with url", err)
	}
	r := SessionCloseResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}
func TestCallSessionValidateServiceNoExist(t *testing.T) {
	resp, err := CallSessionValidate(serverDown.URL, SessionValidateRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m, _ := regexp.MatchString("Post http://127.0.0.1.*", err.Error())
	if !m {
		t.Error("Expect error to match string with url", err)
	}
	r := SessionValidateResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}

func TestCallSessionCreateBadBody(t *testing.T) {
	resp, err := CallSessionCreate(serverBadBody.URL, SessionCreateRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m := "XML syntax error on line 1: invalid character entity & (no semicolon)"
	if err.Error() != m {
		t.Errorf("CallSessionCreate XML error expect: %s, got: %s", m, err.Error())
	}
	r := SessionCreateResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}
func TestCallSessionCloseBadBody(t *testing.T) {
	resp, err := CallSessionClose(serverBadBody.URL, SessionCloseRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m := "XML syntax error on line 1: invalid character entity & (no semicolon)"
	if err.Error() != m {
		t.Errorf("CallSessionCreate XML error expect: %s, got: %s", m, err.Error())
	}
	r := SessionCloseResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}
func TestCallSessionValidateBadBody(t *testing.T) {
	resp, err := CallSessionValidate(serverBadBody.URL, SessionValidateRequest{})
	if err == nil {
		t.Error("Expect error", err)
	}
	m := "XML syntax error on line 1: invalid character entity & (no semicolon)"
	if err.Error() != m {
		t.Errorf("CallSessionCreate XML error expect: %s, got: %s", m, err.Error())
	}
	r := SessionValidateResponse{}
	if resp != r {
		t.Error("Expect empty response", resp)
	}
}

func TestCallSessionCreateSuccess(t *testing.T) {
	//Mock Sabre Web Services
	req := BuildSessionCreateRequest(sampleSessionConf)
	resp, err := CallSessionCreate(serverCreateRQ.URL, req)
	if err != nil {
		t.Error("Error making request CallSessionCreate", err)
	}
	/*
		things we need from session create to run SessionPool....
		response fro Sabre will have a To field, which co-ordinates
		with the From field we sent in the Request... pull the To
		field from Sabre Response and use to for the From field in
		the next request for Avail,Validate,Close, etc...
	*/
	if resp.Header.MessageHeader.To.PartyID.Value != samplefrom {
		t.Errorf("Header.MessageHeader.To.PartyID.Value epxect: %s, got: %s", samplefrom, resp.Header.MessageHeader.To.PartyID.Value)
	}
	if resp.Header.MessageHeader.CPAID != samplepcc {
		t.Errorf("Header.MessageHeader.CPAID epxect: %s, got: %s", samplepcc, resp.Header.MessageHeader.CPAID)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("Header.Security.BinarySecurityToken.Value epxect: %s, got: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Header.MessageHeader.ConversationID != sampleconvid {
		t.Errorf("Header.MessageHeader.ConversationID epxect: %s, got: %s", sampleconvid, resp.Header.MessageHeader.ConversationID)
	}
	if resp.Header.MessageHeader.MessageData.RefToMessageID != samplemid {
		t.Errorf("Header.MessageHeader.MessageData.RefToMessageID epxect: %s, got: %s", samplemid, resp.Header.MessageHeader.MessageData.RefToMessageID)
	}
}
func BenchmarkCallSessionCreate(b *testing.B) {
	req := BuildSessionCreateRequest(sampleSessionConf)
	for n := 0; n < b.N; n++ {
		CallSessionCreate(serverCreateRQ.URL, req)
	}
}

func TestCallSessionClose(t *testing.T) {
	req := BuildSessionCloseRequest(sampleSessionConf)
	resp, err := CallSessionClose(serverCloseRQ.URL, req)
	if err != nil {
		t.Error("Error making request CallSessionClose", err)
	}
	if resp.Body.SessionCloseRS.Status != "Approved" {
		t.Errorf("Body.SessionCloseRS.Status expect: %s, got: %s", "Approved", resp.Body.SessionCloseRS.Status)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
}
func BenchmarkCallSessionClose(b *testing.B) {
	req := BuildSessionCloseRequest(sampleSessionConf)
	for n := 0; n < b.N; n++ {
		CallSessionClose(serverCloseRQ.URL, req)
	}
}

func TestCallSessionValidate(t *testing.T) {
	req := BuildSessionValidateRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime)
	resp, err := CallSessionValidate(serverValidateRQ.URL, req)
	if err != nil {
		t.Error("Error making request CallSessionValidate", err)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
	if resp.Header.Security.BinarySecurityToken.Value != samplebinsectoken {
		t.Errorf("Header.Security.BinarySecurityToken.Value expect: %s, got: %s", samplebinsectoken, resp.Header.Security.BinarySecurityToken.Value)
	}
	if resp.Header.MessageHeader.ConversationID != sampleconvid {
		t.Errorf("Header.MessageHeader.ConversationID exptc: %s, got: %s", sampleconvid, resp.Header.MessageHeader.ConversationID)
	}
	if resp.Header.MessageHeader.MessageData.RefToMessageID != samplemid {
		t.Errorf("Header.MessageHeader.MessageData.RefToMessageID expect: %s, got: %s", samplemid, resp.Header.MessageHeader.MessageData.RefToMessageID)
	}
	if resp.Header.MessageHeader.MessageData.Timestamp != sampletime {
		t.Errorf("Header.MessageHeader.MessageData.Timestamp expect: %s, got: %s", sampletime, resp.Header.MessageHeader.MessageData.Timestamp)
	}
}
func BenchmarkCallSessionValidate(b *testing.B) {
	req := BuildSessionValidateRequest(samplefrom, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime)
	for n := 0; n < b.N; n++ {
		CallSessionValidate(serverValidateRQ.URL, req)
	}
}
