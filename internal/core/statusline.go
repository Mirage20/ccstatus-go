package core

import (
	"context"
	"sort"
	"strings"
	"sync"
)

// StatusLine orchestrates providers and components.
type StatusLine struct {
	providers  []Provider
	components []Component
	cache      Cache
	config     *Config
}

// NewStatusLine creates a new status line.
func NewStatusLine(config *Config, cache Cache) *StatusLine {
	return &StatusLine{
		config: config,
		cache:  cache,
	}
}

// AddProvider registers a provider.
func (sl *StatusLine) AddProvider(p Provider) {
	sl.providers = append(sl.providers, p)
}

// AddComponent registers a component.
func (sl *StatusLine) AddComponent(c Component) {
	sl.components = append(sl.components, c)
}

// Render generates the complete status line.
func (sl *StatusLine) Render(ctx context.Context) string {
	// Create render context
	renderCtx := NewRenderContext(sl.config)

	// Gather data from all providers in parallel
	sl.gatherData(ctx, renderCtx)

	// Sort components by priority
	sort.Slice(sl.components, func(i, j int) bool {
		return sl.components[i].Priority() < sl.components[j].Priority()
	})

	// Render components
	var outputs []string
	for _, component := range sl.components {
		if !component.Enabled(sl.config) {
			continue
		}

		// Check optional condition
		if optional, ok := component.(OptionalComponent); ok {
			if !optional.ShouldRender(renderCtx) {
				continue
			}
		}

		if output := component.Render(renderCtx); output != "" {
			outputs = append(outputs, output)
		}
	}

	// Join with separator
	separator := sl.config.GetString("display.separator", " | ")
	return strings.Join(outputs, separator)
}

// gatherData fetches data from all providers in parallel.
func (sl *StatusLine) gatherData(ctx context.Context, renderCtx *RenderContext) {
	var wg sync.WaitGroup

	for _, provider := range sl.providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()

			// Just call Provide - caching is handled by CachingProvider if wrapped
			data, err := p.Provide(ctx)
			if err != nil {
				renderCtx.SetError(p.Key(), err)
				return
			}

			renderCtx.Set(p.Key(), data)
		}(provider)
	}

	wg.Wait()
}
