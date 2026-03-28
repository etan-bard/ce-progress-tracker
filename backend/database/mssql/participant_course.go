package mssql

import (
	"time"
)

// ParticipantCourse represents the [progress_tracker].[participant_to_course_map] table
type ParticipantCourse struct {
	ID                int        `db:"Id"`
	ParticipantID     int        `db:"ParticipantId"`
	CourseID          int        `db:"CourseId"`
	DateFirstAccessed *time.Time `db:"DateFirstAccessed"`
	DateLastAccessed  time.Time  `db:"DateLastAccessed"`
	CourseCompletion  float32    `db:"CourseCompletion"`
	Action            string     `db:"Action"`
}

// TableName overrides the default table name to include the schema
func (ParticipantCourse) TableName() string {
	return "ProgressTracker.ParticipantToCourseMap"
}
