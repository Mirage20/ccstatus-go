package fastmode

// Config defines configuration for the fastmode component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Icon}}    - Icon for the current state (Icon when active, OffIcon when off)
	//   {{.Enabled}} - Boolean, for use with {{if .Enabled}}...{{end}}
	Template string `yaml:"template"`

	// Icon to display while fast mode is active.
	Icon string `yaml:"icon,omitempty"`

	// Icon to display while fast mode is off (only shown when HideWhenOff is false).
	OffIcon string `yaml:"off_icon,omitempty"`

	// Color while fast mode is active.
	Color string `yaml:"color,omitempty"`

	// Color while fast mode is off.
	OffColor string `yaml:"off_color,omitempty"`

	// HideWhenOff controls whether to render anything while fast mode is off.
	// Fast mode is opt-in (toggled with /fast) and off by default, so the
	// informative event is it being turned ON. When true (default), the
	// component is invisible until fast mode is activated.
	HideWhenOff bool `yaml:"hide_when_off,omitempty"`
}

// defaultConfig returns the default configuration for the fastmode component.
func defaultConfig() *Config {
	return &Config{
		Template:    "{{.Icon}}",
		Icon:        "\uf0e7", // nf-fa-flash
		OffIcon:     "",
		Color:       "yellow",
		OffColor:    "gray",
		HideWhenOff: true,
	}
}
