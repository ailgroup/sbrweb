package itin

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

/*
This service is used in conjunction with Formats and Filters and handles the transaction to move/copy profile data into a PNR to create bookings. It allows a user to move any number of data elements from the profiles with a single transaction. As a result, the functions to copy profile information into the PNR is minimized by a single scan of the system.

	Filters: Filters are used to store information about which data elements, associated profiles, and formats are copied into a PNR. Users are able to pre-define filters to designate which filter is used as the profile default.
	Formats: Formats are used to store custom-defined Sabre (TPF) entries into a PNR. In turn, this allows a user to move profile migrate/copy data from an external system into Sabre Profiles.

http://webservices.sabre.com/drc/providerdoc/PPP/Sabre_Profiles_Technical_User_Guide.pdf
https://richmedia.sabre.com/docs/teletraining/FlexBasicsSG.pdf
*/

// ProfileToPNRBody holds namespaced body
type ProfileToPNRBody struct {
	XMLName        xml.Name `xml:"soap-env:Body"`
	ProfileToPNRRQ ProfileToPNRRQ
}

// ProfileToPNRRequest wrapper for soap payload.
type ProfileToPNRRequest struct {
	srvc.Envelope
	Header srvc.SessionHeader
	Body   ProfileToPNRBody
}

// AssociatedProfile is any profile that does not exist within your profile, but has been linked to your existing profile. These are custom and not all profiles will have associations; they will be created by whoever manages your sabre account.
type AssociatedProfile struct {
	XMLName              xml.Name `xml:"AssociatedProfiles"`
	AssocUniqueID        string   `xml:"AssocUniqueID,attr"`
	AssocProfileName     string   `xml:"AssocProfileName,attr"`
	AssocProfileTypeCode string   `xml:"AssocProfileTypeCode,attr"`
	DomainID             string   `xml:"DomainID,attr"`
	OrderSequenceNo      string   `xml:"OrderSequenceNo,attr"`
}

// Filter is custom created subset of profile information. Not all profiles have them, they must be created by whoever manages your sabrea account.
type Filter struct {
	XMLName           xml.Name `xml:"Filter"`
	FilterID          string   `xml:"FilterID,attr"`
	DomainID          string   `xml:"DomainID,attr"`
	ClientCode        string   `xml:"ClientCode,attr"`
	ClientContextCode string   `xml:"ClientContextCode,attr"`
	FilterName        string   `xml:"FilterName,attr"`
}

// Profile is main element containing ClientCode: ; ClientContext: ; PofileTypeCode: agency[AGY], agent[AGT], corporate[CRP], group[GRP], operational[OPX], traveler[TVL] (likely CRP is what you want). DomainID is your PCC/PsuedoCityCode. The values for UniqueID (your ProfileID); ProfileName (name of profile) and PNRMoveOrderSeqNo exist as values on your profile in sabre account.
type Profile struct {
	XMLName            xml.Name `xml:"Profile"`
	ClientCode         string   `xml:"ClientCode,attr"`
	ClientContextCode  string   `xml:"ClientContextCode,attr"`
	DomainID           string   `xml:"DomainID,attr"`
	ProfileTypeCode    string   `xml:"ProfileTypeCode,attr"`
	UniqueID           string   `xml:"UniqueID,attr,omitempty"`
	ProfileName        string   `xml:"ProfileName,attr"`
	PNRMoveOrderSeqNo  string   `xml:"PNRMoveOrderSeqNo,attr"`
	Filter             *Filter
	AssociatedProfiles []*AssociatedProfile
}

type FilterPath struct {
	XMLName xml.Name `xml:"FilterPath"`
	Profile Profile
}

// BuildFilterPathForProfileOnly constructs a typical Profile with all necessary attributes. It does not build a Filter or AssociatedProfiles.
// NOTE: DomainID is your PsuedoCityCode (PCC, or IPCC), UniqueID is in ProfileIDs, ProfileName is in ProfileNames ... all of these values are can be retrieved from the srvc.SessionConf struct OR they can be passed free-hand into this function.
// FilterPath is returned as pointer in case one needs to append more data (e.g., Filter); this is consistent with paramter that BuildProfileToPNRRequest expects.
func BuildFilterPathForProfileOnly(clientCode, clientContext, domID, profileType, uniqueID, profileName, moveSeq string) *FilterPath {
	return &FilterPath{
		Profile: Profile{
			ClientCode:        clientCode,
			ClientContextCode: clientContext,
			DomainID:          domID,
			ProfileTypeCode:   profileType,
			UniqueID:          uniqueID,
			ProfileName:       profileName,
			PNRMoveOrderSeqNo: moveSeq,
		},
	}
}

// ProfileToPNRRQ root element
type ProfileToPNRRQ struct {
	XMLName xml.Name `xml:"Sabre_OTA_ProfileToPNRRQ"`
	Target  string   `xml:"Target,attr"`
	//Timestamp  string   `xml:"Timestamp,attr"`
	Version    string `xml:"Version,attr"`
	XSISchema  string `xml:"xsi:schemaLocation,attr"`
	XMLNS      string `xml:"xmlns,attr"`
	XMLXSI     string `xml:"xmlns:xsi,attr"`
	FilterPath FilterPath
}

