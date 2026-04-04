package mssql

import (
	"fmt"
	"os"
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
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, entry.ParticipantID, entry.CourseID, entry.DateFirstAccessed, entry.DateLastAccessed, entry.CourseCompletion)
	}

	// Load query from SQL file
	queryBytes, err := os.ReadFile("./database/mssql/queries/upsert_participant_courses.sql")
	if err != nil {
		return fmt.Errorf("failed to read SQL query file: %w", err)
	}
	query := fmt.Sprintf(string(queryBytes), r.tableName, strings.Join(valueStrings, ","))

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
