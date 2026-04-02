package services

import (
	"context"
	"testing"

	mongo "ce-progress-tracker/database/mongo"
	"ce-progress-tracker/database/mssql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type DataMigrationStrategyTestSuite struct {
	suite.Suite
	mockMapper     *MockParticipantMapperInterface
	mockRepository *mssql.MockParticipantCourseRepositoryInterface
	mockMongoRepo  *mongo.MockTakesAnonymizedRepositoryInterface
	strategy       *BatchDataMigrationStrategy
	logger         *zap.Logger
}

func (s *DataMigrationStrategyTestSuite) SetupTest() {
	s.mockMapper = NewMockParticipantMapperInterface(s.T())
	s.mockRepository = mssql.NewMockParticipantCourseRepositoryInterface(s.T())
	s.mockMongoRepo = mongo.NewMockTakesAnonymizedRepositoryInterface(s.T())

	// Create a test logger
	var err error
	s.logger, err = zap.NewDevelopment()
	s.Require().NoError(err)

	s.strategy = NewBatchDataMigrationStrategy(
		s.mockMapper,
		s.mockRepository,
		s.mockMongoRepo,
		s.logger,
		10, // batchSize
		2,  // maxGoroutines
	)
}

func TestDataMigrationStrategyTestSuite(t *testing.T) {
	suite.Run(t, new(DataMigrationStrategyTestSuite))
}

