# Git Command Deduplication Gameplan

## Overview

This document outlines the step-by-step plan to eliminate Git command duplication across the gbm codebase by creating centralized utility functions and refactoring existing code to use them.

## Current State

Based on the analysis report, the codebase has good centralization but still contains duplication in three main areas:

1. **Branch Status Checking** (HIGH duplication) - 8+ occurrences
2. **Repository Introspection** (MEDIUM duplication) - 15+ occurrences  
3. **Info Command Git Calls** (MEDIUM duplication) - 10 direct calls

## Implementation Strategy

### Phase 1: High Priority - Branch Status Utilities

**Goal**: Eliminate duplicate branch status checking patterns

**New Utility Functions** (add to `internal/git.go`):
```go
// GetCurrentBranch returns the current branch name for a given worktree path
func (gm *GitManager) GetCurrentBranch(worktreePath string) (string, error)

// GetUpstreamBranch returns the upstream branch name for a given worktree path  
func (gm *GitManager) GetUpstreamBranch(worktreePath string) (string, error)

// GetAheadBehindCount returns ahead/behind commit counts vs upstream
func (gm *GitManager) GetAheadBehindCount(worktreePath string) (int, int, error)
```

**Files to Refactor**:
- `cmd/info.go:291,300,309` - Replace direct git calls
- `cmd/sync_test.go:210,403,410,875,882` - Use new utilities
- `internal/git.go:573,635,666` - Update existing similar functions

### Phase 2: Medium Priority - Repository Introspection

**Goal**: Consolidate `git rev-parse` operations

**New Utility Functions** (add to `internal/git.go`):
```go
// VerifyRef checks if a git reference exists
func (gm *GitManager) VerifyRef(ref string) error

// GetCommitHash returns the commit hash for a given reference
func (gm *GitManager) GetCommitHash(ref string) (string, error)
```

**Files to Refactor**:
- `cmd/info.go:338` - Replace direct rev-parse calls
- `internal/git_add.go:48,54,89` - Use new utilities

### Phase 3: Low Priority - Info Command Extraction

**Goal**: Extract remaining specialized Git operations from info command

**New Utility Functions** (add to `internal/git.go`):
```go
// GetCommitHistory returns recent commit information
func (gm *GitManager) GetCommitHistory(worktreePath string, count int) ([]CommitInfo, error)

// GetFileChanges returns file change statistics (staged and unstaged)
func (gm *GitManager) GetFileChanges(worktreePath string) ([]FileChange, error)
```

**Files to Refactor**:
- `cmd/info.go:161-338` - Extract git log and git diff calls

## Implementation Guidelines

### Adding New Utility Functions

1. **Location**: All new functions go in `internal/git.go`
2. **Error Handling**: Use `enhanceGitError()` for consistent error messages
3. **Function Signature**: Include `worktreePath string` parameter for context
4. **Return Types**: Use existing structs (`CommitInfo`, `FileChange`) or create new ones

### Refactoring Existing Code

1. **One Function at a Time**: Implement and test each utility function individually
2. **Replace Incrementally**: Update one file at a time to use new utilities
3. **Test After Each Change**: Run relevant tests to ensure no regressions
4. **Keep Interfaces Consistent**: Maintain existing function signatures in manager methods

### Testing Strategy

1. **Unit Tests**: Add tests for each new utility function
2. **Integration Tests**: Ensure existing tests continue to pass
3. **Manual Testing**: Verify common workflows still work correctly

## Success Criteria

- [ ] All branch status checking uses centralized utilities
- [ ] All rev-parse operations use centralized utilities  
- [ ] Info command uses extracted utility functions
- [ ] Full test suite passes
- [ ] No direct `exec.Command("git", ...)` calls outside of utility functions
- [ ] Code is more maintainable and testable

## Risk Mitigation

1. **Incremental Approach**: Small, testable changes reduce risk of breaking existing functionality
2. **Existing Test Coverage**: Leverage existing test suite to catch regressions
3. **Backward Compatibility**: Keep existing manager method interfaces unchanged
4. **Rollback Plan**: Each change can be reverted independently if issues arise

## Expected Benefits

1. **Reduced Code Duplication**: Eliminate 20+ duplicate Git command patterns
2. **Improved Maintainability**: Single place to update Git command logic
3. **Better Error Handling**: Consistent error messages across all Git operations
4. **Enhanced Testability**: Centralized functions are easier to mock and test
5. **Future-Proofing**: New commands can easily reuse existing utilities