// BuildProfileToPNRRequest construct payload for request to copy a profile into a PNR. The FilterPath is parameter since it may need to be constructed in various ways depending on need; for example BuildFilterPathForProfileOnly buils a simple FilterPath Profile
func BuildProfileToPNRRequest(c *srvc.SessionConf, fp *FilterPath) ProfileToPNRRequest {
	return ProfileToPNRRequest{
		Envelope: srvc.CreateEnvelope(),
		Header: srvc.SessionHeader{
			MessageHeader: srvc.MessageHeader{
				MustUnderstand: srvc.SabreMustUnderstand,
				EbVersion:      srvc.SabreEBVersion,
				From: srvc.FromElem{
					PartyID: srvc.CreatePartyID(c.From, srvc.PartyIDTypeURN),
				},
				To: srvc.ToElem{
					PartyID: srvc.CreatePartyID(srvc.SabreToBase, srvc.PartyIDTypeURN),
				},
				CPAID:          c.PCC,
				ConversationID: c.Convid,
				Service:        srvc.ServiceElem{Value: "EPS_ProfileToPNR", Type: "sabreXML"},
				Action:         "EPS_ProfileToPNRRQ",
				MessageData: srvc.MessageDataElem{
					MessageID: c.Msgid,
					Timestamp: c.Timestr,
				},
			},
			Security: srvc.Security{
				XMLNSWsseBase:       srvc.BaseWsse,
				XMLNSWsu:            srvc.BaseWsuNameSpace,
				BinarySecurityToken: c.Binsectok,
			},
		},
		Body: ProfileToPNRBody{
			ProfileToPNRRQ: ProfileToPNRRQ{
				//http://www.sabre.com/eps/schemas
				XMLNS:     srvc.BaseXINamespace,
				XMLXSI:    srvc.BaseXSINamespace,
				XSISchema: srvc.BaseProfileSchema,
				Target:    "Production",
				//Timestamp:  srvc.SabreTimeNowFmt(),
				Version:    "6.50",
				FilterPath: *fp,
			},
		},
	}
}

func (r ResponseMessage) Ok() bool {
	return r.Success.Val == "Success"
}

type ResponseMsgSuccess struct {
	Val string `xml:",chardata"`
}

type ResponseMessage struct {
	Success ResponseMsgSuccess
}

type ProfileToPNRRS struct {
	XMLName         xml.Name `xml:"Sabre_OTA_ProfileToPNRRS"`
	ResponseMessage ResponseMessage
}
type ProfileToPNRResponse struct {
	Envelope srvc.EnvelopeUnMarsh
	Header   srvc.SessionHeaderUnmarsh
	Body     struct {
		ProfileToPNRRS ProfileToPNRRS
		Fault          srvc.SOAPFault
	}
	ErrorSabreService sbrerr.ErrorSabreService
	ErrorSabreXML     sbrerr.ErrorSabreXML
}

// CallProfileToPNR to execute ProfileToPNRRequest, which must be done in order to finish the booking transaction.
func CallProfileToPNR(serviceURL string, req ProfileToPNRRequest) (ProfileToPNRResponse, error) {
	endT := ProfileToPNRResponse{}
	byteReq, _ := xml.Marshal(req)
	srvc.LogSoap.Printf("CallProfileToPNR-REQUEST %s \n\n", byteReq)

	//post payload
	resp, err := http.Post(serviceURL, "text/xml", bytes.NewBuffer(byteReq))
	if err != nil {
		endT.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallProfileToPNR,
			sbrerr.BadService,
		)
		return endT, endT.ErrorSabreService
	}
	// parse payload body into []byte buffer from net Response.ReadCloser
	// note ioutil.ReadAll(resp.Body) has no cap on size and can create memory problems
	bodyBuffer := new(bytes.Buffer)
	_, err = io.Copy(bodyBuffer, resp.Body)
	srvc.LogSoap.Printf("CallProfileToPNR-RESPONSE %s \n\n", bodyBuffer)
	//close body no defer
	resp.Body.Close()
	//handle and return error if bad body
	if err != nil {
		endT.ErrorSabreService = sbrerr.NewErrorSabreService(
			err.Error(),
			sbrerr.ErrCallProfileToPNR,
			sbrerr.BadParse,
		)
		return endT, endT.ErrorSabreService
	}

	//marshal bytes sabre response body into availResp response struct
	err = xml.Unmarshal(bodyBuffer.Bytes(), &endT)
	if err != nil {
		endT.ErrorSabreXML = sbrerr.NewErrorSabreXML(
			err.Error(),
			sbrerr.ErrCallProfileToPNR,
			sbrerr.BadParse,
		)
		return endT, endT.ErrorSabreXML
	}
	if !endT.Body.Fault.Ok() {
		return endT, sbrerr.NewErrorSoapFault(endT.Body.Fault.Format().ErrMessage)
	}

	if !endT.Body.ProfileToPNRRS.ResponseMessage.Ok() {
		return endT, errors.New("CallProfileToPNR no Success")
	}
	return endT, nil
}
