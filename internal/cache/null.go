package cache

import "time"

// NullCache is a no-op cache implementation.
type NullCache struct{}

// NewNullCache creates a new null cache.
func NewNullCache() *NullCache {
	return &NullCache{}
}

// Get always returns false (cache miss).
func (c *NullCache) Get(key string) (interface{}, bool) {
	return nil, false
}

// Set does nothing.
func (c *NullCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

// Delete does nothing.
func (c *NullCache) Delete(key string) error {
	return nil
}
