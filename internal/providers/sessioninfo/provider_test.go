package sessioninfo

import (
	"context"
	"testing"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

// TestProvideCopiesSessionModeFields guards Provide()'s hand-maintained
// field-by-field copy - the classic place for a newly added session field to
// be silently forgotten. It covers both the present case (effort, thinking
// and fast_mode all populated) and the absent case (effort/thinking omitted
// by Claude Code on older versions or unsupported models).
func TestProvideCopiesPresentSessionModeFields(t *testing.T) {
	session := &core.ClaudeSession{
		FastMode: true,
		Effort:   &core.EffortInfo{Level: "xhigh"},
		Thinking: &core.ThinkingInfo{Enabled: true},
	}
	info := provideSessionInfo(t, session)

	if !info.FastMode {
		t.Errorf("FastMode = %v, want true", info.FastMode)
	}
	if info.Effort == nil || info.Effort.Level != "xhigh" {
		t.Errorf("Effort = %+v, want Level=xhigh", info.Effort)
	}
	if info.Thinking == nil || !info.Thinking.Enabled {
		t.Errorf("Thinking = %+v, want Enabled=true", info.Thinking)
	}
}

func TestProvideCopiesAbsentSessionModeFields(t *testing.T) {
	// Zero value: no effort, no thinking, fast_mode false - the shape Claude
	// Code sends for models that don't support effort, or on older versions
	// that don't report thinking at all.
	info := provideSessionInfo(t, &core.ClaudeSession{})

	if info.FastMode {
		t.Errorf("FastMode = %v, want false", info.FastMode)
	}
	if info.Effort != nil {
		t.Errorf("Effort = %+v, want nil", info.Effort)
	}
	if info.Thinking != nil {
		t.Errorf("Thinking = %+v, want nil", info.Thinking)
	}
}

// provideSessionInfo runs Provider.Provide for the given session and
// unwraps the result, failing the test on any error or unexpected type.
func provideSessionInfo(t *testing.T, session *core.ClaudeSession) *SessionInfo {
	t.Helper()

	p := &Provider{session: session}
	got, err := p.Provide(context.Background())
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}

	info, ok := got.(*SessionInfo)
	if !ok {
		t.Fatalf("Provide() returned %T, want *SessionInfo", got)
	}
	return info
}

// TestNewReturnsProvider is a minimal smoke test - no provider_test.go
// existed for this package before.
func TestNewReturnsProvider(t *testing.T) {
	session := &core.ClaudeSession{}
	provider, _ := New(config.NewReader(""), session)

	if provider.Key() != Key {
		t.Errorf("Key() = %v, want %v", provider.Key(), Key)
	}
}
