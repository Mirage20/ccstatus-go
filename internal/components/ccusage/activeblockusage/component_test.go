package activeblockusage

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

// TestRender tests the active block usage component rendering.
func TestRender(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		blockUsage *blockusage.BlockUsage
		want       string
	}{
		{
			name:       "returns empty when block usage is missing",
			config:     defaultConfig(),
			blockUsage: nil,
			want:       "",
		},
		{
			name:   "returns empty when total tokens is zero",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    0,
				MaxBlockTokens: 30000000,
			},
			want: "",
		},
		{
			name:   "renders with green color when usage is under 60%",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    250000, // 250k tokens
				MaxBlockTokens: 500000, // 500k max
			},
			want: "\033[32m\U000F0E7A 250k 50%\033[0m", // Green color
		},
		{
			name:   "renders with yellow color when usage is between 60-80%",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    350000, // 350k tokens
				MaxBlockTokens: 500000, // 500k max = 70%
			},
			want: "\033[33m\U000F0E7A 350k 70%\033[0m", // Yellow color
		},
		{
			name:   "renders with red color when usage is over 80%",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    450000, // 450k tokens
				MaxBlockTokens: 500000, // 500k max = 90%
			},
			want: "\033[31m\U000F0E7A 450k 90%\033[0m", // Red color
		},
		{
			name: "uses custom template without percentage",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}}",
				Icon:              "📊",
				BlockLimit:        500000,
				WarningThreshold:  60.0,
				CriticalThreshold: 80.0,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    100000, // 100k
				MaxBlockTokens: 0,      // Will use BlockLimit from config
			},
			want: "\033[32m📊 100k\033[0m", // Green color, no percentage
		},
		{
			name: "uses custom colors",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}} ({{printf \"%.0f\" .UsagePercentage}}%)",
				Icon:              "\U000F0E7A",
				BlockLimit:        500000,
				WarningThreshold:  60.0,
				CriticalThreshold: 80.0,
				NormalColor:       "cyan",
				WarningColor:      "magenta",
				CriticalColor:     "blue",
			},
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    150000, // 150k (30% of 500k)
				MaxBlockTokens: 0,      // Will use BlockLimit from config
			},
			want: "\033[36m\U000F0E7A 150k (30%)\033[0m", // Cyan color
		},
		{
			name: "uses custom thresholds",
			config: &Config{
				Template:          "{{.Icon}} {{.Formatted}} ({{printf \"%.0f\" .UsagePercentage}}%)",
				Icon:              "\U000F0E7A",
				BlockLimit:        500000,
				WarningThreshold:  40.0, // Lower warning threshold
				CriticalThreshold: 60.0, // Lower critical threshold
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    250000, // 250k (50% - warning with custom threshold)
				MaxBlockTokens: 0,      // Will use BlockLimit from config
			},
			want: "\033[33m\U000F0E7A 250k (50%)\033[0m", // Yellow (warning)
		},
		{
			name: "template with all variables",
			config: &Config{
				Template:          "Tokens: {{.TotalTokens}} ({{.Formatted}}) - {{printf \"%.1f\" .UsagePercentage}}% of {{.Limit}}",
				Icon:              "",
				BlockLimit:        500000,
				WarningThreshold:  60.0,
				CriticalThreshold: 80.0,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
			},
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    123456,
				MaxBlockTokens: 0, // Will use BlockLimit from config
			},
			want: "\033[32mTokens: 123456 (123k) - 24.7% of 500000\033[0m",
		},
		{
			name:   "handles edge case at exactly 60% threshold",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    300000, // 300k
				MaxBlockTokens: 500000, // 500k max = 60%
			},
			want: "\033[32m\U000F0E7A 300k 60%\033[0m", // Green (exactly at threshold)
		},
		{
			name:   "handles edge case at exactly 80% threshold",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    400000, // 400k
				MaxBlockTokens: 500000, // 500k max = 80%
			},
			want: "\033[33m\U000F0E7A 400k 80%\033[0m", // Yellow (exactly at threshold)
		},
		{
			name:   "formats small numbers without units",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    750,
				MaxBlockTokens: 500000, // 750/500k = 0.15%
			},
			want: "\033[32m\U000F0E7A 750 0%\033[0m", // Green color, no unit
		},
		{
			name:   "formats millions with M suffix",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				TotalTokens:    1500000, // 1.5M
				MaxBlockTokens: 500000,  // 1.5M/500k = 300%
			},
			want: "\033[31m\U000F0E7A 1.5M 300%\033[0m", // Red color
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			if tt.blockUsage != nil {
				ctx.Set(blockusage.Key, tt.blockUsage)
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

	if len(providers) != 1 || providers[0] != "blockusage" {
		t.Errorf("RequiredProviders() = %v, want [blockusage]", providers)
	}
}
