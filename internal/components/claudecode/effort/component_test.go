package effort

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		sessionInfo *sessioninfo.SessionInfo
		want        string
	}{
		{
			name:        "returns empty when session info is missing",
			config:      defaultConfig(),
			sessionInfo: nil,
			want:        "",
		},
		{
			name:        "hides when model does not support effort",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: nil},
			want:        "",
		},
		{
			name:        "hides when effort level is empty",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: ""}},
			want:        "",
		},
		{
			name: "renders unsupported label when hide_unsupported is false",
			config: &Config{
				Template:         "{{.Icon}} {{.Level}}",
				Icons:            defaultConfig().Icons,
				Colors:           defaultConfig().Colors,
				UnsupportedLabel: "--",
				HideUnsupported:  false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Effort: nil},
			want:        "\033[90m\uf420 --\033[0m",
		},
		{
			name:        "renders low with cloud_moon icon in gray",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "low"}},
			want:        "\033[90m\ueeef low\033[0m",
		},
		{
			name:        "renders medium with walking icon in cyan",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "medium"}},
			want:        "\033[36m\uee1d medium\033[0m",
		},
		{
			name:        "renders high with running icon in yellow",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "high"}},
			want:        "\033[33m\uef0c high\033[0m",
		},
		{
			// Guards the substring hazard: "high" is a substring of "xhigh".
			// Must use exact map lookup, never model's matchPattern.
			name:        "renders xhigh with brain icon in magenta, not high's running/yellow",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "xhigh"}},
			want:        "\033[35m\uee9c xhigh\033[0m",
		},
		{
			name:        "renders max with fire icon in red",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "max"}},
			want:        "\033[31m\uf06d max\033[0m",
		},
		{
			name:        "renders unknown future level with defaults instead of hiding",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "ultra"}},
			want:        "\033[90m\uf420 ultra\033[0m",
		},
		{
			name:        "level lookup is exact, not case-insensitive",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "HIGH"}},
			want:        "\033[90m\uf420 HIGH\033[0m",
		},
		{
			name: "empty template returns empty",
			config: &Config{
				Template: "",
				Icons:    defaultConfig().Icons,
				Colors:   defaultConfig().Colors,
			},
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "max"}},
			want:        "",
		},
		{
			name: "custom icon-only template",
			config: &Config{
				Template: "{{.Icon}}",
				Icons:    defaultConfig().Icons,
				Colors:   defaultConfig().Colors,
			},
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "max"}},
			want:        "\033[31m\uf06d\033[0m",
		},
		{
			name: "unmapped level falls back to the hardcoded fallback color, not a config value",
			config: &Config{
				Template: "{{.Level}}",
				Icons:    map[string]string{},
				Colors:   map[string]string{"max": "red"},
			},
			sessionInfo: &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "low"}},
			want:        "\033[90mlow\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			if tt.sessionInfo != nil {
				ctx.Set(sessioninfo.Key, tt.sessionInfo)
			}

			got := c.Render(ctx)
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestRenderIsDeterministic guards against a future refactor to substring
// matching over a map, which would make xhigh flip between running and brain.
func TestRenderIsDeterministic(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()
	ctx.Set(sessioninfo.Key, &sessioninfo.SessionInfo{Effort: &core.EffortInfo{Level: "xhigh"}})

	want := c.Render(ctx)
	for range 100 {
		if got := c.Render(ctx); got != want {
			t.Fatalf("Render() = %q, want %q (non-deterministic level lookup)", got, want)
		}
	}
}

func TestRequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "sessioninfo" {
		t.Errorf("RequiredProviders() = %v, want [sessioninfo]", providers)
	}
}
