package stash

import (
	"strconv"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func init() {
	core.RegisterComponent("git.stash", New)
}

// Component displays git stash count.
type Component struct {
	config *Config
}

// New is the factory function for git.stash component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "git.stash", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the git stash display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := gitprovider.GetInfo(ctx)
	if !ok || !info.IsRepo {
		return ""
	}

	// If no stashes, return empty
	if info.Stash == 0 {
		return ""
	}

	// Build template data
	data := map[string]any{
		"Icon":  c.config.Icon,
		"Count": strconv.Itoa(info.Stash),
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
