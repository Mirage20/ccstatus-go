package model

import (
	"fmt"
	"strings"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

const (
	iconModel = "\uf2db" // Nerd Font: Microchip icon
)

// Component displays the Claude model information.
type Component struct {
	priority int
}

// New creates a new model component.
func New(priority int) *Component {
	return &Component{priority: priority}
}

// Name returns the component name.
func (c *Component) Name() string {
	return "model"
}

// Render generates the model display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	sessionInfo, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok || sessionInfo.Model.DisplayName == "" {
		return ""
	}

	// Extract model name - always shorten to just the model type
	modelName := c.extractModelName(sessionInfo.Model.DisplayName)

	// Determine color based on model
	color := c.getModelColor(sessionInfo.Model.DisplayName)

	return format.Colorize(color, fmt.Sprintf("%s %s", iconModel, modelName))
}

// Enabled checks if the component should be rendered.
func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.model.enabled", true)
}

// Priority returns the component priority.
func (c *Component) Priority() int {
	return c.priority
}

// extractModelName extracts the short model name from display name.
func (c *Component) extractModelName(displayName string) string {
	switch {
	case strings.Contains(displayName, "Opus"):
		return "Opus"
	case strings.Contains(displayName, "Sonnet"):
		return "Sonnet"
	case strings.Contains(displayName, "Haiku"):
		return "Haiku"
	default:
		// If unknown, return a shortened version
		if len(displayName) > 20 {
			return "Claude"
		}
		return displayName
	}
}

// getModelColor returns the appropriate color for the model.
func (c *Component) getModelColor(displayName string) format.Color {
	switch {
	case strings.Contains(displayName, "Opus"):
		return format.ColorMagenta
	case strings.Contains(displayName, "Sonnet"):
		return format.ColorCyan
	case strings.Contains(displayName, "Haiku"):
		return format.ColorGreen
	default:
		return format.ColorYellow
	}
}
