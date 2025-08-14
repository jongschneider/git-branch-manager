## Current Implementation Analysis of cmd/list.go

### Overview
The `cmd/list.go` file implements a comprehensive worktree listing command that shows all managed worktrees and their status. It follows the established pattern of other commands in the codebase that use interfaces for testing and modularity.

### Current Function Signatures and Methods Used

The list command currently calls these methods on the `*internal.Manager`:

1. **GetSyncStatus() (*SyncStatus, error)**
   - Returns comprehensive sync status including missing worktrees, orphaned worktrees, and branch changes
   - Used to determine sync status for display in the table

2. **GetAllWorktrees() (map[string]*WorktreeListInfo, error)**
   - Returns all existing worktrees (both tracked and ad-hoc)
   - Includes git status information for each worktree

3. **GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string**
   - Sorts worktree names with tracked worktrees first (alphabetically), then ad-hoc worktrees by creation time
   - Takes the worktree map as input and returns sorted names

4. **GetWorktreeMapping() (map[string]string, error)**
   - Returns mapping of worktree names to expected branches from gbm.branchconfig.yaml
   - Used to determine if worktrees are tracked vs untracked

### Key Data Types

```go
type WorktreeListInfo struct {
    Path           string
    ExpectedBranch string
    CurrentBranch  string
    GitStatus      *GitStatus
}

type SyncStatus struct {
    InSync             bool
    MissingWorktrees   []string
    OrphanedWorktrees  []string
    BranchChanges      map[string]BranchChange
    WorktreePromotions []WorktreePromotion
}

type BranchChange struct {
    OldBranch string
    NewBranch string
}
```

### Current Command Logic Flow

1. **Initialize Manager**: Creates manager using `createInitializedManager()`
2. **Get Sync Status**: Calls `GetSyncStatus()` to understand current state
3. **Get All Worktrees**: Calls `GetAllWorktrees()` to get comprehensive worktree information
4. **Sort Worktrees**: Uses `GetSortedWorktreeNames()` for proper display order
5. **Build Table**: Creates a table with columns: WORKTREE, BRANCH, GIT STATUS, SYNC STATUS, PATH
6. **Determine Sync Status**: For each worktree:
   - Checks for branch changes from `status.BranchChanges`
   - Checks for orphaned worktrees from `status.OrphanedWorktrees`
   - Calls `GetWorktreeMapping()` to determine if worktree is tracked
   - Formats status as IN_SYNC, OUT_OF_SYNC, or UNTRACKED

### Dependencies and Imports

```go
import (
    "errors"
    "fmt"
    "slices"
    "gbm/internal"
    "github.com/spf13/cobra"
)
```

### Existing Test Structure

The `cmd/list_test.go` file contains comprehensive integration tests:

- **Test Scenarios**: Empty repository, GBM config worktrees, untracked worktrees, orphaned worktrees, git status display, expected branch display, sorted output
- **Test Helpers**: `parseListOutput()` and `findWorktreeInRows()` for table parsing
- **Setup Functions**: Uses `testutils.NewBasicRepo()`, `testutils.NewGBMConfigRepo()`, etc.

### Current Error Handling

- Handles `internal.ErrNoRootNodesFound` gracefully 
- Returns errors for git repository issues
- Returns errors for missing gbm.branchconfig.yaml files

### Interface Extraction Requirements

Based on the analysis, the `worktreeLister` interface should include:

```go
type worktreeLister interface {
    GetSyncStatus() (*internal.SyncStatus, error)
    GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
    GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeListInfo) []string
    GetWorktreeMapping() (map[string]string, error)
}
```

### Dependencies to Consider

1. **Internal Types**: The interface will need to reference `internal.SyncStatus` and `internal.WorktreeListInfo`
2. **Error Handling**: Must maintain compatibility with `internal.ErrNoRootNodesFound`
3. **Manager Creation**: The `createInitializedManager()` function pattern should be maintained
4. **Verbose Logging**: Uses `PrintVerbose()` for debug output

### File Locations and Code Snippets

**Main Implementation** (`/Users/jschneider/code/scratch/worktree-manager/cmd/list.go`):
- Lines 21-41: Manager creation and sync status retrieval
- Lines 37-40: Worktree retrieval
- Lines 52: Sorted worktree names
- Lines 70-79: Worktree mapping for sync status determination

**Manager Methods** (`/Users/jschneider/code/scratch/worktree-manager/internal/manager.go`):
- Lines 107-169: `GetSyncStatus()` implementation
- Lines 473-527: `GetAllWorktrees()` implementation  
- Lines 834-881: `GetSortedWorktreeNames()` implementation
- Lines 380-392: `GetWorktreeMapping()` implementation

**Test Structure** (`/Users/jschneider/code/scratch/worktree-manager/cmd/list_test.go`):
- Lines 19-65: Table parsing utilities
- Lines 67-426: Comprehensive test scenarios

The interface extraction should follow the established pattern seen in other commands like `validate.go` and `sync.go`, using the `//go:generate` directive for mock generation and maintaining the same function signatures and return types.

## Established Interface Extraction Patterns

### 1. Interface Definition Patterns

