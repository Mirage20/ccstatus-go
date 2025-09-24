---
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git log:*)
description: Create a git commit with staged changes
argument-hint: [commit message]
---

## Context

- Current git status: !`git status --short`
- Staged changes: !`git diff --cached`
- Current branch: !`git branch --show-current`
- Recent commits: !`git log --oneline -5`

## Task

Create a single git commit with the staged changes.

$ARGUMENTS

If a commit message was provided above, use it. Otherwise, generate an appropriate message that:
- Summarizes what changed
- Follows the repository's commit format
- Is clear and concise

## Important Notes

- Only commit staged changes
- If no changes are staged, suggest using `git add` first
- Verify the commit was created successfully
