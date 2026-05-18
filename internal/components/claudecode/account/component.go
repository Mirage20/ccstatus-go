package account

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	accountprovider "github.com/mirage20/ccstatus-go/internal/providers/account"
)

func init() {
	core.RegisterComponent("account", New)
}

// Component displays which Claude Code account is currently active.
type Component struct {
	config *Config
}

// New is the factory function for the account component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "account", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the account label.
func (c *Component) Render(ctx *core.RenderContext) string {
	acc, ok := accountprovider.GetAccount(ctx)
	if !ok || acc == nil {
		return ""
	}

	profile := acc.Profile
	if profile == "" {
		return ""
	}
	if profile == "unknown" && c.config.HideUnknown {
		return ""
	}

	data := map[string]interface{}{
		"Profile": profile,
		"Icon":    c.config.Icon,
	}
	result := format.RenderTemplate(c.config.Template, data)

	colorName, ok := c.config.Colors[profile]
	if !ok {
		colorName = c.config.DefaultColor
	}
	return format.Colorize(format.ParseColor(colorName), result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"account"}
}
