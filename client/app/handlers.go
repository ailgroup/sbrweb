package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ailgroup/sbrweb/apperr"
	"github.com/ailgroup/sbrweb/engine/hotelws"
	"github.com/go-playground/form"
)

/*

	TODO... implement hotel availability requests by city_code, lat/long, and address


HOTELEREF:
	hqids               = make(HotelRefCriterion)
	hqltln              = make(HotelRefCriterion)
	hqcity              = make(HotelRefCriterion)

	search[hotelws.HotelidQueryField] 	= []string{"0012", "19876", "1109", "445098", "000034"}
	search[hotelws.LatlngQueryField] 	= []string{"32.78,-96.81", "54.87,-102.96"}
	search[hotelws.CountryCodeQueryField] = []string{"DFW", "CHC", "LA"}

	//TODO...
type HotelParamsLatLng struct {
	Base   AvailParamsBase
	LatLng []string
}
	//TODO...
type HotelParamsCityCodes struct {
	Base      AvailParamsBase
	CityCodes []string
}


ADDRESS:
	addr        				= make(AddressCriterion)
	addr[StreetQueryField] 		= "2031 N. 100 W"
	addr[CityQueryField] 		= "Aneheim"
	addr[PostalQueryField] 		= "90458"
	addr[CountryCodeQueryField] = "US"

	//TODO...
type HotelParamsAddress struct {
	Base      	AvailParamsBase
	Street 		string
	City 		string
	PostalCode 	string
	CountryCode string
}
*/

// HotelParamsBase hold guest, arrive, depart that are needed for any hotel request. Distinguish between incoming and outgoing params to allow mutliple time formats; that is, sabre onlyl accepts month-day formats but we want client API to force using the year as well.
type HotelParamsBase struct {
	GuestCount  int    `form:"guest_count"`
	InputArrive string `form:"arrive"`
	InputDepart string `form:"depart"`
	OutArrive   string
	OutDepart   string
}

// HotelParamsIDs holds 1..n hotel ids for making queries
type HotelParamsIDs struct {
	*HotelParamsBase
	Max      int
	HotelIDs []string `form:"hotel_ids"`
}

// AvailHotelIDSResponse for sabre hotel availability
type AvailHotelIDSResponse struct {
	RequestParams       HotelParamsIDs
	SabreEngineErrors   interface{} `json:",omitempty"`
	AvailabilityOptions hotelws.AvailabilityOptions
}

// PropertyHotelIDResponse for sabre property description
type PropertyHotelIDResponse struct {
	RequestParams     HotelParamsIDs
	SabreEngineErrors interface{} `json:",omitempty"`
	RoomStay          hotelws.RoomStay
}

// Validate AvailParamsBase fields. Time date formats arrive/depart are using app timezone location aware validations and set the outgoing arrive/depart formats for sabre. Integer guest_count checks against min/max.
func (b *HotelParamsBase) ValidateAndFormat(loc *time.Location) error {
	//check for null or empty values first
	if b.InputArrive == "" {
		return ErrArriveNull
	}
	if b.InputDepart == "" {
		return ErrDepartNull
	}
	if b.GuestCount == 0 {
		return ErrGuestCountNullOrZero
	}

	tArrive, ok, err := StayFormat(b.InputArrive, loc)
	if !ok {
		return ErrStayFormat(ErrArriveFmtMsg, err.Error(), tArrive.String(), timeShortForm)
	}
	//get app time zone location
	today := BeginOfDay(time.Now().In(loc))
	if ArriveNotInPast(tArrive, today) {
		return ErrArriveInPast(ErrStayInPastMsg, tArrive.String(), today)
	}
	tDepart, ok, err := StayFormat(b.InputDepart, loc)
	if !ok {
		return ErrStayFormat(ErrDepartFmtMsg, err.Error(), tDepart.String(), timeShortForm)
	}
	if DepartBeforeArrive(tDepart, tArrive) {
		return ErrStayRange(ErrStayRangeMsg, tDepart.String(), tArrive.String())
	}
	b.OutArrive = tArrive.Format(hotelws.TimeFormatMD)
	b.OutDepart = tDepart.Format(hotelws.TimeFormatMD)

	if Gt(b.GuestCount, GuestMax) {
		return ErrLtGt(ErrGuestMaxMsg, b.GuestCount, GuestMax)
	}
	if Lt(b.GuestCount, GuestMin) {
		return ErrLtGt(ErrGuestMinMsg, b.GuestCount, GuestMin)
	}
	return nil
}

