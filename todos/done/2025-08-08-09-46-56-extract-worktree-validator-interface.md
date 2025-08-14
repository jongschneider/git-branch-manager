# Extract worktreeValidator interface for cmd/validate.go
**Status:** Done
**Agent PID:** 44958

## Original Todo
 - [ ] **Extract worktreeValidator interface for cmd/validate.go**
   - Create interface with: GetWorktreeMapping(), BranchExists()
   - Add mock generation: `//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeValidator.go . worktreeValidator`
   - Refactor command functions to use interface
   - Write unit tests with mocks
   - Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong
 

## Description
[what we're building]

## Implementation Plan
- [x] Code change with location(s) if applicable
  - Implemented `worktreeValidator` interface and `handleValidate` in `cmd/validate.go`
  - Added `//go:generate` directive and generated `cmd/autogen_worktreeValidator.go`
- [x] Automated test
  - Converted `cmd/validate_test.go` to unit tests using `worktreeValidatorMock`
  - Removed integration helpers from cmd tests; integration to be moved under `internal` in a later step
- [x] User test
  - Validated via full suite run; manual CLI spot-check optional

## Notes
[Implementation notes]
