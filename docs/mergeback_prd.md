# Product Requirements Document: Merge-back Detection and Alerting System

## Executive Summary

This PRD outlines the implementation of an automated merge-back detection and alerting system for the Git Branch Manager (gbm) CLI tool. The system will proactively monitor tracked branches for incomplete merge-back workflows and alert users when they are responsible for completing the merge-back chain.

## Background

### Current State
- gbm manages git worktrees based on environment variable mappings in `.envrc` files
- Long-lived branches represent different environments (MAIN=dev, PREVIEW=staging, PROD=production)
- Manual merge-back process: PROD → PREVIEW → MAIN when hotfixes are applied
- Engineers frequently forget to complete the merge-back chain, leading to drift between environments

### Problem Statement
When hotfixes are merged into PROD, they must be cascaded back through PREVIEW to MAIN to keep environments in sync. Currently:
- No automated detection of incomplete merge-backs
- Engineers are unaware when they're responsible for merge-backs
- Manual tracking is error-prone and time-consuming
- Environment drift goes unnoticed until deployment issues occur

## Goals & Success Criteria

### Primary Goals
1. **Automated Detection**: Detect when commits exist in lower-hierarchy branches that haven't been merged to higher branches
2. **User Responsibility**: Identify when the current user authored commits requiring merge-back
3. **Transparent Alerting**: Surface alerts during normal gbm usage without requiring new commands
4. **Clear Guidance**: Provide actionable information about what needs to be merged back

### Success Criteria
- 100% detection rate for merge-back requirements within tracked branches
- Alerts display within 500ms of gbm command execution
- Zero false positives for branches that are properly synchronized
- Clear identification of user-authored commits requiring action

## User Stories

### Primary User: Software Engineer

**Story 1: Hotfix Author Alert**
```
As a software engineer who created a hotfix in PROD,
I want to be alerted when I run gbm commands that my commits need to be merged back,
So that I can complete the merge-back chain before other work continues.
```

**Story 2: General Awareness**
```
As a software engineer working in the repository,
I want to see when any merge-backs are needed in tracked branches,
So that I'm aware of the overall state of environment synchronization.
```

**Story 3: Clear Responsibility**
```
As a software engineer,
I want to clearly see which commits are mine vs. others requiring merge-back,
So that I can focus on my responsibilities while being aware of team needs.
```

## Functional Requirements

### FR1: Branch Hierarchy Detection
- **Description**: System must determine merge-back hierarchy from `.envrc` file order
- **Acceptance Criteria**:
  - Parse `.envrc` file preserving line order
  - Ignore comments and empty lines
  - Build hierarchy chain where bottom entries merge into entries above them
  - Handle dynamic hierarchy (branches may be added/removed)

**Example**:
```bash
# .envrc file order determines hierarchy
MAIN=main           # Top of chain
PREVIEW=preview-v2  # Merges into MAIN
PROD=prod-v1.0      # Merges into PREVIEW
```
Merge flow: PROD → PREVIEW → MAIN

### FR2: Commit Analysis
- **Description**: Identify commits that exist in source branches but not in target branches
- **Acceptance Criteria**:
  - Use `git log TARGET..SOURCE` to find missing commits
  - Extract commit hash, message, author, email, and timestamp
  - Handle empty results gracefully
  - Skip non-existent branches without failing

### FR3: User Identification
- **Description**: Determine which commits were authored by the current user
- **Acceptance Criteria**:
  - Primary: Match against `git config user.email`
  - Fallback: Match against `git config user.name` if email unavailable
  - Handle missing git config gracefully
  - Case-sensitive exact matching

### FR4: Alert Display
- **Description**: Display merge-back alerts before command output
- **Acceptance Criteria**:
  - Trigger on every gbm command execution
  - Display before primary command output
  - Show warning icon and clear messaging
  - Include commit details for user-authored commits
  - Format output for terminal readability

**Alert Format**:
```
⚠️  Merge-back required in tracked branches:

PROD → PREVIEW: 2 commits need merge-back (1 by you)
• abc1234 - Fix critical auth bug (you - 2 days ago)

PREVIEW → MAIN: 5 commits need merge-back (0 by you)

```

### FR5: Error Handling
- **Description**: Handle error conditions gracefully without breaking gbm functionality
- **Acceptance Criteria**:
  - Silent skip when not in git repository
  - Silent skip when `.envrc` doesn't exist
  - Brief warning for git command failures
  - Brief warning for branch access issues
  - Never prevent primary gbm command from executing

