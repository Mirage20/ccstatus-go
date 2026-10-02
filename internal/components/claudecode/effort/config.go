package effort

// fallbackIcon and fallbackColor are used for an effort level with no entry
// in Icons/Colors (an unrecognized level, or the unsupported-label text).
// Not configurable: a missing map entry is an edge case, not something worth
// a config key for.
const (
	fallbackIcon  = "\uf420" // nf-fa-question
	fallbackColor = "gray"
)

// Config defines configuration for the effort component.
type Config struct {
	// Display template
	// Available template parameters:
	//   {{.Level}} - The effective effort level ("low", "medium", "high", "xhigh", "max")
	//   {{.Icon}}  - The icon mapped to the current level
	Template string `yaml:"template"`

	// Icons per effort level. Lookup is an exact match on the level string -
	// substring matching would be ambiguous because "high" is a substring of
	// "xhigh". Levels not listed fall back to fallbackIcon.
	//
	// Only ever one of low/medium/high/xhigh/max: Claude Code resolves
	// meta-modes like "ultracode" (-> xhigh) and "auto" (-> model default)
	// to a concrete level before it ever reaches stdin, so there's no
	// separate entry to add for them here.
	Icons map[string]string `yaml:"icons,omitempty"`

	// Colors per effort level. Levels not listed fall back to fallbackColor.
	Colors map[string]string `yaml:"colors,omitempty"`

	// UnsupportedLabel stands in for the level when the model does not
	// support effort (Claude Code omits the "effort" key entirely in that
	// case). Only used when HideUnsupported is false.
	//
	// It doubles as the Icons/Colors lookup key, so keep it to a value no
	// real level uses - otherwise it borrows that level's icon and color.
	UnsupportedLabel string `yaml:"unsupported_label,omitempty"`

	// HideUnsupported controls whether to render anything when no effort
	// level is reported. When true (default), the component renders an empty
	// string and disappears from the status line, separator included.
	HideUnsupported bool `yaml:"hide_unsupported,omitempty"`
}

// defaultConfig returns the default configuration for the effort component.
func defaultConfig() *Config {
	return &Config{
		Template: "{{.Icon}} {{.Level}}",
		Icons: map[string]string{
			"low":    "\ueeef", // nf-fa-cloud_moon
			"medium": "\uee1d", // nf-fa-person_walking
			"high":   "\uef0c", // nf-fa-person_running
			"xhigh":  "\uee9c", // nf-fa-brain
			"max":    "\uf06d", // nf-fa-fire
		},
		Colors: map[string]string{
			"low":    "gray",
			"medium": "cyan",
			"high":   "yellow",
			"xhigh":  "magenta",
			"max":    "red",
		},
		UnsupportedLabel: "--",
		HideUnsupported:  true,
	}
}
