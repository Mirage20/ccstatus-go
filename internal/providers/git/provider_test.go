package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupGitRepo creates a temporary git repository for testing.
func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Initialize git repo
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test.com")
	runGit(t, dir, "config", "user.name", "Test User")

	return dir
}

// runGit executes a git command in the given directory.
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, output)
	}
	return string(output)
}

// createFile creates a file with the given content.
func createFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", name, err)
	}
}

func TestProvider_NotAGitRepo(t *testing.T) {
	dir := t.TempDir() // Empty dir, not a git repo

	p := &Provider{workDir: dir}
	result, err := p.Provide(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info := result.(*Info)
	if info.IsRepo {
		t.Error("expected IsRepo to be false for non-git directory")
	}
}

func TestProvider_EmptyRepo(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit so we have a branch
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", "README.md")
	runGit(t, dir, "commit", "-m", "initial")

	p := &Provider{workDir: dir}
	result, err := p.Provide(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info := result.(*Info)
	if !info.IsRepo {
		t.Error("expected IsRepo to be true")
	}
	if info.Branch == "" {
		t.Error("expected Branch to be set")
	}
	if info.Staged != 0 || info.Modified != 0 || info.Untracked != 0 {
		t.Errorf("expected clean repo, got staged=%d modified=%d untracked=%d",
			info.Staged, info.Modified, info.Untracked)
	}
}

func TestProvider_Branch(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	// Default branch (main or master depending on git version)
	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if info.Branch == "" {
		t.Error("expected branch name to be set")
	}

	// Create and checkout new branch
	runGit(t, dir, "checkout", "-b", "feature/test-branch")

	result, _ = p.Provide(context.Background())
	info = result.(*Info)

	if info.Branch != "feature/test-branch" {
		t.Errorf("expected branch 'feature/test-branch', got %q", info.Branch)
	}
}

func TestProvider_DetachedHead(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	// Get commit hash and checkout detached
	hash := runGit(t, dir, "rev-parse", "HEAD")
	runGit(t, dir, "checkout", "--detach", "HEAD")

	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	// Should start with @ for detached HEAD
	if len(info.Branch) == 0 || info.Branch[0] != '@' {
		t.Errorf("expected branch to start with '@' for detached HEAD, got %q (hash: %s)", info.Branch, hash)
	}
}

func TestProvider_StatusCounts(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	// Create untracked file
	createFile(t, dir, "untracked.txt", "untracked")

	// Create modified file (not staged)
	createFile(t, dir, "README.md", "# Modified")

	// Create staged file
	createFile(t, dir, "staged.txt", "staged")
	runGit(t, dir, "add", "staged.txt")

	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if info.Staged != 1 {
		t.Errorf("expected 1 staged file, got %d", info.Staged)
	}
	if info.Modified != 1 {
		t.Errorf("expected 1 modified file, got %d", info.Modified)
	}
	if info.Untracked != 1 {
		t.Errorf("expected 1 untracked file, got %d", info.Untracked)
	}
}

func TestProvider_Conflicts(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Initial")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	// Create a branch with conflicting changes
	runGit(t, dir, "checkout", "-b", "feature")
	createFile(t, dir, "README.md", "# Feature branch change")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "feature change")

	// Go back to main and make conflicting change
	runGit(t, dir, "checkout", "-")
	createFile(t, dir, "README.md", "# Main branch change")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "main change")

	// Try to merge (will conflict)
	cmd := exec.Command("git", "merge", "feature")
	cmd.Dir = dir
	_ = cmd.Run() // Ignore error, merge will fail due to conflict

	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if info.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", info.Conflicts)
	}
}

func TestProvider_NoUpstream(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if info.HasUpstream {
		t.Error("expected HasUpstream to be false for local-only branch")
	}
	if info.Ahead != 0 || info.Behind != 0 {
		t.Errorf("expected ahead=0 behind=0 without upstream, got ahead=%d behind=%d",
			info.Ahead, info.Behind)
	}
}

func TestProvider_AheadBehind(t *testing.T) {
	// Create "remote" repo
	remoteDir := t.TempDir()
	runGit(t, remoteDir, "init", "--bare")

	// Create local repo
	dir := setupGitRepo(t)

	// Create initial commit and push
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")
	runGit(t, dir, "remote", "add", "origin", remoteDir)
	runGit(t, dir, "push", "-u", "origin", "HEAD")

	// Create local commit (ahead by 1)
	createFile(t, dir, "local.txt", "local")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "local commit")

	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if !info.HasUpstream {
		t.Error("expected HasUpstream to be true")
	}
	if info.Ahead != 1 {
		t.Errorf("expected ahead=1, got %d", info.Ahead)
	}
	if info.Behind != 0 {
		t.Errorf("expected behind=0, got %d", info.Behind)
	}
}

func TestProvider_Stash(t *testing.T) {
	dir := setupGitRepo(t)

	// Create initial commit
	createFile(t, dir, "README.md", "# Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "initial")

	// No stash initially
	p := &Provider{workDir: dir}
	result, _ := p.Provide(context.Background())
	info := result.(*Info)

	if info.Stash != 0 {
		t.Errorf("expected 0 stash entries, got %d", info.Stash)
	}

	// Create stash
	createFile(t, dir, "README.md", "# Modified")
	runGit(t, dir, "stash", "push", "-m", "test stash")

	result, _ = p.Provide(context.Background())
	info = result.(*Info)

	if info.Stash != 1 {
		t.Errorf("expected 1 stash entry, got %d", info.Stash)
	}

	// Add another stash
	createFile(t, dir, "README.md", "# Modified again")
	runGit(t, dir, "stash", "push", "-m", "test stash 2")

	result, _ = p.Provide(context.Background())
	info = result.(*Info)

	if info.Stash != 2 {
		t.Errorf("expected 2 stash entries, got %d", info.Stash)
	}
}
