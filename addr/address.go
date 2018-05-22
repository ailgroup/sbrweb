package Addr

// Address represents typical building addresses
type Address struct {
	AddressLine   string `xml:"AddressLine,omitempty"`
	Street        string `xml:"StreetNumber,omitempty"`
	City          string `xml:"CityName,omitempty"`
	StateProvince struct {
		StateCode string `xml:"StateCode,attr"`
	} `xml:"StateCountyProv,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
}
