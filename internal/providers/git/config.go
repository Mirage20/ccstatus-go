package git

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

const (
	// Default cache TTL for git operations.
	defaultCacheTTL = 10 * time.Second
)

// Config defines configuration for the git provider.
type Config struct {
	// Cache configuration
	Cache core.CacheConfig `yaml:"cache"`
}

// defaultConfig returns the default configuration for git provider.
func defaultConfig() *Config {
	return &Config{
		Cache: core.CacheConfig{
			TTL: defaultCacheTTL,
		},
	}
}
