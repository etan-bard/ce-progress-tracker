package mongo

// TakesAnonymized represents an object from the takes_anonymized collection
type TakesAnonymized struct {
	CourseData      *CourseData      `bson:"course_data"`
	ParticipantData *ParticipantData `bson:"participant_data"`
}

type CourseData struct {
	CourseID int `bson:"course_id"`
}

type ParticipantData struct {
	ParticipantID     int     `bson:"participant_id"`
	DateFirstAccessed float64 `bson:"date_first_accessed"`
	DateLastAccessed  float64 `bson:"date_last_accessed"`
	CourseCompletion  float64 `bson:"course_completion"`
}
