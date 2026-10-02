package thinking

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
			name:   "returns empty when CLI does not report thinking, even with hide_when_enabled false",
			config: &Config{Template: "{{.Icon}}", OffIcon: "z", OffColor: "gray", HideWhenEnabled: false},
			sessionInfo: &sessioninfo.SessionInfo{
				Thinking: nil,
			},
			want: "",
		},
		{
			name:        "hides bubble by default when thinking is enabled",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: true}},
			want:        "",
		},
		{
			name:        "renders comment_slash icon in gray when disabled",
			config:      defaultConfig(),
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: false}},
			want:        "\033[90m\ued96\033[0m",
		},
		{
			name: "renders enabled icon when hide_when_enabled is false",
			config: &Config{
				Template:        "{{.Icon}}",
				Icon:            "\uf075",
				OffIcon:         "\ued96",
				Color:           "cyan",
				OffColor:        "gray",
				HideWhenEnabled: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: true}},
			want:        "\033[36m\uf075\033[0m",
		},
		{
			name: "renders Enabled boolean in a conditional template",
			config: &Config{
				Template:        "{{if .Enabled}}on{{else}}off{{end}}",
				Color:           "cyan",
				OffColor:        "gray",
				HideWhenEnabled: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: true}},
			want:        "\033[36mon\033[0m",
		},
		{
			name: "empty template returns empty",
			config: &Config{
				Template:        "",
				OffColor:        "gray",
				HideWhenEnabled: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: false}},
			want:        "",
		},
		{
			// An empty icon hides the component on its own, without the hide
			// flag - Colorize returns "" for empty text and statusline drops
			// empty output. Lets a user surface only the off state.
			name: "empty icon hides the on state even with hide_when_enabled false",
			config: &Config{
				Template:        "{{.Icon}}",
				Icon:            "",
				OffIcon:         "\uf420",
				OffColor:        "red",
				HideWhenEnabled: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: true}},
			want:        "",
		},
		{
			// ...while the off state still renders, so only it is visible.
			name: "empty icon still renders the off state",
			config: &Config{
				Template:        "{{.Icon}}",
				Icon:            "",
				OffIcon:         "\uf420",
				OffColor:        "red",
				HideWhenEnabled: false,
			},
			sessionInfo: &sessioninfo.SessionInfo{Thinking: &core.ThinkingInfo{Enabled: false}},
			want:        "\033[31m\uf420\033[0m",
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
