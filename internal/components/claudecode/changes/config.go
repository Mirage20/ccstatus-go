package changes

// Config defines configuration for the changes component.
type Config struct {
	// Display template
	// Available template parameters (all pre-colored):
	//   {{.Icon}}        - The configured icon (uses component color)
	//   {{.Added}}       - Number of lines added (uses added color)
	//   {{.Removed}}     - Number of lines removed (uses removed color)
	//   {{.AddedSign}}   - The configured added sign (uses added color)
	//   {{.RemovedSign}} - The configured removed sign (uses removed color)
	Template string `yaml:"template"`

	// Icon to display with changes (optional)
	Icon string `yaml:"icon,omitempty"`

	// Color for the changes display (optional)
	Color string `yaml:"color,omitempty"`

	// Signs/icons for added and removed lines
	AddedSign   string `yaml:"added_sign,omitempty"`
	RemovedSign string `yaml:"removed_sign,omitempty"`

	// Colors for added/removed lines
	AddedColor   string `yaml:"added_color,omitempty"`
	RemovedColor string `yaml:"removed_color,omitempty"`

	// Whether to show zero values (e.g. "+0")
	ShowZero bool `yaml:"show_zero,omitempty"`
}

// defaultConfig returns the default configuration for changes component.
func defaultConfig() *Config {
	return &Config{
		Template:     "{{.AddedSign}}{{.Added}}{{.RemovedSign}}{{.Removed}}",
		Icon:         "", // No icon by default
		AddedSign:    "+",
		RemovedSign:  "-",
		AddedColor:   "green",
		RemovedColor: "red",
		ShowZero:     false, // Don't show component if both are zero
	}
}
