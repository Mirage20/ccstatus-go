package thinking

// Config defines configuration for the thinking component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}    - Icon for the current state (Icon when enabled, OffIcon when disabled)
	//   {{.Enabled}} - Boolean, for use with {{if .Enabled}}...{{end}}
	Template string `yaml:"template"`

	// Icon to display while thinking is enabled (only shown when
	// HideWhenEnabled is false).
	Icon string `yaml:"icon,omitempty"`

	// Icon to display while thinking is disabled.
	OffIcon string `yaml:"off_icon,omitempty"`

	// Color while thinking is enabled.
	Color string `yaml:"color,omitempty"`

	// Color while thinking is disabled.
	OffColor string `yaml:"off_color,omitempty"`

	// HideWhenEnabled controls whether to render anything while thinking is
	// enabled. Extended thinking is on by default in Claude Code, so a badge
	// that lights up when enabled would be visible in almost every session
	// and carry no information; the informative event is thinking being
	// turned OFF. When true (default), the component is invisible in the
	// common case and only appears once thinking is disabled.
	HideWhenEnabled bool `yaml:"hide_when_enabled,omitempty"`
}

// defaultConfig returns the default configuration for the thinking component.
func defaultConfig() *Config {
	return &Config{
		Template:        "{{.Icon}}",
		Icon:            "\uf075", // nf-fa-comment
		OffIcon:         "\ued96", // nf-fa-comment_slash
		Color:           "cyan",
		OffColor:        "gray",
		HideWhenEnabled: true,
	}
}
