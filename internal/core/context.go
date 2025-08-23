package core

import "sync"

// RenderContext holds all data and utilities for rendering
type RenderContext struct {
	data      map[ProviderKey]interface{}
	errors    map[ProviderKey]error
	formatter Formatter
	config    *Config
	mu        sync.RWMutex
}

// NewRenderContext creates a new render context
func NewRenderContext(config *Config, formatter Formatter) *RenderContext {
	return &RenderContext{
		data:      make(map[ProviderKey]interface{}),
		errors:    make(map[ProviderKey]error),
		formatter: formatter,
		config:    config,
	}
}

// Get retrieves typed data from a provider (generic function)
func Get[T any](ctx *RenderContext, key ProviderKey) (T, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	var zero T
	value, exists := ctx.data[key]
	if !exists {
		return zero, false
	}

	typed, ok := value.(T)
	return typed, ok
}

// Set stores provider data
func (ctx *RenderContext) Set(key ProviderKey, data interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.data[key] = data
}

// SetError stores provider error
func (ctx *RenderContext) SetError(key ProviderKey, err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.errors[key] = err
}

// GetError retrieves provider error
func (ctx *RenderContext) GetError(key ProviderKey) (error, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	err, exists := ctx.errors[key]
	return err, exists
}

// Formatter returns the formatter
func (ctx *RenderContext) Formatter() Formatter {
	return ctx.formatter
}

// Config returns the configuration
func (ctx *RenderContext) Config() *Config {
	return ctx.config
}