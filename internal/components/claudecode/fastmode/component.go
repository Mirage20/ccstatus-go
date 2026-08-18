package fastmode

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the fastmode component factory
	core.RegisterComponent("fastmode", New)
}

// Component displays whether fast mode is active for the session.
type Component struct {
	config *Config
}

// New is the factory function for the fastmode component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "fastmode", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the fast mode display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok || info == nil {
		return ""
	}

	enabled := info.FastMode
	if !enabled && c.config.HideWhenOff {
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
