package format

// Color represents an ANSI color code
type Color string

// ANSI color codes
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

// Colorize applies a color to text and resets at the end
func Colorize(color Color, text string) string {
	return string(color) + text + string(ColorReset)
}

// Dimmed applies gray color for dimmed text
func Dimmed(text string) string {
	return Colorize(ColorGray, text)
}

// Green applies green color
func Green(text string) string {
	return Colorize(ColorGreen, text)
}

// Red applies red color
func Red(text string) string {
	return Colorize(ColorRed, text)
}
