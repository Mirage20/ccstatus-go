package status

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

func TestComponent_Render_CleanWorkingTree(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:    true,
		Staged:    0,
		Modified:  0,
		Untracked: 0,
	})

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string for clean working tree, got %q", result)
	}
}

func TestComponent_Render_StagedOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Staged: 3,
	})

	result := c.Render(ctx)

	if !containsText(result, "3") {
		t.Errorf("expected result to contain '3', got %q", result)
	}
}

func TestComponent_Render_ModifiedOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:   true,
		Modified: 2,
	})

	result := c.Render(ctx)

	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2', got %q", result)
	}
}

func TestComponent_Render_UntrackedOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:    true,
		Untracked: 5,
	})

	result := c.Render(ctx)

	if !containsText(result, "5") {
		t.Errorf("expected result to contain '5', got %q", result)
	}
}

func TestComponent_Render_AllStatuses(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:    true,
		Staged:    3,
		Modified:  2,
		Untracked: 1,
		Conflicts: 1,
	})

	result := c.Render(ctx)

	// Check that counts are present (with nerd font icons)
	if !containsText(result, "3") {
		t.Errorf("expected result to contain '3' for staged, got %q", result)
	}
	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2' for modified, got %q", result)
	}
	if !containsText(result, "1") {
		t.Errorf("expected result to contain '1' for untracked/conflicts, got %q", result)
	}
}

func TestComponent_Render_ConflictsOnly(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:    true,
		Conflicts: 2,
	})

	result := c.Render(ctx)

	if result == "" {
		t.Error("expected non-empty result for conflicts")
	}
	if !containsText(result, "2") {
		t.Errorf("expected result to contain '2', got %q", result)
	}
}

func TestComponent_Render_CustomIcons(t *testing.T) {
	cfg := defaultConfig()
	cfg.StagedIcon = "S"
	cfg.ModifiedIcon = "M"
	cfg.UntrackedIcon = "U"
	cfg.ConflictIcon = "C"
	c := &Component{config: cfg}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo:    true,
		Staged:    1,
		Modified:  2,
		Untracked: 3,
		Conflicts: 4,
	})

	result := c.Render(ctx)

	if !containsText(result, "S1") {
		t.Errorf("expected result to contain 'S1', got %q", result)
	}
	if !containsText(result, "M2") {
		t.Errorf("expected result to contain 'M2', got %q", result)
	}
	if !containsText(result, "U3") {
		t.Errorf("expected result to contain 'U3', got %q", result)
	}
	if !containsText(result, "C4") {
		t.Errorf("expected result to contain 'C4', got %q", result)
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
	result := c.formatCount(0, "+", testColor)
	if result != "" {
		t.Errorf("formatCount(0, ...) = %q, want empty", result)
	}

	// Test non-zero count returns formatted string with color
	result = c.formatCount(5, "+", testColor)
	if result == "" {
		t.Error("formatCount(5, ...) returned empty, want non-empty")
	}
	if !containsText(result, "+5") {
		t.Errorf("formatCount(5, '+', ...) = %q, want to contain '+5'", result)
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
