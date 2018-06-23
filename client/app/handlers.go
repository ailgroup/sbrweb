package app

import (
	"fmt"
	"net/http"
)

// HotelAvailHandler wraps SOAP call to sabre hotel availability service
func (s *Server) HotelAvailHandler() http.HandlerFunc {
	pcc := s.Config.GetString("SABRE_PCC")
	clientURL := s.Config.GetString("sessions.client.url")
	sessExpireMin := s.Config.GetInt("sessions.expire.min")
	msg := fmt.Sprintf("PCC:%s\n From:%s\n SessionExpireMin:%d\n \n %s", pcc, clientURL, sessExpireMin, "hello hotel availability")
	return func(w http.ResponseWriter, r *http.Request) {
		sess := s.SessionPool.Pick()
		defer s.SessionPool.Put(sess)
		info := fmt.Sprintf("%s \n\n SessionID: %s\n ExpireTime: %v\n Started: %v\n ConvID: %s\n",
			msg,
			sess.ID,
			sess.ExpireTime,
			sess.TimeStarted,
			sess.Sabre.Body.SessionCreateRS.ConversationID,
		)
		//1 way:
		//param1 := r.URL.Query().Get("param1")

		//2 way:
		//value := FormValue("field")
		fmt.Fprint(w, info)
	}
}
