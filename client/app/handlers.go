package app

import (
	"fmt"
	"net/http"
)

// HotelAvailHandler wraps SOAP call to sabre hotel availability service
func (s *Server) HotelAvailHandler() http.HandlerFunc {
	msg := "this"
	return func(w http.ResponseWriter, r *http.Request) {
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)

		//1 way:
		//param1 := r.URL.Query().Get("param1")

		//2 way:
		//value := FormValue("field")
		fmt.Fprint(w, msg)
	}
}
