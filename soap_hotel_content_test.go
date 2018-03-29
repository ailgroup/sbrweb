package sbrweb

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestBuildGetHotelContentRequestMarshal(t *testing.T) {
	ids := []string{"1", "002", "03"}
	content, _ := BuildGetContentRequest(10, ids)
	fmt.Printf("content %+v\n", content)
	if content.XMLNSXsi != baseXsiNamespace {
		t.Errorf("GetHotelContent XMLNSXsi expect: %s, got %s", baseXsiNamespace, content.XMLNSXsi)
	}
	if content.XMLNSSchemaLocation != baseGetHotelContentSchema {
		t.Errorf("GetHotelContent XMLNSSchemaLocation expect: %s, got %s", baseGetHotelContentSchema, content.XMLNSSchemaLocation)
	}

	if len(content.SearchCriteria.HotelRefs) != len(ids) {
		t.Errorf("HotelRefs '%d' not size of input ids '%d'. expect %v, got %v", len(content.SearchCriteria.HotelRefs), len(ids), ids, content.SearchCriteria.HotelRefs)

	}

	b, err := xml.Marshal(content)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	fmt.Printf("content marshal \n%s\n", b)
}
