package sbrhotel

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/ailgroup/sbrweb"
)

var (
	hqbad               = make(HotelQueryParams)
	hqcity              = make(HotelQueryParams)
	hqids               = make(HotelQueryParams)
	sampleCorpID        = "12345"
	sampleNoCorpID      = ""
	sampleHotelCode     = []string{"0012", "19876", "1109", "445098", "000034"}
	sampleHotelCityCode = []string{"DFW", "CHC"}
	sampleGuestCount    = 2

	//sampleAvailRQHotelIDS = []byte(``)

	//sampleAvailRQCities = []byte(``)
)

func init() {
	hqbad[cityQueryField] = sampleHotelCityCode
	hqbad[hotelidQueryField] = sampleHotelCode

	hqcity[cityQueryField] = sampleHotelCityCode
	hqids[hotelidQueryField] = sampleHotelCode
}

func TestBuildHotelSearchReturnError(t *testing.T) {
	_, err := buildHotelSearch(hqbad)
	if err == nil {
		t.Error("HotelQueryParams with more than 1 key should have error")
	}
	if len(hqbad) <= 1 {
		t.Errorf("HotelQueryParams with error should have multiple map keys, got error message: '%v' for map: %v", err, hqbad)
	}
}

func TestBuildHotelSearchCity(t *testing.T) {
	q, err := buildHotelSearch(hqcity)
	if err != nil {
		t.Error("buildHotelSearch with good params should not have error!")
	}
	_, ok := hqcity[cityQueryField]
	if !ok {
		t.Error("cityQueryField should not be empty")
	}

	if len(q.Criterion.HotelRef) < 1 {
		t.Error("HotelSearchCriteria.Criterion.HotelRef should not be empty")
	}

	if q.Criterion.HotelRef[0].HotelCityCode != "DFW" {
		t.Errorf("HotelSearchCriteria.Criterion.HotelRef[0].HotelCityCode expect: %s, got: %s", sampleHotelCityCode[0], q.Criterion.HotelRef[0].HotelCityCode)
	}
	if q.Criterion.HotelRef[1].HotelCityCode != "CHC" {
		t.Errorf("HotelSearchCriteria.Criterion.HotelRef[1].HotelCityCode expect: %s, got: %s", sampleHotelCityCode[1], q.Criterion.HotelRef[1].HotelCityCode)
	}
}

func TestBuildHotelSearchHotelIDs(t *testing.T) {
	q, err := buildHotelSearch(hqids)
	if err != nil {
		t.Error("buildHotelSearch with good params should not have error!")
	}
	_, ok := hqids[hotelidQueryField]
	if !ok {
		t.Error("hotelidQueryField should not be empty")
	}

	if len(q.Criterion.HotelRef) < 1 {
		t.Error("HotelSearchCriteria.Criterion.HotelRef should not be empty")
	}

	if q.Criterion.HotelRef[0].HotelCode != "0012" {
		t.Errorf("HotelSearchCriteria.Criterion.HotelRef[0].HotelCode expect: %s, got: %s", sampleHotelCode[0], q.Criterion.HotelRef[0].HotelCode)
	}
	if q.Criterion.HotelRef[1].HotelCode != "19876" {
		t.Errorf("HotelSearchCriteria.Criterion.HotelRef[1].HotelCityCode expect: %s, got: %s", sampleHotelCode[1], q.Criterion.HotelRef[1].HotelCode)
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
	/*
		if avail.Avail.Customer.Corporate.ID != sampleCorpID {
			t.Errorf("BuildHotelAvailRq Customer.Corporate.ID expect: %s, got %s", sampleCorpID, avail.Avail.Customer.Corporate.ID)
		}
	*/

	_, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
}

func TestBuildHotelSearchWithIDSMarshal(t *testing.T) {
	q, _ := buildHotelSearch(hqids)
	gcount := 4
	avail := BuildHotelAvailRq(sampleCorpID, gcount, q)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	/*
		if string(b) != string(sampleAvailRQHotelIDS) {
			t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQHotelIDS), string(b))
		}
	*/
	fmt.Printf("content marshal \n%s\n", b)
}
func TestBuildHotelSearchWithCitiesMarshal(t *testing.T) {
	q, _ := buildHotelSearch(hqcity)
	gcount := 4
	avail := BuildHotelAvailRq(sampleNoCorpID, gcount, q)

	if avail.Avail.GuestCounts.Count != gcount {
		t.Errorf("BuildHotelAvailRq GuestCounts.Count expect: %d, got %d", gcount, avail.Avail.GuestCounts.Count)
	}

	b, err := xml.Marshal(avail)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	/*
		if string(b) != string(sampleAvailRQCities) {
			t.Errorf("Expected marshal hotel avail for hotel ids \n sample: %s \n result: %s", string(sampleAvailRQCities), string(b))
		}
	*/
	fmt.Printf("content marshal \n%s\n", b)
}
