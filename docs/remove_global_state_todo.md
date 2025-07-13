# Remove Global State - Command-by-Command TODO List

This document tracks the elimination of global state for command flags. Each command should be refactored individually and tested before moving to the next.

## Implementation Order (Simple to Complex)

### ✅ Completed Commands
- [x] **cmd/push.go** - Removed `pushAll` global variable, updated to use `cmd.Flags().GetBool("all")`

### 🔄 In Progress
- [ ] None currently

### ⏳ Pending Commands


#### 1. cmd/pull.go (PRIORITY: HIGH - Simple)
**Global State to Remove:**
```go
var (
    pullAll bool
)
```

**Files to Modify:**
- [ ] **cmd/pull.go**:
  - [ ] Remove `var (pullAll bool)` (lines 12-14)
  - [ ] Update `RunE` function: Add `pullAll, _ := cmd.Flags().GetBool("all")` at start
  - [ ] Update `init()`: Change `pullCmd.Flags().BoolVar(&pullAll, "all", false, "Pull all worktrees")` to `pullCmd.Flags().Bool("all", false, "Pull all worktrees")`
- [ ] **cmd/pull_test.go**:
  - [ ] Remove all 9 instances of `pullAll = false` lines
  - [ ] Verify tests still pass

**Testing Checklist:**
- [ ] Run `go test ./cmd -run TestPull` before changes
- [ ] Run `go test ./cmd -run TestPull` after changes
- [ ] Run full test suite `go test ./...`
- [ ] Manual test: `gbm pull --all`
- [ ] Manual test: `gbm pull <worktree-name>`

---

#### 3. cmd/remove.go (PRIORITY: HIGH - Simple)
**Global State to Remove:**
```go
var (
    force bool
)
```

**Files to Modify:**
- [ ] **cmd/remove.go**:
  - [ ] Remove `var (force bool)` (lines 10-12)
  - [ ] Update `RunE` function: Add `force, _ := cmd.Flags().GetBool("force")` at start
  - [ ] Update `init()`: Change `removeCmd.Flags().BoolVarP(&force, "force", "f", false, "...")` to `removeCmd.Flags().BoolP("force", "f", false, "...")`
- [ ] **cmd/remove_test.go**:
  - [ ] Remove all 13 instances of `force = false` lines
  - [ ] Verify tests still pass

**Testing Checklist:**
- [ ] Run `go test ./cmd -run TestRemove` before changes
- [ ] Run `go test ./cmd -run TestRemove` after changes
- [ ] Run full test suite `go test ./...`
- [ ] Manual test: `gbm remove <worktree-name>`
- [ ] Manual test: `gbm remove <worktree-name> --force`

---

#### 4. cmd/switch.go (PRIORITY: HIGH - Simple)
**Global State to Remove:**
```go
var (
    printPath bool
)
```

**Files to Modify:**
- [ ] **cmd/switch.go**:
  - [ ] Remove `var (printPath bool)` (lines 14-16)
  - [ ] Update `RunE` function: Add `printPath, _ := cmd.Flags().GetBool("print-path")` at start
  - [ ] Update `init()`: Change flag registration to use `.Bool()` instead of `.BoolVar()`
- [ ] **cmd/switch_test.go**:
  - [ ] Remove 2 instances of `printPath = false` lines
  - [ ] Verify tests still pass

**Testing Checklist:**
- [ ] Run `go test ./cmd -run TestSwitch` before changes
- [ ] Run `go test ./cmd -run TestSwitch` after changes
- [ ] Run full test suite `go test ./...`
- [ ] Manual test: `gbm switch <worktree-name>`
- [ ] Manual test: `gbm switch <worktree-name> --print-path`

---

#### 5. cmd/sync.go (PRIORITY: MEDIUM - Two flags)
**Global State to Remove:**
```go
var (
    syncDryRun bool
    syncForce  bool
)
```

**Files to Modify:**
- [ ] **cmd/sync.go**:
  - [ ] Remove `var (syncDryRun bool, syncForce bool)` (lines 12-15)
  - [ ] Update `RunE` function: Add at start:
    ```go
    syncDryRun, _ := cmd.Flags().GetBool("dry-run")
    syncForce, _ := cmd.Flags().GetBool("force")
    ```
  - [ ] Update `init()`: Change to use `.Bool()` instead of `.BoolVar()`
