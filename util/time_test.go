// +build unit
package util

import (
	"testing"
	"time"
)

// TestGetCurrentDateTime tests the GetCurrentDateTime function.
func TestGetCurrentDateTime(t *testing.T) {
	dt := GetCurrentDateTime()

	// Load the Asia/Singapore location
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		t.Fatal(err)
	}

	// Get the current time in Asia/Singapore
	now := time.Now().In(loc)

	// Check the year, month, and day
	if int(dt.Year) != now.Year() || time.Month(dt.Month) != now.Month() || int(dt.Day) != now.Day() {
		t.Errorf("GetCurrentDateTime() date = %v-%v-%v, want %v-%v-%v",
			dt.Year, dt.Month, dt.Day, now.Year(), now.Month(), now.Day())
	}

	// Since hours, minutes, and seconds are very likely to change between the time now and when the function is called
	// Thus, check if they are within a reasonable range (e.g., a few seconds).
	if abs(int32(now.Hour())-dt.Hours) > 1 ||
		abs(int32(now.Minute())-dt.Minutes) > 1 ||
		abs(int32(now.Second())-dt.Seconds) > 1 {
		t.Errorf("GetCurrentDateTime() time = %v:%v:%v, want approximately %v:%v:%v",
			dt.Hours, dt.Minutes, dt.Seconds, now.Hour(), now.Minute(), now.Second())
	}
}

// abs returns the absolute value of an int32.
func abs(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}
