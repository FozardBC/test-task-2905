package api

import (
	"app/internal/api/handlers/list"
	"app/internal/api/handlers/save"
	"app/internal/api/middleware/json"
	mwLogger "app/internal/api/middleware/logger"
	requestid "app/internal/api/middleware/requestID"
	"app/internal/services/quteos"
	"app/internal/storage"

	"log/slog"

	"github.com/gorilla/mux"
)

type API struct {
	Router  mux.Router
	Storage storage.Storage
	Service *quteos.Service
	Log     *slog.Logger
}

func New(storage storage.Storage, log *slog.Logger) *API {
	router := mux.NewRouter()

	api := &API{
		Router:  *router,
		Storage: storage,
		Log:     log,
	}

	api.Service = quteos.New(&storage, log)

	api.Middlewares()

	api.Endpoints()

	return api
}

func (a *API) Endpoints() {

	v1 := a.Router.PathPrefix("/api/v1").Subrouter()

	v1.Handle("/quotes", json.JSONContentTypeMW(save.New(a.Log, a.Service))).Methods("POST")
	v1.Handle("/quetos", json.JSONContentTypeMW(list.New(a.Log, a.Service))).Methods("GET")
	// a.Router.HandleFunc("/quote/{author=l;;l}", a.Servie.Delete).Methods("GET")
	// a.Router.HandleFunc("/quetos/{id}", a.Servie.Update).Methods("DELETE")
}

func (a *API) Middlewares() {
	a.Router.Use(
		requestid.RequestIdMw,
		mwLogger.New(a.Log),
	)

}
