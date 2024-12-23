package main

import (
	"context"
	"fmt"
	"os"

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
	client := github.NewTokenClient(ctx, os.Getenv("INPUT_TOKEN"))
	a := antiscam.New(ctx, client)

	switch eventType {
	case "issue_comment":
		if err := a.ProcessIssueComment(eventData); err != nil {
			fail(err)
		}
	}
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "Action failed: %s\n", err)
	os.Exit(1)
}
