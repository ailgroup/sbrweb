package app

import (
	"github.com/ailgroup/sbrweb/engine/srvc"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

type Server struct {
	Mux         *chi.Mux
	SessionPool *srvc.SessionPool
	VConfig     *viper.Viper
	SConfig     *srvc.SessionConf
	//DStore datastore.Datastore
	//FStore *sessions.FilesystemStore
	//logging...
	//tracking....
	//analytics...
}
