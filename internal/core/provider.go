package core

import (
	"context"
	"time"
)

// ProviderKey uniquely identifies a provider
type ProviderKey string

// Provider interface - minimal and focused
type Provider interface {
	// Key returns the unique identifier for this provider
	Key() ProviderKey

	// Provide fetches and returns data
	Provide(ctx context.Context) (interface{}, error)
}

// CacheableProvider is an optional interface for providers that support caching
type CacheableProvider interface {
	Provider

	// CacheTTL returns how long to cache this provider's data
	CacheTTL() time.Duration

	// CacheKey returns a unique cache key (for file-based caching)
	CacheKey() string
}