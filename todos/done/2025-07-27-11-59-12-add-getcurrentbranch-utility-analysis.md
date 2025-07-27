Now I have a comprehensive understanding of the codebase. Let me create a detailed research report:

## Research Report: Git Command Patterns and GitManager Structure

### 1. Current "git rev-parse --abbrev-ref HEAD" Usage

The pattern `git rev-parse --abbrev-ref HEAD` is currently used in **4 locations**:

#### **Direct exec.Command Usage (needs refactoring):**
- **`cmd/info.go:291`** - In `getBaseBranchInfo()` function, gets current branch but result is discarded (commented as "not used for base branch detection anymore")

#### **Already using ExecGitCommand utility:**
- **`internal/git.go:573`** - In `GitManager.GetCurrentBranch()` method (operates on repo path)
- **`internal/git.go:635`** - In `GitManager.PushWorktree()` method (operates on worktree path)  
- **`internal/git.go:666`** - In `GitManager.PullWorktree()` method (operates on worktree path)

### 2. GitManager Structure Analysis

**Location:** `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`

**Key Components:**
```go
type GitManager struct {
	repo           *git.Repository
	repoPath       string
	worktreePrefix string
}
```

**Existing Current Branch Methods:**
- `GetCurrentBranch() (string, error)` - Works on repo path only
- **Missing:** `GetCurrentBranch(worktreePath string) (string, error)` - Needed for worktree-specific operations

### 3. Utility Function Patterns

**Git Command Execution Utilities:**
- `ExecGitCommand(dir string, args ...string) ([]byte, error)` - Standard output capture
- `execGitCommandRun(dir string, args ...string) error` - No output capture
- `ExecGitCommandCombined(dir string, args ...string) ([]byte, error)` - Combined stdout/stderr
- `ExecGitCommandInteractive(dir string, args ...string) error` - Live terminal output

**Error Handling Pattern:**
- `enhanceGitError(err error, operation string) error` - Provides context-specific error messages based on exit codes

### 4. Testing Infrastructure

**Test Patterns:**
- Uses `testify` for assertions (`require.NoError`, `assert.Equal`)
- Test utilities in `/Users/jschneider/code/scratch/worktree-manager/internal/testutils/`
- `GitTestRepo` struct provides isolated git repository testing
- Helper functions like `must(t *testing.T, err error)` for error handling
- `verifyWorktreeLinked()` for worktree validation

**Testing File Naming:**
- Pattern: `*_test.go` files co-located with source
- Example: `internal/git_add_test.go` tests `GitManager.AddWorktree()`

### 5. Files Requiring Modification

Based on the research and planned refactoring:

#### **Primary Target:**
- **`cmd/info.go`** - Line 291 needs to use new `GetCurrentBranch(worktreePath)` utility

#### **Enhancement Needed:**
- **`internal/git.go`** - Add new `GetCurrentBranch(worktreePath string)` method to support worktree-specific operations

#### **Test Files to Create/Update:**
- **`internal/git_test.go`** or **`cmd/info_test.go`** - Add tests for the new utility function

### 6. Implementation Plan

The planned refactoring aligns with the existing codebase patterns:

1. **Add new utility method** to `GitManager` that accepts `worktreePath` parameter
2. **Use existing `ExecGitCommand`** pattern for consistency
3. **Apply `enhanceGitError`** for consistent error handling
4. **Replace direct `exec.Command`** usage in `cmd/info.go:291`
5. **Add comprehensive tests** following existing test patterns

### 7. Error Handling Requirements

The new utility should:
- Use `ExecGitCommand()` for execution
- Apply `enhanceGitError()` for consistent error messages
- Return descriptive errors when not in a git repository
- Handle detached HEAD scenarios appropriately

This research confirms that the codebase has good architectural patterns in place, and the proposed refactoring will eliminate the remaining direct git command usage while maintaining consistency with existing utilities.