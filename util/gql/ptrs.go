package gql

import "github.com/shurcooL/graphql"

func ToGqlStringP(s string) *graphql.String {
	return (*graphql.String)(&s)
}