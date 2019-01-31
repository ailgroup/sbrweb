package itin

import (
	"encoding/xml"
	"fmt"

	"github.com/ailgroup/sbrweb/engine/sbrerr"
)

type DateTime struct {
	//XMLName xml.Name `xml:"DateTime"`
	Val string `xml:",chardata"`
}
type CityCode struct {
	//XMLName xml.Name `xml:"CityCode"`
	Name string `xml:"name,attr"`
	Val  string `xml:",chardata"`
}
type Departure struct {
	//XMLName  xml.Name `xml:"Departure"`
	DateTime DateTime
	CityCode CityCode
}
type Arrival struct {
	//XMLName  xml.Name `xml:"Arrival"`
	DateTime DateTime
	CityCode CityCode
}

type ReservationItem struct {
}
type ItineraryInfo struct {
	XMLName          xml.Name          `xml:"ItineraryInfo"`
	ReservationItems []ReservationItem `xml:"ReservationItems"`
}
type ItineraryRef struct {
	XMLName     xml.Name `xml:"ItineraryRef"`
	ID          string   `xml:"ID,attr"`
	AirExtras   bool     `xml:"AirExtras,attr"`
	InhibitCode string   `xml:"InhibitCode,attr"`
	PartitionID string   `xml:"PartitionID,attr"`
	PrimeHostID string   `xml:"PrimeHostID,attr"`
	Source      struct {
		PseudoCityCode string `xml:"PseudoCityCode,attr"`
		CreateDateTime string `xml:"CreateDateTime,attr"`
	} `xml:"Source"`
}
type TravelItinerary struct {
	XMLName       xml.Name `xml:"TravelItinerary"`
	Customer      CustomerInfo
	ItineraryInfo ItineraryInfo
	ItineraryRef  ItineraryRef
}

type SystemMessage struct {
	Code string `xml:"code"`
	Val  string `xml:",chardata"`
}
type SystemResult struct {
	Messages []SystemMessage `xml:"Message"`
}
type AppResWarning struct {
	Type          string         `xml:"type,attr"`
	Timestamp     string         `xml:"timeStamp,attr"`
	SystemResults []SystemResult `xml:"SystemSpecificResults"`
}
type AppResError struct {
	Type          string         `xml:"type,attr"`
	Timestamp     string         `xml:"timeStamp,attr"`
	SystemResults []SystemResult `xml:"SystemSpecificResults"`
}

type ApplicationResults struct {
	XMLName xml.Name `xml:"ApplicationResults"`
	Status  string   `xml:"status,attr"`
	Success struct {
		Timestamp string `xml:"timeStamp,attr"`
	} `xml:"Success"`
	Warnings []AppResWarning `xml:"Warning"`
	Errors   []AppResError   `xml:"Error"`
}

func (result ApplicationResults) Ok() bool {
	switch result.Status {
	case sbrerr.StatusNotProcess(): //queries
		return false
	case sbrerr.StatusComplete(): //queries, pnr
		if len(result.Warnings) > 0 {
			return false
		}
		if len(result.Errors) > 0 {
			return false
		}
		return true
	default:
		return false
	}
}
func (result ApplicationResults) ErrFormat() sbrerr.ErrorSabreResult {
	var wmsg string
	for i, w := range result.Warnings {
		var msg string
		for is, s := range w.SystemResults {
			for ms, m := range s.Messages {
				msg += fmt.Sprintf("SystemResult%d:Message-%d:Code-%s:Val-%s. ", is, ms, m.Code, m.Val)
			}
		}
		wmsg += fmt.Sprintf("Warning%d:Type-%s:Msg|%s", i, w.Type, msg)
	}
	for i, w := range result.Errors {
		var msg string
		for is, s := range w.SystemResults {
			for ms, m := range s.Messages {
				msg += fmt.Sprintf("SystemResult%d:Message-%d:Code-%s:Val-%s. ", is, ms, m.Code, m.Val)
			}
		}
		wmsg += fmt.Sprintf("Error%d:Type-%s:Msg|%s", i, w.Type, msg)
	}
	return sbrerr.ErrorSabreResult{
		Code:       sbrerr.SabreEngineStatusCode(result.Status),
		AppMessage: wmsg,
	}
}
