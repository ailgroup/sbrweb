package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ailgroup/sbrweb/engine/srvc"
	trns "github.com/ailgroup/sbrweb/transmission"
	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/spf13/viper"
)

const (
	//ConfSabreUsername sets env variable
	confSabreUsername = "SABRE_USERNAME"
	//ConfSabrePassword sets env variable
	confSabrePassword = "SABRE_PASSWORD"
	//ConfSabrePCC sets env variable
	confSabrePCC = "SABRE_PCC"
	//ConfFile name of file
	confFile = "config"
	//ConfTimeZone timezeone variable
	confTimeZone = "app_timezone"
	//ConfExpireMin key for toml file
	confExpireMin = "sessions.expire.min"
	//ConfExpireMax key for toml filee
	confExpireMax = "sessions.expire.max"
	//ConfPoolSize key for toml filee
	confPoolSize = "sessions.client.pool_size"
	//ConfSabreURL key for toml filee
	confSabreURL = "sessions.sabre_url"
	//ConfClientURL key for toml file
	confClientURL = "sessions.client.url"
)

var (
	vRepeatEvery      = time.Minute * 3
	sessConf          = &srvc.SessionConf{}
	vipConf           = &viper.Viper{}
	port              = ":8080"
	clientAppTimeZone = &time.Location{}
)

func init() {
	clientAppTimeZone = time.UTC
	vipConf = setConfig()
	sessConf = &srvc.SessionConf{
		ServiceURL:  vipConf.GetString(confSabreURL),
		From:        vipConf.GetString(confClientURL),
		PCC:         vipConf.GetString(confSabrePCC),
		Convid:      srvc.GenerateConversationID(vipConf.GetString(confClientURL)),
		Msgid:       srvc.GenerateMessageID(),
		Timestr:     srvc.SabreTimeNowFmt(),
		Username:    vipConf.GetString(confSabreUsername),
		Password:    vipConf.GetString(confSabrePassword),
		AppTimeZone: clientAppTimeZone,
	}
}

func setConfig() *viper.Viper {
	vip := viper.GetViper()

	vip.SetConfigName(confFile)
	vip.AddConfigPath("$HOME")
	vip.AddConfigPath(".")
	err := vip.BindEnv(confSabreUsername)
	if err != nil {
		log.Printf("ERROR ConfSabreUsername not found %v", err)
	}
	err = vip.BindEnv(confSabrePassword)
	if err != nil {
		log.Printf("ERROR ConfSabrePassword not found %v", err)
	}
	err = vip.BindEnv(confSabrePCC)
	if err != nil {
		log.Printf("ERROR ConfSabrePCC not found: %v", err)
	}
	vip.SetDefault(confTimeZone, "UTC")
	err = vip.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("FATAL ERROR CONFIG FILE: %s", err))
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
		Min: vipConf.GetInt(confExpireMin),
		Max: vipConf.GetInt(confExpireMax),
	}

	// create and populate pool, start keepalive, watch for signal close down
	pool := srvc.NewPool(scheme, sessConf, vipConf.GetInt(confPoolSize))
	runPool(pool)

	router := chi.NewRouter()
	//render.SetContentType Sets content-type as... application/json
	//Logger Logs the start and end of each request with the elapsed processing time
	//Heartbeat Monitoring endpoint to check the servers pulse
	//Recoverer Recovers from panics without crashing
	//DefaultCompress Compress results gzipping assets and json
	//RedirectSlashes Redirect slashes to no slash URL versions
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.Heartbeat("/heartbeat"),
		middleware.Recoverer,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
	)

	server := trns.Server{
		VConfig:     vipConf,
		SConfig:     sessConf,
		Mux:         router,
		SessionPool: pool,
	}
	server.RegisterRoutes()
	fmt.Println("Begin on port:", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		panic(fmt.Errorf("FATAL HTTP ERROR: %s", err))
	}
}
