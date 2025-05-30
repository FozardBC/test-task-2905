package api

import (
	"app/internal/api/handlers/delete"
	"app/internal/api/handlers/list"
	"app/internal/api/handlers/random"
	"app/internal/api/handlers/save"
	"app/internal/api/middleware/json"
	mwLogger "app/internal/api/middleware/logger"
	requestid "app/internal/api/middleware/requestID"
	"app/internal/services/quteos"
	"app/internal/storage"
	"fmt"
	"net/http"

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

	v1.Handle("/quotes", json.JSONContentTypeMW(save.New(a.Log, a.Service))).Methods(http.MethodPost)
	v1.Handle("/quotes", json.JSONContentTypeMW(list.New(a.Log, a.Service))).Methods(http.MethodGet)
	v1.Handle("/quotes/random", json.JSONContentTypeMW(random.New(a.Log, a.Service))).Methods(http.MethodGet)
	v1.Handle("/quotos/{id:[0-9]+}", json.JSONContentTypeMW(delete.New(a.Log, a.Service))).Methods(http.MethodDelete)

	Routes(a.Log, &a.Router)
}

func (a *API) Middlewares() {
	a.Router.Use(
		requestid.RequestIdMw,
		mwLogger.New(a.Log),
	)

}

func Routes(log *slog.Logger, router *mux.Router) {
	log.Info("Available routes:")
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Printf("→ %-6s %s\n", methods[0], pathTemplate)
			} else {
				fmt.Printf("→ ALL   %s\n", pathTemplate) // Если методы не указаны (например, HandleFunc без Methods())
			}
		}
		return nil
	})
	if err != nil {
		log.Info("Error printing routes:", err)
	}
}
