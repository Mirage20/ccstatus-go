package context

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/tokenusage"
)

func init() {
	// Register the context component factory
	core.RegisterComponent("context", New)
}

// Component displays session token usage.
type Component struct {
	config *Config
}

// New is the factory function for context component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "context", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the token usage display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	usage, ok := tokenusage.GetTokenUsage(ctx)
	if !ok {
		return ""
	}

	total := usage.Total()
	if total == 0 {
		return ""
	}

	// Calculate percentage
	percentage := float64(total) / float64(c.config.ContextLimit) * 100

	// Format token count
	formatted := format.WithUnit(total)

	// Build template data
	data := map[string]interface{}{
		"Icon":       c.config.Icon,
		"Total":      total,
		"Formatted":  formatted,
		"Percentage": percentage, // Raw float for template formatting
		"Limit":      c.config.ContextLimit,
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Determine color based on usage percentage
	color := c.getUsageColor(percentage)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"tokenusage"}
}

// getUsageColor returns color based on usage percentage and configured thresholds.
func (c *Component) getUsageColor(percentage float64) format.Color {
	switch {
	case percentage > c.config.CriticalThreshold:
		return format.ParseColor(c.config.CriticalColor)
	case percentage > c.config.WarningThreshold:
		return format.ParseColor(c.config.WarningColor)
	default:
		return format.ParseColor(c.config.NormalColor)
	}
}
