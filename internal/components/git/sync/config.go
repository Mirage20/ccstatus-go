package sync

// Config defines configuration for the git.sync component.
type Config struct {
	// Template for display.
	// Available variables (all pre-colored):
	//   {{.Ahead}}  - Commits ahead of upstream (e.g., "↑2")
	//   {{.Behind}} - Commits behind upstream (e.g., "↓3")
	Template string `yaml:"template"`

	// Icons for ahead/behind.
	AheadIcon  string `yaml:"ahead_icon,omitempty"`
	BehindIcon string `yaml:"behind_icon,omitempty"`

	// Colors for ahead/behind.
	AheadColor  string `yaml:"ahead_color,omitempty"`
	BehindColor string `yaml:"behind_color,omitempty"`
}

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		Template:    "{{if .Ahead}} {{.Ahead}}{{end}}{{if .Behind}} {{.Behind}}{{end}}",
		AheadIcon:   "\ueaa1 ",
		BehindIcon:  "\uea9a ",
		AheadColor:  "green",
		BehindColor: "red",
	}
}