## Technical Requirements

### TR1: Performance
- **Target**: Merge-back check completes within 500ms for typical repositories
- **Constraint**: Acceptable to run on every command (V1 - no optimization required)

### TR2: Integration
- **Hook Point**: `PersistentPreRun` function in root command
- **Dependencies**: Existing internal package structure
- **Compatibility**: Must not break existing gbm functionality

### TR3: Git Operations
- **Commands Used**:
  - `git log TARGET..SOURCE --format=%H|%s|%an|%ae|%ct`
  - `git config user.email`
  - `git config user.name`
- **Error Handling**: Graceful degradation on git command failures

### TR4: File Processing
- **Input**: `.envrc` file in repository root
- **Parsing**: Line-by-line with environment variable extraction
- **Format**: `ENVVAR=branch-name` pairs

## Implementation Details

### Architecture

```
┌──────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Root Command    │───▶│  Merge-back      │───▶│  Alert Display  │
│  PersistentPreRun│    │  Detection Logic │    │  (stderr)       │
└──────────────── ─┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  Git Operations  │
                       │  (log, config)   │
                       └──────────────────┘
```

### Data Structures

```go
type MergeBackStatus struct {
    MergeBacksNeeded []MergeBackInfo
    HasUserCommits   bool
}

type MergeBackInfo struct {
    FromBranch   string        // Environment variable name (e.g., "PROD")
    ToBranch     string        // Environment variable name (e.g., "PREVIEW")
    Commits      []CommitInfo  // All commits needing merge-back
    UserCommits  []CommitInfo  // Commits authored by current user
    TotalCount   int
    UserCount    int
}

type CommitInfo struct {
    Hash      string
    Message   string
    Author    string
    Timestamp time.Time
    IsUser    bool
}
```

### Algorithm Flow

1. **Initialization**: Check if in git repo with `.envrc` file
2. **Hierarchy Building**: Parse `.envrc` to get ordered environment variables
3. **Branch Validation**: Verify each mapped branch exists
4. **Commit Analysis**: For each adjacent pair in hierarchy:
   - Execute `git log UPPER..LOWER` to find missing commits
   - Parse commit information
   - Identify user-authored commits
5. **Alert Generation**: Format and display results to stderr
6. **Continuation**: Allow primary gbm command to proceed

## Testing Strategy

### Unit Tests
- **File Parsing**: Various `.envrc` formats and edge cases
- **Git Command Parsing**: Mock git command outputs
- **User Identification**: Different git config scenarios
- **Error Handling**: Missing files, invalid git repos, command failures

### Integration Tests
- **Real Repository**: Test with actual git repositories
- **Multi-branch Scenarios**: Create merge-back situations
- **User Simulation**: Set git config and test detection
- **Command Integration**: Verify alerts appear before command output

### Edge Cases
- Empty `.envrc` files
- Single branch configurations
- Non-existent branches in hierarchy
- Corrupted git repositories
- Missing git configuration

### Sample Alert Outputs

**Clean State (No Alerts)**:
```bash
$ gbm list
# Normal gbm list output with no merge-back alerts
```

**Single User Commit**:
```bash
$ gbm list
⚠️  Merge-back required in tracked branches:

PROD → PREVIEW: 1 commit needs merge-back (1 by you)
• a1b2c3d - Fix login timeout issue (you - 1 day ago)

# Normal gbm list output follows
```

**Multiple Branches with Mixed Authors**:
```bash
$ gbm sync
⚠️  Merge-back required in tracked branches:

PROD → PREVIEW: 3 commits need merge-back (2 by you)
• a1b2c3d - Fix critical security bug (you - 2 days ago)
• e4f5g6h - Update API endpoint timeout (you - 1 day ago)

PREVIEW → MAIN: 7 commits need merge-back (0 by you)

# Normal gbm sync output follows
```

### .envrc Examples

**Standard Three-Tier**:
```bash
MAIN=main
PREVIEW=preview-v2.1
PROD=production-v2.0
```

**Dynamic Two-Tier** (after PREVIEW promotion):
```bash
MAIN=main
PROD=production-v2.1
```

**Complex Environment Setup**:
```bash
MAIN=main
STAGING=staging-branch
PREVIEW=preview-v3.0
PROD=production-v2.5
HOTFIX=hotfix-v2.5.1
```

