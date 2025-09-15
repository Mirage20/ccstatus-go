package sessioninfo

import (
	"context"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	// Self-register with type factory
	core.RegisterProvider(string(Key), New, func() interface{} {
		return &SessionInfo{}
	})
}

// Provider provides session information from the Claude session.
type Provider struct {
	session *core.ClaudeSession
}

// New creates a new session info provider with config.
func New(cfgReader *config.Reader, session *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	// Load provider config with defaults
	cfg := config.GetProvider(cfgReader, "sessioninfo", defaultConfig())

	return &Provider{
		session: session,
	}, cfg.Cache
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide returns the session information.
func (p *Provider) Provide(_ context.Context) (interface{}, error) {
	return &SessionInfo{
		Model:       p.session.Model,
		SessionID:   p.session.SessionID,
		Version:     p.session.Version,
		Workspace:   p.session.Workspace,
		Cost:        p.session.Cost,
		Exceeds200K: p.session.Exceeds200K,
	}, nil
}