// Validate AvailParamsIDs runs params base validations and for hotel_ids
func (a HotelParamsIDs) Validate(loc *time.Location) error {
	if err := a.HotelParamsBase.ValidateAndFormat(loc); err != nil {
		return err
	}
	if len(a.HotelIDs) == 0 {
		return ErrHotelIDNullOrZero
	}
	if Gt(len(a.HotelIDs), a.Max) {
		return ErrLtGt(ErrHotelIDsMaxsg, len(a.HotelIDs), HotelIDsMax)
	}
	//defense against no param or weird values like 'hotel_ids=', 'hotel_ids="', 'hotel_ids=""'
	if len(a.HotelIDs) == 1 {
		if (a.HotelIDs[0] == "") || (a.HotelIDs[0] == "\"") || (a.HotelIDs[0] == "\"\"") {
			return ErrHotelIDNullOrZero
		}
	}
	return nil
}

/*
PropertyDescriptionIDsHandler wraps SOAP call to sabre property description service. This SOAP service is the primary service for returning room rates. It accepts one hotel ref criterion and returns one hotel with one room stay object containing 0..n room rates.

	Example:
*/
func (s *Server) PropertyDescriptionIDsHandler() http.HandlerFunc {
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &HotelParamsIDs{Max: 1}
		decoder := form.NewDecoder()
		response := PropertyHotelIDResponse{}
		// decode params, check errors
		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apperr.DecodeBadInput("PropertyDescriptionIDsHandler", r.URL.Query(), err, http.StatusBadRequest))
			return
		}
		// validate query params
		if err := params.Validate(s.SConfig.AppTimeZone); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(apperr.DecodeInvalid("PropertyDescriptionIDsHandler", err, http.StatusUnprocessableEntity))
			return
		}
		response.RequestParams = *params
		//get session, defer close
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)
		//get binary security token
		s.SConfig.SetBinSec(sess.Sabre)
		// parse incoming params as JSON
		hotelid := make(hotelws.HotelRefCriterion)
		hotelid[hotelws.HotelidQueryField] = params.HotelIDs

		//no need to handle error in these two functions since api validations give similar guarantee
		q, _ := hotelws.NewHotelSearchCriteria(
			hotelws.HotelRefSearch(hotelid),
		)
		body, _ := hotelws.SetHotelPropDescBody(
			params.HotelParamsBase.GuestCount,
			q,
			params.HotelParamsBase.OutArrive,
			params.HotelParamsBase.OutDepart,
		)

		req := hotelws.BuildHotelPropDescRequest(s.SConfig, body)
		call, err := hotelws.CallHotelPropDesc(s.SConfig.ServiceURL, req)
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write(apperr.DecodeUnknown("CallHotelPropDesc::PropertyDescriptionIDsHandler", r.URL.Query(), err, http.StatusFailedDependency))
			return
		}

		response.RoomStay = call.Body.HotelDesc.RoomStay

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

/*
HotelAvailIDsHandler wraps SOAP call to sabre hotel availability service. This SOAP service does not return rates for rooms. It accepts many hotel ref criteria and returns many hotel options.

	Example:
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/avail/hotel/id?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_ids=10'
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/avail/hotel/id?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_ids=10,12'
*/
func (s *Server) HotelAvailIDsHandler() http.HandlerFunc {
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &HotelParamsIDs{Max: HotelIDsMax}
		decoder := form.NewDecoder()
		response := AvailHotelIDSResponse{}
		// decode params, check errors
		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apperr.DecodeBadInput("HotelAvailIDsHandler", r.URL.Query(), err, http.StatusBadRequest))
			return
		}
		// validate query params
		if err := params.Validate(s.SConfig.AppTimeZone); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(apperr.DecodeInvalid("HotelAvailIDsHandler", err, http.StatusUnprocessableEntity))
			return
		}
		response.RequestParams = *params
		//get session, defer close
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)
		//get binary security token
		s.SConfig.SetBinSec(sess.Sabre)
		// parse incoming params as JSON
		searchids := make(hotelws.HotelRefCriterion)
		searchids[hotelws.HotelidQueryField] = params.HotelIDs

		//no need to handle error in these two functions since api validations give similar guarantee
		q, _ := hotelws.NewHotelSearchCriteria(
			hotelws.HotelRefSearch(searchids),
		)
		availBody := hotelws.SetHotelAvailBody(
			params.HotelParamsBase.GuestCount,
			q,
			params.HotelParamsBase.OutArrive,
			params.HotelParamsBase.OutDepart,
		)

		req := hotelws.BuildHotelAvailRequest(s.SConfig, availBody)
		call, err := hotelws.CallHotelAvail(s.SConfig.ServiceURL, req)
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write(apperr.DecodeUnknown("CallHotelAvail::HotelAvailIDsHandler", r.URL.Query(), err, http.StatusFailedDependency))
			return
		}

		response.AvailabilityOptions = call.Body.HotelAvail.AvailOpts

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
