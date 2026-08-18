package effort

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the effort component factory
	core.RegisterComponent("effort", New)
}

// Component displays the effective reasoning effort level for the session.
type Component struct {
	config *Config
}

// New is the factory function for the effort component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "effort", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the effort display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok || info == nil {
		return ""
	}

	// Claude Code omits the whole "effort" key on models that do not support
	// effort selection, so a nil pointer means "unsupported". An empty level
	// is treated the same way defensively - the CLI never sends one.
	level := ""
	if info.Effort != nil {
		level = info.Effort.Level
	}
	if level == "" {
		if c.config.HideUnsupported {
			return ""
		}
		level = c.config.UnsupportedLabel
	}

	data := map[string]interface{}{
		"Level": level,
		"Icon":  c.iconFor(level),
	}
	result := format.RenderTemplate(c.config.Template, data)

	return format.Colorize(format.ParseColor(c.colorFor(level)), result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}

// iconFor returns the icon configured for an exact effort level, using an
// exact map lookup rather than substring matching - "high" is a substring of
// "xhigh" and would be mismatched otherwise.
func (c *Component) iconFor(level string) string {
	if icon, found := c.config.Icons[level]; found {
		return icon
	}
	return fallbackIcon
}

// colorFor returns the color name configured for an exact effort level.
func (c *Component) colorFor(level string) string {
	if color, found := c.config.Colors[level]; found {
		return color
	}
	return fallbackColor
}
