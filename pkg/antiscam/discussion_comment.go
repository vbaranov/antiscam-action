package antiscam

import (
	"encoding/json"
	"fmt"

	"github.com/shurcooL/githubv4"
)

type DiscussionCommentPayload struct {
	Comment struct {
		NodeID string `json:"node_id"`
		Body   string `json:"body"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"comment"`
	Discussion struct {
		NodeID string `json:"node_id"`
	} `json:"discussion"`
}

type deleteCommentMutation struct {
	DeleteDiscussionComment struct {
		ClientMutationID *string
	} `graphql:"deleteDiscussionComment(input: $input)"`
}

type addCommentMutation struct {
	AddDiscussionComment struct {
		Comment struct {
			ID   githubv4.ID
			Body githubv4.String
			URL  githubv4.URI
		}
		ClientMutationID *string
	} `graphql:"addDiscussionComment(input: $input)"`
}

func (a *Antiscam) ProcessDiscussionComment(payload []byte) error {
	var event DiscussionCommentPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	detections := checkComment(event.Comment.Body, event.Comment.User.Login)
	if len(detections) == 0 {
		return nil
	}

	for _, detection := range detections {
		fmt.Printf("Detected scam in %s: %s\n", detection.Location, detection.DebugInfo)
	}

	if err := a.deleteDiscussionScamComment(event.Comment.NodeID); err != nil {
		return fmt.Errorf("failed to delete scam comment: %w", err)
	}

	if err := a.addWarningDiscussionComment(event.Discussion.NodeID, event.Comment.User.Login); err != nil {
		return fmt.Errorf("failed to add warning comment: %w", err)
	}

	return nil
}

func (a *Antiscam) deleteDiscussionScamComment(commentID string) error {
	var mutation deleteCommentMutation
	input := githubv4.DeleteDiscussionCommentInput{
		ID: githubv4.ID(commentID),
	}

	return a.graphqlClient.Mutate(a.ctx, &mutation, input, nil)
}

func (a *Antiscam) addWarningDiscussionComment(discussionID string, userLogin string) error {
	var mutation addCommentMutation
	body := fmt.Sprintf(
		"@%s The previous user tried to scam you by providing a fake support link. Don't interact with it.\n",
		userLogin,
	)

	input := githubv4.AddDiscussionCommentInput{
		DiscussionID: githubv4.ID(discussionID),
		Body:         githubv4.String(body),
	}

	return a.graphqlClient.Mutate(a.ctx, &mutation, input, nil)
}
