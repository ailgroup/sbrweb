// Package havail is a set of commmon functionality for interacting with Sabre hotel avail rest endpoints
package havail

import (
	"net/url"
	"path"
)

var (
	apiVersion       = "v1.0.0"
	baseDevHavailURL = "https://api-crt.cert.havail.sabre.com"
)

type EndpointFunc func() *url.URL
type UOMT string
type CategoryT string
type RadiusT int32
type OffsetT int32
type MaxSearchResultsT int32
type LatitudeT float64
type LongitudeT float64

type StreetT string
type ZipT string
type CityT string
type StateT string
type CountryCodeT string
type AttributeT interface{}

type Success struct {
	TimeStamp string `json:"timeStamp"`
}
type ApplicationResults struct {
	Success []Success
}

// sabreURL sets the base rest URL with a version.
// returns "https://api.havail.sabre.com/v1.0.0"
func sabreURL() *url.URL {
	base, _ := url.Parse(baseDevHavailURL)
	base.Path = path.Join(base.Path, apiVersion)
	return base
}
