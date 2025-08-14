## Analysis of cmd/info.go Structure

### **Current Structure Overview**

The `cmd/info.go` file contains the following key functions:
- `newInfoCommand()` - Creates the cobra command
- `runInfoCommand()` - Main entry point (lines 36-64)
- `getWorktreeInfo()` - Core logic for gathering worktree information (lines 67-145)
- `displayWorktreeInfo()` - Handles rendering (lines 147-154)
- `getWorktreeCreationTime()` - File system interaction (lines 156-162)
- `getBaseBranchInfo()` - Git operations for branch information (lines 164-224)
- `getJiraTicketDetails()` - External JIRA CLI integration (lines 276-377)
- Various JSON parsing structs and helper functions

### **Manager Methods Being Called**

From the analysis of `getWorktreeInfo()` function, I identified these Manager method calls:

1. **Line 68**: `manager.GetGitManager()` - Gets the GitManager instance
2. **Line 70**: `gitManager.GetWorktrees()` - Gets all worktrees
3. **Line 89**: `gitManager.GetWorktreeStatus()` - Gets git status for specific worktree
4. **Line 102-104**: `manager.GetGitManager().GetCommitHistory()` - Gets recent commits
5. **Line 110-113**: `manager.GetGitManager().GetFileChanges()` - Gets modified files
6. **Line 62**: `manager.GetConfig()` - Gets configuration

**In `getBaseBranchInfo()` function (lines 164-224):**
7. **Line 166**: `manager.GetGitManager().GetCurrentBranchInPath()`
8. **Line 173**: `manager.GetGitManager().GetUpstreamBranch()`
9. **Line 179**: `manager.GetGitManager().GetAheadBehindCount()`
10. **Line 188**: `manager.GetState().GetWorktreeBaseBranch()`
11. **Line 196**: `manager.GetConfig().Settings.CandidateBranches`
12. **Line 202**: `manager.GetGitManager().VerifyRefInPath()`

### **External System Interactions**

1. **Git operations** (via GitManager):
   - Getting worktrees, status, commit history, file changes
   - Branch verification, upstream tracking, ahead/behind counts
   - Direct git command execution (line 208-214 for merge-base)

2. **File system operations**:
   - `os.Getwd()` (line 39-43)
   - `os.Stat()` (line 157 for creation time)

3. **JIRA CLI integration**:
   - `exec.LookPath("jira")` (line 278) 
   - `exec.Command("jira", "issue", "view", ...)` (line 283)

### **Comparison with cmd/pull.go and cmd/push.go Patterns**

Both `cmd/pull.go` and `cmd/push.go` follow a similar interface extraction pattern:

**cmd/pull.go (lines 16-21)**:
```go
type worktreePuller interface {
    PullAllWorktrees() error
    PullWorktree(worktreeName string) error
    IsInWorktree(currentPath string) (bool, string, error)
    GetAllWorktrees() (map[string]*WorktreeListInfo, error)
}
```

**cmd/push.go (lines 16-21)**:
```go
type worktreePusher interface {
    PushAllWorktrees() error
    PushWorktree(worktreeName string) error
    IsInWorktree(currentPath string) (bool, string, error)
    GetAllWorktrees() (map[string]*WorktreeListInfo, error)
}
```

Both interfaces:
- Are generated with `//go:generate` comments for mock generation
- Use specific operation-focused names (`worktreePuller`, `worktreePusher`)
- Extract only the Manager methods needed for their specific operations
- Follow a consistent pattern of validation + operation execution

### **Recommended worktreeInfoProvider Interface**

Based on the analysis, the `worktreeInfoProvider` interface should include:

```go
type worktreeInfoProvider interface {
    GetGitManager() *internal.GitManager
    GetConfig() *internal.Config
    GetState() *internal.State
}
```

However, this might be too broad. Looking at the specific GitManager methods being called, a more focused approach would be:

```go
type worktreeInfoProvider interface {
    GetWorktreeInfo(worktreeName string) (*internal.WorktreeInfoData, error)
    GetConfig() *internal.Config
}
```

### **No Integration Tests Found**

There is no `cmd/info_test.go` file in the codebase, so there are no existing integration tests to move to the internal package.

### **Key Insights for Interface Extraction**

