package itin

import (
	"encoding/xml"
	"testing"
)

var (
	sampleLocator          = "IJKZUQ"
	sampleRequestType      = "Stateless"
	sampleSubjectArea      = "PRICE_QUOTE"
	sampleViewName         = "Simple"
	sampleResponseFormat   = "STL"
	sampleGetReservationRQ = []byte(`<soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/" xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" xmlns:xlink="http://www.w3.org/2001/xlink" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><soap-env:Header><eb:MessageHeader soap-env:mustUnderstand="1" eb:version="2.0.0"><eb:From><eb:PartyId type="urn:x12.org:IO5:01">z.com</eb:PartyId></eb:From><eb:To><eb:PartyId type="urn:x12.org:IO5:01">webservices.sabre.com</eb:PartyId></eb:To><eb:CPAId>ABCD1</eb:CPAId><eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId><eb:Service eb:type="sabreXML">GetReservationRQ</eb:Service><eb:Action>GetReservationLLSRQ</eb:Action><eb:MessageData><eb:MessageId>mid:20180216-07:18:42.3|14oUa</eb:MessageId><eb:Timestamp>2018-05-25T19:29:20Z</eb:Timestamp></eb:MessageData></eb:MessageHeader><wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext" xmlns:wsu="http://schemas.xmlsoap.org/ws/2002/12/utility"><wsse:BinarySecurityToken>Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESD!ICESMSLB\/RES.LB!-3142912682934961782!1421699!</wsse:BinarySecurityToken></wsse:Security></soap-env:Header><soap-env:Body><GetReservationRQ xmlns="http://webservices.sabre.com/pnrbuilder/v1_19" Version="1.19.0"><Locator>IJKZUQ</Locator><RequestType>Stateless</RequestType><ReturnOptions PriceQuoteServiceVersion="3.2.0"><SubjectAreas><SubjectArea>PRICE_QUOTE</SubjectArea></SubjectAreas><ViewName>Simple</ViewName><ResponseFormat>STL</ResponseFormat></ReturnOptions></GetReservationRQ></soap-env:Body></soap-env:Envelope>`)
	sampleGetReservationRS = []byte(`<stl19:GetReservationRS xmlns:stl19="http://webservices.sabre.com/pnrbuilder/v1_19" xmlns:ns6="http://services.sabre.com/res/orr/v0" xmlns:or114="http://services.sabre.com/res/or/v1_14" xmlns:raw="http://tds.sabre.com/itinerary" xmlns:ns4="http://webservices.sabre.com/pnrconn/ReaccSearch" Version="1.19.0"> <stl19:Reservation numberInParty="2" numberOfInfants="0" NumberInSegment="2"> <stl19:BookingDetails> <stl19:RecordLocator>IJKZUQ</stl19:RecordLocator> <stl19:CreationTimestamp>2018-07-05T06:39:00</stl19:CreationTimestamp> <stl19:SystemCreationTimestamp>2018-07-05T06:39:00</stl19:SystemCreationTimestamp> <stl19:CreationAgentID>MYE</stl19:CreationAgentID> <stl19:UpdateTimestamp>2018-07-20T02:01:36</stl19:UpdateTimestamp> <stl19:PNRSequence>6</stl19:PNRSequence> <stl19:DivideSplitDetails/> <stl19:EstimatedPurgeTimestamp>2018-08-30T00:00:00</stl19:EstimatedPurgeTimestamp> <stl19:UpdateToken>-12c37dd4c51e882df40aef3b324db514b62ddf6d93dd1c12</stl19:UpdateToken> </stl19:BookingDetails> <stl19:POS AirExtras="false" InhibitCode="U"> <stl19:Source BookingSource="A0A0" AgentSine="MYE" PseudoCityCode="BKK" ISOCountry="TH" AgentDutyCode="5" AirlineVendorID="PG" HomePseudoCityCode="BKK"/> </stl19:POS> <stl19:PassengerReservation> <stl19:Passengers/> <stl19:Segments/> <stl19:TicketingInfo/> <stl19:ItineraryPricing/> </stl19:PassengerReservation> <stl19:ReceivedFrom/> <stl19:EmailAddresses/> </stl19:Reservation> <or114:PriceQuote> <PriceQuoteInfo xmlns="http://www.sabre.com/ns/Ticketing/pqs/1.0"> <Reservation updateToken="eNc:::3qIVVfgBM9uZ0A7HUjkopQ==">IJKZUQ</Reservation> <Summary> <NameAssociation firstName="WAIVE" lastName="TWOPAX" nameId="1" nameNumber="1.1"> <PriceQuote number="1" pricingType="S" status="I" type="PQ"> <Indicators itineraryChange="true"/> <Passenger passengerTypeCount="1" requestedType="ADT" type="ADT"/> <ItineraryType>I</ItineraryType> <Fee code="T01" itemId="1" type="OB"> <Amount currencyCode="THB" decimalPlace="0">140</Amount> </Fee> <ValidatingCarrier>PG</ValidatingCarrier> <Amounts> <Total currencyCode="THB">8890</Total> </Amounts> <LocalCreateDateTime>2018-07-05T18:35:00</LocalCreateDateTime> </PriceQuote> </NameAssociation> <NameAssociation firstName="CHARGE" lastName="TWOPAX" nameId="2" nameNumber="2.1"> <PriceQuote number="2" pricingType="S" status="I" type="PQ"> <Indicators itineraryChange="true"/> <Passenger passengerTypeCount="1" requestedType="ADT" type="ADT"/> <ItineraryType>I</ItineraryType> <Fee code="T01" itemId="1" type="OB"> <Amount currencyCode="THB" decimalPlace="0">140</Amount> <Waiver>04</Waiver> </Fee> <ValidatingCarrier>PG</ValidatingCarrier> <Amounts> <Total currencyCode="THB">8890</Total> </Amounts> <LocalCreateDateTime>2018-07-05T18:35:00</LocalCreateDateTime> </PriceQuote> </NameAssociation> </Summary> <Details number="1" passengerType="ADT" pricingType="S" status="I" type="PQ"> <AgentInfo duty="5" sine="MYE"> <HomeLocation>BKK</HomeLocation> <WorkLocation>BKK</WorkLocation> </AgentInfo> <TransactionInfo> <CreateDateTime>2018-07-05T06:35:00</CreateDateTime> <LocalCreateDateTime>2018-07-05T18:35:00</LocalCreateDateTime> <ExpiryDateTime>2019-06-01T00:00:00</ExpiryDateTime> <InputEntry>WPN1.1</InputEntry> </TransactionInfo> <NameAssociationInfo firstName="WAIVE" lastName="TWOPAX" nameId="1" nameNumber="1.1"/> <SegmentInfo number="1" segmentStatus="OK"> <Flight connectionIndicator="O"> <MarketingFlight number="707">PG</MarketingFlight> <ClassOfService>Y</ClassOfService> <Departure> <DateTime>2018-08-30T13:45:00</DateTime> <CityCode name="BANGKOK SUVARNABH">BKK</CityCode> </Departure> <Arrival> <DateTime>2018-08-30T14:35:00</DateTime> <CityCode name="YANGON">RGN</CityCode> </Arrival> </Flight> <FareBasis>YOWW</FareBasis> <NotValidAfter>2019-08-30</NotValidAfter> <Baggage allowance="20" type="K"/> </SegmentInfo> <FareInfo> <FareIndicators/> <BaseFare currencyCode="THB">7030</BaseFare> <TotalTax currencyCode="THB">1860</TotalTax> <TotalFare currencyCode="THB">8890</TotalFare> <TaxInfo> <CombinedTax code="TS"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <CombinedTax code="E7"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <CombinedTax code="XT"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <Tax code="TS"> <Amount currencyCode="THB">700</Amount> </Tax> <Tax code="E7"> <Amount currencyCode="THB">35</Amount> </Tax> <Tax code="G8"> <Amount currencyCode="THB">15</Amount> </Tax> <Tax code="YQ"> <Amount currencyCode="THB">1000</Amount> </Tax> <Tax code="C7"> <Amount currencyCode="THB">110</Amount> </Tax> </TaxInfo> <FareCalculation>BKK PG RGN218.11NUC218.11END ROE32.231</FareCalculation> <FareComponent fareBasisCode="YOWW" number="1"> <FlightSegmentNumbers> <SegmentNumber>1</SegmentNumber> </FlightSegmentNumbers> <FareDirectionality oneWay="true"/> <Departure> <DateTime>2018-08-30T13:45:00</DateTime> <CityCode name="BANGKOK SUVARNABH">BKK</CityCode> </Departure> <Arrival> <DateTime>2018-08-30T14:35:00</DateTime> <CityCode name="YANGON">RGN</CityCode> </Arrival> <Amount decimalPlace="2">218.11</Amount> <GoverningCarrier>PG</GoverningCarrier> </FareComponent> </FareInfo> <FeeInfo> <OBFee code="T01" type="OB"> <Amount currencyCode="THB" decimalPlace="0">140</Amount> <Description>CARRIER TICKETING FEE01</Description> </OBFee> </FeeInfo> <MiscellaneousInfo> <ValidatingCarrier>PG</ValidatingCarrier> <ItineraryType>I</ItineraryType> </MiscellaneousInfo> <MessageInfo> <Message number="301" type="INFO">One or more form of payment fees may apply</Message> <Message number="302" type="INFO">Actual total will be based on form of payment used</Message> <Message number="201" type="WARNING">Fare not guaranteed until ticketed</Message> <Message type="WARNING">PRIVATE FARE APPLIED - CHECK RULES FOR CORRECT TICKETING</Message> <Message type="WARNING">VALIDATING CARRIER SPECIFIED - PG</Message> <Remarks type="ENS">REFUND AND CHANGE/RESTRICTED/NON ENDORSE</Remarks> <PricingParameters>WPN1.1</PricingParameters> </MessageInfo> </Details> <Details number="2" passengerType="ADT" pricingType="S" status="I" type="PQ"> <AgentInfo duty="5" sine="MYE"> <HomeLocation>BKK</HomeLocation> <WorkLocation>BKK</WorkLocation> </AgentInfo> <TransactionInfo> <CreateDateTime>2018-07-05T06:35:00</CreateDateTime> <LocalCreateDateTime>2018-07-05T18:35:00</LocalCreateDateTime> <ExpiryDateTime>2019-06-01T00:00:00</ExpiryDateTime> <InputEntry>WPN2.1</InputEntry> </TransactionInfo> <NameAssociationInfo firstName="CHARGE" lastName="TWOPAX" nameId="2" nameNumber="2.1"/> <SegmentInfo number="1" segmentStatus="OK"> <Flight connectionIndicator="O"> <MarketingFlight number="707">PG</MarketingFlight> <ClassOfService>Y</ClassOfService> <Departure> <DateTime>2018-08-30T13:45:00</DateTime> <CityCode name="BANGKOK SUVARNABH">BKK</CityCode> </Departure> <Arrival> <DateTime>2018-08-30T14:35:00</DateTime> <CityCode name="YANGON">RGN</CityCode> </Arrival> </Flight> <FareBasis>YOWW</FareBasis> <NotValidAfter>2019-08-30</NotValidAfter> <Baggage allowance="20" type="K"/> </SegmentInfo> <FareInfo> <FareIndicators/> <BaseFare currencyCode="THB">7030</BaseFare> <TotalTax currencyCode="THB">1860</TotalTax> <TotalFare currencyCode="THB">8890</TotalFare> <TaxInfo> <CombinedTax code="TS"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <CombinedTax code="E7"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <CombinedTax code="XT"> <Amount currencyCode="THB">21811</Amount> </CombinedTax> <Tax code="TS"> <Amount currencyCode="THB">700</Amount> </Tax> <Tax code="E7"> <Amount currencyCode="THB">35</Amount> </Tax> <Tax code="G8"> <Amount currencyCode="THB">15</Amount> </Tax> <Tax code="YQ"> <Amount currencyCode="THB">1000</Amount> </Tax> <Tax code="C7"> <Amount currencyCode="THB">110</Amount> </Tax> </TaxInfo> <FareCalculation>BKK PG RGN218.11NUC218.11END ROE32.231</FareCalculation> <FareComponent fareBasisCode="YOWW" number="1"> <FlightSegmentNumbers> <SegmentNumber>1</SegmentNumber> </FlightSegmentNumbers> <FareDirectionality oneWay="true"/> <Departure> <DateTime>2018-08-30T13:45:00</DateTime> <CityCode name="BANGKOK SUVARNABH">BKK</CityCode> </Departure> <Arrival> <DateTime>2018-08-30T14:35:00</DateTime> <CityCode name="YANGON">RGN</CityCode> </Arrival> <Amount decimalPlace="2">218.11</Amount> <GoverningCarrier>PG</GoverningCarrier> </FareComponent> </FareInfo> <FeeInfo> <OBFee code="T01" type="OB"> <Amount currencyCode="THB" decimalPlace="0">140</Amount> <Description>CARRIER TICKETING FEE01</Description> </OBFee> </FeeInfo> <MiscellaneousInfo> <ValidatingCarrier>PG</ValidatingCarrier> <ItineraryType>I</ItineraryType> </MiscellaneousInfo> <MessageInfo> <Message number="301" type="INFO">One or more form of payment fees may apply</Message> <Message number="302" type="INFO">Actual total will be based on form of payment used</Message> <Message number="201" type="WARNING">Fare not guaranteed until ticketed</Message> <Message type="WARNING">PRIVATE FARE APPLIED - CHECK RULES FOR CORRECT TICKETING</Message> <Message type="WARNING">VALIDATING CARRIER SPECIFIED - PG</Message> <Remarks type="ENS">REFUND AND CHANGE/RESTRICTED/NON ENDORSE</Remarks> <PricingParameters>WPN2.1</PricingParameters> </MessageInfo> </Details> </PriceQuoteInfo> </or114:PriceQuote> </stl19:GetReservationRS>`)
)

