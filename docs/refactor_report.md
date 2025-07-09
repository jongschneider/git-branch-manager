# Refactor Analysis Report

This document contains analysis of functions and methods in the codebase to identify potential duplicates and redundancy patterns.

# 001_cmd_directory

## Directory: cmd/

### Functions Found

#### Execute
- **File**: `cmd/root.go:36`
- **Signature**: `func Execute() error`
- **Description**: Returns the result of rootCmd.Execute() - main entry point for CLI
- **Usage**: Called from main.go

#### init (root)
- **File**: `cmd/root.go:40`
- **Signature**: `func init()`
- **Description**: Initializes persistent flags for root command (config, worktree-dir, debug)
- **Usage**: Called automatically by Go runtime

#### GetConfigPath
- **File**: `cmd/root.go:46`
- **Signature**: `func GetConfigPath() string`
- **Description**: Returns the config path from flag or default ".envrc"
- **Usage**: Called throughout cmd package to get config path

#### GetWorktreeDir
- **File**: `cmd/root.go:53`
- **Signature**: `func GetWorktreeDir() string`
- **Description**: Returns the worktree directory from flag or default "worktrees"
- **Usage**: Called throughout cmd package to get worktree directory

#### InitializeLogging
- **File**: `cmd/root.go:60`
- **Signature**: `func InitializeLogging()`
- **Description**: Sets up debug logging to gbm.log file if debug flag is enabled
- **Usage**: Called in PersistentPreRun hook

#### PrintInfo
- **File**: `cmd/root.go:70`
- **Signature**: `func PrintInfo(format string, args ...interface{})`
- **Description**: Prints formatted info message to stderr and optionally to log file
- **Usage**: Called throughout cmd package for info messages

#### PrintVerbose
- **File**: `cmd/root.go:80`
- **Signature**: `func PrintVerbose(format string, args ...interface{})`
- **Description**: Prints formatted verbose message to stderr if debug enabled, always to log file
- **Usage**: Called throughout cmd package for debug messages

#### PrintError
- **File**: `cmd/root.go:92`
- **Signature**: `func PrintError(format string, args ...interface{})`
- **Description**: Prints formatted error message to stderr and optionally to log file
- **Usage**: Called throughout cmd package for error messages

#### CloseLogFile
- **File**: `cmd/root.go:102`
- **Signature**: `func CloseLogFile()`
- **Description**: Closes the log file if it exists
- **Usage**: Called for cleanup (not visible in current files)

#### checkAndDisplayMergeBackAlerts
- **File**: `cmd/root.go:108`
- **Signature**: `func checkAndDisplayMergeBackAlerts()`
- **Description**: Checks for merge-back status and displays alerts if needed
- **Usage**: Called in PersistentPreRun hook

#### handleInteractive
- **File**: `cmd/add.go:84`
- **Signature**: `func handleInteractive(manager *internal.Manager) (string, error)`
- **Description**: Handles interactive mode for adding worktrees - shows branch selection menu
- **Usage**: Called from add command when interactive flag is set

#### generateBranchName
- **File**: `cmd/add.go:123`
- **Signature**: `func generateBranchName(worktreeName string) string`
- **Description**: Generates a branch name from worktree name, using JIRA API if worktree is JIRA key
- **Usage**: Called from add command when creating new branch

#### init (add)
- **File**: `cmd/add.go:151`
- **Signature**: `func init()`
- **Description**: Initializes add command with flags and completion function
- **Usage**: Called automatically by Go runtime

#### runGitBareClone
- **File**: `cmd/clone.go:58`
- **Signature**: `func runGitBareClone(repoUrl string) error`
- **Description**: Clones a repository as bare repo, configures remote, and fetches branches
- **Usage**: Called from clone command

#### extractRepoName
- **File**: `cmd/clone.go:113`
- **Signature**: `func extractRepoName(repoUrl string) string`
- **Description**: Extracts repository name from URL by removing .git suffix and taking last path component
- **Usage**: Called from runGitBareClone

#### getDefaultBranch
- **File**: `cmd/clone.go:126`
- **Signature**: `func getDefaultBranch() (string, error)`
- **Description**: Determines the default branch of the repository using git remote set-head and symbolic-ref
- **Usage**: Called from clone command

#### createMainWorktree
- **File**: `cmd/clone.go:171`
- **Signature**: `func createMainWorktree(defaultBranch string) error`
- **Description**: Creates the main worktree directory and adds it using git worktree add
- **Usage**: Called from clone command

#### setupEnvrc
- **File**: `cmd/clone.go:188`
- **Signature**: `func setupEnvrc(defaultBranch string) error`
- **Description**: Sets up .envrc file by copying from worktree or creating default
- **Usage**: Called from clone command

#### copyFile
- **File**: `cmd/clone.go:219`
- **Signature**: `func copyFile(src, dst string) error`
- **Description**: Copies file from src to dst using io.Copy
- **Usage**: Called from setupEnvrc

#### createDefaultEnvrc
- **File**: `cmd/clone.go:236`
- **Signature**: `func createDefaultEnvrc(path, defaultBranch string) error`
- **Description**: Creates a default .envrc file with MAIN=defaultBranch
- **Usage**: Called from setupEnvrc

#### initializeWorktreeManagement
- **File**: `cmd/clone.go:242`
- **Signature**: `func initializeWorktreeManagement() error`
- **Description**: Creates manager, loads .envrc, and performs initial sync
- **Usage**: Called from clone command

#### init (clone)
- **File**: `cmd/clone.go:272`
- **Signature**: `func init()`
- **Description**: Initializes clone command
- **Usage**: Called automatically by Go runtime

#### init (completion)
- **File**: `cmd/completion.go:66`
- **Signature**: `func init()`
- **Description**: Initializes completion command
- **Usage**: Called automatically by Go runtime

#### runInfoCommand
- **File**: `cmd/info.go:30`
- **Signature**: `func runInfoCommand(cmd *cobra.Command, args []string) error`
- **Description**: Main function for info command - handles current directory reference and calls getWorktreeInfo
- **Usage**: Used as RunE for info command

#### getWorktreeInfo
- **File**: `cmd/info.go:66`
- **Signature**: `func getWorktreeInfo(gitManager *internal.GitManager, worktreeName string) (*internal.WorktreeInfoData, error)`
- **Description**: Gathers comprehensive worktree information including git status, commits, files, and JIRA details
- **Usage**: Called from runInfoCommand

#### displayWorktreeInfo
- **File**: `cmd/info.go:139`
- **Signature**: `func displayWorktreeInfo(data *internal.WorktreeInfoData)`
- **Description**: Renders worktree information using InfoRenderer
- **Usage**: Called from runInfoCommand

#### getWorktreeCreationTime
- **File**: `cmd/info.go:145`
- **Signature**: `func getWorktreeCreationTime(worktreePath string) (time.Time, error)`
- **Description**: Gets worktree creation time from file system stat
- **Usage**: Called from getWorktreeInfo

#### getRecentCommits
- **File**: `cmd/info.go:153`
- **Signature**: `func getRecentCommits(worktreePath string, count int) ([]internal.CommitInfo, error)`
- **Description**: Gets recent commits from worktree using git log with custom format
- **Usage**: Called from getWorktreeInfo

#### getModifiedFiles
- **File**: `cmd/info.go:188`
- **Signature**: `func getModifiedFiles(worktreePath string) ([]internal.FileChange, error)`
- **Description**: Gets modified files from worktree using git diff --numstat for both staged and unstaged
- **Usage**: Called from getWorktreeInfo

#### getBaseBranchInfo
- **File**: `cmd/info.go:282`
- **Signature**: `func getBaseBranchInfo(worktreePath string) (*internal.BranchInfo, error)`
- **Description**: Gets base branch information including upstream and ahead/behind counts
- **Usage**: Called from getWorktreeInfo

#### getJiraTicketDetails
- **File**: `cmd/info.go:390`
- **Signature**: `func getJiraTicketDetails(jiraKey string) (*internal.JiraTicketDetails, error)`
- **Description**: Fetches JIRA ticket details using jira CLI with raw JSON output
- **Usage**: Called from getWorktreeInfo if worktree name is JIRA key

#### formatJiraURL
- **File**: `cmd/info.go:496`
- **Signature**: `func formatJiraURL(selfURL, key string) string`
- **Description**: Converts JIRA API URL to user-friendly browse URL
- **Usage**: Called from getJiraTicketDetails

#### init (info)
- **File**: `cmd/info.go:515`
- **Signature**: `func init()`
- **Description**: Initializes info command
- **Usage**: Called automatically by Go runtime

#### init (list)
- **File**: `cmd/list.go:120`
- **Signature**: `func init()`
- **Description**: Initializes list command
- **Usage**: Called automatically by Go runtime

#### handlePullAll
- **File**: `cmd/pull.go:54`
- **Signature**: `func handlePullAll(manager *internal.Manager) error`
- **Description**: Handles pulling all worktrees using manager.PullAllWorktrees()
- **Usage**: Called from pull command when --all flag is set

