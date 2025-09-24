package model

import (
	"strings"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

const (
	// Maximum length for display name before falling back to "Claude".
	maxDisplayNameLength = 20
)

func init() {
	// Register the model component factory
	core.RegisterComponent("model", New)
}

// Component displays the Claude model information.
type Component struct {
	config *Config
}

// New is the factory function for model component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "model", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the model display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	sessionInfo, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	// Build template data
	data := map[string]interface{}{
		"ID":        sessionInfo.Model.ID,
		"Name":      sessionInfo.Model.DisplayName,
		"ShortName": c.getShortName(sessionInfo),
		"Icon":      c.getIcon(sessionInfo.Model.ID),
	}

	// Render template
	result := format.RenderTemplate(c.config.Template, data)

	// Apply color to the output
	colorName := c.getColorName(sessionInfo.Model.ID)
	color := format.ParseColor(colorName) // Always returns a valid color (gray if unknown)
	return format.Colorize(color, result)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}

// getShortName returns the short name for the model.
func (c *Component) getShortName(info *sessioninfo.SessionInfo) string {
	// Extract from display name (case-insensitive)
	displayName := info.Model.DisplayName
	displayNameLower := strings.ToLower(displayName)

	switch {
	case strings.Contains(displayNameLower, "opus"):
		return "Opus"
	case strings.Contains(displayNameLower, "sonnet"):
		return "Sonnet"
	case strings.Contains(displayNameLower, "haiku"):
		return "Haiku"
	default:
		// If unknown, return a shortened version
		if len(displayName) > maxDisplayNameLength {
			return "Claude"
		}
		return displayName
	}
}

// matchPattern finds the first pattern that matches the modelID.
func (c *Component) matchPattern(modelID string, patterns map[string]string) string {
	modelLower := strings.ToLower(modelID)
	for pattern, value := range patterns {
		if strings.Contains(modelLower, strings.ToLower(pattern)) {
			return value
		}
	}
	return "" // No match found
}

// getIcon returns the icon for the model based on config patterns.
func (c *Component) getIcon(modelID string) string {
	return c.matchPattern(modelID, c.config.Icons)
}

// getColorName returns the color name for the model based on config patterns.
func (c *Component) getColorName(modelID string) string {
	return c.matchPattern(modelID, c.config.Colors)
}
