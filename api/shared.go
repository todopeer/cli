package api

import (
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

func NewClientWithToken(token string) *graphql.Client {
	return graphql.NewClient(gqlAPI, &http.Client{
		Transport: &transport{token: token},
	})
}
