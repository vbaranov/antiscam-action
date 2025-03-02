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
		scammyPatterns := []string{
			".web.app",
			"https://",
			"http://",
			".com",
			"telegram",
			"@",
		}
		exceptions := []string{
			"https://github.com",
			"https://discord.gg/blockscout",
			"https://docs.blockscout.com",
		}
		supportPatterns := []string{
			"support",
			"forum",
			"help",
			"dapps portal",
		}

		isScammy := false
		for _, pattern := range scammyPatterns {
			isScammy = true
			if strings.Contains(body_lower_case, pattern) {
				if pattern == "https://" {
					for _, exception := range exceptions {
						if strings.Contains(body_lower_case, exception) {
							isScammy = false
						}
					}
				} else {
					break
				}
			}
		}

		if isScammy {
			for _, supportPattern := range supportPatterns {
				if strings.Contains(body_lower_case, supportPattern) {
					detections = append(detections, Detection{
						Location:  "body",
						DebugInfo: "Comment body contains scammy text.",
					})
					break
				}
			}
		}
	} else {
		fmt.Printf("Author is whitelisted: %s\n", comment_author)
	}
	return detections
}
