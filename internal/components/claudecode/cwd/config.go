package cwd

const (
	// Default max length for directory name truncation.
	defaultMaxLength = 20
)

// Config defines configuration for the cwd component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Dir}}  - The current directory basename
	//   {{.Icon}} - The configured icon
	Template string `yaml:"template"`

	// Icon to display with directory
	Icon string `yaml:"icon,omitempty"`

	// Color for the display
	Color string `yaml:"color,omitempty"`

	// Ignore patterns (regex) - hide component when directory matches
	Ignore []string `yaml:"ignore,omitempty"`

	// Maximum length for directory name (0 = no limit)
	// Truncates from middle with ellipsis: "my-very…ectory"
	MaxLength int `yaml:"max_length,omitempty"`
}

// defaultConfig returns the default configuration for cwd component.
func defaultConfig() *Config {
	return &Config{
		Template:  "{{.Icon}} {{.Dir}}",
		Icon:      "\uf07b", // Folder icon
		Color:     "gray",
		Ignore:    []string{},
		MaxLength: defaultMaxLength,
	}
}
