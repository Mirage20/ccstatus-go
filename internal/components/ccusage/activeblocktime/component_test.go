package activeblocktime

import (
	"testing"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/blockusage"
)

// TestRender tests the active block time component rendering.
func TestRender(t *testing.T) {
	// Create a test end time (2 hours from now)
	testEndTime := time.Now().Add(2 * time.Hour).Format(time.RFC3339)

	tests := []struct {
		name       string
		config     *Config
		blockUsage *blockusage.BlockUsage
		want       string
		wantPrefix string // For time-based tests where exact match is hard
	}{
		{
			name:       "returns empty when block usage is missing",
			config:     defaultConfig(),
			blockUsage: nil,
			want:       "",
		},
		{
			name:   "renders expired when remaining minutes is zero",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 0,
				EndTime:          "",
			},
			want: "\033[90m\uf017 0m\033[0m", // Gray color
		},
		{
			name:   "renders expired when remaining minutes is negative",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: -30,
				EndTime:          "",
			},
			want: "\033[90m\uf017 0m\033[0m", // Gray color
		},
		{
			name:   "renders minutes only when less than an hour",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 45,
				EndTime:          "",
			},
			want: "\033[90m\uf017 45m\033[0m", // Gray color
		},
		{
			name:   "renders hours only when exactly on the hour",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 120, // Exactly 2 hours
				EndTime:          "",
			},
			want: "\033[90m\uf017 2h\033[0m", // Gray color
		},
		{
			name:   "renders hours and minutes",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 150, // 2h 30m
				EndTime:          "",
			},
			want: "\033[90m\uf017 2h30m\033[0m", // Gray color
		},
		{
			name:   "renders with end time when available",
			config: defaultConfig(),
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 120,
				EndTime:          testEndTime,
			},
			wantPrefix: "\033[90m\uf017 2h ", // Will have time appended
		},
		{
			name: "uses custom template without end time",
			config: &Config{
				Template:      "{{.Icon}} {{.Remaining}}",
				Icon:          "⏰",
				Color:         "cyan",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 90,
				EndTime:          testEndTime,
			},
			want: "\033[36m⏰ 1h30m\033[0m", // Cyan color, no end time
		},
		{
			name: "uses custom expired text",
			config: &Config{
				Template:      "{{.Icon}} {{if .Expired}}finished{{else}}{{.Remaining}}{{end}}",
				Icon:          "\uf017",
				Color:         "red",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: -1,
				EndTime:          "",
			},
			want: "\033[31m\uf017 finished\033[0m", // Red color with custom expired text
		},
		{
			name: "template with all variables",
			config: &Config{
				Template:      "{{.Icon}} {{.RemainingMinutes}} minutes ({{.Remaining}}){{if .Expired}} - EXPIRED{{end}}",
				Icon:          "🕐",
				Color:         "yellow",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 75,
				EndTime:          "",
			},
			want: "\033[33m🕐 75 minutes (1h15m)\033[0m", // Yellow
		},
		{
			name: "shows expired in template",
			config: &Config{
				Template:      "{{if .Expired}}{{.Icon}} Block has expired{{else}}{{.Icon}} {{.Remaining}}{{end}}",
				Icon:          "⚠️",
				Color:         "red",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: -10,
				EndTime:          "",
			},
			want: "\033[31m⚠️ Block has expired\033[0m", // Red color
		},
		{
			name: "handles 24-hour time format",
			config: &Config{
				Template:      "{{.Icon}} {{.Remaining}}",
				Icon:          "\uf017",
				Color:         "gray",
				EndTimeFormat: "15:04", // 24-hour format
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 30,
				EndTime:          testEndTime,
			},
			wantPrefix: "\033[90m\uf017 30m", // Will have 24-hour time
		},
		{
			name: "handles single digit minutes",
			config: &Config{
				Template:      "{{.Icon}} {{.Remaining}}",
				Icon:          "\uf017",
				Color:         "gray",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 5,
				EndTime:          "",
			},
			want: "\033[90m\uf017 5m\033[0m",
		},
		{
			name: "handles very long duration",
			config: &Config{
				Template:      "{{.Icon}} {{.Remaining}}",
				Icon:          "\uf017",
				Color:         "gray",
				EndTimeFormat: "3:04 PM",
			},
			blockUsage: &blockusage.BlockUsage{
				RemainingMinutes: 305, // 5h 5m
				EndTime:          "",
			},
			want: "\033[90m\uf017 5h5m\033[0m",
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

			if tt.wantPrefix != "" {
				// For time-based tests, just check the prefix
				if len(got) < len(tt.wantPrefix) || got[:len(tt.wantPrefix)] != tt.wantPrefix {
					t.Errorf("Render() prefix = %q, want prefix %q", got, tt.wantPrefix)
				}
			} else if got != tt.want {
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
