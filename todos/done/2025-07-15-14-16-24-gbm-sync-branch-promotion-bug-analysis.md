Now I have enough information to provide a comprehensive analysis. Based on my research of the codebase, I can identify the key findings about the `gbm sync` bug when promoting branches to new names.

## Key Findings

### 1. **Sync Command Implementation**
- **Location**: `/Users/jschneider/code/scratch/worktree-manager/cmd/sync.go`
- **Entry Point**: `newSyncCommand()` function calls `manager.SyncWithConfirmation()`
- **Core Logic**: The command fetches from remote, validates config, then processes missing worktrees, branch changes, and orphaned worktrees

### 2. **Worktree Update Logic** 
- **Location**: `/Users/jschneider/code/scratch/worktree-manager/internal/manager.go:212-218`
- **Process**: For branch changes, calls `m.gitManager.UpdateWorktree(worktreePath, change.NewBranch)`
- **Update Method**: `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:349-359`

### 3. **Root Cause of Exit Status 128**

The `UpdateWorktree` function follows this sequence:
1. **Remove existing worktree**: `gm.RemoveWorktree(worktreePath)` 
2. **Recreate with new branch**: `gm.CreateWorktree(envVar, newBranch, relativeWorktreeDir)`

**The bug occurs in step 2** when `CreateWorktree` tries to create a worktree for a branch that may:

#### CreateWorktree Process (`internal/git.go:299-339`):
1. Checks if directory exists (line 302-304)
2. Validates branch exists with `BranchExists()` (line 306-312)
3. **Critical Issue**: Executes `git worktree add <path> <branch>` (lines 321, 333)

#### Exit Status 128 Sources:
Git returns exit status 128 when:
- **Branch is already checked out in another worktree** (most likely cause in promotion scenario)
- **Branch doesn't exist locally** (even if it exists remotely)
- **Directory already exists** (though this is checked)
- **Git internal conflicts** with worktree state

### 4. **Specific Bug in Branch Promotion Scenario**

From the task description, the issue occurs when:
- **Before**: `production` worktree → `production-2025-05-1` branch  
- **After**: `production` worktree → `production-2025-07-1` branch

**Problem**: When `UpdateWorktree` removes the old worktree and tries to create a new one:
1. `RemoveWorktree` successfully removes `production` worktree from `production-2025-05-1`
2. `CreateWorktree` tries to create `production` worktree on `production-2025-07-1`
3. **Git fails with exit status 128** if:
   - The branch `production-2025-07-1` is already checked out elsewhere
   - The branch doesn't exist locally (only remotely)
   - Git worktree internal state is inconsistent

### 5. **Missing Validation Logic**

**Current Issues**:
- `BranchExists()` checks both local and remote branches but may not handle all edge cases
- No validation that the target branch isn't already checked out in another worktree
- No cleanup of Git's internal worktree references that might conflict
- No retry mechanism if the first attempt fails

### 6. **Error Handling Gaps**

**Current Error Path**:
```
manager.go:216 → failed to update worktree for %s: %w
git.go:322/334 → failed to create worktree: %w
exec.Command()  → exit status 128
```

**Missing**:
- Specific handling for "branch already checked out" errors
- Better error messages distinguishing different failure modes
- Validation before attempting the operation

### 7. **Potential Solutions**

The bug can be fixed by:

1. **Enhanced validation** before worktree operations
2. **Better error handling** for Git exit status 128 cases  
3. **Cleanup of Git state** between remove and recreate operations
4. **Checking for branch conflicts** before attempting updates

The core issue is that `UpdateWorktree` assumes a simple remove+recreate will work, but Git's worktree state management has edge cases that cause the recreation to fail with exit status 128.

Based on my comprehensive analysis of the codebase, here are my findings about Git worktree operations and exit status 128 errors:

## Git Worktree Operations in the Codebase

### 1. Git Worktree Command Usage

The codebase uses several Git worktree commands executed via `exec.Command`:

**Creation Commands:**
- `git worktree add <path> <branch>` - Add worktree for existing branch
- `git worktree add -b <branch> <path> [base-branch]` - Create new branch and worktree
- `git worktree add --track -b <branch> <path> <remote-branch>` - Create with remote tracking

