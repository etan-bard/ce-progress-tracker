package mssql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ParticipantCourseRepositoryTestSuite struct {
	suite.Suite
	mockDB *MockDBServiceInterface
	repo   *ParticipantCourseRepository
}

// SetupSubTest initializes the test suite for each sub-test by creating mock instances of MockDBServiceInterface.
// It then initializes the ParticipantCourseRepository with the mock database service for testing purposes.
func (s *ParticipantCourseRepositoryTestSuite) SetupSubTest() {
	s.mockDB = NewMockDBServiceInterface(s.T())
	s.mockDB.EXPECT().Rebind(mock.Anything).Return("rebound query").Maybe()
	s.repo = NewParticipantCourseRepository(s.mockDB)
}

func TestParticipantCourseRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ParticipantCourseRepositoryTestSuite))
}

func (s *ParticipantCourseRepositoryTestSuite) TestUpsertAll() {
	now := time.Now()
	entries := []ParticipantCourse{
		{
			ParticipantID:    1,
			CourseID:         101,
			DateLastAccessed: &now,
			CourseCompletion: 1.0,
		},
		{
			ParticipantID:    2,
			CourseID:         102,
			DateLastAccessed: &now,
			CourseCompletion: 0.0,
		},
	}

	// Test repository initialization
	s.Run("repository initialized", func() {
		s.NotNil(s.repo)
		s.NotNil(s.repo.db)
	})

	// Table-driven tests for different upsert scenarios
	testCases := []struct {
		name             string
		entries          *[]ParticipantCourse
		setupMock        func()
		expectedInserted int
		expectedUpdated  int
		expectedSkipped  int
		expectError      bool
	}{
		{
			name:    "successful upsert with insert and update",
			entries: &entries,
			setupMock: func() {
				s.mockDB.EXPECT().
					Select(mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Run(func(dest interface{}, query string, args ...interface{}) {
						d := dest.(*[]struct {
							Action string `db:"Action"`
						})
						*d = append(*d, struct {
							Action string `db:"Action"`
						}{Action: "INSERT"})
						*d = append(*d, struct {
							Action string `db:"Action"`
						}{Action: "UPDATE"})
					}).
					Return(nil).Once()
			},
			expectedInserted: 1,
			expectedUpdated:  1,
			expectedSkipped:  0,
			expectError:      false,
		},
		{
			name:    "skipped records when no change",
			entries: &entries,
			setupMock: func() {
				s.mockDB.EXPECT().
					Select(mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Run(func(dest interface{}, query string, args ...interface{}) {
						d := dest.(*[]struct {
							Action string `db:"Action"`
						})
						*d = []struct {
							Action string `db:"Action"`
						}{}
					}).
					Return(nil).Once()
			},
			expectedInserted: 0,
			expectedUpdated:  0,
			expectedSkipped:  2,
			expectError:      false,
		},
		{
			name:             "empty entries",
			entries:          &[]ParticipantCourse{},
			setupMock:        func() {}, // No mock setup needed
			expectedInserted: 0,
			expectedUpdated:  0,
			expectedSkipped:  0,
			expectError:      false,
		},
		{
			name:             "nil entries",
			entries:          nil,
			setupMock:        func() {}, // No mock setup needed
			expectedInserted: 0,
			expectedUpdated:  0,
			expectedSkipped:  0,
			expectError:      false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.setupMock()

			var insertedCount, updatedCount, skippedCount int
			err := s.repo.UpsertAll(tc.entries, &insertedCount, &updatedCount, &skippedCount)

			if tc.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tc.expectedInserted, insertedCount)
				s.Equal(tc.expectedUpdated, updatedCount)
				s.Equal(tc.expectedSkipped, skippedCount)
			}

			s.mockDB.AssertExpectations(s.T())
		})
	}

	// Error scenario tests
	s.Run("database error during upsert", func() {
		s.mockDB.EXPECT().
			Select(mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		var i, u, sk int
		err := s.repo.UpsertAll(&entries, &i, &u, &sk)
		s.mockDB.AssertExpectations(s.T())
		s.Error(err)
		s.Contains(err.Error(), "assert.AnError")
	})

	// Table-driven tests for query verification
	verificationTestCases := []struct {
		name      string
		setupMock func()
	}{
		{
			name: "query includes update condition",
			setupMock: func() {
				s.mockDB.EXPECT().
					Select(mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything,
						mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).Once()
			},
		},
	}

	for _, vtc := range verificationTestCases {
		s.Run(vtc.name, func() {
			vtc.setupMock()

			var i, u, sk int
			err := s.repo.UpsertAll(&entries, &i, &u, &sk)
			s.NoError(err)
			s.mockDB.AssertExpectations(s.T())
		})
	}
}
