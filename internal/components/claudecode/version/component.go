package version

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the version component factory
	core.RegisterComponent("version", New)
}

// Component displays the Claude Code version.
type Component struct {
	config *Config
}

// New is the factory function for version component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "version", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the version display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	if info.Version == "" {
		return ""
	}

	// Build template data
	data := map[string]interface{}{
		"Version": info.Version,
		"Icon":    c.config.Icon,
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
