package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TakesAnonymizedRepositoryInterface interface {
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
	var filter = make(map[string]interface{})

	// Add the course IDs if present, otherwise we retrieve everything
	if courseIDs != nil && len(*courseIDs) > 0 {
		filter["course_data.course_id"] = map[string]interface{}{
			"$in": *courseIDs,
		}
	}

	// Add filter for participant_data.date_first_accessed to be greater than 0
	// since there are 6 duplicate records with a null date_first_accessed
	filter["participant_data.date_first_accessed"] = map[string]interface{}{
		"$gt": float64(0),
	}

	return t.collection.Find(
		ctx,
		filter,
		options.Find().SetBatchSize(int32(batchSize)),
	)
}
