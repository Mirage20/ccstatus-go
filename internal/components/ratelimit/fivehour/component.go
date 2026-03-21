package fivehour

import (
	"fmt"
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	core.RegisterComponent("ratelimit.fivehour", New)
}

// Component displays the 5-hour rate limit usage.
type Component struct {
	config *Config
}

// New is the factory function for the 5-hour rate limit component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "ratelimit.fivehour", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the rate limit display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	infoColor := format.ParseColor(c.config.Color)

	// When rate limit data is not available, show placeholder
	if info.RateLimits == nil || info.RateLimits.FiveHour == nil {
		data := map[string]interface{}{
			"Icon":        format.Colorize(infoColor, c.config.Icon),
			"Utilization": format.Colorize(infoColor, "--"),
			"Remaining":   "",
			"EndTime":     "",
			"EndTimeRaw":  (*time.Time)(nil),
		}
		return format.RenderTemplate(c.config.Template, data)
	}

	fiveHour := info.RateLimits.FiveHour

	// Convert unix epoch to time.Time for remaining/end time calculation
	var resetsAt *time.Time
	if fiveHour.ResetsAt != nil {
		t := time.Unix(*fiveHour.ResetsAt, 0)
		resetsAt = &t
	}

	// Calculate remaining time
	remaining := ""
	endTime := ""
	if resetsAt != nil {
		remaining = formatRemaining(*resetsAt)
		endTime = resetsAt.Local().Format(c.config.EndTimeFormat)
	}

	// Determine colors
	statusColor := c.getUsageColor(fiveHour.UsedPercentage)

	// Build template data with pre-colored values
	// Icon and Utilization use status color (green/yellow/red)
	// Remaining and EndTime use info color (gray) as supplementary info
	data := map[string]interface{}{
		"Icon":        format.Colorize(statusColor, c.config.Icon),
		"Utilization": format.Colorize(statusColor, fmt.Sprintf("%.0f%%", fiveHour.UsedPercentage)),
		"Remaining":   format.Colorize(infoColor, remaining),
		"EndTime":     format.Colorize(infoColor, endTime),
		"EndTimeRaw":  resetsAt,
	}

	// Render template (values are pre-colored)
	return format.RenderTemplate(c.config.Template, data)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}

// getUsageColor returns color based on utilization and configured thresholds.
func (c *Component) getUsageColor(utilization float64) format.Color {
	switch {
	case utilization >= c.config.CriticalThreshold:
		return format.ParseColor(c.config.CriticalColor)
	case utilization >= c.config.WarningThreshold:
		return format.ParseColor(c.config.WarningColor)
	default:
		return format.ParseColor(c.config.NormalColor)
	}
}

// formatRemaining formats the time until reset as a human-readable duration.
func formatRemaining(resetTime time.Time) string {
	remaining := time.Until(resetTime)
	if remaining <= 0 {
		return "0m"
	}

	minutes := int(remaining.Minutes())
	return format.DurationMinutes(minutes)
}
