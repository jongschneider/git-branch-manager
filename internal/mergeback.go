package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
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

type EnvVarMapping struct {
	Name   string
	Branch string
	Order  int
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

	gitManager, err := NewGitManager(gitRoot)
	if err != nil {
		return nil, nil
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	envMappings, err := parseEnvrcFile(configPath)
	if err != nil {
		return nil, nil
	}

	if len(envMappings) <= 1 {
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

	for i := len(envMappings) - 1; i > 0; i-- {
		fromMapping := envMappings[i]
		toMapping := envMappings[i-1]

		fromExists, _ := gitManager.BranchExists(fromMapping.Branch)
		toExists, _ := gitManager.BranchExists(toMapping.Branch)
		if !fromExists || !toExists {
			continue
		}

		commits, err := getCommitsNeedingMergeBack(gitRoot, toMapping.Branch, fromMapping.Branch)
		if err != nil {
			continue
		}

		if len(commits) == 0 {
			continue
		}

		userCommits := []MergeBackCommitInfo{}
		for _, commit := range commits {
			if isUserCommit(commit, userEmail, userName) {
				commit.IsUser = true
				userCommits = append(userCommits, commit)
				status.HasUserCommits = true
			}
		}

		mergeBackInfo := MergeBackInfo{
			FromBranch:  fromMapping.Name,
			ToBranch:    toMapping.Name,
			Commits:     commits,
			UserCommits: userCommits,
			TotalCount:  len(commits),
			UserCount:   len(userCommits),
		}

		status.MergeBacksNeeded = append(status.MergeBacksNeeded, mergeBackInfo)
	}

	return status, nil
}

func parseEnvrcFile(configPath string) ([]EnvVarMapping, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mappings []EnvVarMapping
	scanner := bufio.NewScanner(file)
	order := 0

	envVarRegex := regexp.MustCompile(`^([A-Z_][A-Z0-9_]*)=(.+)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := envVarRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			mappings = append(mappings, EnvVarMapping{
				Name:   matches[1],
				Branch: matches[2],
				Order:  order,
			})
			order++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if mappings == nil {
		mappings = []EnvVarMapping{}
	}

	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].Order < mappings[j].Order
	})

	return mappings, nil
}

func getUserInfo(repoPath string) (string, string, error) {
	emailCmd := exec.Command("git", "config", "user.email")
	emailCmd.Dir = repoPath
	emailBytes, err := emailCmd.Output()
	email := strings.TrimSpace(string(emailBytes))
	if err != nil {
		email = ""
	}

	nameCmd := exec.Command("git", "config", "user.name")
	nameCmd.Dir = repoPath
	nameBytes, err := nameCmd.Output()
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
	cmd := exec.Command("git", "log", targetBranch+".."+sourceBranch, "--format=%H|%s|%an|%ae|%ct")
	cmd.Dir = repoPath
	output, err := cmd.Output()
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
