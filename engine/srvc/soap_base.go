package srvc

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

const (
	BaseEBNameSpace       = "http://www.ebxml.org/namespaces/messageHeader"
	BaseNS                = "http://schemas.xmlsoap.org/soap/envelope/"
	BaseWebServicesNS     = "http://webservices.sabre.com/sabreXML/2011/10"
	BaseWsse              = "http://schemas.xmlsoap.org/ws/2002/12/secext"
	BaseWsuNameSpace      = "http://schemas.xmlsoap.org/ws/2002/12/utility"
	BaseXSDNameSpace      = "http://www.w3.org/2001/XMLSchema"
	BaseXSINamespace      = "http://www.w3.org/2001/XMLSchema-instance"
	PartyIDTypeURN        = "urn:x12.org:IO5:01"
	SabreEBVersion        = "2.0.0"
	SabreMustUnderstand   = "1"
	SabreToBase           = "webservices.sabre.com"
	StandardTimeFormatter = "2006-01-02T15:04:05Z"
	//StatusErrorRS is the string value error response when SOAP request->response had an error, typically found in the Header.MessageHeader.Action. Usually, any SOAP response with Action="ErrorRS" will also have a SOAPFault body with more informative error codes. This can be used as an easy way to identify and error.
	StatusErrorRS = "ErrorRS"

	baseOTANameSpace   = "http://www.opentravel.org/OTA/2002/11"
	baseXlinkNameSpace = "http://www.w3.org/2001/xlink"
	letterBytes        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits      = 6                    // 6 bits to represent a letter index
	letterIdxMask      = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax       = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	sabreDefaultDomain = "DEFAULT"
)

var (
	logSoap            = &log.Logger{}
	logSession         = &log.Logger{}
	binaryTokenMatcher = regexp.MustCompile(`-\d.*$`)
)

func init() {
	setUpLogging()
}

func setUpLogging() {
	soapL, err := os.Create("sabre_web_soap.log")
	if err != nil {
		log.Fatal("no soap log file")
	}
	sessL, err := os.Create("sabre_web_session.log")
	if err != nil {
		log.Fatal("no session log file")
	}

	logSoap = log.New(soapL, "[sabre-soap] ", (log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile))
	logSession = log.New(sessL, "[sabre-sess] ", (log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile))
}

// Envelope is wrapper with namespace prefix definitions for payload
type Envelope struct {
	XMLName    xml.Name `xml:"soap-env:Envelope"`
	XMLNSbase  string   `xml:"xmlns:soap-env,attr"`
	XMLNSeb    string   `xml:"xmlns:eb,attr"`
	XMLNSxlink string   `xml:"xmlns:xlink,attr"`
	XMLNSxsd   string   `xml:"xmlns:xsd,attr"`
}

// EnvelopeUnMarsh is wrapper to unmarshal with namespace prefix
type EnvelopeUnMarsh struct {
	XMLName    xml.Name `xml:"Envelope"`
	XMLNSbase  string   `xml:"soap-env,attr,omitempty"`
	XMLNSeb    string   `xml:"eb,attr,omitempty"`
	XMLNSxlink string   `xml:"xlink,attr,omitempty"`
	XMLNSxsd   string   `xml:"xsd,attr,omitempty"`
}

// PartyIDElem is the parties id
type PartyIDElem struct {
	Value string `xml:",chardata"`
	Type  string `xml:"type,attr"`
}

// FromElem who is sending message
type FromElem struct {
	XMLName xml.Name    `xml:"eb:From"`
	PartyID PartyIDElem `xml:"eb:PartyId"`
}

// FromElemUnmarsh is wrapper to unmarshal with namespace prefix
type FromElemUnmarsh struct {
	XMLName xml.Name    `xml:"From"`
	PartyID PartyIDElem `xml:"PartyId"`
}

// ToElem whom message message is sent
type ToElem struct {
	XMLName xml.Name    `xml:"eb:To"`
	PartyID PartyIDElem `xml:"eb:PartyId"`
}

// ToElemUnmarsh is wrapper to unmarshal with namespace prefix
type ToElemUnmarsh struct {
	XMLName xml.Name    `xml:"To"`
	PartyID PartyIDElem `xml:"PartyId"`
}

//MessageDataElem holds unique identifiers coupled with time
type MessageDataElem struct {
	XMLName    xml.Name `xml:"eb:MessageData"`
	MessageID  string   `xml:"eb:MessageId"`
	Timestamp  string   `xml:"eb:Timestamp,omitempty"`
	TimeToLive string   `xml:"eb:TimeToLive,omitempty"`
}

