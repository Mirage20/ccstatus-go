package status

import (
	"fmt"
	"strings"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func init() {
	core.RegisterComponent("git.status", New)
}

// Component displays git working tree status (staged, modified, untracked).
type Component struct {
	config *Config
}

// New is the factory function for git.status component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "git.status", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the git status display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := gitprovider.GetInfo(ctx)
	if !ok || !info.IsRepo {
		return ""
	}

	// If all counts are zero, return empty (clean working tree)
	if info.Staged == 0 && info.Modified == 0 && info.Untracked == 0 && info.Conflicts == 0 {
		return ""
	}

	// Parse colors
	stagedColor := format.ParseColor(c.config.StagedColor)
	modifiedColor := format.ParseColor(c.config.ModifiedColor)
	untrackedColor := format.ParseColor(c.config.UntrackedColor)
	conflictColor := format.ParseColor(c.config.ConflictColor)

	// Build pre-colored indicators
	staged := c.formatCount(info.Staged, c.config.StagedIcon, stagedColor)
	modified := c.formatCount(info.Modified, c.config.ModifiedIcon, modifiedColor)
	untracked := c.formatCount(info.Untracked, c.config.UntrackedIcon, untrackedColor)
	conflicts := c.formatCount(info.Conflicts, c.config.ConflictIcon, conflictColor)

	// Build template data with pre-colored values
	data := map[string]any{
		"Staged":    staged,
		"Modified":  modified,
		"Untracked": untracked,
		"Conflicts": conflicts,
	}

	// Render template (values are pre-colored) and trim leading space
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
