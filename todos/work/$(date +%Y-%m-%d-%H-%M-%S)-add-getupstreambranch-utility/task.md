# Add GetUpstreamBranch utility function to GitManager
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Add GetUpstreamBranch utility function to GitManager** 
- Location: `internal/git.go`
- Replace pattern: `git rev-parse --abbrev-ref @{upstream}`
- Files affected: `cmd/info.go:300`, `internal/git.go:643,676`

## Description
Add a centralized `GetUpstreamBranch` utility method to GitManager that replaces 4 duplicate instances of `git rev-parse --abbrev-ref @{upstream}` command execution across the codebase. The method will provide consistent error handling and follow established GitManager patterns while supporting both informational usage (info command) and conditional logic usage (push/pull operations).

## Implementation Plan
- [x] Add GetUpstreamBranch method to GitManager in internal/git.go (after line ~400 with other branch methods)
- [x] Replace direct git command in cmd/info.go:298-304 with new utility method
- [x] Replace direct git command in internal/git.go:653 (PushWorktree method) with new utility method  
- [x] Replace direct git command in internal/git.go:686 (PullWorktree method) with new utility method
- [x] Add unit tests for GetUpstreamBranch method covering upstream exists, no upstream, and error cases
- [x] Run validation commands: `just lint`, `just format`, `just test-changed`

## Notes
[Implementation notes]