func TestBuildGetReservationMarshal(t *testing.T) {
	req := BuildGetReservationRequest(sampleConf, sampleLocator, sampleRequestType, sampleSubjectArea, sampleViewName, sampleResponseFormat)
	b, err := xml.Marshal(req)
	if err != nil {
		t.Error("Error marshaling pnr read request", err)
	}
	if string(b) != string(sampleGetReservationRQ) {
		t.Errorf("Expected marshal get reservation request \n given: %s \n built: %s", string(sampleGetReservationRQ), string(b))
	}
}

func TestGetReservationRQUnmarshal(t *testing.T) {
	getRes := GetReservationRS{}
	err := xml.Unmarshal(sampleGetReservationRS, &getRes)
	if err != nil {
		t.Fatal("Error unmarshal get reservation response", err)
	}
}

func genPriceQuote(itc bool, ptc, rqt, pt, it, fcode, itid, ft, facc, dp, aval, valc, tcc, tval, lcdt string) PriceQuoteNameAssocElem {
	return PriceQuoteNameAssocElem{
		Indicators: Indicators{ItineraryChange: itc},
		Passenger: Passenger{
			PassengerTypeCount: ptc,
			RequestedType:      rqt,
			Ptype:              pt,
		},
		ItineraryType: ItineraryType{it},
		Fee: Fee{
			Code:   fcode,
			ItemID: itid,
			Ftype:  ft,
			Amount: Amount{
				CurrencyCode: facc,
				DecimalPlace: dp,
				Val:          aval,
			},
		},
		ValidatingCarrier:   ValidatingCarrier{valc},
		Amounts:             Amounts{Total: Total{CurrencyCode: tcc, Val: tval}},
		LocalCreateDateTime: lcdt,
	}
	/*
		PriceQuote: Indicators: {ItineraryChange: true}, Passenger: {PassengerTypeCount: 1, RequestedType: "ADT", Ptype: "ADT"}, ItineraryType: "I", Fee: {Code: "T01", ItemID: "1", Ftype: "OB", Amount: {CurrencyCode: "THB", DecimalPlace: "0"}}, ValidatingCarrier: "PG", Amounts: {Total: {CurrencyCode: "THB"}}, LocalCreateDateTime: "2018-07-05T18:35:00",}

		PriceQuote: Indicators:{ItineraryChange:true} Passenger:{PassengerTypeCount:1 RequestedType:ADT Ptype:ADT} ItineraryType:I Fee:{Code:T01 ItemID:1 Ftype:OB Amount:{CurrencyCode:THB DecimalPlace:0}} ValidatingCarrier:PG Amounts:{Total:{CurrencyCode:THB}}LocalCreateDateTime: "2018-07-05T18:35:00",}
	*/
}

