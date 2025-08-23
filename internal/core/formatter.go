package core

// ColorStyle represents text colors
type ColorStyle string

const (
	ColorReset   ColorStyle = "\x1b[0m"
	ColorRed     ColorStyle = "\x1b[31m"
	ColorGreen   ColorStyle = "\x1b[32m"
	ColorYellow  ColorStyle = "\x1b[33m"
	ColorBlue    ColorStyle = "\x1b[34m"
	ColorMagenta ColorStyle = "\x1b[35m"
	ColorCyan    ColorStyle = "\x1b[36m"
	ColorGray    ColorStyle = "\x1b[90m"
)

// Formatter handles text formatting
type Formatter interface {
	// Color applies color to text
	Color(style ColorStyle, text string) string

	// FormatTokens formats token count (e.g., "1.2k", "3.5M")
	FormatTokens(count int64) string

	// FormatDuration formats time duration
	FormatDuration(minutes int) string

	// FormatPercentage formats percentage
	FormatPercentage(value float64) string

	// Icon returns icon for given name
	Icon(name string) string
}