package services

import (
	mongo2 "ce-progress-tracker/internal/database/mongo"
	"ce-progress-tracker/internal/database/mssql"
	"context"
	"sync"

	"go.uber.org/zap"
)

type DataMigrationStrategyInterface interface {
	Execute(ctx context.Context) (int, int, error)
}

type BatchDataMigrationStrategy struct {
	participantCourseMapper     ParticipantMapperInterface
	participantCourseRepository mssql.ParticipantCourseRepositoryInterface
	takesAnonymizedRepository   mongo2.TakesAnonymizedRepositoryInterface
	logger                      *zap.Logger
	batchSize                   int
}

func NewBatchDataMigrationStrategy(
	participantCourseMapper ParticipantMapperInterface,
	participantCourseRepository mssql.ParticipantCourseRepositoryInterface,
	takesAnonymizedRepository mongo2.TakesAnonymizedRepositoryInterface,
	logger *zap.Logger,
	batchSize int,
) *BatchDataMigrationStrategy {
	return &BatchDataMigrationStrategy{
		participantCourseMapper:     participantCourseMapper,
		participantCourseRepository: participantCourseRepository,
		takesAnonymizedRepository:   takesAnonymizedRepository,
		logger:                      logger,
		batchSize:                   batchSize,
	}
}

func (b *BatchDataMigrationStrategy) Execute(ctx context.Context, courseIDs *[]int) (int, int, error) {
	// Retrieves data from MongoDB and patches into SQL Server in batches
	// This will not snapshot. Options: pass a snapshot=true flag (if supported in MongoDB)
	cursor, err := b.takesAnonymizedRepository.GetCourseIDCursor(ctx, courseIDs, b.batchSize)
	defer cursor.Close(context.Background())
	if err != nil {
		b.logger.Fatal("Failed to query MongoDB collection", zap.Error(err))
	}

	var totalInserted, totalUpdated int
	batch := make([]mssql.ParticipantCourse, 0, b.batchSize)

	wg := sync.WaitGroup{}
	totalsMutex := sync.Mutex{}

	// Iterate over the input data
	for cursor.Next(context.Background()) {
		var mongoRecord mongo2.TakesAnonymized
		if err := cursor.Decode(&mongoRecord); err != nil {
			return 0, 0, err
		}

		// Map from MongoDB to SQL structure
		sqlRecord := b.participantCourseMapper.MongoToSQL(&mongoRecord)
		if sqlRecord == nil {
			b.logger.Warn("Skipping record with missing required fields")
			continue
		}

		batch = append(batch, *sqlRecord)
		if len(batch) == b.batchSize || cursor.RemainingBatchLength() == 0 {
			wg.Add(1)
			go func(currentBatch []mssql.ParticipantCourse) {
				defer wg.Done()
				b.processBatch(ctx, currentBatch, &totalsMutex, &totalInserted, &totalUpdated)
			}(batch)

			// Create a new batch for the next iteration
			batch = make([]mssql.ParticipantCourse, 0, b.batchSize)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return totalInserted, totalUpdated, nil
}

func (b *BatchDataMigrationStrategy) processBatch(ctx context.Context, batch []mssql.ParticipantCourse, mutex *sync.Mutex, totalInserted, totalUpdated *int) {
	if len(batch) == 0 {
		return
	}

	var inserted, updated int
	if err := b.participantCourseRepository.UpsertAll(&batch, &inserted, &updated); err != nil {
		b.logger.Error("Error processing batch", zap.Error(err))
		return
	}

	mutex.Lock()
	*totalInserted += inserted
	*totalUpdated += updated
	mutex.Unlock()

	b.logger.Info("Batch processed", zap.Int("inserted", inserted), zap.Int("updated", updated))
}
