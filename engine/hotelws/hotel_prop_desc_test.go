package hotelws

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
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
	hotelid[HotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, err := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	if err != nil {
		t.Error("Error SetHotelPropDescRqStruct: ", err)
	}
	req := BuildHotelPropDescRequest(sconf, prop)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}

	if string(b) != string(samplePropRQIDs) {
		t.Errorf("Expected marshal SOAP hotel property description for hotel ids \n sample: %s \n result: %s", string(samplePropRQIDs), string(b))
	}
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
}

func TestPropDescCall(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[HotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(sconf, prop)
	resp, err := CallHotelPropDesc(serverHotelPropertyDesc.URL, req)
	if err != nil {
		t.Error("Error making request CallHotelProperty", err)
	}

	if !resp.Body.HotelDesc.Result.Ok() {
		t.Error("CallHotelPropDesc Ok should be true")
	}
	sabreErrFmt := resp.Body.HotelDesc.Result.ErrFormat()
	if sabreErrFmt.Code.String() != "Complete" {
		t.Errorf("CallHotelPropDesc code expected %s, got %s", "Complete", sabreErrFmt.Code.String())
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

func TestNewParsedRoomMeta(t *testing.T) {
	//NotUrlSafeString := "some data with \x00 and \ufeff"
	//NotUrlSafeStringExpected := "some data with  and "
	errExpect := errors.New("illegal base64 data at input byte 31")
	b64NotUrlSafe := "c29tZSBkYXRhIHdpdGggACBhbmQg77u/"
	_, err := NewParsedRoomMeta(b64NotUrlSafe)
	if err == nil {
		t.Errorf("NewParsedRoomMeta expected error")
	}
	if err.Error() != errExpect.Error() {
		t.Errorf("NewParsedRoomMeta expected error %v, got %v", errExpect, err)
	}
}

var proptrack = []struct {
	b64str string
	expect ParsedRoomMeta
}{
	{"YXJ2OjA0LTAyfGRwdDowNC0wNXxnc3Q6MnxoYzpIT0Q0LzExTUFZLTEyTUFZMnxoaWQ6MTB8cnBoOjAwMXxybXQ6UDFLUkFDfFtjdXI6U0dELXJxczpmYWxzZS1hbXQ6MzM1LjQ1O2N1cjpVU0QtcnFzOmZhbHNlLWFtdDo0MzUuNDVd", ParsedRoomMeta{
		Arrive:         "04-02",
		Depart:         "04-05",
		Guest:          "2",
		Hc:             "HOD4/11MAY-12MAY2",
		HotelID:        "10",
		Rmt:            "P1KRAC",
		Rph:            "001",
		StayRatesCache: []string{"cur:SGD-rqs:false-amt:335.45", "cur:USD-rqs:false-amt:435.45"},
		ParsedStayRatesCache: []parsedStayRateCache{
			parsedStayRateCache{
				Amt: "335.45",
				Cur: "SGD",
				Rqs: false,
			},
			parsedStayRateCache{
				Amt: "435.45",
				Cur: "USD",
				Rqs: false,
			},
		},
	}},
	{"YXJ2OjA0LTAyfGRwdDowNC0wNXxnc3Q6MnxoYzpIT0Q0LzExTUFZLTEyTUFZMnxoaWQ6MTB8cnBoOjAwMnxybXQ6RDFLUkFDfFtjdXI6U0dELXJxczpmYWxzZS1hbXQ6MTkwLjQ1XQ==",
		ParsedRoomMeta{
			Arrive:         "04-02",
			Depart:         "04-05",
			Guest:          "2",
			Hc:             "HOD4/11MAY-12MAY2",
			HotelID:        "10",
			Rmt:            "D1KRAC",
			Rph:            "002",
			StayRatesCache: []string{"rtx:0-ttl:190.45-nxt:false"},
			ParsedStayRatesCache: []parsedStayRateCache{
				parsedStayRateCache{
					Amt: "190.45",
					Cur: "SGD",
					Rqs: false,
				},
			},
		}},
}

func TestSetRoomMetaDataPropDesc(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[HotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(sconf, prop)
	resp, _ := CallHotelPropDesc(serverHotelPropertyDesc.URL, req)
	resp.SetRoomMetaData(sampleGuestCount, sampleArrive, sampleDepart, "10")
	for i, rate := range resp.Body.HotelDesc.RoomStay.RoomRates {
		//only test the first 2
		if i > 1 {
			break
		}
		if rate.B64RoomMetaData != proptrack[i].b64str {
			t.Errorf("B64RoomMetaData expect: '%s', got '%s'", proptrack[i].b64str, rate.B64RoomMetaData)
		}
		prm, err := NewParsedRoomMeta(rate.B64RoomMetaData)
		if err != nil {
			t.Errorf("Error on DecodeTrackedEncoding() %v", err)
		}
		if prm.Arrive != proptrack[i].expect.Arrive {
			t.Errorf("Arrive expect %s, got %s", proptrack[i].expect.Arrive, prm.Arrive)
		}
		if prm.Depart != proptrack[i].expect.Depart {
			t.Errorf("Depart expect %s, got %s", proptrack[i].expect.Depart, prm.Depart)
		}
		if prm.Guest != proptrack[i].expect.Guest {
			t.Errorf("Guest expect %s, got %s", proptrack[i].expect.Guest, prm.Guest)
		}
		if prm.Hc != proptrack[i].expect.Hc {
			t.Errorf("Hc expect %s, got %s", proptrack[i].expect.Hc, prm.Hc)
		}
		if prm.HotelID != proptrack[i].expect.HotelID {
			t.Errorf("HotelID expect %s, got %s", proptrack[i].expect.HotelID, prm.HotelID)
		}

		for idx, psrc := range prm.ParsedStayRatesCache {
			amt := proptrack[i].expect.ParsedStayRatesCache[idx].Amt
			cur := proptrack[i].expect.ParsedStayRatesCache[idx].Cur
			rqs := proptrack[i].expect.ParsedStayRatesCache[idx].Rqs
			if psrc.Amt != amt {
				t.Errorf("Amt expect %s, got %s", amt, psrc.Amt)
			}
			if psrc.Cur != cur {
				t.Errorf("Cur expect %s, got %s", cur, psrc.Cur)
			}
			if psrc.Rqs != rqs {
				t.Errorf("Amt expect %v, got %v", rqs, psrc.Rqs)
			}
		}

	}
}

func TestHotelPropDescCallDown(t *testing.T) {
	var hotelid = make(HotelRefCriterion)
	hotelid[HotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(sconf, prop)
	resp, err := CallHotelPropDesc(serverHotelDown.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverHotelDown")
	}
	if !resp.Body.HotelDesc.Result.Ok() {
		t.Error("CallHotelPropDesc Ok should be true")
	}
	sabreErrFmt := resp.Body.HotelDesc.Result.ErrFormat()
	if sabreErrFmt.Code.String() != "Unknown" {
		t.Errorf("CallHotelPropDesc code expected %s, got %s", "Unknown", sabreErrFmt.Code.String())
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
	hotelid[HotelidQueryField] = []string{"10"}
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hotelid),
	)
	prop, _ := SetHotelPropDescBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelPropDescRequest(sconf, prop)
	resp, err := CallHotelPropDesc(serverBadBody.URL, req)
	if err == nil {
		t.Error("Expected error making request to sserverBadBody")
	}
	if !resp.Body.HotelDesc.Result.Ok() {
		t.Error("CallHotelPropDesc Ok should be true")
	}
	sabreErrFmt := resp.Body.HotelDesc.Result.ErrFormat()
	if sabreErrFmt.Code.String() != "Unknown" {
		t.Errorf("CallHotelPropDesc code expected %s, got %s", "Unknown", sabreErrFmt.Code.String())
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
