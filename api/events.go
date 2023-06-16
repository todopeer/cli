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
