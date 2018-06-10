package main

import (
	"github.com/ailgroup/sbrweb/client/common"
	"github.com/ailgroup/sbrweb/client/handlers"
	"github.com/gorilla/mux"
)

// registerRoutes is responsible for regisetering the server-side request handlers
func registerRoutes(env *common.Env, r *mux.Router) {
	r.Handle("/avail", handlers.HotelAvailHandler(env)).Methods("GET")
}

func main() {
	env := common.Env{}
	r := mux.NewRouter()
	registerRoutes(&env, r)
}
