package antiscam

import (
	"context"

	"github.com/google/go-github/v50/github"
)

type Antiscam struct {
	ctx    context.Context
	client *github.Client
}

func New(ctx context.Context, client *github.Client) *Antiscam {
	return &Antiscam{
		ctx:    ctx,
		client: client,
	}
}
