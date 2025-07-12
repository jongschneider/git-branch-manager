# Git Branch Manager (gbm)

A command-line tool that manages Git repository branches and worktrees based on configurations defined in a `.gbm.config.yaml` file. The tool provides automated worktree management, JIRA integration, and intelligent file copying between worktrees.

## Installation

### Git Branch Manager

```bash
go build -o gbm
```

### JIRA CLI (Optional)

For enhanced JIRA integration features, install the official JIRA CLI from [ankitpokhrel/jira-cli](https://github.com/ankitpokhrel/jira-cli).

**Note:** gbm will automatically detect if the JIRA CLI is available and fall back gracefully if not installed. All core worktree functionality works without JIRA integration.

## Quick Start

1. Clone a repository with automatic worktree setup:
```bash
gbm clone <repository-url>
```

2. Or configure an existing repository by creating `.gbm.config.yaml`:
```yaml
worktrees:
  main:
    branch: main
    description: "Main production branch"
  staging:
    branch: feature/staging-env
    description: "Staging environment"
```

3. Sync worktrees with configuration:
```bash
gbm sync
```

4. List all worktrees:
```bash
gbm list
```

## Commands

### Core Worktree Management

- `gbm add <worktree-name> [branch-name]` - Add a new worktree
  - `gbm add feature-work existing-branch` - Create worktree on existing branch
  - `gbm add feature-work new-branch -b` - Create worktree with new branch
  - `gbm add feature-work --interactive` - Interactive branch selection

- `gbm list` - List all managed worktrees with sync status
- `gbm sync` - Synchronize worktrees with `.gbm.config.yaml` definitions
- `gbm remove <worktree-name>` - Remove worktrees with safety checks
- `gbm switch [worktree-name]` - Switch between worktrees with fuzzy matching

### Repository Operations

- `gbm clone <repository-url>` - Clone repository as bare repo with worktree setup
- `gbm pull [worktree-name]` - Pull changes from remote (current/named/all worktrees)
- `gbm push [worktree-name]` - Push changes to remote (current/named/all worktrees)
- `gbm info <worktree-name>` - Display detailed worktree information

### Validation and Utilities

- `gbm validate` - Validate `.gbm.config.yaml` syntax and branch references

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

#### Directory Switching
Add automatic directory switching and status checking:

```bash
# Add to .bashrc/.zshrc
eval "$(gbm shell-integration)"
```

This provides functions like `gcd <worktree-name>` for quick navigation.

## Configuration

### Primary Configuration: `.gbm.config.yaml`

Define tracked worktrees and their branch mappings:

```yaml
worktrees:
  main:
    branch: main
    description: "Main production branch"
    merge_into: ""  # Optional merge strategy
  staging:
    branch: feature/staging-env
    description: "Staging environment branch"
  prod:
    branch: production/2025-07-1
    description: "Production release branch"
```

### Tool Configuration: `.gbm/config.toml`

The tool creates a `.gbm/config.toml` file for settings and metadata:

```toml
[settings]
worktree_prefix = "worktrees"
auto_fetch = true
create_missing_branches = false
merge_back_alerts = false

[jira]
me = "cached-username"

[file_copy]
[[file_copy.rules]]
source_worktree = "main"
files = [".env", "config/local.json", "scripts/"]
```

### File Copying for Ad-Hoc Worktrees

Configure automatic file copying when creating new **ad-hoc worktrees** (created with `gbm add`, not tracked in `.gbm.config.yaml`):

```toml
[file_copy]
# Copy .env file from the main worktree to all new ad-hoc worktrees
[[file_copy.rules]]
source_worktree = "main"
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
├── .git/           # Bare repository (after gbm clone)
├── .gbm.config.yaml
├── .gbm/
│   └── config.toml
├── worktrees/
│   ├── main/           # Contains main branch
│   ├── staging/        # Contains feature/staging-env branch
│   ├── prod/           # Contains production/2025-07-1 branch
│   └── feature-work/   # Ad-hoc worktree created with gbm add
```

## Global Flags

- `--config, -c`: Specify custom `.gbm.config.yaml` path
- `--worktree-dir, -w`: Override worktree directory location
- `--verbose, -v`: Enable verbose output
- `--quiet, -q`: Suppress non-error output

## Features

### Sync Status Tracking
- Automatic detection of configuration drift
- Visual indicators for worktree sync status
- Orphaned worktree detection

### Git Integration
- Bare repository support for efficient storage
- Branch validation and existence checking
- Remote tracking and push/pull operations

### JIRA Workflow
- Intelligent ticket detection and branch naming
- Cached user information for performance
- Seamless integration with existing JIRA workflows

### File Management
- Configurable file copying between worktrees
- Permission preservation
- Conflict handling and safety checks