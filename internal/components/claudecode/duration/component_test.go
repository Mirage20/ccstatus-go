package duration

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRender tests the duration component rendering.
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
			name:   "returns empty when total duration is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    0,
					TotalAPIDurationMs: 0,
				},
			},
			want: "",
		},
		{
			name:   "renders total duration only when API duration is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    45000, // 45 seconds
					TotalAPIDurationMs: 0,
				},
			},
			want: "\033[90m\uF520 45s\033[0m", // Gray color + icon + 45s
		},
		{
			name:   "renders both total and API duration",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    120000, // 2 minutes
					TotalAPIDurationMs: 30000,  // 30 seconds
				},
			},
			want: "\033[90m\uF520 2m \U000F1616 30s\033[0m", // Gray color + icon + 2m + API icon + 30s
		},
		{
			name: "uses custom template with only total duration",
			config: &Config{
				Template:        "Duration: {{.TotalDuration}}",
				Icon:            "⏱️",
				Color:           "cyan",
				ShowAPIDuration: true,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    30000, // 30 seconds
					TotalAPIDurationMs: 15000, // 15 seconds (but not shown in template)
				},
			},
			want: "\033[36mDuration: 30s\033[0m", // Cyan + custom template
		},
		{
			name: "uses custom template with API duration",
			config: &Config{
				Template:        "{{.Icon}} Total: {{.TotalDuration}} | API: {{.APIDuration}}",
				Icon:            "⏱️",
				APIIcon:         "🔌",
				Color:           "yellow",
				ShowAPIDuration: true,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    90000, // 1m 30s
					TotalAPIDurationMs: 45000, // 45s
				},
			},
			want: "\033[33m⏱️ Total: 1m30s | API: 45s\033[0m", // Yellow + custom template
		},
		{
			name: "hides API duration when ShowAPIDuration is false",
			config: &Config{
				Template:        "{{.Icon}} {{.TotalDuration}}{{if .APIDuration}} {{.APIDuration}}{{end}}",
				Icon:            "\uF520",
				Color:           "gray",
				ShowAPIDuration: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    60000, // 1 minute
					TotalAPIDurationMs: 30000, // 30 seconds (should be hidden)
				},
			},
			want: "\033[90m\uF520 1m\033[0m", // Gray + only total duration
		},
		{
			name: "handles milliseconds",
			config: &Config{
				Template:        "{{.Icon}} {{.TotalDuration}}",
				Icon:            "⏱️",
				Color:           "gray",
				ShowAPIDuration: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs: 500, // 500ms rounds down to 0s
				},
			},
			want: "\033[90m⏱️ 0s\033[0m", // Gray + 0s (rounds down from 500ms)
		},
		{
			name: "handles hours",
			config: &Config{
				Template: "{{.Icon}} {{.TotalDuration}}",
				Icon:     "⏱️",
				Color:    "gray",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs: 7200000, // 2 hours
				},
			},
			want: "\033[90m⏱️ 2h\033[0m", // Gray + 2h
		},
		{
			name: "template with raw milliseconds",
			config: &Config{
				Template: "{{.TotalMs}}ms total, {{.APIMs}}ms API",
				Color:    "magenta",
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    12345,
					TotalAPIDurationMs: 6789,
				},
			},
			want: "\033[35m12345ms total, 6789ms API\033[0m", // Magenta + raw values
		},
		{
			name:   "handles long durations with minutes and seconds",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    185000, // 3m5s
					TotalAPIDurationMs: 65000,  // 1m5s
				},
			},
			want: "\033[90m\uF520 3m5s \U000F1616 1m5s\033[0m", // Gray + 3m5s + API icon + 1m5s
		},
		{
			name:   "handles exact minute durations",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalDurationMs:    60000,  // Exactly 1 minute
					TotalAPIDurationMs: 120000, // Exactly 2 minutes
				},
			},
			want: "\033[90m\uF520 1m \U000F1616 2m\033[0m", // Gray + 1m + API icon + 2m
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
