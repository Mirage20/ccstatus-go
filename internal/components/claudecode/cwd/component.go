package cwd

import (
	"path/filepath"
	"regexp"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	core.RegisterComponent("cwd", New)
}

// Component displays the current working directory basename.
type Component struct {
	config         *Config
	ignorePatterns []*regexp.Regexp
}

// New is the factory function for cwd component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "cwd", defaultConfig())

	// Pre-compile ignore patterns
	var patterns []*regexp.Regexp
	for _, pattern := range cfg.Ignore {
		if re, err := regexp.Compile(pattern); err == nil {
			patterns = append(patterns, re)
		}
	}

	return &Component{
		config:         cfg,
		ignorePatterns: patterns,
	}
}

// Render generates the cwd display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	currentDir := info.Workspace.CurrentDir
	if currentDir == "" {
		return ""
	}

	dir := filepath.Base(currentDir)

	// Check if directory matches any ignore pattern
	for _, re := range c.ignorePatterns {
		if re.MatchString(dir) {
			return ""
		}
	}

	// Truncate from middle if exceeds max length
	// Keeps prefix and suffix: "my-very-long-directory" → "my-very…ectory"
	if c.config.MaxLength > 0 && len(dir) > c.config.MaxLength {
		// Split evenly: first half + … + second half
		halfLen := (c.config.MaxLength - 1) / 2 //nolint:mnd // split in half
		firstLen := halfLen
		lastLen := c.config.MaxLength - 1 - firstLen
		dir = dir[:firstLen] + "…" + dir[len(dir)-lastLen:]
	}

	// Build template data
	data := map[string]any{
		"Dir":  dir,
		"Icon": c.config.Icon,
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Apply color to the output
	color := format.ParseColor(c.config.Color)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}
