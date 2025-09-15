package tokenusage

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
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
			TTL: 2 * time.Second, // Changes frequently, short TTL
		},
	}
}
