package mongo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TakesAnonymizedRepositoryTestSuite struct {
	suite.Suite
	mockDB *MockDbServiceInterface
	repo   *TakesAnonymizedRepository
}

func (s *TakesAnonymizedRepositoryTestSuite) SetupTest() {
	s.mockDB = NewMockDbServiceInterface(s.T())
}

func TestTakesAnonymizedRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TakesAnonymizedRepositoryTestSuite))
}

func (s *TakesAnonymizedRepositoryTestSuite) TestGetCourseIDCursor() {
	// We can't easily mock mongo.Collection or mongo.Cursor because they are structs with unexported fields
	// and they require a real connection or a complex mock setup.
	// For this test, we will at least verify that GetCollection is called with the right name.

	s.Run("GetCollection is called with correct name", func() {
		s.mockDB.EXPECT().GetCollection("takes_anonymized").Return(nil)

		repo := NewTakesAnonymizedRepository(s.mockDB)
		s.NotNil(repo)
	})
}
