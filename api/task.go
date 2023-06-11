package api

import (
	"context"
	"time"

	"github.com/shurcooL/graphql"
)

type TaskStatus string

var (
	TaskStatusNotStarted TaskStatus = "NOT_STARTED"
	TaskStatusDoing      TaskStatus = "DOING"
	TaskStatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID          graphql.Int
	Name        graphql.String
	Description graphql.String
	Status      TaskStatus
	CreatedAt   graphql.String
	UpdatedAt   graphql.String
	DueDate     *graphql.String
}

type TaskUpdateInput struct {
	TaskID      graphql.Int     `json:"taskId"`
	Name        *graphql.String `json:"name"`
	Description *graphql.String `json:"description"`
	Status      *TaskStatus     `json:"status"`
	DueDate     *graphql.String `json:"dueDate"`
}

type QueryTaskInput struct {
	Status *TaskStatus `graphql:"status"`
}

func QueryTasks(ctx context.Context, token string, input *QueryTaskInput) ([]*Task, error) {
	client := NewClientWithToken(token)

	query := struct {
		Tasks []*Task `graphql:"tasks(input:{})"`
	}{}

	err := client.Query(ctx, &query, nil)
	if err != nil {
		return nil, err
	}

	return query.Tasks, nil
}

type Time time.Time

func (t *Time) String() string {
	if t == nil {
		return ""
	}

	return (*time.Time)(t).Format(time.DateTime)
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

type TaskRemoveInput struct {
	TaskID graphql.String
}

type ID int64

func RemoveTask(ctx context.Context, token string, taskID ID) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskRemove Task `graphql:"taskRemove(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskRemove, nil
}

func StartTask(ctx context.Context, token string, taskID ID) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskStart Task `graphql:"taskStart(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskStart, nil
}

func UpdateTask(ctx context.Context, token string, input TaskUpdateInput) (*Task, error) {
	client := NewClientWithToken(token)

	var mutation struct {
		TaskUpdate Task `graphql:"taskUpdate(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskUpdate, nil
}
