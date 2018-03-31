package sbrhotel

/*
DEPRECATE: use rest instead
import (
	"encoding/xml"
	"testing"
)


var (
	dt                      = DescriptionTypes{diningTypeField, alertTypeField}
	dq                      = make(DescriptiveQuery)
	sampleGetHotelContentRQ = []byte(`<GetHotelContentRQ xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:schemaLocation="http://services.sabre.com/hotel/content/v1 GetHotelContentRQ.xsd" version="1.0.0"><SearchCriteria><ImageRef MaxImages="10"></ImageRef><HotelRefs><HotelRef HotelCode="1"></HotelRef><HotelRef HotelCode="002"></HotelRef><HotelRef HotelCode="03"></HotelRef></HotelRefs><DescriptiveInfoRef><PropertyInfo>true</PropertyInfo><LocationInfo>true</LocationInfo><Amenities>true</Amenities><Airports>true</Airports><AcceptedCreditCards>true</AcceptedCreditCards><Descriptions><Description Type="Dining"></Description><Description Type="Alerts"></Description></Descriptions></DescriptiveInfoRef></SearchCriteria></GetHotelContentRQ>`)
)

func init() {
	dq[propertyQueryField] = true
	dq[locationQueryField] = true
	dq[amenityQueryField] = true
	dq[airportQueryField] = true
	dq[creditQueryField] = true
}

func TestBuildDescriptions(t *testing.T) {
	di := buildDescriptions(dq, dt)

	if di.Descriptions[0].Type != dt[0] {
		t.Errorf("DescriptiveInfo.Descriptions[0].Type expect %s, got %s", diningTypeField, di.Descriptions[0].Type)
	}
	if di.Descriptions[1].Type != dt[1] {
		t.Errorf("DescriptiveInfo.Descriptions[1].Type expect %s, got %s", alertTypeField, di.Descriptions[1].Type)
	}

	if !di.Property {
		t.Error("DescriptiveInfo.Property should be true")
	}
	if !di.Location {
		t.Error("DescriptiveInfo.Location should be true")
	}
	if !di.Amenities {
		t.Error("DescriptiveInfo.Amenities should be true")
	}
	if !di.Airports {
		t.Error("DescriptiveInfo.Airport should be true")
	}
	if !di.AcceptedCreditCards {
		t.Error("DescriptiveInfo.AcceptedCreditCards should be true")
	}
}

func TestBuildGetHotelContentReturnError(t *testing.T) {
	ids := []string{}
	_, err := BuildGetHotelContent(1, ids, DescriptiveInfo{})
	if err == nil {
		t.Error("When passing no hotel ids to BuildGetContentRequest we should have error")
	}
}

func TestBuildGetHotelContentIDS(t *testing.T) {
	ids := []string{"0000001", "002", "03", "000124"}
	con, err := BuildGetHotelContent(1, ids, DescriptiveInfo{})
	if err != nil {
		t.Error("BuildGetContentRequest should not have error")
	}
	refs := con.SearchCriteria.HotelRefs

	if len(refs) != len(ids) {
		t.Errorf("HotelRefs '%d' not size of input ids '%d'. expect %v, got %v", len(refs), len(ids), ids, refs)
	}
	for i, id := range ids {
		if id != refs[i].HotelCode {
			t.Error("HotelRef id not properly formatted. expect", id, "got", refs[i].HotelCode)
		}
	}
}

func TestBuildGetHotelContentMarshal(t *testing.T) {
	ids := []string{"1", "002", "03"}
	content, _ := BuildGetHotelContent(10, ids, buildDescriptions(dq, dt))
	if content.XMLNSXsi != baseXsiNamespace {
		t.Errorf("BuildGetContentRequest XMLNSXsi expect: %s, got %s", baseXsiNamespace, content.XMLNSXsi)
	}
	if content.XMLNSSchemaLocation != baseGetHotelContentSchema {
		t.Errorf("BuildGetContentRequest XMLNSSchemaLocation expect: %s, got %s", baseGetHotelContentSchema, content.XMLNSSchemaLocation)
	}

	b, err := xml.Marshal(content)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	if string(b) != string(sampleGetHotelContentRQ) {
		t.Errorf("Expected marshal get hotel content \n sample: %s \n result: %s", string(sampleGetHotelContentRQ), string(b))
	}
	//fmt.Printf("content marshal \n%s\n", b)
}
*/
