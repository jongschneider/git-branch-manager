package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gbm/internal"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <worktree-name>",
	Short: "Display detailed information about a worktree",
	Long: `Display comprehensive information about a specific worktree including:
- Worktree metadata (name, path, branch, creation date)
- Git status and branch information
- JIRA ticket details (if the worktree name matches a JIRA key)
- Recent commits and modified files`,
	Args: cobra.ExactArgs(1),
	RunE: runInfoCommand,
}

func runInfoCommand(cmd *cobra.Command, args []string) error {
	worktreeName := args[0]

	// Handle current directory reference
	if worktreeName == "." {
		currentPath, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		worktreeName = filepath.Base(currentPath)
	}

	// Find git root
	gitRoot, err := internal.FindGitRoot(".")
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// Initialize git manager
	gitManager, err := internal.NewGitManager(gitRoot)
	if err != nil {
		return fmt.Errorf("failed to initialize git manager: %w", err)
	}

	// Get worktree information
	worktreeInfo, err := getWorktreeInfo(gitManager, worktreeName)
	if err != nil {
		return fmt.Errorf("failed to get worktree info: %w", err)
	}

	// Display the information
	displayWorktreeInfo(worktreeInfo)

	return nil
}

func getWorktreeInfo(gitManager *internal.GitManager, worktreeName string) (*internal.WorktreeInfoData, error) {
	// Get all worktrees
	worktrees, err := gitManager.GetWorktrees()
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
	gitStatus, err := gitManager.GetWorktreeStatus(targetWorktree.Path)
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
	commits, err := getRecentCommits(targetWorktree.Path, 5)
	if err != nil {
		PrintVerbose("Failed to get recent commits for worktree %s: %v", worktreeName, err)
	}

	// Get modified files
	modifiedFiles, err := getModifiedFiles(targetWorktree.Path)
	if err != nil {
		PrintVerbose("Failed to get modified files for worktree %s: %v", worktreeName, err)
	}

	// Get base branch info
	baseInfo, err := getBaseBranchInfo(targetWorktree.Path)
	if err != nil {
		PrintVerbose("Failed to get base branch info for worktree %s: %v", worktreeName, err)
	}

	// Try to get JIRA ticket details if the worktree name matches a JIRA key
	var jiraTicket *internal.JiraTicketDetails
	if internal.IsJiraKey(worktreeName) {
		jiraTicket, err = getJiraTicketDetails(worktreeName)
		if err != nil {
			PrintVerbose("Failed to get JIRA ticket details for %s: %v", worktreeName, err)
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

func displayWorktreeInfo(data *internal.WorktreeInfoData) {
	renderer := internal.NewInfoRenderer()
	output := renderer.RenderWorktreeInfo(data)
	fmt.Print(output)
}

func getWorktreeCreationTime(worktreePath string) (time.Time, error) {
	stat, err := os.Stat(worktreePath)
	if err != nil {
		return time.Time{}, err
	}
	return stat.ModTime(), nil
}

func getRecentCommits(worktreePath string, count int) ([]internal.CommitInfo, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline", "--format=%H|%s|%an|%ct")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []internal.CommitInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			continue
		}

		commits = append(commits, internal.CommitInfo{
			Hash:      parts[0],
			Message:   parts[1],
			Author:    parts[2],
			Timestamp: time.Unix(timestamp, 0),
		})
	}

	return commits, nil
}

func getModifiedFiles(worktreePath string) ([]internal.FileChange, error) {
	// Get unstaged changes
	cmd := exec.Command("git", "diff", "--numstat")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var files []internal.FileChange
	
	// Parse unstaged changes
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}

		additions, _ := strconv.Atoi(parts[0])
		deletions, _ := strconv.Atoi(parts[1])
		path := parts[2]

		// Determine status based on changes
		status := "M"
		if additions > 0 && deletions == 0 {
			status = "A"
		} else if additions == 0 && deletions > 0 {
			status = "D"
		}

		files = append(files, internal.FileChange{
			Path:      path,
			Status:    status,
			Additions: additions,
			Deletions: deletions,
		})
	}

	// Get staged changes
	cmd = exec.Command("git", "diff", "--cached", "--numstat")
	cmd.Dir = worktreePath
	output, err = cmd.Output()
	if err == nil {
		lines = strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				continue
			}

			additions, _ := strconv.Atoi(parts[0])
			deletions, _ := strconv.Atoi(parts[1])
			path := parts[2]

			// Check if file already exists in our list
			found := false
			for i, existing := range files {
				if existing.Path == path {
					// Update existing entry with staged changes
					files[i].Additions += additions
					files[i].Deletions += deletions
					found = true
					break
				}
			}

			if !found {
				status := "M"
				if additions > 0 && deletions == 0 {
					status = "A"
				} else if additions == 0 && deletions > 0 {
					status = "D"
				}

				files = append(files, internal.FileChange{
					Path:      path,
					Status:    status,
					Additions: additions,
					Deletions: deletions,
				})
			}
		}
	}

	return files, nil
}

