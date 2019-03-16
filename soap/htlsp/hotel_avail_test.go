package htlsp

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/ailgroup/sbrweb/sbrerr"
	"github.com/ailgroup/sbrweb/soap/srvc"
)

func TestAddressSearchReturnError(t *testing.T) {
	_, err := NewHotelSearchCriteria(
		AddressSearch(AddressCriterion{}),
	)
	if err == nil {
		t.Errorf("AddressSearch empty params should return error: '%v'", err)
	}
}

func TestAddressSearchCriteria(t *testing.T) {
	a, err := NewHotelSearchCriteria(
		AddressSearch(addr),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with AddressSearch error %v", err)
	}
	if a.Criterion.AddressSearch.StreetNmbr != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if a.Criterion.AddressSearch.CityName != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if a.Criterion.AddressSearch.PostalCode != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if a.Criterion.AddressSearch.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
	}
}

func TestPackageSearchCriteria(t *testing.T) {
	p, err := NewHotelSearchCriteria(
		PackageSearch(samplePackages),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with PackageSearch error %v", err)
	}
	if len(p.Criterion.Packages) != len(samplePackages) {
		t.Errorf("PackageSearch wrong number of results, expected %d got %d", len(samplePackages), len(p.Criterion.Packages))
	}
}

func TestPropertyTypeSearchCriteria(t *testing.T) {
	p, err := NewHotelSearchCriteria(
		PropertyTypeSearch(samplePropertyTypes),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with PropertyTypeSearch error %v", err)
	}
	if len(p.Criterion.PropertyTypes) != len(samplePropertyTypes) {
		t.Errorf("PackageSearch wrong number of results, expected %d got %d", len(samplePropertyTypes), len(p.Criterion.PropertyTypes))
	}
}

func TestPointOfInterestSearch(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		PointOfInterestSearch(pointOfInt),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with PointOfInterestSearch error %v", err)
	}
	if r.Criterion.PointOfInterest.NonUS {
		t.Error("NonUS should be false")
	}
	if r.Criterion.PointOfInterest.CountryStateCode != sampleState {
		t.Errorf("State expected: %s, got: %s", sampleState, r.Criterion.PointOfInterest.CountryStateCode)
	}
	if r.Criterion.PointOfInterest.Val != sampleCity {
		t.Errorf("POI Val expected: %s, got: %s", sampleCity, r.Criterion.PointOfInterest.Val)
	}
}

func TestHotelRefSearchReturnError(t *testing.T) {
	_, err := NewHotelSearchCriteria(
		HotelRefSearch(hqbad),
	)
	if err == nil {
		t.Errorf("HotelRefSearch empty params should return error: '%v'", err)
	}
}

func TestHotelRefSearchCityCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleHotelCityCode {
		if r.Criterion.HotelRefs[i].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", i, code, r.Criterion.HotelRefs[i].HotelCityCode)
		}
	}
}

func TestHotelRefSearchHotelCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleHotelCode {
		if r.Criterion.HotelRefs[i].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", i, code, r.Criterion.HotelRefs[i].HotelCode)
		}

	}
}

func TestHotelRefSearchLatLngCodeCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with HotelRefSearch error %v", err)
	}
	for i, code := range sampleLatLang {
		ll := strings.Split(code, ",")
		if r.Criterion.HotelRefs[i].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", i, ll[0], r.Criterion.HotelRefs[i].Latitude)
		}
		if r.Criterion.HotelRefs[i].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", i, ll[1], r.Criterion.HotelRefs[i].Longitude)
		}
	}
}

func TestMultipleHotelCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
		HotelRefSearch(hqcity),
		AddressSearch(addr),
		HotelRefSearch(hqltln),
		PackageSearch(samplePackages),
		PropertyTypeSearch(samplePropertyTypes),
	)
	if err != nil {
		t.Errorf("NewHotelSearchCriteria with basic criteria error %v", err)
	}
	counter := 0
	for _, code := range sampleHotelCode {
		if r.Criterion.HotelRefs[counter].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", counter, code, r.Criterion.HotelRefs[counter].HotelCode)
		}
		counter++
	}
	for _, code := range sampleHotelCityCode {
		if r.Criterion.HotelRefs[counter].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", counter, code, r.Criterion.HotelRefs[counter].HotelCityCode)
		}
		counter++
	}
	if r.Criterion.AddressSearch.StreetNmbr != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if r.Criterion.AddressSearch.CityName != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if r.Criterion.AddressSearch.PostalCode != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if r.Criterion.AddressSearch.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
	}

	for _, code := range sampleLatLang {
		ll := strings.Split(code, ",")
		if r.Criterion.HotelRefs[counter].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", counter, ll[0], r.Criterion.HotelRefs[counter].Latitude)
		}
		if r.Criterion.HotelRefs[counter].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", counter, ll[1], r.Criterion.HotelRefs[counter].Longitude)
		}
		counter++
	}
}

func TestSetHotelAvailRqStructMarshal(t *testing.T) {
	availBody := SetHotelAvailBody(sampleGuestCount, &HotelSearchCriteria{}, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)
	if avail.XMLNSXsi != srvc.BaseXSINamespace {
		t.Errorf("SetHotelAvailRqStruct XMLNSXsi expect: %s, got %s", srvc.BaseXSINamespace, avail.XMLNSXsi)
	}
	if avail.Version != hotelRQVersion {
		t.Errorf("SetHotelAvailRqStruct Version expect: %s, got %s", hotelRQVersion, avail.Version)
	}
	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("SetHotelAvailRqStruct GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("SetHotelAvailRqStruct Customer.Corporate.ID expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}
	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestSetHotelAvailRqStructCorpID(t *testing.T) {
	availBody := SetHotelAvailBody(sampleGuestCount, &HotelSearchCriteria{}, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("SetHotelAvailRqStruct Corporate.ID  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}
	avail.addCustomerID(sampleCID)
	if avail.Avail.Customer.CustomerID.Number != sampleCID {
		t.Errorf("SetHotelAvailRqStruct CustomerID.Number  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}
}

func TestAvailIdsMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	gcount := 4
	availBody := SetHotelAvailBody(gcount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCorporateID(sampleCID)
	avail.Avail.RatePlanCandidates = SetRatePlanCandidate([]RatePlan{RatePlan{CurrencyCode: "USD", DCA_ProductCode: "I7A"}})
	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("SetHotelAvailRqStruct GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}
	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqids[HotelidQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqids[HotelidQueryField]))
	}
	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQHotelIDSCoprIDRatePlans) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDSCoprIDRatePlans), string(b))
	}
}

func TestAvailCitiesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	gcount := 3
	availBody := SetHotelAvailBody(gcount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	avail.addCustomerID(sampleCID)
	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}
	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqcity[CityQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqcity[CityQueryField]))
	}
	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQCitiesCustNumber) {
		t.Errorf("Expected marshal hotel avail for hotel cities \n sample: %s \n result: %s", string(sampleAvailRQCitiesCustNumber), string(b))
	}
}

func TestAvailLatLngMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	availBody := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}
	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs) != len(hqltln[LatlngQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRefs), len(hqltln[LatlngQueryField]))
	}
	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQLatLng) {
		t.Errorf("Expected marshal set hotel avail for hotel lat/lng \n sample: %s \n result: %s", string(sampleAvailRQLatLng), string(b))
	}
}

func TestAvailPropertyTypesPackagesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		PackageSearch(samplePackages),
		PropertyTypeSearch(samplePropertyTypes),
	)
	availBody := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	avail := availBody.OTAHotelAvailRQ
	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQPropPackages) {
		t.Errorf("Expected marshal set hotel avail for hotel packages and property types \n sample: %s \n result: %s", string(sampleAvailRQLatLng), string(b))
	}
}

func TestBuildHotelAvailRequestMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	avail := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelAvailRequest(sconf, avail)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}

	if string(b) != string(sampleAvailRQHotelIDS) {
		t.Errorf("Expected marshal SOAP hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDS), string(b))
	}
}