#### handlePullCurrent
- **File**: `cmd/pull.go:59`
- **Signature**: `func handlePullCurrent(manager *internal.Manager, currentPath string) error`
- **Description**: Handles pulling current worktree by checking if in worktree and calling PullWorktree
- **Usage**: Called from pull command when no arguments provided

#### handlePullNamed
- **File**: `cmd/pull.go:74`
- **Signature**: `func handlePullNamed(manager *internal.Manager, worktreeName string) error`
- **Description**: Handles pulling named worktree by validating existence and calling PullWorktree
- **Usage**: Called from pull command when worktree name is provided

#### init (pull)
- **File**: `cmd/pull.go:89`
- **Signature**: `func init()`
- **Description**: Initializes pull command with flags and completion function
- **Usage**: Called automatically by Go runtime

#### handlePushAll
- **File**: `cmd/push.go:56`
- **Signature**: `func handlePushAll(manager *internal.Manager) error`
- **Description**: Handles pushing all worktrees using manager.PushAllWorktrees()
- **Usage**: Called from push command when --all flag is set

#### handlePushCurrent
- **File**: `cmd/push.go:61`
- **Signature**: `func handlePushCurrent(manager *internal.Manager, currentPath string) error`
- **Description**: Handles pushing current worktree by checking if in worktree and calling PushWorktree
- **Usage**: Called from push command when no arguments provided

#### handlePushNamed
- **File**: `cmd/push.go:76`
- **Signature**: `func handlePushNamed(manager *internal.Manager, worktreeName string) error`
- **Description**: Handles pushing named worktree by validating existence and calling PushWorktree
- **Usage**: Called from push command when worktree name is provided

#### init (push)
- **File**: `cmd/push.go:91`
- **Signature**: `func init()`
- **Description**: Initializes push command with flags and completion function
- **Usage**: Called from push command

#### init (remove)
- **File**: `cmd/remove.go:83`
- **Signature**: `func init()`
- **Description**: Initializes remove command with flags and completion function
- **Usage**: Called automatically by Go runtime

#### init (shell-integration)
- **File**: `cmd/shell-integration.go:90`
- **Signature**: `func init()`
- **Description**: Initializes shell-integration command
- **Usage**: Called automatically by Go runtime

#### switchToWorktree
- **File**: `cmd/switch.go:86`
- **Signature**: `func switchToWorktree(manager *internal.Manager, worktreeName string) error`
- **Description**: Switches to a worktree by exact or fuzzy matching and outputs appropriate response
- **Usage**: Called from switch command

#### findFuzzyMatch
- **File**: `cmd/switch.go:130`
- **Signature**: `func findFuzzyMatch(manager *internal.Manager, target string) string`
- **Description**: Finds fuzzy match for worktree name using case-insensitive substring matching
- **Usage**: Called from switchToWorktree

#### listWorktrees
- **File**: `cmd/switch.go:168`
- **Signature**: `func listWorktrees(manager *internal.Manager) error`
- **Description**: Lists all available worktrees with status and branch information
- **Usage**: Called from switch command when no arguments provided

#### init (switch)
- **File**: `cmd/switch.go:204`
- **Signature**: `func init()`
- **Description**: Initializes switch command with flags and completion function
- **Usage**: Called automatically by Go runtime

#### getWorktreeNames
- **File**: `cmd/switch.go:217`
- **Signature**: `func getWorktreeNames() []string`
- **Description**: Gets list of worktree names for completion
- **Usage**: Called from completion functions in multiple commands

#### init (sync)
- **File**: `cmd/sync.go:96`
- **Signature**: `func init()`
- **Description**: Initializes sync command with flags
- **Usage**: Called automatically by Go runtime

#### init (validate)
- **File**: `cmd/validate.go:86`
- **Signature**: `func init()`
- **Description**: Initializes validate command
- **Usage**: Called automatically by Go runtime

### Potential Duplicates/Redundancy

#### Pattern 1: Handle Functions for Push/Pull Commands
- **Functions**: `handlePushAll`, `handlePushCurrent`, `handlePushNamed`, `handlePullAll`, `handlePullCurrent`, `handlePullNamed`
- **Similarity**: All follow identical pattern: validate inputs, call appropriate manager method
- **Suggestion**: Could be consolidated into generic handlers or use method composition

#### Pattern 2: Init Functions
- **Functions**: Multiple `init()` functions across all cmd files
- **Similarity**: All follow pattern of `rootCmd.AddCommand(cmdVar)` and flag setup
- **Suggestion**: No consolidation needed - this is standard Go pattern

#### Pattern 3: Completion Functions
- **Functions**: Multiple `ValidArgsFunction` implementations in add.go, pull.go, push.go, remove.go, switch.go
- **Similarity**: All call `getWorktreeNames()` for completion
- **Suggestion**: Could extract to shared completion helper function

#### Pattern 4: Manager Creation Pattern
- **Functions**: Nearly all commands create manager with identical pattern: `FindGitRoot` -> `NewManager` -> `LoadEnvMapping`
- **Similarity**: Same initialization sequence repeated across commands
- **Suggestion**: Could extract to helper function like `createInitializedManager()`

#### Pattern 5: Error Handling Pattern
- **Functions**: All command functions use similar error wrapping patterns
- **Similarity**: `fmt.Errorf("failed to X: %w", err)` pattern repeated throughout
- **Suggestion**: Could use consistent error wrapping helpers

#### Pattern 6: Worktree Existence Validation
- **Functions**: `handlePullNamed`, `handlePushNamed`, remove command validation
- **Similarity**: All check if worktree exists before operating
- **Suggestion**: Could extract to shared validation function

# 002_internal_core

## Directory: internal/

### Functions Found

#### DefaultConfig
- **File**: `internal/config.go:59`
- **Signature**: `func DefaultConfig() *Config`
- **Description**: Returns default configuration with predefined settings, state, and icons
- **Usage**: Called from LoadConfig when no config file exists

#### LoadConfig
- **File**: `internal/config.go:95`
- **Signature**: `func LoadConfig(gbmDir string) (*Config, error)`
- **Description**: Loads configuration from config.toml file or returns default if not found
- **Usage**: Called from NewManager to initialize configuration

#### Save (Config)
- **File**: `internal/config.go:110`
- **Signature**: `func (c *Config) Save(gbmDir string) error`
- **Description**: Saves configuration to config.toml file using TOML encoder
- **Usage**: Called from manager sync operations to persist state

#### ParseEnvrc
- **File**: `internal/config.go:130`
- **Signature**: `func ParseEnvrc(path string) (*EnvMapping, error)`
- **Description**: Parses .envrc file using regex to extract environment variable mappings
- **Usage**: Called from LoadEnvMapping to read .envrc file

#### HasChanges (GitStatus)
- **File**: `internal/git.go:39`
- **Signature**: `func (gs *GitStatus) HasChanges() bool`
- **Description**: Returns true if git status has any changes (dirty, untracked, modified, staged)
- **Usage**: Called to check if worktree has uncommitted changes

#### execCommand
- **File**: `internal/git.go:44`
- **Signature**: `func execCommand(cmd *exec.Cmd) ([]byte, error)`
- **Description**: Executes command and returns output - wrapper around cmd.Output()
- **Usage**: Called throughout git.go for commands that return output

#### execCommandRun
- **File**: `internal/git.go:50`
- **Signature**: `func execCommandRun(cmd *exec.Cmd) error`
- **Description**: Executes command using Run() instead of Output() - wrapper around cmd.Run()
- **Usage**: Called throughout git.go for commands that don't need output

#### FindGitRoot
- **File**: `internal/git.go:55`
- **Signature**: `func FindGitRoot(startPath string) (string, error)`
- **Description**: Complex function to find git repository root, handles worktrees, bare repos, and subdirectories
- **Usage**: Called from commands to locate git repository root

#### NewGitManager
- **File**: `internal/git.go:174`
- **Signature**: `func NewGitManager(repoPath string) (*GitManager, error)`
- **Description**: Creates new GitManager instance with go-git repository
- **Usage**: Called from NewManager to initialize git operations

#### IsGitRepository
- **File**: `internal/git.go:186`
- **Signature**: `func (gm *GitManager) IsGitRepository() bool`
- **Description**: Checks if path is a valid git repository
- **Usage**: Called to validate repository status

#### BranchExists
- **File**: `internal/git.go:191`
- **Signature**: `func (gm *GitManager) BranchExists(branchName string) (bool, error)`
- **Description**: Checks if branch exists locally or remotely using git references
- **Usage**: Called before creating worktrees to validate branch existence

#### GetWorktrees
- **File**: `internal/git.go:226`
- **Signature**: `func (gm *GitManager) GetWorktrees() ([]*WorktreeInfo, error)`
- **Description**: Gets all worktrees using git worktree list --porcelain
- **Usage**: Called from manager to get current worktree status

