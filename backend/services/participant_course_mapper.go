package services

import (
	"ce-progress-tracker/database/mongo"
	"ce-progress-tracker/database/mssql"
	"time"
)

type ParticipantMapperInterface interface {
	MongoToSQL(record *mongo.TakesAnonymized) *mssql.ParticipantCourse
}

// ParticipantCourseMapper handles conversion between MongoDB and SQL data structures
type ParticipantCourseMapper struct{}

func NewParticipantCourseMapper() ParticipantMapperInterface {
	return &ParticipantCourseMapper{}
}

// MongoToSQL converts a MongoDB TakesAnonymized record to a SQL ParticipantCourse
func (m *ParticipantCourseMapper) MongoToSQL(record *mongo.TakesAnonymized) *mssql.ParticipantCourse {
	if record == nil || record.ParticipantData == nil || record.CourseData == nil {
		return nil
	}

	return &mssql.ParticipantCourse{
		ParticipantID:     record.ParticipantData.ParticipantID,
		CourseID:          record.CourseData.CourseID,
		DateFirstAccessed: m.getUnixTimestamp(record.ParticipantData.DateFirstAccessed),
		DateLastAccessed:  m.getUnixTimestamp(record.ParticipantData.DateLastAccessed),
		CourseCompletion:  float32(record.ParticipantData.CourseCompletion),
	}
}

// getUnixTimestamp Convert float64 timestamp (assumed milliseconds) to time.Time
func (m *ParticipantCourseMapper) getUnixTimestamp(timestamp float64) *time.Time {
	ret := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
	return &ret
}
