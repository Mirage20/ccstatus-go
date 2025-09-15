package duration

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the duration component factory
	core.RegisterComponent("duration", New)
}

// Component displays the session and API duration.
type Component struct {
	config *Config
}

// New is the factory function for duration component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "duration", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the duration display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	// Skip if no duration data
	if info.Cost.TotalDurationMs == 0 {
		return ""
	}

	// Format durations
	totalDuration := format.DurationMs(info.Cost.TotalDurationMs)

	var apiDuration string
	if c.config.ShowAPIDuration && info.Cost.TotalAPIDurationMs > 0 {
		apiDuration = format.DurationMs(info.Cost.TotalAPIDurationMs)
	}

	// Build template data
	data := map[string]interface{}{
		"Icon":          c.config.Icon,
		"APIIcon":       c.config.APIIcon,
		"TotalDuration": totalDuration,
		"APIDuration":   apiDuration,
		"TotalMs":       info.Cost.TotalDurationMs,
		"APIMs":         info.Cost.TotalAPIDurationMs,
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
