package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbServiceInterface interface {
	GetCollection(name string) *mongo.Collection
	Close(ctx context.Context) error
}

type DBService struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDBService(ctx context.Context, uri, dbName string) (DbServiceInterface, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &DBService{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (s *DBService) GetCollection(name string) *mongo.Collection {
	return s.db.Collection(name)
}

func (s *DBService) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
