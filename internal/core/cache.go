package core

import "time"

// Cache provides caching functionality.
type Cache interface {
	// Get retrieves cached data
	Get(key string) (interface{}, bool)

	// Set stores data in cache
	Set(key string, value interface{}, ttl time.Duration) error

	// Delete removes cached data
	Delete(key string) error
}
