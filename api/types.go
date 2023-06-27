package api

import (
	"encoding/json"
	"time"
)

type Time time.Time
type ID int64

func (t *Time) EventTimeOnly() string {
	if t == nil {
		return "doing"
	}

	return (*time.Time)(t).Local().Format(time.TimeOnly)
}

func (t *Time) DateOnly() string {
	if t == nil {
		return "-"
	}

	return (*time.Time)(t).Local().Format(time.DateOnly)
}

func (t *Time) DateTime() string {
	if t == nil {
		return ""
	}

	return (*time.Time)(t).Local().Format(time.DateTime)
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	n, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}

	*t = Time(n)
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	var s string
	if t != nil {
		s = (*time.Time)(t).Format(time.RFC3339Nano)
	}
	return json.Marshal(s)
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}

	return (*time.Time)(t).Format(time.DateTime)
}
