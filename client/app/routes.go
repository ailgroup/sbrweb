package app

import (
	"net/http"
)

// registerRoutes is responsible for registering the server-side request handlers
func (s *Server) RegisterRoutes() {
	s.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hi")) })
	s.Mux.Handle("/avail", s.HotelAvailHandler())
}
