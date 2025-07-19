package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type MergeBackStatus struct {
	MergeBacksNeeded []MergeBackInfo
	HasUserCommits   bool
}

type MergeBackInfo struct {
	FromBranch  string
	ToBranch    string
	Commits     []MergeBackCommitInfo
	UserCommits []MergeBackCommitInfo
	TotalCount  int
	UserCount   int
}

type MergeBackCommitInfo struct {
	Hash      string
	Message   string
	Author    string
	Email     string
	Timestamp time.Time
	IsUser    bool
}

func CheckMergeBackStatus(configPath string) (*MergeBackStatus, error) {
	// Initialize default empty status
	status := &MergeBackStatus{
		MergeBacksNeeded: []MergeBackInfo{},
		HasUserCommits:   false,
	}

	// Try to find git root and create GitManager
	cwd, err := os.Getwd()
	if err != nil {
		return status, nil
	}

	gitRoot, err := FindGitRoot(cwd)
	if err != nil {
		return status, nil
	}

	gitManager, err := NewGitManager(gitRoot, DefaultWorktreeDirname)
	if err != nil {
		return status, nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	config, err := parseConfigFile(configPath)
	if err != nil {
		return status, nil
	}

	if len(config.Worktrees) <= 1 {
		return status, nil
	}

	userEmail, userName, err := getUserInfo(gitRoot)
	if err != nil {
		// Can't get user info, but continue anyway
		userEmail, userName = "", ""
	}

	// Use tree-based logic to find only the most urgent mergeback (like the mergeback command)
	if config.Tree == nil {
		return status, nil
	}

	// Start from deepest leaf nodes and work up the tree
	deepestLeaves := config.Tree.GetAllDeepestLeafNodes()
	checkedNodes := make(map[string]bool)

	for _, leaf := range deepestLeaves {
		// Walk up from this leaf, checking each node for mergebacks
		current := leaf
		for current != nil && current.Parent != nil {
			// Skip if we've already checked this node
			if checkedNodes[current.Name] {
				current = current.Parent
				continue
			}
			checkedNodes[current.Name] = true

			// Check if both branches exist
			fromExists, _ := gitManager.BranchExistsLocalOrRemote(current.Config.Branch)
			toExists, _ := gitManager.BranchExistsLocalOrRemote(current.Parent.Config.Branch)
			if !fromExists || !toExists {
				current = current.Parent
				continue
			}

			// Get commits that need to be merged back
			commits, err := getCommitsNeedingMergeBack(gitRoot, current.Parent.Config.Branch, current.Config.Branch)
			if err != nil {
				fmt.Println("⚠️  Warning:", err)
				current = current.Parent
				continue
			}

			if len(commits) > 0 {
				// Identify user commits
				userCommits := []MergeBackCommitInfo{}
				for _, commit := range commits {
					if isUserCommit(commit, userEmail, userName) {
						commit.IsUser = true
						userCommits = append(userCommits, commit)
						status.HasUserCommits = true
					}
				}

				mergeBackInfo := MergeBackInfo{
					FromBranch:  current.Name,
					ToBranch:    current.Parent.Name,
					Commits:     commits,
					UserCommits: userCommits,
					TotalCount:  len(commits),
					UserCount:   len(userCommits),
				}

				status.MergeBacksNeeded = append(status.MergeBacksNeeded, mergeBackInfo)
			}

			current = current.Parent
		}
	}

	return status, nil
}

func parseConfigFile(configPath string) (*GBMConfig, error) {
	// Use the existing ParseGBMConfig function that properly builds the tree
	return ParseGBMConfig(configPath)
}

func getUserInfo(repoPath string) (string, string, error) {
	emailBytes, err := ExecGitCommand(repoPath, "config", "user.email")
	email := strings.TrimSpace(string(emailBytes))
	if err != nil {
		email = ""
	}

	nameBytes, err := ExecGitCommand(repoPath, "config", "user.name")
	name := strings.TrimSpace(string(nameBytes))
	if err != nil {
		name = ""
	}

	if email == "" && name == "" {
		return "", "", fmt.Errorf("no git user configuration found")
	}

	return email, name, nil
}

func getCommitsNeedingMergeBack(repoPath, targetBranch, sourceBranch string) ([]MergeBackCommitInfo, error) {
	// First, try to fetch to ensure we have the latest remote state
	ExecGitCommand(repoPath, "fetch", "--quiet")

	// Use remote branches for mergeback detection
	remoteTargetBranch := Remote(targetBranch)
	remoteSourceBranch := Remote(sourceBranch)

	output, err := ExecGitCommand(repoPath, "log", remoteTargetBranch+".."+remoteSourceBranch, "--format=%H|%s|%an|%ae|%ct")
	if err != nil {
		// If remote branch doesn't exist, this indicates a configuration error
		return nil, fmt.Errorf("remote branch '%s' or '%s' does not exist - check your gbm.branchconfig.yaml configuration", remoteTargetBranch, remoteSourceBranch)
	}


	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []MergeBackCommitInfo{}, nil
	}

	var commits []MergeBackCommitInfo
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 5 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[4], 10, 64)
		if err != nil {
			continue
		}

		commit := MergeBackCommitInfo{
			Hash:      parts[0],
			Message:   parts[1],
			Author:    parts[2],
			Email:     parts[3],
			Timestamp: time.Unix(timestamp, 0),
			IsUser:    false,
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func isUserCommit(commit MergeBackCommitInfo, userEmail, userName string) bool {
	if userEmail != "" && commit.Email == userEmail {
		return true
	}
	if userName != "" && commit.Author == userName {
		return true
	}
	return false
}

func FormatMergeBackAlert(status *MergeBackStatus) string {
	if status == nil || len(status.MergeBacksNeeded) == 0 {
		return ""
	}

	var output strings.Builder
	output.WriteString("⚠️  Merge-back required in tracked branches:\n\n")

	for _, info := range status.MergeBacksNeeded {
		output.WriteString(fmt.Sprintf("%s → %s: %d commits need merge-back",
			info.FromBranch, info.ToBranch, info.TotalCount))

		if info.UserCount > 0 {
			output.WriteString(fmt.Sprintf(" (%d by you)", info.UserCount))
		} else {
			output.WriteString(" (0 by you)")
		}
		output.WriteString("\n")

		for _, commit := range info.UserCommits {
			relativeTime := FormatRelativeTime(commit.Timestamp)
			output.WriteString(fmt.Sprintf("• %s - %s (you - %s)\n",
				commit.Hash[:7], commit.Message, relativeTime))
		}
		output.WriteString("\n")
	}

	return output.String()
}
