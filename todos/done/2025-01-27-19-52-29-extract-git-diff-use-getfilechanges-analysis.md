## Analysis Results

Based on my analysis of cmd/info.go, here's what I found regarding the git diff usage that needs to be replaced:

### Current Implementation

**Function**: `getModifiedFiles(worktreePath string)` at line 162

**Git Command Patterns**:
1. **Line 164**: `git diff --numstat` (unstaged changes)
2. **Line 205**: `git diff --cached --numstat` (staged changes)

**Usage**: Called once at line 111: `modifiedFiles, err := getModifiedFiles(targetWorktree.Path)`

### Current Logic Analysis

The existing function:
1. **Gets unstaged changes** using `git diff --numstat`
2. **Parses output** into FileChange structs with additions/deletions/path
3. **Gets staged changes** using `git diff --cached --numstat` 
4. **Merges results** - combines staged and unstaged changes for same files
5. **Returns combined list** of all file changes

### Duplicate Logic Found

The current `getModifiedFiles` function contains 92 lines (162-254) of git diff parsing logic that **exactly matches** what our new `GetFileChanges` utility does:

- Same `--numstat` format parsing
- Same additions/deletions extraction
- Same status determination logic ("A", "M", "D")
- Same error handling patterns

### Refactoring Strategy

Since we have the `GetFileChanges` utility that supports both staged and unstaged changes, we can replace the entire `getModifiedFiles` function:

**Option 1: Direct Replacement**
```go
// Replace line 111
modifiedFiles, err := manager.GetGitManager().GetFileChanges(targetWorktree.Path, internal.FileChangeOptions{
    Staged:   true,
    Unstaged: true,
})
```

**Option 2: Keep Function but Simplify**
```go
func getModifiedFiles(worktreePath string, manager *internal.Manager) ([]internal.FileChange, error) {
    return manager.GetGitManager().GetFileChanges(worktreePath, internal.FileChangeOptions{
        Staged:   true,
        Unstaged: true,
    })
}
```

### Key Differences to Address

The current function has **merging logic** that combines staged and unstaged changes for the same file (lines 223-233). However, our GetFileChanges utility returns **separate entries** for staged and unstaged changes of the same file.

**Solution**: Since the info command likely wants to show all changes (both staged and unstaged), the separate entries approach from GetFileChanges is actually better for comprehensive display.

### Benefits of Refactoring

1. **Eliminate Duplication**: Remove 92 lines of duplicate git diff parsing logic
2. **Consistent Error Handling**: Use centralized `enhanceGitError()`
3. **Better Separation**: Clearly distinguish between staged and unstaged changes
4. **Enhanced Features**: Get access to improved parsing and format options
5. **Proven Testing**: Leverage comprehensive unit tests already written for GetFileChanges

### Files Affected

- `cmd/info.go:111` - Update function call
- `cmd/info.go:162-254` - Remove or simplify getModifiedFiles function

### Implementation Recommendation

Use **Option 1** (direct replacement) since:
- Eliminates the most duplicate code
- GetFileChanges is well-tested and handles all edge cases
- Provides clearer separation of staged vs unstaged changes
- Maintains manager pattern access through existing `manager` variable