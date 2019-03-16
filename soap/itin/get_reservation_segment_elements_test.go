package itin

import (
	"encoding/xml"
	"testing"
)

var (
	sampleGetResSegmentRS = []byte(`<?xml version="1.0" encoding="UTF-8"?> <soap-env:Envelope xmlns:soap-env="http://schemas.xmlsoap.org/soap/envelope/"> <soap-env:Header> <eb:MessageHeader xmlns:eb="http://www.ebxml.org/namespaces/messageHeader" eb:version="1.0" soap-env:mustUnderstand="1"> <eb:From> <eb:PartyId eb:type="URI">webservices.sabre.com</eb:PartyId> </eb:From> <eb:To> <eb:PartyId eb:type="URI">www.z.com</eb:PartyId> </eb:To> <eb:CPAId>ABDC1</eb:CPAId> <eb:ConversationId>fds8789h|dev@z.com</eb:ConversationId> <eb:Service eb:type="sabreXML">GetReservationRQ</eb:Service> <eb:Action>GetReservationRQ</eb:Action> <eb:MessageData> <eb:MessageId>6652823509217390550</eb:MessageId> <eb:Timestamp>2019-01-28T14:08:41</eb:Timestamp> <eb:RefToMessageId>mid:20180216-07:18:42.3|14oUa</eb:RefToMessageId> </eb:MessageData> </eb:MessageHeader> <wsse:Security xmlns:wsse="http://schemas.xmlsoap.org/ws/2002/12/secext"> <wsse:BinarySecurityToken valueType="String" EncodingType="wsse:Base64Binary">Shared/IDL:IceSess\/SessMgr:1\.0.IDL/Common/!ICESMS\/RESE!ICESMSLB\/RES.LB!1548684521533!5933!41</wsse:BinarySecurityToken> </wsse:Security> </soap-env:Header> <soap-env:Body> <stl19:GetReservationRS xmlns:stl19="http://webservices.sabre.com/pnrbuilder/v1_19" xmlns:ns6="http://services.sabre.com/res/orr/v0" xmlns:or114="http://services.sabre.com/res/or/v1_14" xmlns:raw="http://tds.sabre.com/itinerary" xmlns:ns4="http://webservices.sabre.com/pnrconn/ReaccSearch" Version="1.19.0"> <stl19:Reservation numberInParty="1" numberOfInfants="0" NumberInSegment="1"> <stl19:BookingDetails> <stl19:RecordLocator>YESWPL</stl19:RecordLocator> <stl19:CreationTimestamp>2019-01-28T07:43:00</stl19:CreationTimestamp> <stl19:SystemCreationTimestamp>2019-01-28T07:43:00</stl19:SystemCreationTimestamp> <stl19:CreationAgentID>AWS</stl19:CreationAgentID> <stl19:UpdateTimestamp>2019-01-28T07:43:28</stl19:UpdateTimestamp> <stl19:PNRSequence>1</stl19:PNRSequence> <stl19:DivideSplitDetails/> <stl19:EstimatedPurgeTimestamp>2019-02-18T00:00:00</stl19:EstimatedPurgeTimestamp> <stl19:UpdateToken>-4d195ac0ba963bff00c17d740b10d2b0ef5bf1c8347447b8</stl19:UpdateToken> </stl19:BookingDetails> <stl19:POS AirExtras="false" InhibitCode="U"> <stl19:Source BookingSource="ABCD1" AgentSine="AWS" PseudoCityCode="ABCD1" ISOCountry="US" AgentDutyCode="*" AirlineVendorID="AA" HomePseudoCityCode="ABCD1" PrimeHostID="1S"/> </stl19:POS> <stl19:PassengerReservation> <stl19:Passengers> <stl19:Passenger id="3" nameType="S" passengerType="ADT" referenceNumber="ABC123" nameId="01.01" nameAssocId="1" elementId="pnr-3.1"> <stl19:LastName>DE MONTAIGNE</stl19:LastName> <stl19:FirstName>MICHEL</stl19:FirstName> <stl19:Seats/> </stl19:Passenger> </stl19:Passengers> <stl19:Segments> <stl19:Segment sequence="1" id="36"> <stl19:Hotel id="36" sequence="1" isPast="false"> <or114:Reservation DayOfWeekInd="1" NumberInParty="02"> <or114:LineNumber>1</or114:LineNumber> <or114:LineType>HHL</or114:LineType> <or114:LineStatus>HK</or114:LineStatus> <or114:POSRequestorID>46595076</or114:POSRequestorID> <or114:RoomType> <or114:RoomTypeCode>CSP</or114:RoomTypeCode> <or114:NumberOfUnits>1</or114:NumberOfUnits> <or114:ShortText>CSP08PG</or114:ShortText> </or114:RoomType> <or114:RoomRates> <or114:AmountBeforeTax>450.00</or114:AmountBeforeTax> <or114:CurrencyCode>EUR</or114:CurrencyCode> </or114:RoomRates> <or114:RateAccessCodeBooked/> <or114:GuestCounts> <or114:GuestCount>2</or114:GuestCount> <or114:ExtraGuestCount>0</or114:ExtraGuestCount> <or114:RollAwayCount>0</or114:RollAwayCount> <or114:CribCount>0</or114:CribCount> </or114:GuestCounts> <or114:TimeSpanStart>2019-02-18T00:00:00</or114:TimeSpanStart> <or114:TimeSpanDuration>1</or114:TimeSpanDuration> <or114:TimeSpanEnd>2019-02-19T00:00:00</or114:TimeSpanEnd> <or114:Guarantee> <or114:Text>GVI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE</or114:Text> </or114:Guarantee> <or114:ChainCode>IC</or114:ChainCode> <or114:HotelCode>1098</or114:HotelCode> <or114:HotelCityCode>AMS</or114:HotelCityCode> <or114:HotelName>INTERCONTINENTAL AMSTEL AMS</or114:HotelName> <or114:HotelTotalPricing> <or114:TotalTax Amount="27.00"> <or114:Tax Id="1">27.00 SERVICE CHARGE</or114:Tax> </or114:TotalTax> <or114:ApproximateTotal AmountAndCurrency="477.00 EUR"/> <or114:Disclaimer Id="1">INCLUDES TAX</or114:Disclaimer> </or114:HotelTotalPricing> </or114:Reservation> <or114:AdditionalInformation> <or114:ConfirmationNumber DirectConnect="true">23164323-</or114:ConfirmationNumber> <or114:Address> <or114:AddressLine>PROFESSOR TULPPLEIN 1</or114:AddressLine> <or114:AddressLine>AMSTERDAM NL 1018 GX</or114:AddressLine> <or114:CountryCode>NL</or114:CountryCode> <or114:City>AMSTERDAM</or114:City> <or114:ZipCode>1018 GX</or114:ZipCode> </or114:Address> <or114:ContactNumbers> <or114:PhoneNumber>31-20-6226060</or114:PhoneNumber> <or114:FaxNumber>31-20-6225808</or114:FaxNumber> </or114:ContactNumbers> <or114:CancelPenaltyPolicyCode>01D</or114:CancelPenaltyPolicyCode> <or114:Commission> <or114:Indicator>C</or114:Indicator> <or114:Text>COMMISSIONABLE</or114:Text> </or114:Commission> </or114:AdditionalInformation> <or114:SegmentText>450.00EUR/RC-@@@-@@@-@/TTX-27.00/HTP-477.00 EUR/TX1-27.00 SERVICE CHARGE/DS1-INCLUDES TAX/CMN-C/CMT-COMMISSIONABLE/AGT-46595076/G-VI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE/C-01D/XS-0001548683004531675D5EA000000000/XT-675D5EA0/XL-0000/DT-28JAN190463/SBS-HS/HSA-PROFESSOR TULPPLEIN 1$AMSTERDAM NL 1018 GX/HFN-INTERCONTINENTAL AMSTEL AMS/HCY-AMSTERDAM/HST-/HCC-NL/HPC-1018 GX/HPH-31-20-6226060/HFX-31-20-6225808/UID-D599525312BD5487675D5E/SI-/CF-23164323-</or114:SegmentText> </stl19:Hotel> <stl19:Product sequence="1" id="36"> <or114:ProductBase> <or114:SegmentReference>36</or114:SegmentReference> </or114:ProductBase> <or114:ProductDetails vendorCode="IC" statusCode="HK" previousStatusCode="SS" startDateTime="2019-02-18T00:00:00" endDateTime="2019-02-19T00:00:00"> <or114:ProductName type="HHL"/> <or114:Hotel id="36" sequence="1" isPast="false"> <or114:Reservation DayOfWeekInd="1" NumberInParty="02"> <or114:LineNumber>1</or114:LineNumber> <or114:LineType>HHL</or114:LineType> <or114:LineStatus>SS</or114:LineStatus> <or114:POSRequestorID>46595076</or114:POSRequestorID> <or114:RoomType> <or114:RoomTypeCode>CSP</or114:RoomTypeCode> <or114:NumberOfUnits>1</or114:NumberOfUnits> <or114:ShortText>CSP08PG</or114:ShortText> </or114:RoomType> <or114:RoomRates> <or114:AmountBeforeTax>450.00</or114:AmountBeforeTax> <or114:CurrencyCode>EUR</or114:CurrencyCode> </or114:RoomRates> <or114:GuestCounts> <or114:GuestCount>2</or114:GuestCount> <or114:ExtraGuestCount>0</or114:ExtraGuestCount> <or114:RollAwayCount>0</or114:RollAwayCount> <or114:CribCount>0</or114:CribCount> </or114:GuestCounts> <or114:TimeSpanStart>2019-02-18T00:00:00</or114:TimeSpanStart> <or114:TimeSpanDuration>1</or114:TimeSpanDuration> <or114:TimeSpanEnd>2019-02-19T00:00:00</or114:TimeSpanEnd> <or114:Guarantee> <or114:Text>GVI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE</or114:Text> </or114:Guarantee> <or114:ChainCode>IC</or114:ChainCode> <or114:HotelCode>1098</or114:HotelCode> <or114:HotelCityCode>AMS</or114:HotelCityCode> <or114:HotelName>INTERCONTINENTAL AMSTEL AMS</or114:HotelName> <or114:HotelTotalPricing> <or114:TotalTax Amount="27.00"> <or114:Tax Id="1">27.00 SERVICE CHARGE</or114:Tax> </or114:TotalTax> <or114:ApproximateTotal AmountAndCurrency="477.00 EUR"/> <or114:Disclaimer Id="1">INCLUDES TAX</or114:Disclaimer> </or114:HotelTotalPricing> </or114:Reservation> <or114:AdditionalInformation> <or114:Address> <or114:AddressLine>PROFESSOR TULPPLEIN 1</or114:AddressLine> <or114:AddressLine>AMSTERDAM NL 1018 GX</or114:AddressLine> <or114:CountryCode>NL</or114:CountryCode> <or114:City>AMSTERDAM</or114:City> <or114:ZipCode>1018 GX</or114:ZipCode> </or114:Address> <or114:ContactNumbers> <or114:PhoneNumber>31-20-6226060</or114:PhoneNumber> <or114:FaxNumber>31-20-6225808</or114:FaxNumber> </or114:ContactNumbers> <or114:CancelPenaltyPolicyCode>01D</or114:CancelPenaltyPolicyCode> <or114:Commission> <or114:Indicator>C</or114:Indicator> <or114:Text>COMMISSIONABLE</or114:Text> </or114:Commission> </or114:AdditionalInformation> <or114:SegmentText>RR450.00EUR/RC-@@@-@@@-@/TTX-27.00/HTP-477.00 EUR/TX1-27.00 SERVICE CHARGE/DS1-INCLUDES TAX/CMN-C/CMT-COMMISSIONABLE/AGT-46595076/G-VI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE/C-01D/XS-0001548683004531675D5EA000000000/XT-675D5EA0/XL-0000/DT-28JAN190463/SBS-HS/HSA-PROFESSOR TULPPLEIN 1$AMSTERDAM NL 1018 GX/HFN-INTERCONTINENTAL AMSTEL AMS/HCY-AMSTERDAM/HST-/HCC-NL/HPC-1018 GX/HPH-31-20-6226060/HFX-31-20-6225808/UID-D599525312BD5487675D5E/SI-/CF-</or114:SegmentText> <or114:RateDescription> <or114:TextLine>BREAKFAST FOR 2 ADULTS FULLY</or114:TextLine> <or114:TextLine>EXECUTIVE CITY VIEW ROOM 35 SQM OR 377 SQFT.THE ROOMS</or114:TextLine> <or114:TextLine>OVERLOOK THE SQUARE.THEY ARE LAID OUT IN CLASSIC FRENCH</or114:TextLine> </or114:RateDescription> <or114:HotelPolicy> <or114:GuaranteePolicy>REQUIRED</or114:GuaranteePolicy> <or114:CancellationPolicy>CANCEL 1 DAYS PRIOR TO ARRIVAL</or114:CancellationPolicy> </or114:HotelPolicy> </or114:Hotel> </or114:ProductDetails> </stl19:Product> </stl19:Segment> </stl19:Segments> <stl19:TicketingInfo/> <stl19:ItineraryPricing/> </stl19:PassengerReservation> <stl19:ReceivedFrom> <stl19:Name>IBE</stl19:Name> </stl19:ReceivedFrom> <stl19:Addresses> <stl19:Address> <stl19:AddressLines> <stl19:AddressLine id="6" type="O"> <stl19:Text>CHTEAU DE MONTAIGNE</stl19:Text> </stl19:AddressLine> <stl19:AddressLine id="7" type="O"> <stl19:Text>RUE 123</stl19:Text> </stl19:AddressLine> <stl19:AddressLine id="8" type="O"> <stl19:Text>GUYENNE,, CA FR</stl19:Text> </stl19:AddressLine> <stl19:AddressLine id="9" type="O"> <stl19:Text>90210</stl19:Text> </stl19:AddressLine> </stl19:AddressLines> </stl19:Address> </stl19:Addresses> <stl19:PhoneNumbers> <stl19:PhoneNumber id="5" index="1" elementId="pnr-5"> <stl19:CityCode>SLC</stl19:CityCode> <stl19:Number>801-428-1231-H-1.1</stl19:Number> </stl19:PhoneNumber> </stl19:PhoneNumbers> <stl19:EmailAddresses/> <stl19:GenericSpecialRequests id="47" type="A" msgType="O"> <stl19:FreeText>HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED</stl19:FreeText> <stl19:AirlineCode>AA</stl19:AirlineCode> <stl19:FullText>AA HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="48" type="A" msgType="O"> <stl19:FreeText>HHL IC CXL AFTER 1800 17FEB FORFEIT ONE NITE STAY</stl19:FreeText> <stl19:AirlineCode>AA</stl19:AirlineCode> <stl19:FullText>AA HHL IC CXL AFTER 1800 17FEB FORFEIT ONE NITE STAY</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="37" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG D BREAKFAST FOR 2 ADULTS FULLY</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG D BREAKFAST FOR 2 ADULTS FULLY</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="38" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG C CANCEL 1 DAYS PRIOR TO ARRIVAL</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG C CANCEL 1 DAYS PRIOR TO ARRIVAL</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="39" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG D EXECUTIVE CITY VIEW ROOM 35 SQ</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG D EXECUTIVE CITY VIEW ROOM 35 SQ</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="40" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG D OVERLOOK THE SQUARE.THEY ARE L</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG D OVERLOOK THE SQUARE.THEY ARE L</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="41" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG TTX 27.00 TTL TAX</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG TTX 27.00 TTL TAX</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="42" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG TX 27.00 SERVICE CHARGE</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG TX 27.00 SERVICE CHARGE</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="43" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG TP 477.00 EUR APPROX. TTL PRICE</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG TP 477.00 EUR APPROX. TTL PRICE</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:GenericSpecialRequests id="44" type="G" msgType="O"> <stl19:FreeText>HHL    1098 CSP08PG TD INCLUDES TAX</stl19:FreeText> <stl19:AirlineCode>IC</stl19:AirlineCode> <stl19:FullText>IC HHL    1098 CSP08PG TD INCLUDES TAX</stl19:FullText> </stl19:GenericSpecialRequests> <stl19:AssociationMatrices> <stl19:AssociationMatrix> <stl19:Name>PssIDType</stl19:Name> <stl19:Parent ref="pnr-36"/> <stl19:Child ref="pnr-or-3"/> </stl19:AssociationMatrix> </stl19:AssociationMatrices> <stl19:OpenReservationElements> <or114:OpenReservationElement id="47" type="SRVC" elementId="pnr-47"> <or114:ServiceRequest airlineCode="AA" serviceType="OSI" ssrType="AFX"> <or114:FreeText>HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED</or114:FreeText> <or114:FullText>AA HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="48" type="SRVC" elementId="pnr-48"> <or114:ServiceRequest airlineCode="AA" serviceType="OSI" ssrType="AFX"> <or114:FreeText>HHL IC CXL AFTER 1800 17FEB FORFEIT ONE NITE STAY</or114:FreeText> <or114:FullText>AA HHL IC CXL AFTER 1800 17FEB FORFEIT ONE NITE STAY</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="37" type="SRVC" elementId="pnr-37"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG D BREAKFAST FOR 2 ADULTS FULLY</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG D BREAKFAST FOR 2 ADULTS FULLY</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="38" type="SRVC" elementId="pnr-38"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG C CANCEL 1 DAYS PRIOR TO ARRIVAL</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG C CANCEL 1 DAYS PRIOR TO ARRIVAL</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="39" type="SRVC" elementId="pnr-39"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG D EXECUTIVE CITY VIEW ROOM 35 SQ</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG D EXECUTIVE CITY VIEW ROOM 35 SQ</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="40" type="SRVC" elementId="pnr-40"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG D OVERLOOK THE SQUARE.THEY ARE L</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG D OVERLOOK THE SQUARE.THEY ARE L</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="41" type="SRVC" elementId="pnr-41"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG TTX 27.00 TTL TAX</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG TTX 27.00 TTL TAX</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="42" type="SRVC" elementId="pnr-42"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG TX 27.00 SERVICE CHARGE</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG TX 27.00 SERVICE CHARGE</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="43" type="SRVC" elementId="pnr-43"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG TP 477.00 EUR APPROX. TTL PRICE</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG TP 477.00 EUR APPROX. TTL PRICE</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> <or114:OpenReservationElement id="44" type="SRVC" elementId="pnr-44"> <or114:ServiceRequest airlineCode="IC" serviceType="OSI" ssrType="GFX"> <or114:FreeText>HHL    1098 CSP08PG TD INCLUDES TAX</or114:FreeText> <or114:FullText>IC HHL    1098 CSP08PG TD INCLUDES TAX</or114:FullText> </or114:ServiceRequest> </or114:OpenReservationElement> </stl19:OpenReservationElements> </stl19:Reservation> </stl19:GetReservationRS> </soap-env:Body> </soap-env:Envelope>`)
)

