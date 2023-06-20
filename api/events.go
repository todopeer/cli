package api

import (
	"context"
	"time"

	"github.com/shurcooL/graphql"
)

type Event struct {
	ID          ID
	TaskID      ID `graphql:"taskID"`
	StartAt     graphql.String
	EndAt       *graphql.String
	Description *graphql.String
}

type QueryEventsResult struct {
	Events []Event
	Tasks  []Task
}

func GetEvent(ctx context.Context, token string, eventID ID) (*Event, error) {
	client := NewClientWithToken(token)

	query := struct {
		Event `graphql:"event(id:$id)"`
	}{}
	variables := map[string]interface{}{
		"id": eventID,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return &query.Event, nil

}

type QueryRunningEventReuslt struct {
	RunningEvent *Event `graphql:"runningEvent"`
}

func QueryRunningEvent(ctx context.Context, token string) (event *Event, err error) {
	client := NewClientWithToken(token)

	query := struct {
		QueryRunningEventReuslt `graphql:"me"`
	}{}

	err = client.Query(ctx, &query, nil)
	if err != nil {
		return
	}

	return query.RunningEvent, nil
}

func QueryLatestEvents(ctx context.Context, token string) (*Event, error) {
	client := NewClientWithToken(token)

	query := struct {
		QueryEventsResult `graphql:"events(since:$since, days:3, limit:2)"`
	}{}

	// only care about event in recent 2 days
	since := time.Now().AddDate(0, 0, -2)

	variables := map[string]interface{}{
		"since": since,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if len(query.Events) == 0 {
		return nil, nil
	}

	return &query.Events[0], nil
}

func QueryEvents(ctx context.Context, token string, since time.Time, days int) (*QueryEventsResult, error) {
	client := NewClientWithToken(token)

	query := struct {
		QueryEventsResult `graphql:"events(since:$since, days: $days)"`
	}{}
	variables := map[string]interface{}{
		"since": since,
		"days":  graphql.Int(days),
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return &query.QueryEventsResult, nil
}

func RemoveEvent(ctx context.Context, token string, eventID ID) (*Event, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		Event `graphql:"eventRemove(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": eventID,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.Event, nil
}

type EventUpdateInput struct {
	Description *graphql.String `json:"description"`
	StartAt     *graphql.String `json:"startAt"`
	EndAt       *graphql.String `json:"endAt"`
	TaskID      *ID             `json:"taskID"`
}

func UpdateEvent(ctx context.Context, token string, eventID ID, input EventUpdateInput) (*Event, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		Event `graphql:"eventUpdate(id:$id, input: $input)"`
	}

	variables := map[string]interface{}{
		"id":    eventID,
		"input": input,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.Event, nil
}