**Management Commands:**
- `git worktree remove <path> --force` - Remove worktree forcefully
- `git worktree list --porcelain` - List all worktrees
- `git worktree prune` - Clean up stale worktree references

**File Locations:**
- `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:506-538` - AddWorktree method
- `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:321-343` - CreateWorktree and RemoveWorktree methods

### 2. Error Handling Patterns

The codebase uses consistent error handling patterns:

```go
// Unified Git command execution functions
func ExecGitCommand(dir string, args ...string) ([]byte, error)
func ExecGitCommandRun(dir string, args ...string) error
func ExecGitCommandCombined(dir string, args ...string) ([]byte, error)
```

All Git operations wrap errors with context:
```go
if err := ExecGitCommandRun(gm.repoPath, "worktree", "add", worktreePath, branchName); err != nil {
    return fmt.Errorf("failed to create worktree: %w", err)
}
```

### 3. Exit Status 128 Sources

Based on the log evidence (`/Users/jschneider/code/scratch/worktree-manager/gbm.log:3`) and task description, exit status 128 occurs specifically during worktree update operations:

**Root Cause Analysis:**
- **Primary Location:** `UpdateWorktree` method in `/Users/jschneider/code/scratch/worktree-manager/internal/git.go:349-359`
- **Failure Point:** The `CreateWorktree` call after successful `RemoveWorktree`
- **Specific Error:** `failed to update worktree for production: failed to create worktree: exit status 128`

**Common Git Exit Status 128 Causes:**
1. **Invalid Repository State** - Attempting Git worktree operations outside a valid Git repository
2. **Branch Reference Issues** - Target branch doesn't exist locally or remotely
3. **Path Conflicts** - Directory already exists or permissions issues
4. **Lock File Conflicts** - Previous Git operations left lock files
5. **Worktree Path Issues** - Invalid characters or path resolution problems

### 4. State Management and Cleanup

**Worktree Lifecycle:**
- **Creation:** Manager.AddWorktree → GitManager.AddWorktree → exec git worktree add
- **Removal:** Manager.RemoveWorktree → GitManager.RemoveWorktree → exec git worktree remove --force
- **Update:** GitManager.UpdateWorktree → RemoveWorktree + CreateWorktree

**State Tracking:**
- Ad-hoc worktrees tracked in `state.AdHocWorktrees`
- Configuration-based worktrees in `gbmConfig.Worktrees`
- Git status monitoring via `GetWorktreeStatus`

### 5. Validation and Error Prevention

**Pre-flight Checks:**
- Branch existence validation via `BranchExists()`
- Directory existence checks before creation
- Remote branch verification with `git rev-parse --verify origin/<branch>`
- Worktree path validation

**Error Recovery:**
- Force removal with `--force` flag to handle stale references
- State cleanup after failed operations
- Configuration validation before sync operations

### 6. Key Findings for Exit Status 128

**Critical Insight:** The exit status 128 occurs during the `gbm sync` operation when:
1. A worktree is being updated to point to a different branch
2. The old worktree is successfully removed
3. The new worktree creation fails with Git exit status 128

**Likely Causes in This Context:**
1. **Branch Reference Issues:** The target branch may not exist locally after the removal
2. **Repository State Corruption:** The removal operation may leave the Git repository in an inconsistent state
3. **Race Conditions:** Timing issues between remove and create operations
4. **Directory Cleanup:** Incomplete directory cleanup preventing recreation

**Error Pattern:**
```
UpdateWorktree → RemoveWorktree (SUCCESS) → CreateWorktree (FAIL: exit status 128)
```

The error handling preserves the underlying Git exit code and wraps it with contextual information, making it traceable through the call stack from the high-level `gbm sync` command down to the specific Git worktree operation that failed.

Based on my examination of the worktree-manager codebase, I can provide a comprehensive analysis of the state management system and the configuration promotion issues. Here are my findings:

## State Management System Analysis

### 1. State Loading and Validation

