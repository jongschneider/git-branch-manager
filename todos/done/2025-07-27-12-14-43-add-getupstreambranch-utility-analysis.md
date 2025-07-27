Based on my comprehensive research of the codebase, here are my detailed findings about the GetUpstreamBranch utility implementation:

## Current Implementation Patterns

### 1. Git Command Execution Patterns

The codebase uses several consistent patterns for executing git commands:

- **ExecGitCommand()**: Main utility function for capturing git command output
- **execGitCommandRun()**: For commands that don't need output capture
- **ExecGitCommandCombined()**: For commands that need combined stdout/stderr
- **ExecGitCommandInteractive()**: For commands with live terminal output

### 2. Current Usage of `git rev-parse --abbrev-ref @{upstream}`

I found **4 instances** where this exact command pattern is currently used:

1. **`/Users/jschneider/code/scratch/worktree-manager/internal/git.go:653`** - In `PushWorktree()` method
2. **`/Users/jschneider/code/scratch/worktree-manager/internal/git.go:686`** - In `PullWorktree()` method  
3. **`/Users/jschneider/code/scratch/worktree-manager/cmd/info.go:298`** - In info command for branch details
4. **`/Users/jschneider/code/scratch/worktree-manager/cmd/push_test.go:96`** - In test helper `checkUpstreamExists()`

### 3. GitManager Structure and Method Patterns

The `GitManager` struct follows these patterns:

- **Method naming**: `Get{Something}()`, `Get{Something}InPath()` for path-specific operations
- **Error handling**: Uses `enhanceGitError()` for consistent error messaging
- **Return patterns**: Returns `(string, error)` for single values, `([]string, error)` for collections

**Existing branch-related methods:**
- `GetCurrentBranch()` - gets current branch from repo path
- `GetCurrentBranchInPath(path string)` - gets current branch from specific path
- `GetDefaultBranch()` - gets repository default branch
- `GetRemoteBranches()` - gets list of remote branches

### 4. Error Handling Patterns

The current implementations handle the upstream command in two ways:

**Pattern 1 - Silent failure (info.go:298-304):**
```go
cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
cmd.Dir = worktreePath
output, err := cmd.Output()
upstream := ""
if err == nil {
    upstream = strings.TrimSpace(string(output))
}
```

**Pattern 2 - Error checking for logic flow (git.go:653, 686):**
```go
_, err = ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")
if err != nil {
    // No upstream set, handle accordingly
}
```

### 5. Testing Patterns

The codebase uses:
- **testify/assert** and **testify/require** for assertions
- **testutils.GitTestRepo** for creating test repositories
- Table-driven tests with setup functions
- Helper functions like `checkUpstreamExists()` for common operations

## Recommended GetUpstreamBranch Implementation

Based on the analysis, here's how `GetUpstreamBranch` should be implemented:

### Method Signature and Placement
```go
// In GitManager struct
func (gm *GitManager) GetUpstreamBranch(worktreePath string) (string, error)
```

### Implementation Approach
1. **Use ExecGitCommand** for consistency with other methods
2. **Apply enhanceGitError** for consistent error messaging  
3. **Return empty string and specific error** when no upstream is set (rather than silent failure)
4. **Handle path validation** like other worktree methods
5. **Follow the same string trimming pattern** as other branch methods

### Error Cases to Handle
- **Path doesn't exist** - consistent with other worktree methods
- **Not a git repository** - enhanced through enhanceGitError
- **No upstream configured** - specific error message
- **Branch exists but no upstream** - differentiate from other errors

### Usage Pattern
The method should support both use cases found in the codebase:
1. **Info gathering** - where upstream may or may not exist
2. **Logic flow control** - where upstream existence determines behavior

### Testing Requirements
Should include tests for:
- Valid worktree with upstream set
- Valid worktree without upstream
- Invalid worktree path
- Non-git directory
- Different upstream naming patterns (origin/branch vs. other remotes)

This implementation would consolidate the 4 duplicate command executions into a single, well-tested utility method that follows the established patterns in the codebase.

---

Based on my analysis of the specific usage contexts for the GetUpstreamBranch utility, here's a comprehensive breakdown:

## Analysis of GetUpstreamBranch Usage Contexts

