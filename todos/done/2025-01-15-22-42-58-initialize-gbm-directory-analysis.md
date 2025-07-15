Based on my analysis of the `gbm clone` command implementation and the codebase, I can now provide a comprehensive summary:

## Current State Analysis

### 1. Current Clone Command Functionality

The `gbm clone` command in `/Users/jschneider/code/scratch/worktree-manager/cmd/clone.go` currently performs these steps:

1. **Clones repository as bare repo** - Creates a bare clone in `.git` directory
2. **Discovers default branch** - Determines the main branch from remote
3. **Creates main worktree** - Creates the primary worktree in `worktrees/` directory  
4. **Sets up gbm.branchconfig.yaml** - Creates or copies the YAML configuration file
5. **Initializes worktree management** - Calls `initializeWorktreeManagement()` which creates a Manager

### 2. Where .gbm Directory Initialization Should Be Added

The `.gbm` directory initialization should be added in the `initializeWorktreeManagement()` function at **line 255-285** in `/Users/jschneider/code/scratch/worktree-manager/cmd/clone.go`. This function already:

- Creates a `Manager` instance via `internal.NewManager(wd)` (line 263)
- The `NewManager` function loads config and state from the `.gbm` directory (line 44-72 in `/Users/jschneider/code/scratch/worktree-manager/internal/manager.go`)

### 3. What Default Config.toml Should Contain

Based on the `DefaultConfig()` function in `/Users/jschneider/code/scratch/worktree-manager/internal/config.go` (lines 80-118), the default `config.toml` should contain:

```toml
[settings]
worktree_prefix = "worktrees"
auto_fetch = true
create_missing_branches = false
merge_back_alerts = false
hotfix_prefix = "HOTFIX"
mergeback_prefix = "MERGE" 
merge_back_check_interval = "6h0m0s"
merge_back_user_commit_interval = "30m0s"

[icons]
success = "✅"
warning = "⚠️"
error = "❌"
info = "💡"
orphaned = "🗑️"
dry_run = "🔍"
missing = "📁"
changes = "🔄"
git_clean = "✓"
git_dirty = "~"
git_ahead = "↑"
git_behind = "↓"
git_diverged = "⇕"
git_unknown = "?"

[jira]
me = ""

[[file_copy.rules]]
# Empty rules array by default
```

### 4. What Default State.toml Should Contain

Based on the `DefaultState()` function in `/Users/jschneider/code/scratch/worktree-manager/internal/state.go` (lines 32-42), the default `state.toml` should contain:

```toml
last_sync = "0001-01-01T00:00:00Z"
tracked_vars = []
ad_hoc_worktrees = []
current_worktree = ""
previous_worktree = ""
last_mergeback_check = "0001-01-01T00:00:00Z"
```

### 5. How Other Commands Handle .gbm Directory Creation

Looking at the codebase patterns:

- **Config and State loading**: The `LoadConfig()` and `LoadState()` functions return default instances if files don't exist, but don't create the files
- **Saving**: Both `Config.Save()` and `State.Save()` methods create the `.gbm` directory via `os.MkdirAll(gbmDir, 0o755)` before writing files
- **Test patterns**: Tests create the `.gbm` directory manually with `os.MkdirAll(".gbm", 0o755)` and then save default config/state

### 6. Existing Patterns for Default Configuration Creation

The current pattern is:
1. `LoadConfig(gbmDir)` and `LoadState(gbmDir)` return defaults when files don't exist
2. The `Manager` gets created with these defaults
3. Files are only written when `SaveConfig()` or `SaveState()` is called

**The missing piece**: The clone command should explicitly save the default config and state files to disk during initialization, rather than just keeping them in memory.

### Recommended Implementation

Add this code to the `initializeWorktreeManagement()` function after creating the manager:

```go
// Save default configuration files to .gbm directory
if err := manager.SaveConfig(); err != nil {
    return fmt.Errorf("failed to save default config: %w", err)
}

if err := manager.SaveState(); err != nil {
    return fmt.Errorf("failed to save default state: %w", err)
}
```

This would ensure that the `.gbm` directory is created with proper default `config.toml` and `state.toml` files during the clone operation.

Based on my examination of the config and state management code, here's what I understand about the structure and default values:

## Configuration Structure (`config.toml`)

### Default Configuration Values:
1. **Settings section**:
   - `worktree_prefix`: "worktrees" (default directory name)
   - `auto_fetch`: true
   - `create_missing_branches`: false
   - `merge_back_alerts`: false
   - `hotfix_prefix`: "HOTFIX"
   - `mergeback_prefix`: "MERGE"
   - `merge_back_check_interval`: 6 hours
   - `merge_back_user_commit_interval`: 30 minutes

2. **Icons section**: Various Unicode symbols for status indicators and Git states
3. **Jira section**: Empty `me` field
4. **FileCopy section**: Empty rules array

### How config.toml is created:
- The `LoadConfig()` function checks if `config.toml` exists in the `.gbm` directory
- If it doesn't exist, it returns `DefaultConfig()` with default values
- The `Save()` method creates the `.gbm` directory and writes the config as TOML

## State Structure (`state.toml`)

