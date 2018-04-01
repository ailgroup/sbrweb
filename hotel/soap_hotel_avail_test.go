package sbrhotel

import (
	"encoding/xml"
	"fmt"
	"strings"
	"testing"

	"github.com/ailgroup/sbrweb"
)

var (
	hqbad               = make(HotelRefCriterion)
	hqcity              = make(HotelRefCriterion)
	hqids               = make(HotelRefCriterion)
	hqltln              = make(HotelRefCriterion)
	addr                = make(AddressCriterion)
	sampleCorpID        = "12345"
	sampleNoCorpID      = ""
	sampleLatLang       = []string{"32.78,-96.81", "54.87,-102.96"}
	sampleHotelCode     = []string{"0012", "19876", "1109", "445098", "000034"}
	sampleHotelCityCode = []string{"DFW", "CHC", "LA"}
	sampleGuestCount    = 2
	sampleStreet        = "2031 N. 100 W"
	sampleCity          = "Nowhere"
	samplePostal        = "999908"
	sampleCountryCode   = "US"

	//sampleAvailRQHotelIDSCoprID = []byte(``)

	//sampleAvailRQCities = []byte(``)

	//sampleAvailRQLatLng = []byte(``)
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

	/*
		avail := BuildHotelAvailRq(sampleCorpID, sampleGuestCount, r)
		b, err := xml.Marshal(avail)
		if err != nil {
			t.Error("Error marshaling get hotel content", err)
		}
		fmt.Printf("\n%s\n", b)
	*/

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
}

func TestBuildHotelSearchMarshal(t *testing.T) {
	avail := BuildHotelAvailRq(sampleCorpID, sampleGuestCount, HotelSearchCriteria{})

	if avail.XMLNSXsi != sbrweb.BaseXSINamespace {
		t.Errorf("BuildHotelAvailRq XMLNSXsi expect: %s, got %s", sbrweb.BaseXSINamespace, avail.XMLNSXsi)
	}
	if avail.Version != hotelAvailVersion {
		t.Errorf("BuildHotelAvailRq Version expect: %s, got %s", hotelAvailVersion, avail.Version)
	}
	if avail.Avail.GuestCounts.Count != sampleGuestCount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", sampleGuestCount, avail.Avail.GuestCounts.Count)
	}
	if avail.Avail.Customer.Corporate.ID != sampleCorpID {
		t.Errorf("BuildHotelAvailRq Customer.Corporate.ID expect: %s, got %s", sampleCorpID, avail.Avail.Customer.Corporate.ID)
	}

	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestBuildHotelSearchNoCorpIDMarshal(t *testing.T) {
	avail := BuildHotelAvailRq(sampleNoCorpID, sampleGuestCount, HotelSearchCriteria{})
	customer := Customer{}

	if avail.Avail.Customer != customer {
		t.Errorf("BuildHotelAvailRq Customer for empty corporate ID should be empty expect: %v, got %v", customer, avail.Avail.Customer)
	}

	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestBuildHotelSearchWithIDSMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqids),
	)
	gcount := 4
	avail := BuildHotelAvailRq(sampleCorpID, gcount, q)

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
	/*
		if string(b) != string(sampleAvailRQHotelIDSCoprID) {
			t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDSCoprID), string(b))
		}
	*/
	fmt.Printf("content marshal \n%s\n", b)
}

/*
func TestBuildHotelSearchWithCitiesMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqcity),
	)
	gcount := 3
	avail := BuildHotelAvailRq(sampleNoCorpID, gcount, q)

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
	if string(b) != string(sampleAvailRQCities) {
		t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQCities), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}

func TestBuildHotelSearchWithLatLngMarshal(t *testing.T) {
	q, _ := NewHotelSearchCriteria(
		HotelRefSearch(hqltln),
	)
	avail := BuildHotelAvailRq(sampleNoCorpID, sampleGuestCount, q)

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
*/
