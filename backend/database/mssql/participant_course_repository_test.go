package mssql

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParticipantCourseRepository_UpsertAll(t *testing.T) {
	mockDB := NewMockDBServiceInterface(t)
	repo := NewParticipantCourseRepository(mockDB)

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

	t.Run("query includes update condition", func(t *testing.T) {
		mockDB.EXPECT().Rebind(mock.MatchedBy(func(query string) bool {
			return strings.Contains(query, "WHEN MATCHED AND (Target.DateLastAccessed <> Source.DateLastAccessed OR Target.CourseCompletion <> Source.CourseCompletion)")
		})).Return("rebound")
		mockDB.EXPECT().Select(mock.Anything, "rebound", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := repo.UpsertAll(&entries, nil, nil, nil)
		assert.NoError(t, err)
	})

	t.Run("successful upsert", func(t *testing.T) {
		var insertedCount, updatedCount int

		mockDB.EXPECT().Rebind(mock.AnythingOfType("string")).Return("rebound query")
		mockDB.EXPECT().Select(mock.Anything, "rebound query", mock.Anything).Run(func(dest interface{}, query string, args ...interface{}) {
			// Simulate SQL Server returning the actions
			d := dest.(*[]struct {
				Action string `db:"Action"`
			})
			*d = append(*d, struct {
				Action string `db:"Action"`
			}{Action: "INSERT"})
			*d = append(*d, struct {
				Action string `db:"Action"`
			}{Action: "UPDATE"})
		}).Return(nil)

		err := repo.UpsertAll(&entries, &insertedCount, &updatedCount)

		assert.NoError(t, err)
		assert.Equal(t, 1, insertedCount)
		assert.Equal(t, 1, updatedCount)
	})

	t.Run("skipped record when no change", func(t *testing.T) {
		var insertedCount, updatedCount, skippedCount int

		mockDB.EXPECT().Rebind(mock.AnythingOfType("string")).Return("rebound query")
		mockDB.EXPECT().Select(mock.Anything, "rebound query", mock.Anything).Run(func(dest interface{}, query string, args ...interface{}) {
			// Simulate SQL Server returning nothing for a record that matched but had no changes
			d := dest.(*[]struct {
				Action string `db:"Action"`
			})
			*d = []struct {
				Action string `db:"Action" `
			}{} // Empty result
		}).Return(nil)

		err := repo.UpsertAll(&entries, &insertedCount, &updatedCount, &skippedCount)

		assert.NoError(t, err)
		assert.Equal(t, 0, insertedCount)
		assert.Equal(t, 0, updatedCount)
		assert.Equal(t, 2, skippedCount)
	})

	t.Run("argument count check", func(t *testing.T) {
		// This is mostly to ensure we haven't introduced a logic error in building valueArgs
		// although the test above already covers it implicitly.
		var insertedCount, updatedCount int
		mockDB.EXPECT().Rebind(mock.Anything).Return("rebound")
		mockDB.EXPECT().Select(mock.Anything, "rebound", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		err := repo.UpsertAll(&entries, &insertedCount, &updatedCount)
		assert.NoError(t, err)
	})

	t.Run("empty entries", func(t *testing.T) {
		var insertedCount, updatedCount int
		err := repo.UpsertAll(&[]ParticipantCourse{}, &insertedCount, &updatedCount)
		assert.NoError(t, err)
		assert.Equal(t, 0, insertedCount)
		assert.Equal(t, 0, updatedCount)
	})

	t.Run("nil entries", func(t *testing.T) {
		var insertedCount, updatedCount int
		err := repo.UpsertAll(nil, &insertedCount, &updatedCount)
		assert.NoError(t, err)
	})
}
