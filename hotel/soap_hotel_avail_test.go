package sbrhotel

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/ailgroup/sbrweb"
)

var (
	hqbad               = make(HotelRefCriterion)
	hqcity              = make(HotelRefCriterion)
	hqids               = make(HotelRefCriterion)
	hqltln              = make(HotelRefCriterion)
	addr                = make(AddressCriterion)
	sampleCID           = "12345"
	sampleLatLang       = []string{"32.78,-96.81", "54.87,-102.96"}
	sampleHotelCode     = []string{"0012", "19876", "1109", "445098", "000034"}
	sampleHotelCityCode = []string{"DFW", "CHC", "LA"}
	sampleGuestCount    = 2
	sampleStreet        = "2031 N. 100 W"
	sampleCity          = "Nowhere"
	samplePostal        = "999908"
	sampleCountryCode   = "US"
	sampleArrive        = time.Now().Add(48 * time.Hour)
	sampleDepart        = sampleArrive.Add(72 * time.Hour)

	sampleAvailRQHotelIDSCoprID = []byte(`<OTA_HotelAvailRQ version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><ReturnHostCommand>true</ReturnHostCommand><AvailRequestSegment><Customer><Corporate><ID>12345</ID></Corporate></Customer><GuestCounts Count="4"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCode="0012"></HotelRef><HotelRef HotelCode="19876"></HotelRef><HotelRef HotelCode="1109"></HotelRef><HotelRef HotelCode="445098"></HotelRef><HotelRef HotelCode="000034"></HotelRef><Address></Address></Criterion></HotelSearchCriteria><TimeSpan><End>04-05</End><Start>04-02</Start></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQCitiesCustNumber = []byte(`<OTA_HotelAvailRQ version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><ReturnHostCommand>true</ReturnHostCommand><AvailRequestSegment><Customer><ID><Number>12345</Number></ID></Customer><GuestCounts Count="3"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef HotelCityCode="DFW"></HotelRef><HotelRef HotelCityCode="CHC"></HotelRef><HotelRef HotelCityCode="LA"></HotelRef><Address></Address></Criterion></HotelSearchCriteria><TimeSpan><End>04-05</End><Start>04-02</Start></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)

	sampleAvailRQLatLng = []byte(`<OTA_HotelAvailRQ version="2.3.0" xmlns="http://webservices.sabre.com/sabreXML/2011/10" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><ReturnHostCommand>true</ReturnHostCommand><AvailRequestSegment><GuestCounts Count="2"></GuestCounts><HotelSearchCriteria><Criterion><HotelRef Latitude="32.78" Longitude="-96.81"></HotelRef><HotelRef Latitude="54.87" Longitude="-102.96"></HotelRef><Address></Address></Criterion></HotelSearchCriteria><TimeSpan><End>04-05</End><Start>04-02</Start></TimeSpan></AvailRequestSegment></OTA_HotelAvailRQ>`)
)

func init() {
	hqcity[cityQueryField] = sampleHotelCityCode
	hqids[hotelidQueryField] = sampleHotelCode
	hqltln[latlngQueryField] = sampleLatLang

	addr[streetQueryField] = sampleStreet
	addr[cityQueryField] = sampleCity
	addr[postalQueryField] = samplePostal
	addr[countryCodeQueryField] = sampleCountryCode
}

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
		t.Errorf("NewHotelSearchCriteria with AddressOption error %v", err)
	}
	if a.Criterion.Address.Street != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if a.Criterion.Address.City != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if a.Criterion.Address.Postal != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if a.Criterion.Address.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
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
		if r.Criterion.HotelRef[i].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", i, code, r.Criterion.HotelRef[i].HotelCityCode)
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
		if r.Criterion.HotelRef[i].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", i, code, r.Criterion.HotelRef[i].HotelCode)
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
		if r.Criterion.HotelRef[i].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", i, ll[0], r.Criterion.HotelRef[i].Latitude)
		}
		if r.Criterion.HotelRef[i].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", i, ll[1], r.Criterion.HotelRef[i].Longitude)
		}
	}
}