1. **Complex Dependencies**: The info command has more complex dependencies than pull/push, requiring git operations, file system access, and external CLI calls
2. **Data Aggregation**: Unlike pull/push which perform single operations, info aggregates data from multiple sources
3. **Optional Operations**: Many operations in `getWorktreeInfo()` are optional (with error handling that continues execution)
4. **Direct Git Commands**: Some functionality uses direct git command execution rather than going through GitManager methods

The interface extraction will need to balance testability with the complexity of the data gathering operations.

Based on my analysis of the internal package structure, here are my findings about where the worktreeInfoProvider business logic should be moved:

## Analysis Summary

### 1. **Manager struct and architecture** (`internal/manager.go`)
The `Manager` struct is the central coordinator that contains:
- **Dependencies**: `*Config`, `*State`, `*GitManager`, `*GBMConfig`
- **Core methods**: Already has `GetWorktreeList()` and `GetAllWorktrees()` methods that aggregate worktree information
- **Data structures**: Defines `WorktreeListInfo` struct used for listing operations

### 2. **GitManager operations** (`internal/git.go`)
The `GitManager` provides lower-level Git operations:
- **Worktree operations**: `GetWorktrees()`, `GetWorktreeStatus()`, `CreateWorktree()`, etc.
- **Git status**: `GetWorktreeStatus()` returns `*GitStatus` with detailed status information
- **Complex data structures**: Defines `WorktreeInfoData`, `CommitInfo`, `BranchInfo`, etc. for comprehensive worktree information
- **Git history methods**: `GetCommitHistory()`, `GetFileChanges()`, `GetRecentMergeableActivity()`

### 3. **Existing data aggregation patterns**
- **Manager methods**: `GetWorktreeList()` (lines 418-457) and `GetAllWorktrees()` (lines 473-527) already aggregate data from Git operations and configuration
- **Status aggregation**: `GetSyncStatus()` method combines Git state with configuration to determine sync status
- **Info rendering**: `InfoRenderer` in `internal/info_renderer.go` handles complex data presentation

### 4. **No existing info-related services**
- No dedicated "Service" classes found in internal package
- No existing InfoService or similar abstraction
- The pattern follows a simpler Manager + GitManager architecture

### 5. **Configuration and State management** 
- **Config** (`internal/config.go`): Handles settings, icons, JIRA config, file copy rules
- **State** (`internal/state.go`): Manages runtime state like tracked worktrees, current/previous worktree, base branch mappings

### 6. **Current interface extraction pattern**
In `cmd/list.go`, there's already an interface extraction pattern:
```go
type worktreeLister interface {
    GetSyncStatus() (*internal.SyncStatus, error)
    GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeListInfo) []string
    GetWorktreeMapping() (map[string]string, error)
}
```

## Recommendations

### **Option 1: Extend Manager (Recommended)**
Move the worktreeInfoProvider logic into the existing `Manager` struct by adding new methods:
- `GetWorktreeInfoData(worktreeName string) (*WorktreeInfoData, error)` 
- This leverages existing `GitManager` methods and data structures
- Follows the established architecture pattern
- `WorktreeInfoData` already exists in `internal/git.go` with all needed fields

### **Option 2: Create InfoService (Alternative)**
Create a new `internal/info_service.go` with:
- `type InfoService struct { manager *Manager, gitManager *GitManager }`
- Methods for aggregating comprehensive worktree information
- This would be a new architectural pattern but provides better separation of concerns

### **Option 3: Extend GitManager**
Add info aggregation methods directly to `GitManager`, but this violates single responsibility since GitManager should focus on Git operations.

## Specific Implementation Path

The `Manager` struct already has the perfect foundation since:
1. **Data aggregation precedent**: `GetWorktreeList()` already aggregates Git status, config, and path information
2. **Dependencies available**: Has access to `gitManager`, `config`, `state`, and `gbmConfig`
3. **Existing data structures**: `WorktreeInfoData` in `git.go` contains all the fields needed (Name, Path, Branch, GitStatus, BaseInfo, Commits, ModifiedFiles, JiraTicket)
4. **Consistent interface pattern**: The `worktreeLister` interface extraction pattern can be extended for the info provider

The business logic should be moved to `Manager` with a new method like:
```go
func (m *Manager) GetWorktreeInfoData(worktreeName string) (*WorktreeInfoData, error)
```

This approach maintains architectural consistency and leverages existing patterns and data structures effectively.