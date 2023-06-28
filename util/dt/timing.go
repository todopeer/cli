package dt

import (
	"strconv"
	"time"
)

func FromTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, timeStr)
}

func FromDate(dateStr string) (time.Time, error) {
	return time.Parse(time.DateOnly, dateStr)
}

func ToDate(t time.Time) string {
	return t.Format(time.DateOnly)
}

func FormatDuration(d time.Duration, withSec bool) string {
	res := ""
	if d >= time.Hour {
		res = strconv.Itoa(int(d/time.Hour)) + "h"
		d %= time.Hour
	}
	if d >= time.Minute {
		res = res + strconv.Itoa(int(d/time.Minute)) + "m"
		d %= time.Minute
	}
	if withSec && d >= time.Second {
		res = res + strconv.Itoa(int(d/time.Second)) + "s"
	}
	return res
}
