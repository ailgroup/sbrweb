package sbrweb

import (
	"encoding/xml"
	"errors"
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
	ImageRef  ImageRef
	HotelRefs []HotelRef `xml:"HotelRefs>HotelRef"`
}

// HotelRef holds hotel ids. Must be a string. Sabre has ids that often need to be passed with prefixed zeros (e.g.... 0012)
type HotelRef struct {
	HotelCode string `xml:"HotelCode,attr"`
}

// ImageReg how many images you want
type ImageRef struct {
	MaxImage int `xml:"MaxImages,attr"`
}

//BuildGetContentRequest payload for getting hotel content
func BuildGetContentRequest(images int, hotelIDs []string) (GetHotelContent, error) {
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
			ImageRef:  ImageRef{MaxImage: images},
			HotelRefs: hotelrefs,
		},
	}, nil
}
