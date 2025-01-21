package main

import (
	"context"
	"fmt"
	"os"
	"golang.org/x/oauth2"
	"github.com/shurcooL/githubv4"
	"github.com/google/go-github/v50/github"
	"github.com/vbaranov/antiscam-action/pkg/antiscam"
)

func main() {
	eventType := os.Getenv("GITHUB_EVENT_NAME")
	var eventData []byte
	if path := os.Getenv("GITHUB_EVENT_PATH"); path != "" {
		var err error
		eventData, err = os.ReadFile(path)
		if err != nil {
			fail(fmt.Errorf("failed to read event data: %w", err))
		}
	} else {
		fail(fmt.Errorf("event data is required"))
	}

	ctx := context.Background()
	restClient := github.NewTokenClient(ctx, os.Getenv("INPUT_TOKEN"))

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("INPUT_TOKEN")},
	)
	httpClient := oauth2.NewClient(ctx, src)

	graphqlClient := githubv4.NewClient(httpClient)

	a := antiscam.New(ctx, restClient, graphqlClient)

	switch eventType {
	case "issue_comment":
		if err := a.ProcessIssueComment(eventData); err != nil {
			fail(err)
		}
	case "discussion_comment":
		if err := a.ProcessDiscussionComment(eventData); err != nil {
			fail(err)
		}
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Action failed: %s\n", err)
	os.Exit(1)
}
