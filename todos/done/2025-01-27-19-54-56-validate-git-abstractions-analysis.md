Based on my comprehensive analysis of the Git command abstractions in the worktree-manager codebase, here's what I found:

## Git Command Abstractions Analysis

### 1. Git Utility Functions Location and Structure

**Primary Location**: `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`

This is the main Git utilities file containing:

**Core Git Command Wrappers:**
- `ExecGitCommand(dir string, args ...string)` - Main utility for executing git commands with output capture
- `execGitCommandRun(dir string, args ...string)` - For commands without output capture  
- `ExecGitCommandCombined(dir string, args ...string)` - Returns combined stdout/stderr
- `ExecGitCommandInteractive(dir string, args ...string)` - For interactive commands

**High-Level Git Operations:**
- `GetCurrentBranch()`, `GetCurrentBranchInPath(path)`
- `GetUpstreamBranch(worktreePath)`, `GetAheadBehindCount(worktreePath)`
- `VerifyRef(ref)`, `VerifyRefInPath(path, ref)`
- `GetCommitHash(ref)`, `GetCommitHashInPath(path, ref)`
- `GetCommitHistory(path, options)` - Flexible commit retrieval
- `GetFileChanges(path, options)` - Staged/unstaged file changes
- `BranchExists()`, `BranchExistsLocalOrRemote()`
- `GetWorktrees()`, `GetWorktreeStatus()`
- Worktree management: `CreateWorktree()`, `MoveWorktree()`, `UpdateWorktree()`

**Specialized Modules:**
- `/Users/jschneider/code/scratch/worktree-manager/internal/git_add.go` - Worktree addition functions
- `/Users/jschneider/code/scratch/worktree-manager/internal/git_remove.go` - Worktree removal functions

### 2. Remaining Direct Git Command Calls

**Production Code Files with Direct Calls:**
- `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go:208` - `git merge-base --is-ancestor` call
- `/Users/jschneider/code/scratch/worktree-manager/cmd/clone.go` - Multiple git commands for repository setup (lines 78, 95, 102, 132, 135, 158, 183)
- `/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback.go:177,181` - `git log` calls for branch comparison
- `/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback.go:603` - `git checkout` call
- `/Users/jschneider/code/scratch/worktree-manager/internal/git.go` - Some specialized calls in `FindGitRoot()`, `PushWorktree()`, `IsInWorktree()`

**Test Files with Direct Calls (Acceptable):**
- Multiple test files in `/cmd/*_test.go` use direct git commands for test setup and verification
- `/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go` - Test harness with git commands

### 3. Test Suite Structure

**Test Organization:**
- 25 test files across the codebase (`*_test.go`)
- Comprehensive test coverage for major commands
- Integration tests using real git repositories
- Test utilities in `/Users/jschneider/code/scratch/worktree-manager/internal/testutils/`

**Test Execution:**
- **Justfile-based**: `just test` (all tests), `just test-changed` (changed packages only)
- **Go standard**: `go test ./...`
- **With timeout**: Tests run with 10-minute timeout
- **Validation pipeline**: `just validate` runs format, vet, lint, build, test-changed

**Test Utilities:**
- `GitTestRepo` - Test repository setup and management
- Mock services for external dependencies
- Scenario-based test configurations
- YAML-based test configurations

### 4. Recent Refactoring Work

**Completed Refactoring (from commit history and todos/done/):**
- ✅ **GetCurrentBranch utilities** - Added both repo and path-specific versions
- ✅ **GetUpstreamBranch utility** - Branch upstream detection  
- ✅ **GetAheadBehindCount utility** - Commit counting vs upstream
- ✅ **VerifyRef utilities** - Reference validation with error handling
- ✅ **GetCommitHash utilities** - Commit hash retrieval for refs
- ✅ **GetCommitHistory utility** - Flexible commit history retrieval
- ✅ **GetFileChanges utility** - File change detection (staged/unstaged)
- ✅ **Info command refactoring** - Extracted git operations to utilities

**Current State:**
The codebase has undergone significant refactoring to centralize Git operations. Most duplicate git command patterns have been eliminated, with remaining direct calls being either:
1. Specialized operations that don't fit the utility pattern yet
2. Test code (which is acceptable)
3. Complex operations in the clone command
4. Edge cases in mergeback operations

**Still Needing Refactoring:**
- Clone command git operations (complex repository setup)
- Some mergeback git operations (branch comparison, checkout)
- Specialized operations in info command (merge-base checks)

### 5. Build and Test Execution

**Primary Commands:**
```bash
# Run all tests
just test

# Run tests for changed files only  
just test-changed

# Full validation pipeline
just validate

# Individual operations
just format    # Format changed files
just vet      # Vet changed packages  
just lint     # Lint changed packages
just build    # Build project
```

**Test Features:**
- Incremental testing (only changed packages)
- Comprehensive validation pipeline
- 10-minute timeout for long-running tests
- Integration with git workflow (detects changed files)

The codebase demonstrates a well-architected approach to Git command abstraction with strong test coverage and systematic refactoring progress. The remaining direct git calls are mostly in specialized areas that require custom handling rather than generic utilities.