package storage

import (
	"distributed.systems.labs/manager/internal/config"
	"distributed.systems.labs/manager/internal/tasks"
	"distributed.systems.labs/shared/pkg/alphabet"
	"github.com/google/uuid"
)

type RequestMetadata struct {
	Alphabet  alphabet.Alphabet
	ID        uuid.UUID            `bson:"_id"`
	Status    config.RequestStatus `bson:"status"`
	Cracks    []string
	Hash      string
	MaxLength int
	Tasks     []tasks.Task
}
