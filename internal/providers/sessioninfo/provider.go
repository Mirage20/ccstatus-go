package sessioninfo

import (
	"context"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// Provider provides session information from the Claude session.
type Provider struct {
	session *core.ClaudeSession
}

// NewProvider creates a new session info provider.
func NewProvider(session *core.ClaudeSession) *Provider {
	return &Provider{
		session: session,
	}
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide returns the session information.
func (p *Provider) Provide(ctx context.Context) (interface{}, error) {
	return &SessionInfo{
		Model:       p.session.Model,
		SessionID:   p.session.SessionID,
		Version:     p.session.Version,
		Workspace:   p.session.Workspace,
		Cost:        p.session.Cost,
		Exceeds200K: p.session.Exceeds200K,
	}, nil
}
