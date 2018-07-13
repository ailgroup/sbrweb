package app

import (
	"net/http"
)

// registerRoutes is responsible for registering the server-side request handlers
func (s *Server) RegisterRoutes() {
	s.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"ping":"p0ng"}`)) })
	s.Mux.Handle("/avail/hotel/id", s.HotelAvailIDsHandler())
	s.Mux.Handle("/rates/hotel/id", s.PropertyDescriptionIDsHandler())
	//s.Mux.Handle("/avail/hotel/latlong", s.HotelAvailIDsHandler())
	//s.Mux.Handle("/avail/hotel/address", s.HotelAvailIDsHandler())
	//s.Mux.Handle("/avail/hotel/citycode", s.HotelAvailIDsHandler())
}