var nameAssocTest = []struct {
	expect NameAssociation
}{
	{NameAssociation{
		FirstName:  "WAIVE",
		LastName:   "TWOPAX",
		NameID:     "",
		NameNumber: "1.1",
		PriceQuote: genPriceQuote(true, "1", "ADT", "ADT", "I", "T01", "1", "OB", "THB", "0", "140", "PG", "THB", "8890", "2018-07-05T18:35:00"),
	}},
	{NameAssociation{
		FirstName:  "CHARGE",
		LastName:   "TWOPAX",
		NameID:     "",
		NameNumber: "2.1",
		PriceQuote: genPriceQuote(true, "1", "ADT", "ADT", "I", "T01", "1", "OB", "THB", "0", "140", "PG", "THB", "8890", "2018-07-05T18:35:00"),
	}},
}

func TestGetReservationRQPriceQuoteInfoResAndSummary(t *testing.T) {
	getRes := GetReservationRS{}
	xml.Unmarshal(sampleGetReservationRS, &getRes)
	quoteInfo := getRes.PriceQuote.PriceQuoteInfo
	summary := quoteInfo.Summary
	//res := quoteInfo.Reservation
	//fmt.Printf("\n PQ RES: %+v \n\n", res)
	//fmt.Printf("\n PQ SUM-NAME-ASSOC: %+v \n\n", summary.NameAssociations[0])
	if len(summary.NameAssociations) != 2 {
		t.Error("PriceQuoteInfo.NameAssociations length should be 2")
	}

	qresVal := "IJKZUQ"
	if quoteInfo.Reservation.Val != qresVal {
		t.Errorf("Reservation value expect: %s, got: %s", qresVal, quoteInfo.Reservation)
	}
	tkn := "eNc:::3qIVVfgBM9uZ0A7HUjkopQ=="
	if quoteInfo.Reservation.UpdateToken != tkn {
		t.Errorf("UpdateToken expect: %s have: %s", tkn, quoteInfo.Reservation.UpdateToken)
	}

	for i, nas := range nameAssocTest {
		if summary.NameAssociations[i].FirstName != nas.expect.FirstName {
			t.Errorf("summary.NameAssociation %d for FirstName expect: %s, got: %s", i, nas.expect.FirstName, summary.NameAssociations[i].FirstName)
		}
		if summary.NameAssociations[i].NameNumber != nas.expect.NameNumber {
			t.Errorf("summary.NameAssociation %d for NameNumber expect: %s, got: %s", i, nas.expect.NameNumber, summary.NameAssociations[i].NameNumber)
		}

		naPq := summary.NameAssociations[i].PriceQuote
		expectNaPq := nas.expect.PriceQuote
		if naPq.ItineraryType != expectNaPq.ItineraryType {
			t.Errorf("summary.NameAssociations.PriceQuote %d for RequestedType expect: %s, got: %s", i, nas.expect.PriceQuote.ItineraryType, summary.NameAssociations[i].PriceQuote.ItineraryType)
		}
		if naPq.Passenger.RequestedType != expectNaPq.Passenger.RequestedType {
			t.Errorf("summary.NameAssociations.PriceQuote %d for RequestedType expect: %s, got: %s", i, nas.expect.PriceQuote.Passenger.RequestedType, summary.NameAssociations[i].PriceQuote.Passenger.RequestedType)
		}
		if naPq.Passenger.PassengerTypeCount != expectNaPq.Passenger.PassengerTypeCount {
			t.Errorf("summary.NameAssociations.PriceQuote %d for PassengerTypeCount expect: %s, got: %s", i, nas.expect.PriceQuote.Passenger.PassengerTypeCount, summary.NameAssociations[i].PriceQuote.Passenger.PassengerTypeCount)
		}
		if naPq.Amounts.Total.Val != expectNaPq.Amounts.Total.Val {
			t.Errorf("summary.NameAssociations.PriceQuote %d for Amounts.Total.Val expect: %s, got: %s", i, nas.expect.PriceQuote.Amounts.Total.Val, naPq.Amounts.Total.Val)
		}
		if naPq.Amounts.Total.CurrencyCode != expectNaPq.Amounts.Total.CurrencyCode {
			t.Errorf("summary.NameAssociations.PriceQuote %d for Amounts.Total.CurrencyCode expect: %s, got: %s", i, nas.expect.PriceQuote.Amounts.Total.CurrencyCode, naPq.Amounts.Total.CurrencyCode)
		}
	}
}