func (s *DataMigrationStrategyTestSuite) TestExecute() {
	ctx := context.Background()
	courseIDs := []int{101, 102}

	testCases := []struct {
		name             string
		setupMocks       func()
		expectedInserted int
		expectedUpdated  int
		expectedSkipped  int
		expectError      bool
		errorContains    string
	}{
		{
			name: "successful migration with all operations",
			setupMocks: func() {
				// Setup MongoDB cursor mock
				mockCursor := mongo.NewMockCursorInterface(s.T())

				// Create test data
				testRecords := []mongo.TakesAnonymized{
					{
						ParticipantData: &mongo.ParticipantData{
							ParticipantID:    1,
							CourseCompletion: 1.0,
						},
						CourseData: &mongo.CourseData{CourseID: 101},
					},
					{
						ParticipantData: &mongo.ParticipantData{
							ParticipantID:    2,
							CourseCompletion: 0.5,
						},
						CourseData: &mongo.CourseData{CourseID: 102},
					},
				}

				// Setup cursor expectations
				mockCursor.EXPECT().Next(mock.Anything).Return(true).Once()
				mockCursor.EXPECT().Decode(mock.Anything).RunAndReturn(func(dest interface{}) error {
					if record, ok := dest.(*mongo.TakesAnonymized); ok {
						*record = testRecords[0]
					}
					return nil
				}).Once()

				mockCursor.EXPECT().Next(mock.Anything).Return(true).Once()
				mockCursor.EXPECT().Decode(mock.Anything).RunAndReturn(func(dest interface{}) error {
					if record, ok := dest.(*mongo.TakesAnonymized); ok {
						*record = testRecords[1]
					}
					return nil
				}).Once()

				mockCursor.EXPECT().Next(mock.Anything).Return(false).Once()
				mockCursor.EXPECT().Close(mock.Anything).Return(nil).Once()

				// Setup MongoDB repository mock
				s.mockMongoRepo.EXPECT().GetCourseIDCursor(ctx, &courseIDs, 10).Return(mockCursor, nil).Once()

				// Setup mapper mocks
				s.mockMapper.EXPECT().MongoToSQL(&testRecords[0]).Return(&mssql.ParticipantCourse{
					ParticipantID:    1,
					CourseID:         101,
					CourseCompletion: 1.0,
				}).Once()

				s.mockMapper.EXPECT().MongoToSQL(&testRecords[1]).Return(&mssql.ParticipantCourse{
					ParticipantID:    2,
					CourseID:         102,
					CourseCompletion: 0.5,
				}).Once()

				// Setup repository mocks for batch processing
				batch1 := []mssql.ParticipantCourse{
					{ParticipantID: 1, CourseID: 101, CourseCompletion: 1.0},
					{ParticipantID: 2, CourseID: 102, CourseCompletion: 0.5},
				}

				s.mockRepository.EXPECT().UpsertAll(&batch1, mock.Anything, mock.Anything, mock.Anything).
					Run(func(batch *[]mssql.ParticipantCourse, inserted, updated, skipped *int) {
						*inserted = 1
						*updated = 1
						*skipped = 0
					}).
					Return(nil).Once()
			},
			expectedInserted: 1,
			expectedUpdated:  1,
			expectedSkipped:  0,
			expectError:      false,
		},
		{
			name: "migration with skipped records",
			setupMocks: func() {
				// Setup MongoDB cursor mock
				mockCursor := mongo.NewMockCursorInterface(s.T())

				// Create test data - one valid, one invalid
				testRecords := []mongo.TakesAnonymized{
					{
						ParticipantData: &mongo.ParticipantData{
							ParticipantID:    1,
							CourseCompletion: 1.0,
						},
						CourseData: &mongo.CourseData{CourseID: 101},
					},
					{}, // Invalid record (missing participant/course data)
				}

				// Setup cursor expectations
				mockCursor.EXPECT().Next(mock.Anything).Return(true).Once()
				mockCursor.EXPECT().Decode(mock.Anything).RunAndReturn(func(dest interface{}) error {
					if record, ok := dest.(*mongo.TakesAnonymized); ok {
						*record = testRecords[0]
					}
					return nil
				}).Once()

				mockCursor.EXPECT().Next(mock.Anything).Return(true).Once()
				mockCursor.EXPECT().Decode(mock.Anything).RunAndReturn(func(dest interface{}) error {
					if record, ok := dest.(*mongo.TakesAnonymized); ok {
						*record = testRecords[1]
					}
					return nil
				}).Once()

				mockCursor.EXPECT().Next(mock.Anything).Return(false).Once()
				mockCursor.EXPECT().Close(mock.Anything).Return(nil).Once()

				// Setup MongoDB repository mock
				s.mockMongoRepo.EXPECT().GetCourseIDCursor(ctx, &courseIDs, 10).Return(mockCursor, nil).Once()

				// Setup mapper mocks - first record maps successfully, second returns nil
				s.mockMapper.EXPECT().MongoToSQL(&testRecords[0]).Return(&mssql.ParticipantCourse{
					ParticipantID:    1,
					CourseID:         101,
					CourseCompletion: 1.0,
				}).Once()

				s.mockMapper.EXPECT().MongoToSQL(&testRecords[1]).Return(nil).Once()

				// Setup repository mocks for batch processing
				batch1 := []mssql.ParticipantCourse{
					{ParticipantID: 1, CourseID: 101, CourseCompletion: 1.0},
				}

				s.mockRepository.EXPECT().UpsertAll(&batch1, mock.Anything, mock.Anything, mock.Anything).
					Run(func(batch *[]mssql.ParticipantCourse, inserted, updated, skipped *int) {
						*inserted = 0
						*updated = 1
						*skipped = 0
					}).
					Return(nil).Once()
			},
			expectedInserted: 0,
			expectedUpdated:  1,
			expectedSkipped:  0,
			expectError:      false,
		},
		{
			name: "migration with empty course IDs",
			setupMocks: func() {
				mockCursor := mongo.NewMockCursorInterface(s.T())

				// Empty cursor - no records to process
				mockCursor.EXPECT().Next(mock.Anything).Return(false).Once()
				mockCursor.EXPECT().Close(mock.Anything).Return(nil).Once()

				s.mockMongoRepo.EXPECT().GetCourseIDCursor(ctx, &courseIDs, 10).Return(mockCursor, nil).Once()
			},
			expectedInserted: 0,
			expectedUpdated:  0,
			expectedSkipped:  0,
			expectError:      false,
		},
		{
			name: "migration with database error",
			setupMocks: func() {
				mockCursor := mongo.NewMockCursorInterface(s.T())

				testRecords := []mongo.TakesAnonymized{
					{
						ParticipantData: &mongo.ParticipantData{
							ParticipantID:    1,
							CourseCompletion: 1.0,
						},
						CourseData: &mongo.CourseData{CourseID: 101},
					},
				}

				mockCursor.EXPECT().Next(mock.Anything).Return(true).Once()
				mockCursor.EXPECT().Decode(mock.Anything).RunAndReturn(func(dest interface{}) error {
					if record, ok := dest.(*mongo.TakesAnonymized); ok {
						*record = testRecords[0]
					}
					return nil
				}).Once()

				mockCursor.EXPECT().Next(mock.Anything).Return(false).Once()
				mockCursor.EXPECT().Close(mock.Anything).Return(nil).Once()

				s.mockMongoRepo.EXPECT().GetCourseIDCursor(ctx, &courseIDs, 10).Return(mockCursor, nil).Once()
				s.mockMapper.EXPECT().MongoToSQL(&testRecords[0]).Return(&mssql.ParticipantCourse{
					ParticipantID:    1,
					CourseID:         101,
					CourseCompletion: 1.0,
				}).Once()

				batch1 := []mssql.ParticipantCourse{
					{ParticipantID: 1, CourseID: 101, CourseCompletion: 1.0},
				}

				s.mockRepository.EXPECT().UpsertAll(&batch1, mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError).Once()
			},
			expectedInserted: 0,
			expectedUpdated:  0,
			expectedSkipped:  0,
			expectError:      true,
			errorContains:    "assert.AnError",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			tc.setupMocks()

			inserted, updated, skipped, err := s.strategy.Execute(ctx, &courseIDs)

			if tc.expectError {
				s.Error(err)
				if tc.errorContains != "" {
					s.Contains(err.Error(), tc.errorContains)
				}
			} else {
				s.NoError(err)
				s.Equal(tc.expectedInserted, inserted)
				s.Equal(tc.expectedUpdated, updated)
				s.Equal(tc.expectedSkipped, skipped)
			}

			s.mockMapper.AssertExpectations(s.T())
			s.mockRepository.AssertExpectations(s.T())
			s.mockMongoRepo.AssertExpectations(s.T())
		})
	}
}
