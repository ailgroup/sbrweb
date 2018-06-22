package main

import (
	"fmt"
	"net/http"

	"github.com/ailgroup/sbrweb/client/common"
	"github.com/ailgroup/sbrweb/client/handlers"
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

func setConfig() *viper.Viper {
	conf := viper.GetViper()

	conf.SetConfigName("config")
	conf.AddConfigPath("$HOME")
	conf.AddConfigPath(".")
	conf.BindEnv("SABRE_USERNAME")
	conf.BindEnv("SABRE_PASSWORD")
	conf.BindEnv("SABRE_PCC")

	err := conf.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return conf
}

// registerRoutes is responsible for registering the server-side request handlers
func registerRoutes(env *common.Env) {
	env.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("hi")) })
	env.Mux.Handle("/avail", handlers.HotelAvailHandler(env))
}

func main() {
	// deal with context??
	r := chi.NewRouter()
	//Logger Logs the start and end of each request with the elapsed processing time
	r.Use(middleware.Logger)
	// Heartbeat Monitoring endpoint to check the servers pulse
	r.Use(middleware.Heartbeat("/heartbeat"))

	env := common.Env{
		Config: setConfig(),
		Mux:    r,
	}
	registerRoutes(&env)
	http.ListenAndServe(":8080", r)
}