var pqDetailsTest = []struct {
	expect DetailsPriceQuoteElem
}{
	{DetailsPriceQuoteElem{
		Number:        "1",
		PassengerType: "ADT",
		PricingType:   "S",
		Status:        "I",
		Dtype:         "PQ",
		AgentInfo: AgentInfo{
			Duty:         "5",
			Sine:         "MYE",
			HomeLocation: HomeLocation{Val: "BKK"},
			WorkLocation: WorkLocation{Val: "BKK"},
		},
		TransactionInfo: TransactionInfo{
			CreateDateTime:      CreateDateTime{"2018-07-05T06:35:00"},
			LocalCreateDateTime: LocalCreateDateTime{"2018-07-05T18:35:00"},
			ExpiryDateTime:      ExpiryDateTime{"2019-06-01T00:00:00"},
			InputEntry:          InputEntry{"WPN1.1"},
		},
		NameAssociationInfo: NameAssociationInfo{
			FirstName:  "WAIVE",
			LastName:   "TWOPAX",
			NameID:     "1",
			NameNumber: "1.1",
		},
		SegmentInfo: SegmentInfo{
			Number:        "1",
			SegmentStatus: "OK",
			Flight: Flight{
				ConnectionIndicator: "O",
				MarketingFlight:     MarketingFlight{Number: "707", Val: "PG"},
				ClassOfService:      ClassOfService{"Y"},
				Departure:           Departure{},
				Arrival:             Arrival{},
			},
			FareBasis:     FareBasis{""},
			NotValidAfter: NotValidAfter{""},
			Baggage:       Baggage{Allowance: "", Btype: ""},
		},
	}},
	{DetailsPriceQuoteElem{
		Number:        "2",
		PassengerType: "ADT",
		PricingType:   "S",
		Status:        "I",
		Dtype:         "PQ",
		AgentInfo: AgentInfo{
			Duty:         "5",
			Sine:         "MYE",
			HomeLocation: HomeLocation{Val: "BKK"},
			WorkLocation: WorkLocation{Val: "BKK"},
		},
		TransactionInfo: TransactionInfo{
			CreateDateTime:      CreateDateTime{""},
			LocalCreateDateTime: LocalCreateDateTime{""},
			ExpiryDateTime:      ExpiryDateTime{""},
			InputEntry:          InputEntry{""},
		},
		NameAssociationInfo: NameAssociationInfo{
			FirstName:  "CHARGE",
			LastName:   "TWOPAX",
			NameID:     "2",
			NameNumber: "2.1",
		},
		SegmentInfo: SegmentInfo{
			Number:        "1",
			SegmentStatus: "",
			Flight: Flight{
				ConnectionIndicator: "",
				MarketingFlight:     MarketingFlight{Number: "", Val: ""},
				ClassOfService:      ClassOfService{""},
				Departure:           Departure{},
				Arrival:             Arrival{},
			},
			FareBasis:     FareBasis{""},
			NotValidAfter: NotValidAfter{""},
			Baggage:       Baggage{Allowance: "", Btype: ""},
		},
	}},
}

