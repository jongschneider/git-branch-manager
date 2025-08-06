Perfect! Now I have a complete understanding of the switch command and its Manager dependencies. Let me provide you with a structured analysis of the findings.

## Analysis of cmd/switch.go

### Current Functions in the Switch Command

1. **`newSwitchCommand()`** - Creates and configures the cobra command
2. **`switchToWorktreeWithFlag(manager, worktreeName, printPath)`** - Core switching logic with flag handling
3. **`findFuzzyMatch(manager, target)`** - Implements fuzzy matching for worktree names
4. **`listWorktrees(manager)`** - Lists all available worktrees when no args provided
5. **`getWorktreeNames()`** - Provides worktree names for shell completion

### Manager Methods Called from Switch Command

Based on the analysis of `/Users/jschneider/code/scratch/worktree-manager/cmd/switch.go`, here are all the Manager methods that would need to be included in a `worktreeSwitcher` interface:

#### Core Switch Methods:
1. **`GetPreviousWorktree() string`** - Returns the previously accessed worktree name
2. **`GetWorktreePath(worktreeName string) (string, error)`** - Gets the filesystem path for a worktree
3. **`SetCurrentWorktree(worktreeName string) error`** - Updates current worktree tracking

#### Listing and Discovery Methods:
4. **`GetAllWorktrees() (map[string]*WorktreeListInfo, error)`** - Gets all available worktrees
5. **`GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string`** - Returns sorted worktree names
6. **`GetStatusIcon(gitStatus *GitStatus) string`** - Gets status icon for git status display

### Function Signatures and Error Handling Patterns

#### Error Handling Patterns:
- **Graceful degradation**: When tracking fails, logs warning but continues operation
- **Fallback behavior**: Exact match → fuzzy match → error
- **User-friendly errors**: Clear messages like "no previous worktree available"
- **Verbose logging**: Uses `PrintVerbose()` for debugging information

#### User Interaction Patterns:
- **Multiple output modes**: Regular display vs. print-path for shell integration
- **Shell integration support**: Checks `GBM_SHELL_INTEGRATION` environment variable
- **Fuzzy matching**: Automatic case-insensitive substring matching with preference for prefix matches
- **Auto-completion**: Provides worktree names for shell completion

### Current Structure Analysis

```go
// Main command function
func newSwitchCommand() *cobra.Command {
    // Command setup with flags and completion
    cmd.Flags().Bool("print-path", false, "Print the worktree path only")
    cmd.ValidArgsFunction = func(...) ([]string, cobra.ShellCompDirective) {
        return getWorktreeNames(), cobra.ShellCompDirectiveNoFileComp
    }
}

// Core switching logic
func switchToWorktreeWithFlag(manager *internal.Manager, worktreeName string, printPath bool) error {
    // 1. Try exact match: manager.GetWorktreePath(worktreeName)
    // 2. Try fuzzy match: findFuzzyMatch(manager, worktreeName)
    // 3. Track switch: manager.SetCurrentWorktree(worktreeName)
    // 4. Handle output format based on flags and environment
}

// Fuzzy matching helper
func findFuzzyMatch(manager *internal.Manager, target string) string {
    // Uses: manager.GetAllWorktrees()
    // Implements case-insensitive substring matching with prefix preference
}

// Listing functionality
func listWorktrees(manager *internal.Manager) error {
    // Uses: manager.GetAllWorktrees(), manager.GetSortedWorktreeNames(), manager.GetStatusIcon()
    // Displays formatted worktree list with status icons and paths
}
```

### Exact Manager Methods for worktreeSwitcher Interface

Based on the analysis, the `worktreeSwitcher` interface should include exactly these methods:

```go
type worktreeSwitcher interface {
    GetPreviousWorktree() string
    GetWorktreePath(worktreeName string) (string, error)
    SetCurrentWorktree(worktreeName string) error
    GetAllWorktrees() (map[string]*WorktreeListInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string
    GetStatusIcon(gitStatus *GitStatus) string
}
```

The switch command is well-structured with clear separation between the command setup, core switching logic, fuzzy matching, and listing functionality. It follows good patterns for error handling and user interaction, with robust support for shell integration and multiple output formats.

Based on my comprehensive analysis of the existing interface extraction patterns in this codebase, here are the established patterns that should be followed for the `worktreeSwitcher` interface:

