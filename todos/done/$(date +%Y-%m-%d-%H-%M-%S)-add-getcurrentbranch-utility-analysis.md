Based on my comprehensive research of the codebase, here is my complete analysis:

## Research Report: Git Command Patterns and GetCurrentBranch Utility Function Analysis

### 1. Current "git rev-parse --abbrev-ref HEAD" Usage Patterns

**Summary of Findings:**
- The pattern `git rev-parse --abbrev-ref HEAD` is used in **4 locations** across the codebase
- **2 locations** already use the `GetCurrentBranchInPath` utility (cmd/info.go:166 is already converted) 
- **2 locations** in internal/git.go still use direct `ExecGitCommand` calls that need refactoring

### 2. Locations Currently Using Direct rev-parse Commands

#### **Direct ExecGitCommand Usage (needs refactoring):**

**Location 1: `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:984`**
```go
// In PushWorktree method
output, err := ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")
if err != nil {
    return fmt.Errorf("failed to get current branch: %w", err)
}
currentBranch := strings.TrimSpace(string(output))
```

**Location 2: `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:1018`**
```go
// In PullWorktree method  
output, err := ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")
if err != nil {
    return fmt.Errorf("failed to get current branch: %w", err)
}
currentBranch := strings.TrimSpace(string(output))
```

### 3. Existing GitManager Interface Analysis

**Current GetCurrentBranch utilities:**
- **`GetCurrentBranch() (string, error)`** - Works on repo path (line 868)
- **`GetCurrentBranchInPath(path string) (string, error)`** - Works on specified path (line 878) ✅ Already implemented

**Location:** `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`

### 4. GitManager Method Patterns and Error Handling

**Naming Conventions:**
- Base methods operate on repo path: `GetCurrentBranch()`
- Path-specific variants add "InPath": `GetCurrentBranchInPath(path string)`
- Follow pattern: `Get[Function][Target]` (e.g., `GetUpstreamBranch`, `GetAheadBehindCount`)

**Error Handling Pattern:**
```go
func (gm *GitManager) GetCurrentBranchInPath(path string) (string, error) {
    output, err := ExecGitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")
    if err != nil {
        return "", enhanceGitError(err, "get current branch")
    }
    return strings.TrimSpace(string(output)), nil
}
```

**Consistent Error Enhancement:**
- Uses `enhanceGitError(err, "operation description")` for context-specific error messages
- Handles common git scenarios like "not a git repository", exit codes 128, etc.

### 5. Analysis of cmd/info.go:291 Context

**Key Finding:** The task mentions cmd/info.go:291, but this line is **already converted**!

**Current state (line 166):**
```go
// Get current branch (not used for base branch detection anymore)
_, err := manager.GetGitManager().GetCurrentBranchInPath(worktreePath)
```

This shows the refactoring was already completed for the info command.

### 6. Locations That Would Benefit from GetCurrentBranch Utility

**Primary Targets for Refactoring:**

1. **`internal/git.go:984` (PushWorktree method)**
   - Replace: `ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")`
   - With: `gm.GetCurrentBranchInPath(worktreePath)`

2. **`internal/git.go:1018` (PullWorktree method)**
   - Replace: `ExecGitCommand(worktreePath, "rev-parse", "--abbrev-ref", "HEAD")`  
   - With: `gm.GetCurrentBranchInPath(worktreePath)`

### 7. Integration with Existing GitManager Interface

**The GetCurrentBranchInPath utility already exists and follows established patterns:**

```go
func (gm *GitManager) GetCurrentBranchInPath(path string) (string, error) {
    output, err := ExecGitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")
    if err != nil {
        return "", enhanceGitError(err, "get current branch")
    }
    return strings.TrimSpace(string(output)), nil
}
```

**Key Interface Characteristics:**
- Uses `ExecGitCommand` utility for execution consistency
- Applies `enhanceGitError` for standardized error handling
- Returns `(string, error)` signature matching other utilities
- Handles path parameter for worktree-specific operations

### 8. Testing Infrastructure

**Existing Test Coverage:**
- **File:** `/Users/jschneider/code/scratch/worktree-manager/internal/git_test.go`
- **Test:** `TestGitManager_GetCurrentBranchInPath` (starting line 14)
- **Pattern:** Uses `testutils.GitTestRepo` for isolated testing
- **Assertions:** Uses `testify` (`require.NoError`, `assert.Equal`)

### 9. Recommended Implementation Approach

Since `GetCurrentBranchInPath` already exists and follows established patterns, the task is to:

1. **Replace direct ExecGitCommand calls** in `PushWorktree` and `PullWorktree` methods
2. **Use existing utility function** instead of creating new ones
3. **Follow existing error handling** patterns with `enhanceGitError`
4. **Maintain current test coverage** (tests already exist)

### 10. Summary of Benefits

**Consistency:** All current branch retrieval will use the same utility function
**Error Handling:** Centralized error enhancement for better user experience  
**Maintainability:** Single point of change for current branch logic
**Testing:** Existing comprehensive test coverage ensures reliability

The codebase already has the necessary infrastructure in place. The remaining work is simply to replace the 2 direct `ExecGitCommand` calls with the existing `GetCurrentBranchInPath` utility method.