package hotelws

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestHotelResSet(t *testing.T) {
	body := SetHotelResBody(1)
	body.NewPropertyResByRPH("12")
	body.NewGuaranteeRes("Testlast", "G", "MC", "2012-12", "1234567890")

	textPrefs := []string{"Tes1", "Test2", "Test3"}
	prefs := &SpecialPrefs{}
	prefs.AddSpecPrefWritConf(true)
	prefs.AddSpecPrefText(textPrefs)
	body.AddSpecialPrefs(prefs)

	b := body.OTAHotelResRQ
	if len(b.Hotel.SpecialPrefs.Text) != len(textPrefs) {
		t.Error("SpecialPrefs.Text is wrong")
	}
	if !b.Hotel.SpecialPrefs.WrittenConfirmation.Ind {
		t.Error("SpecialPrefs.WrittenConfirmation is wrong")
	}
	if b.Hotel.RoomType.NumberOfUnits != 1 {
		t.Error("RoomType.NumberOfUnits is wrong")
	}
	if b.Hotel.BasicPropertyRes.RPH != "12" {
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
	if ccinfo.PersonName.Surname.Val != "Testlast" {
		t.Error("PersonName.Last is wrong")
	}
}

func TestHotelResByHotel(t *testing.T) {
	body := SetHotelResBody(1)
	body.NewPropertyResByHotel("SL", "00004")
	b := body.OTAHotelResRQ
	if b.Hotel.BasicPropertyRes.ChainCode != "SL" {
		t.Error("BasicPropertyRes.ChainCode is wrong")
	}
	if b.Hotel.BasicPropertyRes.HotelCode != "00004" {
		t.Error("BasicPropertyRes.HotelCode is wrong")
	}
}

func TestHotelResBuild(t *testing.T) {
	body := SetHotelResBody(1)
	//body.NewPropertyResByHotel("SL", "00004")
	body.NewPropertyResByRPH("004")
	body.NewGuaranteeRes("Testlast", "GDPST", "MC", "2012-12", "1234567890")

	req := BuildHotelResRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, body)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling hotel resercation request", err)
	}
	fmt.Printf("%s\n", b)
}
