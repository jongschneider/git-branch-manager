# Remove Global State - Command-by-Command TODO List

This document tracks the elimination of global state for command flags. Each command should be refactored individually and tested before moving to the next.

## Implementation Order (Simple to Complex)

### ✅ Completed Commands
- [x] **cmd/push.go** - Removed `pushAll` global variable, updated to use `cmd.Flags().GetBool("all")`
- [x] **cmd/pull.go** - Removed `pullAll` global variable, updated to use `cmd.Flags().GetBool("all")`, fixed flag state pollution in tests

### 🔄 In Progress
- [ ] None currently

### ⏳ Pending Commands


#### ~~1. cmd/pull.go (PRIORITY: HIGH - Simple) - COMPLETED~~
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

**CRITICAL: Convert to Factory Function Pattern (like cmd/push.go and cmd/pull.go)**

**Files to Modify:**
- [ ] **cmd/remove.go**:
  - [ ] Convert `var removeCmd = &cobra.Command{...}` to `func newRemoveCommand() *cobra.Command { cmd := &cobra.Command{...}`
  - [ ] Move flag registration and completion setup inside the function before `return cmd`
  - [ ] Update `RunE` function: Add `force, _ := cmd.Flags().GetBool("force")` at start
  - [ ] Change flag registration to `cmd.Flags().BoolP("force", "f", false, "...")`
- [ ] **cmd/root.go**:
  - [ ] Update `rootCmd.AddCommand(removeCmd)` to `rootCmd.AddCommand(newRemoveCommand())`
- [ ] **cmd/remove_test.go**:
  - [ ] Remove all 13 instances of `force = false` lines (no longer needed with factory function)
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

**CRITICAL: Convert to Factory Function Pattern (like cmd/push.go and cmd/pull.go)**

**Files to Modify:**
- [ ] **cmd/switch.go**:
  - [ ] Convert `var switchCmd = &cobra.Command{...}` to `func newSwitchCommand() *cobra.Command { cmd := &cobra.Command{...}`
  - [ ] Move flag registration and completion setup inside the function before `return cmd`
  - [ ] Update `RunE` function: Add `printPath, _ := cmd.Flags().GetBool("print-path")` at start
  - [ ] Change flag registration to `cmd.Flags().Bool("print-path", false, "...")`
- [ ] **cmd/root.go**:
  - [ ] Update `rootCmd.AddCommand(switchCmd)` to `rootCmd.AddCommand(newSwitchCommand())`
- [ ] **cmd/switch_test.go**:
  - [ ] Remove 2 instances of `printPath = false` lines (no longer needed with factory function)
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

**CRITICAL: Convert to Factory Function Pattern (like cmd/push.go and cmd/pull.go)**

**Files to Modify:**
- [ ] **cmd/sync.go**:
  - [ ] Convert `var syncCmd = &cobra.Command{...}` to `func newSyncCommand() *cobra.Command { cmd := &cobra.Command{...}`
  - [ ] Move flag registration inside the function before `return cmd`
  - [ ] Update `RunE` function: Add at start:
    ```go
    syncDryRun, _ := cmd.Flags().GetBool("dry-run")
    syncForce, _ := cmd.Flags().GetBool("force")
    ```
  - [ ] Change flag registration to use `cmd.Flags().Bool()`
- [ ] **cmd/root.go**:
  - [ ] Update `rootCmd.AddCommand(syncCmd)` to `rootCmd.AddCommand(newSyncCommand())`
- [ ] **cmd/sync_test.go**:
  - [ ] Remove `resetSyncFlags()` function (no longer needed with factory function)
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

**CRITICAL: Convert to Factory Function Pattern (like cmd/push.go and cmd/pull.go)**

**Files to Modify:**
- [ ] **cmd/add.go**:
  - [ ] Convert `var addCmd = &cobra.Command{...}` to `func newAddCommand() *cobra.Command { cmd := &cobra.Command{...}`
  - [ ] Move flag registration and completion setup inside the function before `return cmd`
  - [ ] Update `RunE` function: Add at start:
    ```go
    newBranch, _ := cmd.Flags().GetBool("new-branch")
    interactive, _ := cmd.Flags().GetBool("interactive")
    ```
  - [ ] Note: `baseBranch` comes from `args[2]`, not flags
  - [ ] Change flag registration to use `cmd.Flags().Bool()`
  - [ ] Update `handleInteractive()` function to pass `newBranch` as parameter instead of using global
- [ ] **cmd/root.go**:
  - [ ] Update `rootCmd.AddCommand(addCmd)` to `rootCmd.AddCommand(newAddCommand())`
- [ ] **cmd/add_test.go**:
  - [ ] Remove all 30+ instances of flag reset lines (no longer needed with factory function):
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

### During Implementation - CRITICAL FACTORY FUNCTION PATTERN:
1. [ ] **Convert to Factory Function**: Change `var cmdName = &cobra.Command{...}` to `func newCmdNameCommand() *cobra.Command { cmd := &cobra.Command{...}`
2. [ ] **Move all setup inside function**: Move flag registration, completion setup, etc. before `return cmd`
3. [ ] **Update RunE function**: Add `flagVar, _ := cmd.Flags().Get*()` calls at start
4. [ ] **Update root.go**: Change `rootCmd.AddCommand(cmdName)` to `rootCmd.AddCommand(newCmdNameCommand())`
5. [ ] **Remove global variables**: Delete all global flag variables
6. [ ] **Remove flag reset lines**: Delete all test flag reset code (no longer needed)
7. [ ] **Test incrementally**: Verify each step works

### Why Factory Function Pattern:
- **Eliminates global state pollution** between test runs
- **Each command execution gets fresh state**
- **No manual flag resets needed** in tests
- **Follows established pattern** from cmd/push.go and cmd/pull.go
- **Cleaner, more testable code**

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

## Lessons Learned from cmd/push.go and cmd/pull.go Implementation

### CRITICAL: Factory Function Pattern is the ONLY Correct Approach
- **Factory Functions Eliminate All Issues**: Convert `var cmdName = &cobra.Command{...}` to `func newCmdNameCommand() *cobra.Command { ... }`
- **Root Cause Solution**: Factory functions create fresh command instances for each execution, completely eliminating global state pollution
- **No Manual Flag Resets Needed**: With factory functions, test flag reset code becomes unnecessary and should be removed

### Why Other Approaches Fail:
- **Flag State Pollution**: Cobra commands maintain flag state between executions when reusing the same command instance
- **Manual Flag Resets are Band-aids**: Using `cmd.Flags().Set("flag-name", "false")` in tests is a workaround, not a solution
- **Global Variables Create Test Dependencies**: Any global command variables cause state pollution between test runs

### Correct Implementation Pattern:
1. **Factory Function**: `func newCmdNameCommand() *cobra.Command { cmd := &cobra.Command{...}; return cmd }`
2. **Flag Access**: `flagValue, _ := cmd.Flags().GetBool("flag-name")` at start of `RunE` function
3. **Root Registration**: `rootCmd.AddCommand(newCmdNameCommand())` in root.go
4. **Clean Tests**: Remove all flag reset code from test files
5. **No Global Variables**: Zero global command or flag variables

### Test Assertion Cleanup:
- Remove test assertions that checked global state (e.g., `assert.True(t, globalFlag, "message")`)
- These assertions are impossible and unnecessary with factory functions

## Final Verification

After all commands are complete:
- [ ] No global variables remain in any cmd/*.go files (except rootCmd itself)
- [ ] All tests pass: `go test ./...`
- [ ] All manual testing scenarios work
- [ ] No flag reset code remains in test files
- [ ] CLI behavior is identical to before refactoring