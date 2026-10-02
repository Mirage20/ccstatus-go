package fastmode

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/providers/sessioninfo"
)

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
			name:        "hides when fast_mode is absent from stdin (decodes as false)",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{},
			want:        "",
		},
		{
			name:        "hides by default when fast mode is off",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{FastMode: false},
			want:        "",
		},
		{
			name:        "renders flash icon in yellow when active",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{FastMode: true},
			want:        "\033[33m\uf0e7\033[0m",
		},
		{
			name: "renders off icon when hide_when_off is false",
			config: &Config{
				Template:    "{{.Icon}}",
				Icon:        "\uf0e7",
				OffIcon:     "z",
				Color:       "yellow",
				OffColor:    "gray",
				HideWhenOff: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{FastMode: false},
			want:        "\033[90mz\033[0m",
		},
		{
			name: "renders Enabled boolean in a conditional template",
			config: &Config{
				Template:    "{{if .Enabled}}fast{{else}}slow{{end}}",
				Color:       "yellow",
				OffColor:    "gray",
				HideWhenOff: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{FastMode: true},
			want:        "\033[33mfast\033[0m",
		},
		{
			name: "empty template returns empty",
			config: &Config{
				Template: "",
				Color:    "yellow",
			},
			sessionInfo: &sessioninfo.SessionInfo{FastMode: true},
			want:        "",
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

func TestRequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "sessioninfo" {
		t.Errorf("RequiredProviders() = %v, want [sessioninfo]", providers)
	}
}
