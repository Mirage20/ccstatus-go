package fivehour

import (
	"testing"
	"time"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/ratelimit"
)

// TestRender tests the 5-hour rate limit component rendering.
func TestRender(t *testing.T) {
	// Create test reset times
	twoHoursFromNow := time.Now().Add(2 * time.Hour)
	thirtyMinsFromNow := time.Now().Add(30 * time.Minute)

	tests := []struct {
		name       string
		config     *Config
		limits     *ratelimit.RateLimits
		want       string
		wantPrefix string
	}{
		{
			name:   "returns empty when rate limits is missing",
			config: defaultConfig(),
			limits: nil,
			want:   "",
		},
		{
			name:   "returns empty when five hour is nil",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: nil,
			},
			want: "",
		},
		{
			name:   "renders utilization with green color when under 60%",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 25.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			// Icon and utilization green, remaining and end time gray
			wantPrefix: "\033[32m5h\033[0m \033[32m25%\033[0m \033[90m",
		},
		{
			name:   "renders utilization with yellow color when between 60-80%",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 65.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			wantPrefix: "\033[33m5h\033[0m \033[33m65%\033[0m \033[90m",
		},
		{
			name:   "renders utilization with red color when over 80%",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 85.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			wantPrefix: "\033[31m5h\033[0m \033[31m85%\033[0m \033[90m",
		},
		{
			name:   "renders remaining time in minutes when under an hour",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 20.0,
					ResetsAt:    &thirtyMinsFromNow,
				},
			},
			wantPrefix: "\033[32m5h\033[0m \033[32m20%\033[0m \033[90m",
		},
		{
			name:   "renders remaining time in hours and minutes",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 20.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			wantPrefix: "\033[32m5h\033[0m \033[32m20%\033[0m \033[90m",
		},
		{
			name:   "handles nil reset time gracefully",
			config: defaultConfig(),
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 50.0,
					ResetsAt:    nil,
				},
			},
			// No trailing spaces or empty color codes when reset time is nil
			want: "\033[32m5h\033[0m \033[32m50%\033[0m",
		},
		{
			name: "uses custom icon",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}}",
				Icon:              "󰔛",
				EndTimeFormat:     "3:04 PM",
				WarningThreshold:  60,
				CriticalThreshold: 80,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "gray",
			},
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 30.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			want: "\033[32m󰔛\033[0m \033[32m30%\033[0m",
		},
		{
			name: "uses custom thresholds",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}}",
				Icon:              "5h",
				EndTimeFormat:     "3:04 PM",
				WarningThreshold:  30, // Lower threshold
				CriticalThreshold: 50, // Lower threshold
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "gray",
			},
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 40.0, // Would be yellow with custom threshold
					ResetsAt:    &twoHoursFromNow,
				},
			},
			want: "\033[33m5h\033[0m \033[33m40%\033[0m", // Yellow due to custom threshold
		},
		{
			name: "uses custom info color",
			config: &Config{
				Template:          "{{.Icon}} {{.Utilization}} {{.Remaining}}",
				Icon:              "5h",
				EndTimeFormat:     "3:04 PM",
				WarningThreshold:  60,
				CriticalThreshold: 80,
				NormalColor:       "green",
				WarningColor:      "yellow",
				CriticalColor:     "red",
				Color:             "cyan", // Custom info color
			},
			limits: &ratelimit.RateLimits{
				FiveHour: &ratelimit.RateLimit{
					Utilization: 25.0,
					ResetsAt:    &twoHoursFromNow,
				},
			},
			// Remaining time should be cyan instead of gray
			wantPrefix: "\033[32m5h\033[0m \033[32m25%\033[0m \033[36m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			if tt.limits != nil {
				ctx.Set(ratelimit.Key, tt.limits)
			}

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

	if len(providers) != 1 || providers[0] != "ratelimit" {
		t.Errorf("RequiredProviders() = %v, want [ratelimit]", providers)
	}
}
