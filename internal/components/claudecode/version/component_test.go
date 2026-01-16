package version

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRender tests the version component rendering.
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
			name:   "renders version with default config",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "1.0.89",
			},
			want: "\033[90mv1.0.89\033[0m", // Gray color + v1.0.89 (no icon)
		},
		{
			name:   "returns empty when version is empty string",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "",
			},
			want: "",
		},
		{
			name:   "handles version with special characters",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "2.0.0-beta.1",
			},
			want: "\033[90mv2.0.0-beta.1\033[0m", // Gray color + v2.0.0-beta.1 (no icon)
		},
		{
			name: "uses custom template",
			config: &Config{
				Template: "Version: {{.Version}}",
				Icon:     "📦",
				Color:    "cyan",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "3.0.0",
			},
			want: "\033[36mVersion: 3.0.0\033[0m", // Cyan + custom template
		},
		{
			name: "uses custom icon in template",
			config: &Config{
				Template: "{{.Icon}} {{.Version}}",
				Icon:     "🏷️",
				Color:    "yellow",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "1.2.3",
			},
			want: "\033[33m🏷️ 1.2.3\033[0m", // Yellow + custom icon
		},
		{
			name: "empty template returns empty string",
			config: &Config{
				Template: "",
				Icon:     "\uf8d5",
				Color:    "gray",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "1.0.0",
			},
			want: "", // Empty template returns empty string (no color codes)
		},
		{
			name: "template with only icon",
			config: &Config{
				Template: "{{.Icon}}",
				Icon:     "⚙️",
				Color:    "magenta",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Version: "1.0.0",
			},
			want: "\033[35m⚙️\033[0m", // Magenta + just icon
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

// TestRequiredProviders tests that the component declares its dependencies.
func TestRequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "sessioninfo" {
		t.Errorf("RequiredProviders() = %v, want [sessioninfo]", providers)
	}
}
