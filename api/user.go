package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shurcooL/graphql"
)

const gqlAPI = "http://localhost:8080/query"

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

func NewClientWithToken(token string) *graphql.Client {
	return graphql.NewClient(gqlAPI, &http.Client{
		Transport: &transport{token: token},
	})
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

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return http.DefaultTransport.RoundTrip(req)
}

// Function to handle the deauthentication process
func Deauthenticate() error {
	// Code to delete stored token and user info

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
	ID    graphql.Int    `graphql:"id"`
	Email graphql.String `graphql:"email"`
	Name  graphql.String `graphql:"name"`
}
