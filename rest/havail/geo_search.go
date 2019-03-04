package havail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

/*
https://developer.sabre.com/docs/read/rest_apis/utility/geo_search

In detail: the API has three different types of location resolution: airport or multi-airport city (MAC) code, geo code (latitude and longitude) and place name, e.g., town or city name, with optional date and country code. Included in the request is also a parameter for defining the required radius search of up to 200 miles or kilometers along with the defined unit of measurement (miles or kilometers). The response will be the identified hotels (Sabre Property ID only) and the straight line distance between the resolved location and the designated Hotel.


Method/Endpoint
	POST /v1.0.0/lists/utilities/geosearch/locations?mode=geosearch HTTP/1.1 HTTP/1.1
	POST https://api.havail.sabre.com/v1.0.0/lists/utilities/geosearch/locations?mode=geosearch
*/

var (
	geoURLPart = "lists/utilities/geosearch"
)

// GeoSearchRQ holds options for geo search requests.
// See https://developer.sabre.com/docs/read/rest_apis/utility/geo_search.
//See http://files.developer.sabre.com/doc/providerdoc/STPS/geo-search/v100/geo-search-v100-request.jsonschema.
type GeoSearchRequest struct {
	GeoSearchRQ GeoSearchRQ
}
type GeoSearchRQ struct {
	Version string `json:"version"`
	GeoRef  GeoRef //`json:"GeoRef"`
}

//GeoRef holds params for geo search
type GeoRef struct {
	Category   CategoryT         `json:",omitempty"` //HOTEL
	UOM        UOMT              `json:",omitempty"` //"MI"
	Radius     RadiusT           `json:",omitempty"` //1.0
	MaxResults MaxSearchResultsT `json:",omitempty"` //300
	OffSet     OffsetT           `json:",omitempty"` //1
	HTTPVerb   string            `json:"-"`          //POST
	Endpoint   EndpointFunc      `json:"-"`          //GeoLocations
	AddressRef AddressRef
}

//AddressRef holds params for address in geo search
type AddressRef struct {
	Street      StreetT      `json:",omitempty"`
	City        CityT        `json:",omitempty"`
	County      string       `json:",omitempty"`
	PostalCode  string       `json:",omitempty"`
	StateProv   StateT       `json:",omitempty"`
	CountryCode CountryCodeT `json:",omitempty"`
}

// GeoSearchResult holds property specific data.
type GeoSearchResult struct {
	Distance  float64
	Latitude  LatitudeT
	Longitude LongitudeT
	Name      string
	ID        string `json:"Id"`
	Street    StreetT
	Zip       ZipT
	City      CityT
	State     StateT
	Country   CountryCodeT
	Attribute AttributeT
}

// GeoSearchResults holds meta-data for the query.
type GeoSearchResults struct {
	Category         CategoryT
	UOM              UOMT
	Radius           RadiusT
	MaxSearchResults MaxSearchResultsT
	OffSet           OffsetT
	Latitude         LatitudeT
	Longitude        LongitudeT
	GeoSearchResult  []GeoSearchResult
}

type GeoSearchRS struct {
	ApplicationResults ApplicationResults
	GeoSearchResults   GeoSearchResults
}
type GeoSearchResponse struct {
	GeoSearchRS GeoSearchRS
}

