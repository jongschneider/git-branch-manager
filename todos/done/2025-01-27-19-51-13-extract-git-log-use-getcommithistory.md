# Extract git log calls from cmd/info.go to use GetCommitHistory
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Extract git log calls from cmd/info.go to use GetCommitHistory**
- Update line 161 to use new utility
- Test: `gbm info` shows correct commit history

## Description
Replace the getRecentCommits function in cmd/info.go with the centralized GetCommitHistory utility. This will eliminate 34 lines of duplicate git log parsing logic while maintaining identical functionality and leveraging enhanced error handling.

## Implementation Plan
- [x] Update line 103 in cmd/info.go to use GetCommitHistory instead of getRecentCommits
- [x] Remove the getRecentCommits function entirely (lines 160-193)
- [x] Ensure proper error handling and GitManager access
- [x] Test that the refactored code compiles without errors
- [x] Automated test: Run `go test ./cmd` to ensure info command tests pass
- [x] User test: Verify `gbm info` shows correct commit history

## Notes
[Implementation notes]