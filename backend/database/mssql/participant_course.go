package mssql

import (
	"time"
)

// ParticipantCourse represents the [progress_tracker].[participant_to_course_map] table
type ParticipantCourse struct {
	ID                int        `db:"Id" json:"-"`
	ParticipantID     int        `db:"ParticipantId" json:"participantId"`
	CourseID          int        `db:"CourseId" json:"courseId"`
	DateFirstAccessed *time.Time `db:"DateFirstAccessed" json:"dateFirstAccessed"`
	DateLastAccessed  *time.Time `db:"DateLastAccessed" json:"dateLastAccessed"`
	CourseCompletion  float32    `db:"CourseCompletion" json:"courseCompletion"`
	Action            string     `db:"Action" json:"-"`
}

// TableName overrides the default table name to include the schema
func (ParticipantCourse) TableName() string {
	return "ProgressTracker.ParticipantToCourseMap"
}
