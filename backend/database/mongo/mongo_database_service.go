package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBServiceInterface interface {
	GetCollection(name string) CollectionInterface
	Close(ctx context.Context) error
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (CursorInterface, error)
}

type CursorInterface interface {
	Next(ctx context.Context) bool
	Decode(result any) error
	Close(ctx context.Context) error
}

type CollectionInterface interface {
	Find(ctx context.Context, filter any, opts ...*options.FindOptions) (CursorInterface, error)
}

type mongoCollection struct {
	*mongo.Collection
}

func (c *mongoCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (CursorInterface, error) {
	return c.Collection.Find(ctx, filter, opts...)
}

type DBService struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDBService(ctx context.Context, uri, dbName string) (DBServiceInterface, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &DBService{
		client: client,
		db:     client.Database(dbName),
	}, nil
}

func (s *DBService) GetCollection(name string) CollectionInterface {
	return &mongoCollection{s.db.Collection(name)}
}

func (s *DBService) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *DBService) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (CursorInterface, error) {
	return s.GetCollection("takes_anonymized").Find(ctx, filter, opts...)
}
