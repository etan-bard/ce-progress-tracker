package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TakesAnonymizedRepositoryTestSuite struct {
	suite.Suite
	mockDB         *MockDbServiceInterface
	mockCollection *MockCollectionInterface
	repo           *TakesAnonymizedRepository
}

// SetupTest initializes the test suite by creating mock instances of DbServiceInterface and CollectionInterface.
// It then initializes the TakesAnonymizedRepository with the mock collection for testing purposes.
func (s *TakesAnonymizedRepositoryTestSuite) SetupTest() {
	s.mockDB = NewMockDbServiceInterface(s.T())
	s.mockCollection = NewMockCollectionInterface(s.T())
	s.mockDB.EXPECT().GetCollection(mock.Anything).Return(s.mockCollection)

	s.repo = NewTakesAnonymizedRepository(s.mockDB)
}

func TestTakesAnonymizedRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TakesAnonymizedRepositoryTestSuite))
}

func (s *TakesAnonymizedRepositoryTestSuite) TestGetCourseIDCursor() {
	ctx := context.Background()
	batchSize := 10

	// Test repository initialization
	s.Run("repository initialized", func() {
		s.NotNil(s.repo)
		s.NotNil(s.repo.collection)
	})

	// Table-driven tests for different course ID scenarios
	testCases := []struct {
		name        string
		courseIDs   *[]int
		setupFilter func(map[string]interface{})
	}{
		{
			name:      "cursor with course IDs",
			courseIDs: &[]int{1, 2, 3},
			setupFilter: func(filter map[string]interface{}) {
				s.Contains(filter, "course_data.course_id")
				s.Equal([]int{1, 2, 3}, filter["course_data.course_id"].(map[string]interface{})["$in"])
			},
		},
		{
			name:      "cursor with nil IDs",
			courseIDs: nil,
			setupFilter: func(filter map[string]interface{}) {
				s.NotContains(filter, "course_data.course_id")
			},
		},
		{
			name:      "cursor with empty IDs",
			courseIDs: &[]int{},
			setupFilter: func(filter map[string]interface{}) {
				s.NotContains(filter, "course_data.course_id")
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockCursor := NewMockCursorInterface(s.T())

			s.mockCollection.EXPECT().
				Find(ctx, mock.Anything, mock.Anything).
				Run(func(ctx context.Context, f any, opts ...*options.FindOptions) {
					filter := f.(map[string]interface{})
					tc.setupFilter(filter)

					// Common verification for all cases
					s.Contains(filter, "participant_data.date_first_accessed")
					s.Equal(float64(0), filter["participant_data.date_first_accessed"].(map[string]interface{})["$gt"])
				}).
				Return(mockCursor, nil).Once()

			cursor, err := s.repo.GetCourseIDCursor(ctx, tc.courseIDs, batchSize)
			s.mockCollection.AssertExpectations(s.T())
			s.NoError(err)
			s.NotNil(cursor)
		})
	}

	// Table-driven tests for batch size scenarios
	batchTestCases := []struct {
		name      string
		batchSize int
	}{
		{
			name:      "batch size 10",
			batchSize: 10,
		},
		{
			// valid batch size (uses mongo default)
			name:      "batch size 0",
			batchSize: 0,
		},
		{
			name:      "batch size 500",
			batchSize: 500,
		},
	}

	for _, bt := range batchTestCases {
		s.Run(bt.name, func() {
			mockCursor := NewMockCursorInterface(s.T())
			courseIDs := []int{1, 2, 3}

			s.mockCollection.EXPECT().
				Find(ctx, mock.Anything, mock.Anything).
				Run(func(ctx context.Context, f any, opts ...*options.FindOptions) {
					s.NotNil(opts)
					s.Equal(1, len(opts))

					findOptions := opts[0]
					s.NotNil(findOptions)
					s.Equal(int32(bt.batchSize), *findOptions.BatchSize)
				}).
				Return(mockCursor, nil).Once()

			cursor, err := s.repo.GetCourseIDCursor(ctx, &courseIDs, bt.batchSize)
			s.mockCollection.AssertExpectations(s.T())
			s.NoError(err)
			s.NotNil(cursor)
		})
	}

	// Error scenario tests
	s.Run("database error", func() {
		courseIDs := []int{1, 2, 3}
		batchSize := 10

		s.mockCollection.EXPECT().
			Find(ctx, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		cursor, err := s.repo.GetCourseIDCursor(ctx, &courseIDs, batchSize)
		s.mockCollection.AssertExpectations(s.T())
		s.Error(err)
		s.Nil(cursor)
		s.Contains(err.Error(), "assert.AnError")
	})

	s.Run("invalid batch size", func() {
		courseIDs := []int{1, 2, 3}
		batchSize := -1

		s.mockCollection.EXPECT().
			Find(ctx, mock.Anything, mock.Anything).
			Return(nil, assert.AnError).Once()

		cursor, err := s.repo.GetCourseIDCursor(ctx, &courseIDs, batchSize)
		s.mockCollection.AssertExpectations(s.T())
		s.Error(err)
		s.Nil(cursor)
	})
}
