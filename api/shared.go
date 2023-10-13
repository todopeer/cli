package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/shurcooL/graphql"
)

const gqlAPI = "https://api.todopeer.com/query"

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
	return c.client.Mutate(c.ctx, m, variables)
}

func (c *Client) Query(m any, variables map[string]any) error {
	return c.client.Query(c.ctx, m, variables)
}
