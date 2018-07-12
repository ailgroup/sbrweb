package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ailgroup/sbrweb/apperr"
	"github.com/ailgroup/sbrweb/engine/hotelws"
	"github.com/go-playground/form"
)

type AvailParamsBase struct {
	GuestCount  int    `form:"guest_count"`
	InputArrive string `form:"arrive"`
	InputDepart string `form:"depart"`
	OutArrive   string
	OutDepart   string
}
type AvailParamsIDs struct {
	*AvailParamsBase
	HotelIDs []string `form:"hotel_ids"`
}

/* TODO...
type AvailParamsCityCodes struct {
	Base      AvailParamsBase
	CityCodes []string
}
type AvailParamsLatLng struct {
	Base   AvailParamsBase
	LatLng []string
}
*/

type AvailabilityResponse struct {
	RequestParams     AvailParamsIDs
	SabreEngineErrors interface{} `json:",omitempty"`
	HotelAvail        hotelws.AvailabilityOptions
}

// Validate AvailParamsBase fields. Time date formats arrive/depart are using app timezone location aware validations and set the outgoing arrive/depart formats for sabre. Integer guest_count checks against min/max.
func (b *AvailParamsBase) ValidateAndFormat(loc *time.Location) error {
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
func (a AvailParamsIDs) Validate(loc *time.Location) error {
	if err := a.AvailParamsBase.ValidateAndFormat(loc); err != nil {
		return err
	}
	//defense against no param or weird values like 'hotel_ids=', 'hotel_ids="', 'hotel_ids=""'
	if len(a.HotelIDs) == 0 {
		return ErrHotelIDNullOrZero
	}
	if len(a.HotelIDs) == 1 {
		if (a.HotelIDs[0] == "") || (a.HotelIDs[0] == "\"") || (a.HotelIDs[0] == "\"\"") {
			return ErrHotelIDNullOrZero
		}
	}
	if Gt(len(a.HotelIDs), HotelIDsMax) {
		return ErrLtGt(ErrHotelIDsMaxsg, len(a.HotelIDs), HotelIDsMax)
	}
	return nil
}

/*
HotelAvailHandler wraps SOAP call to sabre hotel availability service

	Example:
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/avail?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_ids=10'
*/
func (s *Server) HotelAvailIDsHandler() http.HandlerFunc {
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &AvailParamsIDs{}
		decoder := form.NewDecoder()
		response := AvailabilityResponse{}
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

		q, _ := hotelws.NewHotelSearchCriteria(
			hotelws.HotelRefSearch(searchids),
		)
		availBody := hotelws.SetHotelAvailBody(
			params.AvailParamsBase.GuestCount,
			q,
			params.AvailParamsBase.OutArrive,
			params.AvailParamsBase.OutDepart,
		)

		req := hotelws.BuildHotelAvailRequest(s.SConfig, availBody)
		call, err := hotelws.CallHotelAvail(s.SConfig.ServiceURL, req)
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write(apperr.DecodeUnknown("CallHotelAvail::HotelAvailIDsHandler", r.URL.Query(), err, http.StatusFailedDependency))
			return
		}

		response.HotelAvail = call.Body.HotelAvail.AvailOpts

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

/*

	arrive := params.Get("arrive")
	depart := params.Get("depart")
	guests, _ := strconv.Atoi(params.Get("guest_count"))
	ids := params["hotel_ids"]
		search[hotelws.HotelidQueryField] = ids
		q, _ := hotelws.NewHotelSearchCriteria(
			hotelws.HotelRefSearch(search),
		)
		availBody := hotelws.SetHotelAvailBody(guests, q, arrive, depart)
		req := hotelws.BuildHotelAvailRequest(s.SConfig, availBody)
		resp, err := hotelws.CallHotelAvail(s.SConfig.ServiceURL, req)
		if err != nil {
			fmt.Printf("CallHotelAvail ERROR: %v", err)
			fmt.Fprint(w, err)
			return
		}
*/

/*
func (b AvailParamsBase) Validate() error {
	if govalidator.IsNull(b.GuestCount) {
		return ErrGuestCountNull
	}
	if govalidator.IsInt(b.GuestCount) {
		return ErrGuestCountNotNumber
	}
	if govalidator.InRangeInt(b.GuestCount, "1", "4") {
		return ErrGuestCountRange
	}
	return nil
}

func (a AvailParamsIDs) Validate() error {
	//if len(a.HotelIDs)
	if len(a.HotelIDs) < 1 || len(a.HotelIDs) > 10 {
		return ErrHotelIDRange
	}
	return nil
}
*/

/*
HOTELEREF:
	hqids               = make(HotelRefCriterion)
	hqltln              = make(HotelRefCriterion)
	hqcity              = make(HotelRefCriterion)

	search[hotelws.HotelidQueryField] 	= []string{"0012", "19876", "1109", "445098", "000034"}
	search[hotelws.LatlngQueryField] 	= []string{"32.78,-96.81", "54.87,-102.96"}
	search[hotelws.CountryCodeQueryField] = []string{"DFW", "CHC", "LA"}

ADDRESS:
	addr        				= make(AddressCriterion)
	addr[StreetQueryField] 		= "2031 N. 100 W"
	addr[CityQueryField] 		= "Aneheim"
	addr[PostalQueryField] 		= "90458"
	addr[CountryCodeQueryField] = "US"
*/
