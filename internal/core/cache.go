package core

import "time"

// Cache provides caching functionality.
type Cache interface {
	// Get retrieves cached data
	Get(key ProviderKey) (interface{}, bool)

	// Set stores data in cache
	Set(key ProviderKey, value interface{}, ttl time.Duration) error

	// Delete removes cached data
	Delete(key ProviderKey) error
}
