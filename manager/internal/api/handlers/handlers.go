package handlers

import (
	"distributed.systems.labs/shared/pkg/middlewares"
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	logger := middlewares.GetLogger(r.Context())
	_, err := fmt.Fprint(w, "Hello, world form manager!")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Println("OK")
}
