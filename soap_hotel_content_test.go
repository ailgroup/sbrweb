package sbrweb

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestContentRQHasFormattedIDS(t *testing.T) {
	ids := []string{"0000001", "002", "03", "000124"}
	con, _ := BuildGetContentRequest(1, ids)
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

	b, err := xml.Marshal(content)
	if err != nil {
		t.Error("Error marshaling get hotel content", err)
	}
	fmt.Printf("content marshal \n%s\n", b)
}
