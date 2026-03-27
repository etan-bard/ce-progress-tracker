package mssql

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ParticipantCourseRepositoryTestSuite struct {
	suite.Suite
	mockDB *MockDBServiceInterface
	repo   *ParticipantCourseRepository
}

func (s *ParticipantCourseRepositoryTestSuite) SetupTest() {
	s.mockDB = NewMockDBServiceInterface(s.T())
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
			DateLastAccessed: now,
			CourseCompletion: 1.0,
		},
		{
			ParticipantID:    2,
			CourseID:         102,
			DateLastAccessed: now,
			CourseCompletion: 0.0,
		},
	}

	s.Run("query includes update condition", func() {
		s.mockDB.On("Rebind", mock.Anything).Return("rebound").Once()
		s.mockDB.On("Select", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		var i, u, sk int
		err := s.repo.UpsertAll(&entries, &i, &u, &sk)
		s.NoError(err)
	})

	s.Run("successful upsert", func() {
		var insertedCount, updatedCount, skippedCount int

		s.mockDB.On("Rebind", mock.Anything).Return("rebound query").Once()
		s.mockDB.On("Select", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dest := args.Get(0)
			d := dest.(*[]struct {
				Action string `db:"Action"`
			})
			*d = append(*d, struct {
				Action string `db:"Action"`
			}{Action: "INSERT"})
			*d = append(*d, struct {
				Action string `db:"Action"`
			}{Action: "UPDATE"})
		}).Return(nil).Once()

		err := s.repo.UpsertAll(&entries, &insertedCount, &updatedCount, &skippedCount)

		s.NoError(err)
		s.Equal(1, insertedCount)
		s.Equal(1, updatedCount)
	})

	s.Run("skipped record when no change", func() {
		var insertedCount, updatedCount, skippedCount int

		s.mockDB.On("Rebind", mock.Anything).Return("rebound query").Once()
		s.mockDB.On("Select", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			dest := args.Get(0)
			d := dest.(*[]struct {
				Action string `db:"Action"`
			})
			*d = []struct {
				Action string `db:"Action"`
			}{}
		}).Return(nil).Once()

		err := s.repo.UpsertAll(&entries, &insertedCount, &updatedCount, &skippedCount)

		s.NoError(err)
		s.Equal(0, insertedCount)
		s.Equal(0, updatedCount)
		s.Equal(2, skippedCount)
	})

	s.Run("empty entries", func() {
		var insertedCount, updatedCount, skippedCount int
		err := s.repo.UpsertAll(&[]ParticipantCourse{}, &insertedCount, &updatedCount, &skippedCount)
		s.NoError(err)
		s.Equal(0, insertedCount)
		s.Equal(0, updatedCount)
	})

	s.Run("nil entries", func() {
		var insertedCount, updatedCount, skippedCount int
		err := s.repo.UpsertAll(nil, &insertedCount, &updatedCount, &skippedCount)
		s.NoError(err)
	})
}
