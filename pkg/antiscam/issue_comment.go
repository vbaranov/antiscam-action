package antiscam

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v69/github"
)

func (a *Antiscam) ProcessIssueComment(payload []byte) error {
	var event github.IssueCommentEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	fmt.Printf("issue comment: %s\n", github.Stringify(event))

	var detections []Detection
	detections = append(detections, checkComment(*event.Comment.Body, *event.Comment.User.Login)...)

	body := fmt.Sprintf("@%s The previous user tried to scam you by providing a fake support link. Don't interact with it.\n", event.GetIssue().GetUser().GetLogin())

	for _, detection := range detections {
		fmt.Printf("Detected scam in %s: %s\n", detection.Location, detection.DebugInfo)
	}

	if len(detections) > 0 {
		a.restClient.Issues.DeleteComment(
			a.ctx,
			event.GetRepo().GetOwner().GetLogin(),
			event.GetRepo().GetName(),
			event.GetComment().GetID(),
		)

		if _, _, err := a.restClient.Issues.CreateComment(
			a.ctx,
			event.GetRepo().GetOwner().GetLogin(),
			event.GetRepo().GetName(),
			event.GetIssue().GetNumber(),
			&github.IssueComment{
				Body: &body,
			},
		); err != nil {
			return err
		}
	}

	return nil
}
