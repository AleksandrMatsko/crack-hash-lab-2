package api

import (
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/shared/pkg/middlewares"
	"distributed.systems.labs/worker/internal/calc"
	"encoding/json"
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	logger := middlewares.GetLogger(r.Context())
	_, err := fmt.Fprint(w, "Hello, world form worker!")
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	logger.Printf("OK")

}

func acceptTask(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLogger(r.Context())

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req contracts.TaskRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Printf("error while decoding json in body: %s", err)
		logger.Printf("responsing with status code %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = req.Validate()
	if err != nil {
		logger.Printf("validation failed with err: %s", err)
		logger.Printf("responsing with status code %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notifier := getNotifier(r.Context())
	logger.Printf("starting crack for hash: %s, request-id: %s ...", req.ToCrack, req.RequestID)
	go calc.ProcessRequest(notifier.Context(), req, notifier.GetResChan())
	logger.Printf("request-id: %s status code %v", req.RequestID, http.StatusOK)
}
