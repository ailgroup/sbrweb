package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ailgroup/sbrweb/client/app"
	"github.com/ailgroup/sbrweb/engine/srvc"
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

const (
	ConfSabreUsername = "SABRE_USERNAME"
	ConfSabrePassword = "SABRE_PASSWORD"
	ConfSabrePCC      = "SABRE_PCC"
	ConfFile          = "config"
	ConfExpireMin     = "sessions.expire.min"
	ConfExpireMax     = "sessions.expire.max"
	ConfPoolSize      = "sessions.client.pool_size"
	ConfSabreURL      = "sessions.sabre_url"
	ConfClientURL     = "sessions.client.url"
)

var (
	vRepeatEvery = time.Minute * 3
	sessConf     = &srvc.SessionConf{}
	vipConf      = &viper.Viper{}
)

func init() {
	vipConf = setConfig()
	sessConf = &srvc.SessionConf{
		ServiceURL: vipConf.GetString(ConfSabreURL),
		From:       vipConf.GetString(ConfClientURL),
		PCC:        vipConf.GetString(ConfSabrePCC),
		Convid:     srvc.GenerateConversationID(vipConf.GetString(ConfClientURL)),
		Msgid:      srvc.GenerateMessageID(),
		Timestr:    srvc.SabreTimeFormat(),
		Username:   vipConf.GetString(ConfSabreUsername),
		Password:   vipConf.GetString(ConfSabrePassword),
	}
}

func setConfig() *viper.Viper {
	vip := viper.GetViper()

	vip.SetConfigName(ConfFile)
	vip.AddConfigPath("$HOME")
	vip.AddConfigPath(".")
	vip.BindEnv(ConfSabreUsername)
	vip.BindEnv(ConfSabrePassword)
	vip.BindEnv(ConfSabrePCC)

	err := vip.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	vip.WatchConfig()
	vip.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	return vip
}

func main() {
	// create scheme for session expirey
	scheme := srvc.ExpireScheme{
		Min: vipConf.GetInt(ConfExpireMin),
		Max: vipConf.GetInt(ConfExpireMax),
	}

	// create pool, background populate pool and sett up keepalive
	p := srvc.NewPool(scheme, sessConf, vipConf.GetInt(ConfPoolSize))
	go func() {
		err := p.Populate()
		if err != nil {
			panic(err)
		}
		srvc.Keepalive(p, vRepeatEvery)
	}()
	// close down session pool properly against sabre re-allocates useable sessiosn
	cls := make(chan os.Signal)
	signal.Notify(cls, os.Interrupt)
	go func() {
		sig := <-cls
		fmt.Printf("Got %s signal. Closing down...\n", sig)
		p.Close()
		os.Exit(1)
	}()

	// pass context through handlers??
	m := chi.NewRouter()
	//Logger Logs the start and end of each request with the elapsed processing time
	m.Use(middleware.Logger)
	// Heartbeat Monitoring endpoint to check the servers pulse
	m.Use(middleware.Heartbeat("/heartbeat"))

	server := app.Server{
		VConfig:     vipConf,
		SConfig:     sessConf,
		Mux:         m,
		SessionPool: p,
	}
	server.RegisterRoutes()
	fmt.Println("Begin on port:", ":8080")
	http.ListenAndServe(":8080", m)
}
