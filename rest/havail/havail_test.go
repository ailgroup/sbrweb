package havail

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"
)

var (
	sabreURLexpect      = "https://api.havail.sabre.com/v1.0.0"
	geoURLexpect        = "https://api.havail.sabre.com/v1.0.0/lists/utilities/geosearch"
	geoLocatonURLexpect = "https://api.havail.sabre.com/v1.0.0/lists/utilities/geosearch/locations?mode=geosearch"
)

func TestRestSabreURL(t *testing.T) {
	if sabreURL().String() != sabreURLexpect {
		t.Errorf("expected: %s, got %s", sabreURLexpect, sabreURL().String())
	}
}

func TestRestGeoURL(t *testing.T) {
	if geoURL().String() != geoURLexpect {
		t.Errorf("expected: %s, got %s", geoURLexpect, geoURL().String())
	}
}

func TestGeoLocations(t *testing.T) {
	if GeoLocations().String() != geoLocatonURLexpect {
		t.Errorf("expected: %s, got %s", geoLocatonURLexpect, GeoLocations().String())
	}
}

func TestEnvVariables(t *testing.T) {
	cenv := getClientID()
	cdef := "SabreClientIDNotFound"
	if cenv == cdef {
		t.Error("ClientID not correct, check SABRE_CLIENT_ID")
	}

	csenv := getClientSecret()
	csdef := "SabreClientSecretNotFound"
	if csenv == csdef {
		t.Error("ClientSecret not correct, check SABRE_CLIENT_SECRET")
	}
}

func TestMakeToken(t *testing.T) {
	tok := makeBasicToken()
	matchPrefix, _ := regexp.MatchString(`Basic`, tok)
	if !matchPrefix {
		t.Error("Prefix to Token is not 'Basic'")
	}
	fmt.Printf("basic token: %v\n", tok)
	cid := base64.StdEncoding.EncodeToString([]byte(getClientID()))
	secret := base64.StdEncoding.EncodeToString([]byte(getClientSecret()))
	concat := fmt.Sprintf("%s:%s", cid, secret)
	matchCid, _ := regexp.MatchString(cid, concat)
	if !matchCid {
		t.Error("Token does not contain b64 client id")
	}
	matchSecret, _ := regexp.MatchString(secret, concat)
	if !matchSecret {
		t.Error("Token does not contain b64 secret")
	}
}

/*
func TestNewSabreClient(t *testing.T) {
	g := GeoLocations()
	c := NewClient(authTokenBearerPrefix, g)

	if c.EndpointURL.String() != geoLocatonURLexpect {
		t.Errorf("expected: %s, got %s", geoLocatonURLexpect, c.EndpointURL.String())
	}
}
*/

/*
{"access_token":"T1RLAQIaXCq0nX34KDdhorrcJVZ3kbLuxBCfIQMNHEn0PgtIIySO6GbEAADAeFqE2VZD8ItHUseeeSkixuzRsSVqij7GIFOuPfKdvi3uyeZcMvRqTZhPLmkpf2GV1M99MBN5MHknIbLnmq4Au65Ece1EpjlAx6BecTToCTqOtzXlfg084eWPXZDazlD1EX4/l7OK+37+9C4kVH0UqMvZgGjzsZFuj0qoOaWeSBMdZCzvF/fDsxhB5L2GXiEsN3ByL11jXtUA65XZFlD5ZiKfO9KdupCUachhcBnN7eUoEjDpvS5bPrDl/JLU4SAV","token_type":"bearer","expires_in":604800}
*/
