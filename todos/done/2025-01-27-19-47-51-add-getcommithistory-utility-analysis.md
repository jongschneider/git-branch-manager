## Analysis Results

Based on my analysis of the worktree-manager codebase, here's what I found regarding git log usage and the opportunity for deduplication:

### Git Log Patterns Found

I identified **4 distinct git log patterns** currently used in the codebase:

1. **Recent Commits (cmd/info.go:161)**:
   ```bash
   git log -5 --oneline --format=%H|%s|%an|%ct
   ```
   - Purpose: Get recent commits for worktree info display
   - Format: `hash|message|author|timestamp`

2. **Mergeback Analysis (internal/mergeback.go:175)**:
   ```bash
   git log origin/target..origin/source --format=%H|%s|%an|%ae|%ct
   ```
   - Purpose: Find commits that need to be merged back between branches
   - Format: `hash|message|author|email|timestamp`

3. **Merge Commits Detection (internal/git.go:948)**:
   ```bash
   git log --merges --since=N.days.ago --pretty=format:%H|%an|%at|%s
   ```
   - Purpose: Find recent merge commits for mergeback activity detection
   - Format: `hash|author|unixtime|message`

4. **Hotfix Activity Detection (internal/git.go:1006)**:
   ```bash
   git log --all --since=N.days.ago --pretty=format:%H|%an|%at|%s|%D --grep=hotfix
   ```
   - Purpose: Find recent hotfix commits for mergeback activity detection
   - Format: `hash|author|unixtime|message|refs`

### Existing Commit-Related Utilities

The codebase already has these commit-related utility functions in `internal/git.go`:
- `GetCommitHash(ref string)` - Get commit hash for a reference
- `GetCommitHashInPath(path, ref string)` - Get commit hash for a reference in specific path
- `getRecentMergeCommits(since string)` - Get recent merge commits (private)
- `extractMergeBranches(commitHash string)` - Extract source/target from merge commit (private)

### Recommended GetCommitHistory Function

Based on the patterns I found, here's the suggested function signature and implementation approach:

```go
// CommitHistoryOptions defines options for retrieving commit history
type CommitHistoryOptions struct {
    // Limit number of commits (equivalent to -N flag)
    Limit int
    
    // Range specification (e.g., "origin/main..origin/feature", "HEAD~5..HEAD")
    Range string
    
    // Since timestamp or relative time (e.g., "7.days.ago", "2023-01-01")
    Since string
    
    // Additional git log flags
    MergesOnly bool  // --merges
    AllBranches bool // --all
    GrepPattern string // --grep=pattern
    
    // Format specification - if empty, uses a sensible default
    CustomFormat string
}

// GetCommitHistory retrieves commit history with flexible options
// If path is empty, uses repository root. Returns commits in chronological order (newest first).
func (gm *GitManager) GetCommitHistory(path string, options CommitHistoryOptions) ([]CommitInfo, error) {
    // Implementation would build git log command based on options
    // Default format: %H|%s|%an|%ae|%ct (hash|message|author|email|timestamp)
}
```

### Centralization Benefits

Creating this utility would:

1. **Eliminate Duplication**: Replace 4 different git log command constructions
2. **Standardize Parsing**: Use consistent pipe-delimited format parsing
3. **Improve Error Handling**: Centralize git command error enhancement
4. **Enhance Maintainability**: Single place to modify git log behavior
5. **Add Flexibility**: Support various use cases through options struct

### Files That Would Benefit

The following files contain git log calls that could be refactored to use the new utility:

- `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go:161` (getRecentCommits function)
- `/Users/jschneider/code/scratch/worktree-manager/internal/mergeback.go:175` (getCommitsNeedingMergeBack function)
- `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:948` (getRecentMergeCommits function)
- `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:1006` (getRecentHotfixActivity function)

This utility would fit well with the existing `ExecGitCommand` pattern already established in the codebase and provide a higher-level abstraction for commit history operations.