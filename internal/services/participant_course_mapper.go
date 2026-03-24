package services

import (
	"ce-progress-tracker/internal/database/mongo"
	"ce-progress-tracker/internal/database/mssql"
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
		ParticipantID:    record.ParticipantData.ParticipantID,
		CourseID:         record.CourseData.CourseID,
		LastAccessedAt:   record.ParticipantData.LastAccessedAt,
		CourseCompletion: record.ParticipantData.CourseCompletion,
	}
}
