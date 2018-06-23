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
	//vEndAfter    = time.Minute * 6
)

func setConfig() *viper.Viper {
	conf := viper.GetViper()

	conf.SetConfigName(ConfFile)
	conf.AddConfigPath("$HOME")
	conf.AddConfigPath(".")
	conf.BindEnv(ConfSabreUsername)
	conf.BindEnv(ConfSabrePassword)
	conf.BindEnv(ConfSabrePCC)

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

func main() {
	c := setConfig()
	scheme := srvc.ExpireScheme{
		Min: c.GetInt(ConfExpireMin),
		Max: c.GetInt(ConfExpireMax),
	}
	p := srvc.NewPool(
		scheme,
		c.GetInt(ConfPoolSize),
		c.GetString(ConfSabreURL),
		c.GetString(ConfClientURL),
		c.GetString(ConfSabrePCC),
		srvc.GenerateConversationID(c.GetString(ConfClientURL)),
		srvc.GenerateMessageID(),
		srvc.SabreTimeFormat(),
		c.GetString(ConfSabreUsername),
		c.GetString(ConfSabrePassword),
	)

	// background populating the pool
	// and setting up the teh keepvalid worker
	go func() {
		err := p.Populate()
		if err != nil {
			panic(err)
		}
		//srvc.Keepalive(p, vRepeatEvery, vEndAfter)
		srvc.Keepalive(p, vRepeatEvery)
	}()
	// close down session pool properly, validates against
	// sabre and re-allocates on sabre side
	cls := make(chan os.Signal)
	signal.Notify(cls, os.Interrupt)
	go func() {
		sig := <-cls
		fmt.Printf("Got %s signal. Closing down...\n", sig)
		p.Close()
		os.Exit(1)
	}()

	// deal with context??
	m := chi.NewRouter()
	//Logger Logs the start and end of each request with the elapsed processing time
	m.Use(middleware.Logger)
	// Heartbeat Monitoring endpoint to check the servers pulse
	m.Use(middleware.Heartbeat("/heartbeat"))

	server := app.Server{
		Config:      c,
		Mux:         m,
		SessionPool: p,
	}
	server.RegisterRoutes()
	http.ListenAndServe(":8080", m)
}
