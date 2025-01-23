package util

import (
	"testing"
	"time"
)

func TestCompareDates(t *testing.T) {

	now := time.Now()
	tests := []struct {
		T1       time.Time
		T2       time.Time
		Expected int
	}{
		{
			T1:       now,
			T2:       now.Add(24 * time.Hour),
			Expected: -1,
		},
		{
			T1:       now,
			T2:       now.Add(-24 * time.Hour),
			Expected: 1,
		},
		{
			T1:       now,
			T2:       now,
			Expected: 0,
		},
		{
			T1:       time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			T2:       time.Date(2024, 12, 10, 1, 0, 0, 0, time.UTC),
			Expected: 0,
		},
		{
			T1:       time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			T2:       time.Date(2024, 12, 10, 23, 0, 0, 0, time.UTC),
			Expected: 0,
		},
		{
			T1:       time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
			T2:       time.Date(2024, 12, 10, 24, 0, 0, 0, time.UTC),
			Expected: -1,
		},
		{
			T1:       time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			T2:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			Expected: -1,
		},
	}

	for i, test := range tests {
		result := CompareDates(test.T1, test.T2)
		if result != test.Expected {
			t.Errorf("test %d: expected %d, got %d", i, test.Expected, result)
		}
	}

}
