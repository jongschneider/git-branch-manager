# Extract worktreeInfoProvider interface for cmd/info.go
**Status:** AwaitingCommit
**Agent PID:** 19462

## Original Todo
**Extract worktreeInfoProvider interface for cmd/info.go**
- Create interface with: GetWorktrees(), GetWorktreeStatus(), GetCommitHistory(), GetFileChanges(), GetCurrentBranchInPath(), GetUpstreamBranch(), GetAheadBehindCount(), VerifyRefInPath(), GetWorktreeBaseBranch(), GetConfig()
- Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider`
- Mock external JIRA CLI calls
- Separate data aggregation logic for better testing
- Write comprehensive unit tests
- Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

## Description
Extract the `worktreeInfoProvider` interface from `cmd/info.go` following the established pattern from `cmd/list.go`. The interface will wrap Manager methods to enable unit testing with mocks. Currently `cmd/info.go` makes direct calls to `manager.GetGitManager().SomeMethod()` and external JIRA CLI commands, making it difficult to test in isolation.

The refactoring will:
1. Create wrapper methods on Manager for GitManager operations (e.g., `manager.GetWorktreeCommitHistory()` wrapping `manager.GetGitManager().GetCommitHistory()`)
2. Move JIRA CLI interactions to Manager methods, leveraging the existing `internal/jira.go` client
3. Create the `worktreeInfoProvider` interface that exposes only the Manager methods needed by `cmd/info.go`
4. Refactor `cmd/info.go` to use the interface, enabling fast unit tests with mocks

## Implementation Plan
- [x] Add wrapper methods to Manager for GitManager operations used by cmd/info.go (internal/manager.go)
  - `GetWorktreeCommitHistory(worktreePath string, limit int) ([]CommitInfo, error)`
  - `GetWorktreeFileChanges(worktreePath string) ([]string, error)`
  - `GetWorktreeCurrentBranch(worktreePath string) (string, error)`
  - `GetWorktreeUpstreamBranch(worktreePath string) (string, error)`  
  - `GetWorktreeAheadBehindCount(worktreePath string, upstream string) (int, int, error)`
  - `VerifyWorktreeRef(ref string, worktreePath string) error`
- [x] Add JIRA interaction methods to Manager (internal/manager.go)
  - `GetJiraTicketDetails(jiraKey string) (*JiraIssue, error)` - wraps existing internal/jira.go functions
  - `IsJiraCliAvailable() bool` - checks if JIRA CLI is installed
- [x] Create worktreeInfoProvider interface in cmd/info.go
  - Include all Manager methods needed: `GetWorktrees()`, `GetWorktreeStatus()`, `GetConfig()`, `GetState()`, plus new wrapper methods
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider`
- [x] Refactor cmd/info.go functions to use interface instead of Manager directly
  - Update `getWorktreeInfo()`, `getBaseBranchInfo()`, `getJiraTicketDetails()` to accept interface
  - Remove direct GitManager access and JIRA CLI calls
  - Use interface wrapper methods instead
- [x] Write comprehensive unit tests for cmd/info.go (cmd/info_test.go)
  - Test `getWorktreeInfo()` with mocked data
  - Test `getBaseBranchInfo()` with various branch scenarios  
  - Test `getJiraTicketDetails()` with mocked JIRA responses
  - Test error handling and edge cases
- [x] Run validation: `go build && go test ./cmd/...` to ensure fast unit tests

## Notes
[Implementation notes]