package ratelimit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const cacheFilename = "ratelimit.json"

// cacheEntry represents a cached rate limit response with timestamp.
type cacheEntry struct {
	Data      *RateLimits `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ErroredAt *time.Time  `json:"errored_at,omitempty"` // Set when saved after an API error
}

// globalCache manages the shared rate limit cache.
type globalCache struct {
	cacheDir string
}

// newGlobalCache creates a new global cache instance.
func newGlobalCache(cacheDir string) *globalCache {
	return &globalCache{cacheDir: cacheDir}
}

// getCachePath returns the cache file path.
func (c *globalCache) getCachePath() string {
	return filepath.Join(c.cacheDir, cacheFilename)
}

// load reads and parses the cache file.
func (c *globalCache) load() (*cacheEntry, bool) {
	path := c.getCachePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry cacheEntry
	if err = json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	return &entry, true
}

// Get retrieves cached rate limits if valid (not expired).
// Returns stale=true if the data was cached after an API error.
func (c *globalCache) Get(ttl time.Duration) (*RateLimits, bool) {
	entry, ok := c.load()
	if !ok {
		return nil, false
	}

	// Check if cache is still valid
	if time.Since(entry.Timestamp) > ttl {
		return nil, false
	}

	entry.Data.Stale = entry.ErroredAt != nil
	return entry.Data, true
}

// InErrorBackoff returns true if we recently had an API error and should
// not retry yet. Uses errorBackoffTTL as the backoff duration.
func (c *globalCache) InErrorBackoff(backoff time.Duration) bool {
	entry, ok := c.load()
	if !ok {
		return false
	}

	if entry.ErroredAt == nil {
		return false
	}

	return time.Since(*entry.ErroredAt) < backoff
}

// GetStale retrieves cached rate limits regardless of TTL expiry.
func (c *globalCache) GetStale() (*RateLimits, bool) {
	entry, ok := c.load()
	if !ok {
		return nil, false
	}

	return entry.Data, true
}

// Set stores rate limits in the cache.
func (c *globalCache) Set(data *RateLimits, erroredAt *time.Time) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(c.cacheDir, 0700); err != nil {
		return err
	}

	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
		ErroredAt: erroredAt,
	}

	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	path := c.getCachePath()
	tempPath := path + ".tmp"

	// Write to temp file first
	if err = os.WriteFile(tempPath, jsonData, 0600); err != nil {
		return err
	}

	// Atomic rename
	if err = os.Rename(tempPath, path); err != nil {
		// Clean up temp file on rename failure
		_ = os.Remove(tempPath)
		return err
	}

	return nil
}
