package api

import (
	"distributed.systems.labs/manager/internal/api/handlers"
	middlewares2 "distributed.systems.labs/manager/internal/api/middlewares"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/shared/pkg/alphabet"
	"distributed.systems.labs/shared/pkg/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func ConfigureEndpoints(s storage.Storage, a alphabet.Alphabet) *mux.Router {
	r := mux.NewRouter()

	sm := middlewares2.StorageMiddleware{S: s}
	am := middlewares2.AlphabetMiddleware{A: a}

	r.Use(middlewares.LoggerMiddleware)
	r.Use(sm.Middleware)
	r.Use(am.Middleware)

	r.HandleFunc("/", handlers.IndexHandler) // for testing

	r.HandleFunc("/api/hash/crack", handlers.HandleCrackHashRequest).Methods(http.MethodPost)
	r.HandleFunc("/api/hash/status", handlers.HandleCrackHashStatus).Methods(http.MethodGet)

	return r
}
