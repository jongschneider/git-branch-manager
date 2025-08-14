## Analysis: cmd/hotfix.go Interface Extraction

### Current Implementation Analysis

**Functions in cmd/hotfix.go:**
- `newHotfixCommand()`: Main command constructor with completion
- `findProductionBranch(manager *internal.Manager)`: Production branch detection via mergeback chain analysis
- `isProductionBranchName(branchName string)`: Helper for production branch detection
- `buildDeploymentChain()` & `buildMergeChain()`: Deployment topology analysis
- `findMergeIntoTarget()`: Merge configuration traversal
- `generateHotfixBranchName`: Function variable for branch name generation

**Direct Manager Dependencies:**
- `manager.GetConfig().Settings.HotfixPrefix`
- `manager.GetGBMConfig()` 
- `manager.AddWorktree()`
- `manager.GetGitManager().GetDefaultBranch()`

**External Dependencies:**
- JIRA CLI integration via `internal.GetJiraIssues()` and `internal.GenerateBranchFromJira()`
- Git operations via `internal.FindGitRoot()`
- Configuration parsing via `internal.ParseGBMConfig()`
- File system operations for config file detection

**Production Branch Detection Logic:**
1. Analyzes `gbm.branchconfig.yaml` mergeback chains
2. Finds branches with `merge_into` targets but no incoming merges (production branches)
3. Falls back to root branches or git default branch

### Established Interface Patterns

**Interface Structure:**
- Defined in command files with `//go:generate` mock generation
- Follow `worktree<Operation>` naming convention
- Include only methods needed by the specific command

**Manager Wrapper Pattern:**
- Manager implements interfaces via wrapper methods
- Wrappers delegate to GitManager or internal logic
- Enable clean separation for testing

**Handler Function Pattern:**
- Commands refactored to use `handle<Operation>(interface, ...)` functions
- Interfaces passed as parameters for dependency injection
- Enables unit testing with mocks

**Testing Approach:**
- Interface-based unit tests with mocks in cmd package
- Integration tests moved to internal package
- Fast execution through mock isolation

### Recommended hotfixCreator Interface

Based on analysis, the interface should include:
- `AddWorktree(worktreeName, branchName string, createBranch bool, baseBranch string) error`
- `GetConfig() *Config` (for HotfixPrefix)
- `GetGBMConfig() *GBMConfig` (for deployment chain analysis)
- `GetDefaultBranch() (string, error)` (wrapper for GitManager)
- `FindProductionBranch() (string, error)` (extract to Manager method)
- `GetJiraIssues() ([]JiraIssue, error)` (for completion)
- `GenerateBranchFromJira(jiraKey string) (string, error)` (for branch naming)

The production branch detection logic should be extracted to a Manager method to support interface testing while preserving the sophisticated deployment chain analysis.