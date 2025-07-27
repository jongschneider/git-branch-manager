# Add GetAheadBehindCount utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Add GetAheadBehindCount utility function to GitManager**
- Location: `internal/git.go` 
- Replace pattern: `git rev-list --left-right --count HEAD...@{upstream}`
- Files affected: `cmd/info.go:309`, `internal/git.go:521`

## Description
Add a centralized `GetAheadBehindCount` utility method to GitManager that replaces 2 duplicate instances of `git rev-list --left-right --count HEAD...@{upstream}` command execution. The method will return (ahead, behind, error) and provide consistent error handling while following established GitManager patterns, replacing silent error handling with proper error propagation.

## Implementation Plan
- [x] Add GetAheadBehindCount method to GitManager in internal/git.go (after GetUpstreamBranch method)
- [x] Replace direct git command in cmd/info.go:304-314 with new utility method, handling errors appropriately
- [x] Replace direct git command in internal/git.go:521-529 (GetWorktreeStatus method) with new utility method
- [x] Add unit tests for GetAheadBehindCount method covering upstream exists, no upstream, and error cases
- [x] Run validation commands: `just lint`, `just format`, `just test-changed`

## Notes
[Implementation notes]