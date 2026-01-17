package stash

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func TestComponent_Render_NotARepo(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string for non-repo, got %q", result)
	}
}

func TestComponent_Render_NoStash(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Stash:  0,
	})

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string when no stash, got %q", result)
	}
}

func TestComponent_Render_WithStash(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Stash:  3,
	})

	result := c.Render(ctx)

	if result == "" {
		t.Error("expected non-empty result for stash")
	}
	if !containsText(result, "3") {
		t.Errorf("expected result to contain '3', got %q", result)
	}
}

func TestComponent_Render_CustomIcon(t *testing.T) {
	cfg := defaultConfig()
	cfg.Icon = "S"
	c := &Component{config: cfg}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Stash:  2,
	})

	result := c.Render(ctx)

	if !containsText(result, "S") {
		t.Errorf("expected result to contain 'S', got %q", result)
	}
	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2', got %q", result)
	}
}

func TestComponent_RequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "git" {
		t.Errorf("expected [\"git\"], got %v", providers)
	}
}

// containsText checks if s contains text.
func containsText(s, text string) bool {
	for i := 0; i <= len(s)-len(text); i++ {
		if s[i:i+len(text)] == text {
			return true
		}
	}
	return false
}
