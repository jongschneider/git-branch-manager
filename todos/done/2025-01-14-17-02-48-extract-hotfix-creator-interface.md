# Extract hotfixCreator interface for cmd/hotfix.go
**Status:** Done
**Agent PID:** 24472

## Original Todo
**Extract hotfixCreator interface for cmd/hotfix.go**
- Create interface with: GetGBMConfig(), GetHotfixPrefix(), AddWorktree(), GenerateBranchFromJira(), GetJiraIssues(), FindProductionBranch()
- Add mock generation: `//go:generate go tool moq -out ./autogen_hotfixCreator.go . hotfixCreator`
- Extract production branch detection logic to internal package
- Interface JIRA integration
- Write unit tests with mocks
- Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

## Description
Extract the `hotfixCreator` interface from `cmd/hotfix.go` following the established pattern from other commands. The interface will enable unit testing by mocking Manager dependencies, while preserving the complex production branch detection and JIRA integration features.

The refactoring will:
1. Extract production branch detection logic to Manager method (enables testing of sophisticated deployment chain analysis)
2. Create wrapper methods on Manager for GitManager operations and JIRA functions
3. Create the `hotfixCreator` interface exposing only the methods needed by cmd/hotfix.go
4. Refactor hotfix command functions to use the interface, enabling fast unit tests with mocks
5. Move any integration tests to internal package

## Implementation Plan
- [x] Extract production branch detection logic to Manager method (internal/manager.go)
  - Move `findProductionBranch()` logic to `Manager.FindProductionBranch() (string, error)`
  - Add wrapper method `Manager.GetDefaultBranch() (string, error)` for GitManager
- [x] Add JIRA wrapper methods to Manager (internal/manager.go)
  - `GetJiraIssues() ([]JiraIssue, error)` wrapping existing internal/jira.go functions
  - `GenerateBranchFromJira(jiraKey string) (string, error)` wrapping existing function
- [x] Create hotfixCreator interface in cmd/hotfix.go
  - Include methods: `AddWorktree()`, `GetConfig()`, `GetGBMConfig()`, `FindProductionBranch()`, `GetJiraIssues()`, `GenerateBranchFromJira()`, `GetDefaultBranch()`
  - Add mock generation: `//go:generate go run github.com/matryer/moq@latest -out ./autogen_hotfixCreator.go . hotfixCreator`
- [x] Refactor cmd/hotfix.go to use interface
  - Create handler function `handleHotfix(creator hotfixCreator, args []string) error`
  - Update command RunE to use handler with manager as interface
  - Remove direct Manager and GitManager access
- [x] Write comprehensive unit tests for cmd/hotfix.go (cmd/hotfix_test.go)
  - Test hotfix creation with mocked production branch detection
  - Test JIRA integration scenarios with mocks
  - Test error handling and edge cases
  - Ensure fast execution (< 1 second)
- [x] Move integration tests to internal package if any exist
  - Check for tests that use real git repos in cmd/hotfix_test.go
  - Move to internal/hotfix_test.go if found
- [x] Run validation: `go build && go test ./cmd/...` to ensure fast unit tests

## Notes
[Implementation notes]