#### CreateWorktree
- **File**: `internal/git.go:270`
- **Signature**: `func (gm *GitManager) CreateWorktree(envVar, branchName, worktreeDir string) error`
- **Description**: Creates new worktree, handles branch tracking and remote setup
- **Usage**: Called from sync operations to create missing worktrees

#### RemoveWorktree
- **File**: `internal/git.go:322`
- **Signature**: `func (gm *GitManager) RemoveWorktree(worktreePath string) error`
- **Description**: Removes worktree using git worktree remove --force
- **Usage**: Called from sync operations and remove command

#### UpdateWorktree
- **File**: `internal/git.go:332`
- **Signature**: `func (gm *GitManager) UpdateWorktree(worktreePath, newBranch string) error`
- **Description**: Updates worktree to new branch by removing and recreating
- **Usage**: Called from sync operations when branch changes

#### FetchAll
- **File**: `internal/git.go:344`
- **Signature**: `func (gm *GitManager) FetchAll() error`
- **Description**: Fetches all branches from remote using go-git
- **Usage**: Called from sync operations when fetch flag is set

#### GetWorktreeStatus
- **File**: `internal/git.go:358`
- **Signature**: `func (gm *GitManager) GetWorktreeStatus(worktreePath string) (*GitStatus, error)`
- **Description**: Gets comprehensive git status for worktree including ahead/behind counts
- **Usage**: Called from list and info commands to show worktree status

#### GetStatusIcon
- **File**: `internal/git.go:419`
- **Signature**: `func (gm *GitManager) GetStatusIcon(gitStatus *GitStatus) string`
- **Description**: Returns appropriate icon string based on git status
- **Usage**: Called from manager to display status icons

#### CreateBranch
- **File**: `internal/git.go:457`
- **Signature**: `func (gm *GitManager) CreateBranch(branchName, baseBranch string) error`
- **Description**: Creates new branch using git branch command
- **Usage**: Called when creating branches for new worktrees

#### AddWorktree
- **File**: `internal/git.go:471`
- **Signature**: `func (gm *GitManager) AddWorktree(worktreeName, branchName string, createBranch bool) error`
- **Description**: Adds worktree with optional branch creation and remote tracking
- **Usage**: Called from add command to create new worktrees

#### GetCurrentBranch
- **File**: `internal/git.go:534`
- **Signature**: `func (gm *GitManager) GetCurrentBranch() (string, error)`
- **Description**: Gets current branch name using git rev-parse --abbrev-ref HEAD
- **Usage**: Called from various commands to get current branch

#### GetRemoteBranches
- **File**: `internal/git.go:545`
- **Signature**: `func (gm *GitManager) GetRemoteBranches() ([]string, error)`
- **Description**: Gets list of remote branches using git branch -r
- **Usage**: Called from add command for branch selection

#### PushWorktree
- **File**: `internal/git.go:571`
- **Signature**: `func (gm *GitManager) PushWorktree(worktreePath string) error`
- **Description**: Pushes worktree changes, sets upstream if needed
- **Usage**: Called from push command and PushAllWorktrees

#### PullWorktree
- **File**: `internal/git.go:606`
- **Signature**: `func (gm *GitManager) PullWorktree(worktreePath string) error`
- **Description**: Pulls worktree changes, sets upstream if needed
- **Usage**: Called from pull command and PullAllWorktrees

#### IsInWorktree
- **File**: `internal/git.go:661`
- **Signature**: `func (gm *GitManager) IsInWorktree(currentPath string) (bool, string, error)`
- **Description**: Checks if current path is within a worktree and returns worktree name
- **Usage**: Called from commands to determine current worktree context

#### NewManager
- **File**: `internal/manager.go:39`
- **Signature**: `func NewManager(repoPath string) (*Manager, error)`
- **Description**: Creates new Manager instance with config, git manager, and icon manager
- **Usage**: Called from all commands to initialize management layer

#### LoadEnvMapping
- **File**: `internal/manager.go:63`
- **Signature**: `func (m *Manager) LoadEnvMapping(envrcPath string) error`
- **Description**: Loads environment variable mapping from .envrc file
- **Usage**: Called from commands to load worktree configuration

#### GetSyncStatus
- **File**: `internal/manager.go:81`
- **Signature**: `func (m *Manager) GetSyncStatus() (*SyncStatus, error)`
- **Description**: Compares current worktrees with .envrc mapping to determine sync status
- **Usage**: Called from sync command to determine what actions are needed

#### Sync
- **File**: `internal/manager.go:129`
- **Signature**: `func (m *Manager) Sync(dryRun, force, fetch bool) error`
- **Description**: Synchronizes worktrees with .envrc mapping, handles creation/removal/updates
- **Usage**: Called from sync command to perform synchronization

#### ValidateEnvrc
- **File**: `internal/manager.go:195`
- **Signature**: `func (m *Manager) ValidateEnvrc() error`
- **Description**: Validates that all branches in .envrc exist in repository
- **Usage**: Called from sync operations to validate configuration

#### GetEnvMapping
- **File**: `internal/manager.go:213`
- **Signature**: `func (m *Manager) GetEnvMapping() (map[string]string, error)`
- **Description**: Returns environment variable mapping
- **Usage**: Called from commands to access mapping

#### BranchExists (Manager)
- **File**: `internal/manager.go:220`
- **Signature**: `func (m *Manager) BranchExists(branchName string) (bool, error)`
- **Description**: Delegate to GitManager.BranchExists
- **Usage**: Called from commands to check branch existence

#### GetWorktreeList
- **File**: `internal/manager.go:224`
- **Signature**: `func (m *Manager) GetWorktreeList() (map[string]*WorktreeListInfo, error)`
- **Description**: Gets worktree list based on .envrc mapping with status information
- **Usage**: Called from list command for .envrc-based worktrees

#### GetStatusIcon (Manager)
- **File**: `internal/manager.go:263`
- **Signature**: `func (m *Manager) GetStatusIcon(gitStatus *GitStatus) string`
- **Description**: Delegate to GitManager.GetStatusIcon
- **Usage**: Called from commands to get status icons

#### GetWorktreePath
- **File**: `internal/manager.go:267`
- **Signature**: `func (m *Manager) GetWorktreePath(worktreeName string) (string, error)`
- **Description**: Gets absolute path to worktree directory
- **Usage**: Called from commands to get worktree paths

#### GetAllWorktrees
- **File**: `internal/manager.go:277`
- **Signature**: `func (m *Manager) GetAllWorktrees() (map[string]*WorktreeListInfo, error)`
- **Description**: Gets all worktrees including ad-hoc ones not in .envrc
- **Usage**: Called from commands to get complete worktree listing

#### AddWorktree (Manager)
- **File**: `internal/manager.go:321`
- **Signature**: `func (m *Manager) AddWorktree(worktreeName, branchName string, createBranch bool) error`
- **Description**: Adds worktree and tracks as ad-hoc if not in .envrc
- **Usage**: Called from add command

#### contains
- **File**: `internal/manager.go:346`
- **Signature**: `func contains(slice []string, item string) bool`
- **Description**: Helper function to check if slice contains string
- **Usage**: Called from AddWorktree to check ad-hoc worktrees

#### GetRemoteBranches (Manager)
- **File**: `internal/manager.go:355`
- **Signature**: `func (m *Manager) GetRemoteBranches() ([]string, error)`
- **Description**: Delegate to GitManager.GetRemoteBranches
- **Usage**: Called from commands to get remote branches

#### GetCurrentBranch (Manager)
- **File**: `internal/manager.go:359`
- **Signature**: `func (m *Manager) GetCurrentBranch() (string, error)`
- **Description**: Delegate to GitManager.GetCurrentBranch
- **Usage**: Called from commands to get current branch

#### PushWorktree (Manager)
- **File**: `internal/manager.go:363`
- **Signature**: `func (m *Manager) PushWorktree(worktreeName string) error`
- **Description**: Pushes specific worktree by name
- **Usage**: Called from push command for individual worktrees

#### PullWorktree (Manager)
- **File**: `internal/manager.go:368`
- **Signature**: `func (m *Manager) PullWorktree(worktreeName string) error`
- **Description**: Pulls specific worktree by name
- **Usage**: Called from pull command for individual worktrees

#### IsInWorktree (Manager)
- **File**: `internal/manager.go:373`
- **Signature**: `func (m *Manager) IsInWorktree(currentPath string) (bool, string, error)`
- **Description**: Delegate to GitManager.IsInWorktree
- **Usage**: Called from commands to check worktree context

#### PushAllWorktrees
- **File**: `internal/manager.go:377`
- **Signature**: `func (m *Manager) PushAllWorktrees() error`
- **Description**: Pushes all worktrees with progress reporting
- **Usage**: Called from push command with --all flag

#### PullAllWorktrees
- **File**: `internal/manager.go:395`
- **Signature**: `func (m *Manager) PullAllWorktrees() error`
- **Description**: Pulls all worktrees with progress reporting
- **Usage**: Called from pull command with --all flag

