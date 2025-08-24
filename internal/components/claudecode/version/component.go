package version

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

	if info.Version == "" {
		return ""
	}

	return format.Dimmed(fmt.Sprintf("v%s", info.Version))
}

func (c *Component) Name() string {
	return "version"
}

func (c *Component) Enabled(config *core.Config) bool {
	return config.GetBool("components.version.enabled", true)
}

func (c *Component) Priority() int {
	return c.priority
}
