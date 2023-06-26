package api

import (
	"context"

	"github.com/shurcooL/graphql"
)

// Function to handle the authentication process
func Login(ctx context.Context, email string, password string) (*AuthPayload, error) {
	client := graphql.NewClient(gqlAPI, nil)

	var mutation struct {
		Login AuthPayload `graphql:"login(input: {email: $email, password: $password})"`
	}

	variables := map[string]interface{}{
		"email":    graphql.String(email),
		"password": graphql.String(password),
	}

	err := client.Mutate(ctx, &mutation, variables)
	if err != nil {
		return nil, err
	}

	return &mutation.Login, nil
}

type UserWithTask struct {
	User
	RunningTask  *Task
	RunningEvent *Event
}

func MeWithTaskEvent(ctx context.Context, token string) (user *User, task *Task, event *Event, err error) {
	client := NewClientWithToken(token)

	var query = &struct {
		Me UserWithTask `graphql:"me"`
	}{}

	err = client.Query(ctx, query, nil)
	if err != nil {
		return
	}

	resp := &query.Me
	user = &resp.User
	task = resp.RunningTask
	event = resp.RunningEvent
	return
}

func Me(ctx context.Context, token string) (*User, error) {
	client := NewClientWithToken(token)

	var query = &struct {
		Me User `graphql:"me"`
	}{}

	err := client.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	return &query.Me, nil
}

// Function to handle the deauthentication process
func Logout(ctx context.Context, token string) error {
	client := NewClientWithToken(token)

	var query = &struct {
		Logout bool `graphql:"logout"`
	}{}

	err := client.Mutate(ctx, query, nil)
	if err != nil {
		return err
	}

	return nil
}

type LoginInput struct {
	Email    graphql.String `graphql:"email"`
	Password graphql.String `graphql:"password"`
}

type AuthPayload struct {
	User  User           `graphql:"user"`
	Token graphql.String `graphql:"token"`
}

type User struct {
	ID            graphql.Int    `graphql:"id"`
	Email         graphql.String `graphql:"email"`
	Name          graphql.String `graphql:"name"`
	RunningTaskID *ID            `graphql:"runningTaskID"`
}
