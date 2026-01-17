package branch

import (
	"testing"

	"github.com/mirage20/ccstatus-go/internal/core"
	gitprovider "github.com/mirage20/ccstatus-go/internal/providers/git"
)

func TestComponent_Render_NotARepo(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	// No git info in context
	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string for non-repo, got %q", result)
	}
}

func TestComponent_Render_NotARepo_EmptyInfo(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	// Git info with IsRepo = false
	ctx.Set(gitprovider.Key, &gitprovider.Info{IsRepo: false})

	result := c.Render(ctx)

	if result != "" {
		t.Errorf("expected empty string for IsRepo=false, got %q", result)
	}
}

func TestComponent_Render_Branch(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Branch: "main",
	})

	result := c.Render(ctx)

	// Result should contain "main" and the icon
	if result == "" {
		t.Error("expected non-empty result")
	}
	// Check that branch name is in output (with ANSI codes)
	if !containsText(result, "main") {
		t.Errorf("expected result to contain 'main', got %q", result)
	}
}

func TestComponent_Render_DetachedHead(t *testing.T) {
	c := &Component{config: defaultConfig()}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Branch: "@abc1234",
	})

	result := c.Render(ctx)

	if !containsText(result, "@abc1234") {
		t.Errorf("expected result to contain '@abc1234', got %q", result)
	}
}

func TestComponent_Render_Truncation(t *testing.T) {
	cfg := defaultConfig()
	cfg.MaxLength = 10
	c := &Component{config: cfg}
	ctx := core.NewRenderContext()

	ctx.Set(gitprovider.Key, &gitprovider.Info{
		IsRepo: true,
		Branch: "feature/very-long-branch-name",
	})

	result := c.Render(ctx)

	// Should be truncated and contain ellipsis
	if !containsText(result, "…") {
		t.Errorf("expected result to contain ellipsis for truncation, got %q", result)
	}
	// Should not contain full branch name
	if containsText(result, "feature/very-long-branch-name") {
		t.Errorf("expected branch to be truncated, got %q", result)
	}
}

func TestTruncateMiddle(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"feature/long-branch", 10, "feat…ranch"},
		{"abcdefghij", 5, "ab…ij"},
		{"abcdefghij", 10, "abcdefghij"},
	}

	for _, tc := range tests {
		result := truncateMiddle(tc.input, tc.maxLen)
		if result != tc.expected {
			t.Errorf("truncateMiddle(%q, %d) = %q, want %q",
				tc.input, tc.maxLen, result, tc.expected)
		}
	}
}

func TestComponent_RequiredProviders(t *testing.T) {
	c := &Component{config: defaultConfig()}
	providers := c.RequiredProviders()

	if len(providers) != 1 || providers[0] != "git" {
		t.Errorf("expected [\"git\"], got %v", providers)
	}
}

// containsText checks if s contains text, ignoring ANSI escape codes.
func containsText(s, text string) bool {
	// Simple check - just look for the text in the string
	// ANSI codes don't affect this since we're looking for substrings
	return len(s) > 0 && len(text) > 0 &&
		(len(s) >= len(text)) &&
		(findText(s, text))
}

func findText(s, text string) bool {
	for i := 0; i <= len(s)-len(text); i++ {
		if s[i:i+len(text)] == text {
			return true
		}
	}
	return false
}
