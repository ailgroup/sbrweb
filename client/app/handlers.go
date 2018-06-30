package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ailgroup/sbrweb/apperr"
	"github.com/go-playground/form"
)

type AvailParamsBase struct {
	GuestCount int    `json:"guest_count" form:"guest_count"`
	Arrive     string `json:"arrive" form:"arrive"`
	Depart     string `json:"depart" form:"depart"`
}
type AvailParamsIDs struct {
	AvailParamsBase
	HotelIDs []string `json:"hotel_ids" form:"hotel_ids"`
}
type AvailParamsCityCodes struct {
	Base      AvailParamsBase
	CityCodes []string
}
type AvailParamsLatLng struct {
	Base   AvailParamsBase
	LatLng []string
}

// Validate AvailParamsBase fields. Time date formats arrive/depart are using app timezone location aware validations. Integer guest_count checks against min/max.
func (b AvailParamsBase) Validate(loc *time.Location) error {
	//check for null or empty values first
	if b.Arrive == "" {
		return ErrArriveNull
	}
	if b.Depart == "" {
		return ErrDepartNull
	}
	if b.GuestCount == 0 {
		return ErrGuestCountNullOrZero
	}

	tArrive, ok, err := StayFormat(b.Arrive, loc)
	if !ok {
		return ErrStayFormat(ErrArriveFmtMsg, err.Error(), tArrive.String(), timeShortForm)
	}
	//get app time zone location
	today := BeginOfDay(time.Now().In(loc))
	if ArriveNotInPast(tArrive, today) {
		return ErrArriveInPast(ErrStayInPastMsg, tArrive.String(), today)
	}
	tDepart, ok, err := StayFormat(b.Depart, loc)
	if !ok {
		return ErrStayFormat(ErrDepartFmtMsg, err.Error(), tDepart.String(), timeShortForm)
	}
	if DepartBeforeArrive(tDepart, tArrive) {
		return ErrStayRange(ErrStayRangeMsg, tDepart.String(), tArrive.String())
	}
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
	if err := a.AvailParamsBase.Validate(loc); err != nil {
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
https://github.com/tidwall/gjson

curl -H "Accept: application/json"  -X GET -d '{"arrive":"06-28","depart":"06-29","guest_count":"2", "hotel_ids":["007"]}' http://localhost:8080/avail | jq
*/
func (s *Server) HotelAvailIDsHandler() http.HandlerFunc {
	var params AvailParamsIDs
	var decoder *form.Decoder
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		params = AvailParamsIDs{}
		decoder = form.NewDecoder()
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
		/*
			//get session
			sess := s.SessionPool.Pick()
			//defer close on session
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
				params.AvailParamsBase.Arrive,
				params.AvailParamsBase.Depart,
			)
			req := hotelws.BuildHotelAvailRequest(s.SConfig, availBody)
			resp, _ := hotelws.CallHotelAvail(s.SConfig.ServiceURL, req)
		*/

		/*
			b, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		*/

		w.Header().Set("Content-Type", "application/json")
		//json.NewEncoder(w).Encode(resp)
		json.NewEncoder(w).Encode(params)
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
