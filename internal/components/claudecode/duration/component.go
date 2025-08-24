package duration

import (
	"fmt"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

const (
	iconDuration    = "\uF520"     // Nerd Font: Stopwatch icon
	iconAPIDuration = "\U000F109B" // Nerd Font: API icon
)

type Component struct {
	priority int
}

func New(priority int) *Component {
	return &Component{
		priority: priority,
	}
}

func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	// Skip if no duration data
	if info.Cost.TotalDurationMs == 0 {
		return ""
	}

	totalDuration := format.Duration(info.Cost.TotalDurationMs)

	output := fmt.Sprintf("%s %s", iconDuration, totalDuration)

	// If we have API duration, show both
	if info.Cost.TotalAPIDurationMs > 0 {
		apiDuration := format.Duration(info.Cost.TotalAPIDurationMs)
		output = fmt.Sprintf("%s %s %s", iconDuration, totalDuration, apiDuration)
	}

	return format.Dimmed(output)
}

func (c *Component) Name() string {
	return "duration"
}

func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.duration.enabled", true)
}

func (c *Component) Priority() int {
	return c.priority
}
