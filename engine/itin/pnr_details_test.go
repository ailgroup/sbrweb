package itin

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

func TestPNRSet(t *testing.T) {
	p := CreatePersonName(sampleFirstName, sampleLastName)
	s := SetPNRDetailBody(samplePhoneReq, p)
	s.AddSpecialDetails()
	s.AddUniqueID("1234ABCD")
	addr := Address{
		AddressLine:   "123 Agora",
		Street:        "Sesame Street",
		City:          "Coalsville",
		StateProvince: StateProvince{StateCode: "XT"},
		CountryCode:   "LI",
		Postal:        "pae098",
	}
	vp := VendorPrefs{
		Airline: Airline{
			Hosted: false,
		},
	}
	s.AddAgencyInfo(addr, vp)

	agencyAddr := s.PassengerDetailsRQ.TravelItinInfo.Agency.Address
	if agencyAddr.Street != "Sesame Street" {
		t.Error("Agency Info street address is wrong")
	}
	if agencyAddr.StateProvince.StateCode != "XT" {
		t.Error("Agency Info StateProvince.StateCode address is wrong")
	}
	venPrefs := s.PassengerDetailsRQ.TravelItinInfo.Agency.VendorPrefs
	if venPrefs.Airline.Hosted {
		t.Error("Agency Info VendorPrefs Airline.Hosted is wrong")
	}

	if s.PassengerDetailsRQ.PreProcess.UniqueID.ID != "1234ABCD" {
		t.Errorf("s.PassengerDetailsRQ.PreProcess.UniqueID.ID given %v, built %v", "1234ABCD", s.PassengerDetailsRQ.PreProcess.UniqueID.ID)
	}

	spd := &SpecialReqDetails{}
	if s.PassengerDetailsRQ.SpecialReq.SpecialServiceRQ.SpecialServiceInfo.AdvancedPassenger.VendorPrefs.Airline.Hosted != spd.SpecialServiceRQ.SpecialServiceInfo.AdvancedPassenger.VendorPrefs.Airline.Hosted {
		t.Errorf("AddSpecialDetails \ngiven: %v \nbuilt: %v", spd, s.PassengerDetailsRQ.SpecialReq)
	}

	pn := s.PassengerDetailsRQ.TravelItinInfo.Customer.PersonName
	if pn.First.Val != sampleFirstName {
		t.Errorf("TravelItinInfo.Customer.PersonName.First expect: %s, got %s", sampleFirstName, pn.First.Val)
	}
	if pn.Last.Val != sampleLastName {
		t.Errorf("TravelItinInfo.Customer.PersonName.Last expect: %s, got %s", sampleLastName, pn.Last.Val)
	}
}

func TestPNRBuildMarshal(t *testing.T) {
	p := CreatePersonName(sampleFirstName, sampleLastName)
	body := SetPNRDetailBody(samplePhoneReq, p)
	req := BuildPNRDetailsRequest(sampleConf, body)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling passenger details request", err)
	}
	if string(b) != string(samplePNRReq) {
		t.Errorf("Expected marshal passenger details request \n given: %s \n built: %s", string(samplePNRReq), string(b))
	}
}

