package api

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/graphql"
)

type TaskStatus string

var (
	TaskStatusNotStarted TaskStatus = "NOT_STARTED"
	TaskStatusDoing      TaskStatus = "DOING"
	TaskStatusDone       TaskStatus = "DONE"
	TaskStatusPaused     TaskStatus = "PAUSED"
)

type Task struct {
	ID          ID
	Name        graphql.String
	Description graphql.String
	Status      TaskStatus
	CreatedAt   Time
	UpdatedAt   Time
	DueDate     *Time
}

func (t *Task) Output() {
	fmt.Printf("%d\t%s\t%s\t%s\n", t.ID, t.Status, t.Name, t.DueDate.DateTime())
}

type TaskUpdateInput struct {
	Name        *graphql.String `json:"name"`
	Description *graphql.String `json:"description"`
	Status      *TaskStatus     `json:"status"`
	DueDate     *graphql.String `json:"dueDate"`
}

type QueryTaskInput struct {
	Status []TaskStatus `json:"status"`
}

func QueryTaskLastEvent(ctx context.Context, token string, taskID ID) (*Event, error) {
	client := NewClientWithToken(token)

	var query struct {
		Task struct {
			Events []Event `graphql:"events(input:{limit:1})"`
		} `graphql:"task(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if len(query.Task.Events) == 0 {
		return nil, nil
	}

	return &query.Task.Events[0], nil
}

func QueryTasks(ctx context.Context, token string, input QueryTaskInput) ([]*Task, error) {
	client := NewClientWithToken(token)

	query := struct {
		Tasks []*Task `graphql:"tasks(input:$input)"`
	}{}
	variables := map[string]interface{}{
		"input": input,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	return query.Tasks, nil
}

type TaskCreateInput struct {
	Name        graphql.String  `json:"name"`
	Description *graphql.String `json:"description"`
	DueDate     *graphql.String `json:"dueDate"`
}

func CreateTask(ctx context.Context, token string, input TaskCreateInput) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskCreate Task `graphql:"taskCreate(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskCreate, nil
}

type TaskDeleteInput struct {
	TaskID graphql.String
}

func DeleteTask(ctx context.Context, token string, taskID ID) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskDelete Task `graphql:"taskDelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskDelete, nil
}

type startTaskOption struct {
	offset time.Duration
	desc   *string
}

type StartTaskOptionFunc func(*startTaskOption)

func StartTaskWithOffset(duration time.Duration) StartTaskOptionFunc {
	return func(o *startTaskOption) {
		o.offset = duration
	}
}

func StartTaskWithDescription(description string) StartTaskOptionFunc {
	return func(o *startTaskOption) {
		o.desc = &description
	}
}

func StartTask(ctx context.Context, token string, taskID ID, options ...StartTaskOptionFunc) (*Task, *Event, error) {
	client := NewClientWithToken(token)

	cfg := startTaskOption{}
	for _, option := range options {
		option(&cfg)
	}

	var mutation struct {
		TaskStart struct {
			Task  Task
			Event *Event
		} `graphql:"taskStart(id: $id, input: {startAt: $startAt, description: $description})"`
	}

	variables := map[string]interface{}{
		"id":          taskID,
		"description": (*graphql.String)(cfg.desc),
		"startAt":     time.Now().Add(-cfg.offset),
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, nil, err
	}

	resp := mutation.TaskStart
	return &resp.Task, resp.Event, nil
}

func UpdateTask(ctx context.Context, token string, taskID ID, input TaskUpdateInput) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskUpdate Task `graphql:"taskUpdate(id:$id, input: $input)"`
	}

	variables := map[string]interface{}{
		"id":    taskID,
		"input": input,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskUpdate, nil
}

func UndeleteTask(ctx context.Context, token string, taskID ID) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskUndelete Task `graphql:"taskUndelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskUndelete, nil
}

func GetTaskEvents(ctx context.Context, token string, taskID ID) (*Task, []Event, error) {
	client := NewClientWithToken(token)

	var query struct {
		Task struct {
			Task
			Events []Event
		} `graphql:"task(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		return nil, nil, err
	}

	return &query.Task.Task, query.Task.Events, nil
}
