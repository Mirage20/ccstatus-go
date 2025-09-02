package core

import "sync"

// RenderContext holds all data and utilities for rendering.
type RenderContext struct {
	data   map[ProviderKey]interface{}
	errors map[ProviderKey]error
	config *Config
	mu     sync.RWMutex
}

// NewRenderContext creates a new render context.
func NewRenderContext(config *Config) *RenderContext {
	return &RenderContext{
		data:   make(map[ProviderKey]interface{}),
		errors: make(map[ProviderKey]error),
		config: config,
	}
}

// Get retrieves typed data from a provider (generic function).
func Get[T any](ctx *RenderContext, key ProviderKey) (T, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	var zero T
	value, exists := ctx.data[key]
	if !exists {
		return zero, false
	}

	// Simple type assertion - cache now handles unmarshaling
	typed, ok := value.(T)
	return typed, ok
}

// Set stores provider data.
func (ctx *RenderContext) Set(key ProviderKey, data interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.data[key] = data
}

// SetError stores provider error.
func (ctx *RenderContext) SetError(key ProviderKey, err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.errors[key] = err
}

// GetError retrieves provider error.
func (ctx *RenderContext) GetError(key ProviderKey) (error, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	err, exists := ctx.errors[key]
	return err, exists
}

// Config returns the configuration.
func (ctx *RenderContext) Config() *Config {
	return ctx.config
}