**Naming Convention:**
- Interface names use camelCase with descriptive prefixes: `worktree` + `Action` (e.g., `worktreePuller`, `worktreePusher`, `worktreeValidator`)
- Interface names are always lowercase to keep them package-private

**Interface Placement:**
- Interfaces are defined directly in the command files (e.g., `/Users/jschneider/code/scratch/worktree-manager/cmd/pull.go`, `/Users/jschneider/code/scratch/worktree-manager/cmd/push.go`)
- Placed immediately after imports, before the command function
- Include clear documentation comments explaining their purpose

**Example Pattern:**
```go
//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreePuller.go . worktreePuller

// worktreePuller interface abstracts the Manager operations needed for pulling worktrees
type worktreePuller interface {
	PullAllWorktrees() error
	PullWorktree(worktreeName string) error
	IsInWorktree(currentPath string) (bool, string, error)
	GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
}
```

### 2. Mock Generation Setup

**Go Generate Directive Pattern:**
- Always placed immediately before the interface definition
- Uses the exact format: `//go:generate go run github.com/matryer/moq@latest -out ./autogen_<interfaceName>.go . <interfaceName>`
- Output file follows naming pattern: `autogen_<interfaceName>.go`

**Generated Files:**
- All generated mock files are named with `autogen_` prefix
- Examples: `autogen_worktreePuller.go`, `autogen_worktreePusher.go`, `autogen_worktreeValidator.go`
- These files provide `<interfaceName>Mock` structs with function fields and call tracking

### 3. Command Function Refactoring

**Handler Function Pattern:**
- Original command logic is extracted into separate handler functions
- Handler functions accept the interface as their first parameter
- Handler functions are named with pattern: `handle<Action><Variant>` (e.g., `handlePullAll`, `handlePullCurrent`, `handlePullNamed`)

**Example from pull.go:**
```go
func handlePullAll(puller worktreePuller) error {
	PrintInfo("Pulling all worktrees...")
	return puller.PullAllWorktrees()
}

func handlePullCurrent(puller worktreePuller, currentPath string) error {
	// Check if we're in a worktree
	inWorktree, worktreeName, err := puller.IsInWorktree(currentPath)
	if err != nil {
		return fmt.Errorf("failed to check if in worktree: %w", err)
	}
	// ... business logic
}
```

### 4. Interface Method Selection

**Common Methods Across Interfaces:**
- `GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)` - appears in most interfaces
- `IsInWorktree(currentPath string) (bool, string, error)` - used for location checking
- Action-specific methods (e.g., `PullWorktree`, `PushWorktree`, `ValidateBranches`)

**Interface Composition:**
- Interfaces include only the minimum methods needed for their specific use case
- Common utility methods are duplicated across interfaces rather than using composition
- Methods match the exact signatures from the Manager struct

### 5. Test Structure and Mock Usage

**Test Organization:**
- Unit tests use mocks and test handler functions directly
- Test files are organized with clear sections separating unit tests from integration tests
- Mock-based tests focus on business logic and error handling

**Mock Usage Pattern:**
```go
func TestHandlePullAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *worktreePullerMock
		expectErr func(t *testing.T, err error)
	}{
		{
			name: "success - pull all worktrees",
			mockSetup: func() *worktreePullerMock {
				return &worktreePullerMock{
					PullAllWorktreesFunc: func() error {
						return nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		// ... more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockSetup()
			err := handlePullAll(mock)
			tt.expectErr(t, err)
		})
	}
}
```

### 6. File Organization

**Generated Files:**
- All auto-generated mock files are in the `cmd/` directory
- Pattern: `cmd/autogen_<interfaceName>.go`
- These files should not be manually edited (contain "DO NOT EDIT" comments)

**Test Files:**
- Unit tests remain in `cmd/<command>_test.go`
- Integration tests are moved to `internal/` package when they test Manager methods directly
- Clear separation between fast unit tests and slower integration tests

### 7. Interface Documentation

**Documentation Pattern:**
- Interfaces include clear purpose statements
- Comments explain what operations the interface abstracts
- Focus on the dependency injection aspect for unit testing

**Examples:**
```go
// worktreePuller interface abstracts the Manager operations needed for pulling worktrees
// worktreeValidator abstracts validation dependencies for unit testing
// worktreeSwitcher interface abstracts the Manager operations needed for switching worktrees
```

## For the worktreeLister Interface

Based on the analysis of `/Users/jschneider/code/scratch/worktree-manager/cmd/list.go`, the `worktreeLister` interface should include these methods:

```go
//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeLister.go . worktreeLister

// worktreeLister interface abstracts the Manager operations needed for listing worktrees
type worktreeLister interface {
	GetSyncStatus() (*internal.SyncStatus, error)
	GetAllWorktrees() (map[string]*internal.WorktreeListInfo, error)
	GetSortedWorktreeNames(worktrees map[string]*internal.WorktreeListInfo) []string
	GetWorktreeMapping() (map[string]string, error)
}
```

The established patterns show a consistent, well-structured approach to interface extraction that prioritizes testability, clear separation of concerns, and maintainable code organization.