package blockusage

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// Config holds configuration for the blockusage provider.
type Config struct {
	// Cache configuration
	Cache core.CacheConfig `yaml:"cache"`
}

// defaultConfig returns the default configuration for blockusage provider.
func defaultConfig() *Config {
	return &Config{
		Cache: core.CacheConfig{
			TTL: 10 * time.Second, // Expensive ccusage command, cache for 10 seconds
		},
	}
}
