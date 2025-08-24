package sessioninfo

import "github.com/mirage20/ccstatus-go/internal/core"

// Key is the unique identifier for the session info provider
const Key = core.ProviderKey("sessioninfo")

// GetSessionInfo is a typed getter for components to use
func GetSessionInfo(ctx *core.RenderContext) (*SessionInfo, bool) {
	return core.Get[*SessionInfo](ctx, Key)
}
