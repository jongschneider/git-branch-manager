# Interface Refactoring Plan for Unit Testing

This document outlines the interface extraction needed to move integration tests out of the `cmd` package and use only unit tests with interface mocking.

## Overview

Each command currently depends on the concrete `*internal.Manager` type. We need to extract interfaces for each command's specific needs and use `//go:generate go tool moq` to generate mocks.

## Command-by-Command Analysis

### 1. cmd/add.go ✅ EXAMPLE

**Current Dependencies:**
- `manager.AddWorktree()`
- `manager.GetGitManager().GetDefaultBranch()`  
- `manager.GetGitManager().BranchExists()`
- `internal.GetJiraIssues(manager)`
- `internal.GenerateBranchFromJira()`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeAdder.go . worktreeAdder
type worktreeAdder interface {
    AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
    GetDefaultBranch() (string, error)
    BranchExists(branch string) (bool, error)
    GetJiraIssues() ([]internal.JiraIssue, error)
    GenerateBranchFromJira(jiraKey string) (string, error)
}
```

### 2. cmd/pull.go

**Current Dependencies:**
- `manager.PullAllWorktrees()`
- `manager.PullWorktree(worktreeName)`
- `manager.IsInWorktree(currentPath)`
- `manager.GetAllWorktrees()`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreePuller.go . worktreePuller
type worktreePuller interface {
    PullAllWorktrees() error
    PullWorktree(worktreeName string) error
    IsInWorktree(path string) (bool, string, error)
    GetAllWorktrees() (map[string]*internal.WorktreeDetailedInfo, error)
}
```

### 3. cmd/push.go

**Current Dependencies:**
- `manager.PushAllWorktrees()`
- `manager.PushWorktree(worktreeName)`
- `manager.IsInWorktree(currentPath)`
- `manager.GetAllWorktrees()`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreePusher.go . worktreePusher
type worktreePusher interface {
    PushAllWorktrees() error
    PushWorktree(worktreeName string) error
    IsInWorktree(path string) (bool, string, error)
    GetAllWorktrees() (map[string]*internal.WorktreeDetailedInfo, error)
}
```

### 4. cmd/remove.go

**Current Dependencies:**
- `manager.GetWorktreePath(worktreeName)`
- `manager.GetWorktreeStatus(worktreePath)`
- `manager.RemoveWorktree(worktreeName)`
- `manager.GetAllWorktrees()` (for completion)

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeRemover.go . worktreeRemover
type worktreeRemover interface {
    GetWorktreePath(worktreeName string) (string, error)
    GetWorktreeStatus(worktreePath string) (*internal.GitStatus, error)
    RemoveWorktree(worktreeName string) error
    GetAllWorktrees() (map[string]*internal.WorktreeDetailedInfo, error)
}
```

### 5. cmd/switch.go

**Current Dependencies:**
- `manager.GetWorktreePath(worktreeName)`
- `manager.SetCurrentWorktree(worktreeName)`
- `manager.GetPreviousWorktree()`
- `manager.GetAllWorktrees()`
- `manager.GetSortedWorktreeNames(worktrees)`
- `manager.GetStatusIcon(gitStatus)`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeSwitcher.go . worktreeSwitcher
type worktreeSwitcher interface {
    GetWorktreePath(worktreeName string) (string, error)
    SetCurrentWorktree(worktreeName string) error
    GetPreviousWorktree() string
    GetAllWorktrees() (map[string]*internal.WorktreeDetailedInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeDetailedInfo) []string
    GetStatusIcon(gitStatus *internal.GitStatus) string
}
```

### 6. cmd/sync.go

**Current Dependencies:**
- `manager.GetSyncStatus()`
- `manager.SyncWithConfirmation(dryRun, force, confirmFunc)`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeSyncer.go . worktreeSyncer
type worktreeSyncer interface {
    GetSyncStatus() (*internal.SyncStatus, error)
    SyncWithConfirmation(dryRun, force bool, confirmFunc func(string) bool) error
}
```

