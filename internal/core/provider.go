package core

import (
	"context"
	"sync"
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
)

// ProviderKey uniquely identifies a provider.
type ProviderKey string

// Provider interface - minimal and focused.
type Provider interface {
	// Key returns the unique identifier for this provider
	Key() ProviderKey

	// Provide fetches and returns data
	Provide(ctx context.Context) (interface{}, error)
}

// ============================================================================
// Provider Registry
// ============================================================================

// CacheConfig represents cache configuration for a provider.
type CacheConfig struct {
	TTL time.Duration `yaml:"ttl"`
}

// ProviderFactory is a function that creates a provider from config and session,
// and returns its cache configuration.
type ProviderFactory func(cfgReader *config.Reader, session *ClaudeSession) (Provider, CacheConfig)

// ProviderRegistration includes both factory and type information.
type ProviderRegistration struct {
	Factory     ProviderFactory
	NewInstance func() interface{} // Creates new instance for unmarshaling
}

// providerRegistry holds all registered provider factories.
type providerRegistry struct {
	mu            sync.RWMutex
	registrations map[string]*ProviderRegistration
}

// global providerRegistryInstance instance.
var providerRegistryInstance = &providerRegistry{
	registrations: make(map[string]*ProviderRegistration),
}

// RegisterProvider registers a provider factory with a name and type factory.
func RegisterProvider(name string, factory ProviderFactory, newInstance func() interface{}) {
	providerRegistryInstance.mu.Lock()
	defer providerRegistryInstance.mu.Unlock()
	providerRegistryInstance.registrations[name] = &ProviderRegistration{
		Factory:     factory,
		NewInstance: newInstance,
	}
}

// CreateProvider creates a provider by name using the registered factory.
func CreateProvider(name string, cfgReader *config.Reader, session *ClaudeSession, cache Cache) (Provider, bool) {
	providerRegistryInstance.mu.RLock()
	registration, exists := providerRegistryInstance.registrations[name]
	providerRegistryInstance.mu.RUnlock()

	if !exists {
		return nil, false
	}

	// Factory returns both provider and cache config
	provider, cacheConfig := registration.Factory(cfgReader, session)
	if provider == nil {
		return nil, false
	}

	// Apply caching if TTL is configured
	if cache != nil && cacheConfig.TTL > 0 {
		provider = NewCachingProvider(provider, cache, cacheConfig.TTL, registration.NewInstance)
	}

	return provider, true
}
