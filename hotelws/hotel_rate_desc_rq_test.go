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

func TestRateDescCall(t *testing.T) {
	// assume RPH is from previous hotel property description call
	rpc := SetRateParams(
		[]RatePlan{
			RatePlan{
				RPH: 12,
			},
		},
	)
	rate, _ := SetHotelRateDescRqStruct(rpc)

	req := BuildHotelRateDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, rate)

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
		t.Errorf("Number of room rates is wrong, expect %d got %d", 1, numRoomRates)
	}

	rr := roomStayRates[0]
	if rr.IATA_Character != "J1KA16" {
		t.Errorf("IATA_Character expected %s, got %s", "J1KA16", rr.IATA_Character)
	}
	if rr.GuaranteeSurcharge != "G" {
		t.Errorf("GuaranteeSurcharge expected %s, got %s", "G", rr.GuaranteeSurcharge)
	}
	cnum := rr.AdditionalInfo.CancelPolicy.Numeric
	copt := rr.AdditionalInfo.CancelPolicy.Option
	if cnum != 2 {
		t.Errorf("RoomRate expected cancel policy numeric %d, got %d", 2, cnum)
	}
	if copt != "D" {
		t.Errorf("RoomRate expected cancel policy option %s, got %s", "D", copt)
	}

	/*
		indexRoomRate := numRoomRates - 1
		numRates := len(roomStayRates[indexRoomRate].Rates)
		if numRates != 1 {
			t.Error("Number of rates is wrong")
		}
		for _, rate := range roomStayRates[indexRoomRate].Rates {
			if rate.Amount != "400.00" {
				t.Errorf("Rate expected %s, got %s", "400.00", rate.Amount)
			}
			if rate.CurrencyCode != "SGD" {
				t.Errorf("CurrencyCode expected %s, got %s", "SGD", rate.CurrencyCode)
			}
			if rate.HRD_RequiredForSell != "false" {
				t.Errorf("CurrencyCode expected %s, got %s", "false", rate.HRD_RequiredForSell)
			}

			hprice := rate.HotelPricing
			if hprice.Amount != "470.80" {
				t.Errorf("HotelPricing expected %s, got %s", "470.80", hprice.Amount)
			}
			if hprice.TotalSurcharges.Amount != "40.00" {
				t.Errorf("TotalSurcharges expected %s, got %s", "40.00", hprice.TotalSurcharges.Amount)
			}
			if hprice.TotalTaxes.Amount != "30.80" {
				t.Errorf("TotalTaxes expected %s, got %s", "30.80", hprice.TotalTaxes.Amount)
			}
		}

	*/
}
