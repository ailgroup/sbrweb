package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ailgroup/sbrweb/apperr"
	"github.com/ailgroup/sbrweb/engine/hotelws"
	"github.com/ailgroup/sbrweb/engine/itin"
	"github.com/go-playground/form"
)

// HotelParamsBase hold guest, arrive, depart that are needed for any hotel request. Distinguish between incoming and outgoing params to allow mutliple time formats; that is, sabre only accepts month-day formats but we want client API to force using the year as well.
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

// HotelParamsID holds 1 hotel id for making queries
type HotelParamsID struct {
	*HotelParamsBase
	HotelID string `form:"hotel_id"`
}

// BookRoomParams hold params for creating pnr and executing a reservation. CCCode is the credit card type code (MC, AMX, etc...); RoomRPH is the is the sabre reference place holder of the room. RPH is gotten from a previous request for room rates through RatesHotelIDHandler.
type BookRoomParams struct {
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	NumRooms  string `form:"num_rooms"`
	CCPhone   string `form:"cc_phone"`
	CCCode    string `form:"cc_code"`
	CCExpire  string `form:"cc_expire"`
	CCNumber  string `form:"cc_number"`
	RoomMeta  string `form:"room_meta"`
}

// AvailHotelIDSResponse for sabre hotel availability
type AvailHotelIDSResponse struct {
	RequestParams       HotelParamsIDs
	SabreEngineErrors   interface{} `json:",omitempty"`
	AvailabilityOptions hotelws.AvailabilityOptions
}

// PropertyHotelIDResponse for sabre property description
type PropertyHotelIDResponse struct {
	RequestParams     HotelParamsID
	SabreEngineErrors interface{} `json:",omitempty"`
	RoomStay          hotelws.RoomStay
}

// Validate AvailParamsBase fields. Time date formats arrive/depart are using app timezone location aware validations and set the outgoing arrive/depart formats for sabre. Integer guest_count checks against min/max.
func (b *HotelParamsBase) ValidateAndFormat(loc *time.Location) error {
	//check for null or empty values first
	if b == nil {
		return ErrBaseParamsNull
	}
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

// Validate HotelParamsID runs params base validations and for hotel_id
func (a HotelParamsID) Validate(loc *time.Location) error {
	if err := a.HotelParamsBase.ValidateAndFormat(loc); err != nil {
		return err
	}
	if a.HotelID == "" {
		return ErrHotelIDNullOrZero
	}
	//defense against no param or weird values like 'hotel_id=', 'hotel_id="', 'hotel_id=""'
	if (a.HotelID == "") || (a.HotelID == "\"") || (a.HotelID == "\"\"") {
		return ErrHotelIDNullOrZero
	}
	return nil
}

// Validate HotelParamsIDs runs params base validations and for hotel_ids
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
	//if len(a.HotelIDs) >= 1 {
	for _, id := range a.HotelIDs {
		if (id == "") || (id == "\"") || (id == "\"\"") {
			return ErrHotelIDNullOrZero

		}
	}
	//}
	return nil
}

// Validate AvailParamsIDs runs params base validations and for hotel_ids
func (b BookRoomParams) Validate() error {
	return nil
}

/*
	BookRoomHandler creates a pnr, fetches rate, books room, ends transaction. It accepts room meta data generated from previous requests. Required params: last_name, num_rooms, rph, cc_phone, cc_code, cc_expire, cc_number

	curl -H "Accept: application/json" -X GET 'http://localhost:8080/book/hotel/room?'
*/

