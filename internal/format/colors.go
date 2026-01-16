package format

import "strings"

// Color represents an ANSI color code.
type Color string

// ANSI color codes.
const (
	ColorReset   Color = "\033[0m"
	ColorRed     Color = "\033[31m"
	ColorGreen   Color = "\033[32m"
	ColorYellow  Color = "\033[33m"
	ColorBlue    Color = "\033[34m"
	ColorMagenta Color = "\033[35m"
	ColorCyan    Color = "\033[36m"
	ColorGray    Color = "\033[90m"
)

// Colorize applies a color to text and resets at the end.
// Returns empty string if text is empty (no color codes for nothing).
func Colorize(color Color, text string) string {
	if text == "" {
		return ""
	}
	return string(color) + text + string(ColorReset)
}

// ParseColor converts a color name string to a Color constant.
// Returns ColorGray as default if the name is not recognized.
func ParseColor(name string) Color {
	colors := map[string]Color{
		"red":     ColorRed,
		"green":   ColorGreen,
		"yellow":  ColorYellow,
		"blue":    ColorBlue,
		"magenta": ColorMagenta,
		"cyan":    ColorCyan,
		"gray":    ColorGray,
		"grey":    ColorGray, // Alternative spelling
	}

	if color, ok := colors[strings.ToLower(name)]; ok {
		return color
	}
	return ColorGray // Default to gray for unknown colors
}