#### RemoveWorktree (Manager)
- **File**: `internal/manager.go:413`
- **Signature**: `func (m *Manager) RemoveWorktree(worktreeName string) error`
- **Description**: Removes worktree and updates ad-hoc tracking
- **Usage**: Called from remove command

#### GetWorktreeStatus (Manager)
- **File**: `internal/manager.go:438`
- **Signature**: `func (m *Manager) GetWorktreeStatus(worktreePath string) (*GitStatus, error)`
- **Description**: Delegate to GitManager.GetWorktreeStatus
- **Usage**: Called from commands to get worktree status

#### SetCurrentWorktree
- **File**: `internal/manager.go:442`
- **Signature**: `func (m *Manager) SetCurrentWorktree(worktreeName string) error`
- **Description**: Sets current worktree and updates previous worktree tracking
- **Usage**: Called from switch command to track current worktree

#### GetPreviousWorktree
- **File**: `internal/manager.go:451`
- **Signature**: `func (m *Manager) GetPreviousWorktree() string`
- **Description**: Gets previously active worktree name
- **Usage**: Called from switch command for switching back

#### GetCurrentWorktree
- **File**: `internal/manager.go:455`
- **Signature**: `func (m *Manager) GetCurrentWorktree() string`
- **Description**: Gets currently active worktree name
- **Usage**: Called from commands to get current worktree context

#### GetSortedWorktreeNames
- **File**: `internal/manager.go:459`
- **Signature**: `func (m *Manager) GetSortedWorktreeNames(worktrees map[string]*WorktreeListInfo) []string`
- **Description**: Sorts worktree names with .envrc worktrees first, then ad-hoc by creation time
- **Usage**: Called from list command to order worktrees

#### IsJiraKey
- **File**: `internal/jira.go:19`
- **Signature**: `func IsJiraKey(s string) bool`
- **Description**: Checks if string matches JIRA key pattern (PROJECT-NUMBER)
- **Usage**: Called throughout codebase to identify JIRA keys

#### GetJiraKeys
- **File**: `internal/jira.go:25`
- **Signature**: `func GetJiraKeys() ([]string, error)`
- **Description**: Fetches all JIRA issue keys for current user using jira CLI
- **Usage**: Called from add command for JIRA integration

#### GetJiraIssues
- **File**: `internal/jira.go:60`
- **Signature**: `func GetJiraIssues() ([]JiraIssue, error)`
- **Description**: Fetches all JIRA issues for current user with full details
- **Usage**: Called from commands that need JIRA issue information

#### GetJiraIssue
- **File**: `internal/jira.go:80`
- **Signature**: `func GetJiraIssue(key string) (*JiraIssue, error)`
- **Description**: Fetches detailed information for specific JIRA issue
- **Usage**: Called from branch name generation and info commands

#### ParseJiraList
- **File**: `internal/jira.go:147`
- **Signature**: `func ParseJiraList(output string) []JiraIssue`
- **Description**: Parses output of 'jira issue list' command into JiraIssue structs
- **Usage**: Called from GetJiraIssues to parse CLI output

#### BranchName (JiraIssue)
- **File**: `internal/jira.go:200`
- **Signature**: `func (j *JiraIssue) BranchName() string`
- **Description**: Generates filesystem-safe branch name from JIRA issue
- **Usage**: Called from branch name generation

#### GenerateBranchFromJira
- **File**: `internal/jira.go:218`
- **Signature**: `func GenerateBranchFromJira(jiraKey string) (string, error)`
- **Description**: Fetches JIRA issue and generates branch name
- **Usage**: Called from add command for JIRA-based branch creation

### Potential Duplicates/Redundancy

#### Pattern 1: Delegation Functions in Manager
- **Functions**: `BranchExists`, `GetStatusIcon`, `GetRemoteBranches`, `GetCurrentBranch`, `PushWorktree`, `PullWorktree`, `IsInWorktree`, `GetWorktreeStatus`
- **Similarity**: All simple delegates to GitManager methods
- **Suggestion**: Consider removing delegation layer or using embedded GitManager

#### Pattern 2: Git Command Execution
- **Functions**: `execCommand`, `execCommandRun`
- **Similarity**: Both are thin wrappers around exec.Cmd methods
- **Suggestion**: Could consolidate into single function with return type flag

#### Pattern 3: JIRA API Functions
- **Functions**: `GetJiraKeys`, `GetJiraIssues`, `GetJiraIssue`
- **Similarity**: All follow pattern of getting current user then listing issues
- **Suggestion**: Could extract common "get current user" functionality

#### Pattern 4: Branch Validation
- **Functions**: `BranchExists` (GitManager), `BranchExists` (Manager), `ValidateEnvrc`
- **Similarity**: All checking if branches exist in repository
- **Suggestion**: Could consolidate branch existence checking

#### Pattern 5: Worktree Path Construction
- **Functions**: `GetWorktreePath`, `PushWorktree`, `PullWorktree`, `RemoveWorktree`
- **Similarity**: All construct paths using same pattern: `filepath.Join(repoPath, prefix, name)`
- **Suggestion**: Could extract path construction to helper function

#### Pattern 6: Git Status Checking
- **Functions**: `HasChanges`, `GetWorktreeStatus`, `GetStatusIcon`
- **Similarity**: All deal with git status information
- **Suggestion**: Could consolidate status checking logic

#### Pattern 7: Push/Pull Operations
- **Functions**: `PushWorktree`, `PullWorktree`, `PushAllWorktrees`, `PullAllWorktrees`
- **Similarity**: Similar patterns for checking upstream and handling remote operations
- **Suggestion**: Could extract common upstream checking logic

# 003_internal_utils

## Directory: internal/

### Functions Found

#### NewIconManager
- **File**: `internal/styles.go:16`
- **Signature**: `func NewIconManager(config *Config) *IconManager`
- **Description**: Creates new icon manager with configuration
- **Usage**: Called from NewManager to initialize icon system

#### SetGlobalIconManager
- **File**: `internal/styles.go:21`
- **Signature**: `func SetGlobalIconManager(manager *IconManager) void`
- **Description**: Sets global icon manager instance
- **Usage**: Called from NewManager to set up global icon access

#### GetGlobalIconManager
- **File**: `internal/styles.go:26`
- **Signature**: `func GetGlobalIconManager() *IconManager`
- **Description**: Returns global icon manager instance with fallback to default
- **Usage**: Called from git status and formatting functions

#### Success/Warning/Error/Info (IconManager)
- **File**: `internal/styles.go:35-48`
- **Signature**: `func (im *IconManager) Success/Warning/Error/Info() string`
- **Description**: Icon getter methods for status icons
- **Usage**: Called from various formatting functions

#### FormatVerbose
- **File**: `internal/styles.go:133`
- **Signature**: `func FormatVerbose(text string) string`
- **Description**: Applies verbose style (subtle, italic) to text
- **Usage**: Called from cmd package for debug output

#### FormatHeader
- **File**: `internal/styles.go:137`
- **Signature**: `func FormatHeader(text string) string`
- **Description**: Applies header style (bold, primary color) to text
- **Usage**: Called from cmd package for section headers

#### FormatSubHeader
- **File**: `internal/styles.go:141`
- **Signature**: `func FormatSubHeader(text string) string`
- **Description**: Applies subheader style (bold, subtle) to text
- **Usage**: Called from cmd package for subsection headers

#### FormatBold
- **File**: `internal/styles.go:145`
- **Signature**: `func FormatBold(text string) string`
- **Description**: Applies bold style to text
- **Usage**: Called from cmd package for emphasis

#### FormatSubtle
- **File**: `internal/styles.go:149`
- **Signature**: `func FormatSubtle(text string) string`
- **Description**: Applies subtle style (gray foreground) to text
- **Usage**: Called from cmd package for secondary text

#### FormatPrompt
- **File**: `internal/styles.go:153`
- **Signature**: `func FormatPrompt(text string) string`
- **Description**: Applies prompt style (primary color, bold) to text
- **Usage**: Called from cmd package for user prompts

#### FormatStatusIcon
- **File**: `internal/styles.go:158`
- **Signature**: `func FormatStatusIcon(icon, text string) string`
- **Description**: Formats status icon with appropriate color based on icon type
- **Usage**: Called from other formatting functions

#### FormatSuccess/Warning/Error/Info
- **File**: `internal/styles.go:184-201`
- **Signature**: `func FormatSuccess/Warning/Error/Info(text string) string`
- **Description**: Helper functions for common status formatting
- **Usage**: Called from cmd package for status messages

#### FormatGitStatus
- **File**: `internal/styles.go:205`
- **Signature**: `func FormatGitStatus(status *GitStatus) string`
- **Description**: Formats git status with appropriate color and icon
- **Usage**: Called from git-related commands

#### NewTable
- **File**: `internal/table.go:18`
- **Signature**: `func NewTable(headers []string) *Table`
- **Description**: Creates new table with headers and terminal width detection
- **Usage**: Called from list command to create worktree tables

