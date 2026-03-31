package analyzer

import "strings"

func Score(diff string) (int, string) {
	score := 0

	if strings.Contains(diff, "/etc") {
		score += 40
	}
	if strings.Contains(diff, "/usr/local") {
		score += 10
	}
	if strings.Contains(diff, "systemd") {
		score += 30
	}
	if strings.Contains(diff, "/var") {
		score += 10
	}

	level := "LOW"
	if score > 60 {
		level = "HIGH"
	} else if score > 30 {
		level = "MEDIUM"
	}

	return score, level
}
