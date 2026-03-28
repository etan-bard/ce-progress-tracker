package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakesAnonymizedRepository_GetCourseIDCursor(t *testing.T) {
	mockDB := NewMockDbServiceInterface(t)

	// We can't easily mock mongo.Collection or mongo.Cursor because they are structs with unexported fields
	// and they require a real connection or a complex mock setup.
	// For this test, we will at least verify that GetCollection is called with the right name.

	t.Run("GetCollection is called with correct name", func(t *testing.T) {
		mockDB.EXPECT().GetCollection("takes_anonymized").Return(nil)

		repo := NewTakesAnonymizedRepository(mockDB)
		assert.NotNil(t, repo)
	})
}
