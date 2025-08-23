package formatter

import (
	"fmt"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// DefaultFormatter implements the Formatter interface
type DefaultFormatter struct {
	icons map[string]string
}

// NewDefaultFormatter creates a new default formatter
func NewDefaultFormatter() *DefaultFormatter {
	return &DefaultFormatter{
		icons: map[string]string{
			"ai":      "🤖",  // AI icon
			"context": "📊",  // Context icon
			"clock":   "⏰",  // Clock icon
			"arrow":   "→",  // Arrow icon
			"folder":  "📁",  // Folder icon
			"git":     "🔀",  // Git icon
			"version": "🏷️", // Version icon
		},
	}
}

// Color applies color to text
func (f *DefaultFormatter) Color(style core.ColorStyle, text string) string {
	return fmt.Sprintf("%s%s%s", style, text, core.ColorReset)
}

// FormatTokens formats token count (e.g., "1.2k", "3.5M")
func (f *DefaultFormatter) FormatTokens(count int64) string {
	switch {
	case count >= 1000000:
		return fmt.Sprintf("%.1fM", float64(count)/1000000)
	case count >= 1000:
		return fmt.Sprintf("%dk", count/1000)
	default:
		return fmt.Sprintf("%d", count)
	}
}

// FormatDuration formats time duration
func (f *DefaultFormatter) FormatDuration(minutes int) string {
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

// FormatPercentage formats percentage
func (f *DefaultFormatter) FormatPercentage(value float64) string {
	return fmt.Sprintf("%.0f%%", value)
}

// Icon returns icon for given name
func (f *DefaultFormatter) Icon(name string) string {
	if icon, exists := f.icons[name]; exists {
		return icon
	}
	return ""
}
