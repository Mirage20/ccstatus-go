// Package ratelimit provides rate limit data from Anthropic OAuth API.
//
// Inspired by:
//   - https://github.com/rz1989s/claude-code-statusline
//   - https://github.com/uppinote20/claude-dashboard
package ratelimit

import (
	"context"
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

func init() {
	// Self-register with type factory
	core.RegisterProvider(string(Key), New, func() interface{} {
		return &RateLimits{}
	})
}

// Provider provides rate limit data from Anthropic OAuth API.
// This provider manages its own global cache instead of using
// the per-session CachingProvider.
type Provider struct {
	cache *globalCache
	ttl   time.Duration
}

// New creates a new rate limit provider with config.
func New(cfgReader *config.Reader, _ *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	// Load provider config with defaults
	cfg := config.GetProvider(cfgReader, "ratelimit", defaultConfig())

	return &Provider{
		cache: newGlobalCache(cfg.CacheDir),
		ttl:   cfg.TTL,
	}, cfg.Cache // TTL: 0 to bypass CachingProvider
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// Provide fetches rate limits from API or cache.
func (p *Provider) Provide(ctx context.Context) (interface{}, error) {
	// Check global cache first
	if cached, found := p.cache.Get(p.ttl); found {
		return cached, nil
	}

	// Get OAuth token
	token, err := getOAuthToken(ctx)
	if err != nil || token == "" {
		// No token available - return empty rate limits (graceful degradation)
		return &RateLimits{}, nil //nolint:nilerr // Intentional graceful degradation
	}

	// Fetch from API
	limits, err := fetchRateLimits(ctx, token)
	if err != nil {
		// API error - return empty rate limits (graceful degradation)
		return &RateLimits{}, nil //nolint:nilerr // Intentional graceful degradation
	}

	// Save to global cache (ignore errors - cache is best-effort)
	_ = p.cache.Set(limits)

	return limits, nil
}

// Key is the provider key.
const Key = core.ProviderKey("ratelimit")

// GetRateLimits is a typed getter for components.
func GetRateLimits(ctx *core.RenderContext) (*RateLimits, bool) {
	return core.Get[*RateLimits](ctx, Key)
}
