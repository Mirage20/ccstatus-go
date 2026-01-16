package sevenday

import (
	"fmt"
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/ratelimit"
)

func init() {
	core.RegisterComponent("ratelimit.sevenday", New)
}

// Component displays the 7-day rate limit usage.
type Component struct {
	config *Config
}

// New is the factory function for the 7-day rate limit component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "ratelimit.sevenday", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the rate limit display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	limits, ok := ratelimit.GetRateLimits(ctx)
	if !ok || limits.SevenDay == nil {
		return ""
	}

	sevenDay := limits.SevenDay

	// Calculate remaining time
	remaining := ""
	endTime := ""
	if sevenDay.ResetsAt != nil {
		remaining = formatRemainingDays(*sevenDay.ResetsAt)
		endTime = sevenDay.ResetsAt.Local().Format(c.config.EndTimeFormat)
	}

	// Determine colors
	statusColor := c.getUsageColor(sevenDay.Utilization)
	infoColor := format.ParseColor(c.config.Color)

	// Build template data with pre-colored values
	// Icon and Utilization use status color (green/yellow/red)
	// Remaining and EndTime use info color (gray) as supplementary info
	data := map[string]interface{}{
		"Icon":        format.Colorize(statusColor, c.config.Icon),
		"Utilization": format.Colorize(statusColor, fmt.Sprintf("%.0f%%", sevenDay.Utilization)),
		"Remaining":   format.Colorize(infoColor, remaining),
		"EndTime":     format.Colorize(infoColor, endTime),
		"EndTimeRaw":  sevenDay.ResetsAt,
	}

	// Render template (values are pre-colored)
	return format.RenderTemplate(c.config.Template, data)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"ratelimit"}
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

// formatRemainingDays formats the time until reset with days support.
// Examples: "2d3h", "5h30m", "45m", "0m".
func formatRemainingDays(resetTime time.Time) string {
	const (
		minutesPerHour = 60
		hoursPerDay    = 24
		minutesPerDay  = hoursPerDay * minutesPerHour
	)

	remaining := time.Until(resetTime)
	if remaining <= 0 {
		return "0m"
	}

	totalMinutes := int(remaining.Minutes())
	days := totalMinutes / minutesPerDay
	hours := (totalMinutes % minutesPerDay) / minutesPerHour
	minutes := totalMinutes % minutesPerHour

	switch {
	case days > 0 && hours > 0:
		return fmt.Sprintf("%dd%dh", days, hours)
	case days > 0:
		return fmt.Sprintf("%dd", days)
	case hours > 0 && minutes > 0:
		return fmt.Sprintf("%dh%dm", hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dm", minutes)
	}
}
