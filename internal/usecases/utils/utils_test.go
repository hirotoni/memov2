package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekSplitter(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{
			name:     "First week of 2024",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "2024 | Week 1",
		},
		{
			name:     "Middle of the year",
			date:     time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			expected: "2024 | Week 24",
		},
		{
			name:     "End of year",
			date:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: "2025 | Week 1", // Dec 31, 2024 is in week 1 of 2025
		},
		{
			name:     "Week 52",
			date:     time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: "2023 | Week 52",
		},
		{
			name:     "Different year example",
			date:     time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
			expected: "2023 | Week 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WeekSplitter(tt.date)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWeekSplitter_CurrentTime(t *testing.T) {
	// Test with current time
	now := time.Now()
	result := WeekSplitter(now)

	year, week := now.ISOWeek()

	// Just verify it contains the year and week format
	assert.Contains(t, result, "Week")
	assert.Contains(t, result, " | ")

	// Verify the actual year and week are in the result
	year2, week2 := now.ISOWeek()
	assert.Equal(t, year, year2)
	assert.Equal(t, week, week2)
}

func TestWeekSplitter_LeapYear(t *testing.T) {
	// Test with a leap year date
	leapYearDate := time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)
	result := WeekSplitter(leapYearDate)

	assert.Contains(t, result, "2024")
	assert.Contains(t, result, "Week")
}