func TestGetReservationRQPriceQuoteInfoDetails(t *testing.T) {
	getRes := GetReservationRS{}
	xml.Unmarshal(sampleGetReservationRS, &getRes)
	details := getRes.PriceQuote.PriceQuoteInfo.Details
	//fmt.Printf("\n PQ DETAILS: %+v \n\n", details[0])
	if len(details) != 2 {
		t.Error("PriceQuoteInfo.Details should be 2")
	}

	for i, test := range pqDetailsTest {
		if details[i].Number != test.expect.Number {
			t.Errorf("Number for idx-%d expect: %s got: %s", i, details[i].Number, test.expect.Number)
		}
		if details[i].PassengerType != test.expect.PassengerType {
			t.Errorf("PassengerType for idx-%d expect: %s got: %s", i, details[i].PassengerType, test.expect.PassengerType)
		}
		if details[i].PricingType != test.expect.PricingType {
			t.Errorf("PricingType for idx-%d expect: %s got: %s", i, details[i].PricingType, test.expect.PricingType)
		}
		if details[i].Status != test.expect.Status {
			t.Errorf("Status for idx-%d expect: %s got: %s", i, details[i].PricingType, test.expect.PricingType)
		}
		if details[i].Dtype != test.expect.Dtype {
			t.Errorf("Dtype for idx-%d expect: %s got: %s", i, details[i].Dtype, test.expect.Dtype)
		}

		if details[i].AgentInfo != test.expect.AgentInfo {
			t.Errorf("AgentInfo for idx-%d no match \nexpect: %+v \ngot %+v", i, test.expect.AgentInfo, details[i].AgentInfo)
		}
		/*
			if details[i].TransactionInfo != test.expect.TransactionInfo {
				t.Errorf("TransactionInfo for idx-%s no match \nexpect: %+v \ngot %+v", i, test.expect.TransactionInfo, details[i].TransactionInfo)
			}
		*/
		if details[i].NameAssociationInfo != test.expect.NameAssociationInfo {
			t.Errorf("NameAssociationInfo for idx-%d no match \nexpect: %+v \ngot %+v", i, test.expect.NameAssociationInfo, details[i].NameAssociationInfo)
		}
		/*
			if details[i].SegmentInfo != test.expect.SegmentInfo {
				t.Errorf("SegmentInfo for idx-%s no match \nexpect: %+v \ngot %+v", i, test.expect.SegmentInfo, details[i].SegmentInfo)
			}
		*/
	}
}
