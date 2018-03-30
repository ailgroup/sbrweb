package sbrweb

import (
	"encoding/xml"
	"errors"
)

const (
	propertyQueryField = "property"
	locationQueryField = "location"
	amenityQueryField  = "amenity"
	airportQueryField  = "airport"
	creditQueryField   = "creditcard"
	diningTypeField    = "Dining"
	alertTypeField     = "Alerts"
)

// GetHotelContent retrieve sabre hotel content using various query criteria, see SearchCriteria
type GetHotelContent struct {
	XMLName             xml.Name `xml:"GetHotelContentRQ"`
	XMLNSXsi            string   `xml:"xmlns:xsi,attr"`
	XMLNSSchemaLocation string   `xml:"xmlns:schemaLocation,attr"`
	Version             string   `xml:"version,attr"`
	SearchCriteria      SearchCriteria
}

// GetHotelContentUnmarsh wrapper to unmarshal with namespace prefix
type GetHotelContentUnmarsh struct {
	XMLName             xml.Name `xml:"GetHotelContentRQ"`
	XMLNSXsi            string   `xml:"xmlns,attr,omitempty"`
	XMLNSSchemaLocation string   `xml:"xmlnsn,attr,omitempty"`
	Version             string   `xml:"version,attr,omitempty"`
}

// SearchCriteria dynamic struct to build query for hotel content
type SearchCriteria struct {
	ImageRef        ImageRef
	HotelRefs       []HotelRef      `xml:"HotelRefs>HotelRef"`
	DescriptiveInfo DescriptiveInfo `xml:"DescriptiveInfoRef"`
}

type DescriptiveInfo struct {
	Property            bool          `xml:"PropertyInfo,omitempty"`
	Location            bool          `xml:"LocationInfo,omitempty"`
	Amenities           bool          `xml:"Amenities,omitempty"`
	Airports            bool          `xml:"Airports,omitempty"`
	AcceptedCreditCards bool          `xml:"AcceptedCreditCards,omitempty"`
	Descriptions        []Description `xml:"Descriptions>Description"`
	//..more?..//
}

// Description value of which descriptions to return
type Description struct {
	Type string `xml:"Type,attr"`
}

// HotelRef holds hotel ids. Must be a string. Sabre has ids that often need to be passed with prefixed zeros (e.g.... 0012)
type HotelRef struct {
	HotelCode string `xml:"HotelCode,attr"`
}

// ImageReg how many images you want
type ImageRef struct {
	MaxImage int `xml:"MaxImages,attr"`
}

type DescriptiveQuery map[string]bool
type DescriptionTypes []string

func buildDescriptions(query DescriptiveQuery, dtypes DescriptionTypes) DescriptiveInfo {
	d := DescriptiveInfo{}
	for k, v := range query {
		if v {
			switch k {
			case propertyQueryField:
				d.Property = true
			case locationQueryField:
				d.Location = true
			case amenityQueryField:
				d.Amenities = true
			case airportQueryField:
				d.Airports = true
			case creditQueryField:
				d.AcceptedCreditCards = true
			}
		}
	}
	for _, t := range dtypes {
		switch t {
		case diningTypeField:
			d.Descriptions = append(d.Descriptions, Description{Type: diningTypeField})
		case alertTypeField:
			d.Descriptions = append(d.Descriptions, Description{Type: alertTypeField})
		}
	}
	return d
}

//BuildGetContentRequest payload for getting hotel content
func BuildGetHotelContent(images int, hotelIDs []string, dinfo DescriptiveInfo) (GetHotelContent, error) {
	if len(hotelIDs) == 0 {
		return GetHotelContent{}, errors.New("Must have more than 1 hotel id")
	}
	var hotelrefs []HotelRef
	for _, id := range hotelIDs {
		hotelrefs = append(hotelrefs, HotelRef{HotelCode: id})
	}
	return GetHotelContent{
		XMLNSXsi:            baseXsiNamespace,
		XMLNSSchemaLocation: baseGetHotelContentSchema,
		Version:             sabreHotelContentVersion,
		SearchCriteria: SearchCriteria{
			HotelRefs:       hotelrefs,
			DescriptiveInfo: dinfo,
			ImageRef:        ImageRef{MaxImage: images},
		},
	}, nil
}