func TestGetResSegmentRSUnmarshal(t *testing.T) {
	getRes := GetReservationResponse{}
	err := xml.Unmarshal(sampleGetResSegmentRS, &getRes)
	if err != nil {
		t.Fatal("Error unmarshal get reservation response", err)
	}
}

func TestGetResReservationRSBasic(t *testing.T) {
	getRes := GetReservationResponse{}
	_ = xml.Unmarshal(sampleGetResSegmentRS, &getRes)
	res := getRes.Body.GetReservationRS.Reservation

	if res.ReceivedFrom.Name != "IBE" {
		t.Errorf("ReceivedFrom.Name wrong, expect: %v, got: %v", "IBE", res.ReceivedFrom.Name)
	}
	var space = "http://webservices.sabre.com/pnrbuilder/v1_19"
	psngr := Passenger{
		XMLName:         xml.Name{Space: space, Local: "Passenger"},
		ID:              "3",
		NameType:        "S",
		PassengerType:   "ADT",
		ReferenceNumber: "ABC123",
		NameID:          "01.01",
		NameAssocID:     "1",
		ElementID:       "pnr-3.1",
		LastName:        "DE MONTAIGNE",
		FirstName:       "MICHEL",
	}
	if res.PassengerReservation.Passengers[0] != psngr {
		t.Errorf("res.PassengerReservation.Passengers[0] \nexp: %v \ngot: %v", psngr, res.PassengerReservation.Passengers[0])
	}
	addr := AddressReservationElem{
		XMLName: xml.Name{Space: space, Local: "Address"},
		AddressLines: AddressLineResElem{
			XMLName: xml.Name{Space: space, Local: "AddressLine"},
			ID:      "9",
			Atype:   "O",
			Text:    "90210",
		},
	}
	if res.Addresses[0] != addr {
		t.Errorf("Addresses \nexp: %v \ngot: %v", addr, res.Addresses[0])
	}

	phn := PhoneNumberReservationElem{
		XMLName:   xml.Name{Space: space, Local: "PhoneNumber"},
		ID:        "5",
		Index:     "1",
		ElementID: "pnr-5",
		CityCode:  "SLC",
		Number:    "801-428-1231-H-1.1",
	}
	if res.PhoneNumbers[0] != phn {
		t.Errorf("PhoneNumbers \nexp: %v \ngot: %v", phn, res.PhoneNumbers[0])
	}

	gnrsp := GenericSpecialRequests{
		XMLName:     xml.Name{Space: space, Local: "GenericSpecialRequests"},
		ID:          "47",
		GType:       "A",
		MsgType:     "O",
		FreeText:    "HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED",
		AirlineCode: "AA",
		FullText:    "AA HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED",
	}
	if res.GenericSpecialRequests[0] != gnrsp {
		t.Errorf("GenericSpecialRequests \nexp: %v \ngot: %v", gnrsp, res.GenericSpecialRequests[0])
	}
	if res.AssociationMatrices[0].Name != "PssIDType" {
		t.Errorf("AssociationMatrix \nexp: %v \ngot: %v", "PssIDType", res.AssociationMatrices[0].Name)
	}

	ore := OpenReservationElement{
		XMLName:   xml.Name{Space: space, Local: "OpenReservationElement"},
		ID:        "47",
		SType:     "A",
		ElementID: "pnr-5",
		ServiceRequest: ServiceRequestOpenRes{
			XMLName:     xml.Name{Space: "http://services.sabre.com/res/or/v1_14", Local: "ServiceRequest"},
			AirlineCode: "AA",
			ServiceType: "OSI",
			SsrType:     "AFX",
			FreeText:    "HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED",
			FullText:    "AA HHL IC TOTAL: 477.00 EUR ALL KNOWN FEES INCLUDED",
		},
	}
	if res.OpenReservationElements[0].ServiceRequest != ore.ServiceRequest {
		t.Errorf("OpenReservationElement \nexp: %v \ngot: %v", ore.ServiceRequest, res.OpenReservationElements[0].ServiceRequest)
	}
}

