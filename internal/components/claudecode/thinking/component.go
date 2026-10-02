package thinking

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the thinking component factory
	core.RegisterComponent("thinking", New)
}

// Component displays whether extended thinking is enabled for the session.
type Component struct {
	config *Config
}

// New is the factory function for the thinking component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "thinking", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the thinking display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok || info == nil || info.Thinking == nil {
		// Older Claude Code versions do not report thinking at all - render
		// nothing rather than assuming it is off.
		return ""
	}

	enabled := info.Thinking.Enabled
	if enabled && c.config.HideWhenEnabled {
		return ""
	}

	icon, colorName := c.config.Icon, c.config.Color
	if !enabled {
		icon, colorName = c.config.OffIcon, c.config.OffColor
	}

	data := map[string]interface{}{
		"Icon":    icon,
		"Enabled": enabled,
	}
	result := format.RenderTemplate(c.config.Template, data)

	return format.Colorize(format.ParseColor(colorName), result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}
