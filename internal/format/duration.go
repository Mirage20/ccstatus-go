package format

import "fmt"

// FormatDuration formats minutes into a human-readable duration
// Examples:
//   45  -> "45m"
//   90  -> "1h30m"
//   120 -> "2h"
//   0   -> "expired"
//   -5  -> "expired"
func FormatDuration(minutes int) string {
	if minutes <= 0 {
		return "expired"
	}

	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}