package antiscam

import (
	"fmt"
	"os"
	"strings"
)

var whitelisted_logins = map[string]bool{}

func checkComment(body string, comment_author string) []Detection {
	fmt.Printf("comment author: %s\n", comment_author)
	fmt.Printf("env SCAM_ACTION_WHITELISTED_LOGINS value: %s\n", os.Getenv("SCAM_ACTION_WHITELISTED_LOGINS"))
	env_whitelisted_logins := os.Getenv("SCAM_ACTION_WHITELISTED_LOGINS")
	if env_whitelisted_logins != "" {
		for _, login := range strings.Split(env_whitelisted_logins, ",") {
			whitelisted_logins[strings.ToLower(strings.TrimSpace(login))] = true
		}
	}
	fmt.Printf("whitelisted_logins: %v\n", whitelisted_logins)
	var detections []Detection
	if !whitelisted_logins[strings.ToLower(comment_author)] {
		body_lower_case := strings.ToLower(body)
		if strings.Contains(body_lower_case, ".web.app") || ((strings.Contains(body_lower_case, "https://") &&
			!strings.Contains(body_lower_case, "https://github.com") &&
			!strings.Contains(body_lower_case, "https://discord.gg/blockscout") &&
			!strings.Contains(body_lower_case, "https://docs.blockscout.com")) ||
			strings.Contains(body_lower_case, "http://")) &&
			(strings.Contains(body_lower_case, "support") ||
				strings.Contains(body_lower_case, "forum") ||
				strings.Contains(body_lower_case, "help") ||
				strings.Contains(body_lower_case, "dapps portal")) {
			detections = append(detections, Detection{
				Location:  "body",
				DebugInfo: "Comment body contains scammy text.",
			})
		}
	} else {
		fmt.Printf("Author is whitelisted: %s\n", comment_author)
	}
	return detections
}
