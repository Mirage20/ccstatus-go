package null

import (
	"time"
)

// Cache is a no-op cache implementation.
type Cache struct{}

// NewCache creates a new null cache.
func NewCache() *Cache {
	return &Cache{}
}

// Get always returns false (cache miss).
func (c *Cache) Get(_ string, _ any) (bool, error) {
	return false, nil
}

// Set does nothing.
func (c *Cache) Set(_ string, _ any, _ time.Duration) error {
	return nil
}

// Delete does nothing.
func (c *Cache) Delete(_ string) error {
	return nil
}

// Close does nothing (no-op implementation).
func (c *Cache) Close() error {
	return nil
}
