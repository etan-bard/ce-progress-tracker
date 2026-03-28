package mssql

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

// ParticipantCourseController manages APIs related to participant-course mappings
type ParticipantCourseController struct {
	participantCourseRepository ParticipantCourseRepositoryInterface
	tags                        []string
}

// NewParticipantCourseController creates a new ParticipantCourseController instance
func NewParticipantCourseController(repository ParticipantCourseRepositoryInterface) *ParticipantCourseController {
	return &ParticipantCourseController{
		participantCourseRepository: repository,
		tags:                        []string{"ParticipantCourse"},
	}
}

// GetAllParticipantCoursesResponse wraps the response for getting all participant courses
type GetAllParticipantCoursesResponse struct {
	Body []ParticipantCourse
}

// GetAllParticipantCourses retrieves all participant-course mappings from the repository
func (m *ParticipantCourseController) GetAllParticipantCourses(_ context.Context, _ *struct{}) (*GetAllParticipantCoursesResponse, error) {
	results, err := m.participantCourseRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return &GetAllParticipantCoursesResponse{Body: *results}, nil
}

// Register registers all routes related to participant-course mappings
func (m *ParticipantCourseController) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-all-participant-courses",
		Method:      "GET",
		Path:        "/participant-courses",
		Summary:     "Gets all ParticipantID/CourseID maps, along with info for completion rate and access dates.",
		Tags:        m.tags,
	}, m.GetAllParticipantCourses)
}
