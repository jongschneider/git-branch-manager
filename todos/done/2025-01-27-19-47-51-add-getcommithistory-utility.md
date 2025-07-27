# Add GetCommitHistory utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Add GetCommitHistory utility function to GitManager**
- Location: `internal/git.go`
- Replace pattern: `git log --oneline --format=...`
- Files affected: `cmd/info.go:161`

## Description
Create a new GetCommitHistory utility function in GitManager that can replace the 4 different git log patterns currently used throughout the codebase. This will centralize commit history retrieval with flexible options for different use cases including recent commits, mergeback analysis, merge detection, and hotfix tracking.

## Implementation Plan
- [x] Design CommitInfo struct to represent commit data (internal/git.go)
- [x] Design CommitHistoryOptions struct for flexible query options (internal/git.go)
- [x] Implement GetCommitHistory utility function in GitManager (internal/git.go)
- [x] Add helper functions for parsing commit output and building git log commands
- [x] Write unit tests for the new GetCommitHistory function
- [x] Automated test: Run `go test ./internal` to ensure new utility works correctly
- [x] User test: Verify existing functionality still works (gbm info command shows commit history)

## Notes
[Implementation notes]