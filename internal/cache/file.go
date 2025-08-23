package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// FileCache implements file-based caching
type FileCache struct {
	baseDir string
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Data      json.RawMessage `json:"data"`
	ExpiresAt time.Time       `json:"expires_at"`
	Version   string          `json:"version"`
}

// NewFileCache creates a new file cache
func NewFileCache(baseDir string) *FileCache {
	// Ensure cache directory exists
	os.MkdirAll(baseDir, 0755)
	return &FileCache{baseDir: baseDir}
}

// Get retrieves cached data
func (fc *FileCache) Get(key string) (interface{}, bool) {
	path := fc.getCachePath(key)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, false
	}
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(path) // Clean up expired entry
		return nil, false
	}

	// Return raw JSON for the caller to unmarshal
	return entry.Data, true
}

// Set stores data in cache
func (fc *FileCache) Set(key string, value interface{}, ttl time.Duration) error {
	// Marshal the value
	valueData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	entry := CacheEntry{
		Data:      valueData,
		ExpiresAt: time.Now().Add(ttl),
		Version:   "1.0",
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path := fc.getCachePath(key)
	tempPath := path + ".tmp"

	// Write to temp file first
	if err := os.WriteFile(tempPath, entryData, 0644); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tempPath, path)
}

// Delete removes cached data
func (fc *FileCache) Delete(key string) error {
	path := fc.getCachePath(key)
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// getCachePath generates a cache file path from key
func (fc *FileCache) getCachePath(key string) string {
	hash := sha256.Sum256([]byte(key))
	filename := hex.EncodeToString(hash[:8]) + ".json"
	return filepath.Join(fc.baseDir, filename)
}