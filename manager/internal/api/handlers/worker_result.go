package handlers

import (
	middlewares2 "distributed.systems.labs/manager/internal/api/middlewares"
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/shared/pkg/middlewares"
	"encoding/json"
	"net/http"
)

func HandleWorkerResult(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLogger(r.Context())

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req contracts.TaskResultRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Printf("failed to decode request body to json: %s", err)
		logger.Printf("status %v", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	store := middlewares2.GetStorage(r.Context())
	err = store.AddCracks(req.RequestID, req.Cracks, req.StartIndex)
	if err != nil {
		logger.Printf("requestId %s failed to add new cracks %s", req.RequestID, err)
		logger.Printf("status %v", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Printf("requestId %s OK", req.RequestID)
}
