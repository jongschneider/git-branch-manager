## Analysis: cmd/add.go Interface Extraction

### Current Structure and Dependencies

The `cmd/add.go` file is well-structured with the following key components:

**Main Structures:**
- `WorktreeArgs` - Contains resolved arguments for creating a worktree
- `ArgsResolver` - Handles complex logic for resolving command arguments and contains a reference to `*internal.Manager`

**Manager Methods Used:**
The add command uses these Manager methods:
- `manager.AddWorktree(worktreeName, branchName, newBranch, baseBranch)` - Primary worktree creation
- `manager.GetGitManager().GetDefaultBranch()` - Via ArgsResolver for base branch resolution  
- `manager.GetGitManager().BranchExists(branchName)` - Via ArgsResolver for branch validation

**JIRA Integration Methods:**
- `internal.GetJiraIssues(manager)` - For command completion
- `internal.GenerateBranchFromJira(jiraKey, manager)` - For auto-generating branch names
- `internal.IsJiraKey(worktreeName)` - For JIRA key detection

### ArgsResolver Function Analysis

The `ArgsResolver` is a sophisticated argument processing system:

**Key Methods:**
- `ResolveArgs(cmdArgs, newBranchFlag)` - Main entry point that orchestrates argument resolution
- `resolveBranchName(cmdArgs, newBranchFlag, worktreeName)` - Handles branch name logic
- `resolveBaseBranch(newBranchFlag, baseBranch)` - Validates and resolves base branches

**Logic Flow:**
1. **Worktree Name**: Always required (first argument)
2. **Branch Name**: Can be specified explicitly (second argument) or auto-generated
3. **Base Branch**: Optional (third argument), defaults to repository's default branch
4. **JIRA Integration**: Auto-detects JIRA keys and suggests/generates branch names

**Smart Features:**
- JIRA key detection with intelligent branch name suggestions
- Validation of base branch existence
- Fallback behavior when JIRA CLI isn't available
- Clear error messages with usage hints

### Manager Methods Called by add.go

**Direct Manager Methods:**
- `AddWorktree(worktreeName, branchName, createBranch, baseBranch)` - Creates the worktree
- `GetGitManager()` - Access to GitManager for lower-level operations

**GitManager Methods (via manager.GetGitManager()):**
- `GetDefaultBranch()` - Gets repository's default branch (main/master)
- `BranchExists(branchName)` - Validates branch existence
- `AddWorktree(worktreeName, branchName, createBranch, baseBranch)` - Core git worktree creation

**JIRA Integration Functions:**
- `GetJiraIssues(manager)` - Fetches JIRA issues for completion
- `GenerateBranchFromJira(jiraKey, manager)` - Auto-generates branch names from JIRA
- `IsJiraKey(worktreeName)` - Detects JIRA key patterns

### Current Testing Setup

**Comprehensive Test Coverage:**
The testing is quite robust with `/Users/jschneider/code/scratch/worktree-manager/cmd/add_test.go` containing:

- **Integration Tests**: Full end-to-end testing with real git operations
- **Edge Case Coverage**: Invalid branches, missing arguments, duplicate worktrees
- **JIRA Integration**: Tests for JIRA key detection and branch generation
- **Completion Testing**: Tests for tab completion functionality
- **Error Scenarios**: Comprehensive error handling validation

**Test Infrastructure:**
- Uses `testutils.NewMultiBranchRepo()` for test repository setup
- Leverages `testutils` package for consistent test environments
- Tests both success and failure scenarios
- Validates git operations and file system changes

**GitManager Tests:**
`/Users/jschneider/code/scratch/worktree-manager/internal/git_add_test.go` provides:
- Lower-level GitManager testing
- Worktree creation with various branch scenarios
- Base branch validation
- Error condition testing

### Integration with Manager Struct

**Current Manager Structure:**
```go
type Manager struct {
    config     *Config
    state      *State  
    gitManager *GitManager
    gbmConfig  *GBMConfig
    repoPath   string
    gbmDir     string
}
```

**Integration Points:**
- **Dependency Injection**: Manager is created in `createInitializedManager()` and passed to commands
- **Error Handling**: Graceful degradation when manager creation fails
- **Configuration**: Manager loads both internal config and gbm.branchconfig.yaml
- **State Management**: Tracks worktree state and metadata

### Existing Interfaces and Patterns

**Current State - No Formal Interfaces:**
The codebase currently uses concrete types (`*internal.Manager`) rather than interfaces. However, there are clear architectural patterns:

**Established Patterns:**
- **Command Pattern**: Each command is self-contained with its logic
- **Dependency Injection**: Manager is injected into commands during initialization
- **Error Propagation**: Consistent error handling and user feedback
- **Graceful Degradation**: Commands handle manager creation failures appropriately

**Documentation Shows Interface Planning:**
The `/Users/jschneider/code/scratch/worktree-manager/docs/interface-refactoring-plan.md` shows planned interfaces:
```go
type worktreeAdder interface {
    AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
    GetDefaultBranch() (string, error)
    BranchExists(branch string) (bool, error)
    GetJiraIssues() ([]internal.JiraIssue, error)
    GenerateBranchFromJira(jiraKey string) (string, error)
}
```

### Key Observations

