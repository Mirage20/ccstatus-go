package core

import (
	"context"
	"time"
)

// CachingProvider wraps any provider with caching behavior
type CachingProvider struct {
	provider Provider
	cache    Cache
	ttl      time.Duration
}

// NewCachingProvider creates a new caching wrapper for a provider
func NewCachingProvider(p Provider, cache Cache, ttl time.Duration) *CachingProvider {
	return &CachingProvider{
		provider: p,
		cache:    cache,
		ttl:      ttl,
	}
}

// Key returns the underlying provider's key
func (cp *CachingProvider) Key() ProviderKey {
	return cp.provider.Key()
}

// Provide fetches data from cache or underlying provider
func (cp *CachingProvider) Provide(ctx context.Context) (interface{}, error) {
	cacheKey := string(cp.provider.Key())

	// Try cache first
	if cp.cache != nil {
		if cached, found := cp.cache.Get(cacheKey); found {
			return cached, nil
		}
	}

	// Fetch from underlying provider
	data, err := cp.provider.Provide(ctx)
	if err != nil {
		return nil, err
	}

	// Cache the result (ignore cache errors - they shouldn't break the flow)
	if cp.cache != nil && cp.ttl > 0 {
		_ = cp.cache.Set(cacheKey, data, cp.ttl)
	}

	return data, nil
}