//MessageDataElemUnmarsh holds unique identifiers coupled with time
type MessageDataElemUnmarsh struct {
	XMLName        xml.Name `xml:"MessageData"`
	MessageID      string   `xml:"MessageId"`
	Timestamp      string   `xml:"Timestamp"`
	TimeToLive     string   `xml:"TimeToLive"`
	RefToMessageID string   `xml:"RefToMessageId"`
}

// ServiceElem defines type of service
type ServiceElem struct {
	Value string `xml:",chardata"`
	Type  string `xml:"eb:type,attr"`
}

// MessageHeader contains message specific data such as credentials, from, to, conversation id, soap service, soap action, etc...
type MessageHeader struct {
	XMLName        xml.Name `xml:"eb:MessageHeader"`
	MustUnderstand string   `xml:"soap-env:mustUnderstand,attr"`
	EbVersion      string   `xml:"eb:version,attr"`
	From           FromElem
	To             ToElem
	CPAID          string      `xml:"eb:CPAId"`
	ConversationID string      `xml:"eb:ConversationId"`
	Service        ServiceElem `xml:"eb:Service"`
	Action         string      `xml:"eb:Action"`
	MessageData    MessageDataElem
}

// MessageHeaderUnmarsh wrapper to unmarshal with namespace prefix
type MessageHeaderUnmarsh struct {
	XMLName        xml.Name `xml:"MessageHeader"`
	MustUnderstand string   `xml:"mustUnderstand,attr"`
	EbVersion      string   `xml:"version,attr"`
	From           FromElemUnmarsh
	To             ToElemUnmarsh
	CPAID          string `xml:"CPAId"`
	ConversationID string `xml:"ConversationId"`
	Service        string `xml:"Service"`
	Action         string `xml:"Action"`
	MessageData    MessageDataElemUnmarsh
}

// UsernameTokenElem contains user security info
type UsernameTokenElem struct {
	XMLName      xml.Name `xml:"wsse:UsernameToken"`
	Username     string   `xml:"wsse:Username"`
	Password     string   `xml:"wsse:Password"`
	Organization string   `xml:"Organization"`
	Domain       string   `xml:"Domain"`
}

// UsernameTokenElemUnmarsh wrapper to unmarshal with namespace prefix
type UsernameTokenElemUnmarsh struct {
	XMLName      xml.Name `xml:"UsernameToken"`
	Username     string   `xml:"Username"`
	Password     string   `xml:"Password"`
	Organization string   `xml:"Organization"`
	Domain       string   `xml:"Domain"`
}

// Security is wrapper with namespace prefix definitions for payload
// Pointer to UsernameTokenElem will ingore it if empty; not used on Close and Validate services.
// Omitempty on BinarySecurityToken will ingore when empty, not used on Create service.
type Security struct {
	XMLName             xml.Name `xml:"wsse:Security"`
	XMLNSWsseBase       string   `xml:"xmlns:wsse,attr"`
	XMLNSWsu            string   `xml:"xmlns:wsu,attr"`
	UserNameToken       *UsernameTokenElem
	BinarySecurityToken string `xml:"wsse:BinarySecurityToken,omitempty"`
}

