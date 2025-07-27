# Add GetFileChanges utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Add GetFileChanges utility function to GitManager**
- Location: `internal/git.go`
- Replace pattern: `git diff --numstat`, `git diff --cached --numstat`
- Files affected: `cmd/info.go:197,238`

## Description
Create a new GetFileChanges utility function in GitManager that can replace the duplicate git diff --numstat patterns currently used in cmd/info.go. This will centralize file change retrieval with flexible options for both staged and unstaged changes, eliminating the duplicate parsing logic.

## Implementation Plan
- [x] Design FileChangeOptions struct for flexible query options (internal/git.go)
- [x] Implement GetFileChanges utility function in GitManager (internal/git.go)
- [x] Add helper function parseNumstatOutput for parsing git diff --numstat output
- [x] Write unit tests for the new GetFileChanges function
- [x] Automated test: Run `go test ./internal` to ensure new utility works correctly
- [x] User test: Verify existing functionality still works (gbm info command shows file changes)

## Notes
[Implementation notes]