### 7. cmd/validate.go

**Current Dependencies:**
- `manager.GetWorktreeMapping()`
- `manager.BranchExists(branchName)`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeValidator.go . worktreeValidator
type worktreeValidator interface {
    GetWorktreeMapping() (map[string]string, error)
    BranchExists(branchName string) (bool, error)
}
```

### 8. cmd/list.go

**Current Dependencies:**
- `manager.GetSyncStatus()`
- `manager.GetAllWorktrees()`
- `manager.GetSortedWorktreeNames(worktrees)`
- `manager.GetWorktreeMapping()`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeLister.go . worktreeLister
type worktreeLister interface {
    GetSyncStatus() (*internal.SyncStatus, error)
    GetAllWorktrees() (map[string]*internal.WorktreeDetailedInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeDetailedInfo) []string
    GetWorktreeMapping() (map[string]string, error)
}
```

### 9. cmd/hotfix.go

**Current Dependencies:**
- `manager.GetGBMConfig()`
- `manager.GetConfig().Settings.HotfixPrefix`
- `manager.AddWorktree()`
- `internal.GenerateBranchFromJira()` (via generateHotfixBranchName)
- `internal.GetJiraIssues(manager)` (for completion)

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_hotfixCreator.go . hotfixCreator
type hotfixCreator interface {
    GetGBMConfig() *internal.GBMConfig
    GetHotfixPrefix() string
    AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
    GenerateBranchFromJira(jiraKey string) (string, error)
    GetJiraIssues() ([]internal.JiraIssue, error)
    FindProductionBranch() (string, error)
}
```

### 10. cmd/mergeback.go

**Current Dependencies:**
- `manager.GetGitManager().GetRecentMergeableActivity(7)`
- `manager.GetConfig().Settings.MergebackPrefix`
- `manager.AddWorktree()`
- Multiple git operations for finding merge targets
- Complex dependency on git operations

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_mergebackCreator.go . mergebackCreator
type mergebackCreator interface {
    GetRecentMergeableActivity(days int) ([]internal.RecentActivity, error)
    GetMergebackPrefix() string
    AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
    FindMergeTargetBranch() (string, string, error) // branch, worktreeName, error
    GetGBMConfig() *internal.GBMConfig
    ValidateActivityForMergeback(activity internal.RecentActivity) bool
}
```

### 11. cmd/info.go

**Current Dependencies:**
- `manager.GetGitManager().GetWorktrees()`
- `manager.GetGitManager().GetWorktreeStatus(path)`
- `manager.GetGitManager().GetCommitHistory(path, options)`
- `manager.GetGitManager().GetFileChanges(path, options)`
- `manager.GetGitManager().GetCurrentBranchInPath(path)`
- `manager.GetGitManager().GetUpstreamBranch(path)`
- `manager.GetGitManager().GetAheadBehindCount(path)`
- `manager.GetGitManager().VerifyRefInPath(path, ref)`
- `manager.GetState().GetWorktreeBaseBranch(name)`
- `manager.GetConfig()`

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider
type worktreeInfoProvider interface {
    GetWorktrees() ([]*internal.WorktreeInfo, error)
    GetWorktreeStatus(path string) (*internal.GitStatus, error)
    GetCommitHistory(path string, options internal.CommitHistoryOptions) ([]internal.Commit, error)
    GetFileChanges(path string, options internal.FileChangeOptions) ([]internal.FileChange, error)
    GetCurrentBranchInPath(path string) (string, error)
    GetUpstreamBranch(path string) (string, error)
    GetAheadBehindCount(path string) (int, int, error)
    VerifyRefInPath(path, ref string) (bool, error)
    GetWorktreeBaseBranch(name string) (string, bool)
    GetConfig() *internal.Config
}
```

### 12. cmd/clone.go

**Current Dependencies:**
- `internal.NewManager(wd)`
- `manager.SaveConfig()`
- `manager.SaveState()`
- `manager.LoadGBMConfig(path)`
- `manager.Sync(false, false)`
- Direct git commands and file operations

**Required Interface:**
```go
//go:generate go tool moq -out ./autogen_repoCloner.go . repoCloner
type repoCloner interface {
    SaveConfig() error
    SaveState() error
    LoadGBMConfig(path string) error
    Sync(allowInvalidBranches, force bool) error
}
```

## Refactoring Complexity Analysis

### High Refactoring Effort ⚠️

#### 1. **cmd/mergeback.go** - MOST COMPLEX
**Why it's complex:**
- **Direct git command execution**: Uses `exec.Command("git", ...)` directly in multiple places
- **Complex git logic**: `hasCommitsBetweenBranches()`, `isBranchAheadOf()`, `getCommitsToMerge()`, `performMerge()`
- **File system operations**: Working directory manipulation, path resolution
- **Interactive user prompts**: Multiple confirmation dialogs that are hard to test
- **Recursive tree traversal**: `findNextMergeTargetInChain()` with complex branching logic
- **Git repo analysis**: Parsing git log output, checking branch relationships

**Additional interfaces needed:**
```go
type gitCommandExecutor interface {
    HasCommitsBetween(target, source string) (bool, error)
    IsBranchAheadOf(source, target string) (bool, error)
    GetCommitsToMerge(target, source string) ([]string, error)
    PerformMerge(worktreePath, source, target string) error
    ExecGitCommand(repoRoot string, args ...string) ([]byte, error)
}

