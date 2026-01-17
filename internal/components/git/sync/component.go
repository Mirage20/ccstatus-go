package sync

import (
	"fmt"
	"strings"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func init() {
	core.RegisterComponent("git.sync", New)
}

// Component displays git sync status (ahead/behind upstream).
type Component struct {
	config *Config
}

// New is the factory function for git.sync component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "git.sync", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the git sync display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := gitprovider.GetInfo(ctx)
	if !ok || !info.IsRepo {
		return ""
	}

	// If no upstream configured, return empty
	if !info.HasUpstream {
		return ""
	}

	// If both ahead and behind are zero, return empty (in sync)
	if info.Ahead == 0 && info.Behind == 0 {
		return ""
	}

	// Parse colors
	aheadColor := format.ParseColor(c.config.AheadColor)
	behindColor := format.ParseColor(c.config.BehindColor)

	// Build pre-colored indicators
	ahead := c.formatCount(info.Ahead, c.config.AheadIcon, aheadColor)
	behind := c.formatCount(info.Behind, c.config.BehindIcon, behindColor)

	// Build template data with pre-colored values
	data := map[string]any{
		"Ahead":  ahead,
		"Behind": behind,
	}

	// Render template and trim leading space
	result := format.RenderTemplate(c.config.Template, data)
	return strings.TrimLeft(result, " ")
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"git"}
}

// formatCount returns colorized "icon+count" if count > 0, empty string otherwise.
func (c *Component) formatCount(count int, icon string, color format.Color) string {
	if count == 0 {
		return ""
	}
	return format.Colorize(color, fmt.Sprintf("%s%d", icon, count))
}
