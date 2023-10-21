package util

import (
	"time"

	"google.golang.org/genproto/googleapis/type/datetime"
)

func GetCurrentDateTime() *datetime.DateTime {
	now := time.Now()

	dt := &datetime.DateTime{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
		Day:   int32(now.Day()),
		Hours: int32(now.Hour()),
		Minutes: int32(now.Minute()),
		Seconds: int32(now.Second()),
		Nanos:   int32(now.Nanosecond()),
	}

	return dt
}