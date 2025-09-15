package activeblocktime

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

func init() {
	// Register the active block time component factory
	core.RegisterComponent("activeblocktime", New)
}

// Component displays the remaining time for the active 5-hour block.
type Component struct {
	config *Config
}

// New is the factory function for active block time component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "activeblocktime", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the block time display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	usage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return ""
	}

	// Format remaining time using the format package
	remaining := format.DurationMinutes(usage.RemainingMinutes)

	// Format end time if available
	endTime := ""
	if usage.EndTime != "" {
		if t, err := time.Parse(time.RFC3339, usage.EndTime); err == nil {
			endTime = t.Local().Format(c.config.EndTimeFormat)
		}
	}

	// Build template data
	data := map[string]interface{}{
		"Icon":             c.config.Icon,
		"RemainingMinutes": usage.RemainingMinutes,
		"Remaining":        remaining,
		"EndTime":          endTime,
		"EndTimeRaw":       usage.EndTime,
		"Expired":          usage.RemainingMinutes <= 0,
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Apply color
	color := format.ParseColor(c.config.Color)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"blockusage"}
}
