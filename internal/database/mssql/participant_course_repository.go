package mssql

import (
	"fmt"
	"strings"
)

type ParticipantCourseRepositoryInterface interface {
	UpsertAll(entries *[]ParticipantCourse, insertedCount *int, updatedCount *int) error
}

// ParticipantCourseRepository handles database operations for ParticipantCourse
type ParticipantCourseRepository struct {
	db DBServiceInterface
}

// NewParticipantCourseRepository creates a new repository instance
func NewParticipantCourseRepository(db DBServiceInterface) *ParticipantCourseRepository {
	return &ParticipantCourseRepository{db: db}
}

// UpsertAll inserts or updates a list of ParticipantCourse entries and determines the action taken for each record
func (r *ParticipantCourseRepository) UpsertAll(entries *[]ParticipantCourse, insertedCount *int, updatedCount *int) error {
	if entries == nil || len(*entries) == 0 {
		return nil
	}

	tableName := ParticipantCourse{}.TableName()
	var valueStrings []string
	var valueArgs []any
	for _, entry := range *entries {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, entry.ParticipantID, entry.CourseID, entry.LastAccessedAt, entry.CourseCompletion)
	}

	// Values seems to have a limit of 1000. For simplicity, will handle this in the main func, but for general use
	// it would be ideal to do batch processing in this function.
	query := fmt.Sprintf(`
		MERGE INTO %s AS Target
		USING (VALUES %s) AS Source (ParticipantId, CourseId, LastAccessedAt, CourseCompletion)
		ON Target.ParticipantId = Source.ParticipantId AND Target.CourseId = Source.CourseId
		WHEN MATCHED THEN
			UPDATE SET 
				Target.LastAccessedAt = Source.LastAccessedAt,
				Target.CourseCompletion = Source.CourseCompletion
		WHEN NOT MATCHED THEN
			INSERT (ParticipantId, CourseId, LastAccessedAt, CourseCompletion)
			VALUES (Source.ParticipantId, Source.CourseId, Source.LastAccessedAt, Source.CourseCompletion)
		OUTPUT $action AS Action;`, tableName, strings.Join(valueStrings, ","))

	query = r.db.Rebind(query)

	var results []struct {
		Action string `db:"Action"`
	}

	if err := r.db.Select(&results, query, valueArgs...); err != nil {
		return err
	}

	for i, entry := range *entries {
		if i < len(results) {
			if entry.Action == "INSERT" && insertedCount != nil {
				*insertedCount++
			} else if entry.Action == "UPDATE" && updatedCount != nil {
				*updatedCount++
			}
		}
	}

	return nil
}
