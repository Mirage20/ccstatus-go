package git

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/mirage20/ccstatus-go/internal/config"
	"github.com/mirage20/ccstatus-go/internal/core"
)

const (
	// Timeout for git commands.
	gitTimeout = 500 * time.Millisecond
)

func init() {
	// Self-register with type factory
	core.RegisterProvider(string(Key), New, func() any {
		return &Info{}
	})
}

// Provider provides git repository information.
type Provider struct {
	workDir string
}

// New creates a new git provider with config.
func New(cfgReader *config.Reader, session *core.ClaudeSession) (core.Provider, core.CacheConfig) {
	cfg := config.GetProvider(cfgReader, "git", defaultConfig())

	return &Provider{
		workDir: session.Workspace.CurrentDir,
	}, cfg.Cache
}

// Key returns the unique identifier for this provider.
func (p *Provider) Key() core.ProviderKey {
	return Key
}

// gitCmd creates a git command with --no-optional-locks flag to prevent lock contention
// with concurrent git operations (e.g., Claude Code running git commands).
func (p *Provider) gitCmd(ctx context.Context, args ...string) *exec.Cmd {
	//nolint:gosec // args are hardcoded within this package, not user input
	cmd := exec.CommandContext(ctx, "git", append([]string{"--no-optional-locks"}, args...)...)
	cmd.Dir = p.workDir
	return cmd
}

// Provide returns git repository information.
func (p *Provider) Provide(ctx context.Context) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, gitTimeout)
	defer cancel()

	info := &Info{}

	// Check if it's a git repo by getting the branch
	branch, err := p.getBranch(ctx)
	if err != nil {
		// Not a git repo or git not available - return empty info (not an error)
		return info, nil
	}

	info.IsRepo = true
	info.Branch = branch

	// Get status counts (staged, modified, untracked, conflicts)
	info.Staged, info.Modified, info.Untracked, info.Conflicts = p.getStatusCounts(ctx)

	// Get ahead/behind (only if we have an upstream)
	info.Ahead, info.Behind, info.HasUpstream = p.getAheadBehind(ctx)

	// Get stash count
	info.Stash = p.getStashCount(ctx)

	return info, nil
}

// getBranch returns the current branch name or "@<hash>" for detached HEAD.
func (p *Provider) getBranch(ctx context.Context) (string, error) {
	// Try to get branch name
	cmd := p.gitCmd(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	branch := strings.TrimSpace(string(output))

	// If detached HEAD, get short commit hash instead
	if branch == "HEAD" {
		return p.getShortHash(ctx)
	}

	return branch, nil
}

// getShortHash returns the short commit hash prefixed with "@".
func (p *Provider) getShortHash(ctx context.Context) (string, error) {
	cmd := p.gitCmd(ctx, "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return "@" + strings.TrimSpace(string(output)), nil
}

// getStatusCounts returns counts of staged, modified, untracked, and conflicted files.
//
//nolint:nonamedreturns // named returns document the meaning of each count
func (p *Provider) getStatusCounts(ctx context.Context) (staged, modified, untracked, conflicts int) {
	cmd := p.gitCmd(ctx, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, 0
	}

	for line := range strings.SplitSeq(strings.TrimSuffix(string(output), "\n"), "\n") {
		if len(line) < 2 { //nolint:mnd // porcelain format has 2 char prefix
			continue
		}
		// Porcelain format: XY filename
		// X = staged status, Y = unstaged status
		x, y := line[0], line[1]

		// Untracked: both X and Y are '?'
		if x == '?' && y == '?' {
			untracked++
			continue
		}

		// Conflicts: U in either position, or DD/AA/AU/UA/DU/UD combinations
		if x == 'U' || y == 'U' || (x == 'D' && y == 'D') || (x == 'A' && y == 'A') {
			conflicts++
			continue
		}

		// Staged: X is not ' ' (already handled '?' above)
		if x != ' ' {
			staged++
		}

		// Modified (unstaged): Y is 'M', 'D', etc. (not ' ')
		if y != ' ' {
			modified++
		}
	}

	return staged, modified, untracked, conflicts
}

// getAheadBehind returns commits ahead/behind upstream and whether upstream exists.
//
//nolint:nonamedreturns // named returns document the meaning of each value
func (p *Provider) getAheadBehind(ctx context.Context) (ahead, behind int, hasUpstream bool) {
	// Check if upstream exists
	cmd := p.gitCmd(ctx, "rev-parse", "--abbrev-ref", "@{upstream}")
	if err := cmd.Run(); err != nil {
		return 0, 0, false
	}

	// Get ahead/behind counts
	cmd = p.gitCmd(ctx, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}

	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) == 2 { //nolint:mnd // expected format: "ahead behind"
		ahead, _ = strconv.Atoi(parts[0])
		behind, _ = strconv.Atoi(parts[1])
	}

	return ahead, behind, true
}

// getStashCount returns the number of stash entries.
func (p *Provider) getStashCount(ctx context.Context) int {
	cmd := p.gitCmd(ctx, "stash", "list")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return 0
	}

	return len(strings.Split(trimmed, "\n"))
}