type userInteractor interface {
    PromptConfirmation(message string) (bool, error)
    DisplayMergeInfo(commits []string, source, target string)
}

type pathResolver interface {
    GetWorkingDirectory() (string, error)
    FindGitRoot(path string) (string, error)
    GetWorktreePath(name string) string
}
```

#### 2. **cmd/clone.go** - HIGH COMPLEXITY
**Why it's complex:**
- **Heavy file system operations**: Directory creation, file copying, path manipulation
- **Direct exec.Command usage**: Multiple git commands executed directly
- **External script execution**: References `git-bare-clone.sh`
- **Working directory changes**: `os.Chdir()` calls that affect global state
- **Complex error handling**: Cleanup operations on failure

**Additional interfaces needed:**
```go
type fileSystemOperator interface {
    MkdirAll(path string, perm os.FileMode) error
    Chdir(dir string) error
    GetWorkingDir() (string, error)
    WriteFile(filename string, data []byte, perm os.FileMode) error
    CopyFile(src, dst string) error
    Stat(name string) (os.FileInfo, error)
}

type gitCloner interface {
    CloneBareRepo(repoUrl, targetDir string) error
    ConfigureRemote() error
    FetchFromOrigin() error
    GetDefaultBranch() (string, error)
    CreateWorktree(name, branch string) error
}
```

### Medium Refactoring Effort ⚡

#### 3. **cmd/hotfix.go** - MEDIUM-HIGH
**Why it needs work:**
- **Complex branch detection logic**: `findProductionBranch()` with config parsing
- **JIRA integration**: Direct calls to JIRA CLI
- **Branch name generation**: Complex logic that's currently hard to test
- **Config file parsing**: Direct file system access

#### 4. **cmd/info.go** - MEDIUM
**Why it needs work:**
- **Multiple data sources**: Git status, file changes, commit history, JIRA details
- **Direct JIRA CLI execution**: JSON parsing of jira command output
- **File system metadata**: Creation time detection
- **Complex data aggregation**: Multiple async operations combined

### Low Refactoring Effort ✅

#### 5. **cmd/add.go** - LOW (Good Example)
- Mostly uses manager methods
- Some JIRA integration but isolated
- ArgsResolver is already somewhat separated

#### 6. **cmd/pull/push/remove/switch/sync/validate/list.go** - LOW-MEDIUM
- Primarily use manager methods
- Limited direct file system access
- Straightforward command patterns

## Implementation Steps

### Phase 1: Establish Patterns (Low Complexity)
Start with commands that have straightforward dependencies to establish refactoring patterns:

1. **cmd/add.go** ✅ (Example implementation)
   - Extract `worktreeAdder` interface
   - Generate mock with `//go:generate go tool moq -out ./autogen_worktreeAdder.go . worktreeAdder`
   - Modify `ArgsResolver` and command function to use interface
   - Write unit tests using the mock

