package account

import "github.com/mirage20/ccstatus-go/internal/core"

// Key is the unique identifier for the account provider.
const Key = core.ProviderKey("account")

// GetAccount is a typed getter for components to use.
func GetAccount(ctx *core.RenderContext) (*Account, bool) {
	return core.Get[*Account](ctx, Key)
}
