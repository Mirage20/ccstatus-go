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

// Get retrieves cached rate limits if valid (not expired).
func (c *globalCache) Get(ttl time.Duration) (*RateLimits, bool) {
	path := c.getCachePath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry cacheEntry
	if err = json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check if cache is still valid
	if time.Since(entry.Timestamp) > ttl {
		return nil, false
	}

	return entry.Data, true
}

// Set stores rate limits in the cache.
func (c *globalCache) Set(data *RateLimits) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(c.cacheDir, 0700); err != nil {
		return err
	}

	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
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
