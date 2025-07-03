# Git Branch Manager (gbm)

A command-line tool that manages Git repository branches and worktrees based on environment variables defined in a `.envrc` file. The tool synchronizes local worktrees with branch definitions and provides notifications when configurations drift out of sync.

## Installation

```bash
go build -o gbm
```

## Quick Start

1. Create a `.envrc` file in your Git repository:
```bash
PROD=production/2025-05-1
PREVIEW=production/2025-06-1  
MAIN=master
STAGING=feature/new-api
```

2. Initialize gbm:
```bash
./gbm init
```

3. Check status:
```bash
./gbm status
```

4. Sync worktrees:
```bash
./gbm sync
```

## Commands

### Core Commands

- `gbm init` - Initialize the repository for branch management
- `gbm sync` - Synchronize worktrees with .envrc definitions
- `gbm status` - Show current sync status
- `gbm check` - Quick check for drift (useful for automation)

### Utility Commands

- `gbm list` - List all managed worktrees
- `gbm clean` - Remove orphaned worktrees
- `gbm validate` - Validate .envrc syntax and branch references

### Shell Integration

Add automatic checking to your shell:

```bash
# Add to .bashrc/.zshrc
eval "$(gbm shell-integration)"
```

## Configuration

The tool creates a `.gbm/config.toml` file for metadata storage with default settings:

```toml
[settings]
worktree_prefix = "worktrees"
auto_fetch = true
create_missing_branches = false

[state]
last_sync = "2025-07-01T10:30:00Z"
tracked_vars = ["PROD", "PREVIEW", "MAIN", "STAGING"]
```

## Directory Structure

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

## Global Flags

- `--config, -c`: Specify custom .envrc path
- `--worktree-dir, -w`: Override worktree directory location
- `--verbose, -v`: Enable verbose output
- `--quiet, -q`: Suppress non-error output