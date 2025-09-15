package core

import (
	"sync"

	"github.com/mirage20/ccstatus-go/internal/config"
)

// Component renders a part of the status line.
type Component interface {
	// Render generates the display string
	Render(ctx *RenderContext) string

	// RequiredProviders returns the list of provider names this component needs
	RequiredProviders() []string
}

// OptionalComponent can be conditionally displayed.
type OptionalComponent interface {
	Component

	// ShouldRender determines if component should render based on context
	ShouldRender(ctx *RenderContext) bool
}

// ============================================================================
// Component Registry
// ============================================================================

// ComponentFactory is a function that creates a component from config.
type ComponentFactory func(cfgReader *config.Reader) Component

// componentRegistry holds all registered component factories.
type componentRegistry struct {
	mu        sync.RWMutex
	factories map[string]ComponentFactory
}

// global componentRegistryInstance instance.
var componentRegistryInstance = &componentRegistry{
	factories: make(map[string]ComponentFactory),
}

// RegisterComponent registers a component factory with a name.
func RegisterComponent(name string, factory ComponentFactory) {
	componentRegistryInstance.mu.Lock()
	defer componentRegistryInstance.mu.Unlock()
	componentRegistryInstance.factories[name] = factory
}

// CreateComponent creates a component by name using the registered factory.
func CreateComponent(name string, cfgReader *config.Reader) (Component, bool) {
	componentRegistryInstance.mu.RLock()
	defer componentRegistryInstance.mu.RUnlock()

	factory, exists := componentRegistryInstance.factories[name]
	if !exists {
		return nil, false
	}

	// Factory knows its own config path
	return factory(cfgReader), true
}