#### NewTestTable
- **File**: `internal/table.go:43`
- **Signature**: `func NewTestTable(headers []string, termWidth int) *Table`
- **Description**: Creates table with specific terminal width for testing
- **Usage**: Called from test files to create predictable table output

#### AddRow
- **File**: `internal/table.go:65`
- **Signature**: `func (t *Table) AddRow(row []string) void`
- **Description**: Adds row to table data
- **Usage**: Called from list command to populate table

#### Print
- **File**: `internal/table.go:69`
- **Signature**: `func (t *Table) Print() void`
- **Description**: Builds and prints table to stdout
- **Usage**: Called from list command to display table

#### String
- **File**: `internal/table.go:77`
- **Signature**: `func (t *Table) String() string`
- **Description**: Builds table and returns as string
- **Usage**: Called from test files to verify table output

#### buildTable
- **File**: `internal/table.go:86`
- **Signature**: `func (t *Table) buildTable() void`
- **Description**: Creates responsive table based on terminal width
- **Usage**: Called internally from Print and String methods

#### getResponsiveHeaders
- **File**: `internal/table.go:115`
- **Signature**: `func (t *Table) getResponsiveHeaders() ([]string, []int)`
- **Description**: Returns headers and column indices that fit in terminal
- **Usage**: Called from buildTable to determine visible columns

#### calculateEstimatedWidth
- **File**: `internal/table.go:164`
- **Signature**: `func (t *Table) calculateEstimatedWidth() int`
- **Description**: Estimates total width needed for all columns
- **Usage**: Called from getResponsiveHeaders to check if all columns fit

#### calculateEstimatedWidthForHeaders
- **File**: `internal/table.go:195`
- **Signature**: `func (t *Table) calculateEstimatedWidthForHeaders(headers []string) int`
- **Description**: Estimates width for specific headers
- **Usage**: Called from getResponsiveHeaders to check partial column fits

#### getTerminalWidth
- **File**: `internal/info_renderer.go:31`
- **Signature**: `func getTerminalWidth() int`
- **Description**: Returns terminal width with multiple fallbacks (term.GetSize, COLUMNS env, tput cols)
- **Usage**: Called from table and info renderer for responsive layout

#### NewInfoRenderer
- **File**: `internal/info_renderer.go:58`
- **Signature**: `func NewInfoRenderer() *InfoRenderer`
- **Description**: Creates new info renderer with adaptive colors and terminal width
- **Usage**: Called from info command to render detailed worktree information

#### RenderWorktreeInfo
- **File**: `internal/info_renderer.go:125`
- **Signature**: `func (r *InfoRenderer) RenderWorktreeInfo(data *WorktreeInfoData) string`
- **Description**: Renders complete worktree info with sections for worktree, JIRA, and git
- **Usage**: Called from info command to display worktree details

#### renderWorktreeSection
- **File**: `internal/info_renderer.go:151`
- **Signature**: `func (r *InfoRenderer) renderWorktreeSection(data *WorktreeInfoData) string`
- **Description**: Renders worktree section with name, path, branch, creation time, and status
- **Usage**: Called from RenderWorktreeInfo

#### renderJiraSection
- **File**: `internal/info_renderer.go:175`
- **Signature**: `func (r *InfoRenderer) renderJiraSection(jira *JiraTicketDetails) string`
- **Description**: Renders JIRA section with ticket details and comments
- **Usage**: Called from RenderWorktreeInfo when JIRA data is available

#### renderGitSection
- **File**: `internal/info_renderer.go:245`
- **Signature**: `func (r *InfoRenderer) renderGitSection(data *WorktreeInfoData) string`
- **Description**: Renders git section with branch info, commits, and modified files
- **Usage**: Called from RenderWorktreeInfo

#### renderKeyValue
- **File**: `internal/info_renderer.go:330`
- **Signature**: `func (r *InfoRenderer) renderKeyValue(key, value string) string`
- **Description**: Renders key-value pair with consistent formatting
- **Usage**: Called from all render section methods

#### formatGitStatus
- **File**: `internal/info_renderer.go:338`
- **Signature**: `func (r *InfoRenderer) formatGitStatus(status *GitStatus) string`
- **Description**: Formats git status for info display (different from styles.go version)
- **Usage**: Called from renderWorktreeSection

#### formatDuration
- **File**: `internal/info_renderer.go:355`
- **Signature**: `func (r *InfoRenderer) formatDuration(d time.Duration) string`
- **Description**: Formats duration as human-readable string (minutes/hours/days)
- **Usage**: Called from render methods to show time ago

#### getStatusIcon
- **File**: `internal/info_renderer.go:365`
- **Signature**: `func (r *InfoRenderer) getStatusIcon(status string) string`
- **Description**: Returns file status icon (A/M/D/?)
- **Usage**: Called from renderGitSection for file status

#### formatPriority
- **File**: `internal/info_renderer.go:378`
- **Signature**: `func (r *InfoRenderer) formatPriority(priority string) string`
- **Description**: Formats JIRA priority with colored icons
- **Usage**: Called from renderJiraSection

#### wrapText
- **File**: `internal/info_renderer.go:396`
- **Signature**: `func (r *InfoRenderer) wrapText(text string, width int) string`
- **Description**: Wraps text to fit within specified width
- **Usage**: Called from render methods for long text content

### Potential Duplicates/Redundancy

#### Pattern 1: Terminal Width Detection
- **Functions**: `getTerminalWidth` (info_renderer.go), terminal width detection in `NewTable`
- **Similarity**: Both detect terminal width with similar fallback strategies
- **Suggestion**: Consolidate into single utility function

#### Pattern 2: Git Status Formatting
- **Functions**: `FormatGitStatus` (styles.go), `formatGitStatus` (info_renderer.go)
- **Similarity**: Both format git status but with different output styles
- **Suggestion**: Could share common logic while maintaining different presentation

#### Pattern 3: Icon Manager Access
- **Functions**: Multiple icon getter methods in IconManager (Success, Warning, Error, etc.)
- **Similarity**: All simple getters that return config icon values
- **Suggestion**: Could use reflection or map-based approach to reduce boilerplate

#### Pattern 4: Format Helper Functions
- **Functions**: `FormatSuccess`, `FormatWarning`, `FormatError`, `FormatInfo`
- **Similarity**: All follow same pattern of calling FormatStatusIcon with appropriate icon
- **Suggestion**: Could generate these or use single parameterized function

#### Pattern 5: Text Styling Functions
- **Functions**: `FormatVerbose`, `FormatHeader`, `FormatSubHeader`, `FormatBold`, `FormatSubtle`, `FormatPrompt`
- **Similarity**: All apply lipgloss style to text and return formatted string
- **Suggestion**: Could use map of style names to styles

#### Pattern 6: Width Calculation
- **Functions**: `calculateEstimatedWidth`, `calculateEstimatedWidthForHeaders`
- **Similarity**: Both calculate table width with similar logic
- **Suggestion**: Could extract common width calculation logic

#### Pattern 7: Responsive Layout Logic
- **Functions**: `getResponsiveHeaders`, `buildTable`, responsive calculations in info renderer
- **Similarity**: All deal with adapting content to terminal width
- **Suggestion**: Could extract common responsive layout utilities

#### Pattern 8: Column Configuration
- **Functions**: `NewTable`, `NewTestTable`
- **Similarity**: Both set up same minimum column widths map
- **Suggestion**: Could extract column configuration to constants or shared function

# 004_internal_testutils

## Directory: internal/testutils/

### Functions Found

#### NewGitTestRepo
- **File**: `internal/testutils/git_harness.go:34`
- **Signature**: `func NewGitTestRepo(t *testing.T, opts ...RepoOption) *GitTestRepo`
- **Description**: Creates new test repository with options, sets up bare remote, local repo, git config, and initial commit
- **Usage**: Called from test scenarios and test files to create test repositories

#### setupBareRemote
- **File**: `internal/testutils/git_harness.go:68`
- **Signature**: `func (r *GitTestRepo) setupBareRemote() error`
- **Description**: Creates bare remote repository and sets default branch
- **Usage**: Called internally from NewGitTestRepo

#### setupLocalRepo
- **File**: `internal/testutils/git_harness.go:89`
- **Signature**: `func (r *GitTestRepo) setupLocalRepo() error`
- **Description**: Clones remote repository to create local working directory
- **Usage**: Called internally from NewGitTestRepo

#### configureGitUser
- **File**: `internal/testutils/git_harness.go:98`
- **Signature**: `func (r *GitTestRepo) configureGitUser() error`
- **Description**: Sets git user name and email in local repository
- **Usage**: Called internally from NewGitTestRepo

#### createInitialCommit
- **File**: `internal/testutils/git_harness.go:110`
- **Signature**: `func (r *GitTestRepo) createInitialCommit() error`
- **Description**: Creates README.md, commits it, renames branch, and pushes to remote
- **Usage**: Called internally from NewGitTestRepo

