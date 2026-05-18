package account

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	accountprovider "github.com/mirage20/ccstatus-go/internal/providers/account"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		account *accountprovider.Account
		want    string
	}{
		{
			name:    "returns empty when account info is missing",
			config:  defaultConfig(),
			account: nil,
			want:    "",
		},
		{
			name:    "returns empty when profile is empty",
			config:  defaultConfig(),
			account: &accountprovider.Account{Profile: ""},
			want:    "",
		},
		{
			name:    "renders personal profile with default config",
			config:  defaultConfig(),
			account: &accountprovider.Account{Profile: "personal"},
			want:    "\033[36mpersonal\033[0m", // cyan
		},
		{
			name:    "renders work profile with default config",
			config:  defaultConfig(),
			account: &accountprovider.Account{Profile: "work"},
			want:    "\033[35mwork\033[0m", // magenta
		},
		{
			name:    "hides unknown by default",
			config:  defaultConfig(),
			account: &accountprovider.Account{Profile: "unknown"},
			want:    "",
		},
		{
			name: "renders unknown when hide_unknown is false",
			config: &Config{
				Template:     "{{.Icon}}{{.Profile}}",
				Colors:       map[string]string{"personal": "cyan", "work": "magenta"},
				DefaultColor: "gray",
				HideUnknown:  false,
			},
			account: &accountprovider.Account{Profile: "unknown"},
			want:    "\033[90munknown\033[0m", // gray default
		},
		{
			name: "falls back to default color for unmapped profile",
			config: &Config{
				Template:     "{{.Profile}}",
				Colors:       map[string]string{"personal": "cyan"},
				DefaultColor: "yellow",
				HideUnknown:  false,
			},
			account: &accountprovider.Account{Profile: "staging"},
			want:    "\033[33mstaging\033[0m", // yellow
		},
		{
			name: "renders custom icon and template",
			config: &Config{
				Template:     "{{.Icon}} {{.Profile}}",
				Icon:         "👤",
				Colors:       map[string]string{"personal": "green"},
				DefaultColor: "gray",
				HideUnknown:  true,
			},
			account: &accountprovider.Account{Profile: "personal"},
			want:    "\033[32m👤 personal\033[0m",
		},
		{
			name: "empty template returns empty",
			config: &Config{
				Template:     "",
				Colors:       map[string]string{"personal": "cyan"},
				DefaultColor: "gray",
				HideUnknown:  true,
			},
			account: &accountprovider.Account{Profile: "personal"},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Component{config: tt.config}
			ctx := core.NewRenderContext()

			if tt.account != nil {
				ctx.Set(accountprovider.Key, tt.account)
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

	if len(providers) != 1 || providers[0] != "account" {
		t.Errorf("RequiredProviders() = %v, want [account]", providers)
	}
}
