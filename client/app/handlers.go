package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ailgroup/sbrweb/engine/hotelws"
)

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
var (
	search = hotelws.HotelRefCriterion{}
	addr   = hotelws.AddressCriterion{}
)

// HotelAvailHandler wraps SOAP call to sabre hotel availability service
// http://localhost:8080/avail?arrive=06-25&depart=06-26&guest_count=2&hotel_id=10&hotel_id=12&hotel_id=13&city_code=DWF&city_code=CHC&city_code=LA&lat_lng=32.78,-96.81&lat_lng=54.87,-102.96

// curl --request GET --url 'http://localhost:8080/avail?arrive=06-25&depart=06-26&guest_count=2&hotel_id=10&hotel_id=12&hotel_id=13&city_code=DWF&city_code=CHC&city_code=LA&lat_lng=32.78,-96.81&lat_lng=54.87,-102.96'
func (s *Server) HotelAvailIDsHandler() http.HandlerFunc {
	//define hotel search criterion
	search = make(hotelws.HotelRefCriterion)
	addr = make(hotelws.AddressCriterion)
	//closure to execute
	return func(w http.ResponseWriter, r *http.Request) {
		//get session
		sess := s.SessionPool.Pick()
		//defer close on session
		defer s.SessionPool.Put(sess)
		//get binary security token
		s.SConfig.SetBinSec(sess.Sabre)
		// parse incoming params
		params := r.URL.Query()
		if ids, ok := params["hotel_id"]; ok {
			search[hotelws.HotelidQueryField] = ids
		}
		if latlng, ok := params["lat_lng"]; ok {
			search[hotelws.LatlngQueryField] = latlng
		}
		if cities, ok := params["city_code"]; ok {
			search[hotelws.CountryCodeQueryField] = cities
		}

		if len(search) > 1 || len(search) < 1 {
			log.Printf("Search Validation ERROR: %v", search)
			fmt.Fprint(w, ErrInvalidSearchCriterion)
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
		fmt.Fprint(w, search)
	}
}
