package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shurcooL/graphql"
)

const gqlAPI = "https://api.todopeer.com/query"

type LogFunc func(string, ...interface{})

var logFunc LogFunc

func SetLogger(logger LogFunc) {
	logFunc = logger
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return http.DefaultTransport.RoundTrip(req)
}

func newClientWithToken(token string) *graphql.Client {
	return graphql.NewClient(gqlAPI, &http.Client{
		Transport: &transport{token: token},
	})
}

type Client struct {
	client *graphql.Client
	ctx    context.Context
}

func NewClient(token string) *Client {
	return &Client{client: newClientWithToken(token), ctx: context.Background()}
}

func (c *Client) Mutate(m any, variables map[string]any) error {
	err := c.client.Mutate(c.ctx, m, variables)

	if logFunc != nil {
		cmdS, _ := json.Marshal(m)
		varS, _ := json.Marshal(variables)

		logFunc("Mutate(%s, %s): %v", string(cmdS), string(varS), err)
	}

	return err
}

func (c *Client) Query(m any, variables map[string]any) error {
	err := c.client.Query(c.ctx, m, variables)

	if logFunc != nil {
		cmdS, _ := json.Marshal(m)
		varS, _ := json.Marshal(variables)

		logFunc("Query(%s, %s): %v", string(cmdS), string(varS), err)
	}

	return err
}