#### runGitCommand
- **File**: `internal/testutils/git_harness.go:135`
- **Signature**: `func (r *GitTestRepo) runGitCommand(args ...string) error`
- **Description**: Executes git command in local directory with error handling and output logging
- **Usage**: Called throughout GitTestRepo methods for git operations

#### runCommand
- **File**: `internal/testutils/git_harness.go:146`
- **Signature**: `func (r *GitTestRepo) runCommand(name string, args ...string) ([]byte, error)`
- **Description**: Executes arbitrary command in local directory and returns output
- **Usage**: Called from ListBranches and other methods needing command output

#### Cleanup
- **File**: `internal/testutils/git_harness.go:152`
- **Signature**: `func (r *GitTestRepo) Cleanup() void`
- **Description**: Empty cleanup method (cleanup handled by t.TempDir())
- **Usage**: Called from tests for cleanup (currently no-op)

#### GetLocalPath
- **File**: `internal/testutils/git_harness.go:155`
- **Signature**: `func (r *GitTestRepo) GetLocalPath() string`
- **Description**: Returns local repository directory path
- **Usage**: Called from tests to get local repo path

#### GetRemotePath
- **File**: `internal/testutils/git_harness.go:159`
- **Signature**: `func (r *GitTestRepo) GetRemotePath() string`
- **Description**: Returns remote repository directory path
- **Usage**: Called from tests to get remote repo path

#### CreateBranch
- **File**: `internal/testutils/git_harness.go:163`
- **Signature**: `func (r *GitTestRepo) CreateBranch(name, content string) error`
- **Description**: Creates new branch with content file, commits, pushes, and returns to default branch
- **Usage**: Called from test scenarios to create branches

#### CreateBranchFrom
- **File**: `internal/testutils/git_harness.go:192`
- **Signature**: `func (r *GitTestRepo) CreateBranchFrom(name, baseBranch, content string) error`
- **Description**: Creates new branch from specific base branch with content
- **Usage**: Called from tests needing branches from specific base

#### SwitchToBranch
- **File**: `internal/testutils/git_harness.go:225`
- **Signature**: `func (r *GitTestRepo) SwitchToBranch(name string) error`
- **Description**: Switches to specified branch using git checkout
- **Usage**: Called from tests to switch between branches

#### WriteFile
- **File**: `internal/testutils/git_harness.go:232`
- **Signature**: `func (r *GitTestRepo) WriteFile(path, content string) error`
- **Description**: Writes file to repository with directory creation
- **Usage**: Called from tests and other methods to create files

#### CreateEnvrc
- **File**: `internal/testutils/git_harness.go:245`
- **Signature**: `func (r *GitTestRepo) CreateEnvrc(mapping map[string]string) error`
- **Description**: Creates .envrc file with deterministic ordering (MAIN, PREVIEW, PROD, then alphabetical)
- **Usage**: Called from test scenarios to create .envrc files

#### CommitChanges
- **File**: `internal/testutils/git_harness.go:291`
- **Signature**: `func (r *GitTestRepo) CommitChanges(message string) error`
- **Description**: Adds all changes and commits with message
- **Usage**: Called from tests to commit changes

#### CommitChangesWithForceAdd
- **File**: `internal/testutils/git_harness.go:303`
- **Signature**: `func (r *GitTestRepo) CommitChangesWithForceAdd(message string) error`
- **Description**: Adds all changes with force flag and commits
- **Usage**: Called from scenarios to commit .envrc (which might be in .gitignore)

#### PushBranch
- **File**: `internal/testutils/git_harness.go:315`
- **Signature**: `func (r *GitTestRepo) PushBranch(branch string) error`
- **Description**: Pushes specified branch to remote
- **Usage**: Called from tests to push branches to remote

#### ConvertToBare
- **File**: `internal/testutils/git_harness.go:322`
- **Signature**: `func (r *GitTestRepo) ConvertToBare() string`
- **Description**: Returns remote directory path (for bare repository access)
- **Usage**: Called from tests needing bare repository path

#### ListBranches
- **File**: `internal/testutils/git_harness.go:326`
- **Signature**: `func (r *GitTestRepo) ListBranches() ([]string, error)`
- **Description**: Lists remote branches excluding HEAD
- **Usage**: Called from tests to verify branch existence

#### WithWorkingDir
- **File**: `internal/testutils/git_harness.go:344`
- **Signature**: `func (r *GitTestRepo) WithWorkingDir(dir string) func()`
- **Description**: Changes working directory and returns restore function
- **Usage**: Called from InLocalRepo to temporarily change directory

#### InLocalRepo
- **File**: `internal/testutils/git_harness.go:352`
- **Signature**: `func (r *GitTestRepo) InLocalRepo(fn func() error) error`
- **Description**: Executes function in local repository directory
- **Usage**: Called from tests to run operations in local repo

#### CreateSynchronizedBranch
- **File**: `internal/testutils/git_harness.go:358`
- **Signature**: `func (r *GitTestRepo) CreateSynchronizedBranch(name string) error`
- **Description**: Creates branch and immediately pushes to remote
- **Usage**: Called from tests needing synchronized branches

#### NewBasicRepo
- **File**: `internal/testutils/scenarios.go:8`
- **Signature**: `func NewBasicRepo(t *testing.T) *GitTestRepo`
- **Description**: Creates basic repository with main branch and test user
- **Usage**: Called from tests and other scenarios as foundation

#### NewMultiBranchRepo
- **File**: `internal/testutils/scenarios.go:15`
- **Signature**: `func NewMultiBranchRepo(t *testing.T) *GitTestRepo`
- **Description**: Creates repository with multiple branches (develop, feature/auth, production/v1.0)
- **Usage**: Called from tests and other scenarios needing multiple branches

#### NewEnvrcRepo
- **File**: `internal/testutils/scenarios.go:33`
- **Signature**: `func NewEnvrcRepo(t *testing.T, mapping map[string]string) *GitTestRepo`
- **Description**: Creates multi-branch repository with .envrc file
- **Usage**: Called from scenarios needing .envrc configuration

#### NewStandardEnvrcRepo
- **File**: `internal/testutils/scenarios.go:51`
- **Signature**: `func NewStandardEnvrcRepo(t *testing.T) *GitTestRepo`
- **Description**: Creates repository with standard .envrc mapping (MAIN, DEV, FEAT, PROD)
- **Usage**: Called from tests needing standard configuration

#### NewRepoWithConflictingBranches
- **File**: `internal/testutils/scenarios.go:62`
- **Signature**: `func NewRepoWithConflictingBranches(t *testing.T) *GitTestRepo`
- **Description**: Creates repository with conflicting file content between branches
- **Usage**: Called from tests needing merge conflicts

#### NewLargeHistoryRepo
- **File**: `internal/testutils/scenarios.go:100`
- **Signature**: `func NewLargeHistoryRepo(t *testing.T) *GitTestRepo`
- **Description**: Creates repository with 10 commits of history
- **Usage**: Called from tests needing commit history

#### NewEmptyRepo
- **File**: `internal/testutils/scenarios.go:120`
- **Signature**: `func NewEmptyRepo(t *testing.T) *GitTestRepo`
- **Description**: Creates empty repository with just initial setup
- **Usage**: Called from tests needing minimal repository

#### WithDefaultBranch
- **File**: `internal/testutils/repo_options.go:5`
- **Signature**: `func WithDefaultBranch(branch string) RepoOption`
- **Description**: Option to set default branch name
- **Usage**: Called from tests to configure default branch

#### WithUser
- **File**: `internal/testutils/repo_options.go:11`
- **Signature**: `func WithUser(name, email string) RepoOption`
- **Description**: Option to set git user name and email
- **Usage**: Called from tests to configure git user

#### WithRemoteName
- **File**: `internal/testutils/repo_options.go:18`
- **Signature**: `func WithRemoteName(name string) RepoOption`
- **Description**: Option to set remote name (default: origin)
- **Usage**: Called from tests to configure remote name

#### NewMockJiraCLI
- **File**: `internal/testutils/mock_services.go:26`
- **Signature**: `func NewMockJiraCLI(t *testing.T) *MockJiraCLI`
- **Description**: Creates mock JIRA CLI for testing
- **Usage**: Called from tests needing JIRA CLI simulation

#### SetResponse
- **File**: `internal/testutils/mock_services.go:33`
- **Signature**: `func (m *MockJiraCLI) SetResponse(command, output string) void`
- **Description**: Sets successful response for JIRA command
- **Usage**: Called from tests to configure JIRA responses

#### AddResponse
- **File**: `internal/testutils/mock_services.go:41`
- **Signature**: `func (m *MockJiraCLI) AddResponse(command, output string, exitCode int) void`
- **Description**: Adds response with specific exit code for JIRA command
- **Usage**: Called from tests to configure JIRA responses with exit codes

#### SimulateFailure
- **File**: `internal/testutils/mock_services.go:49`
- **Signature**: `func (m *MockJiraCLI) SimulateFailure(command string, errorMsg string) void`
- **Description**: Simulates JIRA command failure with error message
- **Usage**: Called from tests to simulate JIRA failures

