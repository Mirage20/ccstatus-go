package branch

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func init() {
	core.RegisterComponent("git.branch", New)
}

// Component displays the git branch name.
type Component struct {
	config *Config
}

// New is the factory function for git.branch component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "git.branch", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the git branch display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := gitprovider.GetInfo(ctx)
	if !ok || !info.IsRepo {
		return ""
	}

	branch := info.Branch

	// Truncate from middle if exceeds max length
	if c.config.MaxLength > 0 && len(branch) > c.config.MaxLength {
		branch = truncateMiddle(branch, c.config.MaxLength)
	}

	// Build template data
	data := map[string]any{
		"Icon":   c.config.Icon,
		"Branch": branch,
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Apply color
	color := format.ParseColor(c.config.Color)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"git"}
}

// truncateMiddle truncates a string from the middle with ellipsis.
func truncateMiddle(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	// Split in half, accounting for ellipsis
	halfLen := (maxLen - 1) / 2 //nolint:mnd // split in half minus ellipsis
	firstLen := halfLen
	lastLen := maxLen - 1 - firstLen
	return s[:firstLen] + "…" + s[len(s)-lastLen:]
}
