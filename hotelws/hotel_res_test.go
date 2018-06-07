package hotelws

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestHotelResSet(t *testing.T) {
	body := SetHotelResBody(
		2,
		TimeSpanFormatter("04-22", "04-25", TimeFormatMD, TimeFormatMDTHM),
	)
	body.NewPropertyResByRPH(12)
	body.NewGuaranteeRes("Testlast", "G", "MC", "2012-12", "1234567890")
	b := body.OTAHotelResRQ
	if b.Hotel.RoomType.NumberOfUnits != 0 {
		t.Error("RoomType.NumberOfUnits is wrong")
	}
	if b.Hotel.RoomType.RoomTypeCode != "" {
		t.Error("RoomType.RoomTypeCode is wrong")
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
	if ccinfo.PersonName.Surname.Val != "Testlast" {
		t.Error("PersonName.Last is wrong")
	}
}

func TestHotelResBuild(t *testing.T) {
	body := SetHotelResBody(
		2,
		TimeSpanFormatter("04-22", "04-25", TimeFormatMD, TimeFormatMDTHM),
	)
	body.NewPropertyResByHotel("SL", "00004")
	body.NewGuaranteeRes("Testlast", "GDPST", "MC", "2012-12", "1234567890")
	body.AddRoomType(1, "ABC123")

	//prefs := &SpecialPrefs{}
	//prefs.AddSpecPrefWritConf(true)
	//prefs.AddSpecPrefText([]string{"Tes1", "Test2", "Test3"})
	//body.AddSpecialPrefs(prefs)

	req := BuildHotelResRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, body)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling hotel resercation request", err)
	}
	fmt.Printf("%s\n", b)
}
