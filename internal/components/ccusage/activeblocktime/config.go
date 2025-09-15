package activeblocktime

// Config defines configuration for the active block time component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}             - The configured icon
	//   {{.RemainingMinutes}} - Raw remaining minutes
	//   {{.Remaining}}        - Formatted remaining time (e.g. "2h30m", "45m", "0m" for expired)
	//   {{.EndTime}}          - Formatted end time (e.g. "11:30 PM")
	//   {{.EndTimeRaw}}       - Raw RFC3339 end time string
	//   {{.Expired}}          - Boolean indicating if block has expired
	Template string `yaml:"template"`

	// Icon to display with time
	Icon string `yaml:"icon,omitempty"`

	// Color for the time display
	Color string `yaml:"color,omitempty"`

	// Time format for end time (Go time format)
	// Examples: "3:04 PM", "15:04", "15:04:05"
	EndTimeFormat string `yaml:"end_time_format,omitempty"`
}

// defaultConfig returns the default configuration for active block time component.
func defaultConfig() *Config {
	return &Config{
		Template:      "{{.Icon}} {{.Remaining}}{{if .EndTime}} {{.EndTime}}{{end}}",
		Icon:          "\uf017", // Nerd Font: Clock icon
		Color:         "gray",
		EndTimeFormat: "3:04 PM", // 12-hour format with AM/PM
	}
}
