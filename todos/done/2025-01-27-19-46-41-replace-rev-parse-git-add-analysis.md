Based on my analysis of the codebase, I can provide you with a comprehensive report on the Git command deduplication opportunities. Here's what I found:

## Analysis of rev-parse calls in internal/git_add.go

Looking at the specific lines you mentioned in `internal/git_add.go`:

### Line 48: `ExecGitCommand(gm.repoPath, "merge-base", branchName, baseBranch)`
This is actually a `merge-base` command, not `rev-parse`. This command finds the common ancestor between two branches.

### Line 54: `gm.GetCommitHash(baseBranch)`
This calls a method that internally uses:
```go
output, err := ExecGitCommand(gm.repoPath, "rev-parse", ref)
```
**This is already using a centralized utility function!**

### Line 89: `Remote(branchName)` and line 90: `gm.VerifyRef(remoteBranch)`
The `VerifyRef` method internally uses:
```go
_, err := ExecGitCommand(gm.repoPath, "rev-parse", "--verify", ref)
```
**This is also already using a centralized utility function!**

## Existing Utility Functions in internal/git.go

The codebase already has excellent centralized Git command utilities:

### Core Execution Functions:
1. **`ExecGitCommand(dir string, args ...string) ([]byte, error)`** - Main utility for git commands with output capture
2. **`execGitCommandRun(dir string, args ...string) error`** - For commands without output capture
3. **`ExecGitCommandCombined(dir string, args ...string) ([]byte, error)`** - Returns combined stdout/stderr
4. **`ExecGitCommandInteractive(dir string, args ...string) error`** - For interactive commands

### Specialized rev-parse Utilities:
1. **`VerifyRef(ref string) (bool, error)`** - Uses `rev-parse --verify`
2. **`VerifyRefInPath(path, ref string) (bool, error)`** - Uses `rev-parse --verify` in specific path
3. **`GetCommitHash(ref string) (string, error)`** - Uses `rev-parse` to get commit hash
4. **`GetCommitHashInPath(path, ref string) (string, error)`** - Gets commit hash in specific path
5. **`GetCurrentBranch() (string, error)`** - Uses `rev-parse --abbrev-ref HEAD`
6. **`GetCurrentBranchInPath(path string) (string, error)`** - Gets current branch in specific path
7. **`GetUpstreamBranch(worktreePath string) (string, error)`** - Uses `rev-parse --abbrev-ref @{upstream}`

## Patterns Found Across the Codebase

After analyzing all rev-parse usage, I found these patterns that could benefit from further centralization:

### 1. Repository Root Detection (Currently in FindGitRoot)
Multiple places use:
```go
cmd := exec.Command("git", "rev-parse", "--git-dir")
cmd := exec.Command("git", "rev-parse", "--show-toplevel") 
cmd := exec.Command("git", "rev-parse", "--is-bare-repository")
```

### 2. Direct exec.Command Usage (Should be migrated)
Some places still use `exec.Command` directly instead of the centralized utilities:
- Lines 134, 151, 165 in `FindGitRoot` function
- Lines 202, 216, 229 in subdirectory detection
- Line 819 in `IsInWorktree` method

## Recommendations for Deduplication

### Immediate Actions:
1. **Migrate FindGitRoot function** to use `ExecGitCommand` instead of direct `exec.Command`
2. **Migrate IsInWorktree method** to use `ExecGitCommand`
3. **Add specialized utility functions** for common rev-parse patterns:

```go
// GetGitDir returns the git directory path
func (gm *GitManager) GetGitDir() (string, error) {
    output, err := ExecGitCommand(gm.repoPath, "rev-parse", "--git-dir")
    if err != nil {
        return "", enhanceGitError(err, "get git dir")
    }
    return strings.TrimSpace(string(output)), nil
}

// GetTopLevel returns the top-level directory of the working tree
func (gm *GitManager) GetTopLevel() (string, error) {
    output, err := ExecGitCommand(gm.repoPath, "rev-parse", "--show-toplevel")
    if err != nil {
        return "", enhanceGitError(err, "get top level")
    }
    return strings.TrimSpace(string(output)), nil
}

// IsBareRepository checks if the repository is a bare repository
func (gm *GitManager) IsBareRepository() (bool, error) {
    output, err := ExecGitCommand(gm.repoPath, "rev-parse", "--is-bare-repository")
    if err != nil {
        return false, enhanceGitError(err, "check bare repository")
    }
    return strings.TrimSpace(string(output)) == "true", nil
}
```

### Key Findings:
- **The `internal/git_add.go` file is already well-optimized** - it uses the existing centralized utilities
- **The main deduplication opportunity** is in the `FindGitRoot` function and `IsInWorktree` method
- **The codebase has excellent centralized Git utilities** that should be used consistently throughout

The specific lines you mentioned (48, 54, 89) are actually already using centralized utilities or are not rev-parse commands, so they don't need immediate changes for your deduplication project.