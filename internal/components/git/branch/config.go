package branch

const (
	// Default max length for branch name truncation.
	defaultMaxLength = 20
)

// Config defines configuration for the git.branch component.
type Config struct {
	// Template for display.
	// Available variables:
	//   {{.Icon}}   - The configured icon
	//   {{.Branch}} - Branch name or @hash for detached HEAD
	Template string `yaml:"template"`

	// Icon to display.
	Icon string `yaml:"icon,omitempty"`

	// Color for the display.
	Color string `yaml:"color,omitempty"`

	// MaxLength for branch name (0 = no limit).
	// Truncates from middle with ellipsis.
	MaxLength int `yaml:"max_length,omitempty"`
}

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		Template:  "{{.Icon}} {{.Branch}}",
		Icon:      "\uE725", // Git branch nerd font icon
		Color:     "gray",
		MaxLength: defaultMaxLength,
	}
}
