package storage

import (
	"context"
	"distributed.systems.labs/manager/internal/config"
	"github.com/google/uuid"
	"log"
	"sync"
)

type InMemoryStorage struct {
	ctx    context.Context
	locker sync.Locker
	data   map[uuid.UUID]RequestMetadata
}

// InMemoryStorage should implement Storage
var _ Storage = &InMemoryStorage{}

func InitInMemoryStorage(ctx context.Context) *InMemoryStorage {
	return &InMemoryStorage{
		ctx:    ctx,
		locker: &sync.Mutex{},
		data:   make(map[uuid.UUID]RequestMetadata),
	}
}

func (s *InMemoryStorage) SaveNew(metadata RequestMetadata) (uuid.UUID, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	id := uuid.New()
	_, ok := s.data[id]
	for ok {
		id = uuid.New()
		_, ok = s.data[id]
	}
	metadata.ID = id
	s.data[id] = metadata
	return id, nil
}

func (s *InMemoryStorage) AddCracks(reqID uuid.UUID, cracks []string, startIndex uint64) error {
	s.locker.Lock()
	data, ok := s.data[reqID]
	s.locker.Unlock()

	if !ok {
		return ErrNoSuchRequest
	}
	for _, c := range cracks {
		if len([]rune(c)) > data.MaxLength {
			return ErrTooLongCrack
		}
	}

	s.locker.Lock()
	defer s.locker.Unlock()

	data, ok = s.data[reqID]
	if !ok {
		return ErrNoSuchRequest
	}
	numDone := 0
	for i := range data.Tasks {
		if data.Tasks[i].StartIndex == startIndex && !data.Tasks[i].Done {
			data.Tasks[i].Done = true
			data.Cracks = append(data.Cracks, cracks...)
		}
		if data.Tasks[i].Done {
			numDone += 1
		}
	}
	defer log.Printf("requestId: %s tasks done %v / %v", data.ID, numDone, len(data.Tasks))

	if numDone == len(data.Tasks) {
		data.Status = config.Ready
	}
	s.data[reqID] = data

	return nil
}

func (s *InMemoryStorage) Get(reqID uuid.UUID) (RequestMetadata, bool, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	data, ok := s.data[reqID]
	return data, ok, nil
}

func (s *InMemoryStorage) Ctx() context.Context {
	return s.ctx
}

func (s *InMemoryStorage) Atomically(reqID uuid.UUID, fn func(req *RequestMetadata)) (RequestMetadata, error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	data, ok := s.data[reqID]
	if ok {
		fn(&data)
		s.data[reqID] = data
		return data, nil
	}
	return RequestMetadata{}, ErrNoSuchRequest
}
