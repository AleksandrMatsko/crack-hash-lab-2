package storage

import (
	"context"
	"distributed.systems.labs/manager/internal/config"
	"github.com/google/uuid"
	"log"
)

type Storage interface {
	Atomically(reqID uuid.UUID, fn func(req *RequestMetadata)) (RequestMetadata, error)
	AddCracks(reqID uuid.UUID, cracks []string, startIndex uint64) error
	Get(reqID uuid.UUID) (RequestMetadata, bool, error)
	Ctx() context.Context
	// SaveNew should save provided metadata to storage and return generated id
	SaveNew(metadata RequestMetadata) (uuid.UUID, error)
}

func SetStatusErrAndSave(logger *log.Logger, S Storage, requestID uuid.UUID) {
	_, err := S.Atomically(requestID, func(req *RequestMetadata) {
		req.Status = config.Error
	})
	if err != nil {
		logger.Printf("failed to save request metadata: %s", err)
	}
}
