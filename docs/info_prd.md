# Product Requirements Document: `gbm info` Command

## Overview

The `gbm info` command provides comprehensive information about a specific worktree, including its Git status, JIRA ticket details (if applicable), and other relevant metadata. This command combines local Git information with external JIRA data to give users a complete picture of their worktree's state.

## Functional Requirements

### Core Command Behavior

**Command Signature:**
```bash
gbm info <worktree-name>
```
or if in a worktree directory
```bash
gbm info .
```

**Primary Function:**
- Display detailed information about a specified worktree in a visually appealing, structured format
- Integrate JIRA ticket information when the worktree name matches a JIRA key pattern
- Show Git repository status, branch information, and recent activity
- Present information in themed sections using `charmbracelet/lipgloss` styling

### Information Sections

Based on the mockup in `docs/info_ascii_mockup.md`, the command should display:

#### 1. Header Section
- Worktree name prominently displayed
- Visual separator/border using lipgloss styling

#### 2. Worktree Information Section
- **Name**: Worktree identifier
- **Path**: Absolute path to worktree directory
- **Branch**: Current branch name
- **Created**: Creation timestamp (relative time format)
- **Status**: Git status indicator (clean/dirty with file count)

#### 3. JIRA Ticket Section (if applicable)
- **Key**: JIRA ticket identifier
- **Summary**: Ticket title/description
- **Status**: Current JIRA status with transition arrows
- **Assignee**: Assigned user
- **Priority**: Priority level with visual indicator
- **Reporter**: Ticket creator
- **Created**: Ticket creation date
- **Due Date**: Due date (if set)
- **Epic**: Parent epic information (if applicable)
- **Link**: Direct URL to JIRA ticket
- **Latest Comment**: Most recent comment with timestamp and author

#### 4. Git Status Section
- **Base Branch**: Branch this was created from with divergence info
- **Upstream**: Remote tracking branch information
- **Position**: Commits ahead/behind upstream
- **Last Commit**: Most recent commit info (hash, message, author, timestamp)
- **Modified Files**: List of changed files with line change counts
- **Recent Commits**: Last 4-5 commits with condensed info

## Technical Requirements

### Dependencies

#### Required Libraries
- `github.com/charmbracelet/lipgloss` - For styling and layout (already in go.mod)
- External `jira` CLI tool - For JIRA integration (must be installed and configured)

#### Integration Points
- **JIRA CLI**: Use existing `internal/jira.go` functions and extend as needed
- **Git**: Leverage existing `internal/git.go` functions
- **Styling**: Use existing `internal/styles.go` patterns

### Data Sources

#### Git Information (Available)
```go
// From existing codebase
- Worktree path and creation info
- Current branch and upstream tracking
- Git status (ahead/behind, dirty state)
- Commit history and file changes
- Branch divergence information
```

#### JIRA Information (Via jira-cli)
```bash
# Available commands to leverage
jira issue view <key>           # Detailed ticket info
jira issue list -a$(jira me)    # User's tickets
jira issue comments <key>       # Ticket comments
```

#### File System Information
```go
// Need to implement
- Worktree creation timestamp (via os.Stat)
- Modified file analysis (git diff --stat)
```

### Implementation Strategy

#### 1. Command Structure
```go
// cmd/info.go
var infoCmd = &cobra.Command{
    Use:   "info <worktree-name>",
    Short: "Display detailed information about a worktree",
    Args:  cobra.ExactArgs(1),
    RunE:  runInfoCommand,
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
    // Validation and data gathering
    // Render using lipgloss components
}
```

#### 2. Data Structures
```go
type WorktreeInfoData struct {
    // Worktree metadata
    Name         string
    Path         string
    Branch       string
    CreatedAt    time.Time
    GitStatus    *GitStatus
    
    // Git details
    BaseInfo     *BranchInfo
    Commits      []CommitInfo
    ModifiedFiles []FileChange
    
    // JIRA integration (optional)
    JiraTicket   *JiraTicketDetails
}

type JiraTicketDetails struct {
    Key          string
    Summary      string
    Status       string
    Assignee     string
    Priority     string
    Reporter     string
    Created      time.Time
    DueDate      *time.Time
    Epic         string
    URL          string
    LatestComment *Comment
}
```

#### 3. Rendering Components
```go
// Use lipgloss for consistent styling
func renderHeader(worktreeName string) string
func renderWorktreeSection(data *WorktreeInfoData) string
func renderJiraSection(jira *JiraTicketDetails) string
func renderGitSection(data *WorktreeInfoData) string
```

## Data Availability Assessment

### Readily Available (✅)
- Worktree name, path, current branch
- Basic Git status (dirty/clean, ahead/behind)
- Recent commits and commit messages
- Current file modifications

### Requires Implementation (🔨)
- **Worktree creation timestamp**: Use `os.Stat()` on worktree directory
- **Base branch detection**: Analyze `git log --graph` or use `git merge-base`
- **File change statistics**: Parse `git diff --stat HEAD~1` output

### External Dependencies (🔗)
- **JIRA ticket details**: Extend existing `internal/jira.go` with:
  ```go
  func GetJiraTicketDetails(key string) (*JiraTicketDetails, error)
  func GetJiraComments(key string) ([]Comment, error)
  ```

### Potentially Unavailable (❓)
- **Epic information**: Depends on JIRA configuration and permissions
- **Due dates**: May not be set on all tickets

## Implementation Plan

### Phase 1: Core Command Structure
1. Create `cmd/info.go` with basic command structure
2. Implement worktree validation and data gathering
3. Create basic lipgloss layouts for each section

### Phase 2: Git Information
1. Extend `internal/git.go` with detailed status functions
2. Implement file change analysis
3. Add branch relationship detection
4. Create Git status rendering components

### Phase 3: JIRA Integration
1. Extend `internal/jira.go` with detailed ticket fetching
2. Add comment retrieval functionality
3. Implement JIRA section rendering
4. Add error handling for missing JIRA CLI or permissions

### Phase 4: Styling and Polish
1. Refine lipgloss styling to match mockup aesthetics
2. Add responsive layout for different terminal widths
3. Implement proper error handling and fallbacks
4. Add comprehensive testing

## User Experience Considerations

### Success Cases
- **Standard worktree**: Shows all Git information cleanly
- **JIRA-linked worktree**: Displays comprehensive ticket details
- **No JIRA access**: Gracefully shows only Git information

### Error Handling
- **Invalid worktree name**: Clear error message with suggestions
- **JIRA CLI unavailable**: Show Git info only with informational note
- **Network issues**: Timeout gracefully, show cached/local data

### Performance
- Cache JIRA responses to avoid repeated API calls
- Lazy-load expensive operations (file diff analysis)
- Provide progress indication for slow operations

## Acceptance Criteria

1. ✅ Command displays comprehensive worktree information
2. ✅ JIRA integration works when CLI is available and configured
3. ✅ Graceful degradation when JIRA is unavailable
4. ✅ Visual layout matches mockup aesthetic and readability
5. ✅ Command completes within 2 seconds for standard cases
6. ✅ Error messages are helpful and actionable
7. ✅ Styling is consistent with existing CLI patterns

## Future Enhancements

- **GitHub/GitLab integration**: Automatic PR detection and status
- **Caching layer**: Store JIRA responses for offline access
- **Interactive mode**: Allow drilling down into specific sections
- **Export options**: JSON output for scripting integration
