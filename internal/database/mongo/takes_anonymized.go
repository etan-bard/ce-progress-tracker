package mongo

import "time"

// TakesAnonymized represents an object from the takes_anonymized collection
type TakesAnonymized struct {
	CourseData      *CourseData      `bson:"course_data"`
	ParticipantData *ParticipantData `bson:"participant_data"`
}

type CourseData struct {
	CourseID int `bson:"course_id"`
}

type ParticipantData struct {
	ParticipantID    int        `bson:"participant_id"`
	LastAccessedAt   *time.Time `bson:"last_accessed_at"`
	CourseCompletion bool       `bson:"course_completion"`
}
