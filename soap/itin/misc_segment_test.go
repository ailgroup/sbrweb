package itin

import (
	"encoding/xml"
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
	segReq := BuildMiscSegmentRequest(sampleConf, samplebinsectoken, oth)
	_, err := xml.Marshal(segReq)
	if err != nil {
		t.Error("Error marshal build end transaction", err)
	}
}
