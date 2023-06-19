package dt

import "time"

func FromTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, timeStr)
}

func FromTimePtr(timeStrP *string) (*time.Time, error) {
	if timeStrP == nil {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339Nano, *timeStrP)
	return &t, err
}

func FromDate(dateStr string) (time.Time, error) {
	return time.Parse(time.DateOnly, dateStr)
}

func ToTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func ToDate(t time.Time) string {
	return t.Format(time.DateOnly)
}