// BinarySecurityTokenUnmarsh returned from sabre in security; see SecurityUnmarsh
type BinarySecurityTokenUnmarsh struct {
	Value        string `xml:",chardata"`
	ValueType    string `xml:"valueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
}

// SecurityUnmarsh wrapper to unmarshal with namespace prefix
type SecurityUnmarsh struct {
	XMLName             xml.Name `xml:"Security"`
	XMLNSWsseBase       string   `xml:"wsse,attr"`
	XMLNSWsu            string   `xml:"xmlns:wsu,attr"`
	UserNameToken       UsernameTokenElemUnmarsh
	BinarySecurityToken BinarySecurityTokenUnmarsh
}

// SessionHeader header of session
type SessionHeader struct {
	XMLName       xml.Name `xml:"soap-env:Header"`
	MessageHeader MessageHeader
	Security      Security
}

// SessionHeaderUnmarsh wrapper to unmarshal with namespace prefix
type SessionHeaderUnmarsh struct {
	XMLName       xml.Name `xml:"Header"`
	MessageHeader MessageHeaderUnmarsh
	Security      SecurityUnmarsh
}

// ReferenceElem for Manifest
type ReferenceElem struct {
	XMLName xml.Name `xml:"eb:Reference"`
	Href    string   `xml:"xlink:href,attr"`
	Type    string   `xml:"xlink:type,attr"`
}

// Manifest for soap body
type Manifest struct {
	XMLName        xml.Name `xml:"eb:Manifest"`
	MustUnderstand string   `xml:"soap-env:mustUnderstand,attr"`
	EbVersion      string   `xml:"eb:version,attr"`
	Reference      ReferenceElem
}

// SourceElem for POS session create
type SourceElem struct {
	PseudoCityCode string `xml:"PseudoCityCode,attr"`
}

// POSElem for session create
type POSElem struct {
	Source SourceElem `xml:"Source"`
}

// Ok simple check for soap fault; return true if empty string
func (fault SOAPFault) Ok() bool {
	return fault.Code == ""
}

// Format on SOAPFault for error checking and printing
func (fault SOAPFault) Format() sbrerr.ErrorSoapFault {
	msg := fault.Detail.ApplicationResults.Error.SystemSpecificResults.Message
	shrt := fault.Detail.ApplicationResults.Error.SystemSpecificResults.ShortText
	return sbrerr.ErrorSoapFault{
		StackTrace: fault.Detail.StackTrace,
		FaultCode:  fault.Code,
		ErrMessage: fmt.Sprintf("%s %s %s", msg, fault.String, shrt),
		Code:       sbrerr.SoapFault,
	}
}

//SOAPFault catching error messages
type SOAPFault struct {
	XMLName xml.Name `xml:"Fault"`
	Code    string   `xml:"faultcode"`
	String  string   `xml:"faultstring"`
	Actor   string   `xml:"faultactor"`
	Detail  struct {
		StackTrace         string `xml:"StackTrace"`
		ApplicationResults struct {
			Error struct {
				Type                  string `xml:"type,attr"`
				SystemSpecificResults struct {
					Message   string `xml:"Message"`
					ShortText string `xml:"ShortText"`
				} `xml:"SystemSpecificResults"`
			} `xml:"Error"`
		} `xml:"ApplicationResults"`
	} `xml:"detail"`
}

// helper to crete message header
func CreateEnvelope() Envelope {
	return Envelope{
		XMLNSbase:  BaseNS,
		XMLNSeb:    BaseEBNameSpace,
		XMLNSxlink: baseXlinkNameSpace,
		XMLNSxsd:   BaseXSDNameSpace,
	}
}

// CreatePartyID helper ot make party ids
func CreatePartyID(partyValue, partyType string) PartyIDElem {
	return PartyIDElem{
		Value: partyValue,
		Type:  partyType,
	}
}

// CreateManifest helper
func CreateManifest() Manifest {
	return Manifest{
		MustUnderstand: SabreMustUnderstand,
		EbVersion:      SabreEBVersion,
		Reference: ReferenceElem{
			Href: "cid:rootelement",
			Type: "simple",
		},
	}
}

// SabreTokenParse rips off redundant info into easy to read chunk, mostly for logging.
// Take this 'Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!-3174053563846592370!1390092!0'
// Return this '-3174053563846592370!1390092!0'
func SabreTokenParse(tok string) string {
	if tok == "" {
		return "none"
	}
	return binaryTokenMatcher.FindAllString(tok, -1)[0]
}

// SabreTimeNowFmt returns time.Now in format: '2017-11-27T09:58:31Z'
func SabreTimeNowFmt() string {
	return time.Now().Format(StandardTimeFormatter)
}

// randStringBytesMaskImprSrc generate random string of specific length
func randStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// GenerateMessageID returns 'mid:20060102-15:04:05|urioe'
func GenerateMessageID() string {
	return "mid:" + time.Now().Format("20060102-15:04:05.99") + "|" + randStringBytesMaskImprSrc(5)
}

// GenerateConversationID returns 'cid:1Fv0Oq65|www.z.com'
func GenerateConversationID(from string) string {
	return "cid:" + randStringBytesMaskImprSrc(8) + "|" + from
}

// SessionCreateRQ for session create request
type SessionCreateRQ struct {
	XMLName xml.Name `xml:"ns:SessionCreateRQ"`
	XMLNS   string   `xml:"xmlns:ns,attr"`
	POS     POSElem  `xml:"POS"`
	//ReturnContextID bool `xml:"returnContextID,attr,omitempty"`
}

// SessionCreateRQBody body of session
type SessionCreateRQBody struct {
	XMLName         xml.Name `xml:"soap-env:Body"`
	Manifest        Manifest
	SessionCreateRQ SessionCreateRQ
}

// SessionCreateRequest is wrapper with namespace prefix definitions for payload
type SessionCreateRequest struct {
	Envelope
	Header SessionHeader
	Body   SessionCreateRQBody
}

type SessionConf struct {
	ServiceURL  string
	From        string
	PCC         string
	Binsectok   string
	Convid      string
	Msgid       string
	Timestr     string
	Username    string
	Password    string
	AppTimeZone *time.Location
}

// SetTime updates the timestamp. Pass around SessionConf and update the timestamp for any new request
func (s *SessionConf) SetTime() *SessionConf {
	s.Timestr = SabreTimeNowFmt()
	return s
}

// SetBinSec updates the timestamp. Pass around SessionConf and update the timestamp for any new request
func (s *SessionConf) SetBinSec(session SessionCreateResponse) *SessionConf {
	s.Binsectok = session.Header.Security.BinarySecurityToken.Value
	return s
}

// BuildSessionCreateRequest build session create envelope for request
// CPAID, Organization, and PseudoCityCode all use the PCC/iPCC. ConversationID is typically a contact email address with unique identifier to the request. MessageID is typically a timestamped identifier to locate specific queries: it should contai a company identifier.
//func BuildSessionCreateRequest(from, pcc, convid, mid, time, username, password string) SessionCreateRequest {
func BuildSessionCreateRequest(c *SessionConf) SessionCreateRequest {
	return SessionCreateRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(c.From, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          c.PCC,
				ConversationID: c.Convid,
				Service:        ServiceElem{"SessionCreateRQ", "OTA"},
				Action:         "SessionCreateRQ",
				MessageData: MessageDataElem{
					MessageID: c.Msgid,
					Timestamp: c.Timestr,
				},
			},
			Security: Security{
				XMLNSWsseBase: BaseWsse,
				XMLNSWsu:      BaseWsuNameSpace,
				UserNameToken: &UsernameTokenElem{
					Username:     c.Username,
					Password:     c.Password,
					Organization: c.PCC,
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
						PseudoCityCode: c.PCC,
					},
				},
			},
		},
	}
}

// SessionCreateResponse is wrapper with namespace prefix definitions for payload
type SessionCreateResponse struct {
	Envelope EnvelopeUnMarsh
	Header   SessionHeaderUnmarsh
	Body     struct {
		SessionCreateRS struct {
			Version        string `xml:"version,attr"`
			Status         string `xml:"status,attr"`
			XMLNS          string `xml:"xmlns,attr"`
			ConversationID string `xml:"ConversationId"`
		} `xml:"SessionCreateRS"`
		Fault SOAPFault
	}
}

// CallSessionCreate to sabre web services.
func CallSessionCreate(serviceURL string, req SessionCreateRequest) (SessionCreateResponse, error) {
	sessionResponse := SessionCreateResponse{}
	//construct payload

	byteReq, _ := xml.Marshal(req)
	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallSessionCreate, sbrerr.BadService)
	}

	//parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//marshal byte body sabre response body into session envelope response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &sessionResponse)
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallSessionCreate, sbrerr.BadParse)
	}
	return sessionResponse, nil
}

// SessionCloseRQ for session create request
type SessionCloseRQ struct {
	XMLName xml.Name `xml:"SessionCloseRQ"`
	POS     POSElem  `xml:"POS"`
}

// SessionCloseRQBody body of session
type SessionCloseRQBody struct {
	XMLName        xml.Name `xml:"soap-env:Body"`
	SessionCloseRQ SessionCloseRQ
}

// SessionCloseRequest is wrapper with namespace prefix definitions for payload
type SessionCloseRequest struct {
	Envelope
	Header SessionHeader
	Body   SessionCloseRQBody
}

// BuildSessionCloseRequest build session Close envelope for request.
// CPAID, Organization, and PseudoCityCode all use the PCC/iPCC. ConversationID, MessageID, BinarySecurityToken must be from the existing session you wish to close.
func BuildSessionCloseRequest(c *SessionConf) SessionCloseRequest {
	return SessionCloseRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(c.From, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          c.PCC,
				ConversationID: c.Convid,
				Service:        ServiceElem{"SessionCloseRQ", "OTA"},
				Action:         "SessionCloseRQ",
				MessageData: MessageDataElem{
					MessageID: c.Msgid,
					Timestamp: c.Timestr,
				},
			},
			Security: Security{
				XMLNSWsseBase:       BaseWsse,
				XMLNSWsu:            BaseWsuNameSpace,
				BinarySecurityToken: c.Binsectok,
			},
		},
		Body: SessionCloseRQBody{
			SessionCloseRQ: SessionCloseRQ{
				POS: POSElem{
					Source: SourceElem{
						PseudoCityCode: c.PCC,
					},
				},
			},
		},
	}
}

// SessionCloseResponse is wrapper with namespace prefix definitions for payload
type SessionCloseResponse struct {
	Envelope EnvelopeUnMarsh
	Header   SessionHeaderUnmarsh
	Body     struct {
		SessionCloseRS struct {
			Version string `xml:"version,attr"`
			Status  string `xml:"status,attr"`
		} `xml:"SessionCloseRS"`
		Fault SOAPFault
	}
}

// CallSessionClose to sabre web services
func CallSessionClose(serviceURL string, e SessionCloseRequest) (SessionCloseResponse, error) {
	sessionResponse := SessionCloseResponse{}
	//construct payload

	byteReq, _ := xml.Marshal(e)
	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallSessionClose, sbrerr.BadService)

	}

	//parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//marshal byte body sabre response body into session envelope response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &sessionResponse)
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallSessionClose, sbrerr.BadParse)
	}
	return sessionResponse, nil
}

// SessionValidateRQ for session create request
type SessionValidateRQ struct {
	XMLName xml.Name `xml:"SessionValidateRQ"`
	POS     POSElem  `xml:"POS"`
}

// SessionValidateRQBody body of session
type SessionValidateRQBody struct {
	XMLName           xml.Name `xml:"soap-env:Body"`
	SessionValidateRQ SessionValidateRQ
}

// SessionValidateRequest is wrapper with namespace prefix definitions for payload
type SessionValidateRequest struct {
	Envelope
	Header SessionHeader
	Body   SessionValidateRQBody
}

// BuildSessionValidateRequest build session Validate envelope for request.
// CPAID, Organization, and PseudoCityCode all use the PCC/iPCC. ConversationID, MessageID, BinarySecurityToken must be from the existing session you wish to validate.
func BuildSessionValidateRequest(from, pcc, binsectoken, convid, mid, time string) SessionValidateRequest {
	return SessionValidateRequest{
		Envelope: CreateEnvelope(),
		Header: SessionHeader{
			MessageHeader: MessageHeader{
				MustUnderstand: SabreMustUnderstand,
				EbVersion:      SabreEBVersion,
				From:           FromElem{PartyID: CreatePartyID(from, PartyIDTypeURN)},
				To:             ToElem{PartyID: CreatePartyID(SabreToBase, PartyIDTypeURN)},
				CPAID:          pcc,
				ConversationID: convid,
				Service:        ServiceElem{"SessionValidateRQ", "OTA"},
				Action:         "SessionValidateRQ",
				MessageData: MessageDataElem{
					MessageID: mid,
					Timestamp: time,
				},
			},
			Security: Security{
				XMLNSWsseBase:       BaseWsse,
				XMLNSWsu:            BaseWsuNameSpace,
				BinarySecurityToken: binsectoken,
			},
		},
		Body: SessionValidateRQBody{
			SessionValidateRQ: SessionValidateRQ{
				POS: POSElem{
					Source: SourceElem{
						PseudoCityCode: pcc,
					},
				},
			},
		},
	}
}

// SessionValidateResponse is wrapper with namespace prefix definitions for payload
type SessionValidateResponse struct {
	Envelope EnvelopeUnMarsh
	Header   SessionHeaderUnmarsh
	Body     struct {
		SessionValidateRS struct {
			Version string `xml:"version,attr"`
			Status  string `xml:"status,attr"`
		} `xml:"SessionValidateRS"`
		Fault SOAPFault
	}
}

// CallSessionValidate to sabre web services
func CallSessionValidate(serviceURL string, req SessionValidateRequest) (SessionValidateResponse, error) {
	sessionResponse := SessionValidateResponse{}
	//construct payload

	byteReq, _ := xml.Marshal(req)
	buffer := bytes.NewBuffer(byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", buffer)
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreService(err.Error(), sbrerr.ErrCallSessionValidate, sbrerr.BadService)
	}

	//parse payload body into []byte buffer from net Response.ReadCloser
	// ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	io.Copy(bodyBuffer, resp.Body)
	resp.Body.Close()

	//marshal byte body sabre response body into session envelope response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &sessionResponse)
	if err != nil {
		return sessionResponse, sbrerr.NewErrorSabreXML(err.Error(), sbrerr.ErrCallSessionValidate, sbrerr.BadParse)
	}
	return sessionResponse, nil
}