func (s *Server) BookRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &BookRoomParams{}
		decoder := form.NewDecoder()

		// decode params, check errors
		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apperr.DecodeBadInput("BookRoomHandler", r.URL.Query(), err, http.StatusBadRequest))
			return
		}

		// validate query params
		if err := params.Validate(); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(apperr.DecodeInvalid("BookRoomHandler", err, http.StatusUnprocessableEntity))
			return
		}

		//PNR
		//build person
		person := itin.CreatePersonName(params.FirstName, params.LastName)
		//build pnr
		pnrBody := itin.SetPNRDetailBody(params.CCPhone, person)
		pnrReq := itin.BuildPNRDetailsRequest(s.SConfig, pnrBody)
		//call pnr
		pnrResp, err := itin.CallPNRDetail(s.SConfig.ServiceURL, pnrReq)
		if err != nil {
			w.WriteHeader(http.StatusFailedDependency)
			w.Write(apperr.DecodeUnknown("CallPNRDetail::BookRoomHandler", r.URL.Query(), err, http.StatusFailedDependency))
			return
		}

		json.NewEncoder(w).Encode(pnrResp)

		//Hotel Reservation
		/*
			reservationBody := hotelws.SetHotelResBody(1, srvc.SabreTimeNowFmt())
			reservationBody.NewGuaranteeRes(
				person.Last.Val,
				rr.GuaranteeSurcharge, //gtype string,
				"MC",               //ccCode string,
				"2019-07",          //ccExpire string,
				"5105105105105100", //ccNumber string,
			)
			reservationReq := hotelws.BuildHotelResRequest(s.SConfig, reservationBody)
			resResp, err := hotelws.CallHotelRes(prodURL, reservationReq)
			if err != nil {
				return err
			}

			call.SetTrackedEncode()
		*/
	}
}

/*
RatesHotelIDHandler wraps SOAP call to sabre property description service. This SOAP service is the primary service for returning room rates. It accepts one hotel ref criterion and returns one hotel with one room stay object containing 0..n room rates.

	Example:
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?guest_count=2&arrive=2018-07-17&depart=2018-07-18&hotel_id=10'
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/rates/hotel/id?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_id=12'
*/
func (s *Server) RatesHotelIDHandler() http.HandlerFunc {
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &HotelParamsID{}
		decoder := form.NewDecoder()
		response := PropertyHotelIDResponse{}
		// decode params, check errors
		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apperr.DecodeBadInput("RatesHotelIDHandler", r.URL.Query(), err, http.StatusBadRequest))
			return
		}
		// validate query params
		if err := params.Validate(s.SConfig.AppTimeZone); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(apperr.DecodeInvalid("RatesHotelIDHandler", err, http.StatusUnprocessableEntity))
			return
		}
		response.RequestParams = *params
		//get session, defer close
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)
		//get binary security token
		s.SConfig.SetBinSec(sess.Sabre)
		// setup hotel ref criterion for 1 hotel id
		hotelid := make(hotelws.HotelRefCriterion)
		hotelid[hotelws.HotelidQueryField] = []string{params.HotelID}

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
			w.Write(apperr.DecodeUnknown("CallHotelPropDesc::RatesHotelIDHandler", r.URL.Query(), err, http.StatusFailedDependency))
			return
		}

		call.SetRoomMetaData(params.GuestCount, params.OutArrive, params.OutDepart, params.HotelID)
		response.RoomStay = call.Body.HotelDesc.RoomStay

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

/*
HotelIDsHandler wraps SOAP call to sabre hotel availability service. This SOAP service does not return rates for rooms. Instead it shows basic availability options along with property information. It accepts many hotel ref criteria and returns many hotel options.

	Example:
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/hotel/ids?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_ids=10'
		curl -H "Accept: application/json" -X GET 'http://localhost:8080/hotel/ids?guest_count=4&arrive=2018-07-17&depart=2018-07-18&hotel_ids=10,12'
*/
func (s *Server) HotelIDsHandler() http.HandlerFunc {
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := &HotelParamsIDs{Max: HotelIDsMax}
		decoder := form.NewDecoder()
		response := AvailHotelIDSResponse{}
		// decode params, check errors
		if err := decoder.Decode(&params, r.URL.Query()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(apperr.DecodeBadInput("HotelIDsHandler", r.URL.Query(), err, http.StatusBadRequest))
			return
		}
		// validate query params
		if err := params.Validate(s.SConfig.AppTimeZone); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write(apperr.DecodeInvalid("HotelIDsHandler", err, http.StatusUnprocessableEntity))
			return
		}
		response.RequestParams = *params
		//get session, defer close
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)
		//get binary security token
		s.SConfig.SetBinSec(sess.Sabre)
		// setup hotel ref serch for hotel ids...
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
