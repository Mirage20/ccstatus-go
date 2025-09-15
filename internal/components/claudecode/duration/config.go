package duration

// Config defines configuration for the duration component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}         - The main duration icon
	//   {{.APIIcon}}      - The API duration icon
	//   {{.TotalDuration}} - Formatted total duration (e.g. "2m 15s")
	//   {{.APIDuration}}   - Formatted API duration (e.g. "1m 30s")
	//   {{.TotalMs}}       - Raw total duration in milliseconds
	//   {{.APIMs}}         - Raw API duration in milliseconds
	Template string `yaml:"template"`

	// Icon for total duration
	Icon string `yaml:"icon,omitempty"`

	// Icon for API duration
	APIIcon string `yaml:"api_icon,omitempty"`

	// Color for the duration display
	Color string `yaml:"color,omitempty"`

	// Whether to show API duration
	ShowAPIDuration bool `yaml:"show_api_duration,omitempty"`
}

// defaultConfig returns the default configuration for duration component.
func defaultConfig() *Config {
	return &Config{
		Template:        "{{.Icon}} {{.TotalDuration}}{{if .APIDuration}} {{.APIIcon}} {{.APIDuration}}{{end}}",
		Icon:            "\uF520",     // Nerd Font: Stopwatch icon
		APIIcon:         "\U000F1616", // Nerd Font: Connected plug icon
		Color:           "gray",       // Dimmed/gray color
		ShowAPIDuration: true,
	}
}
