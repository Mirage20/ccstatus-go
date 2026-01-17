package status

// Config defines configuration for the git.status component.
type Config struct {
	// Template for display.
	// Available variables (all pre-colored):
	//   {{.Staged}}    - Staged files indicator (e.g., " 3")
	//   {{.Modified}}  - Modified files indicator (e.g., " 2")
	//   {{.Untracked}} - Untracked files indicator (e.g., " 1")
	//   {{.Conflicts}} - Conflicted files indicator (e.g., " 2")
	Template string `yaml:"template"`

	// Icons/prefixes for each status type.
	StagedIcon    string `yaml:"staged_icon,omitempty"`
	ModifiedIcon  string `yaml:"modified_icon,omitempty"`
	UntrackedIcon string `yaml:"untracked_icon,omitempty"`
	ConflictIcon  string `yaml:"conflict_icon,omitempty"`

	// Colors for each status type.
	StagedColor    string `yaml:"staged_color,omitempty"`
	ModifiedColor  string `yaml:"modified_color,omitempty"`
	UntrackedColor string `yaml:"untracked_color,omitempty"`
	ConflictColor  string `yaml:"conflict_color,omitempty"`
}

// defaultConfig returns the default configuration.
func defaultConfig() *Config {
	return &Config{
		Template:       "{{if .Staged}} {{.Staged}}{{end}}{{if .Modified}} {{.Modified}}{{end}}{{if .Untracked}} {{.Untracked}}{{end}}{{if .Conflicts}} {{.Conflicts}}{{end}}",
		StagedIcon:     "\uf05d ", // nf-cod-diff_added
		ModifiedIcon:   "\uf044 ", // nf-fa-pencil
		UntrackedIcon:  "\uf420 ", // nf-fa-question
		ConflictIcon:   "\uf421 ", // nf-fa-exclamation
		StagedColor:    "green",
		ModifiedColor:  "yellow",
		UntrackedColor: "gray",
		ConflictColor:  "red",
	}
}
