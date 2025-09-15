package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache implements file-based caching with session isolation.
type Cache struct {
	baseDir   string
	sessionID string
	entries   map[string]*entry
	dirty     bool
	mu        sync.RWMutex
}

// entry represents a cached item.
type entry struct {
	Data      json.RawMessage `json:"data"` // Store as raw JSON to avoid base64 encoding
	ExpiresAt time.Time       `json:"expires_at"`
	CachedAt  time.Time       `json:"cached_at"`
}

// data represents the structure of the cache file.
type data struct {
	SessionID   string            `json:"session_id"`
	LastUpdated time.Time         `json:"last_updated"`
	Providers   map[string]*entry `json:"providers"`
	Version     string            `json:"version"`
}

// NewCache creates a new file cache with session isolation.
func NewCache(baseDir, sessionID string) *Cache {
	fc := &Cache{
		baseDir:   baseDir,
		sessionID: sessionID,
		entries:   make(map[string]*entry),
	}

	// Load existing cache if available
	_ = fc.load()

	return fc
}

// load reads the entire cache file into memory.
func (fc *Cache) load() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	path := fc.getCachePath()
	fileData, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// No cache file yet, that's fine
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	var cacheData data
	if err = json.Unmarshal(fileData, &cacheData); err != nil {
		// Corrupted cache file, start fresh
		return nil
	}

	// Validate session ID matches
	if cacheData.SessionID != fc.sessionID {
		// Different session, start fresh
		return nil
	}

	// Load entries, filtering out expired ones
	now := time.Now()
	for key, e := range cacheData.Providers {
		if now.Before(e.ExpiresAt) {
			fc.entries[key] = e
		}
	}

	return nil
}

// Save writes all cache entries to disk.
func (fc *Cache) Save() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if !fc.dirty {
		return nil // Nothing to save
	}

	// Clean expired entries before saving
	now := time.Now()
	activeEntries := make(map[string]*entry)
	for key, e := range fc.entries {
		if now.Before(e.ExpiresAt) {
			activeEntries[key] = e
		}
	}

	cacheData := data{
		SessionID:   fc.sessionID,
		LastUpdated: now,
		Providers:   activeEntries,
		Version:     "2.0",
	}

	jsonData, err := json.MarshalIndent(cacheData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	path := fc.getCachePath()
	tempPath := path + ".tmp"

	// Ensure cache directory exists (lazy creation)
	if err = os.MkdirAll(fc.baseDir, 0700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write to temp file first
	if err = os.WriteFile(tempPath, jsonData, 0600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Atomic rename
	if err = os.Rename(tempPath, path); err != nil {
		// Try to clean up temp file on failure
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to rename cache file: %w", err)
	}

	// Only mark as clean after successful save
	fc.dirty = false
	return nil
}

// Get retrieves cached data from in-memory cache and unmarshals into target.
func (fc *Cache) Get(key string, target any) (bool, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	e, exists := fc.entries[key]
	if !exists {
		return false, nil
	}

	// Check if expired
	if time.Now().After(e.ExpiresAt) {
		return false, nil
	}

	// Unmarshal directly into the target
	err := json.Unmarshal(e.Data, target)
	return err == nil, err
}

// Set stores data in in-memory cache.
func (fc *Cache) Set(key string, value any, ttl time.Duration) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Marshal the value to JSON
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	fc.entries[key] = &entry{
		Data:      jsonData, // json.RawMessage is just []byte
		ExpiresAt: time.Now().Add(ttl),
		CachedAt:  time.Now(),
	}
	fc.dirty = true

	return nil
}

// Delete removes cached data.
func (fc *Cache) Delete(key string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.entries[key]; exists {
		delete(fc.entries, key)
		fc.dirty = true
	}

	return nil
}

// Cleanup removes old cache files (older than 24 hours).
func (fc *Cache) Cleanup() error {
	fc.mu.RLock()
	baseDir := fc.baseDir
	sessionID := fc.sessionID
	fc.mu.RUnlock()

	pattern := filepath.Join(baseDir, "ccstatus_*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	cutoff := time.Now().Add(-24 * time.Hour)
	currentFile := filepath.Join(baseDir, fmt.Sprintf("ccstatus_%s.json", sessionID))

	for _, f := range files {
		// Don't delete current session's cache
		if f == currentFile {
			continue
		}

		info, statErr := os.Stat(f)
		if statErr != nil {
			continue
		}

		// Remove files older than 24 hours
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(f)
		}
	}

	return nil
}

// Close performs save and cleanup operations.
// It saves any pending changes and optionally cleans up old cache files.
func (fc *Cache) Close() error {
	// Save any pending changes
	if err := fc.Save(); err != nil {
		return fmt.Errorf("failed to save cache: %w", err)
	}

	// Cleanup old cache files occasionally (10% chance based on session ID)
	if len(fc.sessionID) > 0 && fc.sessionID[len(fc.sessionID)-1]%10 == 0 {
		// Ignore cleanup errors - they shouldn't prevent close
		_ = fc.Cleanup()
	}

	return nil
}

// getCachePath generates the cache file path for this session.
func (fc *Cache) getCachePath() string {
	// Use session ID in filename for easy identification
	filename := fmt.Sprintf("ccstatus_%s.json", fc.sessionID)
	return filepath.Join(fc.baseDir, filename)
}
