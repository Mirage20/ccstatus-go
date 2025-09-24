package tokenusage

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

const (
	// Default cache TTL for tokenusage provider (changes frequently).
	defaultCacheTTL = 2 * time.Second
)

// Config holds configuration for the tokenusage provider.
type Config struct {
	// Cache configuration
	Cache core.CacheConfig `yaml:"cache"`
}

// defaultConfig returns the default configuration for tokenusage provider.
func defaultConfig() *Config {
	return &Config{
		Cache: core.CacheConfig{
			TTL: defaultCacheTTL,
		},
	}
}
