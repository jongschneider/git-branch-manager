# Analysis for VerifyRef Utility Function

## Current "git rev-parse --verify" Usage

Based on my examination of the codebase, here's a comprehensive analysis of how "git rev-parse --verify" is currently used:

### 1. Current Occurrences

I found **6 distinct usages** of "git rev-parse --verify" across the codebase:

#### **cmd/info.go:327** 
```go
cmd := exec.Command("git", "rev-parse", "--verify", candidate)
cmd.Dir = worktreePath
if _, err := cmd.Output(); err == nil {
    // Branch exists, check if it's actually a base
    cmd = exec.Command("git", "merge-base", "--is-ancestor", candidate, "HEAD")
    cmd.Dir = worktreePath
    if err := cmd.Run(); err == nil {
        baseBranch = candidate
        break
    }
}
```
**Purpose**: Verifying existence of candidate base branches ("main", "master", "develop", "dev")

#### **internal/git_add.go:89**
```go
_, err = ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)
if err == nil {
    // Remote tracking branch exists, use --track but don't create new branch
    finalArgs = append(finalArgs, "worktree", "add", "--track", worktreePath, remoteBranch)
} else {
    // No remote tracking branch, create worktree without tracking
    finalArgs = append(finalArgs, "worktree", "add", worktreePath, branchName)
}
```
**Purpose**: Checking if remote tracking branch exists before creating worktree

#### **internal/git.go:317** (in BranchExistsLocalOrRemote method)
```go
// Check if remote branch exists
remoteBranch := Remote(branchName)
_, err := ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)
return err == nil, nil
```
**Purpose**: Checking if remote branch exists for a given branch name

#### **internal/git.go:406** (in CreateWorktree method)
```go
// Check if remote tracking branch exists
remoteBranch := Remote(branchName)
_, err = ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)
if err == nil {
    // Remote tracking branch exists, create worktree and set up tracking
```
**Purpose**: Checking remote branch existence before creating worktree with tracking

#### **internal/git.go:738** (in PullFromUpstream method)
```go
// Check if remote branch exists
_, err = ExecGitCommand(worktreePath, "rev-parse", "--verify", remoteBranch)
if err == nil {
    // Remote branch exists, set upstream and pull
    _, err = ExecGitCommand(worktreePath, "branch", "--set-upstream-to", remoteBranch)
```
**Purpose**: Verifying remote branch exists before setting upstream and pulling

#### **cmd/mergeback.go:568 and :619** (two occurrences)
```go
output, err := internal.ExecGitCommand(repoRoot, "rev-parse", "--verify", branch)
if err == nil && strings.TrimSpace(string(output)) != "" {
    sourceBranch = branch
    break
}
```
**Purpose**: Finding first existing branch from a list of possible source branches

### 2. Current Command Structure & Returns

The current usage patterns follow these structures:

**Pattern 1: Using exec.Command directly (cmd/info.go)**
```go
cmd := exec.Command("git", "rev-parse", "--verify", candidate)
cmd.Dir = worktreePath
if _, err := cmd.Output(); err == nil {
    // Branch exists
}
```

**Pattern 2: Using ExecGitCommand utility (most common)**
```go
_, err := ExecGitCommand(gm.repoPath, "rev-parse", "--verify", remoteBranch)
if err == nil {
    // Branch/ref exists
}
```

**Pattern 3: Using output for additional validation**
```go
output, err := internal.ExecGitCommand(repoRoot, "rev-parse", "--verify", branch)
if err == nil && strings.TrimSpace(string(output)) != "" {
    // Branch exists and has valid output
}
```

### 3. Error Handling Patterns

Current error handling follows these patterns:

1. **Simple boolean check**: Most cases only check if `err == nil` to determine existence
2. **No error processing**: Errors are typically ignored, treating any error as "ref doesn't exist"
3. **Output validation**: In mergeback.go, additional validation ensures output is not empty

### 4. Existing GitManager Structure

From `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`, the GitManager has these key characteristics:

**Core Structure:**
```go
type GitManager struct {
    repo           *git.Repository
    repoPath       string
    worktreePrefix string
}
```

**Utility Function Patterns:**
- Methods follow naming convention: `func (gm *GitManager) MethodName(...) (..., error)`
- Consistent error wrapping with `enhanceGitError()` function
- Use of `ExecGitCommand()` for git operations
- Return tuple patterns like `(string, error)` or `(bool, error)`

**Recent Utility Examples:**
```go
// GetUpstreamBranch returns the upstream branch name for a given worktree path.
// Returns empty string if no upstream is set (not an error condition).
func (gm *GitManager) GetUpstreamBranch(worktreePath string) (string, error)

// GetAheadBehindCount returns the number of commits ahead and behind the upstream branch.
// Returns (0, 0, nil) if no upstream is set (not an error condition).
func (gm *GitManager) GetAheadBehindCount(worktreePath string) (int, int, error)
```

### 5. Recommended VerifyRef Utility Function Design

Based on the analysis, here's the recommended signature and behavior for a `VerifyRef` utility function:

```go
// VerifyRef verifies that a git reference (branch, tag, commit) exists and is valid.
// Returns the resolved SHA if successful, empty string and nil if ref doesn't exist (not an error).
// Returns error only for git command failures or repository issues.
func (gm *GitManager) VerifyRef(ref string) (string, error)

// VerifyRefInPath verifies that a git reference exists in a specific worktree/repository path.
// Returns the resolved SHA if successful, empty string and nil if ref doesn't exist (not an error).
// Returns error only for git command failures or repository issues.
func (gm *GitManager) VerifyRefInPath(path, ref string) (string, error)
```

This would standardize all the current `git rev-parse --verify` usage while providing:
- Consistent error handling
- Proper SHA resolution
- Clear distinction between "doesn't exist" vs "git error"
- Path flexibility for worktree operations