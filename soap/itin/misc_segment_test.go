package itin

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestMiscSegmentXML(t *testing.T) {
	oth := MiscSegment{
		DepartureDateTime: "05-16",
		NumberInParty:     2,
		Status:            "GK",
		Typ:               "OTH",
		OriginLocation:    OriginLocation{},
		VendorPrefs:       VendorPrefs{},
		//OriginLocation:    itin.OriginLocation{LocationCode: "YYY"},
		//VendorPrefs:       itin.VendorPrefs{Airline: itin.Airline{Code: "XX"}},
	}
	segReq := BuildMiscSegmentRequest(sampleConf, oth)
	b, err := xml.Marshal(segReq)
	if err != nil {
		t.Error("Error marshal build end transaction", err)
	}
	fmt.Printf("%s\n\n", b)
	/*
		if string(b) != string(sampleEndTReq) {
			t.Errorf("Expect end transaction \n given: %s \n built: %s", sampleEndTReq, b)
		}
	*/
}
