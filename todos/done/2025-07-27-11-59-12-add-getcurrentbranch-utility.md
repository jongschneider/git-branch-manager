# Add GetCurrentBranch utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
Add GetCurrentBranch utility function to GitManager
- Location: `internal/git.go`
- Replace pattern: `git rev-parse --abbrev-ref HEAD`
- Files affected: `cmd/info.go:291`, test files

## Description
Add a new `GetCurrentBranchInPath(path string)` utility method to GitManager that executes `git rev-parse --abbrev-ref HEAD` in a specified directory. This will replace the direct `exec.Command` usage in `cmd/info.go:291` and provide a consistent, reusable way to get the current branch from any git repository path.

## Implementation Plan
- [x] Add `GetCurrentBranchInPath(path string) (string, error)` method to GitManager in `internal/git.go`
- [x] Replace direct `exec.Command` usage in `cmd/info.go:291` with new utility method
- [x] Add unit tests for the new utility method in `internal/git_test.go`
- [x] User test: Verify `gbm info` command still works correctly and shows proper branch information

## Notes
Implementation notes