- [ ] **cmd/sync_test.go**:
  - [ ] Remove `resetSyncFlags()` function (lines 18-21)
  - [ ] Remove all calls to `resetSyncFlags()`
  - [ ] Verify tests still pass

**Testing Checklist:**
- [ ] Run `go test ./cmd -run TestSync` before changes
- [ ] Run `go test ./cmd -run TestSync` after changes
- [ ] Run full test suite `go test ./...`
- [ ] Manual test: `gbm sync`
- [ ] Manual test: `gbm sync --dry-run`
- [ ] Manual test: `gbm sync --force`

---

#### 6. cmd/add.go (PRIORITY: MEDIUM - Complex logic)
**Global State to Remove:**
```go
var (
    newBranch   bool
    baseBranch  string
    interactive bool
)
```

**Files to Modify:**
- [ ] **cmd/add.go**:
  - [ ] Remove `var (newBranch bool, baseBranch string, interactive bool)` (lines 12-16)
  - [ ] Update `RunE` function: Add at start:
    ```go
    newBranch, _ := cmd.Flags().GetBool("new-branch")
    interactive, _ := cmd.Flags().GetBool("interactive")
    ```
  - [ ] Note: `baseBranch` comes from `args[2]`, not flags
  - [ ] Update `init()`: Change to use `.Bool()` instead of `.BoolVar()`
  - [ ] Update `handleInteractive()` function to pass `newBranch` as parameter instead of using global
- [ ] **cmd/add_test.go**:
  - [ ] Remove all 30+ instances of flag reset lines:
    ```go
    newBranch = false
    interactive = false
    baseBranch = ""
    ```
  - [ ] Verify tests still pass

**Testing Checklist:**
- [ ] Run `go test ./cmd -run TestAdd` before changes
- [ ] Run `go test ./cmd -run TestAdd` after changes
- [ ] Run full test suite `go test ./...`
- [ ] Manual test: `gbm add <worktree> <branch>`
- [ ] Manual test: `gbm add <worktree> <branch> -b`
- [ ] Manual test: `gbm add <worktree> --interactive`

---

#### 7. cmd/root.go (PRIORITY: LOW - Most Complex, Do Last)
**Global State to Remove:**
```go
var (
    configPath  string
    worktreeDir string
    debug       bool
    logFile     *os.File
)
```

**Files to Modify:**
- [ ] **cmd/root.go**:
  - [ ] Remove global variables (lines 15-20)
  - [ ] Update persistent flag registration in `init()`:
    ```go
    rootCmd.PersistentFlags().String("config", "", "specify custom .envrc path")
    rootCmd.PersistentFlags().String("worktree-dir", "", "override worktree directory location")
    rootCmd.PersistentFlags().Bool("debug", false, "enable debug logging to ./gbm.log")
    ```
  - [ ] Create helper functions that accept `cmd *cobra.Command`:
    ```go
    func getConfigPath(cmd *cobra.Command) string
    func getWorktreeDir(cmd *cobra.Command) string
    func isDebugEnabled(cmd *cobra.Command) bool
    ```
  - [ ] Update all existing helper functions to accept cmd parameter:
    - [ ] `createInitializedManager(cmd *cobra.Command)`
    - [ ] `createInitializedManagerStrict(cmd *cobra.Command)`
    - [ ] `createInitializedGitManager(cmd *cobra.Command)`
  - [ ] Update print functions to check debug flag from cmd
  - [ ] Update `PersistentPreRun` to pass cmd to functions
  - [ ] Handle `logFile` appropriately (local variable or context)

**Testing Checklist:**
- [ ] Run full test suite before changes
- [ ] Run full test suite after changes
- [ ] Manual test all commands with --debug flag
- [ ] Manual test all commands with --config flag
- [ ] Manual test all commands with --worktree-dir flag

**⚠️ WARNING: This change affects ALL other commands. Only do this after commands 1-6 are complete.**

---

#### 8. Commands Using Helper Functions (Do After root.go)

These can be done in any order after root.go is complete:

