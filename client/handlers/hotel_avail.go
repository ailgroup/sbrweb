package handlers

import (
	"fmt"
	"net/http"

	"github.com/ailgroup/sbrweb/client/common"
)

//https://stackoverflow.com/questions/15407719/in-gos-http-package-how-do-i-get-the-query-string-on-a-post-request
func HotelAvailHandler(env *common.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//1 way:
		//param1 := r.URL.Query().Get("param1")

		//2 way:
		//value := FormValue("field")
		fmt.Fprint(w, "hello hotel availability")
	})
}
