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
	ConfTimeZone      = "app_timezone"
	ConfExpireMin     = "sessions.expire.min"
	ConfExpireMax     = "sessions.expire.max"
	ConfPoolSize      = "sessions.client.pool_size"
	ConfSabreURL      = "sessions.sabre_url"
	ConfClientURL     = "sessions.client.url"
)

var (
	vRepeatEvery      = time.Minute * 3
	sessConf          = &srvc.SessionConf{}
	vipConf           = &viper.Viper{}
	port              = ":8080"
	ClientAppTimeZone = &time.Location{}
)

func init() {
	ClientAppTimeZone = time.UTC
	vipConf = setConfig()
	sessConf = &srvc.SessionConf{
		ServiceURL:  vipConf.GetString(ConfSabreURL),
		From:        vipConf.GetString(ConfClientURL),
		PCC:         vipConf.GetString(ConfSabrePCC),
		Convid:      srvc.GenerateConversationID(vipConf.GetString(ConfClientURL)),
		Msgid:       srvc.GenerateMessageID(),
		Timestr:     srvc.SabreTimeFormat(),
		Username:    vipConf.GetString(ConfSabreUsername),
		Password:    vipConf.GetString(ConfSabrePassword),
		AppTimeZone: ClientAppTimeZone,
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

	//loc, _ := time.LoadLocation("Europe/Berlin")
	vip.SetDefault("app_timezone", "UTC")

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

func runPool(p *srvc.SessionPool) {
	keepKill := make(chan os.Signal, 1)
	signal.Notify(keepKill, os.Interrupt)
	shutDown := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt)
	go func() {
		err := p.Populate()
		if err != nil {
			fmt.Printf("Error popluating session pool %v\n", err)
			os.Exit(1)
		}
		srvc.Keepalive(p, vRepeatEvery, keepKill)
		sig := <-shutDown
		fmt.Printf("\nGot '%s' SIGNAL. Shutdown keepalive, session pool, and exit program...\n", sig)
		p.Close()
		os.Exit(1)
	}()
}

func main() {
	// create scheme for session expirey
	scheme := srvc.ExpireScheme{
		Min: vipConf.GetInt(ConfExpireMin),
		Max: vipConf.GetInt(ConfExpireMax),
	}

	// create and populate pool, start keepalive, watch for signal close down
	pool := srvc.NewPool(scheme, sessConf, vipConf.GetInt(ConfPoolSize))
	runPool(pool)

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
		SessionPool: pool,
	}
	server.RegisterRoutes()
	fmt.Println("Begin on port:", port)
	http.ListenAndServe(port, m)
}
