package services

import (
	mongo2 "ce-progress-tracker/database/mongo"
	"ce-progress-tracker/database/mssql"
	"context"
	"sync"

	"go.uber.org/zap"
)

type DataMigrationStrategyInterface interface {
	Execute(ctx context.Context, courseIDs *[]int) (int, int, int, error)
}

type BatchDataMigrationStrategy struct {
	participantCourseMapper     ParticipantMapperInterface
	participantCourseRepository mssql.ParticipantCourseRepositoryInterface
	takesAnonymizedRepository   mongo2.TakesAnonymizedRepositoryInterface
	logger                      *zap.Logger
	batchSize                   int
	maxGoroutines               int
}

func NewBatchDataMigrationStrategy(
	participantCourseMapper ParticipantMapperInterface,
	participantCourseRepository mssql.ParticipantCourseRepositoryInterface,
	takesAnonymizedRepository mongo2.TakesAnonymizedRepositoryInterface,
	logger *zap.Logger,
	batchSize int,
	maxGoroutines int,
) *BatchDataMigrationStrategy {
	return &BatchDataMigrationStrategy{
		participantCourseMapper:     participantCourseMapper,
		participantCourseRepository: participantCourseRepository,
		takesAnonymizedRepository:   takesAnonymizedRepository,
		logger:                      logger,
		batchSize:                   batchSize,
		maxGoroutines:               maxGoroutines,
	}
}

func (b *BatchDataMigrationStrategy) Execute(ctx context.Context, courseIDs *[]int) (int, int, int, error) {
	// Retrieves data from MongoDB and patches into SQL Server in batches
	// This will not snapshot. Options: pass a snapshot=true flag (if supported in MongoDB)
	cursor, err := b.takesAnonymizedRepository.GetCourseIDCursor(ctx, courseIDs, b.batchSize)
	if err != nil {
		b.logger.Fatal("Failed to query MongoDB collection", zap.Error(err))
	}
	defer cursor.Close(context.Background())

	var totalInserted, totalUpdated, totalSkipped int
	batch := make([]mssql.ParticipantCourse, 0, b.batchSize)

	var wg sync.WaitGroup
	totalsMutex := sync.Mutex{}

	// Semaphore to limit concurrent goroutines
	sem := make(chan struct{}, b.maxGoroutines)

	// Iterate over the input data
	for cursor.Next(context.Background()) {
		var mongoRecord mongo2.TakesAnonymized
		if err := cursor.Decode(&mongoRecord); err != nil {
			return 0, 0, 0, err
		}

		// Map from MongoDB to SQL structure
		sqlRecord := b.participantCourseMapper.MongoToSQL(&mongoRecord)
		if sqlRecord == nil {
			b.logger.Warn("Skipping record with missing required fields")
			continue
		}

		batch = append(batch, *sqlRecord)

		// Process the batch if it hits the size limit
		if len(batch) == b.batchSize {
			currentBatch := batch
			wg.Add(1)
			sem <- struct{}{} // Acquire semaphore
			go func() {
				defer wg.Done()
				defer func() { <-sem }() // Release semaphore
				b.processBatch(ctx, currentBatch, &totalsMutex, &totalInserted, &totalUpdated, &totalSkipped)
			}()

			// Create a new batch for the next iteration
			batch = make([]mssql.ParticipantCourse, 0, b.batchSize)
		}
	}

	// Process any remaining records in the last batch
	if len(batch) > 0 {
		currentBatch := batch
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore
		go func() {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore
			b.processBatch(ctx, currentBatch, &totalsMutex, &totalInserted, &totalUpdated, &totalSkipped)
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return totalInserted, totalUpdated, totalSkipped, nil
}

func (b *BatchDataMigrationStrategy) processBatch(ctx context.Context, batch []mssql.ParticipantCourse, mutex *sync.Mutex, totalInserted, totalUpdated, totalSkipped *int) {
	if len(batch) == 0 {
		return
	}

	var inserted, updated, skipped int
	if err := b.participantCourseRepository.UpsertAll(&batch, &inserted, &updated, &skipped); err != nil {
		b.logger.Error("Error processing batch", zap.Error(err))
		return
	}

	mutex.Lock()
	*totalInserted += inserted
	*totalUpdated += updated
	*totalSkipped += skipped
	mutex.Unlock()

	b.logger.Info("Batch processed", zap.Int("inserted", inserted), zap.Int("updated", updated), zap.Int("skipped", skipped))
}
