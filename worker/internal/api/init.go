package api

import (
	"distributed.systems.labs/shared/pkg/middlewares"
	"distributed.systems.labs/worker/internal/notify"
	"github.com/gorilla/mux"
	"net/http"
)

func ConfigureEndpoints(mn notify.Notifier) *mux.Router {
	r := mux.NewRouter()

	m := notifierMiddleware{mn: mn}
	r.Use(middlewares.LoggerMiddleware)
	r.Use(m.Middleware)

	r.HandleFunc("/", indexHandler) // for testing
	r.HandleFunc("/internal/api/worker/hash/crack/task", acceptTask).Methods(http.MethodPost)
	return r
}
