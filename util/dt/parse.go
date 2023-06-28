package dt

import "time"

func FromTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, timeStr)
}

func FromDate(dateStr string) (time.Time, error) {
	return time.Parse(time.DateOnly, dateStr)
}

func ToDate(t time.Time) string {
	return t.Format(time.DateOnly)
}
