package mssql

import (
	"time"
)

// ParticipantCourse represents the [progress_tracker].[participant_to_course_map] table
type ParticipantCourse struct {
	ID               int        `db:"Id"`
	ParticipantID    int        `db:"ParticipantId"`
	CourseID         int        `db:"CourseId"`
	LastAccessedAt   *time.Time `db:"LastAccessedAt"`
	CourseCompletion bool       `db:"CourseCompletion"`
	Action           string     `db:"Action"`
}

// TableName overrides the default table name to include the schema
func (ParticipantCourse) TableName() string {
	return "progress_tracker.participant_to_course_map"
}
