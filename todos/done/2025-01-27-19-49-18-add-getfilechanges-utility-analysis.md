## Analysis Results

Based on my analysis of the worktree-manager codebase, here's what I found regarding git diff usage and the opportunity for deduplication:

### Git Diff Patterns Found

I identified **2 distinct git diff patterns** currently used in `cmd/info.go`:

1. **Unstaged Changes (cmd/info.go:197)**:
   ```bash
   git diff --numstat
   ```
   - Purpose: Get unstaged file changes with addition/deletion counts
   - Format: `additions deletions filename`

2. **Staged Changes (cmd/info.go:238)**:
   ```bash
   git diff --cached --numstat
   ```
   - Purpose: Get staged file changes with addition/deletion counts
   - Format: `additions deletions filename`

### Current Implementation Analysis

Both patterns are used in the `getModifiedFiles()` function in `cmd/info.go` with nearly identical parsing logic:

```go
// Lines 197-235: Unstaged changes
cmd := exec.Command("git", "diff", "--numstat")
// Parse output into FileChange structs

// Lines 238-270: Staged changes  
cmd = exec.Command("git", "diff", "--cached", "--numstat")
// Parse output into FileChange structs (duplicate logic)
```

### Existing FileChange Struct

The `FileChange` struct already exists in `internal/git.go`:

```go
type FileChange struct {
    Path      string
    Status    string
    Additions int
    Deletions int
}
```

### Recommended GetFileChanges Function

Based on the patterns I found, here's the suggested function signature and implementation approach:

```go
// FileChangeOptions defines options for retrieving file changes
type FileChangeOptions struct {
    // Include staged changes (--cached)
    Staged bool
    
    // Include unstaged changes (default: true)
    Unstaged bool
    
    // Show only names (--name-only)
    NamesOnly bool
    
    // Show status (--name-status) 
    ShowStatus bool
    
    // Custom diff options
    ExtraArgs []string
}

// GetFileChanges retrieves file changes with flexible options
// If path is empty, uses repository root. Returns all requested changes.
func (gm *GitManager) GetFileChanges(path string, options FileChangeOptions) ([]FileChange, error) {
    // Implementation would:
    // 1. Build git diff command based on options
    // 2. Execute command and parse output
    // 3. Return structured FileChange objects
}
```

### Implementation Benefits

Creating this utility would:

1. **Eliminate Duplication**: Replace duplicate git diff parsing logic in `getModifiedFiles()`
2. **Standardize Parsing**: Use consistent numstat output parsing
3. **Improve Error Handling**: Centralize git command error enhancement
4. **Add Flexibility**: Support both staged and unstaged changes in one call
5. **Enhance Maintainability**: Single place to modify git diff behavior

### Files That Would Benefit

The following files contain git diff calls that could be refactored:

- `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go:197` (unstaged changes)
- `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go:238` (staged changes)

### Proposed Implementation Strategy

1. Create `FileChangeOptions` struct for flexible queries
2. Implement `GetFileChanges` utility function in `GitManager`
3. Add helper function `parseNumstatOutput` for parsing git diff --numstat output
4. Refactor `getModifiedFiles()` in `cmd/info.go` to use the new utility
5. Add comprehensive unit tests

This utility would fit well with the existing `ExecGitCommand` pattern and provide a higher-level abstraction for file change operations, similar to the recently implemented `GetCommitHistory` function.