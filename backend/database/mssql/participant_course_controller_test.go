package mssql

import (
	"net/http/httptest"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestParticipantCourseController_GetAllParticipantCourses(t *testing.T) {
	mockRepo := &MockParticipantCourseRepositoryInterface{}
	controller := NewParticipantCourseController(mockRepo)

	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("Test", "1.0.0"))
	controller.Register(api)

	t.Run("Returns 200 with content", func(t *testing.T) {
		results := []ParticipantCourse{
			{ParticipantID: 1, CourseID: 2},
		}
		mockRepo.On("GetAll").Return(&results, nil).Once()

		req := httptest.NewRequest("GET", "/participant-courses", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, 200, res.Code)
		assert.Contains(t, res.Body.String(), "1")
	})

	t.Run("Returns 500 on error", func(t *testing.T) {
		mockRepo.On("GetAll").Return(nil, assert.AnError).Once()

		req := httptest.NewRequest("GET", "/participant-courses", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, 500, res.Code)
	})
}
