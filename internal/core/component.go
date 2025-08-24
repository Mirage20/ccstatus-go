package core

// Component renders a part of the status line
type Component interface {
	// Name returns the component identifier
	Name() string

	// Render generates the display string
	Render(ctx *RenderContext) string

	// Enabled checks if component should be rendered
	Enabled(config *Config) bool

	// Priority determines rendering order (lower = earlier)
	Priority() int
}

// OptionalComponent can be conditionally displayed
type OptionalComponent interface {
	Component

	// ShouldRender determines if component should render based on context
	ShouldRender(ctx *RenderContext) bool
}
