package account

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	core.RegisterProvider(string(Key), New, func() interface{} {
		return &Account{}
	})
}

// Provider determines which Claude Code account is currently active by
// comparing the live Keychain credentials against saved profile snapshots.
type Provider struct {
	activeSlot string
	profiles   map[string]string
}

// New creates a new account provider with config.
func New(cfgReader *config.Reader, _ *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	cfg := config.GetProvider(cfgReader, "account", defaultConfig())
	return &Provider{
		activeSlot: cfg.ActiveSlot,
		profiles:   cfg.Profiles,
	}, cfg.Cache
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide returns the currently active account profile.
func (p *Provider) Provide(ctx context.Context) (interface{}, error) {
	// Keychain is macOS-only; on other platforms there's nothing to detect.
	if runtime.GOOS != "darwin" {
		return &Account{Profile: ""}, nil
	}

	active, ok := readKeychain(ctx, p.activeSlot)
	if !ok || active == "" {
		return &Account{Profile: ""}, nil
	}

	for name, slot := range p.profiles {
		saved, found := readKeychain(ctx, slot)
		if !found {
			continue
		}
		if saved == active {
			return &Account{Profile: name}, nil
		}
	}

	return &Account{Profile: "unknown"}, nil
}

// readKeychain shells out to `security` to read a generic password by service
// name. Returns the secret and whether the read succeeded.
func readKeychain(ctx context.Context, service string) (string, bool) {
	if service == "" {
		return "", false
	}
	out, err := exec.CommandContext(ctx, "security", "find-generic-password", "-s", service, "-w").Output()
	if err != nil {
		return "", false
	}
	return strings.TrimRight(string(out), "\n"), true
}
