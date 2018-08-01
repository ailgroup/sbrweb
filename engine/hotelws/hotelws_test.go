package hotelws

import "testing"

var dirtyStrings = []struct {
	dirty string
	clean string
}{
	{"slash/ whitespace", "slashwhitespace"},
	{"boo/this/bad boo", "boothisbadboo"},
}

func TestSanatizeString(t *testing.T) {
	for _, ds := range dirtyStrings {
		clean := sanatize(ds.dirty)
		if clean != ds.clean {
			t.Errorf("expected %s, got %s", ds.clean, clean)
		}
	}
}
