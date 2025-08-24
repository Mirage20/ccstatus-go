package context

import (
	"fmt"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/tokenusage"
)

const (
	iconContext = "\uea7b" // Nerd Font: Context/Tokens icon
)

// Component displays session token usage.
type Component struct {
	priority int
}

// New creates a new tokens component.
func New(priority int) *Component {
	return &Component{priority: priority}
}

// Name returns the component name.
func (c *Component) Name() string {
	return "tokens"
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

	// Get context limit from config (default 200k for Claude models)
	contextLimit := ctx.Config().GetInt64("components.tokens.context_limit", 200000)
	percentage := float64(total) / float64(contextLimit) * 100

	// Determine color based on usage percentage
	color := c.getUsageColor(percentage)

	formatted := format.FormatWithUnit(total)

	return format.Colorize(color, fmt.Sprintf("%s %s", iconContext, formatted))
}

// Enabled checks if the component should be rendered.
func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.tokens.enabled", true)
}

// Priority returns the component priority.
func (c *Component) Priority() int {
	return c.priority
}

// ShouldRender implements OptionalComponent for conditional display.
func (c *Component) ShouldRender(ctx *core.RenderContext) bool {
	usage, ok := tokenusage.GetTokenUsage(ctx)
	if !ok {
		return false
	}
	return usage.Total() > 0
}

// getUsageColor returns color based on usage percentage.
func (c *Component) getUsageColor(percentage float64) format.Color {
	switch {
	case percentage > 90:
		return format.ColorRed
	case percentage > 80:
		return format.ColorYellow
	default:
		return format.ColorGreen
	}
}