2. **cmd/pull.go, cmd/push.go, cmd/remove.go, cmd/switch.go**
   - Extract specific interfaces for each command's needs
   - Generate mocks
   - Refactor command functions to use interfaces
   - Write unit tests

### Phase 2: Build on Patterns (Medium-Low Complexity)
Apply established patterns to slightly more complex commands:

3. **cmd/sync.go, cmd/validate.go, cmd/list.go**
   - Follow established interface extraction patterns
   - Focus on manager method dependencies
   - Create comprehensive unit tests

### Phase 3: Handle Complex Data (Medium Complexity)
Tackle commands with multiple data sources and external dependencies:

4. **cmd/info.go**
   - Extract interfaces for git operations, JIRA integration
   - Mock external CLI calls (jira command)
   - Separate data aggregation logic for better testing

5. **cmd/hotfix.go**
   - Extract production branch detection logic to internal package
   - Interface JIRA integration
   - Simplify branch name generation

### Phase 4: Architectural Refactoring (High Complexity)
Save the most complex commands for last, may require architectural changes:

6. **cmd/clone.go**
   - Consider moving heavy logic to internal package
   - Create file system abstraction layer
   - Abstract git command execution
   - **Recommendation**: Move most logic to `internal.CloneManager` and keep cmd thin

7. **cmd/mergeback.go**
   - Break down into smaller, testable components
   - Move git analysis logic to internal package
   - Create user interaction abstraction
   - **Recommendation**: Major refactoring - split into multiple internal services

### Special Considerations for High-Complexity Commands

**For cmd/mergeback.go and cmd/clone.go:**
- Consider moving 80% of logic into the `internal` package where it can be tested with real git operations
- Create higher-level abstractions that hide complexity from cmd layer
- Break down monolithic functions into smaller, focused components
- Use composition over large interfaces

**Alternative approach for these commands:**
Instead of trying to mock everything, create a thin cmd layer that delegates to well-tested internal services:

```go
// cmd/mergeback.go becomes thin
func newMergebackCommand() *cobra.Command {
    return &cobra.Command{
        RunE: func(cmd *cobra.Command, args []string) error {
            service := &mergebackService{manager: manager}
            return service.ExecuteMergeback(args)
        },
    }
}

// internal/mergeback_service.go - fully testable with integration tests
type mergebackService struct {
    manager *Manager
}

func (s *mergebackService) ExecuteMergeback(args []string) error {
    // All the complex logic here, tested in internal package
}
```

## Migration Strategy

1. **Start with Phase 1** to validate the approach and establish patterns
2. **Complete Phases 2-3** to handle the majority of commands  
3. **Evaluate Phase 4** - consider architectural changes vs. interface extraction
4. **Move integration tests** from `cmd/*_test.go` to `internal/*_test.go` as interfaces are created
5. **Update helper functions** like `createInitializedManager()` to return interfaces

## Benefits

- **Faster tests**: Unit tests with mocks vs integration tests with real git repos
- **Better isolation**: Each command can be tested independently
- **Clearer dependencies**: Interfaces make dependencies explicit
- **Easier maintenance**: Changes to internal implementation don't break cmd tests
- **Better coverage**: Can test edge cases and error conditions easily with mocks

## Notes

- Some commands share similar interface methods (e.g., `GetAllWorktrees()`) - consider common interfaces
- Complex commands like `mergeback.go` may need additional refactoring to reduce dependencies
- JIRA integration should also be interfaced for testing
- File operations in `clone.go` may need additional abstraction