### Default State Values:
- `last_sync`: Zero time value
- `tracked_vars`: Empty array
- `ad_hoc_worktrees`: Empty array
- `current_worktree`: Empty string
- `previous_worktree`: Empty string
- `last_mergeback_check`: Zero time value

### How state.toml is created:
- The `LoadState()` function checks if `state.toml` exists in the `.gbm` directory
- If it doesn't exist, it returns `DefaultState()` with default values
- The `Save()` method creates the `.gbm` directory and writes the state as TOML

## Key Observations:

1. **File Locations**: Both files are stored in a `.gbm` directory at the repository root
2. **TOML Format**: Both files use TOML encoding/decoding
3. **Default Behavior**: If files don't exist, the code returns default structs rather than creating files immediately
4. **Auto-creation**: The `Save()` methods will create the `.gbm` directory if it doesn't exist
5. **Migration Support**: The code includes a `ConfigState` struct for migration purposes
6. **Additional Config**: There's also support for a YAML-based `gbm.branchconfig.yaml` file for worktree definitions

The initialization logic is designed to work with missing files by providing sensible defaults, and files are only created when explicitly saved through the `Save()` methods.

Based on my analysis of the codebase, I can now provide a comprehensive overview of the existing .gbm directory creation and management patterns:

## Summary of .gbm Directory Creation and Management Patterns

### 1. **Current .gbm Directory Creation Pattern**

The .gbm directory is created **lazily** when configuration or state files are saved, not during initialization. This happens in two key places:

- **`internal/config.go:140-143`**: The `Config.Save()` method creates the .gbm directory with `os.MkdirAll(gbmDir, 0o755)` before writing config.toml
- **`internal/state.go:62-65`**: The `State.Save()` method creates the .gbm directory with `os.MkdirAll(gbmDir, 0o755)` before writing state.toml

### 2. **Directory Path Determination**

The .gbm directory path is determined by:
- **`internal/config.go:135-138`**: `GetGBMDir(repoRoot)` returns `filepath.Join(repoRoot, ".gbm")`
- **`internal/manager.go:45`**: `gbmDir := filepath.Join(repoPath, ".gbm")`

### 3. **Configuration and State Loading Pattern**

Both configuration and state files follow a consistent pattern:

**Configuration (`LoadConfig` in `internal/config.go:120-133`)**:
- Checks if `config.toml` exists
- If not found, returns `DefaultConfig()` (doesn't create files)
- If found, loads and parses the TOML file

**State (`LoadState` in `internal/state.go:44-59`)**:
- Checks if `state.toml` exists  
- If not found, returns `DefaultState()` (doesn't create files)
- If found, loads and parses the TOML file

### 4. **Manager Initialization Pattern**

The `NewManager` function (`internal/manager.go:44-72`):
1. Constructs gbmDir path
2. Calls `LoadConfig(gbmDir)` - returns defaults if files don't exist
3. Calls `LoadState(gbmDir)` - returns defaults if files don't exist
4. Creates GitManager and IconManager
5. **Does NOT create .gbm directory or files at this point**

### 5. **Error Handling Pattern**

Both `Save` methods use consistent error handling:
```go
if err := os.MkdirAll(gbmDir, 0o755); err != nil {
    return fmt.Errorf("failed to create .gbm directory: %w", err)
}
```

### 6. **Current Clone Command Behavior**

The `gbm clone` command (`cmd/clone.go:255-284`):
1. Creates manager with `NewManager(wd)` - loads defaults, doesn't create .gbm
2. Calls `manager.LoadGBMConfig()` - loads branch config
3. Calls `manager.Sync()` - may trigger state saves that create .gbm directory

### 7. **Testing Patterns**

Tests manually create .gbm directories and files:
- **`cmd/root_test.go:33-34`**: `os.MkdirAll(gbmDir, 0755)` followed by `config.Save(gbmDir)`
- **`cmd/mergeback_test.go:520`**: `os.MkdirAll(".gbm", 0o755)` followed by file creation

### 8. **Gap Identified**

The current `gbm clone` command does **NOT** explicitly initialize the .gbm directory with default config.toml and state.toml files. The directory and files are only created when:
- Configuration is saved via `SaveConfig()`
- State is saved via `SaveState()`
- This happens indirectly through operations like `Sync()`

### 9. **Recommended Implementation Pattern**

To implement the requested feature, the pattern should be:

1. **Add initialization call in `initializeWorktreeManagement()`**:
   ```go
   // Save default config and state to create .gbm directory
   if err := manager.SaveConfig(); err != nil {
       return fmt.Errorf("failed to initialize .gbm/config.toml: %w", err)
   }
   if err := manager.SaveState(); err != nil {
       return fmt.Errorf("failed to initialize .gbm/state.toml: %w", err)
   }
   ```

2. **This leverages existing patterns**:
   - Uses existing `SaveConfig()` and `SaveState()` methods
   - Maintains consistent error handling
   - Follows the established directory creation pattern
   - Creates both config.toml and state.toml with defaults

3. **The existing `Save` methods already handle**:
   - Directory creation (`MkdirAll`)
   - File creation with proper permissions
   - TOML encoding with proper error handling

This approach ensures the .gbm directory is initialized with default files immediately after cloning, following the established patterns in the codebase.