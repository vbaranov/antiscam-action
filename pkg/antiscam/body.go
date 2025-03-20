package antiscam

import (
	"fmt"
	"os"
	"strings"
)

const ScammyTextDebugInfo = "Comment body contains scammy text."

var whitelisted_logins = map[string]bool{}

func checkComment(body string, comment_author string) []Detection {
	fmt.Printf("comment author: %s\n", comment_author)
	env_whitelisted_logins := os.Getenv("SCAM_ACTION_WHITELISTED_LOGINS")
	fmt.Printf("env SCAM_ACTION_WHITELISTED_LOGINS value: %s\n", env_whitelisted_logins)
	if env_whitelisted_logins != "" {
		for _, login := range strings.Split(env_whitelisted_logins, ",") {
			whitelisted_logins[strings.ToLower(strings.TrimSpace(login))] = true
		}
	}
	fmt.Printf("whitelisted_logins: %v\n", whitelisted_logins)
	var detections []Detection
	if !whitelisted_logins[strings.ToLower(comment_author)] {
		body_lower_case := strings.ToLower(body)
		solidScammyPatterns := []string{
			".web.app",
			"telegram",
			"@",
		}
		scammyPatterns := append(solidScammyPatterns, "https://", "http://", ".com")
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

		isPotentiallyScammy := false
		isScammy := false
		for _, pattern := range scammyPatterns {
			isPotentiallyScammy = true
			if strings.Contains(body_lower_case, pattern) {
				if contains(solidScammyPatterns, pattern) {
					detections = appendDetection(detections)
					isScammy = true
					break
				}
				if pattern == "https://" {
					for _, exception := range exceptions {
						if strings.Contains(body_lower_case, exception) {
							isPotentiallyScammy = false
						}
					}
				} else {
					break
				}
			}
		}

		if isPotentiallyScammy && !isScammy {
			for _, supportPattern := range supportPatterns {
				if strings.Contains(body_lower_case, supportPattern) {
					detections = appendDetection(detections)
					break
				}
			}
		}
	} else {
		fmt.Printf("Author is whitelisted: %s\n", comment_author)
	}
	return detections
}

func appendDetection(detections []Detection) []Detection {
	return append(detections, Detection{
		Location:  "body",
		DebugInfo: ScammyTextDebugInfo,
	})
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
