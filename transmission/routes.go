package transmission

import "net/http"

/*
// RegisterRoutes is responsible for registering the server-side request handlers
func (s *Server) RegisterRoutes() {
	s.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"ping":"p0ng"}`)) })
	s.Mux.Handle("/hotel/ids", s.HotelIDsHandler())
	//s.Mux.Handle("/hotel/latlong", s.Handler())
	//s.Mux.Handle("/hotel/address", s.Handler())
	//s.Mux.Handle("/hotel/citycodes", s.Handler())
	s.Mux.Handle("/rates/hotel/id", s.RatesHotelIDHandler())
	s.Mux.Handle("/book/hotel/room", s.BookRoomHandler())
}
*/

// RegisterRoutes is responsible for registering the server-side request handlers
func (s *Server) RegisterRoutes() {
	s.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"ping":"p0ng"}`)) })
	s.Mux.Get("/hotel/ids", s.HotelIDsHandler())
	s.Mux.Get("/rates/hotel/id", s.RatesHotelIDHandler())
	s.Mux.Post("/book/hotel/room", s.BookRoomHandler())
}
