package sevenday

const (
	// Default thresholds for usage levels (in percentage).
	defaultWarningThreshold  = 60.0
	defaultCriticalThreshold = 80.0
)

// Config defines configuration for the 7-day rate limit component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}         - The configured icon/label (default: "7d")
	//   {{.Utilization}}  - Usage percentage (e.g. 26)
	//   {{.Remaining}}    - Formatted remaining time (e.g. "2d3h")
	//   {{.EndTime}}      - Formatted reset time (e.g. "Sat 4:29 PM")
	//   {{.EndTimeRaw}}   - Raw reset time for custom formatting
	Template string `yaml:"template"`

	// Icon/label to display (default: "7d")
	Icon string `yaml:"icon,omitempty"`

	// Time format for end time (Go time format)
	// Examples: "Mon 3:04 PM", "Jan 2 15:04"
	EndTimeFormat string `yaml:"end_time_format,omitempty"`

	// Color thresholds (percentages)
	WarningThreshold  float64 `yaml:"warning_threshold,omitempty"`
	CriticalThreshold float64 `yaml:"critical_threshold,omitempty"`

	// Colors for different usage levels
	NormalColor   string `yaml:"normal_color,omitempty"`
	WarningColor  string `yaml:"warning_color,omitempty"`
	CriticalColor string `yaml:"critical_color,omitempty"`

	// Color for supplementary info (remaining time, end time)
	Color string `yaml:"color,omitempty"`
}

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		Template:          "{{.Icon}} {{.Utilization}} {{.Remaining}}",
		Icon:              "7d",
		EndTimeFormat:     "Mon 3:04 PM",
		WarningThreshold:  defaultWarningThreshold,
		CriticalThreshold: defaultCriticalThreshold,
		NormalColor:       "green",
		WarningColor:      "yellow",
		CriticalColor:     "red",
		Color:             "gray",
	}
}
