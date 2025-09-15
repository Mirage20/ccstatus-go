package version

// Config defines configuration for the version component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Version}} - The Claude Code version (e.g. "1.0.89")
	//   {{.Icon}}    - The configured icon (if any)
	Template string `yaml:"template"`

	// Icon to display with version
	Icon string `yaml:"icon,omitempty"`

	// Color for the version display
	Color string `yaml:"color,omitempty"`
}

// defaultConfig returns the default configuration for version component.
func defaultConfig() *Config {
	return &Config{
		Template: "v{{.Version}}",
		Icon:     "",     // No icon by default
		Color:    "gray", // Dimmed/gray color for version
	}
}
