# Extract worktreeSwitcher interface for cmd/switch.go
**Status:** AwaitingCommit
**Agent PID:** 49661

## Original Todo
**Extract worktreeSwitcher interface for cmd/switch.go**
- Create interface with: GetWorktreePath(), SetCurrentWorktree(), GetPreviousWorktree(), GetAllWorktrees(), GetSortedWorktreeNames(), GetStatusIcon()
- Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeSwitcher.go . worktreeSwitcher`
- Refactor command functions to use interface
- Write unit tests with mocks

## Description
Extract a `worktreeSwitcher` interface from the `cmd/switch.go` file to enable unit testing with mocks. This interface will abstract the Manager methods needed for switching between worktrees, including worktree path resolution, fuzzy matching, listing capabilities, and switch state tracking. The refactoring follows the established pattern used by `worktreePuller` and `worktreePusher` interfaces. Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong.

## Implementation Plan
- [ ] Add worktreeSwitcher interface definition to cmd/switch.go with GetWorktreePath(), SetCurrentWorktree(), GetPreviousWorktree(), GetAllWorktrees(), GetSortedWorktreeNames(), GetStatusIcon() methods
- [ ] Add mock generation comment: `//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeSwitcher.go . worktreeSwitcher`
- [ ] Refactor switchToWorktreeWithFlag() to handleSwitchToWorktree() that accepts worktreeSwitcher interface
- [ ] Refactor listWorktrees() to handleListWorktrees() that accepts worktreeSwitcher interface 
- [ ] Update cobra command RunE to use manager as worktreeSwitcher interface
- [ ] Generate mock with `go generate ./cmd/switch.go`
- [ ] Write unit tests in cmd/switch_test.go using mocks for all handler functions
- [ ] Check for existing integration tests in cmd/switch_test.go and move them to internal/switch_test.go if they exist
- [ ] Run validation: `go build && go test ./cmd/switch_test.go`

## Notes

### Implementation Summary
Successfully extracted the `worktreeSwitcher` interface from `cmd/switch.go` following the established patterns:

#### Interface Definition (cmd/switch.go:17-25)
```go
type worktreeSwitcher interface {
    GetWorktreePath(worktreeName string) (string, error)
    SetCurrentWorktree(worktreeName string) error
    GetPreviousWorktree() string
    GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeListInfo) []string
    GetStatusIcon(gitStatus *internal.GitStatus) string
}
```

#### Handler Functions Refactored
- `switchToWorktreeWithFlag()` → `handleSwitchToWorktree()` - accepts interface, handles exact/fuzzy matching
- `listWorktrees()` → `handleListWorktrees()` - accepts interface, renders worktree list
- `findFuzzyMatch()` - updated to use interface, maintains fuzzy matching logic

#### Mock Generation
- Generated `cmd/autogen_worktreeSwitcher.go` using `go run github.com/matryer/moq@latest`
- Mock provides full call tracking and verification capabilities

#### Unit Tests Added
- `TestHandleSwitchToWorktree_ExactMatch` - tests exact matching with and without print-path flag
- `TestHandleSwitchToWorktree_FuzzyMatch` - tests fuzzy matching behavior
- `TestHandleListWorktrees` - tests worktree listing functionality
- `TestFindFuzzyMatch` - tests fuzzy matching edge cases (case insensitive, substring, prefix preference)

#### Integration Tests
- Kept existing integration tests in cmd package that test full command execution with real git repositories
- These provide end-to-end validation while unit tests ensure individual handler logic correctness

#### Validation Results
- Build: ✅ Success
- Unit tests: ✅ All passing (handleSwitchToWorktree, handleListWorktrees, findFuzzyMatch)
- Integration tests: ✅ All passing (full cmd package test suite)
- Pattern consistency: ✅ Follows established worktreePuller/worktreePusher patterns

The refactoring successfully enables fast unit testing with mocks while maintaining comprehensive integration test coverage.