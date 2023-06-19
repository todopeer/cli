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
	Description graphql.String
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
	TaskID      *ID `json:"taskID"`
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
