# Refactoring Todo Checklist

## Phase 1: Easy Wins (Low Risk, High Impact)

### 1. Extract Manager Creation Helper
**Difficulty:** Easy | **Impact:** High | **Files:** 10+ cmd files | **Est. Time:** 30 minutes

#### Tasks:
- [ ] Create `createInitializedManager()` function in `cmd/` package
- [ ] Replace identical 4-line sequence in `cmd/add.go`
- [ ] Replace identical 4-line sequence in `cmd/clone.go`
- [ ] Replace identical 4-line sequence in `cmd/info.go`
- [ ] Replace identical 4-line sequence in `cmd/list.go`
- [ ] Replace identical 4-line sequence in `cmd/pull.go`
- [ ] Replace identical 4-line sequence in `cmd/push.go`
- [ ] Replace identical 4-line sequence in `cmd/remove.go`
- [ ] Replace identical 4-line sequence in `cmd/switch.go`
- [ ] Replace identical 4-line sequence in `cmd/sync.go`
- [ ] Replace identical 4-line sequence in `cmd/validate.go`

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] No functional changes to application behavior
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

### 2. Consolidate Terminal Width Detection
**Difficulty:** Easy | **Impact:** Medium | **Files:** 2 files | **Est. Time:** 15 minutes

#### Tasks:
- [ ] Extract `getTerminalWidth()` to shared utility file
- [ ] Update `internal/table.go` to use shared function
- [ ] Update `internal/info_renderer.go` to use shared function
- [ ] Remove duplicate terminal width detection code

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Terminal width detection works identically to before
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

### 3. Unify Time Formatting
**Difficulty:** Easy | **Impact:** Medium | **Files:** 2 files | **Est. Time:** 20 minutes

#### Tasks:
- [ ] Extract unified time formatting function to shared utility
- [ ] Merge `formatDuration()` and `formatRelativeTime()` logic
- [ ] Update `internal/info_renderer.go` to use shared function
- [ ] Update `internal/mergeback.go` to use shared function
- [ ] Remove duplicate time formatting implementations

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Time formatting output remains identical
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

## Phase 2: Medium Changes (Moderate Risk, Good Impact)

### 4. Consolidate Git Command Execution
**Difficulty:** Medium | **Impact:** High | **Files:** 3 locations | **Est. Time:** 45 minutes

#### Tasks:
- [ ] Create unified `execGitCommand()` function in `internal/git.go`
- [ ] Update `internal/mergeback.go` to use shared git execution
- [ ] Update `internal/testutils/git_harness.go` to use shared git execution
- [ ] Remove duplicate git command execution implementations
- [ ] Ensure consistent error handling across all git operations

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Git command execution behavior remains identical
- [ ] Error handling consistency maintained
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

### 5. Generic Push/Pull Command Handlers
**Difficulty:** Medium | **Impact:** Medium | **Files:** 2 cmd files | **Est. Time:** 40 minutes

#### Tasks:
- [ ] Create `handleCommandAll()` generic function
- [ ] Create `handleCommandCurrent()` generic function
- [ ] Create `handleCommandNamed()` generic function
- [ ] Replace `handlePushAll()` with generic version
- [ ] Replace `handlePullAll()` with generic version
- [ ] Replace `handlePushCurrent()` with generic version
- [ ] Replace `handlePullCurrent()` with generic version
- [ ] Replace `handlePushNamed()` with generic version
- [ ] Replace `handlePullNamed()` with generic version

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Push/pull command behavior remains identical
- [ ] Error messages remain consistent
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

### 6. Consolidate .envrc Parsing
**Difficulty:** Medium | **Impact:** Medium | **Files:** 2 files | **Est. Time:** 35 minutes

#### Tasks:
- [ ] Merge `ParseEnvrc()` and `parseEnvrcFile()` into single implementation
- [ ] Update `internal/mergeback.go` to use shared parser
- [ ] Remove duplicate .envrc parsing code
- [ ] Ensure consistent regex patterns and error handling

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] .envrc parsing behavior remains identical
- [ ] Both return types supported (map vs array)
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

