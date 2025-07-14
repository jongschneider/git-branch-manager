# Refactor Fishy Tests Report

This document identifies potentially problematic tests after the refactoring from `.envrc` to `gbm.branchconfig.yaml`. These tests may contain workarounds or behavior changes that mask underlying issues rather than properly testing the intended functionality.

## 🚨 High Priority Issues

### 1. TestListCommand_UntrackedWorktrees (cmd/list_test.go:167-168)

**Location**: `cmd/list_test.go:167-168`

**Issue**: Test assertion was changed from expecting "IN_SYNC" to "UNTRACKED" with a suspicious comment explaining the behavior change.

**Original vs New**:
```go
// OLD BEHAVIOR (likely):
assert.Contains(t, mainWorktree.SyncStatus, "IN_SYNC")

// NEW BEHAVIOR:
// When cloning without gbm.branchconfig.yaml, MAIN worktree starts as UNTRACKED until config is properly set up
assert.Contains(t, mainWorktree.SyncStatus, "UNTRACKED")
```

**Why it's fishy**: The comment suggests this is a behavioral change rather than a test fix. If the MAIN worktree should be tracked after cloning, this might indicate the new implementation has a bug where it doesn't properly track the main worktree, and the test was adjusted to match the broken behavior instead of fixing the underlying issue.

**Recommendation**: Investigate if the main worktree should actually be properly tracked after cloning with a default config.

### 2. TestValidateCommand_DuplicateWorktrees (cmd/validate_test.go:346-371)

**Location**: `cmd/validate_test.go:346-371`

**Issue**: Test completely changed behavior from testing `.envrc` duplicate environment variables to testing YAML duplicate keys.

**Original vs New**:
```go
// OLD BEHAVIOR (from git diff comments):
// Create .envrc with duplicate environment variables (last one wins)
// err := repo.WriteFile(".envrc", "MAIN=main\nMAIN=develop\nTEST=feature/auth")
// Should fail because 'develop' branch doesn't exist
// Verify that the last occurrence wins (MAIN=develop, not MAIN=main)

// NEW BEHAVIOR:
// Create gbm.branchconfig.yaml with duplicate worktree names (invalid YAML)
configContent := `worktrees:
  main:
    branch: main
  main:              # Duplicate key
    branch: develop
`
// Should fail due to duplicate YAML keys
assert.Contains(t, err.Error(), "mapping key \"main\" already defined")
```

**Why it's fishy**: This is testing completely different behavior. The original test was validating that duplicate environment variables behave correctly (last one wins), while the new test is just checking that YAML parsers reject duplicate keys. This masks the loss of functionality around handling duplicate entries and doesn't test the equivalent behavior in the new format.

**Recommendation**: Consider if there should be equivalent behavior for handling duplicate worktree definitions, or if this represents a valid breaking change that should be documented.

## ⚠️ Medium Priority Issues

### 3. Naming Convention Inconsistency

**Location**: Multiple test files

**Issue**: Inconsistent naming conventions between old and new test data.

**Details**:
- `NewStandardEnvrcRepo` uses uppercase keys: "MAIN", "DEV", "FEAT", "PROD"
- `NewStandardGBMConfigRepo` uses lowercase keys: "main", "dev", "feat", "prod"
- Some tests still reference uppercase worktree names ("MAIN") while using lowercase config keys ("main")

**Examples from git diff**:
```go
// Tests still look for "MAIN" worktree:
mainWorktree, found := findWorktreeInRows(rows, "MAIN")

// But config uses "main":
"main": "main",
"dev":  "develop",
```

**Why it's fishy**: This suggests the naming convention change may not have been fully thought through. It's unclear if worktree names should be case-sensitive, case-insensitive, or if there's a canonical format.

**Recommendation**: Clarify the intended naming convention and ensure consistency between config keys, worktree directory names, and test expectations.

### 4. Test Harness Ordering Changes

**Location**: `internal/testutils/git_harness.go`

**Issue**: The ordered keys for test reproducibility changed significantly.

**Original vs New**:
```go
// OLD (.envrc):
orderedKeys := []string{"MAIN", "PREVIEW", "PROD"}

// NEW (gbm.branchconfig.yaml):
orderedKeys := []string{"main", "preview", "staging", "dev", "feat", "prod", "hotfix"}
```

**Why it's fishy**: The test harness completely changed its expected ordering and added new standard keys. This suggests tests might be passing not because the functionality is correct, but because the test infrastructure was modified to match the new implementation.

**Recommendation**: Verify that the new ordering matches actual user expectations and isn't just convenient for making tests pass.

## 📝 Minor Issues

### 5. Error Message Changes

**Location**: Various test files

**Issue**: Many tests updated error message expectations from ".envrc" to "gbm.branchconfig.yaml" without verifying the new error messages are equally helpful.

**Example**:
```go
// OLD:
assert.Contains(t, err.Error(), ".envrc")

// NEW:
assert.Contains(t, err.Error(), "gbm.branchconfig.yaml")
```

**Why it's worth noting**: While these changes are expected, some error messages might be less clear than before, and the tests might be passing without anyone verifying the user experience is equivalent or better.

## Summary

The refactoring appears to have several concerning patterns:

1. **Behavioral workarounds**: Tests changed expectations to match new implementation rather than ensuring equivalent functionality
2. **Lost functionality**: Some features (like duplicate handling) seem to have been lost rather than translated
3. **Naming inconsistencies**: Unclear conventions around case sensitivity and naming
4. **Test infrastructure changes**: The test harness was significantly modified, which could mask implementation issues

## Recommendations

1. **Review TestListCommand_UntrackedWorktrees**: Investigate if main worktree tracking is actually broken
2. **Restore duplicate handling**: Decide if duplicate worktree handling should be implemented or document as breaking change
3. **Standardize naming**: Establish clear naming conventions for worktree names and configuration keys
4. **Verify error messages**: Ensure new error messages are as helpful as the old ones
5. **Add integration tests**: Consider adding end-to-end tests that verify the migration path from `.envrc` to `gbm.branchconfig.yaml` works correctly