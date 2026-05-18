package account

import (
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
)

// Config holds configuration for the account provider.
type Config struct {
	// Cache configuration
	Cache core.CacheConfig `yaml:"cache"`

	// ActiveSlot is the Keychain service name holding the currently-active
	// OAuth credentials. Claude Code writes to "Claude Code-credentials".
	ActiveSlot string `yaml:"active_slot"`

	// Profiles maps profile name -> Keychain service name for that profile's
	// saved token snapshot.
	Profiles map[string]string `yaml:"profiles"`
}

// defaultConfig returns the default configuration for the account provider.
func defaultConfig() *Config {
	return &Config{
		Cache: core.CacheConfig{
			TTL: 60 * time.Second,
		},
		ActiveSlot: "Claude Code-credentials",
		Profiles: map[string]string{
			"personal": "Claude Code-personal",
			"work":     "Claude Code-work",
		},
	}
}
