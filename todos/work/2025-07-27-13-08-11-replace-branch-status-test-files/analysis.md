Based on my analysis, here are my specific findings about the branch status calls in `cmd/sync_test.go` that need to be replaced with GitManager utilities:

## Current State Analysis

I found the exact git commands mentioned in the todo at the specified line numbers:

### Line 210 (approx):
```go
cmd := exec.Command("git", "branch", "--show-current")
cmd.Dir = filepath.Join(repoPath, "worktrees", "feat")
branchOutput, err := cmd.Output()
require.NoError(t, err)
assert.Equal(t, "develop", strings.TrimSpace(string(branchOutput)))
```

### Lines 403 & 410:
```go
// Line 403: Verify feat worktree was updated to develop branch  
cmd := exec.Command("git", "branch", "--show-current")
cmd.Dir = filepath.Join(repoPath, "worktrees", "feat")
featOutput, err := cmd.Output()
require.NoError(t, err)
assert.Equal(t, "develop", strings.TrimSpace(string(featOutput)))

// Line 410: manual worktree should still be on main branch (unchanged)
manualCmd := exec.Command("git", "branch", "--show-current")
manualCmd.Dir = filepath.Join(repoPath, "worktrees", "manual")
manualOutput, err := manualCmd.Output()
require.NoError(t, err)
assert.Equal(t, "main", strings.TrimSpace(string(manualOutput)))
```

### Lines 875 & 882:
```go
// Line 875: Check main worktree is on main branch
mainCmd := exec.Command("git", "branch", "--show-current")
mainCmd.Dir = filepath.Join(repoPath, "worktrees", "main")
mainOutput, err := mainCmd.Output()
require.NoError(t, err)
assert.Equal(t, "main", strings.TrimSpace(string(mainOutput)), "main worktree should be on main branch")

// Line 882: Check production worktree is on production-2025-07-1 branch (promoted from preview)
prodCmd := exec.Command("git", "branch", "--show-current")
prodCmd.Dir = filepath.Join(repoPath, "worktrees", "production")
prodOutput, err := prodCmd.Output()
require.NoError(t, err)
assert.Equal(t, "production-2025-07-1", strings.TrimSpace(string(prodOutput)), "production worktree should be on production-2025-07-1 branch")
```

## Commands and Appropriate Utilities

### 1. `git branch --show-current` → `GetCurrentBranchInPath`
All the identified usages are running `git branch --show-current` in specific worktree directories to verify what branch that worktree is currently on. This directly maps to:

```go
// Replace: git branch --show-current
// With: manager.GetGitManager().GetCurrentBranchInPath(worktreePath)
```

The `GetCurrentBranchInPath` utility internally uses `git rev-parse --abbrev-ref HEAD` which is equivalent to `git branch --show-current`.

### 2. Usage Contexts
All the usages are in **test validation functions** where they:
- Verify that worktrees were correctly updated to expected branches after sync operations
- Confirm branch state after worktree promotions and configuration changes
- Test both positive cases (correct branch) and edge cases (unchanged worktrees)

### 3. Other Git Commands Found (Not in Todo)
I found two additional git commands in the test file:
- **Line 622**: `git worktree prune` - Administrative command, probably doesn't need replacement
- **Line 660**: `git worktree add --force` - Worktree creation command, already has utilities in GitManager

## Additional Patterns Not Mentioned in Todo

I found **no other instances** of the patterns mentioned in the todo:
- No `git rev-parse --abbrev-ref @{upstream}` calls (which would use `GetUpstreamBranch`)
- No `git rev-list --left-right --count HEAD...@{upstream}` calls (which would use `GetAheadBehindCount`)

## Conversion Strategy

All 5 instances should be converted to:

```go
// Instead of:
cmd := exec.Command("git", "branch", "--show-current")
cmd.Dir = worktreePath
output, err := cmd.Output()
require.NoError(t, err)
branchName := strings.TrimSpace(string(output))

// Use:
manager, err := createInitializedManager()
require.NoError(t, err)
branchName, err := manager.GetGitManager().GetCurrentBranchInPath(worktreePath)
require.NoError(t, err)
```

The tests already have access to `createInitializedManager()` function and show examples of using it for GitManager operations, so this conversion should be straightforward and consistent with existing patterns in the codebase.