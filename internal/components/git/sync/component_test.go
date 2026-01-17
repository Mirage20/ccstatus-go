package sync

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	"github.com/mirage20/ccstatus-go/internal/format"
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

func TestComponent_Render_NoUpstream(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: false,
	})

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string when no upstream, got %q", result)
	}
}

func TestComponent_Render_InSync(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: true,
		Ahead:       0,
		Behind:      0,
	})

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string when in sync, got %q", result)
	}
}

func TestComponent_Render_AheadOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: true,
		Ahead:       3,
		Behind:      0,
	})

	result := c.Render(ctx)

	if result == "" {
		t.Error("expected non-empty result for ahead")
	}
	if !containsText(result, "3") {
		t.Errorf("expected result to contain '3', got %q", result)
	}
}

func TestComponent_Render_BehindOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: true,
		Ahead:       0,
		Behind:      2,
	})

	result := c.Render(ctx)

	if result == "" {
		t.Error("expected non-empty result for behind")
	}
	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2', got %q", result)
	}
}

func TestComponent_Render_AheadAndBehind(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: true,
		Ahead:       3,
		Behind:      2,
	})

	result := c.Render(ctx)

	if !containsText(result, "3") {
		t.Errorf("expected result to contain '3' for ahead, got %q", result)
	}
	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2' for behind, got %q", result)
	}
}

func TestComponent_Render_CustomIcons(t *testing.T) {
	cfg := defaultConfig()
	cfg.AheadIcon = "A"
	cfg.BehindIcon = "B"
	c := &Component{config: cfg}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:      true,
		HasUpstream: true,
		Ahead:       1,
		Behind:      2,
	})

	result := c.Render(ctx)

	if !containsText(result, "A1") {
		t.Errorf("expected result to contain 'A1', got %q", result)
	}
	if !containsText(result, "B2") {
		t.Errorf("expected result to contain 'B2', got %q", result)
	}
}

func TestComponent_RequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "git" {
		t.Errorf("expected [\"git\"], got %v", providers)
	}
}

func TestFormatCount(t *testing.T) {
	c := &Component{config: defaultConfig()}
	testColor := format.ParseColor("green")

	// Test zero count returns empty
	result := c.formatCount(0, "^", testColor)
	if result != "" {
		t.Errorf("formatCount(0, ...) = %q, want empty", result)
	}

	// Test non-zero count returns formatted string
	result = c.formatCount(5, "^", testColor)
	if result == "" {
		t.Error("formatCount(5, ...) returned empty, want non-empty")
	}
	if !containsText(result, "^5") {
		t.Errorf("formatCount(5, '^', ...) = %q, want to contain '^5'", result)
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
