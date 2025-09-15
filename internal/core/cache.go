package core

import "time"

// Cache provides caching functionality with value-based storage.
type Cache interface {
	// Get retrieves cached data and unmarshals into target
	// target must be a pointer to the desired type
	Get(key string, target any) (bool, error)

	// Set stores data in cache
	Set(key string, value any, ttl time.Duration) error

	// Delete removes cached data
	Delete(key string) error

	// Close performs any necessary cleanup and persistence
	// Should be called when done using the cache
	Close() error
}
