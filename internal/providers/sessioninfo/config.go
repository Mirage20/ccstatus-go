package sessioninfo

import (
	"github.com/mirage20/ccstatus-go/internal/core"
)

// Config holds configuration for the sessioninfo provider.
type Config struct {
	// Cache configuration
	Cache core.CacheConfig `yaml:"cache"`

	// Provider-specific config would go here
	// (sessioninfo doesn't need any currently)
}

// defaultConfig returns the default configuration for sessioninfo provider.
func defaultConfig() *Config {
	return &Config{
		Cache: core.CacheConfig{
			TTL: 0, // No caching - sessioninfo is cheap, just returns struct fields
		},
	}
}
