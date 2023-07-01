package api

import (
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

func (c *Client) QueryTaskLastEvent(taskID ID) (*Event, error) {
	var query struct {
		Task struct {
			Events []Event `graphql:"events(input:{limit:1})"`
		} `graphql:"task(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := c.client.Query(c.ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if len(query.Task.Events) == 0 {
		return nil, nil
	}

	return &query.Task.Events[0], nil
}

func (c *Client) QueryTasks(input QueryTaskInput) ([]*Task, error) {
	query := struct {
		Tasks []*Task `graphql:"tasks(input:$input)"`
	}{}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.client.Query(c.ctx, &query, variables)
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

func (c *Client) CreateTask(input TaskCreateInput) (*Task, error) {
	var mutation struct {
		TaskCreate Task `graphql:"taskCreate(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.client.Mutate(c.ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskCreate, nil
}

type TaskDeleteInput struct {
	TaskID graphql.String
}

func (c *Client) DeleteTask(taskID ID) (*Task, error) {
	var mutation struct {
		TaskDelete Task `graphql:"taskDelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := c.client.Mutate(c.ctx, &mutation, variables)
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

func (c *Client) StartTask(taskID ID, options ...StartTaskOptionFunc) (*Task, *Event, error) {
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

	err := c.client.Mutate(c.ctx, &mutation, variables)
	if err != nil {
		return nil, nil, err
	}

	resp := mutation.TaskStart
	return &resp.Task, resp.Event, nil
}

func (c *Client) UpdateTask(taskID ID, input TaskUpdateInput) (*Task, error) {
	var mutation struct {
		TaskUpdate Task `graphql:"taskUpdate(id:$id, input: $input)"`
	}

	variables := map[string]interface{}{
		"id":    taskID,
		"input": input,
	}

	err := c.client.Mutate(c.ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskUpdate, nil
}

func (c *Client) UndeleteTask(taskID ID) (*Task, error) {
	var mutation struct {
		TaskUndelete Task `graphql:"taskUndelete(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := c.client.Mutate(c.ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.TaskUndelete, nil
}

func (c *Client) GetTaskEvents(taskID ID) (*Task, []Event, error) {
	var query struct {
		Task struct {
			Task
			Events []Event
		} `graphql:"task(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": taskID,
	}

	err := c.client.Query(c.ctx, &query, variables)
	if err != nil {
		return nil, nil, err
	}

	return &query.Task.Task, query.Task.Events, nil
}