## Phase 3: Structural Changes (Higher Risk, Good Impact)

### 7. Simplify Manager Delegation Layer
**Difficulty:** Hard | **Impact:** High | **Files:** Multiple internal files | **Est. Time:** 90 minutes

#### Tasks:
- [ ] Analyze current delegation patterns in Manager
- [ ] Choose approach: Remove delegation OR embed GitManager
- [ ] Update all Manager delegation methods:
  - [ ] `BranchExists()`
  - [ ] `GetStatusIcon()`
  - [ ] `GetRemoteBranches()`
  - [ ] `GetCurrentBranch()`
  - [ ] `PushWorktree()`
  - [ ] `PullWorktree()`
  - [ ] `IsInWorktree()`
  - [ ] `GetWorktreeStatus()`
- [ ] Update all calling code to use new pattern
- [ ] Remove redundant delegation methods

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] All manager functionality works identically
- [ ] No breaking changes to public API
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

## Phase 4: Test Infrastructure (Medium Risk, Quality Impact)

### 8. Consolidate Test Repository Creation
**Difficulty:** Medium | **Impact:** Medium | **Files:** testutils package | **Est. Time:** 75 minutes

#### Tasks:
- [ ] Analyze current test scenario functions
- [ ] Design builder pattern for test repository configuration
- [ ] Merge similar scenario functions:
  - [ ] `NewBasicRepo()`
  - [ ] `NewMultiBranchRepo()`
  - [ ] `NewEnvrcRepo()`
  - [ ] `NewStandardEnvrcRepo()`
  - [ ] `NewRepoWithConflictingBranches()`
  - [ ] `NewLargeHistoryRepo()`
  - [ ] `NewEmptyRepo()`
- [ ] Update all test files to use new builder pattern
- [ ] Remove duplicate scenario creation code

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Test repository creation behavior remains identical
- [ ] No test functionality lost
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

### 9. Unify Test Git Command Execution
**Difficulty:** Medium | **Impact:** Medium | **Files:** testutils package | **Est. Time:** 30 minutes

#### Tasks:
- [ ] Merge `runGitCommand()` and `runCommand()` in test utilities
- [ ] Use shared git execution from Phase 2 task #4
- [ ] Update all test files using old command execution
- [ ] Remove duplicate test command execution code

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Test git command execution behavior remains identical
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

## Maybe Someday Items

### 10. Unify Git Status Formatting (Moved from Phase 3)
**Difficulty:** Medium-Hard | **Impact:** Medium | **Files:** 2 files | **Est. Time:** 60 minutes

#### Tasks:
- [ ] Extract common git status logic while preserving different presentation
- [ ] Create `GitStatusFormatter` interface with multiple implementations
- [ ] Update `internal/styles.go` to use new interface
- [ ] Update `internal/info_renderer.go` to use new interface
- [ ] Remove duplicate git status formatting code

#### Completion Requirements:
- [ ] All existing tests pass without modification
- [ ] Code compiles without errors
- [ ] Git status formatting output remains identical in both contexts
- [ ] **COMPLETE ONLY WHEN ALL TESTS PASS**

---

## Important Notes

1. **TEST INTEGRITY**: Never modify existing tests during refactoring unless the test itself is being refactored
2. **COMPLETION RULE**: Tasks are only complete when ALL tests pass
3. **NO CHEATING**: Do not change tests to make them pass - fix the code instead
4. **REGRESSION PREVENTION**: Run full test suite after each task
5. **INCREMENTAL PROGRESS**: Complete one task fully before moving to the next

## Test Commands to Run After Each Task

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./cmd/...
go test ./internal/...
```

---

**Total Estimated Time:** ~6.5 hours across 9 tasks (excluding maybe someday items)