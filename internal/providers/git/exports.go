package git

import "github.com/mirage20/ccstatus-go/internal/core"

// Key is the unique identifier for the git provider.
const Key = core.ProviderKey("git")

// GetInfo is a typed getter for components to use.
func GetInfo(ctx *core.RenderContext) (*Info, bool) {
	return core.Get[*Info](ctx, Key)
}
