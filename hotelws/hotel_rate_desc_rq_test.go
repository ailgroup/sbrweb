package hotelws

import (
	"encoding/xml"
	"testing"
)

func TestHotelRateDescMarshal(t *testing.T) {
	rpc := SetRateParams(
		[]RatePlan{
			RatePlan{
				RPH: 12,
			},
		},
	)
	rate, err := SetHotelRateDescRqStruct(rpc)
	if err != nil {
		t.Error("Error SetHotelRateDescRqStruct:", err)
	}
	req := BuildHotelRateDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, rate)

	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}

	if string(b) != string(sampleHotelRateDescRQRPH) {
		t.Errorf("Expected marshal SOAP hotel rate description for rph \n sample: %s \n result: %s", string(sampleHotelRateDescRQRPH), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

var ccards = []string{"DS", "CA", "MC", "CB", "VI", "VS", "AX", "JC", "DC"}

func TestRateDescCall(t *testing.T) {
	// assume RPH is from previous hotel property description call
	rpc := SetRateParams(
		[]RatePlan{
			RatePlan{
				RPH: 12,
			},
		},
	)
	raterq, _ := SetHotelRateDescRqStruct(rpc)

	req := BuildHotelRateDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, raterq)

	resp, err := CallHotelRateDesc(serverHotelRateDesc.URL, req)
	if err != nil {
		t.Error("Error making request CallHotelRateDesc", err)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
	roomStayRates := resp.Body.HotelDesc.RoomStay.RoomRates
	numRoomRates := len(roomStayRates)
	if numRoomRates != 1 {
		t.Error("Number of room rates is wrong")
	}

	rr := roomStayRates[0]
	if rr.IATA_Character != "J1KA16" {
		t.Errorf("IATA_Character expected %s, got %s", "J1KA16", rr.IATA_Character)
	}
	if rr.GuaranteeSurcharge != "G" {
		t.Errorf("GuaranteeSurcharge expected %s, got %s", "G", rr.GuaranteeSurcharge)
	}
	if len(rr.AdditionalInfo.PaymentCard) != 9 {
		t.Errorf("AdditionalInfo.PaymentCard count is wrong: %v", rr.AdditionalInfo.PaymentCard)
	}
	for idx, card := range rr.AdditionalInfo.PaymentCard {
		if card.Code != ccards[idx] {
			t.Errorf("AdditionalInfo.PaymentCard expect: %s, got: %s", ccards[idx], card.Code)
		}
	}

	cnum := rr.AdditionalInfo.CancelPolicy.Numeric
	copt := rr.AdditionalInfo.CancelPolicy.Option
	if cnum != 2 {
		t.Errorf("RoomRate expected cancel policy numeric %d, got %d", 2, cnum)
	}
	if copt != "D" {
		t.Errorf("RoomRate expected cancel policy option %s, got %s", "D", copt)
	}

	numRates := len(roomStayRates[0].Rates)
	if numRates != 1 {
		t.Error("Number of rates is wrong")
	}
	rate := rr.Rates[0]
	if rate.Amount != "274.55" {
		t.Errorf("Rate expected %s, got %s", "274.55", rate.Amount)
	}
	if rate.CurrencyCode != "USD" {
		t.Errorf("CurrencyCode expected %s, got %s", "USD", rate.CurrencyCode)
	}
	if rate.HRD_RequiredForSell != "false" {
		t.Errorf("CurrencyCode expected %s, got %s", "false", rate.HRD_RequiredForSell)
	}

	hprice := rate.HotelPricing
	if hprice.Amount != "307.50" {
		t.Errorf("HotelPricing expected empty %s, got %s", "307.50", hprice.Amount)
	}
	if hprice.TotalSurcharges.Amount != "" {
		t.Errorf("TotalSurcharges expected %s, got %s", "", hprice.TotalSurcharges.Amount)
	}
	if hprice.TotalTaxes.Amount != "32.95" {
		t.Errorf("TotalTaxes expected %s, got %s", "32.95", hprice.TotalTaxes.Amount)
	}

	if hprice.TotalTaxes.TaxFieldOne != "19.22" {
		t.Errorf("TaxeFieldOnw expected %s, got %s", "19.22", hprice.TotalTaxes.TaxFieldOne)
	}
	if hprice.TotalTaxes.TaxFieldTwo != "13.73" {
		t.Errorf("TaxFieldTwo expected %s, got %s", "13.73", hprice.TotalTaxes.TaxFieldTwo)
	}

}
