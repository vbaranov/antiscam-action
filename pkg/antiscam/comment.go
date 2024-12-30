package antiscam

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v50/github"
)

func (a *Antiscam) ProcessIssueComment(payload []byte) error {
	var event github.IssueCommentEvent
	var event2 github.DiscussionEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}
	if err2 := json.Unmarshal(payload, &event2); err2 != nil {
		return err2
	}

	var detections []Detection
	detections = append(detections, checkComment(event.GetComment().GetBody())...)

	// body := fmt.Sprintf("@%s The previous user tried to scam you by providing a fake support link. Don't interact with it.\n", event.GetIssue().GetUser().GetLogin())

	for _, detection := range detections {
		fmt.Printf("Detected scam in %s: %s\n", detection.Location, detection.DebugInfo)
	}

	fmt.Printf("Organization ID %s\n", event.GetOrganization().GetName())
	fmt.Printf("Team ID %s\n", event.GetRepo().GetName())
	fmt.Printf("Issue number %d\n", event.Issue.GetNumber())
	fmt.Printf("Discussion number %d\n", event2.Discussion.GetNumber())
	fmt.Printf("Discussion comment number %d\n", int(event.GetComment().GetID()))

	if len(detections) > 0 {
		if _, err := a.client.Teams.DeleteCommentBySlug(
			a.ctx,
			event.GetOrganization().GetName(),
			event.GetRepo().GetName(),
			event2.Discussion.GetNumber(),
			int(event.GetComment().GetID()),
		); err != nil {
			return err
		}

		// if _, _, err := a.client.Issues.CreateComment(
		// 	a.ctx,
		// 	event.GetRepo().GetOwner().GetLogin(),
		// 	event.GetRepo().GetName(),
		// 	event.GetIssue().GetNumber(),
		// 	&github.IssueComment{
		// 		Body: &body,
		// 	},
		// ); err != nil {
		// 	return err
		// }
	}

	return nil
}
