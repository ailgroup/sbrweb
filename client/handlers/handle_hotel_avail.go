package handlers

import (
	"fmt"
	"net/http"

	"github.com/ailgroup/sbrweb/client/common"
)

// HotelAvailHandler wraps SOAP call to sabre hotel availability service
func HotelAvailHandler(env *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1 way:
		//param1 := r.URL.Query().Get("param1")

		//2 way:
		//value := FormValue("field")
		pcc := env.Config.GetString("SABRE_PCC")
		clientURL := env.Config.GetString("sessions.client.url")
		sessExpireMin := env.Config.GetInt("sessions.expire.min")
		msg := fmt.Sprintf("PCC:%s\n From:%s\n SessionExpireMin:%d\n \n %s", pcc, clientURL, sessExpireMin, "hello hotel availability")
		fmt.Fprint(w, msg)
	})
}
