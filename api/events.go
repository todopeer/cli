package api

import (
	"fmt"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/todopeer/cli/util/dt"
)

type Event struct {
	ID          ID
	TaskID      ID `graphql:"taskID"`
	StartAt     Time
	EndAt       *Time
	Description *graphql.String
}

type EventFormatter struct {
	Prefix          string
	WithDate        bool
	DurationFromNow *time.Time
}

func (f EventFormatter) Output(e *Event) {
	fmt.Printf("%s[%d]", f.Prefix, e.ID)
	if f.WithDate {
		fmt.Print(e.StartAt.DateOnly(), " ")
	}
	fmt.Printf("%s - %s", e.StartAt.EventTimeOnly(), e.EndAt.EventTimeOnly())
	if f.DurationFromNow != nil {
		end := (*time.Time)(e.EndAt)
		if end == nil {
			end = f.DurationFromNow
		}
		duration := end.Sub((time.Time)(e.StartAt))
		fmt.Print(" (", dt.FormatDuration(duration, false), ")")
	}
	if e.Description != nil {
		fmt.Println(":", *e.Description)
	} else {
		fmt.Println()
	}
}

type QueryEventsResult struct {
	Events []Event
	Tasks  []Task
}

func (c *Client) GetEvent(eventID ID) (*Event, error) {
	query := struct {
		Event `graphql:"event(id:$id)"`
	}{}
	variables := map[string]interface{}{
		"id": eventID,
	}

	err := c.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Event, nil

}

type QueryRunningEventReuslt struct {
	RunningEvent *Event `graphql:"runningEvent"`
}

func (c *Client) QueryRunningEvent() (event *Event, err error) {
	query := struct {
		QueryRunningEventReuslt `graphql:"me"`
	}{}

	err = c.Query(&query, nil)
	if err != nil {
		return
	}

	return query.RunningEvent, nil
}

func (c *Client) QueryLatestEvents() (*Event, error) {
	query := struct {
		QueryEventsResult `graphql:"events(since:$since, days:3, limit:2)"`
	}{}

	// only care about event in recent 2 days
	since := time.Now().AddDate(0, 0, -2)

	variables := map[string]interface{}{
		"since": since,
	}

	err := c.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	if len(query.Events) == 0 {
		return nil, nil
	}

	return &query.Events[0], nil
}

func (c *Client) QueryEvents(since time.Time, days int) (*QueryEventsResult, error) {
	query := struct {
		QueryEventsResult `graphql:"events(since:$since, days: $days)"`
	}{}
	variables := map[string]interface{}{
		"since": since,
		"days":  graphql.Int(days),
	}

	err := c.Query(&query, variables)
	if err != nil {
		return nil, err
	}

	return &query.QueryEventsResult, nil
}

func (c *Client) DeleteEvent(eventID ID) (*Event, error) {
	var mutation struct {
		Event `graphql:"eventDelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": eventID,
	}

	err := c.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.Event, nil
}

type EventUpdateInput struct {
	Description *graphql.String `json:"description"`
	StartAt     *Time           `json:"startAt"`
	EndAt       *Time           `json:"endAt"`
	TaskID      *ID             `json:"taskID"`
}

func (c *Client) UpdateEvent(eventID ID, input EventUpdateInput) (*Event, error) {
	var mutation struct {
		Event `graphql:"eventUpdate(id:$id, input: $input)"`
	}

	variables := map[string]interface{}{
		"id":    eventID,
		"input": input,
	}

	err := c.Mutate(&mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.Event, nil
}
