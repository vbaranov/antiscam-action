package antiscam

import (
	"context"

	"github.com/google/go-github/v69/github"
	"github.com/shurcooL/githubv4"
)

type Antiscam struct {
	ctx           context.Context
	restClient    *github.Client
	graphqlClient *githubv4.Client
}

func New(ctx context.Context, restClient *github.Client, graphqlClient *githubv4.Client) *Antiscam {
	return &Antiscam{
		ctx:           ctx,
		restClient:    restClient,
		graphqlClient: graphqlClient,
	}
}
