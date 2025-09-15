package model

// Config defines configuration for the model component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.ID}}        - The model ID (e.g. "claude-opus-4-1-20250805")
	//   {{.Name}}      - The display name (e.g. "Opus 4.1")
	//   {{.ShortName}} - The extracted short name (e.g. "Opus", "Sonnet", "Haiku")
	//   {{.Icon}}      - The matched icon based on model ID patterns
	Template string `yaml:"template"`

	// Visual customization (pattern -> value)
	// Pattern matching is done via substring search on model ID
	Icons  map[string]string `yaml:"icons,omitempty"`
	Colors map[string]string `yaml:"colors,omitempty"`
}

// defaultConfig returns the default configuration for model component.
func defaultConfig() *Config {
	return &Config{
		Template: "{{.Icon}} {{.ShortName}}",
		Icons: map[string]string{
			"opus":   "\uf2db", // Nerd Font: Microchip icon
			"sonnet": "\uf2db", // Nerd Font: Microchip icon
			"haiku":  "\uf2db", // Nerd Font: Microchip icon
		},
		Colors: map[string]string{
			"opus":   "magenta",
			"sonnet": "cyan",
			"haiku":  "green",
		},
	}
}
