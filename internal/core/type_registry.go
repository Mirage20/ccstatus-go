package core

import "sync"

// TypeFactory creates a new instance of a provider's data type.
type TypeFactory func() interface{}

// TypeRegistry manages type factories for creating provider data types.
type TypeRegistry struct {
	factories map[ProviderKey]TypeFactory
	mu        sync.RWMutex
}

// NewTypeRegistry creates a new type registry.
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		factories: make(map[ProviderKey]TypeFactory),
	}
}

// Register registers a type factory for a provider key.
func (tr *TypeRegistry) Register(key ProviderKey, factory TypeFactory) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.factories[key] = factory
}

// CreateInstance creates a new instance for the given provider key.
func (tr *TypeRegistry) CreateInstance(key ProviderKey) (interface{}, bool) {
	tr.mu.RLock()
	defer tr.mu.RUnlock()

	factory, exists := tr.factories[key]
	if !exists {
		return nil, false
	}
	return factory(), true
}