##### 8a. cmd/clone.go
- [ ] **cmd/clone.go**:
  - [ ] Update `GetConfigPath()` → `getConfigPath(cmd)` calls
  - [ ] Update any print function calls that now need cmd parameter
- [ ] **Testing**: Manual test `gbm clone <repo>` with various flags

##### 8b. cmd/info.go
- [ ] **cmd/info.go**:
  - [ ] Update helper function calls to pass `cmd` parameter
  - [ ] Update print function calls if needed
- [ ] **Testing**: Manual test `gbm info`

##### 8c. cmd/list.go
- [ ] **cmd/list.go**:
  - [ ] Update `createInitializedManager()` → `createInitializedManager(cmd)`
  - [ ] Update any print function calls
- [ ] **Testing**: Manual test `gbm list`

##### 8d. cmd/validate.go
- [ ] **cmd/validate.go**:
  - [ ] Update helper function calls to pass `cmd` parameter
- [ ] **Testing**: Manual test `gbm validate`

##### 8e. cmd/hotfix.go
- [ ] **cmd/hotfix.go**:
  - [ ] Update helper function calls to pass `cmd` parameter
  - [ ] Update print function calls
- [ ] **Testing**: Manual test `gbm hotfix` commands

##### 8f. cmd/mergeback.go
- [ ] **cmd/mergeback.go**:
  - [ ] Update helper function calls to pass `cmd` parameter
  - [ ] Update print function calls
- [ ] **Testing**: Manual test `gbm mergeback` commands

---

## General Guidelines for Each Command

### Before Starting Any Command:
1. [ ] Read the current implementation
2. [ ] Run existing tests to ensure they pass
3. [ ] Identify all global variables used
4. [ ] Identify all test files that reset flags

### During Implementation:
1. [ ] Remove global variable declarations
2. [ ] Add `cmd.Flags().Get*()` calls at start of `RunE` function
3. [ ] Update flag registration in `init()` function
4. [ ] Remove flag reset lines from test files
5. [ ] Test incrementally

### After Completing Each Command:
1. [ ] Run command-specific tests
2. [ ] Run full test suite
3. [ ] Manual testing with various flag combinations
4. [ ] Commit changes before moving to next command
5. [ ] Update this TODO list to mark command as completed

## Success Criteria

A command is considered complete when:
- [ ] All global variables for that command are removed
- [ ] Flag access uses `cmd.Flags().Get*()` methods only
- [ ] All tests pass for that command
- [ ] Full test suite still passes
- [ ] Manual testing confirms flags work correctly
- [ ] No global state remains for that command

## Notes for Future Implementation

- **Error Handling**: `cmd.Flags().Get*()` returns `(value, error)`. Since we control all flag definitions, using `_` to ignore the error is acceptable.
- **Backwards Compatibility**: This refactoring maintains all existing flag behavior and CLI interfaces.
- **Testing Strategy**: The existing test suite uses `rootCmd.SetArgs()` and `cmd.Execute()`, which naturally works with the new flag approach.
- **Rollback Strategy**: Each command can be reverted independently if issues arise.

## Lessons Learned from cmd/push.go Implementation

- **Flag State Pollution**: Cobra commands maintain flag state between executions when reusing the same command instance. Tests that fail unexpectedly may be due to flags remaining set from previous test runs.
- **Test Isolation Fix**: Add explicit flag resets in tests where needed using `commandName.Flags().Set("flag-name", "default-value")` before executing commands that expect clean flag state.
- **Test Assertion Cleanup**: When removing global variables, also remove test assertions that checked the global state (e.g., `assert.True(t, globalFlag, "message")`). These assertions are no longer possible or necessary after the refactoring.
- **Flag Access Pattern**: The pattern `flagValue, _ := cmd.Flags().GetBool("flag-name")` should be added at the very beginning of the `RunE` function, before any other logic.

## Final Verification

After all commands are complete:
- [ ] No global variables remain in any cmd/*.go files (except rootCmd itself)
- [ ] All tests pass: `go test ./...`
- [ ] All manual testing scenarios work
- [ ] No flag reset code remains in test files
- [ ] CLI behavior is identical to before refactoring