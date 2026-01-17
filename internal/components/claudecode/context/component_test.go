package context

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRequiredProviders tests that the component declares its dependencies.
func TestRequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "sessioninfo" {
		t.Errorf("RequiredProviders() = %v, want [sessioninfo]", providers)
	}
}

// TestRender tests the context component rendering.
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
			name:   "shows zero when current_usage is nil (session start)",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage:      nil,
				},
			},
			want: "\033[32m\uea7b 0\033[0m", // Green, 0 tokens
		},
		{
			name:   "shows zero when total tokens is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens:              0,
						OutputTokens:             0,
						CacheCreationInputTokens: 0,
						CacheReadInputTokens:     0,
					},
				},
			},
			want: "\033[32m\uea7b 0\033[0m", // Green, 0 tokens
		},
		{
			name:   "renders with green color when usage is under 80%",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 60000, // 60k tokens (30% of 200k)
					},
				},
			},
			want: "\033[32m\uea7b 60k\033[0m", // Green color
		},
		{
			name:   "renders with yellow color when usage is between 60-75%",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 130000, // 130k tokens (65% of 200k)
					},
				},
			},
			want: "\033[33m\uea7b 130k\033[0m", // Yellow color
		},
		{
			name:   "renders with red color when usage is over 75%",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 160000, // 160k tokens (80% of 200k)
					},
				},
			},
			want: "\033[31m\uea7b 160k\033[0m", // Red color
		},
		{
			name:   "includes all tokens in total (input + output + cache)",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens:              20000, // 20k
						OutputTokens:             10000, // 10k - included
						CacheCreationInputTokens: 10000, // 10k
						CacheReadInputTokens:     10000, // 10k
						// Total: 20k + 10k + 10k + 10k = 50k (25% of 200k)
					},
				},
			},
			want: "\033[32m\uea7b 50k\033[0m", // Green color
		},
		{
			name:   "uses dynamic context_window_size from session",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 100000, // 100k limit from session
					CurrentUsage: &core.ContextUsage{
						InputTokens: 70000, // 70k tokens (70% of 100k)
					},
				},
			},
			want: "\033[33m\uea7b 70k\033[0m", // Yellow color (70% of dynamic 100k)
		},
		{
			name:   "falls back to config limit when context_window_size is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 0, // No size from session
					CurrentUsage: &core.ContextUsage{
						InputTokens: 60000, // 60k tokens (30% of 200k fallback)
					},
				},
			},
			want: "\033[32m\uea7b 60k\033[0m", // Green color (30% of 200k fallback)
		},
		{
			name: "uses custom template with percentage",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}} ({{.Percentage}}%)",
				Icon:              "📊",
				ContextLimit:      200000,
				WarningThreshold:  80.0,
				CriticalThreshold: 90.0,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 100000, // 100k tokens (50% of limit)
					},
				},
			},
			want: "\033[32m📊 100k (50%)\033[0m", // Green with percentage
		},
		{
			name: "respects custom context limit as fallback",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}}",
				Icon:              "\uea7b",
				ContextLimit:      100000, // 100k limit
				WarningThreshold:  80.0,
				CriticalThreshold: 90.0,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 0, // Triggers fallback to config limit
					CurrentUsage: &core.ContextUsage{
						InputTokens: 85000, // 85k tokens (85% of 100k custom limit)
					},
				},
			},
			want: "\033[33m\uea7b 85k\033[0m", // Yellow color (85% of 100k)
		},
		{
			name: "uses custom colors",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}}",
				Icon:              "\uea7b",
				ContextLimit:      200000,
				WarningThreshold:  80.0,
				CriticalThreshold: 90.0,
				NormalColor:       "cyan",
				WarningColor:      "magenta",
				CriticalColor:     "blue",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 50000, // 50k (25% - normal)
					},
				},
			},
			want: "\033[36m\uea7b 50k\033[0m", // Cyan color
		},
		{
			name: "uses custom thresholds",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}}",
				Icon:              "\uea7b",
				ContextLimit:      200000,
				WarningThreshold:  60.0, // Lower warning threshold
				CriticalThreshold: 70.0, // Lower critical threshold
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 130000, // 130k (65% - warning with custom threshold)
					},
				},
			},
			want: "\033[33m\uea7b 130k\033[0m", // Yellow (warning)
		},
		{
			name: "template with all variables",
			config: &Config{
				Template:          "{{.Icon}} {{.Total}}/{{.Limit}} tokens ({{.Formatted}}, {{.Percentage}}%)",
				Icon:              "🎯",
				ContextLimit:      200000,
				WarningThreshold:  80.0,
				CriticalThreshold: 90.0,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens:  40000, // 40k
						OutputTokens: 10000, // 10k
						// Total: 50k (25%)
					},
				},
			},
			want: "\033[32m🎯 50000/200000 tokens (50k, 25%)\033[0m",
		},
		{
			name:   "handles edge case at exactly 60% threshold",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 120000, // Exactly 60% of 200k
					},
				},
			},
			want: "\033[32m\uea7b 120k\033[0m", // Green color (exactly 60% is not > 60%)
		},
		{
			name:   "handles edge case at exactly 75% threshold",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				ContextWindow: core.ContextWindow{
					ContextWindowSize: 200000,
					CurrentUsage: &core.ContextUsage{
						InputTokens: 150000, // Exactly 75% of 200k
					},
				},
			},
			want: "\033[33m\uea7b 150k\033[0m", // Yellow color (exactly 75% is not > 75%)
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