func TestPNRDetailCall(t *testing.T) {
	body := SetPNRDetailBody(samplePhoneReq, CreatePersonName(sampleFirstName, sampleLastName))
	req := BuildPNRDetailsRequest(sampleConf, body)
	resp, err := CallPNRDetail(serverPNRDetails.URL, req)
	if err != nil {
		t.Error("Error making request CallPNRDetailsRequest", err)
	}
	appres := resp.Body.PassengerDetailsRS.AppResults
	if appres.Status != "Complete" {
		t.Errorf("AppResults.Status expect: %s, got %s", "Complete", appres.Status)
	}
	if appres.Success.Timestamp != sampletimeOffset {
		t.Errorf("AppResults.Success.Timestamp expect: %s, got %s", sampletimeOffset, appres.Success.Timestamp)
	}
	if len(appres.Warnings) != 0 {
		t.Errorf("AppResults.Warnings expect: %d, got %d", 0, len(appres.Warnings))
	}
	travelItin := resp.Body.PassengerDetailsRS.TravelItineraryReadRS.TravelItinerary
	customer := travelItin.Customer

	if len(customer.ContactNumbers) != 1 {
		t.Error("wrong number of contact numbers")
	}
	numbersOne := customer.ContactNumbers[0]
	if numbersOne.Phone != samplePhoneRes {
		t.Errorf("customer.ContactNumbers[0].Phone expect: %s, got %s", samplePhoneRes, numbersOne.Phone)
	}
	if numbersOne.LocationCode != "SLC" {
		t.Errorf("customer.ContactNumbers[0].LocationCode expect: %s, got %s", "SLC", numbersOne.LocationCode)
	}
	if numbersOne.RPH != 1 {
		t.Errorf("customer.ContactNumbers[0].RPH expect: %d, got %d", 1, numbersOne.RPH)
	}
	person := customer.PersonName
	if person.WithInfant {
		t.Error("PersonName.WithInfant should be false")
	}
	if person.RPH != 1 {
		t.Errorf("PersonName.RPH expect: %d, got %d", 1, person.RPH)
	}
	if person.NameNumber != "01.01" {
		t.Errorf("PersonName.NameNumber expect: %s, got %s", "01.01", person.NameNumber)
	}
	if person.First.Val != strings.ToUpper(sampleFirstName) {
		t.Errorf("person.First expect: %s, got %s", sampleFirstName, person.First.Val)
	}
	if person.Last.Val != strings.ToUpper(sampleLastName) {
		t.Errorf("person.First expect: %s, got %s", sampleLastName, person.Last.Val)
	}

	if len(resp.Body.PassengerDetailsRS.TravelItineraryReadRS.TravelItinerary.ItineraryInfo.ReservationItems) > 1 {
		t.Error("ReservationItems wrong number")
	}

	itinRef := resp.Body.PassengerDetailsRS.TravelItineraryReadRS.TravelItinerary.ItineraryRef
	if itinRef.AirExtras {
		t.Error("ItineraryRef.AirExtras should be false")
	}

	if itinRef.InhibitCode != "U" {
		t.Errorf("ItineraryRef.InhibitCode expect: %s, got %s", "U", itinRef.InhibitCode)
	}
	if itinRef.PartitionID != "AA" {
		t.Errorf("ItineraryRef.PartitionID expect: %s, got %s", "AA", itinRef.PartitionID)
	}
	if itinRef.PrimeHostID != "1S" {
		t.Errorf("ItineraryRef.PrimeHostID expect: %s, got %s", "1S", itinRef.PrimeHostID)
	}
	if itinRef.Source.PseudoCityCode != samplepcc {
		t.Errorf("ItineraryRef.Source.PseudoCityCode expect: %s, got %s", samplepcc, itinRef.Source.PseudoCityCode)
	}
}

func TestPNRDetailCallWarn(t *testing.T) {
	body := SetPNRDetailBody(samplePhoneReq, CreatePersonName(sampleFirstName, sampleLastName))
	req := BuildPNRDetailsRequest(sampleConf, body)
	resp, err := CallPNRDetail(serverBizLogic.URL, req)
	if err == nil {
		t.Error("CallPNRDetailsRequest Should have errors", err)
	}
	if !resp.Body.Fault.Ok() {
		t.Error("Soap Fault be Ok() since errors was nil")
	}
	appRes := resp.Body.PassengerDetailsRS.AppResults
	if appRes.Ok() {
		t.Error("Application Results should not be Ok()")
	}
	if len(appRes.Warnings) != 2 {
		t.Errorf("Wrong number of warnings, want: %d, got %d", 2, len(appRes.Warnings))
	}
}

func TestPNRCallBadBodyResponseBody(t *testing.T) {
	p := CreatePersonName(sampleFirstName, sampleLastName)
	body := SetPNRDetailBody(samplePhoneReq, p)
	req := BuildPNRDetailsRequest(sampleConf, body)
	resp, err := CallPNRDetail(serverBadBody.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverBadBody")
	}
	if err.Error() != resp.ErrorSabreXML.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreXML.Code != sbrerr.BadParse {
		t.Errorf("Expect %d got %d", sbrerr.BadParse, resp.ErrorSabreXML.Code)
	}
	if resp.ErrorSabreXML.AppMessage != sbrerr.ErrCallPNRDetails {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallPNRDetails, resp.ErrorSabreXML.AppMessage)
	}
}

func TestPNRDetailsCallDown(t *testing.T) {
	body := SetPNRDetailBody(samplePhoneReq, CreatePersonName(sampleFirstName, sampleLastName))
	req := BuildPNRDetailsRequest(sampleConf, body)
	resp, err := CallPNRDetail(serverDown.URL, req)
	if err == nil {
		t.Error("Expected error making request to serverHotelDown")
	}
	if err.Error() != resp.ErrorSabreService.ErrMessage {
		t.Error("Error() message should match resp.ErrorSabreService.ErrMessage")
	}
	if resp.ErrorSabreService.Code != sbrerr.BadService {
		t.Errorf("Expect %d got %d", sbrerr.BadService, resp.ErrorSabreService.Code)
	}
	if resp.ErrorSabreService.AppMessage != sbrerr.ErrCallPNRDetails {
		t.Errorf("Expect %s got %s", sbrerr.ErrCallPNRDetails, resp.ErrorSabreService.AppMessage)
	}
}
