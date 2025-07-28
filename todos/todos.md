## Interface Refactoring for Unit Testing Project

**Reference**: See `docs/interface-refactoring-plan.md` for complete analysis and strategy.

### Phase 1: Establish Patterns (Low Complexity Commands)


- [ ] **Extract worktreePuller interface for cmd/pull.go**
  - Create interface with: PullAllWorktrees(), PullWorktree(), IsInWorktree(), GetAllWorktrees()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreePuller.go . worktreePuller`
  - Refactor command functions to use interface
  - Write unit tests with mocks

- [ ] **Extract worktreePusher interface for cmd/push.go**
  - Create interface with: PushAllWorktrees(), PushWorktree(), IsInWorktree(), GetAllWorktrees()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreePusher.go . worktreePusher`
  - Refactor command functions to use interface
  - Write unit tests with mocks

- [ ] **Extract worktreeRemover interface for cmd/remove.go**
  - Create interface with: GetWorktreePath(), GetWorktreeStatus(), RemoveWorktree(), GetAllWorktrees()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeRemover.go . worktreeRemover`
  - Refactor command functions to use interface
  - Write unit tests with mocks

- [ ] **Extract worktreeSwitcher interface for cmd/switch.go**
  - Create interface with: GetWorktreePath(), SetCurrentWorktree(), GetPreviousWorktree(), GetAllWorktrees(), GetSortedWorktreeNames(), GetStatusIcon()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeSwitcher.go . worktreeSwitcher`
  - Refactor command functions to use interface
  - Write unit tests with mocks

### Phase 2: Build on Patterns (Medium-Low Complexity)

- [ ] **Extract worktreeSyncer interface for cmd/sync.go**
  - Create interface with: GetSyncStatus(), SyncWithConfirmation()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeSyncer.go . worktreeSyncer`
  - Refactor command functions to use interface
  - Write unit tests with mocks

- [ ] **Extract worktreeValidator interface for cmd/validate.go**
  - Create interface with: GetWorktreeMapping(), BranchExists()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeValidator.go . worktreeValidator`
  - Refactor command functions to use interface
  - Write unit tests with mocks

- [ ] **Extract worktreeLister interface for cmd/list.go**
  - Create interface with: GetSyncStatus(), GetAllWorktrees(), GetSortedWorktreeNames(), GetWorktreeMapping()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeLister.go . worktreeLister`
  - Refactor command functions to use interface
  - Write unit tests with mocks

### Phase 3: Handle Complex Data (Medium Complexity)

- [ ] **Extract worktreeInfoProvider interface for cmd/info.go**
  - Create interface with: GetWorktrees(), GetWorktreeStatus(), GetCommitHistory(), GetFileChanges(), GetCurrentBranchInPath(), GetUpstreamBranch(), GetAheadBehindCount(), VerifyRefInPath(), GetWorktreeBaseBranch(), GetConfig()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider`
  - Mock external JIRA CLI calls
  - Separate data aggregation logic for better testing
  - Write comprehensive unit tests

- [ ] **Extract hotfixCreator interface for cmd/hotfix.go**
  - Create interface with: GetGBMConfig(), GetHotfixPrefix(), AddWorktree(), GenerateBranchFromJira(), GetJiraIssues(), FindProductionBranch()
  - Add mock generation: `//go:generate go tool moq -out ./autogen_hotfixCreator.go . hotfixCreator`
  - Extract production branch detection logic to internal package
  - Interface JIRA integration
  - Write unit tests with mocks

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

- [ ] **Move integration tests to internal package**
  - Move cmd/*_test.go integration tests to internal/*_test.go
  - Ensure actual functionality is tested in internal package
  - Keep only interface-based unit tests in cmd package

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
