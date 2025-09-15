package core

import (
	"context"
	"strings"
	"sync"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/format"
)

// SeparatorConfig defines configuration for the status line separator.
type SeparatorConfig struct {
	Symbol string `yaml:"symbol"`
	Color  string `yaml:"color"`
}

// StatusLine orchestrates providers and components.
type StatusLine struct {
	providers  []Provider
	components []Component
	separator  SeparatorConfig
}

// NewStatusLine creates a new status line with configuration.
func NewStatusLine(cfgReader *config.Reader) *StatusLine {
	// Load separator config with defaults
	separator := config.Get(cfgReader, "separator", SeparatorConfig{
		Symbol: " | ",
		Color:  "gray",
	})

	return &StatusLine{
		separator: separator,
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
	renderCtx := NewRenderContext()

	// Gather data from all providers in parallel
	sl.gatherData(ctx, renderCtx)

	// Render components in the order they were added (determined by layout config)
	var outputs []string
	for _, component := range sl.components {
		// Components now manage their own enabled state internally
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

	// Join with colored separator
	separatorColor := format.ParseColor(sl.separator.Color)
	coloredSeparator := format.Colorize(separatorColor, sl.separator.Symbol)
	return strings.Join(outputs, coloredSeparator)
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
