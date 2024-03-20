package handlers

import (
	middlewares2 "distributed.systems.labs/manager/internal/api/middlewares"
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/processing"
	"distributed.systems.labs/manager/internal/sending"
	"distributed.systems.labs/manager/internal/storage"
	"distributed.systems.labs/manager/internal/tasks"
	"distributed.systems.labs/shared/pkg/contracts"
	"distributed.systems.labs/shared/pkg/middlewares"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type crackHashRequest struct {
	Hash      string `json:"hash"`
	MaxLength int    `json:"maxLength"`
}

func (req crackHashRequest) validate() error {
	if req.Hash == "" {
		return contracts.ErrEmptyHashToCrack
	}
	if req.MaxLength < 0 {
		return contracts.ErrNegativeMaxLength
	}
	return nil
}

type crackHashResponse struct {
	RequestID uuid.UUID `json:"requestId"`
}

func HandleCrackHashRequest(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.GetLogger(r.Context())

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var req crackHashRequest
	err := decoder.Decode(&req)
	if err != nil {
		logger.Printf("failed to decode json body with err: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = req.validate()
	if err != nil {
		logger.Printf("validation failed with err: %s", err)
	}

	S := middlewares2.GetStorage(r.Context())

	metadata := storage.RequestMetadata{
		Hash:      req.Hash,
		MaxLength: req.MaxLength,
		Alphabet:  middlewares2.GetAlphabet(r.Context()),
		Status:    config.InProgress,
	}

	workers, err := config.GetWorkers()
	if err != nil {
		logger.Printf("failed to get workers: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	numWorkers := uint64(len(workers))
	if numWorkers == 0 {
		logger.Printf("no workers in config")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	taskSize, err := config.GetTaskSize()
	if err != nil {
		logger.Printf("failed to get task size: %s", err)
	}
	logger.Printf("taskSize = %v", taskSize)

	//preparedTasks := tasks.CalcTasksWithFixedLength(m.Alphabet.Length(), m.MaxLength, taskSize)
	preparedTasks := tasks.CalcTasksWithNumWorkers(metadata.Alphabet.Length(), metadata.MaxLength, numWorkers, 10)
	metadata.Tasks = preparedTasks
	builder := strings.Builder{}
	builder.WriteString("calculated tasks:\n")
	for i, t := range preparedTasks {
		_, _ = builder.WriteString(fmt.Sprintf("task %v: StartIndex = %v PartCount %v\n",
			i, t.StartIndex, t.PartCount))
	}
	logger.Println(builder.String())

	id, err := S.SaveNew(metadata)
	if err != nil {
		logger.Printf("failed to save request metadata: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	metadata.ID = id

	defaultLogger := log.Default()
	requestLogger := log.New(
		defaultLogger.Writer(),
		fmt.Sprintf("requestId: %s ", metadata.ID),
		defaultLogger.Flags()|log.Lmsgprefix)

	go func() {
		sending.SendLoop(requestLogger, preparedTasks, metadata)
		processing.RequestChecker(
			S.Ctx(),
			requestLogger,
			metadata.ID,
			processing.CalcTimeoutsWithNumWorkers(preparedTasks[0].PartCount, numWorkers),
			S)
	}()

	encoder := json.NewEncoder(w)
	rsp := crackHashResponse{RequestID: metadata.ID}
	_ = encoder.Encode(rsp)
	logger.Printf("requestId: %s OK", metadata.ID)
}
