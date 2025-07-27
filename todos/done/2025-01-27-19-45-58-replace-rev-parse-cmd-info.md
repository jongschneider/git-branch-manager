# Replace duplicate rev-parse calls in cmd/info.go with new utilities
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Replace duplicate rev-parse calls in cmd/info.go with new utilities**
- Update line 338 to use VerifyRef
- Test: `gbm info` command still works correctly

## Description
This todo appears to be already completed. The cmd/info.go file is already using the VerifyRef utilities correctly at line 327, and there are no direct git rev-parse calls remaining in this file that need to be replaced.

## Implementation Plan
- [x] Verify the current state of cmd/info.go line 338 and surrounding context
- [x] Confirm that VerifyRef utilities are properly used throughout cmd/info.go
- [x] Check if there are any remaining direct git rev-parse calls in cmd/info.go that still need replacement
- [x] Update the todo list to reflect the actual current state
- [x] Automated test: Run `go test ./cmd` to ensure info command tests pass
- [x] User test: Run `gbm info` command to verify it works correctly

## Notes
[Implementation notes]