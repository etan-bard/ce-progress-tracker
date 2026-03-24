package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TakesAnonymizedRepositoryInterface interface{
	GetCourseIDCursor(ctx context.Context, courseIDs *[]int, batchSize int) (*mongo.Cursor, error)
}

type TakesAnonymizedRepository struct {
	collection *mongo.Collection
}

func NewTakesAnonymizedRepository(mongo DbServiceInterface) *TakesAnonymizedRepository {
	return &TakesAnonymizedRepository{
		collection: mongo.GetCollection("takes_anonymized"),
	}
}

func (t *TakesAnonymizedRepository) GetCourseIDCursor(ctx context.Context, courseIDs *[]int, batchSize int) (*mongo.Cursor, error) {
	// Create filter based on courseIDs
	var filter map[string]interface{}
	
	if courseIDs != nil && len(*courseIDs) > 0 {
		filter = map[string]interface{}{
			"course_id": map[string]interface{}{
				"$in": *courseIDs,
			},
		}
	} else {
		// If courseIDs is nil or empty, return all documents
		filter = map[string]interface{}{}
	}

	return t.collection.Find(
		ctx,
		filter,
		options.Find().SetBatchSize(int32(batchSize)),
	)
}
