package changes

import (
	"strconv"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func init() {
	// Register the changes component factory
	core.RegisterComponent("changes", New)
}

// Component displays the line changes (added/removed).
type Component struct {
	config *Config
}

// New is the factory function for changes component.
func New(cfgReader *config.Reader) core.Component {
	cfg := config.GetComponent(cfgReader, "changes", defaultConfig())
	return &Component{
		config: cfg,
	}
}

// Render generates the changes display string.
func (c *Component) Render(ctx *core.RenderContext) string {
	info, ok := sessioninfo.GetSessionInfo(ctx)
	if !ok {
		return ""
	}

	// Skip if both are zero and ShowZero is false
	if !c.config.ShowZero && info.Cost.TotalLinesAdded == 0 && info.Cost.TotalLinesRemoved == 0 {
		return ""
	}

	// Parse colors
	addedColor := format.ParseColor(c.config.AddedColor)
	removedColor := format.ParseColor(c.config.RemovedColor)
	componentColor := format.ParseColor(c.config.Color)

	// Build template data with pre-colored values
	data := map[string]interface{}{
		"Icon":        format.Colorize(componentColor, c.config.Icon),
		"Added":       format.Colorize(addedColor, strconv.Itoa(info.Cost.TotalLinesAdded)),
		"Removed":     format.Colorize(removedColor, strconv.Itoa(info.Cost.TotalLinesRemoved)),
		"AddedSign":   format.Colorize(addedColor, c.config.AddedSign),
		"RemovedSign": format.Colorize(removedColor, c.config.RemovedSign),
	}

	// Render template
	return format.RenderTemplate(c.config.Template, data)
}

// RequiredProviders returns the list of provider names this component needs.
func (c *Component) RequiredProviders() []string {
	return []string{"sessioninfo"}
}
