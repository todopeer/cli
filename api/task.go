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
	CreatedAt   graphql.String
	UpdatedAt   graphql.String
	DueDate     *graphql.String
}

func (t *Task) Output() {
	fmt.Printf("%d\t%s\t%s\t", t.ID, t.Status, t.Name)
	if t.DueDate != nil {
		fmt.Printf("%s\n", *t.DueDate)
	} else {
		fmt.Println()
	}
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
		TaskStart Task `graphql:"taskUpdate(id: $id, input: {status: DOING})"`
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
