package git

// Info represents git repository information.
type Info struct {
	// Branch is the current branch name, or "@<hash>" for detached HEAD
	Branch string

	// IsRepo indicates if the directory is a git repository
	IsRepo bool

	// Staged is the number of staged files
	Staged int

	// Modified is the number of modified (unstaged) files
	Modified int

	// Untracked is the number of untracked files
	Untracked int

	// Conflicts is the number of files with merge conflicts
	Conflicts int

	// Ahead is the number of commits ahead of upstream (0 if no upstream)
	Ahead int

	// Behind is the number of commits behind upstream (0 if no upstream)
	Behind int

	// HasUpstream indicates if the branch has a tracking upstream
	HasUpstream bool

	// Stash is the number of stash entries
	Stash int
}
