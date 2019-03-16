package havail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	v2DevAuthTokenURL       = "https://api-crt.cert.havail.sabre.com/v2/auth/token"
	authTokenBasicPrefix    = "Basic"
	tickerTimeDefaultSecond = time.Duration(1200 * time.Second)
)

//BasicAuthRS accepts and parses response from sabre for an acces token.
// Expiration is: 7 days; 168 hours; 604800 seconds.
type BasicAuthRS struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int32  `json:"expires_in"`
}

// SabreClient contains values for a common client to Sabre rest endpoints.
type SabreClient struct {
	BasicAuthRS *BasicAuthRS
	ticker      *time.Ticker
}

// GetBasicAuthToken requests a 7day access_token from sabre
func (s *SabreClient) SetAccessToken() error {
	httpClient := &http.Client{}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequest("POST", v2DevAuthTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", makeBasicToken())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &s.BasicAuthRS)
	if err != nil {
		return err
	}
	return nil
}

// superviseAccessToken runs in the background to refresh the access_token.
// It has not graceful shutdown nor does it stop the ticker; it is expected to run
// during the life of the program.
func (s *SabreClient) superviseAccessToken() {
	s.setTicker()
	go func(tckr *time.Ticker) {
		defer tckr.Stop()
		//for range tckr.C {
		for t := range tckr.C {
			fmt.Println("tick at", t)
			//debug
			fmt.Printf("\nOld access token %s\n\n", s.BasicAuthRS.AccessToken)
			//TODO handle error??
			_ = s.SetAccessToken()
			//debug
			fmt.Printf("New access token %s\n\n", s.BasicAuthRS.AccessToken)
		}
	}(s.ticker)
}

// setTicker defines the time duration to wait before refreshing the rest api access token.
func (s *SabreClient) setTicker() {
	if s.BasicAuthRS.ExpiresIn < 100 {
		s.ticker = time.NewTicker(tickerTimeDefaultSecond * time.Second)
		return
	}
	//debug....
	val := percentOf(0.02, s.BasicAuthRS.ExpiresIn)
	fmt.Printf("Expires in %v, ticker set to %v\n", s.BasicAuthRS.ExpiresIn, val*time.Second)
	s.ticker = time.NewTicker(val * time.Second)

	//...production
	//0.02 == 2minutes; 0.2 == 20minutes
	// wait 90% of alloted time before refreshing token
	//s.ticker = time.NewTicker(percentOf(90.0, s.BasicAuthRS.ExpiresIn) * time.Second)
}

// percentOf calculates n percent of d.
func percentOf(n float64, d int32) time.Duration {
	val := time.Duration((n * float64(d)) / 100.0)
	if val < 100 {
		return time.Duration(1000 * time.Second)
	}
	return val
}

// NewClient returns SabreClient with for endpoint. Each client will share access token from
// BasicAuthRS
func NewClient() (*SabreClient, error) {
	c := SabreClient{}
	err := c.SetAccessToken()
	if err != nil {
		return &c, err
	}
	c.superviseAccessToken()
	return &c, nil
}

// makeToken builds the api token for Sabre according to instructions.
// See https://developer.sabre.com/resources/getting_started_with_sabre_apis/sabre_apis_101/how_to_guides/rest_apis_token_credentials.
func makeBasicToken() string {
	cid := base64.StdEncoding.EncodeToString([]byte(getClientID()))
	secret := base64.StdEncoding.EncodeToString([]byte(getClientSecret()))
	encString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cid, secret)))
	return fmt.Sprintf("%s %s", authTokenBasicPrefix, encString)
}

// getClientID sets the ENV variable else defaults.
// TODO update to read from a config file...
func getClientID() string {
	v, ok := os.LookupEnv("SABRE_CLIENT_ID")
	if !ok {
		return "SabreClientIDNotFound"
	}
	return v
}

// getClientSecret sets the ENV variable else defaults.
// TODO update to read from a config file...
func getClientSecret() string {
	v, ok := os.LookupEnv("SABRE_CLIENT_SECRET")
	if !ok {
		return "SabreClientSecretNotFound"
	}
	return v
}
