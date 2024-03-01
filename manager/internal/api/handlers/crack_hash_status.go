package handlers

import (
	middlewares2 "distributed.systems.labs/manager/internal/api/middlewares"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/shared/pkg/middlewares"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type statusResponse struct {
	Status config.RequestStatus `json:"status"`
	Data   *[]string            `json:"data"`
}

func HandleCrackHashStatus(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	logger := middlewares.GetLogger(r.Context())

	str := r.URL.Query().Get("requestId")
	if str == "" {
		logger.Println("no query parameters provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqID, err := uuid.Parse(str)
	if err != nil {
		logger.Printf("error while parsing uuid: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storage := middlewares2.GetStorage(r.Context())
	metadata, ok, err := storage.Get(reqID)
	if err != nil {
		logger.Printf("error while getting request metadata: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		logger.Printf("no request with id %s", reqID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var rsp statusResponse
	rsp.Status = metadata.Status
	rsp.Data = nil
	if rsp.Status == config.Ready {
		rsp.Data = &metadata.Cracks
	}
	_ = encoder.Encode(rsp)
	logger.Printf("requestId %s OK", reqID)
}
