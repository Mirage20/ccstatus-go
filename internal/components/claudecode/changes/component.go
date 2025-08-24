package changes

import (
	"fmt"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
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

	// Skip if both are zero
	if info.Cost.TotalLinesAdded == 0 && info.Cost.TotalLinesRemoved == 0 {
		return ""
	}

	added := format.Green(fmt.Sprintf("+%d", info.Cost.TotalLinesAdded))
	removed := format.Red(fmt.Sprintf("-%d", info.Cost.TotalLinesRemoved))

	return fmt.Sprintf("%s%s", added, removed)
}

func (c *Component) Name() string {
	return "changes"
}

func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.changes.enabled", true)
}

func (c *Component) Priority() int {
	return c.priority
}
