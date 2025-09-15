package activeblockusage

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

func init() {
	// Register the active block usage component factory
	core.RegisterComponent("activeblockusage", New)
}

// Component displays the active 5-hour block token usage.
type Component struct {
	config *Config
}

// New is the factory function for active block usage component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "activeblockusage", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the block usage display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	usage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return ""
	}

	// Skip if no tokens used
	if usage.TotalTokens == 0 {
		return ""
	}

	// Format token count
	formatted := format.WithUnit(usage.TotalTokens)

	// Determine which limit to use: user config takes precedence
	var limit int64
	if c.config.BlockLimit > 0 {
		// User has explicitly configured a limit
		limit = c.config.BlockLimit
	} else {
		// Use dynamic max from historical data
		limit = usage.MaxBlockTokens
	}

	// Calculate usage percentage based on the chosen limit
	var usagePercentage float64
	if limit > 0 {
		usagePercentage = float64(usage.TotalTokens) / float64(limit) * 100
	}

	// Build template data
	data := map[string]interface{}{
		"Icon":            c.config.Icon,
		"TotalTokens":     usage.TotalTokens,
		"Formatted":       formatted,
		"UsagePercentage": usagePercentage,
		"Limit":           limit,
		"MaxBlockTokens":  usage.MaxBlockTokens, // Also expose the dynamic max if users want it
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Determine color based on usage percentage
	color := c.getUsageColor(usagePercentage)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"blockusage"}
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
