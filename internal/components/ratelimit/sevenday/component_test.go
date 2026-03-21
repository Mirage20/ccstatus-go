package sevenday

import (
	"testing"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

// TestRender tests the 7-day rate limit component rendering.
func TestRender(t *testing.T) {
	// Create test reset times as unix epoch seconds
	twoDaysFromNow := time.Now().Add(2*24*time.Hour + 3*time.Hour).Unix()
	fiveHoursFromNow := time.Now().Add(5*time.Hour + 30*time.Minute).Unix()

	tests := []struct {
		name       string
		config     *Config
		rateLimits *core.SessionRateLimits
		want       string
		wantPrefix string
	}{
		{
			name:       "shows placeholder when rate limits is missing",
			config:     defaultConfig(),
			rateLimits: nil,
			want:       "\033[90m7d\033[0m \033[90m--\033[0m",
		},
		{
			name:   "shows placeholder when seven day is nil",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: nil,
			},
			want: "\033[90m7d\033[0m \033[90m--\033[0m",
		},
		{
			name:   "renders utilization with green color when under 60%",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 25.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			// Icon and utilization green, remaining gray
			wantPrefix: "\033[32m7d\033[0m \033[32m25%\033[0m \033[90m",
		},
		{
			name:   "renders utilization with yellow color when between 60-80%",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 65.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			wantPrefix: "\033[33m7d\033[0m \033[33m65%\033[0m \033[90m",
		},
		{
			name:   "renders utilization with red color when over 80%",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 85.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			wantPrefix: "\033[31m7d\033[0m \033[31m85%\033[0m \033[90m",
		},
		{
			name:   "renders remaining time with days and hours",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 20.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			wantPrefix: "\033[32m7d\033[0m \033[32m20%\033[0m \033[90m2d",
		},
		{
			name:   "renders remaining time in hours when under a day",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 20.0,
					ResetsAt:       &fiveHoursFromNow,
				},
			},
			wantPrefix: "\033[32m7d\033[0m \033[32m20%\033[0m \033[90m5h",
		},
		{
			name:   "handles nil reset time gracefully",
			config: defaultConfig(),
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 50.0,
					ResetsAt:       nil,
				},
			},
			// No trailing spaces or empty color codes when reset time is nil
			want: "\033[32m7d\033[0m \033[32m50%\033[0m",
		},
		{
			name: "uses custom icon",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}}",
				Icon:              "󰔛",
				EndTimeFormat:     "Mon 3:04 PM",
				WarningThreshold:  60,
				CriticalThreshold: 80,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "gray",
			},
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 30.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			want: "\033[32m󰔛\033[0m \033[32m30%\033[0m",
		},
		{
			name: "uses custom thresholds",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}}",
				Icon:              "7d",
				EndTimeFormat:     "Mon 3:04 PM",
				WarningThreshold:  30, // Lower threshold
				CriticalThreshold: 50, // Lower threshold
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "gray",
			},
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 40.0, // Would be yellow with custom threshold
					ResetsAt:       &twoDaysFromNow,
				},
			},
			want: "\033[33m7d\033[0m \033[33m40%\033[0m", // Yellow due to custom threshold
		},
		{
			name: "uses custom info color",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}} {{.Remaining}}",
				Icon:              "7d",
				EndTimeFormat:     "Mon 3:04 PM",
				WarningThreshold:  60,
				CriticalThreshold: 80,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "cyan", // Custom info color
			},
			rateLimits: &core.SessionRateLimits{
				SevenDay: &core.SessionRateLimit{
					UsedPercentage: 25.0,
					ResetsAt:       &twoDaysFromNow,
				},
			},
			// Remaining time should be cyan instead of gray
			wantPrefix: "\033[32m7d\033[0m \033[32m25%\033[0m \033[36m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			info := &sessioninfo.SessionInfo{
				RateLimits: tt.rateLimits,
			}
			ctx.Set(sessioninfo.Key, info)

			got := c.Render(ctx)

			if tt.wantPrefix != "" {
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

	if len(providers) != 1 || providers[0] != "sessioninfo" {
		t.Errorf("RequiredProviders() = %v, want [sessioninfo]", providers)
	}
}
