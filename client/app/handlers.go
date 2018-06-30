package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ailgroup/sbrweb/apperr"
	"github.com/go-playground/form"
)

/*
type AvailParamsBase struct {
	GuestCount int    `json:"guest_count" form:"guest_count" validate:"gt=0,lt=5"`
	Arrive     string `json:"arrive" form:"arrive" validate:"required"`
	Depart     string `json:"depart" form:"depart" validate:"required"`
}
type AvailParamsIDs struct {
	AvailParamsBase
	HotelIDs []string `json:"hotel_ids" form:"hotel_ids" validate:"required,dive,min=1,max=10"`
}
*/
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

// Validate AvailParamsBase fields arrive, depart, guest_count looking at date formats and counts.
func (b AvailParamsBase) Validate() error {
	tArrive, ok, err := ValidArriveDepart(b.Arrive)
	if !ok {
		return ErrStayFormat(ErrArriveFmtMsg, err.Error(), b.Arrive, timeShortForm)
	}

	tDepart, ok, err := ValidArriveDepart(b.Depart)
	if !ok {
		return ErrStayFormat(ErrDepartFmtMsg, err.Error(), b.Depart, timeShortForm)
	}

	if !startBeforeEnd(tArrive, tDepart) {
		return ErrStayRange(ErrStayRangeMsg, b.Arrive, b.Depart)
	}

	if b.GuestCount == 0 {
		return ErrGuestCountNullOrZero
	}
	if Gt(b.GuestCount, GuestGTE) {
		return ErrLtGt(ErrGuestLTEMsg, b.GuestCount, GuestGTE)
	}
	if Lt(b.GuestCount, GuestLTE) {
		return ErrLtGt(ErrGuestGTEMsg, b.GuestCount, GuestLTE)
	}
	return nil
}

func (a AvailParamsIDs) Validate() error {
	if err := a.AvailParamsBase.Validate(); err != nil {
		return err
	}
	fmt.Printf("%v %d\n", a.HotelIDs, len(a.HotelIDs))
	if len(a.HotelIDs) == 1 {
		if (a.HotelIDs[0] == "") || (a.HotelIDs[0] == "\"") || (a.HotelIDs[0] == "\"\"") {
			return ErrHotelIDNullOrZero
		}
	}

	if Gt(len(a.HotelIDs), HotelIDsLTE) {
		return ErrLtGt(ErrHotelIDsLTEMsg, len(a.HotelIDs), HotelIDsLTE)
	}
	return nil
}

/*
HotelAvailHandler wraps SOAP call to sabre hotel availability service
https://github.com/tidwall/gjson

curl -H "Accept: application/json"  -X GET -d '{"arrive":"06-28","depart":"06-29","guest_count":"2", "hotel_ids":["007"]}' http://localhost:8080/avail | jq
*/
func (s *Server) HotelAvailIDsHandler() http.HandlerFunc {
	var availPIDs AvailParamsIDs
	var decoder *form.Decoder
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		availPIDs = AvailParamsIDs{}
		decoder = form.NewDecoder()
		if err := decoder.Decode(&availPIDs, r.URL.Query()); err != nil {
			appMsg := fmt.Sprintf("Cannot Decode Query: %v", r.URL.Query())
			b, _ := json.Marshal(
				apperr.NewErrorBadInput(
					err.Error(),
					appMsg,
					apperr.BadInput,
					http.StatusBadRequest,
				),
			)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(b)
			return
		}
		fmt.Printf("\n%+v\n", availPIDs)
		if err := availPIDs.Validate(); err != nil {
			b, _ := json.Marshal(
				apperr.NewErrorInvalid(
					err.Error(),
					"Invalid",
					apperr.Invalid,
					http.StatusUnprocessableEntity,
				),
			)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(b)
			return
		}
		/*
			// decode, validate params
			if err := decoder.Decode(&params, r.URL.Query()); err != nil {
				appMsg := fmt.Sprintf("Cannot Decode Query: %v", r.URL.Query())
				b, _ := json.Marshal(
					apperr.NewErrorBadInput(err.Error(), appMsg, apperr.BadInput, http.StatusBadRequest),
				)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(b)
				return
			}
			if err := params.Validate(); err != nil {
				b, _ := json.Marshal(
					apperr.NewErrorInvalid(err.Error(), "Invalid Request", apperr.BadInput, http.StatusUnprocessableEntity),
				)
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write(b)
				return
			}
		*/
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
			g, _ := strconv.Atoi(params.Base.GuestCount)
			availBody := hotelws.SetHotelAvailBody(g, q, params.Base.Arrive, params.Base.Depart)
			req := hotelws.BuildHotelAvailRequest(s.SConfig, availBody)
			resp, _ := hotelws.CallHotelAvail(s.SConfig.ServiceURL, req)

			b, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		*/
		//b, _ := json.Marshal(params)
		//w.WriteHeader(http.StatusOK)
		//w.Write(b)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(availPIDs)
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
