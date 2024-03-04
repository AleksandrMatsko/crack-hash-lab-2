package storage

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoStorage struct {
	ctx    context.Context
	client *mongo.Client
}

var _ Storage = &MongoStorage{}

func InitMongoStorage(ctx context.Context, connStr string) (*MongoStorage, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return &MongoStorage{
		ctx:    ctx,
		client: client,
	}, nil
}

func (m *MongoStorage) Atomically(reqID uuid.UUID, fn func(req *RequestMetadata)) (RequestMetadata, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoStorage) AddCracks(reqID uuid.UUID, cracks []string, startIndex uint64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MongoStorage) Get(reqID uuid.UUID) (RequestMetadata, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoStorage) Ctx() context.Context {
	return m.Ctx()
}

func (m *MongoStorage) SaveNew(metadata RequestMetadata) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoStorage) Close() {
	_ = m.client.Disconnect(m.ctx)
}
