package antiscam

import "strings"

func checkComment(body string) []Detection {
	var detections []Detection
	body_lower_case := strings.ToLower(body)
	if ((strings.Contains(body_lower_case, "https://") || strings.Contains(body_lower_case, "http://")) && (strings.Contains(body_lower_case, "support") || strings.Contains(body_lower_case, "forum") || strings.Contains(body_lower_case, "help"))) {
		detections = append(detections, Detection{
			Location:       "body",
			DebugInfo:      "Comment body contains scammy text.",
		})
	}
	return detections
}
