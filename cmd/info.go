package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gbm/internal"

	"github.com/spf13/cobra"
)

//go:generate go run github.com/matryer/moq@latest -out ./autogen_worktreeInfoProvider.go . worktreeInfoProvider

// worktreeInfoProvider interface abstracts the Manager operations needed for getting worktree info
type worktreeInfoProvider interface {
	// Core worktree operations
	GetWorktrees() ([]*internal.WorktreeInfo, error)
	GetWorktreeStatus(worktreePath string) (*internal.GitStatus, error)

	// Configuration and state access
	GetConfig() *internal.Config
	GetState() *internal.State

	// Wrapper methods for GitManager operations
	GetWorktreeCommitHistory(worktreePath string, limit int) ([]internal.CommitInfo, error)
	GetWorktreeFileChanges(worktreePath string) ([]internal.FileChange, error)
	GetWorktreeCurrentBranch(worktreePath string) (string, error)
	GetWorktreeUpstreamBranch(worktreePath string) (string, error)
	GetWorktreeAheadBehindCount(worktreePath string) (int, int, error)
	VerifyWorktreeRef(ref string, worktreePath string) (bool, error)

	// JIRA integration
	GetJiraTicketDetails(jiraKey string) (*internal.JiraTicketDetails, error)
}

func newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <worktree-name>",
		Short: "Display detailed information about a worktree",
		Long: `Display comprehensive information about a specific worktree including:
- Worktree metadata (name, path, branch, creation date)
- Git status and branch information
- JIRA ticket details (if the worktree name matches a JIRA key)
- Recent commits and modified files`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfoCommand(args[0])
		},
	}

	return cmd
}

func runInfoCommand(worktreeName string) error {
	// Handle current directory reference
	if worktreeName == "." {
		currentPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		worktreeName = filepath.Base(currentPath)
	}

	manager, err := createInitializedManager()
	if err != nil {
		if !errors.Is(err, ErrLoadGBMConfig) {
			return err
		}

		PrintVerbose("%v", err)
	}

	// Get worktree information
	worktreeInfo, err := getWorktreeInfo(manager, worktreeName)
	if err != nil {
		return fmt.Errorf("failed to get worktree info: %w", err)
	}

	// Display the information
	displayWorktreeInfo(worktreeInfo, manager.GetConfig())

	return nil
}

func getWorktreeInfo(provider worktreeInfoProvider, worktreeName string) (*internal.WorktreeInfoData, error) {
	// Get all worktrees
	worktrees, err := provider.GetWorktrees()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktrees: %w", err)
	}

	// Find the specific worktree
	var targetWorktree *internal.WorktreeInfo
	for _, wt := range worktrees {
		if wt.Name == worktreeName {
			targetWorktree = wt
			break
		}
	}

	if targetWorktree == nil {
		return nil, fmt.Errorf("worktree '%s' not found", worktreeName)
	}

	// Get git status for the worktree
	gitStatus, err := provider.GetWorktreeStatus(targetWorktree.Path)
	if err != nil {
		PrintVerbose("Failed to get git status for worktree %s: %v", worktreeName, err)
		gitStatus = nil
	}

	// Get worktree creation time
	createdAt, err := getWorktreeCreationTime(targetWorktree.Path)
	if err != nil {
		PrintVerbose("Failed to get creation time for worktree %s: %v", worktreeName, err)
	}

	// Get recent commits
	commits, err := provider.GetWorktreeCommitHistory(targetWorktree.Path, 5)
	if err != nil {
		PrintVerbose("Failed to get recent commits for worktree %s: %v", worktreeName, err)
	}

	// Get modified files
	modifiedFiles, err := provider.GetWorktreeFileChanges(targetWorktree.Path)
	if err != nil {
		PrintVerbose("Failed to get modified files for worktree %s: %v", worktreeName, err)
	}

	// Get base branch info
	baseInfo, err := getBaseBranchInfo(targetWorktree.Path, worktreeName, provider)
	if err != nil {
		PrintVerbose("Failed to get base branch info for worktree %s: %v", worktreeName, err)
	}

	// Try to get JIRA ticket details if the worktree name contains a JIRA key
	var jiraTicket *internal.JiraTicketDetails
	jiraKey := internal.ExtractJiraKey(worktreeName)
	if jiraKey != "" {
		jiraTicket, err = provider.GetJiraTicketDetails(jiraKey)
		if err != nil {
			if errors.Is(err, internal.ErrJiraCliNotFound) {
				PrintVerbose("JIRA CLI not available, skipping ticket details for %s", jiraKey)
			} else {
				PrintVerbose("Failed to get JIRA ticket details for %s: %v", jiraKey, err)
			}
		}
	}

	return &internal.WorktreeInfoData{
		Name:          worktreeName,
		Path:          targetWorktree.Path,
		Branch:        targetWorktree.Branch,
		CreatedAt:     createdAt,
		GitStatus:     gitStatus,
		BaseInfo:      baseInfo,
		Commits:       commits,
		ModifiedFiles: modifiedFiles,
		JiraTicket:    jiraTicket,
	}, nil
}

func displayWorktreeInfo(data *internal.WorktreeInfoData, config *internal.Config) {
	if config == nil {
		config = internal.DefaultConfig()
	}
	renderer := internal.NewInfoRenderer(config)
	output := renderer.RenderWorktreeInfo(data)
	fmt.Println(output)
}

func getWorktreeCreationTime(worktreePath string) (time.Time, error) {
	stat, err := os.Stat(worktreePath)
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

func getBaseBranchInfo(worktreePath, worktreeName string, provider worktreeInfoProvider) (*internal.BranchInfo, error) {
	// Get current branch (not used for base branch detection anymore)
	_, err := provider.GetWorktreeCurrentBranch(worktreePath)
	if err != nil {
		return nil, err
	}
	// Not needed for base branch detection

	// Get upstream branch
	upstream, err := provider.GetWorktreeUpstreamBranch(worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get upstream branch: %w", err)
	}

	// Get ahead/behind count
	aheadBy, behindBy, err := provider.GetWorktreeAheadBehindCount(worktreePath)
	if err != nil {
		// Maintain backward compatibility - use 0,0 if error occurs
		aheadBy, behindBy = 0, 0
	}

	// Try to determine actual base branch - first check stored information
	baseBranch := ""
	if storedBaseBranch, exists := provider.GetState().GetWorktreeBaseBranch(worktreeName); exists {
		baseBranch = storedBaseBranch
	}

	// If no stored information, fall back to git merge-base detection
	if baseBranch == "" {
		// Try configured candidate branches in order of preference
		candidateBranches := provider.GetConfig().Settings.CandidateBranches
		if len(candidateBranches) == 0 {
			// Fallback to default if not configured
			candidateBranches = []string{"main", "master", "develop", "dev"}
		}
		for _, candidate := range candidateBranches {
			exists, err := provider.VerifyWorktreeRef(candidate, worktreePath)
			if err != nil {
				continue // Skip candidates that cause git errors
			}
			if exists {
				// Branch exists, check if it's actually a base
				cmd := exec.Command("git", "merge-base", "--is-ancestor", candidate, "HEAD")
				cmd.Dir = worktreePath
				if err := cmd.Run(); err == nil {
					baseBranch = candidate
					break
				}
			}
		}
	}

	return &internal.BranchInfo{
		Name:     baseBranch,
		Upstream: upstream,
		AheadBy:  aheadBy,
		BehindBy: behindBy,
	}, nil
}

// JSON structs for parsing jira --raw output

func init() {
}
