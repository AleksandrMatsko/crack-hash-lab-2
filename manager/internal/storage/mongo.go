package storage

import (
	"context"
	"distributed.systems.labs/manager/internal/config"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
)

type MongoStorage struct {
	ctx            context.Context
	client         *mongo.Client
	dbName         string
	collectionName string
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
		ctx:            ctx,
		client:         client,
		dbName:         config.GetMongoDBName(),
		collectionName: config.GetMongoDBCollectionName(),
	}, nil
}

func (m *MongoStorage) Atomically(reqID uuid.UUID, fn func(req *RequestMetadata) error) (RequestMetadata, error) {
	wc := writeconcern.Majority()
	txOptions := options.Transaction().
		SetReadPreference(readpref.Primary()).
		SetWriteConcern(wc)

	session, err := m.client.StartSession(
		options.Session().
			SetDefaultReadPreference(readpref.Primary()).
			SetDefaultWriteConcern(wc))
	if err != nil {
		return RequestMetadata{}, err
	}
	defer session.EndSession(m.Ctx())

	var result RequestMetadata
	err = mongo.WithSession(m.Ctx(), session, func(mongoCtx mongo.SessionContext) error {
		err = session.StartTransaction(txOptions)
		if err != nil {
			return err
		}

		filter := bson.D{{"_id", reqID}}
		collection := m.getCollection()
		var metadata RequestMetadata
		err = collection.FindOne(mongoCtx, filter).Decode(&metadata)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return ErrNoSuchRequest
			}
			return err
		}
		// call function to do
		err = fn(&metadata)
		if err != nil {
			return err
		}

		_, err := collection.ReplaceOne(mongoCtx, filter, metadata)
		if err != nil {
			return err
		}

		err = session.CommitTransaction(mongoCtx)
		if err != nil {
			return err
		}
		result = metadata
		return nil
	})
	if err != nil {
		_ = session.AbortTransaction(m.Ctx())
		log.Printf("err type = %T", err)
		return RequestMetadata{}, fmt.Errorf("requestID: %s error while executing transaction: %s", reqID, err)
	}
	return result, nil
}

func (m *MongoStorage) AddCracks(reqID uuid.UUID, cracks []string, startIndex uint64) error {
	var numDoneTasks, tasksAmount int
	_, err := m.Atomically(reqID, func(metadata *RequestMetadata) error {
		for _, c := range cracks {
			if len([]rune(c)) > metadata.MaxLength {
				return ErrTooLongCrack
			}
		}

		numDone := 0
		for i := range metadata.Tasks {
			if metadata.Tasks[i].StartIndex == startIndex && !metadata.Tasks[i].Done {
				metadata.Tasks[i].Done = true
				metadata.Cracks = append(metadata.Cracks, cracks...)
			}
			if metadata.Tasks[i].Done {
				numDone += 1
			}
		}

		if numDone == len(metadata.Tasks) {
			metadata.Status = config.Ready
		}
		numDoneTasks = numDone
		tasksAmount = len(metadata.Tasks)
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("requestId: %s tasks done %v / %v", reqID, numDoneTasks, tasksAmount)
	return nil
}

func (m *MongoStorage) Get(reqID uuid.UUID) (RequestMetadata, bool, error) {
	filter := bson.D{{"_id", reqID}}

	collection := m.getCollection()
	var metadata RequestMetadata
	err := collection.FindOne(m.Ctx(), filter).Decode(&metadata)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return RequestMetadata{}, false, nil
		}
		return RequestMetadata{}, false, err
	}
	return metadata, true, nil
}

func (m *MongoStorage) Ctx() context.Context {
	return m.ctx
}

func (m *MongoStorage) SaveNew(metadata RequestMetadata) (uuid.UUID, error) {
	id := uuid.New()
	metadata.ID = id

	collection := m.getCollection()
	_, err := collection.InsertOne(m.Ctx(), metadata)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (m *MongoStorage) Close() {
	_ = m.client.Disconnect(m.ctx)
}

func (m *MongoStorage) getCollection() *mongo.Collection {
	return m.client.
		Database(
			m.dbName,
			options.Database().
				SetReadPreference(readpref.Primary()),
		).
		Collection(
			m.collectionName,
			options.Collection().
				SetReadPreference(readpref.Primary()))
}