### 1. Usage in `cmd/info.go` (around line 300)

**Context**: The `getBaseBranchInfo` function that builds branch information for the info command.

**Current code pattern**:
```go
// Get upstream branch
cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
cmd.Dir = worktreePath
output, err := cmd.Output()
upstream := ""
if err == nil {
    upstream = strings.TrimSpace(string(output))
}
```

**Analysis**:
- **Variable assignment**: `upstream` variable stores the result
- **Error handling**: Errors are silently ignored - if no upstream exists, `upstream` remains empty string
- **Usage**: The `upstream` value is assigned to `BranchInfo.Upstream` field and returned
- **Context**: This is part of gathering comprehensive branch information for display

### 2. Usage in `internal/git.go` PushWorktree method (around line 653)

**Context**: Determining whether to push with `-u` flag or do a simple push.

**Current code pattern**:
```go
// Check if upstream is set
_, err = ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")

var cmd *exec.Cmd
if err != nil {
    // No upstream set, push with -u flag
    cmd = exec.Command("git", "push", "-u", "origin", currentBranch)
} else {
    // Upstream is set, simple push
    cmd = exec.Command("git", "push")
}
```

**Analysis**:
- **Variable assignment**: Only the error is checked, actual upstream value is discarded
- **Error handling**: Error indicates no upstream exists, which drives conditional logic
- **Usage**: Determines push strategy - set upstream vs. simple push
- **Context**: Part of pushing changes from a worktree

### 3. Usage in `internal/git.go` PullWorktree method (around line 686)

**Context**: Determining whether upstream is set before attempting pull.

**Current code pattern**:
```go
// Check if upstream is set
_, err = ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")
if err != nil {
    // No upstream set, try to set it and pull
    remoteBranch := Remote(currentBranch)
    
    // Check if remote branch exists
    _, err = ExecGitCommand(worktreePath, "rev-parse", "--verify", remoteBranch)
    if err == nil {
        // Remote branch exists, set upstream and pull
        _, err = ExecGitCommand(worktreePath, "branch", "--set-upstream-to", remoteBranch)
        // ...
    }
}
```

**Analysis**:
- **Variable assignment**: Only the error is checked, actual upstream value is discarded
- **Error handling**: Error indicates no upstream, triggers upstream setup logic
- **Usage**: Determines whether to set up upstream tracking before pulling
- **Context**: Part of pulling changes into a worktree

## Interface Design Recommendations

Based on these usage patterns, the `GetUpstreamBranch` method should have this interface:

```go
// GetUpstreamBranch returns the upstream branch name for a given worktree path.
// Returns empty string if no upstream is set (not an error condition).
func (gm *GitManager) GetUpstreamBranch(worktreePath string) (string, error)
```

### Key Design Decisions:

1. **Return empty string, not error for missing upstream**: In the info.go usage, missing upstream is handled gracefully by setting `upstream = ""`. The method should follow this pattern.

2. **Only return error for actual git failures**: Real errors (repository issues, invalid paths) should be returned as errors, but missing upstream is a normal state.

3. **Follow existing GitManager patterns**: Like `GetCurrentBranchInPath`, it takes a `worktreePath` parameter and uses `ExecGitCommand` internally.

4. **Use consistent error handling**: Should use `enhanceGitError()` for real git command failures.

### Implementation should be:

```go
func (gm *GitManager) GetUpstreamBranch(worktreePath string) (string, error) {
    output, err := ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "@{upstream}")
    if err != nil {
        // Check if this is a "no upstream" error vs a real git error
        if strings.Contains(err.Error(), "no upstream") || strings.Contains(err.Error(), "HEAD does not point to a branch") {
            return "", nil // No upstream set - not an error
        }
        return "", enhanceGitError(err, "get upstream branch")
    }
    return strings.TrimSpace(string(output)), nil
}
```

### Refactoring Impact:

1. **info.go**: Replace the direct git command with `manager.GetGitManager().GetUpstreamBranch(worktreePath)`
2. **PushWorktree**: Use `upstream, err := gm.GetUpstreamBranch(worktreePath); if upstream == "" {`
3. **PullWorktree**: Same pattern as PushWorktree

This design maintains backward compatibility while centralizing the git command logic and providing consistent error handling across all usage contexts.