func getBaseBranchInfo(worktreePath string) (*internal.BranchInfo, error) {
	// Get current branch (not used for base branch detection anymore)
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	_ = strings.TrimSpace(string(output)) // Not needed for base branch detection

	// Get upstream branch
	cmd = exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	cmd.Dir = worktreePath
	output, err = cmd.Output()
	upstream := ""
	if err == nil {
		upstream = strings.TrimSpace(string(output))
	}

	// Get ahead/behind count
	cmd = exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = worktreePath
	output, err = cmd.Output()
	aheadBy, behindBy := 0, 0
	if err == nil {
		parts := strings.Fields(string(output))
		if len(parts) == 2 {
			aheadBy, _ = strconv.Atoi(parts[0])
			behindBy, _ = strconv.Atoi(parts[1])
		}
	}

	// Try to determine actual base branch using git merge-base
	baseBranch := ""
	// Try common base branches in order of preference
	candidateBranches := []string{"main", "master", "develop", "dev"}
	for _, candidate := range candidateBranches {
		cmd = exec.Command("git", "rev-parse", "--verify", candidate)
		cmd.Dir = worktreePath
		if _, err := cmd.Output(); err == nil {
			// Branch exists, check if it's actually a base
			cmd = exec.Command("git", "merge-base", "--is-ancestor", candidate, "HEAD")
			cmd.Dir = worktreePath
			if err := cmd.Run(); err == nil {
				baseBranch = candidate
				break
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
type JiraRawResponse struct {
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary     string `json:"summary"`
		Created     string `json:"created"`
		DueDate     *string   `json:"duedate"`
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			Name string `json:"name"`
		} `json:"priority"`
		Reporter struct {
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
		} `json:"reporter"`
		Assignee *struct {
			DisplayName  string `json:"displayName"`
			EmailAddress string `json:"emailAddress"`
		} `json:"assignee"`
		Parent *struct {
			Key string `json:"key"`
		} `json:"parent"`
		Description *struct {
			Content []struct {
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"content"`
		} `json:"description"`
		Comment struct {
			Comments []struct {
				Author struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Body struct {
					Content []struct {
						Content []struct {
							Text string `json:"text"`
						} `json:"content"`
					} `json:"content"`
				} `json:"body"`
				Created string `json:"created"`
			} `json:"comments"`
		} `json:"comment"`
	} `json:"fields"`
}

func getJiraTicketDetails(jiraKey string) (*internal.JiraTicketDetails, error) {
	// Check if jira CLI is available
	if _, err := exec.LookPath("jira"); err != nil {
		return nil, fmt.Errorf("jira CLI not found: %w", err)
	}

	// Get raw JSON data using jira CLI
	cmd := exec.Command("jira", "issue", "view", jiraKey, "--raw")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get JIRA issue: %w", err)
	}

	// Parse the JSON response
	var jiraResponse JiraRawResponse
	if err := json.Unmarshal(output, &jiraResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JIRA response: %w", err)
	}

	// Build the ticket details from parsed JSON
	ticket := &internal.JiraTicketDetails{
		Key:     jiraResponse.Key,
		Summary: jiraResponse.Fields.Summary,
		Status:  jiraResponse.Fields.Status.Name,
		URL:     formatJiraURL(jiraResponse.Self, jiraResponse.Key),
	}

	// Parse created date
	if jiraResponse.Fields.Created != "" {
		if createdDate, err := time.Parse(time.RFC3339, jiraResponse.Fields.Created); err == nil {
			ticket.Created = createdDate
		}
	}

	// Add priority
	if jiraResponse.Fields.Priority.Name != "" {
		ticket.Priority = jiraResponse.Fields.Priority.Name
	}

	// Add reporter
	if jiraResponse.Fields.Reporter.DisplayName != "" {
		reporter := jiraResponse.Fields.Reporter.DisplayName
		if jiraResponse.Fields.Reporter.EmailAddress != "" {
			reporter += " (" + jiraResponse.Fields.Reporter.EmailAddress + ")"
		}
		ticket.Reporter = reporter
	}

	// Add assignee (can be null)
	if jiraResponse.Fields.Assignee != nil {
		assignee := jiraResponse.Fields.Assignee.DisplayName
		if jiraResponse.Fields.Assignee.EmailAddress != "" {
			assignee += " (" + jiraResponse.Fields.Assignee.EmailAddress + ")"
		}
		ticket.Assignee = assignee
	}

	// Add due date (can be null)
	if jiraResponse.Fields.DueDate != nil && *jiraResponse.Fields.DueDate != "" {
		if dueDate, err := time.Parse("2006-01-02", *jiraResponse.Fields.DueDate); err == nil {
			ticket.DueDate = &dueDate
		}
	}

	// Add epic information
	if jiraResponse.Fields.Parent != nil {
		ticket.Epic = jiraResponse.Fields.Parent.Key
	}

	// Add latest comment
	if len(jiraResponse.Fields.Comment.Comments) > 0 {
		latest := jiraResponse.Fields.Comment.Comments[0]
		
		// Extract comment text from nested structure
		var commentText strings.Builder
		for _, content := range latest.Body.Content {
			for _, textContent := range content.Content {
				if textContent.Text != "" {
					commentText.WriteString(textContent.Text)
				}
			}
		}

		if commentText.Len() > 0 {
			comment := &internal.Comment{
				Author:  latest.Author.DisplayName,
				Content: commentText.String(),
			}
			
			// Parse comment timestamp
			if latest.Created != "" {
				if timestamp, err := time.Parse(time.RFC3339, latest.Created); err == nil {
					comment.Timestamp = timestamp
				}
			}
			
			ticket.LatestComment = comment
		}
	}

	return ticket, nil
}

// formatJiraURL converts the API URL to user-friendly browse URL
// Input: "https://thetalake.atlassian.net/rest/api/3/issue/45305", "INGSVC-4929"
// Output: "https://thetalake.atlassian.net/browse/INGSVC-4929"
func formatJiraURL(selfURL, key string) string {
	if selfURL == "" || key == "" {
		return ""
	}

	// Find the position of "/rest" in the URL
	restIndex := strings.Index(selfURL, "/rest")
	if restIndex == -1 {
		// If "/rest" not found, return the original URL
		return selfURL
	}

	// Extract the base URL (everything before "/rest")
	baseURL := selfURL[:restIndex]
	
	// Construct the browse URL
	return baseURL + "/browse/" + key
}

func init() {
	rootCmd.AddCommand(infoCmd)
}