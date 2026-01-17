package newline

import (
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	core.RegisterComponent("newline", New)
}

// Component outputs a newline for multi-line status layouts.
type Component struct{}

// New is the factory function for newline component.
func New(_ *config.Reader) core.Component {
	return &Component{}
}

// Render returns a newline character.
func (c *Component) Render(_ *core.RenderContext) string {
	return "\n"
}

// RequiredProviders returns no providers since newline needs no data.
func (c *Component) RequiredProviders() []string {
	return nil
}
