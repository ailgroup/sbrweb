package itin

import (
	"encoding/xml"
	"testing"
)

var (
	//sampleProfileToPNR             = []byte(``)
	sampleProfileClientCode        = "TN"  // X
	sampleProfileClientContext     = "TML" // X
	sampleProfileTypeCode          = "TRP" // X
	sampleProfileUniqueID          = "102598202"
	sampleProfileName              = "TestProfile_2013-04-30_08:24:42_GXE"
	sampleProfilePNRMoveOrderSeqNo = "1"
)

func TestProfileToPNRXML(t *testing.T) {
	filterPath := BuildFilterPathForProfileOnly(
		sampleProfileClientCode,
		sampleProfileClientContext,
		samplepcc, //PCC from sessionconf
		sampleProfileTypeCode,
		sampleProfileUniqueID, //ProfileIDs from sessionconf
		sampleProfileName,     //ProfileNames from sessionconf
		sampleProfilePNRMoveOrderSeqNo,
	)
	et := BuildProfileToPNRRequest(sampleConf, samplebinsectoken, filterPath)
	_, err := xml.Marshal(et)
	if err != nil {
		t.Error("Error marshal build end transaction", err)
	}
	/*
		if string(b) != string(sampleEndTReq) {
			t.Errorf("Expect end transaction \n given: %s \n built: %s", sampleEndTReq, b)
		}
	*/
}
