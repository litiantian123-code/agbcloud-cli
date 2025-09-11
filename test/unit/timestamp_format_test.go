// Copyright 2025 AgbCloud CLI Contributors
// SPDX-License-Identifier: Apache-2.0

package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// formatTimestamp is the function we're testing (copied from cmd/image.go)
func formatTimestamp(timestamp string) string {
	if timestamp == "" {
		return "-"
	}
	// Try to parse the UTC timestamp and convert to local timezone
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		// Convert UTC time to local timezone
		localTime := t.Local()
		return localTime.Format("2006-01-02 15:04")
	}
	// If parsing fails, return the original string truncated
	return truncateString(timestamp, 20)
}

// truncateString helper function
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectEmpty bool
		checkFormat bool
	}{
		{
			name:        "Empty timestamp",
			input:       "",
			expectEmpty: true,
		},
		{
			name:        "Valid UTC timestamp",
			input:       "2025-09-11T05:48:08Z",
			checkFormat: true,
		},
		{
			name:        "Valid UTC timestamp with milliseconds",
			input:       "2025-09-11T05:48:08.123Z",
			checkFormat: true,
		},
		{
			name:        "Valid UTC timestamp with timezone offset",
			input:       "2025-09-11T05:48:08+00:00",
			checkFormat: true,
		},
		{
			name:  "Invalid timestamp format",
			input: "invalid-timestamp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestamp(tt.input)

			if tt.expectEmpty {
				assert.Equal(t, "-", result, "Empty timestamp should return '-'")
				return
			}

			if tt.checkFormat {
				// For valid timestamps, check that the result follows the expected format
				assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`, result, "Result should match YYYY-MM-DD HH:MM format")

				// Verify that timezone conversion actually happened
				// Parse the original UTC time
				utcTime, err := time.Parse(time.RFC3339, tt.input)
				assert.NoError(t, err, "Should be able to parse input timestamp")

				// Convert to local time
				localTime := utcTime.Local()
				expectedResult := localTime.Format("2006-01-02 15:04")

				assert.Equal(t, expectedResult, result, "Should correctly convert UTC to local time")
				return
			}

			// For invalid timestamps, should return truncated original string
			if len(tt.input) > 20 {
				assert.Equal(t, tt.input[:17]+"...", result, "Invalid long timestamp should be truncated")
			} else {
				assert.Equal(t, tt.input, result, "Invalid short timestamp should be returned as-is")
			}
		})
	}
}

func TestTimezoneConversion(t *testing.T) {
	// Test with a known UTC timestamp
	utcTimestamp := "2025-09-11T05:48:08Z"

	// Parse the UTC time
	utcTime, err := time.Parse(time.RFC3339, utcTimestamp)
	assert.NoError(t, err, "Should parse UTC timestamp")

	// Get the expected local time
	localTime := utcTime.Local()
	expectedFormat := localTime.Format("2006-01-02 15:04")

	// Test our function
	result := formatTimestamp(utcTimestamp)

	// Should match the expected local time format
	assert.Equal(t, expectedFormat, result, "Should convert UTC to local timezone correctly")

	// Verify that conversion actually happened (unless we're in UTC timezone)
	if localTime.Location() != time.UTC {
		utcFormat := utcTime.Format("2006-01-02 15:04")
		assert.NotEqual(t, utcFormat, result, "Local time should be different from UTC time (unless in UTC timezone)")
	}
}
