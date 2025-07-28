## Complete Analysis of Duplicate Branch Status Calls in Test Files

After thoroughly researching the codebase, I've identified the following patterns and opportunities for improvement:

### 1. Available Utility Functions

The codebase provides several utility functions that can replace direct git commands:

**GitManager methods:**
- `GetCurrentBranch() (string, error)` - replacement for `git rev-parse --abbrev-ref HEAD`
- `GetCurrentBranchInPath(path string) (string, error)` - same as above but for specific path
- `GetUpstreamBranch(worktreePath string) (string, error)` - replacement for `git rev-parse --abbrev-ref @{upstream}`
- `GetAheadBehindCount(worktreePath string) (int, int, error)` - replacement for ahead/behind status checks
- `GetCommitHash(ref string) (string, error)` - replacement for `git rev-parse HEAD` or other refs
- `GetCommitHashInPath(path, ref string) (string, error)` - same as above but for specific path

**Manager wrapper methods:**
- `Manager.GetCurrentBranch() (string, error)` - wraps GitManager.GetCurrentBranch()
- `Manager.GetGitManager() *GitManager` - provides access to GitManager for other methods

### 2. Current Duplicate Patterns Found

#### **A. Branch Status Calls in sync_test.go (Lines mentioned in task)**

**Lines 210, 403, 410, 875, 882 - Getting current branch:**
```go
cmd := exec.Command("git", "branch", "--show-current")
cmd.Dir = filepath.Join(repoPath, "worktrees", "feat")
branchOutput, err := cmd.Output()
require.NoError(t, err)
assert.Equal(t, "develop", strings.TrimSpace(string(branchOutput)))
```

**Replacement approach:** Use `manager.GetGitManager().GetCurrentBranchInPath(worktreePath)`

#### **B. Upstream Branch Checks in push_test.go**

**Line 96 - Checking upstream:**
```go
gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
gitCmd.Dir = worktreePath
err := gitCmd.Run()
```

**Replacement approach:** Use `manager.GetGitManager().GetUpstreamBranch(worktreePath)`

#### **C. Commit Hash Retrieval in push_test.go and pull_test.go**

**push_test.go line 79, pull_test.go line 46:**
```go
gitCmd = exec.Command("git", "rev-parse", "HEAD")
gitCmd.Dir = repo.GetLocalPath()
output, err := gitCmd.Output()
```

**Replacement approach:** Use `manager.GetGitManager().GetCommitHashInPath(path, "HEAD")`

#### **D. Additional Patterns in Various Test Files**

1. **Branch listing patterns** (mergeback_test.go line 469)
2. **Worktree operations** (sync_test.go line 660) - these use git worktree commands that don't have utility replacements yet
3. **Basic git operations** (add, commit, pull) - these are operational commands, not status queries

### 3. Specific Locations That Need Changes

#### **High Priority - Direct Branch Status Queries:**

1. **cmd/sync_test.go:**
   - Line 210: `validateResult` function in "branch reference changed" test
   - Line 403: `validateResult` function in "tracked worktrees updated" test  
   - Line 410: Same function, checking manual worktree branch
   - Line 875: `TestSyncCommand_WorktreePromotion` main worktree check
   - Line 882: Same test, production worktree check

2. **cmd/push_test.go:**
   - Line 96: `checkUpstreamExists` function checking upstream configuration
   - Line 79: `getRemoteCommitHash` function getting commit hash

3. **cmd/pull_test.go:**
   - Line 46: Getting commit hash for comparison

#### **Medium Priority - Indirect Status Checks:**

1. **internal/git_test.go:**
   - Lines 502, 569: Using `ExecGitCommand` for commit hash retrieval (already uses utility)

### 4. Implementation Strategy

#### **Pattern 1: Branch Status in Worktree Paths**
Replace:
```go
cmd := exec.Command("git", "branch", "--show-current")
cmd.Dir = worktreePath
branchOutput, err := cmd.Output()
branch := strings.TrimSpace(string(branchOutput))
```

With:
```go
manager, err := createInitializedManager()
require.NoError(t, err)
branch, err := manager.GetGitManager().GetCurrentBranchInPath(worktreePath)
require.NoError(t, err)
```

#### **Pattern 2: Upstream Branch Checks**
Replace:
```go
gitCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
gitCmd.Dir = worktreePath
err := gitCmd.Run()
```

With:
```go
manager, err := createInitializedManager()
require.NoError(t, err)
upstream, err := manager.GetGitManager().GetUpstreamBranch(worktreePath)
// Note: empty string means no upstream, not an error
```

#### **Pattern 3: Commit Hash Retrieval**
Replace:
```go
gitCmd = exec.Command("git", "rev-parse", "HEAD")
gitCmd.Dir = path
output, err := gitCmd.Output()
hash := strings.TrimSpace(string(output))
```

With:
```go
manager, err := createInitializedManager()
require.NoError(t, err)
hash, err := manager.GetGitManager().GetCommitHashInPath(path, "HEAD")
require.NoError(t, err)
```

### 5. Benefits of Replacement

1. **Consistency:** All git operations use the same error handling and command execution patterns
2. **Error handling:** Better error messages through `enhanceGitError` function
3. **Testability:** Utilities handle edge cases like missing upstream branches gracefully
4. **Maintainability:** Single source of truth for git command execution
5. **Robustness:** Built-in handling for common git failure scenarios

### 6. What NOT to Replace

Some git commands should remain as direct exec calls:
- **Operational commands:** `git add`, `git commit`, `git push`, `git pull` (these modify state)
- **Worktree management:** `git worktree add`, `git worktree remove` (no utilities exist yet)
- **Repository setup commands:** Used in test setup that creates the test environment

### 7. Recommended Implementation Order

1. **Start with sync_test.go lines 210, 403, 410, 875, 882** - these are the specific lines mentioned in the task
2. **Continue with push_test.go upstream checks (line 96)**
3. **Handle commit hash retrievals in push_test.go (line 79) and pull_test.go (line 46)**
4. **Review any additional patterns found during implementation**

This analysis provides a complete roadmap for replacing duplicate git command patterns with the available utility functions while maintaining test functionality and improving code consistency.