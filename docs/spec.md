# Git Branch Manager CLI Design Specification

## Overview

A command-line tool that manages Git repository branches and worktrees based on environment variables defined in a `.envrc` file. The tool synchronizes local worktrees with branch definitions and provides notifications when configurations drift out of sync.

## Core Concept

The `.envrc` file serves as the source of truth for branch management, where each environment variable maps to a specific branch:

```bash
PROD=production/2025-05-1
PREVIEW=production/2025-06-1
MAIN=master
STAGING=feature/new-api
```

## Command Structure

### Primary Commands

**`gbm init`**
- Initializes the current Git repository for branch management
- Creates initial worktree structure based on `.envrc`
- Validates that the repository is a valid Git repo
- Creates `.gbm/` directory for metadata storage

**`gbm sync`**
- Synchronizes all worktrees with current `.envrc` definitions
- Creates missing worktrees for new environment variables
- Updates existing worktrees if branch references have changed
- Removes orphaned worktrees (no longer in `.envrc`)

**`gbm status`**
- Shows current sync status between `.envrc` and actual worktrees
- Displays which branches are out of sync
- Lists missing worktrees
- Shows orphaned worktrees

**`gbm check`**
- Runs automatically on shell prompt or can be called manually
- Quick check for drift between `.envrc` and worktrees
- Returns non-zero exit code if out of sync (useful for shell integration)

### Utility Commands

**`gbm list`**
- Lists all managed worktrees and their status
- Shows environment variable mappings
- Indicates sync status for each entry


**`gbm validate`**
- Validates `.envrc` syntax and branch references
- Checks if referenced branches exist locally or remotely

## Configuration File Structure

**`.envrc`** - Primary configuration
```bash
# Long-living environment branches
PROD=production/2025-05-1
PREVIEW=production/2025-06-1
MAIN=master
STAGING=feature/new-api
```

**`.gbm/config.toml`** - Tool metadata (created automatically)
```toml
[settings]
worktree_prefix = "worktrees"
auto_fetch = true
create_missing_branches = false

[state]
last_sync = "2025-07-01T10:30:00Z"
tracked_vars = ["PROD", "PREVIEW", "MAIN", "STAGING"]
```

## Worktree Management

### Directory Structure
```
project-root/
├── .git/
├── .envrc
├── .gbm/
│   └── config.toml
├── worktrees/
│   ├── PROD/           # Contains production/2025-05-1 branch
│   ├── PREVIEW/        # Contains production/2025-06-1 branch
│   ├── MAIN/           # Contains master branch
│   └── STAGING/        # Contains feature/new-api branch
└── main-workspace/     # Original repository workspace
```

### Worktree Naming
- Worktrees are named after the environment variable (e.g., `PROD`, `PREVIEW`)
- Stored in `worktrees/` subdirectory by default
- Configurable via `.gbm/config.toml`

## Sync Detection and Notification

### Automatic Checking
Integration with shell prompt to show sync status:
```bash
# Example shell integration
export PS1="$(gbm check --format=prompt)$PS1"
```

### Sync Status Indicators
- ✅ All worktrees in sync
- ⚠️  Some worktrees out of sync
- ❌ Major configuration drift detected
- 🔄 Sync in progress

### Interactive Sync Resolution
When drift is detected, prompt user with options:
```
⚠️  Configuration drift detected:

Changes needed:
  • PROD: production/2025-05-1 → production/2025-07-1 (branch changed)
  • STAGING: worktree missing (new environment variable)
  • OLD_FEATURE: orphaned worktree (variable removed)

Actions:
  [s] Sync all changes
  [r] Review changes individually
  [i] Ignore for this session
  [q] Quit
```

## Command-Line Interface

### Flags and Options

**Global Flags:**
- `--config, -c`: Specify custom `.envrc` path
- `--worktree-dir, -w`: Override worktree directory location
- `--verbose, -v`: Enable verbose output
- `--quiet, -q`: Suppress non-error output

**Command-Specific Options:**

`gbm init`:
- `--force`: Initialize even if worktrees already exist
- `--fetch`: Fetch remote branches during initialization

`gbm sync`:
- `--dry-run`: Show what would be changed without making changes
- `--force`: Skip confirmation prompts
- `--fetch`: Update remote tracking before sync

`gbm check`:
- `--format=prompt|json|text`: Output format for integration
- `--exit-code`: Return status code only

## Error Handling

### Common Error Scenarios
1. **Not a Git repository**: Clear error message with suggestion to run `git init`
2. **Missing `.envrc`**: Prompt to create one or specify path
3. **Invalid branch references**: List invalid branches and suggest alternatives
4. **Permission issues**: Clear guidance on directory permissions
5. **Conflicting worktrees**: Handle existing worktree conflicts gracefully

### Recovery Mechanisms
- Backup worktree state before major operations
- Rollback capability for failed sync operations
- Repair mode for corrupted `.gbm/` metadata

## Integration Points

### Shell Integration
```bash
# Add to .bashrc/.zshrc for automatic checking
eval "$(gbm shell-integration)"
```

### Git Hooks Integration
- Optional pre-commit hook to validate `.envrc` changes
- Post-merge hook to check for sync requirements

### CI/CD Integration
- `gbm validate` can be used in CI pipelines
- JSON output format for programmatic consumption

## Implementation Considerations

### Dependencies
- Standard library only where possible
- Git command-line tool (shell out to `git` commands)
- TOML parsing library for configuration

### Performance
- Lazy loading of Git repository state
- Caching of branch existence checks
- Minimal filesystem operations during status checks

### Cross-Platform Support
- Windows, macOS, and Linux compatibility
- Handle path separators and permissions appropriately
- Shell integration for major shells (bash, zsh, fish)

## Future Enhancements

- **Remote branch creation**: Option to create missing branches remotely
- **Template support**: Branch naming templates and conventions
- **Multi-repository support**: Manage multiple repositories from single config
- **Plugin system**: Custom sync behaviors and integrations
- **Web interface**: Optional local web UI for visual management
