package hotelws

import (
	"encoding/xml"
	"testing"

	"github.com/ailgroup/sbrweb/sbrerr"
)

func TestPropDescValidReqCityCode(t *testing.T) {
	//city code make this fail
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	_, err := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	if err != sbrerr.ErrPropDescCityCode {
		t.Error("ErrPropDescCityCode should return")
	}
}

func TestPropDescValidReqLatLng(t *testing.T) {
	//lat,lng will make this fail
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	_, err := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	if err != sbrerr.ErrPropDescLatLng {
		t.Error("ErrPropDescLatLng should return")
	}
}

func TestPropDescValidMultiHotelRefs(t *testing.T) {
	//multiple ids (multiple criterion) will make this fail
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	_, err := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	if err != sbrerr.ErrPropDescHotelRefs {
		t.Error("ErrPropDescHotelRefs should return error")
	}
}

func TestPropDescBuildHotelPropDescMarshal(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[hotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, err := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	if err != nil {
		t.Error("Error SetHotelPropDescRqStruct: ", err)
	}
	req := BuildHotelPropDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, prop)

	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}

	if string(b) != string(samplePropRQIDs) {
		t.Errorf("Expected marshal SOAP hotel property description for hotel ids \n sample: %s \n result: %s", string(samplePropRQIDs), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestSetHotelPropDescRqStructCorpID(t *testing.T) {
	body, _ := SetHotelPropDescBody(sampleGuestCount, &HotelSearchCriteria{}, sampleArrive, sampleDepart)
	prop := body.HotelPropDescRQ
	prop.addCorporateID(sampleCID)

	if prop.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("SetHotelPropDescRqStruct Corporate.ID  expect: %s, got %s", sampleCID, prop.Avail.Customer.Corporate.ID)
	}

	prop.addCustomerID(sampleCID)
	if prop.Avail.Customer.CustomerID.Number != sampleCID {
		t.Errorf("SetHotelPropDescRqStruct CustomerID.Number  expect: %s, got %s", sampleCID, prop.Avail.Customer.Corporate.ID)
	}

}

func TestPropDescUnmarshal(t *testing.T) {
	prop := HotelPropDescResponse{}
	err := xml.Unmarshal(sampleHotelPropDescRSgood, &prop)
	if err != nil {
		t.Fatalf("Error unmarshaling hotel avail %s \nERROR: %v", sampleHotelAvailRSgood, err)
	}
	reqError := prop.Body.HotelDesc.Result.Error
	if reqError.Type != "" {
		t.Errorf("Request error %v should not have message %s", reqError, reqError.System.Message)
	}
	success := prop.Body.HotelDesc.Result.Success
	if success.System.HostCommand.LNIATA != "222222" {
		t.Errorf("System.HostCommand.LNIATA for success expect: %v, got: %v", "222222", success.System.HostCommand.LNIATA)
	}
	roomStayRates := prop.Body.HotelDesc.RoomStay.RoomRates
	numRates := len(roomStayRates)
	if numRates != 16 {
		t.Error("Number of rates is wrong")
	}

	rate0 := roomStayRates[0]
	sample0 := rateSamples[0]
	if rate0.RPH != sample0.rph {
		t.Errorf("RPH expected %s, got %s", sample0.rph, rate0.RPH)
	}
	if rate0.GuaranteedRate != sample0.guarrate {
		t.Errorf("GuaranteedRate expected %s, got %s", sample0.guarrate, rate0.GuaranteedRate)
	}
	if rate0.IATA_Product != sample0.iataprod {
		t.Errorf("IATA_Product expected %s, got %s", sample0.iataprod, rate0.IATA_Product)
	}
	if rate0.IATA_Character != sample0.iatachar {
		t.Errorf("IATA_Character expected %s, got %s", sample0.iatachar, rate0.IATA_Character)
	}
	rrate0 := rate0.Rates[0]
	samplerrate0 := sample0.rates[0]
	if rrate0.Amount != samplerrate0.Amount {
		t.Errorf("Rate expected %s, got %s", samplerrate0.Amount, rrate0.Amount)
	}
	if rrate0.CurrencyCode != samplerrate0.CurrencyCode {
		t.Errorf("Currency code expected %s, got %s", samplerrate0.CurrencyCode, rrate0.CurrencyCode)
	}
	xtra0 := rrate0.AdditionalGuestAmounts[0].Charges[0].ExtraPerson
	sampleXtra0 := samplerrate0.AdditionalGuestAmounts[0].Charges[0].ExtraPerson
	if xtra0 != sampleXtra0 {
		t.Errorf("ExtraPerson expected %s, got %s", sampleXtra0, xtra0)
	}
	hprice0 := rrate0.HotelPricing.Amount
	sampleHprice0 := samplerrate0.HotelPricing.Amount
	if hprice0 != sampleHprice0 {
		t.Errorf("HotelPricing expected %s, got %s", sampleHprice0, hprice0)
	}
	surcharge0 := rrate0.HotelPricing.TotalSurcharges.Amount
	sampleSurcharge0 := samplerrate0.HotelPricing.TotalSurcharges.Amount
	if surcharge0 != sampleSurcharge0 {
		t.Errorf("TotalSurcharges expected %s, got %s", sampleSurcharge0, surcharge0)
	}
	taxes0 := rrate0.HotelPricing.TotalTaxes.Amount
	sampleTaxes0 := samplerrate0.HotelPricing.TotalTaxes.Amount
	if taxes0 != sampleTaxes0 {
		t.Errorf("TotalTaxes expected %s, got %s", sampleTaxes0, taxes0)
	}
	//fmt.Printf("CURRENT: %+v\n", prop)
	//fmt.Printf("RATES COUNT: %d\n", len(prop.Body.HotelDesc.RoomStay.RoomRates))
}

func TestPropDescCall(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[hotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, prop)

	resp, err := CallHotelPropDesc(serverHotelPropertyDesc.URL, req)
	if err != nil {
		t.Error("Error making request CallHotelProperty", err)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}
	roomStayRates := resp.Body.HotelDesc.RoomStay.RoomRates
	numRoomRates := len(roomStayRates)
	if numRoomRates != 16 {
		t.Error("Number of rates is wrong")
	}

	for idx, rr := range roomStayRates {
		if rr.IATA_Character != iataCharSample[idx] {
			t.Errorf("IATA_Character %d expected %s, got %s", idx, iataCharSample[idx], rr.IATA_Character)
		}
		if rr.GuaranteeSurcharge != "G" {
			t.Errorf("GuaranteeSurcharge %d expected %s, got %s", idx, "G", rr.GuaranteeSurcharge)
		}
		cnum := rr.AdditionalInfo.CancelPolicy.Numeric
		copt := rr.AdditionalInfo.CancelPolicy.Option
		if cnum != 2 {
			t.Errorf("RoomRate %d expected cancel policy numeric %d, got %d", idx, 2, cnum)
		}
		if copt != "D" {
			t.Errorf("RoomRate %d expected cancel policy option %s, got %s", idx, "D", copt)
		}
	}

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
}

func TestHotelPropDescCallDown(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[hotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, prop)

	resp, err := CallHotelPropDesc(serverHotelDown.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverHotelDown")
	}
	if err.Error() != resp.ErrorSabreService.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreService.Code != sbrerr.BadService {
		t.Errorf("Expect %d got %d", sbrerr.BadService, resp.ErrorSabreService.Code)
	}
	if resp.ErrorSabreService.AppMessage != sbrerr.ErrCallHotelPropDesc {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallHotelPropDesc, resp.ErrorSabreService.AppMessage)
	}
}

func TestHotelPropDescCallBadResponseBody(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[hotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(samplesite, samplepcc, samplebinsectoken, sampleconvid, samplemid, sampletime, prop)

	resp, err := CallHotelPropDesc(serverBadBody.URL, req)
	if err == nil {
		t.Error("Expected error making request to sserverBadBody")
	}
	if err.Error() != resp.ErrorSabreXML.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreXML.Code != sbrerr.BadParse {
		t.Errorf("Expect %d got %d", sbrerr.BadParse, resp.ErrorSabreXML.Code)
	}
	if resp.ErrorSabreXML.AppMessage != sbrerr.ErrCallHotelPropDesc {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallHotelPropDesc, resp.ErrorSabreXML.AppMessage)
	}
}