#### InstallMock
- **File**: `internal/testutils/mock_services.go:57`
- **Signature**: `func (m *MockJiraCLI) InstallMock() error`
- **Description**: Installs mock JIRA CLI in PATH by creating temp script
- **Usage**: Called from tests to activate JIRA mocking

#### RemoveMock
- **File**: `internal/testutils/mock_services.go:85`
- **Signature**: `func (m *MockJiraCLI) RemoveMock() error`
- **Description**: Removes mock JIRA CLI from filesystem
- **Usage**: Called from tests to cleanup JIRA mocking

#### generateMockScript
- **File**: `internal/testutils/mock_services.go:95`
- **Signature**: `func (m *MockJiraCLI) generateMockScript() string`
- **Description**: Generates platform-specific mock script
- **Usage**: Called internally from InstallMock

#### generateUnixScript
- **File**: `internal/testutils/mock_services.go:102`
- **Signature**: `func (m *MockJiraCLI) generateUnixScript() string`
- **Description**: Generates Unix/Linux bash script for JIRA mocking
- **Usage**: Called from generateMockScript for Unix platforms

#### generateWindowsScript
- **File**: `internal/testutils/mock_services.go:119`
- **Signature**: `func (m *MockJiraCLI) generateWindowsScript() string`
- **Description**: Generates Windows batch script for JIRA mocking
- **Usage**: Called from generateMockScript for Windows platform

#### IsCommandAvailable
- **File**: `internal/testutils/mock_services.go:136`
- **Signature**: `func IsCommandAvailable(command string) bool`
- **Description**: Checks if command is available in PATH
- **Usage**: Called from tests to check command availability

#### MockCommandUnavailable
- **File**: `internal/testutils/mock_services.go:141`
- **Signature**: `func MockCommandUnavailable(t *testing.T, command string) void`
- **Description**: Modifies PATH to make command unavailable for testing
- **Usage**: Called from tests to simulate missing commands

### Potential Duplicates/Redundancy

#### Pattern 1: Git Command Execution
- **Functions**: `runGitCommand`, `runCommand`
- **Similarity**: Both execute commands in local directory with error handling
- **Suggestion**: Could consolidate into single command execution function

#### Pattern 2: Branch Creation Operations
- **Functions**: `CreateBranch`, `CreateBranchFrom`, `CreateSynchronizedBranch`
- **Similarity**: All create branches with similar patterns (checkout, commit, push, return to default)
- **Suggestion**: Could extract common branch creation logic

#### Pattern 3: Repository Setup Methods
- **Functions**: `setupBareRemote`, `setupLocalRepo`, `configureGitUser`, `createInitialCommit`
- **Similarity**: All setup methods called in sequence from NewGitTestRepo
- **Suggestion**: Could consolidate setup logic or use builder pattern

#### Pattern 4: Scenario Factory Functions
- **Functions**: `NewBasicRepo`, `NewMultiBranchRepo`, `NewEnvrcRepo`, `NewStandardEnvrcRepo`, `NewRepoWithConflictingBranches`, `NewLargeHistoryRepo`, `NewEmptyRepo`
- **Similarity**: All create repositories with different configurations
- **Suggestion**: Could use builder pattern or configuration structs

#### Pattern 5: File Operations
- **Functions**: `WriteFile`, `CreateEnvrc`, `CommitChanges`, `CommitChangesWithForceAdd`
- **Similarity**: All involve file creation and committing
- **Suggestion**: Could extract common file handling patterns

#### Pattern 6: Mock Script Generation
- **Functions**: `generateUnixScript`, `generateWindowsScript`
- **Similarity**: Both generate scripts but for different platforms
- **Suggestion**: Could use template-based approach to reduce duplication

#### Pattern 7: Repository Options
- **Functions**: `WithDefaultBranch`, `WithUser`, `WithRemoteName`
- **Similarity**: All follow same option pattern for configuring repository
- **Suggestion**: Could generate options or use struct embedding

#### Pattern 8: Mock Response Management
- **Functions**: `SetResponse`, `AddResponse`, `SimulateFailure`
- **Similarity**: All manipulate the responses map with similar patterns
- **Suggestion**: Could consolidate response management

#### Pattern 9: Directory Path Methods
- **Functions**: `GetLocalPath`, `GetRemotePath`, `ConvertToBare`
- **Similarity**: All return directory paths from repository
- **Suggestion**: Could use property-based access or unified path provider

# 005_internal_mergeback

## Directory: internal/

### Functions Found

#### CheckMergeBackStatus
- **File**: `internal/mergeback.go:44`
- **Signature**: `func CheckMergeBackStatus(configPath string) (*MergeBackStatus, error)`
- **Description**: Main function to check for merge-back requirements across environment branches, returns status with needed merge-backs and user commits
- **Usage**: Called from cmd/root.go to check and display merge-back alerts

#### parseEnvrcFile
- **File**: `internal/mergeback.go:127`
- **Signature**: `func parseEnvrcFile(configPath string) ([]EnvVarMapping, error)`
- **Description**: Parses .envrc file to extract environment variable mappings with order preservation
- **Usage**: Called from CheckMergeBackStatus to read environment configuration

#### getUserInfo
- **File**: `internal/mergeback.go:173`
- **Signature**: `func getUserInfo(repoPath string) (string, string, error)`
- **Description**: Gets git user email and name from repository configuration
- **Usage**: Called from CheckMergeBackStatus to identify user commits

#### getCommitsNeedingMergeBack
- **File**: `internal/mergeback.go:197`
- **Signature**: `func getCommitsNeedingMergeBack(repoPath, targetBranch, sourceBranch string) ([]MergeBackCommitInfo, error)`
- **Description**: Uses git log to find commits that exist in source branch but not in target branch
- **Usage**: Called from CheckMergeBackStatus to find commits needing merge-back

#### isUserCommit
- **File**: `internal/mergeback.go:241`
- **Signature**: `func isUserCommit(commit MergeBackCommitInfo, userEmail, userName string) bool`
- **Description**: Determines if a commit was made by the current user based on email or name match
- **Usage**: Called from CheckMergeBackStatus to identify user's commits

#### FormatMergeBackAlert
- **File**: `internal/mergeback.go:251`
- **Signature**: `func FormatMergeBackAlert(status *MergeBackStatus) string`
- **Description**: Formats merge-back status into user-friendly alert message with commit details
- **Usage**: Called from cmd/root.go to display merge-back alerts

#### formatRelativeTime
- **File**: `internal/mergeback.go:281`
- **Signature**: `func formatRelativeTime(t time.Time) string`
- **Description**: Formats time duration as human-readable relative time (minutes/hours/days ago)
- **Usage**: Called from FormatMergeBackAlert to show commit timestamps

### Potential Duplicates/Redundancy

#### Pattern 1: Git Command Execution
- **Functions**: `getUserInfo`, `getCommitsNeedingMergeBack`
- **Similarity**: Both execute git commands using exec.Command with similar error handling patterns
- **Suggestion**: Could use the existing `execCommand` function from git.go for consistency

#### Pattern 2: File Parsing with Regex
- **Functions**: `parseEnvrcFile` uses regex similar to `ParseEnvrc` in config.go
- **Similarity**: Both parse .envrc files with environment variable regex patterns
- **Suggestion**: Could consolidate .envrc parsing logic or share regex patterns

#### Pattern 3: Time Formatting
- **Functions**: `formatRelativeTime` similar to `formatDuration` in info_renderer.go
- **Similarity**: Both format time durations into human-readable strings
- **Suggestion**: Could extract common time formatting utilities

#### Pattern 4: User Identification
- **Functions**: `getUserInfo`, `isUserCommit`
- **Similarity**: Both deal with git user identification and comparison
- **Suggestion**: Could extract user identification logic to shared utility

#### Pattern 5: Git Repository Operations
- **Functions**: `CheckMergeBackStatus` creates GitManager similarly to other functions
- **Similarity**: Same pattern of FindGitRoot -> NewGitManager as seen in other files
- **Suggestion**: Could use the existing manager creation pattern from other cmd functions

#### Pattern 6: Error Handling with Nil Returns
- **Functions**: `CheckMergeBackStatus` returns nil on various errors
- **Similarity**: Same pattern of returning nil status on errors as seen in other functions
- **Suggestion**: Could use consistent error handling patterns

#### Pattern 7: String Building for Output
- **Functions**: `FormatMergeBackAlert` uses strings.Builder similar to other formatting functions
- **Similarity**: Same pattern of building formatted output strings
- **Suggestion**: Could extract common string building patterns

#### Pattern 8: Branch Existence Checking
- **Functions**: `CheckMergeBackStatus` checks branch existence using gitManager.BranchExists
- **Similarity**: Same pattern as other functions that validate branch existence
- **Suggestion**: Could use existing branch validation utilities

# 006_main_entry

## Directory: /

### Functions Found

#### main
- **File**: `main.go:12`
- **Signature**: `func main()`
- **Description**: Application entry point that calls cmd.Execute() and handles errors with proper cleanup
- **Usage**: Called by Go runtime when application starts

