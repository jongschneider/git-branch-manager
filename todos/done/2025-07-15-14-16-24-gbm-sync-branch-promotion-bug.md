# Fix gbm sync bug when promoting a branch to a new name
**Status:** Done
**Agent PID:** 12164

## Original Todo
fix `gbm sync` bug when promoting a branch to a new name. Example:
   * From
   ```toml
   # state.toml
# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  # Primary worktree - no merge_into (root of merge chain)
  master:
    branch: master
    description: "Main production branch"
  preview:
    branch: production-2025-07-1 
    description: "Blade Runner"
    merge_into: preview
  production:
    branch: production-2025-05-1 
    description: "Arrival"
    merge_into: preview
   ```
   TO
   ```toml
# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  # Primary worktree - no merge_into (root of merge chain)
  master:
    branch: master
    description: "Main production branch"
  production:
    branch: production-2025-07-1 
    description: "Blade Runner"
    merge_into: master
   ```
  ```sh
 gbm sync
Error: failed to update worktree for production: failed to create worktree: exit status 128
Usage:
  gbm sync [flags]

Flags:
     --dry-run   show what would be changed without making changes
     --force     skip confirmation prompts and remove orphaned worktrees
 -h, --help      help for sync

Global Flags:
     --debug                 enable debug logging to ./gbm.log
     --worktree-dir string   override worktree directory location

❌ ERROR: Error: failed to update worktree for production: failed to create worktree: exit status 128
   ```

## Description
Fix the `gbm sync` bug that occurs when promoting a branch to a new name in the configuration. The issue happens during worktree updates when the `UpdateWorktree` function tries to remove an existing worktree and recreate it with a different branch. The bug manifests as "exit status 128" when Git cannot create the new worktree, typically because the target branch is already checked out in another worktree or due to Git's internal state conflicts during the remove-recreate process.

## Implementation Plan
Fix the git worktree creation failure by implementing proper worktree promotion logic and adding clear destructive action confirmations:
- [x] Add worktree promotion detection logic in `internal/manager.go` to identify when a branch is being moved between worktrees (not just branch changes)
- [x] Implement new PromoteWorktree method in `internal/git.go` that uses `git worktree move` instead of remove+recreate for promotion scenarios
- [x] Add confirmation prompt in sync command that clearly describes which worktree will be promoted and which will be destroyed (e.g., "Worktree preview (production-2025-07-1) will be promoted to production. Worktree production (production-2025-05-1) will be removed.")
- [x] Update UpdateWorktree method in `internal/git.go:349-359` to route to PromoteWorktree when appropriate, falling back to current logic for simple branch changes
- [x] Add branch availability validation to prevent conflicts before any destructive operations
- [x] Improve error handling to provide specific messages for different failure scenarios
- [x] Automated test: Create test case that reproduces the branch promotion scenario and verifies both the confirmation prompt and the git worktree move operation
- [x] User test: Test the exact scenario from the bug report with proper confirmation dialog and successful promotion

## Notes
[Implementation notes]