func TestHotelAvailUnmarshal(t *testing.T) {
	avail := HotelAvailResponse{}
	err := xml.Unmarshal(sampleHotelAvailRSgood, &avail)
	if err != nil {
		t.Errorf("Error unmarshaling hotel avail %s \nERROR: %v", sampleHotelAvailRSgood, err)
	}
	reqError := avail.Body.HotelAvail.Result.Error
	if reqError.Type != "" {
		t.Errorf("Request error %v should not have message %s", reqError, reqError.System.Message)
	}
	success := avail.Body.HotelAvail.Result.Success
	if success.System.HostCommand.LNIATA != "222222" {
		t.Errorf("System.HostCommand.LNIATA for success expect: %v, got: %v", "222222", success.System.HostCommand.LNIATA)
	}

	options := avail.Body.HotelAvail.AvailOpts.AvailableOptions[0]
	if options.RPH != "001" {
		t.Errorf("First Availability option should be 1")
	}
	rr := options.PropertyInfo.RoomRateAvail
	if rr.RateLevelCode != "RAC" {
		t.Errorf("RateLevelCode should be: %s, got: %s", "RAC", rr.RateLevelCode)
	}
	if rr.HotelRateCode != "RAC" {
		t.Errorf("HotelRateCode should be: %s, got: %s", "RAC", rr.HotelRateCode)
	}

	rateRange := options.PropertyInfo.RateRange
	if rateRange.CurrencyCode != "USD" {
		t.Errorf("RateRange CurrencyCode should be %s, got %s", "USD", rateRange.CurrencyCode)
	}
	if rateRange.Max != "289.00" {
		t.Errorf("RateRange Max should be %s, got %s", "USD", rateRange.Max)
	}
	if rateRange.Min != "134.00" {
		t.Errorf("RateRange Min should be %s, got %s", "USD", rateRange.Min)
	}
}

func TestHotelAvailCallByIDs(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	avail := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelAvailRequest(sconf, avail)
	resp, err := CallHotelAvail(serverHotelAvailability.URL, req)
	if err != nil {
		t.Error("Error making request CallHotelAvail", err)
	}
	if resp.Body.Fault.String != "" {
		t.Errorf("Body.Fault.String expect empty: '%s', got: %s", "", resp.Body.Fault.String)
	}

	for idx, o := range resp.Body.HotelAvail.AvailOpts.AvailableOptions {
		if o.RPH != fmt.Sprintf("00%d", idx+1) {
			t.Errorf("AvailableOptions %d RPH expected %d, got %s", idx, idx+1, o.RPH)
		}
		if o.PropertyInfo.HotelCityCode != "TUL" {
			t.Errorf("AvailableOptions %d HotelCityCode expected %s, got %s", idx, "TUL", o.PropertyInfo.HotelCityCode)
		}

	}
}

func TestHotelAvailCallDown(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	avail := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelAvailRequest(sconf, avail)
	resp, err := CallHotelAvail(serverHotelDown.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverHotelDown")
	}
	if err.Error() != resp.ErrorSabreService.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreService.Code != sbrerr.BadService {
		t.Errorf("Expect %d got %d", sbrerr.BadService, resp.ErrorSabreService.Code)
	}
	if resp.ErrorSabreService.AppMessage != sbrerr.ErrCallHotelAvail {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallHotelAvail, resp.ErrorSabreService.AppMessage)
	}
}

func TestHotelAvailCallBadResponseBody(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	avail := SetHotelAvailBody(sampleGuestCount, q, sampleArrive, sampleDepart)
	req := BuildHotelAvailRequest(sconf, avail)
	resp, err := CallHotelAvail(serverBadBody.URL, req)
	if err == nil {
		t.Error("Expected error making request to sserverBadBody")
	}
	if err.Error() != resp.ErrorSabreXML.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreXML.Code != sbrerr.BadParse {
		t.Errorf("Expect %d got %d", sbrerr.BadParse, resp.ErrorSabreXML.Code)
	}
	if resp.ErrorSabreXML.AppMessage != sbrerr.ErrCallHotelAvail {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallHotelAvail, resp.ErrorSabreXML.AppMessage)
	}
}
