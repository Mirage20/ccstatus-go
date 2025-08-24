package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileCache implements file-based caching with session isolation.
type FileCache struct {
	baseDir   string
	sessionID string
	entries   map[string]*CacheEntry
	dirty     bool
	mu        sync.RWMutex
}

// CacheEntry represents a cached item.
type CacheEntry struct {
	Data      json.RawMessage `json:"data"`
	ExpiresAt time.Time       `json:"expires_at"`
	CachedAt  time.Time       `json:"cached_at"`
}

// CacheFile represents the structure of the cache file.
type CacheFile struct {
	SessionID   string                 `json:"session_id"`
	LastUpdated time.Time              `json:"last_updated"`
	Providers   map[string]*CacheEntry `json:"providers"`
	Version     string                 `json:"version"`
}

// NewFileCache creates a new file cache with session isolation.
func NewFileCache(baseDir, sessionID string) *FileCache {
	// Ensure cache directory exists
	os.MkdirAll(baseDir, 0755)

	fc := &FileCache{
		baseDir:   baseDir,
		sessionID: sessionID,
		entries:   make(map[string]*CacheEntry),
	}

	// Load existing cache if available
	fc.Load()

	return fc
}

// Load reads the entire cache file into memory.
func (fc *FileCache) Load() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	path := fc.getCachePath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// No cache file yet, that's fine
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	var cacheFile CacheFile
	if err := json.Unmarshal(data, &cacheFile); err != nil {
		// Corrupted cache file, start fresh
		return nil
	}

	// Validate session ID matches
	if cacheFile.SessionID != fc.sessionID {
		// Different session, start fresh
		return nil
	}

	// Load entries, filtering out expired ones
	now := time.Now()
	for key, entry := range cacheFile.Providers {
		if now.Before(entry.ExpiresAt) {
			fc.entries[key] = entry
		}
	}

	return nil
}

// Save writes all cache entries to disk.
func (fc *FileCache) Save() error {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	if !fc.dirty {
		return nil // Nothing to save
	}

	// Clean expired entries before saving
	now := time.Now()
	activeEntries := make(map[string]*CacheEntry)
	for key, entry := range fc.entries {
		if now.Before(entry.ExpiresAt) {
			activeEntries[key] = entry
		}
	}

	cacheFile := CacheFile{
		SessionID:   fc.sessionID,
		LastUpdated: now,
		Providers:   activeEntries,
		Version:     "2.0",
	}

	data, err := json.MarshalIndent(cacheFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	path := fc.getCachePath()
	tempPath := path + ".tmp"

	// Write to temp file first
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to rename cache file: %w", err)
	}

	fc.dirty = false
	return nil
}

// Get retrieves cached data from in-memory cache.
func (fc *FileCache) Get(key string) (interface{}, bool) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	entry, exists := fc.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Return the data - could be json.RawMessage (from disk) or actual type (from Set)
	return entry.Data, true
}

// Set stores data in in-memory cache.
func (fc *FileCache) Set(key string, value interface{}, ttl time.Duration) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Marshal the value
	valueData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fc.entries[key] = &CacheEntry{
		Data:      valueData,
		ExpiresAt: time.Now().Add(ttl),
		CachedAt:  time.Now(),
	}
	fc.dirty = true

	return nil
}

// Delete removes cached data.
func (fc *FileCache) Delete(key string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.entries[key]; exists {
		delete(fc.entries, key)
		fc.dirty = true
	}

	return nil
}

// Cleanup removes old cache files (older than 24 hours).
func (fc *FileCache) Cleanup() error {
	pattern := filepath.Join(fc.baseDir, "ccstatus_*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-24 * time.Hour)
	currentFile := fc.getCachePath()

	for _, file := range files {
		// Don't delete current session's cache
		if file == currentFile {
			continue
		}

		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// Remove files older than 24 hours
		if info.ModTime().Before(cutoff) {
			os.Remove(file)
		}
	}

	return nil
}

// getCachePath generates the cache file path for this session.
func (fc *FileCache) getCachePath() string {
	// Use session ID in filename for easy identification
	filename := fmt.Sprintf("ccstatus_%s.json", fc.sessionID)
	return filepath.Join(fc.baseDir, filename)
}
