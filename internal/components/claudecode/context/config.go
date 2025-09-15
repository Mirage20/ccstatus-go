package context

// Config defines configuration for the context (token usage) component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}       - The configured icon
	//   {{.Total}}      - Total token count (raw number)
	//   {{.Formatted}}  - Formatted token count with unit (e.g. "22k")
	//   {{.Percentage}} - Usage percentage as float (use {{printf "%.0f" .Percentage}} for whole number)
	//   {{.Limit}}      - Context limit
	Template string `yaml:"template"`

	// Icon to display with context
	Icon string `yaml:"icon,omitempty"`

	// Context limit (default 200k for Claude models)
	ContextLimit int64 `yaml:"context_limit,omitempty"`

	// Color thresholds (percentages)
	WarningThreshold  float64 `yaml:"warning_threshold,omitempty"`  // Yellow color threshold
	CriticalThreshold float64 `yaml:"critical_threshold,omitempty"` // Red color threshold

	// Colors for different usage levels
	NormalColor   string `yaml:"normal_color,omitempty"`
	WarningColor  string `yaml:"warning_color,omitempty"`
	CriticalColor string `yaml:"critical_color,omitempty"`
}

// defaultConfig returns the default configuration for context component.
func defaultConfig() *Config {
	return &Config{
		Template:          "{{.Icon}} {{.Formatted}}",
		Icon:              "\uea7b", // Nerd Font: Context/Tokens icon
		ContextLimit:      200000,   // 200k default for Claude
		WarningThreshold:  80.0,
		CriticalThreshold: 90.0,
		NormalColor:       "green",
		WarningColor:      "yellow",
		CriticalColor:     "red",
	}
}
