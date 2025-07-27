# Extract git diff calls from cmd/info.go to use GetFileChanges
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Extract git diff calls from cmd/info.go to use GetFileChanges**
- Update lines 197, 238 to use new utility
- Test: `gbm info` shows correct file changes

## Description
Replace the getModifiedFiles function in cmd/info.go with the centralized GetFileChanges utility. This will eliminate 92 lines of duplicate git diff parsing logic while maintaining file change functionality and providing clearer separation between staged and unstaged changes.

## Implementation Plan
- [x] Update line 111 in cmd/info.go to use GetFileChanges instead of getModifiedFiles
- [x] Remove the getModifiedFiles function entirely (lines 162-254)
- [x] Ensure proper GitManager access and error handling
- [x] Test that the refactored code compiles without errors
- [x] Automated test: Run `go test ./cmd` to ensure info command tests pass
- [x] User test: Verify `gbm info` shows correct file changes

## Notes
[Implementation notes]