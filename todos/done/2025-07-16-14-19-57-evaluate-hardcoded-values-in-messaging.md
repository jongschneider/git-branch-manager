# Evaluate hardcoded values in messaging
**Status:** Done
**Agent PID:** 50911

## Original Todo
- evaluate hardcoded values in messaging
    * `gbm info`
    * all messages printing out info to the user

## Description
The codebase contains numerous hardcoded values in user-facing messaging that should be dynamic or configurable. Key issues include outdated `.envrc` references in help text, hardcoded section headers in `gbm info`, hardcoded branch names in candidate detection, and hardcoded status/priority messages throughout the system.

## Implementation Plan
- [x] Fix outdated `.envrc` references in root command help text to use `gbm.branchconfig.yaml`
- [x] Replace hardcoded candidate branches in `cmd/info.go` with configurable list or dynamic detection
- [x] Update hardcoded section headers in `internal/info_renderer.go` to use configurable values
- [x] Fix hardcoded file paths in completion help text to be more generic
- [x] Create constants for commonly referenced configuration filenames to improve consistency
- [x] Test that all messaging changes work correctly and maintain functionality

## Notes
[Implementation notes]