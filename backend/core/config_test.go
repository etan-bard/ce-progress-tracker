package core

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigGetScriptBatchSize tests the GetScriptBatchSize method with bounds checking
func TestConfigGetScriptBatchSize(t *testing.T) {
	testCases := []struct {
		name        string
		batchSize   int
		expected    int
		description string
	}{
		{
			name:        "batch size within bounds",
			batchSize:   100,
			expected:    100,
			description: "Batch size of 100 should be returned as-is",
		},
		{
			name:        "batch size at lower bound",
			batchSize:   1,
			expected:    1,
			description: "Batch size of 1 should be returned as-is",
		},
		{
			name:        "batch size at upper bound",
			batchSize:   400,
			expected:    400,
			description: "Batch size of 400 should be returned as-is",
		},
		{
			name:        "batch size too small",
			batchSize:   0,
			expected:    400,
			description: "Batch size of 0 should default to 400",
		},
		{
			name:        "batch size too large",
			batchSize:   401,
			expected:    400,
			description: "Batch size of 401 should default to 400",
		},
		{
			name:        "negative batch size",
			batchSize:   -1,
			expected:    400,
			description: "Negative batch size should default to 400",
		},
		{
			name:        "very large batch size",
			batchSize:   1000,
			expected:    400,
			description: "Batch size of 1000 should default to 400",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test config
			v := viper.New()
			v.Set("SCRIPT_BATCH_SIZE", tc.batchSize)

			config := &Config{
				Viper: v,
			}

			result := config.GetScriptBatchSize()
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

// TestConfigGetMaxGoroutines tests the GetMaxGoroutines method with bounds checking
func TestConfigGetMaxGoroutines(t *testing.T) {
	testCases := []struct {
		name          string
		maxGoroutines int
		expected      int
		description   string
	}{
		{
			name:          "max goroutines within bounds",
			maxGoroutines: 5,
			expected:      5,
			description:   "Max goroutines of 5 should be returned as-is",
		},
		{
			name:          "max goroutines at lower bound",
			maxGoroutines: 1,
			expected:      1,
			description:   "Max goroutines of 1 should be returned as-is",
		},
		{
			name:          "max goroutines too small",
			maxGoroutines: 0,
			expected:      10,
			description:   "Max goroutines of 0 should default to 10",
		},
		{
			name:          "negative max goroutines",
			maxGoroutines: -1,
			expected:      10,
			description:   "Negative max goroutines should default to 10",
		},
		{
			name:          "large max goroutines",
			maxGoroutines: 100,
			expected:      100,
			description:   "Max goroutines of 100 should be returned as-is",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test config
			v := viper.New()
			v.Set("MAX_GOROUTINES", tc.maxGoroutines)

			config := &Config{
				Viper: v,
			}

			result := config.GetMaxGoroutines()
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

// TestConfigGetCourseIDs tests the GetCourseIDs method
func TestConfigGetCourseIDs(t *testing.T) {
	t.Run("course IDs from viper", func(t *testing.T) {
		// Set up test with viper directly
		v := viper.New()
		v.Set("COURSE_IDS", []int{1, 2, 3, 4, 5})

		config := &Config{
			Viper: v,
		}

		result := config.GetCourseIDs()
		require.NotNil(t, result)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, *result)
	})

	t.Run("nil course IDs", func(t *testing.T) {
		// Test when no course IDs are set
		v := viper.New()

		config := &Config{
			Viper: v,
		}

		result := config.GetCourseIDs()
		assert.Nil(t, result)
	})

	t.Run("empty course IDs", func(t *testing.T) {
		// Test with empty slice
		v := viper.New()
		v.Set("COURSE_IDS", []int{})

		config := &Config{
			Viper: v,
		}

		result := config.GetCourseIDs()
		require.NotNil(t, result)
		assert.Empty(t, *result)
	})
}

// TestConfigGetAPIPort tests the GetAPIPort method
func TestConfigGetAPIPort(t *testing.T) {
	testCases := []struct {
		name        string
		port        string
		expected    string
		description string
	}{
		{
			name:        "custom port",
			port:        "3000",
			expected:    "3000",
			description: "Custom port should be returned",
		},
		{
			name:        "empty port (should use default)",
			port:        "",
			expected:    "8080",
			description: "Empty port should default to 8080",
		},
		{
			name:        "standard port",
			port:        "8080",
			expected:    "8080",
			description: "Standard port 8080 should be returned",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := viper.New()
			if tc.port != "" {
				v.Set("API_PORT", tc.port)
			}

			config := &Config{
				Viper: v,
			}

			result := config.GetAPIPort()
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