//func (c *SabreClient) GeoSearchFor(ref GeoRef) (GeoSearchRS, error) {
func (c *SabreClient) GeoSearchFor(srch GeoSearchRequest) (*GeoSearchResponse, error) {
	fmt.Printf("\n\nGEO-REF %+v \n\n", srch)
	httpClient := &http.Client{}
	reqByte, _ := json.Marshal(srch)
	fmt.Printf("\n\nJSON_MARSHAL %s \n\n", reqByte)
	req, _ := http.NewRequest(
		srch.GeoSearchRQ.GeoRef.HTTPVerb,
		srch.GeoSearchRQ.GeoRef.Endpoint().String(),
		bytes.NewBuffer(reqByte),
	)
	req.Header.Add("Authorization", strings.Join([]string{c.BasicAuthRS.TokenType, c.BasicAuthRS.AccessToken}, " "))
	req.Header.Add("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	geo := &GeoSearchResponse{}
	err = json.Unmarshal(body, geo)
	resp.Body.Close()
	return geo, err

	//fmt.Printf("ERROR? %v \n\n", err)
	//fmt.Printf("RESPONSE: %v\n", resp)
	// fmt.Printf("BODY: %v\n", body)
	//fmt.Printf("BODY: %s\n", body)
}

// GeoLocations returns url for geosearch with locations in Sabre rest,
// returns "https://api.havail.sabre.com/v1.0.0/lists/utilities/geosearch/locations?mode=geosearch"
func GeoLocations() *url.URL {
	geo := geoURL()
	geo.Path = path.Join(geo.Path, "locations")
	q := geo.Query()
	q.Set("mode", "geosearch")
	geo.RawQuery = q.Encode()
	return geo
}

// geoURL builds the base urlor making geo requests in Sabre rest.
// returns "https://api.havail.sabre.com/v1.0.0/lists/utilities/geosearch"
func geoURL() *url.URL {
	u := sabreURL()
	u.Path = path.Join(u.Path, geoURLPart)
	return u
}

/*

ERROR:
	{"status":"NotProcessed","reportingSystem":"RAF","timeStamp":"2019-03-07T01:42:43+00:00","type":"Validation","errorCode":"ERR.RAF.VALIDATION","instance":"raf-darhlc005.sabre.com-9080","message":"[{\"level\":\"error\",\"schema\":{\"loadingURI\":\"#\",\"pointer\":\"/definitions/com.sabre.services.util.geo.v1.GeoRef\"},\"instance\":{\"pointer\":\"/GeoSearchRQ/GeoRef\"},\"domain\":\"validation\",\"keyword\":\"additionalProperties\",\"message\":\"object instance has properties which are not allowed by the schema: [\\\"AddresRef\\\"]\",\"unwanted\":[\"AddresRef\"]}]"}

BODY:
	{"GeoSearchRS":{"ApplicationResults":{"Success":[{"timeStamp":"2019-03-06T19:45:28.423-06:00"}]},"GeoSearchResults":{"Radius":1.0,"UOM":"MI","Category":"HOTEL","Latitude":32.877416,"Longitude":-96.959879,"MaxSearchResults":17,"OffSet":1,"GeoSearchResult":[{"Distance":0.28,"Latitude":32.881352,"Longitude":-96.961089,"Name":"HOMESTEAD DALLAS-LAS COLINAS","Id":"42006","Street":"5315 CARNABY STREET","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.29,"Latitude":32.88159,"Longitude":-96.959395,"Name":"MARRIOTT EXECUSTAY BEAVER CREEK","Id":"86034","Street":"1000 Meadow Creek Drive","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.45,"Latitude":32.88242,"Longitude":-96.95483,"Name":"EXTENDEDSTAYDELUXE MEADOW CRK","Id":"42911","Street":"605 MEADOW CREEK DR","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.46,"Latitude":32.879841,"Longitude":-96.952458,"Name":"CANDLEWOOD SUITES DALLAS","Id":"52286","Street":"5300 GREEN PARK DRIVE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.52,"Latitude":32.884876,"Longitude":-96.959822,"Name":"WINGATE BY WYNDHAM LAS COLINAS","Id":"30960","Street":"850 W WALNUT HILL LANE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.52,"Latitude":32.88486,"Longitude":-96.96088,"Name":"TOWNEPLACE SUITES LAS COLINAS","Id":"44599","Street":"900 W WALNUT HILL LANE","Zip":"75038-2613","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.52,"Latitude":32.88486,"Longitude":-96.961448,"Name":"HAMPTON INN DALLAS LAS COLINAS","Id":"37676","Street":"820 WALNUT HILL LANE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.54,"Latitude":32.88486,"Longitude":-96.962509,"Name":"RESIDENCE INN LAS COLINAS","Id":"21274","Street":"950 WALNUT HILL LANE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.55,"Latitude":32.883969,"Longitude":-96.95447,"Name":"EXTENDEDSTAYDELUXE LAS COLINAS","Id":"43307","Street":"5401 GREEN PARK DRIVE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.58,"Latitude":32.884701,"Longitude":-96.954984,"Name":"HYATT PLACE DALLA LAS COLINAS","Id":"43702","Street":"5455 GREEN PARK DR","Zip":"75038","City":"IRVING","State":"TX","Country":"US","Attribute":{}},{"Distance":0.64,"Latitude":32.884293,"Longitude":-96.952418,"Name":"FAIRFIELD INN LAS COLINAS","Id":"41741","Street":"630 W JOHN CARPENTER FREEWAY","Zip":"75039","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.65,"Latitude":32.88505,"Longitude":-96.96631,"Name":"COURTYARD LAS COLINAS","Id":"17064","Street":"1151 W WALNUT HILL LANE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.71,"Latitude":32.88565,"Longitude":-96.96716,"Name":"STAYBRIDGE SUITES LA COLINAS","Id":"48961","Street":"1201 EXECUTIVE CIRCLE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.8,"Latitude":32.879379,"Longitude":-96.946315,"Name":"HOLIDAY INN EXP STES LAS COLI","Id":"41673","Street":"333 W JOHN CARPENTER FREEWAY","Zip":"75039","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.85,"Latitude":32.865174,"Longitude":-96.959835,"Name":"LA QUINTA IS LAS COLINAS","Id":"4665","Street":"4225 MACARTHUR BLVD","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.95,"Latitude":32.890227,"Longitude":-96.965772,"Name":"COMFORT SUITES LAS COLINAS","Id":"2238","Street":"1223 GREENWAY CIRCLE","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}},{"Distance":0.99,"Latitude":32.863111,"Longitude":-96.960096,"Name":"FOUR SEASONS DALLAS","Id":"12132","Street":"4150 N MACARTHUR BLVD","Zip":"75038","City":"Irving","State":"TX","Country":"US","Attribute":{}}]}},"Links":[{"rel":"self","href":"https://api-crt.cert.havail.sabre.com/v1.0.0/lists/utilities/geosearch/locations?mode=geosearch"},{"rel":"linkTemplate","href":"https://api-crt.cert.havail.sabre.com/<version>/lists/utilities/geosearch/locations?mode=<mode>"}]}
*/
