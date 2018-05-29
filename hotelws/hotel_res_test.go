package hotelws

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestHotelResSet(t *testing.T) {
	body := SetHotelResBody()
	req := BuildHotelResRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, body)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling hotel resercation request", err)
	}
	fmt.Printf("%s\n", b)
}
