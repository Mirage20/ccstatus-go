package context

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
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
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	cw := info.ContextWindow

	// Determine context limit: use dynamic size from session, fallback to config
	contextLimit := cw.ContextWindowSize
	if contextLimit == 0 {
		contextLimit = c.config.ContextLimit
	}

	// Calculate total context usage including OutputTokens
	// OutputTokens are included because they become part of the conversation
	// history sent to the next API call
	// When current_usage is nil (session just started), total is 0
	var total int64
	if cw.CurrentUsage != nil {
		usage := cw.CurrentUsage
		total = usage.InputTokens + usage.OutputTokens + usage.CacheCreationInputTokens + usage.CacheReadInputTokens
	}

	// Calculate percentage
	percentage := float64(total) / float64(contextLimit) * 100

	// Format token count
	formatted := format.WithUnit(total)

	// Build template data
	data := map[string]interface{}{
		"Icon":       c.config.Icon,
		"Total":      total,
		"Formatted":  formatted,
		"Percentage": percentage, // Raw float for template formatting
		"Limit":      contextLimit,
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Determine color based on usage percentage
	color := c.getUsageColor(percentage)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
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
