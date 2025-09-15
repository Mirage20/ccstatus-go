package format

import "fmt"

// DurationMinutes formats minutes into a human-readable duration
// Examples:
//
//	45  -> "45m"
//	90  -> "1h30m"
//	120 -> "2h"
//	0   -> "0m"
//	-5  -> "0m"
func DurationMinutes(minutes int) string {
	if minutes <= 0 {
		return "0m"
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

// DurationMs formats milliseconds into a human-readable duration
// Examples:
//
//	45000  -> "45s"
//	90000  -> "1m30s"
//	120000 -> "2m"
//	3600000 -> "1h"
func DurationMs(ms int64) string {
	if ms <= 0 {
		return "0s"
	}

	seconds := ms / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	secs := seconds % 60

	if minutes < 60 {
		if secs == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, secs)
	}

	hours := minutes / 60
	mins := minutes % 60

	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}
