# Analysis: Extract worktreeValidator interface for cmd/validate.go

## Where changes are needed
- cmd/validate.go: Define `worktreeValidator` interface and refactor command execution to depend on it.
- cmd/validate.go: Add `//go:generate` directive to create `autogen_worktreeValidator.go` with moq.
- cmd/validate.go: Extract core logic into a small function (e.g., `handleValidate(v worktreeValidator) error`) mirroring pull/push patterns.
- cmd/validate_test.go: Replace real-repo integration tests with fast unit tests using the generated mock; move existing integration tests to internal.
- internal/tests: Add or adapt integration tests that validate `Manager` behavior (mapping + branch existence) without going through the cmd layer.

## Existing patterns to follow
- cmd/pull.go defines `worktreePuller` with `//go:generate` and isolates logic in `handlePull*` functions.
- cmd/push.go defines `worktreePusher` similarly; both have generated mocks: `autogen_worktreePuller.go`, `autogen_worktreePusher.go`.
- We will mirror this pattern for `validate`.

## Which files need modification
- cmd/validate.go (introduce interface, go:generate, refactor RunE to delegate to handler)
- cmd/validate_test.go (convert to unit tests with mocks, relocate integration tests)
- New: cmd/autogen_worktreeValidator.go (generated)
- New or moved: internal/validate_integration_test.go (or grouped within an existing internal integration suite)

## Related functionality/code
- internal/manager.go
  - `GetWorktreeMapping() (map[string]string, error)` — used at cmd/list.go:70 and validate.go currently
  - `BranchExists(branchName string) (bool, error)` — wrapper over `gitManager.BranchExists`
- internal/git.go
  - `GitManager.BranchExists(branchName string) (bool, error)`
- internal/table.go and internal/styles.go are used by the table rendering in validate flow

## Notes / Considerations
- Keep output behavior identical (messages, table columns, success/failure conditions) to avoid breaking tests/UX.
- Unit tests in cmd should verify:
  - When all branches exist: success header and VALID per row, no error.
  - When some branches missing: failure header and NOT FOUND for those, returns error.
  - Empty mapping: prints an empty table (or no rows) and succeeds.
  - Error propagation from `GetWorktreeMapping` and `BranchExists`.
- Integration tests moved to internal should validate real git behaviors (branch creation, missing branches, malformed YAML) without relying on CLI.
