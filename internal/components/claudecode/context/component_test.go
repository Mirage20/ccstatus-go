package context

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/tokenusage"
)

// TestRender tests the context component rendering.
func TestRender(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		tokenUsage *tokenusage.TokenUsage
		want       string
	}{
		{
			name:       "returns empty when token usage is missing",
			config:     defaultConfig(),
			tokenUsage: nil,
			want:       "",
		},
		{
			name:   "returns empty when total tokens is zero",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:              0,
				OutputTokens:             0,
				CacheCreationInputTokens: 0,
				CacheReadInputTokens:     0,
			},
			want: "",
		},
		{
			name:   "renders with green color when usage is under 80%",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  50000, // 50k tokens
				OutputTokens: 10000, // 10k tokens
				// Total: 60k, which is 30% of 200k default limit
			},
			want: "\033[32m\uea7b 60k\033[0m", // Green color
		},
		{
			name:   "renders with yellow color when usage is between 80-90%",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  150000, // 150k tokens
				OutputTokens: 20000,  // 20k tokens
				// Total: 170k, which is 85% of 200k default limit
			},
			want: "\033[33m\uea7b 170k\033[0m", // Yellow color
		},
		{
			name:   "renders with red color when usage is over 90%",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  180000, // 180k tokens
				OutputTokens: 15000,  // 15k tokens
				// Total: 195k, which is 97.5% of 200k default limit
			},
			want: "\033[31m\uea7b 195k\033[0m", // Red color
		},
		{
			name:   "includes cache tokens in total",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:              20000, // 20k
				OutputTokens:             10000, // 10k
				CacheCreationInputTokens: 15000, // 15k
				CacheReadInputTokens:     5000,  // 5k
				// Total: 50k
			},
			want: "\033[32m\uea7b 50k\033[0m", // Green color (25% of 200k)
		},
		{
			name:   "formats small numbers without units",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  500,
				OutputTokens: 250,
				// Total: 750
			},
			want: "\033[32m\uea7b 750\033[0m", // Green color, no unit
		},
		{
			name:   "formats millions with M suffix",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  1500000, // 1.5M
				OutputTokens: 500000,  // 0.5M
				// Total: 2M
			},
			want: "\033[31m\uea7b 2.0M\033[0m", // Red color (1000% of 200k limit!)
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
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  100000, // 100k tokens (50% of limit)
				OutputTokens: 0,
			},
			want: "\033[32m📊 100k (50%)\033[0m", // Green with percentage
		},
		{
			name: "respects custom context limit",
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
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  85000, // 85k tokens
				OutputTokens: 0,
				// Total: 85k, which is 85% of 100k custom limit
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
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  50000, // 50k (25% - normal)
				OutputTokens: 0,
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
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  130000, // 130k (65% - warning with custom threshold)
				OutputTokens: 0,
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
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  40000, // 40k
				OutputTokens: 10000, // 10k
				// Total: 50k (25%)
			},
			want: "\033[32m🎯 50000/200000 tokens (50k, 25%)\033[0m",
		},
		{
			name:   "handles edge case at exactly 80% threshold",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  160000, // Exactly 80% of 200k
				OutputTokens: 0,
			},
			want: "\033[32m\uea7b 160k\033[0m", // Green color (exactly 80%)
		},
		{
			name:   "handles edge case at exactly 90% threshold",
			config: defaultConfig(),
			tokenUsage: &tokenusage.TokenUsage{
				InputTokens:  180000, // Exactly 90% of 200k
				OutputTokens: 0,
			},
			want: "\033[33m\uea7b 180k\033[0m", // Yellow color (exactly 90%)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			if tt.tokenUsage != nil {
				ctx.Set(tokenusage.Key, tt.tokenUsage)
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

	if len(providers) != 1 || providers[0] != "tokenusage" {
		t.Errorf("RequiredProviders() = %v, want [tokenusage]", providers)
	}
}
