package ratelimit

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

const (
	// Default cache TTL for rate limits.
	defaultCacheTTL = 5 * time.Minute

	// Default backoff duration after an API error before retrying.
	defaultErrorBackoff = 5 * time.Minute

	// Default cache directory.
	defaultCacheSubdir = "ccstatus"
)

// Config holds configuration for the ratelimit provider.
type Config struct {
	// TTL for the global rate limit cache (not the per-session cache).
	TTL time.Duration `yaml:"ttl"`

	// ErrorBackoff is how long to wait before retrying after an API error (e.g. 429).
	// During this period, stale cached data is served without making API calls.
	ErrorBackoff time.Duration `yaml:"error_backoff"`

	// CacheDir is the directory for the global cache.
	// Defaults to ~/.cache/ccstatus
	CacheDir string `yaml:"cache_dir"`

	// Cache configuration for the CachingProvider wrapper.
	// Set to 0 TTL to bypass the session-based CachingProvider.
	Cache core.CacheConfig `yaml:"cache"`
}

// defaultConfig returns the default configuration for ratelimit provider.
func defaultConfig() *Config {
	return &Config{
		TTL:          defaultCacheTTL,
		ErrorBackoff: defaultErrorBackoff,
		CacheDir:     defaultCacheDir(),
		Cache: core.CacheConfig{
			TTL: 0, // Bypass CachingProvider - we use our own global cache
		},
	}
}

// defaultCacheDir returns the default cache directory path.
func defaultCacheDir() string {
	// Try XDG_CACHE_HOME first
	if cacheHome := os.Getenv("XDG_CACHE_HOME"); cacheHome != "" {
		return filepath.Join(cacheHome, defaultCacheSubdir)
	}

	// Fall back to ~/.cache
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/tmp", defaultCacheSubdir)
	}

	return filepath.Join(homeDir, ".cache", defaultCacheSubdir)
}
