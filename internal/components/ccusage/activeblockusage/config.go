package activeblockusage

// Config defines configuration for the active block usage component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}            - The configured icon
	//   {{.TotalTokens}}     - Raw token count
	//   {{.Formatted}}       - Formatted token count with unit (e.g. "350k")
	//   {{.UsagePercentage}} - Usage percentage (e.g. 75.5)
	//   {{.Limit}}           - Active limit (user config or dynamic max)
	//   {{.MaxBlockTokens}}  - Dynamic max from historical data
	Template string `yaml:"template"`

	// Icon to display with block usage
	Icon string `yaml:"icon,omitempty"`

	// Block token limit override (takes precedence over dynamic calculation)
	// If not set, uses historical maximum from all blocks
	BlockLimit int64 `yaml:"block_limit,omitempty"`

	// Color thresholds (percentages)
	WarningThreshold  float64 `yaml:"warning_threshold,omitempty"`  // Yellow color threshold
	CriticalThreshold float64 `yaml:"critical_threshold,omitempty"` // Red color threshold

	// Colors for different usage levels
	NormalColor   string `yaml:"normal_color,omitempty"`
	WarningColor  string `yaml:"warning_color,omitempty"`
	CriticalColor string `yaml:"critical_color,omitempty"`
}

// defaultConfig returns the default configuration for active block usage component.
func defaultConfig() *Config {
	return &Config{
		Template:          "{{.Icon}} {{.Formatted}} {{printf \"%.0f\" .UsagePercentage}}%",
		Icon:              "\U000F0E7A", // Nerd Font: Up-down arrow icon
		BlockLimit:        0,            // 0 means use dynamic calculation
		WarningThreshold:  60.0,
		CriticalThreshold: 80.0,
		NormalColor:       "green",
		WarningColor:      "yellow",
		CriticalColor:     "red",
	}
}
