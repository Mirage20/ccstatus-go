package cache

import (
	"os"

	"github.com/mirage20/ccstatus-go/internal/cache/file"
	"github.com/mirage20/ccstatus-go/internal/cache/null"
	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

// New creates a cache instance based on configuration.
// Returns NullCache if cache.enabled is false, otherwise returns FileCache.
// Default behavior is to enable cache.
func New(cfg *config.Reader, sessionID string) core.Cache {
	// Default true - cache enabled unless explicitly disabled
	if !config.Get(cfg, "cache.enabled", true) {
		return null.NewCache()
	}

	// Use file cache with configured or default directory
	dir := config.Get(cfg, "cache.dir", os.TempDir())
	return file.NewCache(dir, sessionID)
}
