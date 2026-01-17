package stash

// Config defines configuration for the git.stash component.
type Config struct {
	// Template for display.
	// Available variables:
	//   {{.Icon}}  - The configured icon
	//   {{.Count}} - Number of stash entries
	Template string `yaml:"template"`

	// Icon for stash.
	Icon string `yaml:"icon,omitempty"`

	// Color for display.
	Color string `yaml:"color,omitempty"`
}

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		Template: "{{.Icon}} {{.Count}}",
		Icon:     "\uf48d", // nf-fa-inbox
		Color:    "cyan",
	}
}