func TestMultipleCriteriaCriteria(t *testing.T) {
	r, err := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
		HotelRefSearch(hqcity),
		AddressSearch(addr),
		HotelRefSearch(hqltln),
	)

	if err != nil {
		t.Errorf("NewHotelSearchCriteria with all criteria error %v", err)
	}

	counter := 0
	for _, code := range sampleHotelCode {
		if r.Criterion.HotelRef[counter].HotelCode != code {
			t.Errorf("HotelRef[%d].HotelCode expect: %s, got: %s", counter, code, r.Criterion.HotelRef[counter].HotelCode)
		}
		counter++
	}
	for _, code := range sampleHotelCityCode {
		if r.Criterion.HotelRef[counter].HotelCityCode != code {
			t.Errorf("HotelRef[%d].HotelCityCode city expect: %s, got: %s", counter, code, r.Criterion.HotelRef[counter].HotelCityCode)
		}
		counter++
	}

	if r.Criterion.Address.Street != sampleStreet {
		t.Error("buildAddress street not correct")
	}
	if r.Criterion.Address.City != sampleCity {
		t.Error("buildAddress city not correct")
	}
	if r.Criterion.Address.Postal != samplePostal {
		t.Error("buildAddress postal not correct")
	}
	if r.Criterion.Address.CountryCode != sampleCountryCode {
		t.Error("buildAddress country code not correct")
	}

	for _, code := range sampleLatLang {
		ll := strings.Split(code, ",")
		if r.Criterion.HotelRef[counter].Latitude != ll[0] {
			t.Errorf("HotelRef[%d].Latitude expect: %s, got: %s", counter, ll[0], r.Criterion.HotelRef[counter].Latitude)
		}
		if r.Criterion.HotelRef[counter].Longitude != ll[1] {
			t.Errorf("HotelRef[%d].Longitude expect: %s, got: %s", counter, ll[1], r.Criterion.HotelRef[counter].Longitude)
		}
		counter++
	}
	/*
		avail := BuildHotelAvailRq(sampleCorpID, sampleGuestCount, r)
		b, err := xml.Marshal(avail)
		if err != nil {
			t.Error("Error marshaling get hotel content", err)
		}
		fmt.Printf("\n%s\n", b)
	*/

}

func TestBuildHotelSearchMarshal(t *testing.T) {
	avail := BuildHotelAvailRq(sampleGuestCount, HotelSearchCriteria{}, sampleArrive, sampleDepart)
	avail.addCorporateID(sampleCID)

	if avail.XMLNSXsi != sbrweb.BaseXSINamespace {
		t.Errorf("BuildHotelAvailRq XMLNSXsi expect: %s, got %s", sbrweb.BaseXSINamespace, avail.XMLNSXsi)
	}
	if avail.Version != hotelAvailVersion {
		t.Errorf("BuildHotelAvailRq Version expect: %s, got %s", hotelAvailVersion, avail.Version)
	}
	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("BuildHotelAvailRq Customer.Corporate.ID expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestBuildHotelSearchCorpID(t *testing.T) {
	avail := BuildHotelAvailRq(sampleGuestCount, HotelSearchCriteria{}, sampleArrive, sampleDepart)

	avail.addCorporateID(sampleCID)
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("BuildHotelAvailRq Customer.Corporate.ID  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

	avail.addCorporateID(sampleCID)
	if avail.Avail.Customer.Corporate.ID != sampleCID {
		t.Errorf("BuildHotelAvailRq Customer.Corporate.ID  expect: %s, got %s", sampleCID, avail.Avail.Customer.Corporate.ID)
	}

}

func TestBuildHotelSearchWithIDSMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	gcount := 4
	avail := BuildHotelAvailRq(gcount, q, sampleArrive, sampleDepart)
	avail.addCorporateID(sampleCID)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef) != len(hqids[hotelidQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef), len(hqids[hotelidQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQHotelIDSCoprID) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDSCoprID), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestBuildHotelSearchWithCitiesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	gcount := 3
	avail := BuildHotelAvailRq(gcount, q, sampleArrive, sampleDepart)
	avail.addCustomerID(sampleCID)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef) != len(hqcity[cityQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef), len(hqcity[cityQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQCitiesCustNumber) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQCitiesCustNumber), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestBuildHotelSearchWithLatLngMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	avail := BuildHotelAvailRq(sampleGuestCount, q, sampleArrive, sampleDepart)

	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}

	if len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef) != len(hqltln[latlngQueryField]) {
		t.Error("HotelRefs shoudl be same length as params", len(avail.Avail.HotelSearchCriteria.Criterion.HotelRef), len(hqltln[latlngQueryField]))
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleAvailRQLatLng) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQLatLng), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}
