package mssql

import (
	"fmt"
	"strings"
)

type ParticipantCourseRepositoryInterface interface {
	GetAll() (*[]ParticipantCourse, error)
	UpsertAll(entries *[]ParticipantCourse, insertedCount *int, updatedCount *int, skippedCount *int) error
}

// ParticipantCourseRepository handles database operations for ParticipantCourse
type ParticipantCourseRepository struct {
	db        DBServiceInterface
	tableName string
}

// NewParticipantCourseRepository creates a new participantCourseRepository instance
func NewParticipantCourseRepository(db DBServiceInterface) *ParticipantCourseRepository {
	return &ParticipantCourseRepository{
		db:        db,
		tableName: ParticipantCourse{}.TableName()}
}

func (r *ParticipantCourseRepository) GetAll() (*[]ParticipantCourse, error) {
	var results []ParticipantCourse
	query := fmt.Sprintf("SELECT * FROM %s", r.tableName)
	if err := r.db.Select(&results, query); err != nil {
		return nil, err
	}
	return &results, nil
}

// UpsertAll inserts or updates a list of ParticipantCourse entries and determines the action taken for each record
func (r *ParticipantCourseRepository) UpsertAll(entries *[]ParticipantCourse, insertedCount *int, updatedCount *int, skippedCount *int) error {
	if entries == nil || len(*entries) == 0 {
		return nil
	}

	// Max of 2000 parameters. For simplicity, will handle this in the main func, but for general use
	// it would be ideal to do batch processing in this function.
	var valueStrings []string
	var valueArgs []any
	for _, entry := range *entries {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, entry.ParticipantID, entry.CourseID, entry.DateLastAccessed, entry.CourseCompletion)
	}

	if len(valueArgs) != len(valueStrings)*4 {
		return fmt.Errorf("argument count mismatch: expected %d, got %d", len(valueStrings)*4, len(valueArgs))
	}

	// Merges values into the DB table, returning a simple count of if a row is updated or inserted.
	query := fmt.Sprintf(`
		MERGE INTO %s AS Target
		USING (
    		SELECT 
        		CAST(v.ParticipantId AS INT), 
        		CAST(v.CourseId AS INT), 
        		CAST(v.DateLastAccessed AS DATETIME2),
        		CAST(v.CourseCompletion AS REAL)
    		FROM (VALUES %s) AS v (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
		) AS Source (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
		ON Target.ParticipantId = Source.ParticipantId AND Target.CourseId = Source.CourseId
		WHEN MATCHED AND (
			Target.DateLastAccessed IS DISTINCT FROM Source.DateLastAccessed 
			OR Target.CourseCompletion IS DISTINCT FROM Source.CourseCompletion) THEN
			UPDATE SET 
				Target.DateLastAccessed = Source.DateLastAccessed,
				Target.CourseCompletion = Source.CourseCompletion
		WHEN NOT MATCHED BY TARGET THEN
			INSERT (ParticipantId, CourseId, DateLastAccessed, CourseCompletion)
			VALUES (Source.ParticipantId, Source.CourseId, Source.DateLastAccessed, Source.CourseCompletion)
		OUTPUT $action AS Action;`, r.tableName, strings.Join(valueStrings, ","))

	query = r.db.Rebind(query)

	var results []struct {
		Action string `db:"Action"`
	}

	if err := r.db.Select(&results, query, valueArgs...); err != nil {
		return err
	}

	for i, result := range results {
		if i < len(*entries) {
			if result.Action == "INSERT" && insertedCount != nil {
				*insertedCount++
			} else if result.Action == "UPDATE" && updatedCount != nil {
				*updatedCount++
			}
		}
	}

	// Any remainders go into the skipped count
	*skippedCount += len(*entries) - *insertedCount - *updatedCount

	return nil
}
