# Add VerifyRef utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
Add VerifyRef utility function to GitManager
- Location: `internal/git.go`
- Replace pattern: `git rev-parse --verify <ref>`
- Files affected: `cmd/info.go:338`, `internal/git_add.go:89`

## Description
Create a `VerifyRef` utility function in GitManager to replace scattered `git rev-parse --verify` calls throughout the codebase. This will centralize reference verification logic and provide consistent error handling. The function should follow existing GitManager patterns and handle both repository-level and worktree-specific verification needs.

## Implementation Plan
Based on the analysis, here's how we'll implement the VerifyRef utility:

- [x] Add `VerifyRef` method to GitManager in `internal/git.go` that takes a ref string and returns `(bool, error)` following existing patterns
- [x] Add `VerifyRefInPath` method to GitManager for worktree-specific verification that takes path and ref parameters  
- [x] Update `cmd/info.go:327` to use new VerifyRefInPath method instead of direct exec.Command
- [x] Update `internal/git_add.go:89` to use new VerifyRef method instead of direct ExecGitCommand
- [x] Automated test: Add unit tests for both VerifyRef methods to ensure they handle valid refs, invalid refs, and git errors correctly
- [x] User test: Run `gbm info` and `gbm add` commands to verify functionality still works correctly

## Notes
[Implementation notes]