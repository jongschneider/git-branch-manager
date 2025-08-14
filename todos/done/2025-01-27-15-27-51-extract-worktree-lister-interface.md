# Extract worktreeLister interface for cmd/list.go
**Status:** Done
**Agent PID:** 52752

## Original Todo
- [ ] **Extract worktreeLister interface for cmd/list.go**
  - Create interface with: GetSyncStatus(), GetAllWorktrees(), GetSortedWorktreeNames(), GetWorktreeMapping()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeLister.go . worktreeLister`
  - Refactor command functions to use interface
  - Write unit tests with mocks
  - Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

## Description
Extract a `worktreeLister` interface from `cmd/list.go` to enable unit testing with mocks and follow the established interface patterns in the codebase. The interface will abstract the four Manager methods needed for worktree listing: `GetSyncStatus()`, `GetAllWorktrees()`, `GetSortedWorktreeNames()`, and `GetWorktreeMapping()`. This follows the successful pattern established in `cmd/pull.go`, `cmd/push.go`, and `cmd/validate.go`.

## Implementation Plan
- [x] Create worktreeLister interface in cmd/list.go with methods: GetSyncStatus(), GetAllWorktrees(), GetSortedWorktreeNames(), GetWorktreeMapping()
- [x] Add go:generate directive for mock generation using github.com/matryer/moq
- [x] Extract handleList function that accepts worktreeLister interface and contains main listing logic
- [x] Refactor listCmd.RunE to use handleList with manager as interface
- [x] Generate mock file using go generate command
- [x] Convert cmd/list_test.go integration tests to unit tests using worktreeListerMock
- [x] Move integration tests to internal/list_test.go (following established pattern)
- [x] Run validation: go test ./cmd -v to verify unit tests pass quickly
- [x] User test: Run gbm list command to verify functionality works correctly

## Notes
[Implementation notes]