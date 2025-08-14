# GBM Init Command Implementation Plan

## Overview
Implement `gbm init` command to initialize a new git repository with gbm worktree management structure, following the same patterns established by `gbm clone`.

## Command Signature
```bash
gbm init [directory] [--branch=<branch-name>]
```

## Requirements Summary
- **Directory Handling**: Create if missing, validate if exists (fail if already has git repo)
- **Git Repository Validation**: Fail if current/target directory is already in a git repository  
- **Branch Name Logic**: Respect native git default branch detection, allow `--branch` override
- **Repository Structure**: Initialize bare git repo with worktree-based development setup
- **Initial Commit**: Create initial commit after all scaffolding is complete
- **Configuration**: Set up `gbm.branchconfig.yaml`, `.gbm/` directory and files

## Architecture Design

### Interface-Based Approach
Following existing patterns (`worktreePuller`, `worktreeAdder`, etc.), create:

```go
//go:generate go tool moq -out ./autogen_repositoryInitializer.go . repositoryInitializer

// repositoryInitializer interface abstracts the Manager operations needed for initializing repositories
type repositoryInitializer interface {
    AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error
    SaveConfig() error
    SaveState() error
    GetRepoPath() string
    GetConfig() *internal.Config
}
```

### Implementation Structure
- **Command Logic**: `cmd/init.go` - handles CLI interface, argument parsing, orchestration
- **Phase 1**: Standalone helper functions for pre-repository setup (validation, git init)
- **Phase 2**: Manager interface methods for gbm-specific setup (worktrees, config, state)
- **Testing**: Interface mocking via moq for unit tests

## Implementation Tasks

### ✅ Completed
- [x] Research git native default branch handling
- [x] Design command structure with directory and --branch flag  
- [x] Study existing command patterns for interface usage
- [x] Design interfaces for init-specific operations

### ✅ Completed
- [x] Create implementation plan document
- [x] Implement git repository detection/validation logic
- [x] Create directory handling logic (create if missing, validate if exists)
- [x] Initialize bare git repository with proper configuration
- [x] Detect/respect native git default branch or use --branch flag
- [x] Reuse existing scaffolding functions from clone.go
- [x] Create initial commit after all scaffolding is complete
- [x] Add init command to root command structure

### ✅ Completed
- [x] Test and debug initial implementation
- [x] Refactor to follow pull.go pattern with clean interfaces
- [x] Use Manager.AddWorktree() for consistency
- [x] Use cmp.Or for succinct code in getNativeDefaultBranch
- [x] Run code quality checks (format, validate, vet) - all passed
- [x] Test in real production environment (api_test directory)
- [x] Test via tmux session for user workflow validation

### ✅ Completed
- [x] Add tests for init command (unit tests with mocked interface)

## ✅ Implementation Complete & Production Ready!

The `gbm init` command has been successfully implemented and tested! 

### What Works:
✅ **Directory validation** - Prevents initialization in existing git repositories  
✅ **Directory creation** - Creates target directory if it doesn't exist  
✅ **Bare repository setup** - Initializes proper `.git` structure  
✅ **Branch name detection** - Respects `git config init.defaultBranch` and `--branch` flag  
✅ **Worktree structure** - Creates `worktrees/` directory and main worktree  
✅ **GBM configuration** - Generates `gbm.branchconfig.yaml` with initial branch setup  
✅ **State management** - Creates `.gbm/` directory with `config.toml` and `state.toml`  
✅ **Initial commit** - Commits the `gbm.branchconfig.yaml` file  

### Example Usage:
```bash
gbm init                    # Initialize in current directory 
gbm init my-project         # Initialize in 'my-project' directory
gbm init --branch=develop   # Use 'develop' as default branch
```

### Test Results:
```
Repository initialized successfully!
Main worktree available at: /tmp/test-init-dir/worktrees/main

Created structure:
├── .git/                    # bare repository
├── worktrees/main/         # main worktree with initial commit
├── gbm.branchconfig.yaml   # branch configuration
└── .gbm/                   # gbm state directory
    ├── config.toml
    └── state.toml
```

## Implementation Flow

### 1. Command Entry Point (`cmd/init.go`)
```go
func newInitCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "init [directory] [--branch=<branch-name>]",
        Short: "Initialize a new git repository with gbm structure",
        Args:  cobra.MaximumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            branchFlag, _ := cmd.Flags().GetString("branch")
            initializer := createRepositoryInitializer()
            targetDir := resolveTargetDirectory(args)
            return handleInit(initializer, targetDir, branchFlag)
        },
    }
    cmd.Flags().String("branch", "", "Override default branch name")
    return cmd
}
```

### 2. Initialization Logic
```go
func handleInit(targetDir, branchFlag string) error {
    // Phase 1: Pre-repository setup (standalone functions)
    validateInitDirectory(targetDir)
    createInitDirectory(targetDir) 
    initializeBareRepository(targetDir)
    branchName := cmp.Or(branchFlag, getNativeDefaultBranch())
    
    // Phase 2: Manager-based setup (interface methods)
    manager := internal.NewManager(targetDir)
    setupWorktreeStructure(manager, branchName)
    createGBMConfig(manager, branchName)
    initializeGBMState(manager)
    createInitialCommit(manager, branchName)
}
```

