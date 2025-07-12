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

### Ad-hoc Worktree Management

- `gbm add <worktree-name> [branch-name]` - Add a new worktree
  - `gbm add feature-work existing-branch` - Create worktree on existing branch
  - `gbm add feature-work new-branch -b` - Create worktree with new branch
  - `gbm add feature-work --interactive` - Interactive branch selection

### JIRA Integration

The `add` command includes intelligent JIRA integration when the `jira` CLI is available:

#### Tab Completion
- **First tab**: Shows JIRA keys with issue summaries for context
- **Second tab**: Suggests branch names based on JIRA issue details

```bash
$ gbm add <TAB>
INGSVC-5544    Zoom Clips - New Integration
INGSVC-5581    MsSharepoint - Authentication Error
INGSVC-3663    MONDAY.COM: handle error response

$ gbm add INGSVC-5544 <TAB>
feature/INGSVC-5544_Zoom_Clips_New_Integration
```

#### Smart Branch Generation
- **Stories** and **Improvements** → `feature/` prefix
- **Bugs** → `bug/` prefix
- Summary text cleaned and formatted for branch names
- Special characters replaced with underscores

#### Workflow Examples
```bash
# Two-tab completion workflow
$ gbm add INGSV<TAB>                    # Complete to INGSVC-5544
$ gbm add INGSVC-5544 <TAB>             # Complete to branch name
$ gbm add INGSVC-5544 feature/INGSVC-5544_Zoom_Clips_New_Integration -b

# Smart suggestion workflow
$ gbm add INGSVC-5544
Error: branch name required. Suggested: feature/INGSVC-5544_Zoom_Clips_New_Integration

Try: gbm add INGSVC-5544 feature/INGSVC-5544_Zoom_Clips_New_Integration -b
```

#### Requirements
- `jira` CLI tool installed and authenticated
- Gracefully falls back to basic completion when JIRA unavailable

### Shell Integration

#### Status Checking
Add automatic checking to your shell:

```bash
# Add to .bashrc/.zshrc
eval "$(gbm shell-integration)"
```

#### Tab Completion
Enable tab completion for enhanced JIRA integration:

```bash
# Bash
gbm completion bash > /etc/bash_completion.d/gbm
# Or for current session: source <(gbm completion bash)

# Zsh
gbm completion zsh > "${fpath[1]}/_gbm"

# Fish
gbm completion fish > ~/.config/fish/completions/gbm.fish
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

[file_copy]
[[file_copy.rules]]
source_worktree = "master"
files = [".env", "config/local.json", "scripts/"]
```

### File Copying for Ad-Hoc Worktrees

You can configure automatic file copying when creating new **ad-hoc worktrees** (created with `gbm add`, not tracked in `.gbm.config.yaml`) by adding `[file_copy]` rules to your `.gbm/config.toml`:

```toml
[file_copy]
# Copy .env file from the master worktree to all new ad-hoc worktrees
[[file_copy.rules]]
source_worktree = "master"
files = [".env"]

# Copy multiple files and directories from a specific worktree
[[file_copy.rules]]
source_worktree = "main"
files = [".env.local", "config/development.json", "scripts/"]
```

**Configuration Options:**
- `source_worktree`: The name of the worktree to copy files from
- `files`: Array of file paths or directory paths to copy (supports both files and directories)

**File Copy Rules:**
- Files are copied **only** when creating new ad-hoc worktrees with `gbm add`
- Tracked worktrees (defined in `.gbm.config.yaml`) do **not** get file copying
- If the source worktree doesn't exist, the rule is skipped with a warning
- If a specific file doesn't exist in the source worktree, that file is skipped with a warning
- If a file already exists in the target worktree, it is skipped
- Directory permissions and file permissions are preserved during copying
- Files are copied recursively for directories

**Use Cases:**
- Copy environment files (`.env`, `.env.local`) to new feature branches
- Copy local configuration files that aren't tracked in git
- Copy development scripts or tools to new worktrees

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