**State Structure (`internal/state.go`):**
- The system uses two main files: `.gbm/state.toml` for runtime state and `.gbm/config.toml` for configuration
- State includes: `LastSync`, `TrackedVars`, `AdHocWorktrees`, `CurrentWorktree`, `PreviousWorktree`, `LastMergebackCheck`
- Configuration includes settings, icons, JIRA config, and file copy rules
- Separate YAML-based branch configuration in `gbm.branchconfig.yaml`

**Loading Process:**
- `LoadState()` loads from `.gbm/state.toml` or returns default state if missing
- `LoadConfig()` loads from `.gbm/config.toml` or returns default config if missing
- `ParseGBMConfig()` loads the YAML branch configuration separately
- No migration logic exists - old state structures are preserved only for reference

### 2. Configuration Change Handling

**Sync Process (`internal/manager.go`, `GetSyncStatus()`):**
- Compares current worktrees against desired state in `gbm.branchconfig.yaml`
- Identifies three types of changes:
  - **Missing worktrees**: Defined in config but don't exist on filesystem
  - **Orphaned worktrees**: Exist on filesystem but not in config
  - **Branch changes**: Worktree exists but points to different branch than configured

**State Tracking:**
- `TrackedVars` in state.toml stores list of managed worktrees
- `AdHocWorktrees` tracks worktrees created outside of configuration
- Updates state after successful sync operations

### 3. Configuration Promotion Bug Analysis

**Root Cause (from task.md and code analysis):**
The bug occurs during branch promotion scenarios like:
```yaml
# From:
worktrees:
  production:
    branch: production-2025-05-1 
    description: "Old version"
  preview:
    branch: production-2025-07-1 
    description: "New version"

# To:
worktrees:
  production:
    branch: production-2025-07-1  # Promoted branch
    description: "New version"
```

**Technical Issue:**
- `UpdateWorktree()` function (`internal/git.go`) uses remove-then-recreate approach
- First calls `RemoveWorktree()` with `git worktree remove --force`
- Then calls `CreateWorktree()` to recreate with new branch
- Error occurs in recreation step: `"failed to create worktree: exit status 128"`

### 4. Potential Issues During Configuration Promotion

**Branch Conflicts:**
- If the target branch is already checked out in another worktree, git will refuse to create the new worktree
- The `UpdateWorktree()` logic assumes removal frees up the branch, but this may not always work

**Directory State Issues:**
- Race conditions between removal and recreation
- Filesystem cleanup delays
- Git's internal worktree registry may not be immediately updated

**Validation Gaps:**
- `ValidateConfig()` only checks if branches exist, not if they're available for checkout
- No pre-flight checks for branch availability during updates

### 5. State Validation and Conflict Handling

**Current Validation:**
- Basic branch existence checking via `BranchExists()`
- Directory existence validation
- YAML syntax validation for configuration files

**Missing Validation:**
- No check for branch checkout conflicts before worktree operations
- No verification that branches are available (not checked out elsewhere)
- No rollback mechanism if partial sync operations fail

**Error Handling:**
- Operations fail fast with descriptive error messages
- State updates only occur after successful operations
- Orphaned worktree handling requires explicit `--force` flag with confirmation

### 6. Recommendations for State Management Improvements

**Immediate Fixes Needed:**
1. **Pre-flight validation**: Check branch availability before attempting worktree updates
2. **Atomic operations**: Implement transaction-like behavior for multi-step updates
3. **Better error recovery**: Add rollback capability for failed operations
4. **Branch conflict detection**: Identify when target branches are already in use

**State Validation Enhancements:**
1. **Comprehensive validation**: Extend validation to check branch checkout status
2. **Conflict resolution**: Provide options for handling branch conflicts
3. **State consistency checks**: Validate state.toml against actual filesystem state
4. **Migration support**: Add proper state migration logic for config changes

**Configuration Change Safety:**
1. **Dry-run improvements**: Better simulation of complex change scenarios
2. **Confirmation dialogs**: More detailed information about risky operations
3. **Backup and restore**: Ability to snapshot and restore worktree states
4. **Incremental updates**: Process changes in safer, smaller steps

The core issue appears to be that the current `UpdateWorktree` implementation is too simplistic for complex branch promotion scenarios, where the same branch may be referenced by multiple worktree configurations during the transition.