### 3. Key Implementation Details
- **Clean Interface**: Only Manager operations that need mocking (`AddWorktree`, `SaveConfig`, `SaveState`, etc.)
- **Standalone Functions**: Pre-repository operations (`validateInitDirectory`, `initializeBareRepository`, etc.)
- **Modern Go**: Uses `cmp.Or` for concise fallback logic
- **Consistency**: Uses `manager.AddWorktree()` same as `gbm add` command

### 4. Integration
- Add command to `cmd/root.go`
- Ensure proper error handling and user feedback
- Add shell completion support

## Technical Decisions

### Branch Name Resolution
1. If `--branch` flag provided → use that branch name
2. If no flag provided → check `git config --get init.defaultBranch`  
3. Fall back to git's built-in default (usually "main")

### Directory Validation
- Current directory: Fail if already in git repository
- Target directory (if specified):
  - Create if doesn't exist
  - If exists and empty → proceed
  - If exists and non-empty → check for git repo:
    - Has git repo → fail
    - No git repo → proceed

### Repository Structure
Mirror `gbm clone` behavior:
```
target-directory/
├── .git/                    # bare repository
├── worktrees/
│   └── main/               # main worktree (or custom branch name)
├── gbm.branchconfig.yaml   # branch configuration
└── .gbm/                   # gbm state directory
    ├── config.toml
    └── state.toml
```

### Initial Commit Content
Create initial commit after all scaffolding with:
- `gbm.branchconfig.yaml` file
- Any other default files created during setup

## Open Questions
~~1. **Initial commit content**: Should we include `gbm.branchconfig.yaml` in the initial commit or keep it empty?~~
✅ **ANSWERED**: Yes, include `gbm.branchconfig.yaml` in the initial commit.

~~2. **Manager dependency**: How to handle manager creation for a new repository that doesn't exist yet?~~
✅ **ANSWERED**: Option 1 - Create Manager after bare repository initialization. Flow: `validate directory → init bare repo → create Manager → worktree setup → gbm config`

~~3. **Error recovery**: If initialization fails partway through, should we clean up created directories/files?~~
✅ **ANSWERED**: Do not clean up on failure - leave files/directories for user to handle manually.

## Quality & Testing

### ✅ Code Quality Validation
- **Format**: `just format` - ✅ All files properly formatted
- **Vet**: `just vet` - ✅ No issues found in cmd/internal packages  
- **Lint**: `golangci-lint` - ✅ 0 issues across all packages
- **Build**: ✅ Compilation successful
- **Existing Tests**: ✅ All existing tests continue to pass

### ✅ Real-World Testing
- **Production Directory**: Successfully tested in `/Users/jschneider/code/scratch/api_test`
- **Existing Code**: Works seamlessly with existing project files and directories
- **Tmux Integration**: Tested via tmux session for real user workflow
- **Repository Validation**: Correctly prevents reinitializing existing repositories

### Testing Strategy for Next Session
- Unit tests using mocked `repositoryInitializer` interface
- Integration tests with temporary directories  
- Error case testing (already in git repo, directory validation failures)
- Branch name resolution testing (native vs flag override)

## Architecture Improvements Made

### 🔄 Refactoring to pull.go Pattern
- **Before**: Complex monolithic interface with 10+ methods
- **After**: Clean interface focused on Manager operations (5 methods)
- **Benefit**: Follows established codebase patterns perfectly

### 🎯 Key Optimizations
1. **Interface Simplification**: `repositoryInitializer` now matches `worktreePuller` pattern
2. **Code Reuse**: Uses `manager.AddWorktree()` instead of duplicating logic  
3. **Modern Go**: Uses `cmp.Or` for concise fallback logic
4. **Clean Separation**: Phase 1 (standalone) vs Phase 2 (Manager interface)

## Final Status: PRODUCTION READY ✅

The `gbm init` command is:
- ✅ **Fully implemented** with all requirements met
- ✅ **Thoroughly tested** in real-world scenarios  
- ✅ **Code quality validated** with all checks passing
- ✅ **Architecture consistent** with existing codebase patterns
- ✅ **User validated** through tmux session testing

## Unit Tests Completed ✅

Added comprehensive unit tests for the `gbm init` command including:
- **Mock generation**: Created `autogen_repositoryInitializer.go` using `github.com/matryer/moq`
- **Function coverage**: All key functions tested with table-driven patterns:
  - `resolveTargetDirectory()` - directory path resolution
  - `validateInitDirectory()` - git repository detection 
  - `setupWorktreeStructure()` - worktree creation via interface
  - `createGBMConfig()` - configuration file generation
  - `initializeGBMState()` - state management
  - `createInitialCommit()` - git commit logic
  - `getNativeDefaultBranch()` - branch name detection
- **Command testing**: Full command structure validation
- **Error scenarios**: Comprehensive error handling coverage
- **Interface validation**: Proper mock interaction verification

All tests pass successfully, maintaining compatibility with existing test suite.

---
*Last Updated: 2025-08-13*