package cache

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// NullCache is a no-op cache implementation.
type NullCache struct{}

// NewNullCache creates a new null cache.
func NewNullCache() *NullCache {
	return &NullCache{}
}

// Get always returns false (cache miss).
func (c *NullCache) Get(key core.ProviderKey) (interface{}, bool) {
	return nil, false
}

// Set does nothing.
func (c *NullCache) Set(key core.ProviderKey, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete does nothing.
func (c *NullCache) Delete(key core.ProviderKey) error {
	return nil
}
