package blockusage

import (
	"fmt"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

// Component displays 5-hour block usage
type Component struct {
	priority int
}

// New creates a new block usage component
func New(priority int) *Component {
	return &Component{priority: priority}
}

// Name returns the component name
func (c *Component) Name() string {
	return "blockusage"
}

// Render generates the block usage display string
func (c *Component) Render(ctx *core.RenderContext) string {
	blockUsage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return ""
	}

	if blockUsage.TotalTokens == 0 {
		return ""
	}

	f := ctx.Formatter()

	// Determine color based on usage percentage
	color := c.getUsageColor(blockUsage.UsagePercentage)
	
	// Format the token usage part
	arrowIcon := f.Icon("arrow")
	formattedTokens := f.FormatTokens(blockUsage.TotalTokens)
	formattedPercentage := f.FormatPercentage(blockUsage.UsagePercentage)
	
	tokenPart := f.Color(color, fmt.Sprintf("%s %s %s", arrowIcon, formattedTokens, formattedPercentage))
	
	// Format the time part with end time
	clockIcon := f.Icon("clock")
	formattedTime := c.formatRemainingTime(blockUsage.RemainingMinutes, blockUsage.EndTime)
	timePart := f.Color(core.ColorGray, fmt.Sprintf("| %s %s", clockIcon, formattedTime))

	return tokenPart + " " + timePart
}

// Enabled checks if the component should be rendered
func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.blockusage.enabled", true)
}

// Priority returns the component priority
func (c *Component) Priority() int {
	return c.priority
}

// ShouldRender implements OptionalComponent for conditional display
func (c *Component) ShouldRender(ctx *core.RenderContext) bool {
	blockUsage, ok := blockusage.GetBlockUsage(ctx)
	if !ok {
		return false
	}
	return blockUsage.TotalTokens > 0
}

// getUsageColor returns color based on usage percentage
func (c *Component) getUsageColor(percentage float64) core.ColorStyle {
	switch {
	case percentage > 80:
		return core.ColorRed
	case percentage > 60:
		return core.ColorYellow
	default:
		return core.ColorGreen
	}
}

// formatRemainingTime formats the remaining time with end time like the TypeScript version
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
			return fmt.Sprintf("%s (%s)", remaining, timeStr)
		}
	}
	
	return remaining
}