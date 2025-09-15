package changes

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRender tests the changes component rendering.
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
			name:   "returns empty when both added and removed are zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   0,
					TotalLinesRemoved: 0,
				},
			},
			want: "",
		},
		{
			name:   "renders both added and removed lines with default colors",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   42,
					TotalLinesRemoved: 17,
				},
			},
			want: "\033[32m+\033[0m\033[32m42\033[0m\033[31m-\033[0m\033[31m17\033[0m", // Green + and 42, Red - and 17
		},
		{
			name:   "renders only added lines when removed is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   100,
					TotalLinesRemoved: 0,
				},
			},
			want: "\033[32m+\033[0m\033[32m100\033[0m\033[31m-\033[0m\033[31m0\033[0m", // Green +100, Red -0
		},
		{
			name:   "renders only removed lines when added is zero",
			config: defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   0,
					TotalLinesRemoved: 50,
				},
			},
			want: "\033[32m+\033[0m\033[32m0\033[0m\033[31m-\033[0m\033[31m50\033[0m", // Green +0, Red -50
		},
		{
			name: "uses custom template with icon",
			config: &Config{
				Template:     "{{.Icon}} {{.AddedSign}}{{.Added}}/{{.RemovedSign}}{{.Removed}}",
				Icon:         "📝",
				Color:        "cyan",
				AddedSign:    "+",
				RemovedSign:  "-",
				AddedColor:   "green",
				RemovedColor: "red",
				ShowZero:     false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   10,
					TotalLinesRemoved: 5,
				},
			},
			want: "\033[36m📝\033[0m \033[32m+\033[0m\033[32m10\033[0m/\033[31m-\033[0m\033[31m5\033[0m", // Cyan icon, green +10, red -5
		},
		{
			name: "uses custom signs",
			config: &Config{
				Template:     "{{.AddedSign}}{{.Added}} {{.RemovedSign}}{{.Removed}}",
				AddedSign:    "↑",
				RemovedSign:  "↓",
				AddedColor:   "green",
				RemovedColor: "red",
				ShowZero:     false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   20,
					TotalLinesRemoved: 15,
				},
			},
			want: "\033[32m↑\033[0m\033[32m20\033[0m \033[31m↓\033[0m\033[31m15\033[0m", // Green ↑20, Red ↓15
		},
		{
			name: "uses custom colors",
			config: &Config{
				Template:     "{{.AddedSign}}{{.Added}}{{.RemovedSign}}{{.Removed}}",
				AddedSign:    "+",
				RemovedSign:  "-",
				AddedColor:   "cyan",
				RemovedColor: "yellow",
				ShowZero:     false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   30,
					TotalLinesRemoved: 10,
				},
			},
			want: "\033[36m+\033[0m\033[36m30\033[0m\033[33m-\033[0m\033[33m10\033[0m", // Cyan +30, Yellow -10
		},
		{
			name: "shows zero values when ShowZero is true",
			config: &Config{
				Template:     "{{.AddedSign}}{{.Added}}{{.RemovedSign}}{{.Removed}}",
				AddedSign:    "+",
				RemovedSign:  "-",
				AddedColor:   "green",
				RemovedColor: "red",
				ShowZero:     true,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   0,
					TotalLinesRemoved: 0,
				},
			},
			want: "\033[32m+\033[0m\033[32m0\033[0m\033[31m-\033[0m\033[31m0\033[0m", // Green +0, Red -0
		},
		{
			name: "handles large numbers",
			config: &Config{
				Template:     "Changes: {{.AddedSign}}{{.Added}}, {{.RemovedSign}}{{.Removed}}",
				AddedSign:    "+",
				RemovedSign:  "-",
				AddedColor:   "green",
				RemovedColor: "red",
				ShowZero:     false,
			},
			sessionInfo: &sessioninfo.SessionInfo{
				Cost: core.CostInfo{
					TotalLinesAdded:   12345,
					TotalLinesRemoved: 6789,
				},
			},
			want: "Changes: \033[32m+\033[0m\033[32m12345\033[0m, \033[31m-\033[0m\033[31m6789\033[0m",
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
