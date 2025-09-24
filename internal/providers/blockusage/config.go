package blockusage

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

const (
	// Default cache TTL for blockusage provider (expensive ccusage command).
	defaultCacheTTL = 10 * time.Second
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
			TTL: defaultCacheTTL,
		},
	}
}
