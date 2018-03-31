package sbrhotel

const (
	hotelAvailVersion = "2.3.0"

	streetQueryField      = "street_qf"
	cityQueryField        = "city_qf"
	postalQueryField      = "postal_qf"
	countryCodeQueryField = "countryCode_qf"
	hotelidQueryField     = "hotelID_qf"
	returnHostCommand     = true
)

type Address struct {
	City        string `xml:"CityName,omitempty"`
	CountryCode string `xml:"CountryCode,omitempty"`
	Postal      string `xml:"PostalCode,omitempty"`
	Street      string `xml:"StreetNumber,omitempty"`
}

// QueryParams is a typed function to support optional query params on creation of new search criterion
type QueryParams func(*HotelSearchCriteria) error

//type HotelQueryParams map[string][]string
type GeoQueryParams map[string][]string
type AddressQueryParams map[string]string

func NewHotelSearchCriteria(queryParams ...QueryParams) (*HotelSearchCriteria, error) {
	criteria := &HotelSearchCriteria{}
	for _, qm := range queryParams {
		err := qm(criteria)
		if err != nil {
			return &HotelSearchCriteria{}, err
		}
	}
	return criteria, nil
}

// AddressOption Sets the starting learning rate; default is 0.025 for skip-gram,  and 0.05 for CBOW.
// Note, if you wanna change the learning rate for skip-gram you should do it AFTER cbow option has been set.
func AddressOption(params AddressQueryParams) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		a := Address{}
		for k, v := range params {
			switch k {
			case streetQueryField:
				a.Street = v
			case cityQueryField:
				a.City = v
			case postalQueryField:
				a.Postal = v
			case countryCodeQueryField:
				a.CountryCode = v
			}
		}
		q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{Address: a})
		return nil
	}
}

func GeoOption(params GeoQueryParams) func(q *HotelSearchCriteria) error {
	return func(q *HotelSearchCriteria) error {
		for k, v := range params {
			switch k {
			case cityQueryField:
				for _, city := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCityCode: city})
				}
			case hotelidQueryField:
				for _, code := range v {
					q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCode: code})
				}
			}
		}
		return nil
	}
}

/*

func buildHotelSearch(params HotelQueryParams) (HotelSearchCriteria, error) {
	q := HotelSearchCriteria{}
	if len(params) > 1 {
		return q, fmt.Errorf("Cannot use more than 1 search type criteria, have %v", params)
	}
	for k, v := range params {
		switch k {
		case cityQueryField:
			for _, city := range v {
				q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCityCode: city})
			}
		case hotelidQueryField:
			for _, code := range v {
				q.Criterion.HotelRef = append(q.Criterion.HotelRef, HotelRef{HotelCode: code})
			}
		}
	}

	return q, nil
}

func buildAddress(params AddressQueryParams) Address {
	a := Address{}
	for k, v := range params {
		switch k {
		case streetQueryField:
			a.Street = v
		case cityQueryField:
			a.City = v
		case postalQueryField:
			a.Postal = v
		case countryCodeQueryField:
			a.CountryCode = v
		}
	}
	return a
}
*/

/*
	propertyQueryField = "property"
	locationQueryField = "location"
	amenityQueryField  = "amenity"
	airportQueryField  = "airport"
	creditQueryField   = "creditcard"
	diningTypeField    = "Dining"
	alertTypeField     = "Alerts"

	sabreHotelContentVersion = "1.0.0"
*/
