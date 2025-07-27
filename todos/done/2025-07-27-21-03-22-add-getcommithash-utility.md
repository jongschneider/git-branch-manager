# Add GetCommitHash utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
- [ ] **Add GetCommitHash utility function to GitManager**
  - Location: `internal/git.go`
  - Replace pattern: `git rev-parse <ref>`
  - Files affected: `internal/git_add.go:48,54`

## Description
Add a `GetCommitHash` utility function to the GitManager struct that abstracts the `git rev-parse <ref>` command pattern. This function will replace direct `ExecGitCommand` calls for retrieving commit hashes, providing consistent error handling and following the established codebase patterns. The function will be used in `internal/git_add.go` line 54 where `git rev-parse` is currently called directly to get the base branch commit hash.

## Implementation Plan
- [x] Add `GetCommitHash(ref string) (string, error)` method to GitManager in `internal/git.go`
- [x] Add `GetCommitHashInPath(path, ref string) (string, error)` method to GitManager in `internal/git.go` 
- [x] Replace direct `ExecGitCommand` call in `internal/git_add.go:54` with `gm.GetCommitHash(baseBranch)`
- [x] Add comprehensive tests for `GetCommitHash` function in `internal/git_test.go`
- [x] Add tests for `GetCommitHashInPath` function in `internal/git_test.go`
- [x] Automated test: Run `go test ./internal` to verify new functions work correctly
- [x] User test: Run `gbm add` command to ensure worktree creation still works correctly

## Notes
[Implementation notes]