### Analysis

The main.go file is extremely minimal and follows Go best practices for CLI applications:

1. **Single Responsibility**: The main function has only one job - to start the application
2. **Proper Error Handling**: Uses os.Exit(1) for error conditions
3. **Resource Cleanup**: Uses defer for log file cleanup with cmd.CloseLogFile()
4. **Delegation Pattern**: All functionality is delegated to the cmd package through cmd.Execute()
5. **Error Reporting**: Uses cmd.PrintError() for consistent error formatting

### Key Characteristics

- **Minimal Implementation**: Only 9 lines of actual code
- **Clean Error Path**: Handles cmd.Execute() errors gracefully
- **Proper Exit Codes**: Uses os.Exit(1) for failures, implicit 0 for success
- **Resource Management**: Ensures log file is closed even on errors
- **Consistent Styling**: Uses the same error formatting as rest of application

### Usage Patterns

The main function is referenced extensively in test files, particularly in cmd package tests where cmd.Execute() is called directly. This shows that the application follows the standard Go CLI pattern where main() is just a thin wrapper around the actual command execution logic.

### Potential Duplicates/Redundancy

#### Pattern 1: Entry Point Pattern
- **Functions**: `main()` in main.go
- **Similarity**: This is the standard Go application entry point pattern
- **Suggestion**: No changes needed - this follows Go best practices

#### Pattern 2: Error Handling with Exit
- **Functions**: `main()` uses os.Exit(1) for errors
- **Similarity**: Standard pattern for CLI applications
- **Suggestion**: No changes needed - this is idiomatic Go

#### Pattern 3: Defer for Cleanup
- **Functions**: `main()` uses defer for log file cleanup
- **Similarity**: Standard Go pattern for resource cleanup
- **Suggestion**: No changes needed - this is proper Go resource management

# 007_cross_package

## Cross-Package Analysis

### Overview

This analysis examines duplicate patterns and consolidation opportunities that span multiple packages and directories in the codebase. After analyzing jobs 001-006, several recurring patterns emerge that cross package boundaries and present opportunities for shared utilities and standardization.

### Key Cross-Package Duplicate Patterns

#### Pattern 1: Git Command Execution
- **Locations**: 
  - `internal/git.go` - `execCommand`, `execCommandRun`
  - `internal/mergeback.go` - `getUserInfo`, `getCommitsNeedingMergeBack`
  - `internal/testutils/git_harness.go` - `runGitCommand`, `runCommand`
- **Similarity**: All execute git commands with similar error handling, output capture, and directory context patterns
- **Suggestion**: Create shared `pkg/gitutils` package with standardized command execution functions

#### Pattern 2: Manager Creation Pattern
- **Locations**: Nearly all cmd files follow identical sequence:
  - `FindGitRoot` -> `NewGitManager` -> `NewManager` -> `LoadEnvMapping`
- **Files**: `cmd/add.go`, `cmd/clone.go`, `cmd/info.go`, `cmd/list.go`, `cmd/pull.go`, `cmd/push.go`, `cmd/remove.go`, `cmd/switch.go`, `cmd/sync.go`, `cmd/validate.go`
- **Similarity**: Same initialization sequence repeated across all commands
- **Suggestion**: Extract to `createInitializedManager()` helper function in cmd package

#### Pattern 3: Error Handling and Wrapping
- **Locations**: 
  - All cmd files use `fmt.Errorf("failed to X: %w", err)` pattern
  - internal/git.go, internal/manager.go, internal/mergeback.go use similar wrapping
- **Similarity**: Consistent error wrapping patterns across all packages
- **Suggestion**: Create shared error handling utilities with standardized messages

#### Pattern 4: Terminal Width Detection
- **Locations**:
  - `internal/table.go` - `NewTable` constructor
  - `internal/info_renderer.go` - `getTerminalWidth` function
- **Similarity**: Both detect terminal width with similar fallback strategies (term.GetSize, COLUMNS env, tput cols)
- **Suggestion**: Extract to shared `pkg/terminal` utility package

#### Pattern 5: Time/Duration Formatting
- **Locations**:
  - `internal/info_renderer.go` - `formatDuration` 
  - `internal/mergeback.go` - `formatRelativeTime`
- **Similarity**: Both format time durations as human-readable strings (minutes/hours/days ago)
- **Suggestion**: Create shared time formatting utilities in `pkg/timeutils`

#### Pattern 6: Branch Existence Validation
- **Locations**:
  - `internal/git.go` - `BranchExists` method
  - `internal/manager.go` - `BranchExists` delegation
  - `internal/manager.go` - `ValidateEnvrc` validation
  - Multiple cmd files validate branch existence before operations
- **Similarity**: Same validation patterns scattered across multiple files
- **Suggestion**: Centralize branch validation in git manager with consistent error messages

#### Pattern 7: File Path Construction
- **Locations**:
  - `internal/manager.go` - `GetWorktreePath` uses `filepath.Join(repoPath, worktreeDir, name)`
  - `internal/git.go` - Multiple functions construct worktree paths
  - cmd files construct paths for operations
- **Similarity**: Repeated patterns for constructing worktree paths
- **Suggestion**: Create path construction utilities in shared package

#### Pattern 8: Configuration Parsing (.envrc)
- **Locations**:
  - `internal/config.go` - `ParseEnvrc` function
  - `internal/mergeback.go` - `parseEnvrcFile` function
- **Similarity**: Both parse .envrc files with similar regex patterns but different return types
- **Suggestion**: Consolidate into single .envrc parsing utility with consistent interface

#### Pattern 9: Git Status Formatting
- **Locations**:
  - `internal/styles.go` - `FormatGitStatus` function
  - `internal/info_renderer.go` - `formatGitStatus` method
- **Similarity**: Both format git status but with different output styles
- **Suggestion**: Share common status logic while maintaining different presentation styles

#### Pattern 10: Command Completion
- **Locations**:
  - `cmd/add.go`, `cmd/pull.go`, `cmd/push.go`, `cmd/remove.go`, `cmd/switch.go` - ValidArgsFunction implementations
  - All use `getWorktreeNames()` for completion
- **Similarity**: Repeated completion patterns across commands
- **Suggestion**: Extract shared completion utilities

### Consolidation Opportunities

#### 1. Shared Utilities Package Structure
```
pkg/
├── gitutils/          # Git command execution, branch validation
├── terminal/          # Terminal width detection, responsive layout
├── timeutils/         # Time formatting, duration display
├── pathutils/         # Path construction, validation
├── configutils/       # .envrc parsing, configuration handling
└── completion/        # Command completion utilities
```

#### 2. Command Factory Pattern
```go
// Extract common manager creation
func createInitializedManager() (*internal.Manager, error) {
    gitRoot, err := internal.FindGitRoot(".")
    if err != nil {
        return nil, fmt.Errorf("failed to find git root: %w", err)
    }
    
    manager, err := internal.NewManager(gitRoot)
    if err != nil {
        return nil, fmt.Errorf("failed to create manager: %w", err)
    }
    
    if err := manager.LoadEnvMapping(cmd.GetConfigPath()); err != nil {
        return nil, fmt.Errorf("failed to load env mapping: %w", err)
    }
    
    return manager, nil
}
```

#### 3. Standardized Error Handling
```go
// Create consistent error wrapping
func WrapError(operation string, err error) error {
    return fmt.Errorf("failed to %s: %w", operation, err)
}

func WrapErrorf(operation string, err error, format string, args ...interface{}) error {
    context := fmt.Sprintf(format, args...)
    return fmt.Errorf("failed to %s (%s): %w", operation, context, err)
}
```

#### 4. Unified Configuration Parsing
```go
// Consolidate .envrc parsing
type EnvrcParser interface {
    ParseFile(path string) (map[string]string, error)
    ParseMappings(path string) ([]EnvVarMapping, error)
}
```

#### 5. Shared Testing Utilities
- Consolidate common test patterns from testutils
- Extract shared git command execution for tests
- Standardize repository scenario creation

### Implementation Priority

1. **High Priority**: Manager creation factory (used in 10+ files)
2. **High Priority**: Git command execution utilities (used in 5+ files)
3. **Medium Priority**: Terminal width detection (used in 2 files but frequently called)
4. **Medium Priority**: Time formatting utilities (used in 2 files)
5. **Low Priority**: Configuration parsing consolidation (used in 2 files)

### Benefits

1. **Reduced Code Duplication**: Eliminate 40+ duplicate function implementations
2. **Consistent Behavior**: Standardize error handling, formatting, and validation
3. **Easier Maintenance**: Single location for common functionality updates
4. **Better Testing**: Shared utilities can be thoroughly tested once
5. **Improved Reliability**: Consistent implementations reduce bugs

### Migration Strategy

1. Create pkg/ directory structure
2. Extract high-priority utilities first
3. Update imports incrementally
4. Add comprehensive tests for shared utilities
5. Remove duplicate implementations
6. Update documentation

---