## Interface Extraction Patterns Analysis

### 1. Interface Definition Patterns

**Location**: Interfaces are defined at the top of the command file, immediately after imports and before the command functions.

**Naming Convention**: 
- Interface names follow the pattern `worktree{Action}` (e.g., `worktreePuller`, `worktreePusher`, `worktreeRemover`)
- Use camelCase starting with lowercase
- Interface name should match the command action

**Interface Structure**:
```go
// worktree{Action} interface abstracts the Manager operations needed for {action} worktrees
type worktree{Action} interface {
    // Only include methods actually used by the command handlers
    Method1() ReturnType
    Method2(param ParamType) (ReturnType, error)
}
```

### 2. Mock Generation Comments

**Two patterns observed**:

**Pattern 1** (used by `pull.go`):
```go
//go:generate go tool moq -out ./autogen_worktreePuller.go . worktreePuller
```

**Pattern 2** (used by `push.go` and `remove.go`):
```go
//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreePusher.go . worktreePusher
```

**Recommendation**: Use Pattern 2 (with `go run github.com/matryer/moq@latest`) for consistency with the newer implementations.

### 3. Interface Method Selection

**Principle**: Only include methods that are actually used by the command handlers.

**Examples**:
- `worktreePuller`: `PullAllWorktrees()`, `PullWorktree()`, `IsInWorktree()`, `GetAllWorktrees()`
- `worktreePusher`: Same methods as puller but with Push instead of Pull
- `worktreeRemover`: `GetWorktreePath()`, `GetWorktreeStatus()`, `RemoveWorktree()`, `GetAllWorktrees()`

### 4. Command Function Refactoring Pattern

**Handler Functions**: Extract command logic into separate functions that accept the interface:
```go
func handle{Action}All({action} worktree{Action}) error { ... }
func handle{Action}Current({action} worktree{Action}, currentPath string) error { ... }
func handle{Action}Named({action} worktree{Action}, worktreeName string) error { ... }
```

**cobra.Command RunE**: Use manager as interface in command execution:
```go
RunE: func(cmd *cobra.Command, args []string) error {
    manager, err := createInitializedManager()
    if err != nil { ... }
    
    return handle{Action}(manager, args...)
}
```

### 5. Test Structure Patterns

**File Organization**:
- Fast unit tests with mocks in `cmd/{command}_test.go`
- Integration tests using real git in `internal/{command}_test.go`

**Test Structure**:
```go
func TestHandle{Action}(t *testing.T) {
    tests := []struct {
        name      string
        mockSetup func() *worktree{Action}Mock
        expectErr func(t *testing.T, err error)
    }{
        {
            name: "success case",
            mockSetup: func() *worktree{Action}Mock {
                return &worktree{Action}Mock{
                    MethodFunc: func(...) (...) {
                        // Mock implementation
                        return nil
                    },
                }
            },
            expectErr: func(t *testing.T, err error) {
                assert.NoError(t, err)
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := tt.mockSetup()
            err := handle{Action}(mock, ...)
            tt.expectErr(t, err)
        })
    }
}
```

**Mock Verification**:
- Use `assert.Len(t, mock.MethodCalls(), expectedCount)` to verify method calls
- Use assertions within mock functions to verify parameters

### 6. For `worktreeSwitcher` Interface Specifically

Based on the `switch.go` analysis, the interface should include:

```go
//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeSwitcher.go . worktreeSwitcher

// worktreeSwitcher interface abstracts the Manager operations needed for switching worktrees
type worktreeSwitcher interface {
    GetWorktreePath(worktreeName string) (string, error)
    SetCurrentWorktree(worktreeName string) error
    GetPreviousWorktree() string
    GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeListInfo) []string
    GetStatusIcon(gitStatus *internal.GitStatus) string
}
```

**Handler Functions to Extract**:
- `switchToWorktreeWithFlag()` → `handleSwitchToWorktree()`
- `findFuzzyMatch()` → Keep as utility function using interface
- `listWorktrees()` → `handleListWorktrees()`

### 7. File Generation

The mock generation creates files with this pattern:
- Filename: `autogen_worktree{Action}.go`
- Generated mock struct: `worktree{Action}Mock`
- Comprehensive call tracking and verification methods
- Thread-safe implementation with mutexes

This established pattern ensures consistency, testability, and maintainability across all command interfaces in the codebase.