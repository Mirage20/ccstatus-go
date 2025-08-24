package activeblock

import (
	"fmt"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

const (
	iconArrow = "\U000F0E7A" // Nerd Font: Up-down arrow icon
	iconClock = "\uf017"     // Nerd Font: Clock icon
)

// Component displays 5-hour block usage.
type Component struct {
	priority int
}

// New creates a new block usage component.
func New(priority int) *Component {
	return &Component{priority: priority}
}

// Name returns the component name.
func (c *Component) Name() string {
	return "blockusage"
}

// Render generates the block usage display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	blockUsage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return ""
	}

	if blockUsage.TotalTokens == 0 {
		return ""
	}

	// Determine color based on usage percentage
	color := c.getUsageColor(blockUsage.UsagePercentage)

	// Format the token usage part
	formattedTokens := format.FormatWithUnit(blockUsage.TotalTokens)
	formattedPercentage := format.FormatPercentage(blockUsage.UsagePercentage)

	tokenPart := format.Colorize(color, fmt.Sprintf("%s %s %s", iconArrow, formattedTokens, formattedPercentage))

	// Format the time part with end time
	formattedTime := c.formatRemainingTime(blockUsage.RemainingMinutes, blockUsage.EndTime)
	timePart := format.Colorize(format.ColorGray, fmt.Sprintf("| %s %s", iconClock, formattedTime))

	return tokenPart + " " + timePart
}

// Enabled checks if the component should be rendered.
func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.blockusage.enabled", true)
}

// Priority returns the component priority.
func (c *Component) Priority() int {
	return c.priority
}

// ShouldRender implements OptionalComponent for conditional display.
func (c *Component) ShouldRender(ctx *core.RenderContext) bool {
	blockUsage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return false
	}
	return blockUsage.TotalTokens > 0
}

// getUsageColor returns color based on usage percentage.
func (c *Component) getUsageColor(percentage float64) format.Color {
	switch {
	case percentage > 80:
		return format.ColorRed
	case percentage > 60:
		return format.ColorYellow
	default:
		return format.ColorGreen
	}
}

// formatRemainingTime formats the remaining time with end time like the TypeScript version.
func (c *Component) formatRemainingTime(minutes int, endTimeStr string) string {
	if minutes <= 0 {
		return "expired"
	}

	// Format remaining time
	var remaining string
	if minutes < 60 {
		remaining = fmt.Sprintf("%dm", minutes)
	} else {
		hours := minutes / 60
		mins := minutes % 60
		if mins == 0 {
			remaining = fmt.Sprintf("%dh", hours)
		} else {
			remaining = fmt.Sprintf("%dh%dm", hours, mins)
		}
	}

	// Parse and format end time if available
	if endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err == nil {
			// Format as "HH:MM AM/PM" in local time
			hour := endTime.Local().Hour()
			minute := endTime.Local().Minute()
			ampm := "AM"
			displayHour := hour

			if hour >= 12 {
				ampm = "PM"
				if hour > 12 {
					displayHour = hour - 12
				}
			} else if hour == 0 {
				displayHour = 12
			}

			timeStr := fmt.Sprintf("%d:%02d %s", displayHour, minute, ampm)
			return fmt.Sprintf("%s %s", remaining, timeStr)
		}
	}

	return remaining
}
