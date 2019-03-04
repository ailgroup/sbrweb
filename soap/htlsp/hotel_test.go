package htlsp

import (
	"errors"
	"testing"
)

var dirtyStringsSample = []struct {
	dirty string
	clean string
}{
	{"slash/ whitespace", "slashwhitespace"},
	{"boo/this/bad boo", "boothisbadboo"},
}

func TestSanatizeString(t *testing.T) {
	for _, ds := range dirtyStringsSample {
		clean := sanatize(ds.dirty)
		if clean != ds.clean {
			t.Errorf("expected %s, got %s", ds.clean, clean)
		}
	}
}

var b64EncodeStringsSample = []struct {
	in  string
	enc string
}{
	{"hello smiley face", "aGVsbG8gc21pbGV5IGZhY2U="},
	{"1234 {\"this\"} and-/that", "MTIzNCB7InRoaXMifSBhbmQtL3RoYXQ="},
}

func TestB64Enc(t *testing.T) {
	for idx, sample := range b64EncodeStringsSample {
		encoded := B64Enc(sample.in)
		if encoded != sample.enc {
			t.Errorf("for sample %d: expected %s, got %s", idx, sample.enc, encoded)
		}
	}
}

var b64DecodeStringsSample = []struct {
	in  string
	dec string
	err error
}{
	{"", "", nil},
	{"MTIzNCB7InRoaXMifSBhbmQtL3RoYXQ=", "1234 {\"this\"} and-/that", nil},
	{"aGVsbG8gc21pbGV5IGZhY2U=", "hello smiley face", nil},
	//not url safe!
	{"c29tZSBkYXRhIHdpdGggACBhbmQg77u/", "some data with \x00 and \ufeff", errors.New("illegal base64 data at input byte 31")},
}

func TestB64Dec(t *testing.T) {
	for idx, sample := range b64DecodeStringsSample {
		decoded, err := B64Dec(sample.in)
		if err != nil {
			if err.Error() != sample.err.Error() {
				t.Errorf("for sample %d: expected %s, got %s", idx, sample.err, err)
			}
			continue
		}
		if decoded != sample.dec {
			t.Errorf("for sample %d: expected %s, got %s", idx, sample.dec, decoded)
		}
	}
}