func genSampleHotelSegTest() HotelSegmentElem {
	var or14 = "http://services.sabre.com/res/or/v1_14"
	return HotelSegmentElem{
		XMLName:  xml.Name{Space: or14, Local: "Hotel"},
		ID:       "36",
		Sequence: "1",
		IsPast:   false,
		Reservation: ReservatioHotel{
			XMLName:       xml.Name{Space: or14, Local: "Reservation"},
			DayOfWeekInd:  "1",
			NumberInParty: "02",
			LineNumber:    "1",
			LineType:      "HHL",
			LineStatus:    "SS",
			POSRequestor:  "46595076",
			RoomType: RoomTypeRes{
				XMLName:       xml.Name{Space: or14, Local: "RoomType"},
				RoomTypeCode:  "CSP",
				NumberOfUnits: "1",
				ShortText:     "CSP08PG",
			},
			RoomRates: RoomRatesRes{
				XMLName:         xml.Name{Space: or14, Local: "RoomRates"},
				AmountBeforeTax: "450.00",
				CurrencyCode:    "EUR",
			},
			TimeSpanStart:    "2019-02-18T00:00:00",
			TimeSpanDuration: "1",
			TimeSpanEnd:      "2019-02-19T00:00:00",
			Guarantee:        GuaranteeRes{Text: "GVI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE"},
		},
		AdditionalInformation: AdditionalInformation{
			CancelPenaltyPolicyCode: "01D",
			Commission: CommissionAdditional{
				Indicator: "C",
				Text:      "COMMISSIONABLE",
			},
		},
		SegmentText: "450.00EUR/RC-@@@-@@@-@/TTX-27.00/HTP-477.00 EUR/TX1-27.00 SERVICE CHARGE/DS1-INCLUDES TAX/CMN-C/CMT-COMMISSIONABLE/AGT-46595076/G-VI4XXXXXXXXXXX4501EXP 10 20-DE MONTAIGNE/C-01D/XS-0001548683004531675D5EA000000000/XT-675D5EA0/XL-0000/DT-28JAN190463/SBS-HS/HSA-PROFESSOR TULPPLEIN 1$AMSTERDAM NL 1018 GX/HFN-INTERCONTINENTAL AMSTEL AMS/HCY-AMSTERDAM/HST-/HCC-NL/HPC-1018 GX/HPH-31-20-6226060/HFX-31-20-6225808/UID-D599525312BD5487675D5E/SI-/CF-23164323-",
		//RateDescription: RateDescriptionSegment{},
		/*
			ONLY ON PRODUCT
		*/
		HotelPolicy: HotelPolicySegment{
			GuaranteePolicy:    "REQUIRED",
			CancellationPolicy: "CANCEL 1 DAYS PRIOR TO ARRIVAL",
		},
	}
}

func TestGetResSegmentHotelRSBasic(t *testing.T) {
	getRes := GetReservationResponse{}
	_ = xml.Unmarshal(sampleGetResSegmentRS, &getRes)
	hot := getRes.Body.GetReservationRS.Reservation.PassengerReservation.Segments.Hotel
	exp := genSampleHotelSegTest()
	if hot.ID != exp.ID {
		t.Errorf("Hotel.Id exp: %s, got: %s", exp.ID, hot.ID)
	}
	if hot.IsPast != exp.IsPast {
		t.Errorf("Hotel.IsPast exp: %t, got: %t", exp.IsPast, hot.IsPast)
	}
	if hot.Reservation.RoomType != exp.Reservation.RoomType {
		t.Errorf("RoomTypeCode \nexp: %+v \ngot: %+v", exp.Reservation.RoomType, hot.Reservation.RoomType)
	}
	if hot.AdditionalInformation.CancelPenaltyPolicyCode != exp.AdditionalInformation.CancelPenaltyPolicyCode {
		t.Errorf("CancelPenaltyPolicyCode \nexp: %+v \ngot: %+v", exp.AdditionalInformation.CancelPenaltyPolicyCode, hot.AdditionalInformation.CancelPenaltyPolicyCode)
	}

}