**Strengths:**
- Well-structured argument resolution logic
- Comprehensive test coverage
- Robust JIRA integration with fallback behavior
- Clear separation of concerns between argument processing and execution
- Good error handling with helpful user messages

**Areas for Improvement:**
- No formal interfaces - currently tightly coupled to concrete Manager type
- JIRA functions are package-level functions rather than methods
- Some functions could benefit from better separation (mixing business logic with presentation)

**Architecture Readiness:**
The code is well-positioned for interface extraction and dependency injection improvements, as evidenced by the clear separation of concerns and existing dependency injection pattern for the Manager.

### Current Architecture Analysis

#### Existing Interfaces: None Currently Defined
- No interfaces are currently defined in the `internal` package
- The codebase uses concrete types (`*internal.Manager`, `*internal.GitManager`) throughout
- All commands directly depend on the concrete `Manager` type

#### Mock Generation Patterns: Ready for go tool moq
- No existing mock files found (no `autogen_*` files)
- The refactoring plan specifies using `//go:generate go tool moq` for mock generation
- Pattern will be: `//go:generate go tool moq -out ./autogen_interfaceName.go . interfaceName`
- Mock files will be generated in the `cmd` package alongside the command files

#### JIRA Integration Functions Structure

Located in `/Users/jschneider/code/scratch/worktree-manager/internal/jira.go`:

**Key Functions:**
- `GetJiraIssues(manager *Manager) ([]JiraIssue, error)` - Fetches JIRA issues for current user
- `GenerateBranchFromJira(jiraKey string, manager *Manager) (string, error)` - Creates branch name from JIRA
- `IsJiraKey(s string) bool` - Validates JIRA key format
- `GetJiraIssue(key string, manager *Manager) (*JiraIssue, error)` - Fetches single issue details

**Current Structure:**
- Functions take `*Manager` parameter to access config and caching
- Use `exec.Command("jira", ...)` to call external JIRA CLI
- Parse tabular and view output formats
- Cache JIRA user information in config

#### Manager and GitManager Types Structure

**Manager (in `/Users/jschneider/code/scratch/worktree-manager/internal/manager.go`):**
```go
type Manager struct {
    config     *Config
    state      *State  
    gitManager *GitManager
    gbmConfig  *GBMConfig
    repoPath   string
    gbmDir     string
}
```

**Key Manager Methods:**
- `AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error`
- `GetAllWorktrees() (map[string]*WorktreeListInfo, error)`
- `GetGitManager() *GitManager`
- `GetSyncStatus() (*SyncStatus, error)`
- `BranchExists(branchName string) (bool, error)`
- `RemoveWorktree(worktreeName string) error`
- File system operations for worktree management

**GitManager (in `/Users/jschneider/code/scratch/worktree-manager/internal/git.go`):**
```go
type GitManager struct {
    repo           *git.Repository
    repoPath       string
    worktreePrefix string
}
```

**Key GitManager Methods:**
- `GetWorktrees() ([]*WorktreeInfo, error)`
- `GetDefaultBranch() (string, error)`
- `BranchExists(branchName string) (bool, error)`
- `GetWorktreeStatus(worktreePath string) (*GitStatus, error)`
- `GetCurrentBranchInPath(path string) (string, error)`
- `GetCommitHistory(path string, options CommitHistoryOptions) ([]CommitInfo, error)`
- `GetFileChanges(path string, options FileChangeOptions) ([]FileChange, error)`

#### Current Testing Patterns

**Testing Structure:**
- Integration tests in `cmd/*_test.go` using real git repositories
- Test utilities in `internal/testutils/` for repo setup and mocking
- Existing mock for JIRA CLI in `internal/testutils/mock_services.go`
- Uses testify for assertions and test structure

**Current Test Dependencies:**
- Tests create real git repositories using `testutils.NewMultiBranchRepo(t)`
- Commands tested through Cobra's command execution
- Tests currently require git operations and file system access

#### Recommended Interface Patterns

Based on the analysis, the interfaces should follow these patterns:

**Command-Specific Interfaces:**
```go
// Example for add command
//go:generate go tool moq -out ./autogen_worktreeAdder.go . worktreeAdder
type worktreeAdder interface {
    AddWorktree(worktreeName, branchName string, newBranch bool, baseBranch string) error
    GetDefaultBranch() (string, error)
    BranchExists(branch string) (bool, error)
    GetJiraIssues() ([]internal.JiraIssue, error)
    GenerateBranchFromJira(jiraKey string) (string, error)
}
```

**JIRA Integration Interface:**
```go
type jiraProvider interface {
    GetJiraIssues() ([]JiraIssue, error)
    GenerateBranchFromJira(jiraKey string) (string, error)
    IsJiraKey(s string) bool
}
```

#### Migration Strategy Recommendations

1. **Start with `cmd/add.go`** as the example implementation
2. **Extract JIRA functions** to interface-based approach first
3. **Create thin wrapper methods** on Manager that implement interface methods
4. **Generate mocks** and convert integration tests to unit tests
5. **Apply pattern** to other commands in complexity order

The codebase is well-structured for interface extraction, with clear separation between git operations, JIRA integration, and command logic. The main work will be creating the interface definitions and updating the ArgsResolver and command functions to use interfaces instead of concrete types.