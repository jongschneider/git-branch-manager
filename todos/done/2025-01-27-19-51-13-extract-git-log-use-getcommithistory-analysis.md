## Analysis Results

Based on my analysis of cmd/info.go, here's what I found regarding the git log usage that needs to be replaced:

### Current Implementation

**Function**: `getRecentCommits(worktreePath string, count int)` at line 160

**Git Command Pattern** (line 161):
```bash
git log -5 --oneline --format=%H|%s|%an|%ct
```

**Current Format**: `%H|%s|%an|%ct` (hash|message|author|timestamp)

**Usage**: Called once at line 103: `commits, err := getRecentCommits(targetWorktree.Path, 5)`

### Current Parsing Logic

The existing function parses the output expecting exactly 4 fields:
- `parts[0]` - Hash
- `parts[1]` - Message  
- `parts[2]` - Author
- `parts[3]` - Timestamp (Unix)

Missing fields compared to GetCommitHistory's default format:
- Email field
- Refs field

### Refactoring Strategy

Since we have a newly implemented `GetCommitHistory` utility that supports flexible options, we can replace the `getRecentCommits` function entirely by:

1. **Direct Replacement**: Update the call at line 103 to use `GetCommitHistory` directly
2. **Remove Function**: Delete the entire `getRecentCommits` function (lines 160-193)
3. **Compatible Format**: Use custom format matching existing expectation or adapt to new enhanced CommitInfo struct

### Implementation Options

**Option 1: Minimal Change (Backward Compatible)**
```go
// Replace line 103
commits, err := manager.GetGitManager().GetCommitHistory(targetWorktree.Path, internal.CommitHistoryOptions{
    Limit: 5,
    CustomFormat: "%H|%s|%an|%ct", // Match existing format
})
```

**Option 2: Enhanced (Use Full CommitInfo)**
```go
// Replace line 103 - use default format with all fields
commits, err := manager.GetGitManager().GetCommitHistory(targetWorktree.Path, internal.CommitHistoryOptions{
    Limit: 5,
})
// No custom format needed - uses default enhanced format
```

### Benefits of Refactoring

1. **Eliminate Duplication**: Remove 34 lines of duplicate git log parsing logic
2. **Consistent Error Handling**: Use centralized `enhanceGitError()` 
3. **Enhanced Features**: Get access to email and refs fields for future use
4. **Better Maintainability**: Single place to modify git log behavior
5. **Proven Testing**: Leverage comprehensive unit tests already written for GetCommitHistory

### Files Affected

- `cmd/info.go:103` - Update function call
- `cmd/info.go:160-193` - Remove getRecentCommits function entirely

This refactoring represents a clean elimination of duplicate code while maintaining identical functionality.