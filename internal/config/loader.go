package config

import (
	"github.com/mirage20/ccstatus-go/internal/core"
	"os"
	"path/filepath"
)

// Load loads configuration from various sources
func Load() (*core.Config, error) {
	config := core.NewConfig()

	// Set default values
	setDefaults(config)

	// Load from environment variables (if any)
	loadFromEnv(config)

	return config, nil
}

// Default returns a configuration with default values
func Default() *core.Config {
	config := core.NewConfig()
	setDefaults(config)
	return config
}

// setDefaults sets default configuration values
func setDefaults(config *core.Config) {
	// Display settings
	config.Set("display.separator", " | ")

	// Component settings
	config.Set("components.model.enabled", true)
	config.Set("components.tokens.enabled", true)
	config.Set("components.tokens.context_limit", int64(200000))
	config.Set("components.blockusage.enabled", true)
	config.Set("components.git.enabled", false)

	// Provider settings
	config.Set("providers.git.enabled", false)

	// Cache settings
	config.Set("cache.enabled", true)
	config.Set("cache.dir", getCacheDir())
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv(config *core.Config) {
	// Example: override cache directory from env
	if cacheDir := os.Getenv("CCSTATUS_CACHE_DIR"); cacheDir != "" {
		config.Set("cache.dir", cacheDir)
	}

	// Example: disable cache from env
	if os.Getenv("CCSTATUS_NO_CACHE") == "1" {
		config.Set("cache.enabled", false)
	}
}

// getCacheDir returns the cache directory path
func getCacheDir() string {
	// Priority order:
	// 1. Session-specific cache (if available)
	if sessionID := os.Getenv("CLAUDE_SESSION_ID"); sessionID != "" {
		return filepath.Join(os.TempDir(), "ccstatus", sessionID)
	}

	// 2. User cache directory
	if cacheHome := os.Getenv("XDG_CACHE_HOME"); cacheHome != "" {
		return filepath.Join(cacheHome, "ccstatus")
	}

	// 3. Home directory cache
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".cache", "ccstatus")
	}

	// 4. Fallback to temp
	return filepath.Join(os.TempDir(), "ccstatus", "global")
}