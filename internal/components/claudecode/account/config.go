package account

// Config defines configuration for the account component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Profile}} - The active profile name (e.g. "personal", "work", "unknown")
	//   {{.Icon}}    - The configured icon (if any)
	Template string `yaml:"template"`

	// Icon to display before the profile name (optional)
	Icon string `yaml:"icon,omitempty"`

	// Colors per profile name. Profiles not listed fall back to DefaultColor.
	Colors map[string]string `yaml:"colors,omitempty"`

	// DefaultColor is used when the active profile has no entry in Colors.
	DefaultColor string `yaml:"default_color,omitempty"`

	// Hide controls whether to render anything when no profile is detected
	// (active token missing, or doesn't match any saved slot).
	// When true (default), the component renders an empty string in those cases.
	HideUnknown bool `yaml:"hide_unknown,omitempty"`
}

// defaultConfig returns the default configuration for the account component.
func defaultConfig() *Config {
	return &Config{
		Template: "{{.Icon}}{{.Profile}}",
		Icon:     "",
		Colors: map[string]string{
			"personal": "cyan",
			"work":     "magenta",
		},
		DefaultColor: "gray",
		HideUnknown:  true,
	}
}
