Now I have a comprehensive understanding of the codebase. Let me provide my detailed analysis:

## Detailed Analysis of Codebase Research

Based on my examination of the worktree-manager codebase, here's a comprehensive analysis of the current state and patterns needed for implementing the GetCommitHash utility function:

### 1. Current git rev-parse Usage Patterns in internal/git_add.go (lines 48,54)

The specific git rev-parse commands being used are:

**Line 48:** `mergeBase, err := ExecGitCommand(gm.repoPath, "merge-base", branchName, baseBranch)`
- This is actually a `git merge-base` command, not `git rev-parse`

**Line 54:** `baseCommit, err := ExecGitCommand(gm.repoPath, "rev-parse", baseBranch)`
- This uses `git rev-parse <ref>` to get the commit hash of the base branch
- Purpose: Retrieve the commit hash of the base branch to compare with merge-base results
- Context: Used in worktree creation to verify branch relationships

### 2. GitManager Structure and Utility Functions

**GitManager Structure:**
```go
type GitManager struct {
    repo           *git.Repository
    repoPath       string
    worktreePrefix string
}
```

**Key Utility Function Patterns:**
- **Location**: All functions are in `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`
- **Naming Convention**: `Get{Something}()` for repository-level operations, `Get{Something}InPath()` for path-specific operations
- **Parameter Patterns**: Methods that work on specific paths take `path string` as first parameter
- **Return Patterns**: 
  - Single values: `(string, error)`
  - Boolean checks: `(bool, error)`
  - Collections: `([]string, error)`

**Existing Utilities that follow the pattern:**
- `GetCurrentBranch()` - Gets current branch from repository root
- `GetCurrentBranchInPath(path string)` - Gets current branch from specific path
- `GetUpstreamBranch(worktreePath string)` - Gets upstream branch info
- `VerifyRef(ref string)` - Verifies if a reference exists
- `VerifyRefInPath(path, ref string)` - Verifies reference in specific path

### 3. Error Handling Patterns with enhanceGitError()

The `enhanceGitError()` function is consistently used throughout GitManager utilities:

**Function Signature:**
```go
func enhanceGitError(err error, operation string) error
```

**Usage Pattern:**
```go
if err != nil {
    return enhanceGitError(err, "operation description")
}
```

**Examples from existing code:**
- `enhanceGitError(err, "get current branch")`
- `enhanceGitError(err, "get upstream branch")`
- `enhanceGitError(err, "verify ref")`
- `enhanceGitError(err, "worktree add")`

**Error Enhancement Features:**
- Analyzes exit codes (128, 1, etc.)
- Provides context-specific error messages
- Handles common git scenarios like "not a git repository", "already checked out", etc.

### 4. Testing Patterns for GitManager Functions

**Test File Structure:**
- Main GitManager tests: `/Users/jschneider/code/scratch/worktree-manager/internal/git_test.go`
- Test helper functions in `/Users/jschneider/code/scratch/worktree-manager/internal/git_add_test.go`

**Testing Infrastructure:**
- Uses `testutils.NewGitTestRepo()` for test repository setup
- Leverages `github.com/stretchr/testify` for assertions (`assert`, `require`)
- Pattern: Table-driven tests with setup/expect/error validation functions

**Test Patterns:**
```go
func TestGitManager_MethodName(t *testing.T) {
    repo := testutils.NewGitTestRepo(t, options...)
    defer repo.Cleanup()
    
    gitManager, err := NewGitManager(repo.GetLocalPath(), "worktrees")
    require.NoError(t, err)
    
    tests := []struct {
        name        string
        setup       func(t *testing.T, repo *testutils.GitTestRepo)
        // test parameters
        expectError bool
    }{
        // test cases
    }
}
```

**Common Test Scenarios:**
- Valid operations on existing branches/refs
- Error cases (non-existent paths, invalid refs)
- Edge cases (detached HEAD, no upstream, etc.)

### 5. Where and How to Add GetCommitHash Function

**Recommended Implementation Location:**
Add to `/Users/jschneider/code/scratch/worktree-manager/internal/git.go` following existing patterns.

**Suggested Function Signatures:**
```go
// GetCommitHash returns the commit hash for a given reference in the repository
func (gm *GitManager) GetCommitHash(ref string) (string, error)

// GetCommitHashInPath returns the commit hash for a given reference in a specific path
func (gm *GitManager) GetCommitHashInPath(path, ref string) (string, error)
```

**Implementation Pattern:**
```go
func (gm *GitManager) GetCommitHash(ref string) (string, error) {
    output, err := ExecGitCommand(gm.repoPath, "rev-parse", ref)
    if err != nil {
        return "", enhanceGitError(err, "get commit hash")
    }
    return strings.TrimSpace(string(output)), nil
}
```

**Replacement Target:**
- Line 54 in `internal/git_add.go`: Replace `ExecGitCommand(gm.repoPath, "rev-parse", baseBranch)` with `gm.GetCommitHash(baseBranch)`

**Test Coverage Needed:**
- Valid references (HEAD, branch names, commit hashes)
- Invalid references (non-existent branches)
- Error cases (not a git repository, invalid ref format)
- Both repository-level and path-specific versions

This analysis shows the codebase follows consistent patterns for git utilities, with strong error handling, comprehensive testing, and clear separation between repository-level and path-specific operations. The new GetCommitHash function should seamlessly integrate with these established patterns.