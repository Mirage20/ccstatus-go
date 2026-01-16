package model

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRender tests the model component rendering.
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
			sessionInfo: nil, // No session info
			want:        "",
		},
		{
			name:   "renders Opus model with default config",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-opus-4-1-20250805",
					DisplayName: "Opus 4.1",
				},
			},
			want: "\033[35m\uf2db Opus\033[0m", // Magenta color + icon + "Opus"
		},
		{
			name:   "renders Sonnet model with default config",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-sonnet-3-5-20241022",
					DisplayName: "Sonnet 3.5",
				},
			},
			want: "\033[36m\uf2db Sonnet\033[0m", // Cyan color + icon + "Sonnet"
		},
		{
			name:   "renders Haiku model with default config",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-haiku-3-20240307",
					DisplayName: "Haiku 3",
				},
			},
			want: "\033[32m\uf2db Haiku\033[0m", // Green color + icon + "Haiku"
		},
		{
			name:   "extracts Opus from complex display name",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-opus-4-1-20250805",
					DisplayName: "Claude 3 Opus (Latest)",
				},
			},
			want: "\033[35m\uf2db Opus\033[0m", // Magenta color + icon + "Opus"
		},
		{
			name:   "extracts Sonnet from complex display name",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-3-5-sonnet-20241022",
					DisplayName: "Claude 3.5 Sonnet (October 2024)",
				},
			},
			want: "\033[36m\uf2db Sonnet\033[0m", // Cyan color + icon + "Sonnet"
		},
		{
			name:   "handles unknown model with short name",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "custom-model",
					DisplayName: "CustomAI",
				},
			},
			want: "\033[90m CustomAI\033[0m", // Gray color (default for unknown) + no icon + full name
		},
		{
			name:   "handles unknown model with long name",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "very-long-model-name",
					DisplayName: "This is a very long model display name that exceeds 20 characters",
				},
			},
			want: "\033[90m Claude\033[0m", // Gray color (default) + no icon + "Claude" (fallback)
		},
		{
			name:   "case insensitive matching for Opus",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "CLAUDE-OPUS",
					DisplayName: "claude opus model",
				},
			},
			want: "\033[35m\uf2db Opus\033[0m", // Magenta color + icon + "Opus"
		},
		{
			name:   "case insensitive matching for Sonnet",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "CLAUDE-SONNET",
					DisplayName: "CLAUDE SONNET MODEL",
				},
			},
			want: "\033[36m\uf2db Sonnet\033[0m", // Cyan color + icon + "Sonnet"
		},
		{
			name:   "handles edge case of exactly 20 character unknown model",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "exactly-20-chars",
					DisplayName: "12345678901234567890", // Exactly 20 chars
				},
			},
			want: "\033[90m 12345678901234567890\033[0m", // Gray + no icon + full name (not > 20)
		},
		{
			name: "uses custom template",
			config: &Config{
				Template: "[{{.ShortName}}]",
				Icons: map[string]string{
					"opus": "🎭",
				},
				Colors: map[string]string{
					"opus": "cyan",
				},
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-opus-4-1-20250805",
					DisplayName: "Opus 4.1",
				},
			},
			want: "\033[36m[Opus]\033[0m", // Cyan + custom template
		},
		{
			name: "uses custom icon and color",
			config: &Config{
				Template: "{{.Icon}} {{.ShortName}}",
				Icons: map[string]string{
					"sonnet": "🎵",
				},
				Colors: map[string]string{
					"sonnet": "yellow",
				},
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-sonnet-3-5-20241022",
					DisplayName: "Sonnet 3.5",
				},
			},
			want: "\033[33m🎵 Sonnet\033[0m", // Yellow + custom icon
		},
		{
			name: "template with all variables",
			config: &Config{
				Template: "{{.Icon}} {{.Name}} ({{.ShortName}}) [{{.ID}}]",
				Icons: map[string]string{
					"haiku": "🌸",
				},
				Colors: map[string]string{
					"haiku": "magenta",
				},
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-haiku-3-20240307",
					DisplayName: "Haiku 3",
				},
			},
			want: "\033[35m🌸 Haiku 3 (Haiku) [claude-haiku-3-20240307]\033[0m",
		},
		{
			name: "empty template returns empty string",
			config: &Config{
				Template: "",
				Colors: map[string]string{
					"opus": "cyan",
				},
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-opus-4-1-20250805",
					DisplayName: "Opus 4.1",
				},
			},
			want: "", // Empty template returns empty string (no color codes)
		},
		{
			name: "specific pattern matching for opus",
			config: &Config{
				Template: "{{.Icon}}",
				Icons: map[string]string{
					"opus": "🎭", // Specific pattern for opus
				},
				Colors: map[string]string{
					"opus": "magenta",
				},
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Model: core.ModelInfo{
					ID:          "claude-opus-4-1-20250805",
					DisplayName: "Opus 4.1",
				},
			},
			want: "\033[35m🎭\033[0m", // Magenta color with opus icon
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
