# Push Integration Testing Implementation

## Overview

The worktreePusher interface extraction is complete with unit tests for the cmd layer handlers, but we need comprehensive integration tests for the actual Manager methods that interact with real git repositories.

## Current State

### ✅ Completed
- **Interface Definition**: `worktreePusher` interface in `cmd/push.go`
- **Mock Generation**: `cmd/autogen_worktreePusher.go` 
- **Unit Tests**: `cmd/push_test.go` with 11 test cases using mocks
- **Handler Refactoring**: All push handlers use interface

### ❌ Missing
- **Integration Tests**: `internal/push_test.go` needs to be implemented
- **Real Git Testing**: Tests against actual repositories using git_harness.go

## Required Implementation: internal/push_test.go

### Test Structure
The integration tests should use the established patterns from other internal tests:
- Use `internal/testutils/git_harness.go` for repo setup
- Use `testutils.NewStandardGBMConfigRepo()` for consistent test scenarios
- Test real git operations, not mocked interfaces

### Methods to Test

#### 1. Manager.PushAllWorktrees() 
**Happy Paths:**
- Push all worktrees when all have local changes
- Push all worktrees when some are up to date
- Push all worktrees when all are up to date (no-op)
- Push with upstream tracking already set
- Push with new upstream tracking needed

**Sad Paths:**
- Push with one worktree in detached HEAD state
- Push with network/remote errors
- Push with merge conflicts
- Push with no worktrees configured

#### 2. Manager.PushWorktree(worktreeName string)
**Happy Paths:**
- Push worktree with local changes
- Push worktree that's already up to date
- Push multiple commits
- Push with upstream tracking
- Push new branch to remote

**Sad Paths:**
- Push nonexistent worktree
- Push worktree in detached HEAD
- Push with network/remote errors
- Push with merge conflicts
- Push worktree with no remote configured

#### 3. Manager.IsInWorktree(currentPath string)
**Happy Paths:**
- Detect when in a valid worktree (return true, worktree_name)
- Detect when not in a worktree (return false, "", nil)
- Detect from subdirectory within worktree
- Detect from root of worktree

**Sad Paths:**
- Path that doesn't exist
- Permission denied on path
- Malformed git repository
- Corrupted worktree

#### 4. Manager.GetAllWorktrees()
**Happy Paths:**
- Get worktrees with standard gbm config
- Get worktrees with complex branching setup
- Get worktrees with various branch states (ahead/behind/up-to-date)
- Get worktrees with untracked/staged/committed changes

**Sad Paths:**
- No worktrees configured
- Worktree directories that don't exist
- Corrupted worktree references
- Invalid gbm.branchconfig.yaml

### Test Utilities Needed

#### Setup Functions
```go
// setupPushTestScenario creates a realistic git repository with:
// - Remote repository (bare)
// - Local bare clone setup (like gbm uses)
// - Multiple worktrees with gbm.branchconfig.yaml
// - Various commit states for testing
func setupPushTestScenario(t *testing.T) (*testutils.GitTestRepo, *Manager, string)

// makeLocalChanges creates commits in a worktree for testing push operations
func makeLocalChanges(t *testing.T, worktreePath, filename, content string)

// verifyRemoteHasCommit verifies that a commit was successfully pushed to remote
func verifyRemoteHasCommit(t *testing.T, repo *testutils.GitTestRepo, branch, commitMessage string)

// createWorktreeWithChanges sets up a worktree and adds commits for testing
func createWorktreeWithChanges(t *testing.T, manager *Manager, worktreeName, branchName string, numCommits int)
```

#### Verification Functions
```go
// verifyPushResults checks that push operations completed successfully
func verifyPushResults(t *testing.T, repo *testutils.GitTestRepo, expectedPushes map[string][]string)

// verifyWorktreeState checks worktree is in expected state after operations
func verifyWorktreeState(t *testing.T, worktreePath, expectedBranch string, expectedFiles []string)

// verifyRemoteSync checks that local and remote branches are in sync
func verifyRemoteSync(t *testing.T, repo *testutils.GitTestRepo, branch string)
```

### Test Scenarios Based on Real Usage

#### Scenario 1: Developer Workflow
```go
// Developer adds commits to multiple worktrees and pushes all
// Should handle: new commits, upstream tracking, success verification
TestManager_PushAllWorktrees_DeveloperWorkflow
```

#### Scenario 2: CI/CD Integration  
```go
// Automated system pushes specific worktrees
// Should handle: single worktree push, error handling, idempotency
TestManager_PushWorktree_CICDIntegration
```

#### Scenario 3: Complex Repository State
```go
// Repository with mixed worktree states
// Should handle: some ahead, some behind, some up-to-date
TestManager_PushAllWorktrees_MixedStates
```

#### Scenario 4: Error Recovery
```go
// Handle various error conditions gracefully
// Should handle: network issues, conflicts, permission errors
TestManager_Push_ErrorRecovery
```

### Integration with Existing Test Infrastructure

#### Use Existing Patterns
- Follow `internal/pull_test.go` patterns for similar git operations
- Use `testutils.GitTestRepo` for consistent repo setup
- Use `require.NoError` / `assert.Error` for proper test failures
- Use table-driven tests for multiple scenarios

#### Test Data Management
- Use temporary directories for all test repositories
- Clean up test data automatically with `t.Cleanup()`
- Use realistic branch names and commit messages
- Create test scenarios that mirror real gbm usage

### Performance Considerations
- Tests should complete in reasonable time (< 5 seconds each)
- Use local git operations (no network calls)
- Parallel test execution where possible
- Efficient test data setup and teardown

### Error Testing Strategy
- Test both expected errors (nonexistent worktree) and unexpected errors (network)
- Verify error messages are helpful and actionable  
- Test error handling doesn't leave repository in bad state
- Test recovery from partial failures

## Implementation Notes

### Git Operations to Test
- `git push` with various flags and scenarios
- `git rev-parse` for worktree detection
- `git status` for change detection
- `git branch` for branch management
- `git remote` for upstream configuration

### Repository States to Cover
- Clean worktrees (no changes)
- Dirty worktrees (uncommitted changes)
- Staged changes
- Untracked files
- Ahead of remote
- Behind remote
- Diverged from remote
- New branches (no upstream)

### Edge Cases
- Empty repositories
- Single commit repositories
- Repositories with no remotes
- Repositories with multiple remotes
- Repositories with complex merge histories

## Success Criteria

### Test Coverage
- All four interface methods fully tested
- All happy paths covered
- All error conditions covered
- Real git operations verified
- No flaky or unreliable tests

### Test Quality
- Tests are fast and reliable
- Clear test names and documentation
- Proper setup and teardown
- Realistic test scenarios
- Good error messages on failure

### Integration
- Tests pass consistently in CI
- Tests work on different platforms
- Tests don't interfere with each other
- Tests clean up properly after execution

## Next Steps

1. **Create test file**: `internal/push_test.go`
2. **Implement setup functions**: Based on existing patterns in `internal/pull_test.go`
3. **Add test cases**: Start with happy paths, then add error cases
4. **Verify git operations**: Ensure pushes actually work with real repositories
5. **Test error handling**: Verify graceful handling of all error conditions
6. **Performance optimization**: Ensure tests run quickly and reliably

## Context for Implementation

### Existing Test Patterns
- Reference `internal/pull_test.go` for git operation testing patterns
- Reference `internal/manager_add_integration_test.go` for Manager testing
- Reference `internal/testutils/` for git harness usage

### Key Requirements
- Must test real git operations, not mocks
- Must be reliable and not flaky
- Must cover all interface methods comprehensively
- Must follow established testing patterns in the codebase