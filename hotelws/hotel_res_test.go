package hotelws

import (
	"encoding/xml"
	"testing"
)

func TestHotelResSet(t *testing.T) {
	body := SetHotelResBody(12, "GDPST", "MC", "2012-12", "1234567890", "Lastname")

	prefs, err := NewSpecialPrefs(
		WrittenConf(true),
	)
	if err != nil {
		t.Errorf("NewSpecialPrefs should not error %v", err)
	}
	body.addSpecialPrefs(prefs)
	body.addRoomTypeUnits(2)

	b := body.OTAHotelResRQ
	if !b.Hotel.SpecialPrefs.WrittenConfirmation.Ind {
		t.Errorf("SpecialPrefs WrittenConfirmation.Ind expcted: %v got: %v", true, b.Hotel.SpecialPrefs)
	}

	if b.Hotel.BasicPropertyRes.RPH != 12 {
		t.Error("RPH is wrong")
	}
	ccinfo := b.Hotel.Guarantee.CCInfo
	if ccinfo.PaymentCard.Code != "MC" {
		t.Error("PaymentCard.Code is wrong")
	}
	if ccinfo.PaymentCard.ExpireDate != "2012-12" {
		t.Error("PaymentCard.ExpireDate is wrong")
	}
	if ccinfo.PaymentCard.Number != "1234567890" {
		t.Error("PaymentCard.Number is wrong")
	}
	if ccinfo.PersonName.Last.Val != "Lastname" {
		t.Error("PersonName.Last is wrong")
	}
}

func TestHotelResBuild(t *testing.T) {
	body := SetHotelResBody(12, "GDPST", "MC", "2012-12", "1234567890", "Lastname")
	req := BuildHotelResRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, body)
	_, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling hotel resercation request", err)
	}
	//fmt.Printf("%s\n", b)
}
