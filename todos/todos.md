## Interface Refactoring for Unit Testing Project

**Reference**: See `docs/interface-refactoring-plan.md` for complete analysis and strategy.

### Phase 1: Establish Patterns (Low Complexity Commands)

**✅ COMPLETED**: Extract worktreePuller interface for cmd/pull.go
- Interface created with: PullAllWorktrees(), PullWorktree(), IsInWorktree(), GetAllWorktrees()
- Mock generation added: `//go:generate go tool moq -out ./autogen_worktreePuller.go . worktreePuller`
- Command functions refactored to use interface
- Unit tests with mocks completed (cmd/pull_test.go)
- Integration tests added (internal/pull_test.go) - TestManager_PullWorktree fully working

**✅ COMPLETED**: Extract worktreePusher interface for cmd/push.go
- Interface created with: PushAllWorktrees(), PushWorktree(), IsInWorktree(), GetAllWorktrees()
- Mock generation added: `//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreePusher.go . worktreePusher`
- Command functions refactored to use interface
- Unit tests with mocks completed (cmd/push_test.go)
- Integration tests added (internal/push_test.go) - TestManager_PushWorktree and TestManager_PushAllWorktrees fully working




### Phase 2: Build on Patterns (Medium-Low Complexity)


- [ ] **Extract worktreeLister interface for cmd/list.go**
  - Create interface with: GetSyncStatus(), GetAllWorktrees(), GetSortedWorktreeNames(), GetWorktreeMapping()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeLister.go . worktreeLister`
  - Refactor command functions to use interface
  - Write unit tests with mocks
  - Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

### Phase 3: Handle Complex Data (Medium Complexity)

- [ ] **Extract worktreeInfoProvider interface for cmd/info.go**
  - Create interface with: GetWorktrees(), GetWorktreeStatus(), GetCommitHistory(), GetFileChanges(), GetCurrentBranchInPath(), GetUpstreamBranch(), GetAheadBehindCount(), VerifyRefInPath(), GetWorktreeBaseBranch(), GetConfig()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider`
  - Mock external JIRA CLI calls
  - Separate data aggregation logic for better testing
  - Write comprehensive unit tests
  - Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

- [ ] **Extract hotfixCreator interface for cmd/hotfix.go**
  - Create interface with: GetGBMConfig(), GetHotfixPrefix(), AddWorktree(), GenerateBranchFromJira(), GetJiraIssues(), FindProductionBranch()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_hotfixCreator.go . hotfixCreator`
  - Extract production branch detection logic to internal package
  - Interface JIRA integration
  - Write unit tests with mocks
  - Additionally, move any integration tests that use real git repos, worktrees, etc. out of the cmd layer and into the internal layer where they belong

### Phase 4: Architectural Refactoring (High Complexity)

- [ ] **Refactor cmd/clone.go - Move logic to internal package**
  - Move 80% of clone logic to internal.CloneManager
  - Create thin cmd layer that delegates to CloneManager
  - Create file system abstraction layer if needed for testing
  - Test CloneManager with integration tests in internal package

- [ ] **Refactor cmd/mergeback.go - Split into internal services**
  - Move git analysis logic to internal.MergebackAnalyzer
  - Move merge execution to internal.MergeExecutor  
  - Move user interaction to internal.UserInteractor
  - Create thin cmd layer that delegates to services
  - Test services with integration tests in internal package

### Common Infrastructure

- [ ] **Update helper functions for interface support**
  - Modify createInitializedManager() to return interfaces where appropriate
  - Add interface wrappers for backward compatibility
  - Update common completion functions to use interfaces

- [ ] **Standardize integration test patterns in internal package**
  - Use `internal/testutils/git_harness.go` GitTestRepo methods consistently
  - Replace manual `execGitCommandRun` calls with GitTestRepo methods: `repo.CreateSynchronizedBranch()`, `repo.WriteFile()`, `repo.CommitChanges()`, `repo.PushBranch()`
  - Create helper functions like `runGitInWorktree()` for worktree-specific operations
  - Follow `internal/git_add_test.go` pattern for test structure and setup
  - Use `testutils.NewGitTestRepo()` with proper options for consistent test environments

- [ ] **Move integration tests to internal package**
  - Move cmd/*_test.go integration tests to internal/*_test.go
  - Ensure actual functionality is tested in internal package
  - Keep only interface-based unit tests in cmd package
  - Pattern: cmd/pull_test.go has fast unit tests with mocks, internal/pull_test.go has integration tests with real git repos

- [ ] **Improve integration test infrastructure**
  - Add `createRemoteChanges()` helper pattern to simulate remote developer changes
  - Standardize branch creation patterns to avoid git worktree tracking conflicts
  - Use proper base branches when creating worktrees to ensure git operations work correctly
  - Add debugging helpers for troubleshooting test failures (git worktree list, branch status, etc.)
  
- [ ] **Fix incomplete integration tests**
  - Complete `TestManager_PullAllWorktrees()` in internal/pull_test.go (currently has issues with GetAllWorktrees returning 0 results)
  - Investigate why Manager.GetAllWorktrees() doesn't find worktrees that exist in git worktree list
  - Ensure integration tests cover both single and multiple worktree scenarios

- [ ] **Document common testing patterns and gotchas**
  - Document git worktree branch checkout conflicts (branch can only be checked out in one place)
  - Create guide for when to use separate clones vs main repo for testing remote changes
  - Add examples for proper base branch selection when creating test worktrees
  - Document GitTestRepo method usage patterns and when to use each method

- [ ] **Add interface generation validation**
  - Ensure `//go:generate` directives are run and mock files are up-to-date in CI
  - Add build validation that checks if mock files are in sync with interfaces
  - Consider automating mock regeneration in pre-commit hooks or CI pipeline

### Validation

- [ ] **Verify all cmd package tests use only mocks**
  - Ensure no real git operations in cmd tests
  - Verify fast test execution (< 1 second for all cmd tests)
  - Confirm test isolation (no shared state between tests)

- [ ] **Validate integration test coverage in internal package**
  - Ensure all core functionality is covered by integration tests
  - Verify real git operations work correctly
  - Test edge cases and error conditions

## Notes

- Work through each todo sequentially to minimize risk
- Run tests after each change to catch regressions early
- Reference the gameplan document for implementation details and guidelines
- Each utility function should use `enhanceGitError()` for consistent error handling
- For high-complexity commands (mergeback, clone), prefer architectural refactoring over extensive mocking
- Start with Phase 1 to establish patterns before tackling complex commands
