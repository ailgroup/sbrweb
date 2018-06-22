package common

//"github.com/ailgroup/sbrweb/client/datastore"
//"github.com/gorilla/sessions"
//"github.com/isomorphicgo/isokit"
//"github.com/isomorphicgo/isokit"
import (
	"github.com/ailgroup/sbrweb/engine/srvc"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

type Env struct {
	//Router *httprouter.Router
	Mux         *chi.Mux
	SessionPool *srvc.SessionPool
	Config      *viper.Viper
	//DStore datastore.Datastore
	//FStore *sessions.FilesystemStore
	//Router *isokit.Router
	//TemplateSet *isokit.TemplateSet
	//logging...
	//tracking....
	//analytics...
}
