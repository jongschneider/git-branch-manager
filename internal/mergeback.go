package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
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
	// Try to find git root and create GitManager
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	gitRoot, err := FindGitRoot(cwd)
	if err != nil {
		return nil, nil
	}

	gitManager, err := NewGitManager(gitRoot, DefaultWorktreeDirname)
	if err != nil {
		return nil, nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	config, err := parseConfigFile(configPath)
	if err != nil {
		return nil, nil
	}

	if len(config.Worktrees) <= 1 {
		return &MergeBackStatus{}, nil
	}

	userEmail, userName, err := getUserInfo(gitRoot)
	if err != nil {
		return &MergeBackStatus{}, nil
	}

	status := &MergeBackStatus{
		MergeBacksNeeded: []MergeBackInfo{},
		HasUserCommits:   false,
	}

	// Check each worktree that has a merge_into target
	for worktreeName, worktreeConfig := range config.Worktrees {
		if worktreeConfig.MergeInto == "" {
			continue // Skip root worktrees (no merge target)
		}

		// Find the target worktree config
		targetConfig, exists := config.Worktrees[worktreeConfig.MergeInto]
		if !exists {
			continue // Skip if target worktree doesn't exist in config
		}

		// Check if both branches exist
		fromExists, _ := gitManager.BranchExists(worktreeConfig.Branch)
		toExists, _ := gitManager.BranchExists(targetConfig.Branch)
		if !fromExists || !toExists {
			continue
		}

		// Get commits that need to be merged back
		commits, err := getCommitsNeedingMergeBack(gitRoot, targetConfig.Branch, worktreeConfig.Branch)
		if err != nil {
			continue
		}

		if len(commits) == 0 {
			continue
		}

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
			FromBranch:  worktreeName,
			ToBranch:    worktreeConfig.MergeInto,
			Commits:     commits,
			UserCommits: userCommits,
			TotalCount:  len(commits),
			UserCount:   len(userCommits),
		}

		status.MergeBacksNeeded = append(status.MergeBacksNeeded, mergeBackInfo)
	}

	return status, nil
}


func parseConfigFile(configPath string) (*GBMConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config GBMConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
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
	output, err := ExecGitCommand(repoPath, "log", targetBranch+".."+sourceBranch, "--format=%H|%s|%an|%ae|%ct")
	if err != nil {
		return nil, err
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
