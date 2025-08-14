package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ErrJiraCliNotFound is returned when the JIRA CLI is not available
var ErrJiraCliNotFound = errors.New("jira CLI not found")

// JiraIssue represents a JIRA issue with its key, type, summary, and status
type JiraIssue struct {
	Type    string
	Key     string
	Summary string
	Status  string
}

// jiraRawResponse represents the raw JSON response from JIRA CLI
type jiraRawResponse struct {
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string  `json:"summary"`
		Created string  `json:"created"`
		DueDate *string `json:"duedate"`
		Status  struct {
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

// IsJiraKey checks if a string matches the JIRA key pattern (PROJECT-NUMBER)
func IsJiraKey(s string) bool {
	matched, _ := regexp.MatchString(`^[A-Z]+-\d+$`, s)
	return matched
}

// ExtractJiraKey extracts a JIRA key from a string, handling prefixed worktree names
// For example: "HOTFIX_INGSVC-5638" returns "INGSVC-5638"
func ExtractJiraKey(s string) string {
	re := regexp.MustCompile(`[A-Z]+-\d+`)
	match := re.FindString(s)
	return match
}

// getJiraUser gets the current JIRA user, using cached value if available
func getJiraUser(manager *Manager) (string, error) {
	config := manager.GetConfig()

	// If we have a cached value, use it
	if config.Jira.Me != "" {
		return config.Jira.Me, nil
	}

	// Otherwise, fetch it from jira CLI and cache it
	meCmd := exec.Command("jira", "me")
	userOutput, err := meCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current JIRA user: %w", err)
	}
	user := strings.TrimSpace(string(userOutput))

	// Cache the result
	config.Jira.Me = user

	// Save the updated config
	if saveErr := manager.SaveConfig(); saveErr != nil {
		// Log warning but don't fail the operation
		fmt.Printf("Warning: failed to save JIRA user to config: %v\n", saveErr)
	}

	return user, nil
}

// GetJiraKeys fetches all JIRA issue keys for the current user
func GetJiraKeys(manager *Manager) ([]string, error) {
	user, err := getJiraUser(manager)
	if err != nil {
		return nil, err
	}

	// Now list issues for the user
	cmd := exec.Command("jira", "issue", "list", "-a"+user, "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JIRA issues: %w", err)
	}

	var keys []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] { // Skip header
		if line = strings.TrimSpace(line); line != "" {
			fields := strings.Split(line, "\t")
			// Find the field that matches JIRA key pattern
			for _, field := range fields {
				trimmedField := strings.TrimSpace(field)
				if IsJiraKey(trimmedField) {
					keys = append(keys, trimmedField)
					break // Only take the first JIRA key found in this line
				}
			}
		}
	}
	return keys, nil
}

// GetJiraIssues fetches all JIRA issues for the current user with full details
func GetJiraIssues(manager *Manager) ([]JiraIssue, error) {
	user, err := getJiraUser(manager)
	if err != nil {
		return nil, err
	}

	// Now list issues for the user
	cmd := exec.Command("jira", "issue", "list", "-a"+user, "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JIRA issues: %w", err)
	}

	return ParseJiraList(string(output)), nil
}

// GetJiraIssue fetches detailed information for a specific JIRA issue
func GetJiraIssue(key string, manager *Manager) (*JiraIssue, error) {
	// For individual issue lookup, use jira issue view to get complete details without truncation
	cmd := exec.Command("jira", "issue", "view", key, "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JIRA issue %s: %w", key, err)
	}

	// Parse the view output to extract issue details
	lines := strings.Split(string(output), "\n")
	var issueType, summary, status string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse the header line with emojis to get status
		if strings.Contains(line, "🐞") && strings.Contains(line, key) {
			if strings.Contains(line, "🐞") {
				issueType = "Bug"
			}
			if strings.Contains(line, "🚧 Open") {
				status = "Open"
			}
		}

		// Parse the title line starting with # - this is the summary
		if strings.HasPrefix(line, "# ") {
			summary = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break // Found the summary, no need to continue
		}
	}

	// Validate that we got the essential fields
	if summary == "" {
		return nil, fmt.Errorf("failed to parse JIRA issue %s: summary not found", key)
	}

	// Default issueType to "Bug" if not found
	if issueType == "" {
		issueType = "Bug"
	}

	return &JiraIssue{
		Type:    issueType,
		Key:     key,
		Summary: summary,
		Status:  status,
	}, nil
}

// ParseJiraList parses the output of 'jira issue list' command
func ParseJiraList(output string) []JiraIssue {
	var issues []JiraIssue
	lines := strings.Split(output, "\n")

	for _, line := range lines[1:] { // Skip header
		if line = strings.TrimSpace(line); line != "" {
			fields := strings.Split(line, "\t")
			if len(fields) >= 3 {
				// Find the JIRA key in this line
				var issueKey, issueType, summary, status string
				keyIndex := -1

				for i, field := range fields {
					trimmedField := strings.TrimSpace(field)
					if IsJiraKey(trimmedField) {
						issueKey = trimmedField
						keyIndex = i
						break
					}
				}

				if issueKey != "" {
					// Type is usually the first field
					issueType = strings.TrimSpace(fields[0])

					// Summary is usually the field after the key
					if keyIndex+1 < len(fields) {
						summary = strings.TrimSpace(fields[keyIndex+1])
					}

					// Status is the last non-empty field
					for i := len(fields) - 1; i >= 0; i-- {
						if trimmed := strings.TrimSpace(fields[i]); trimmed != "" {
							status = trimmed
							break
						}
					}

					issue := JiraIssue{
						Type:    issueType,
						Key:     issueKey,
						Summary: summary,
						Status:  status,
					}
					issues = append(issues, issue)
				}
			}
		}
	}
	return issues
}

// BranchName generates a branch name from a JIRA issue
func (j *JiraIssue) BranchName() string {
	summary := strings.ReplaceAll(j.Summary, " ", "_")
	summary = strings.ReplaceAll(summary, "-", "_")
	// Remove special characters and make it filesystem-safe
	summary = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(summary, "_")
	// Clean up multiple underscores
	summary = regexp.MustCompile(`_+`).ReplaceAllString(summary, "_")
	summary = strings.Trim(summary, "_")

	issueType := strings.ToLower(j.Type)
	if issueType == "story" || issueType == "improvement" {
		issueType = "feature"
	}

	branchName := fmt.Sprintf("%s/%s_%s", issueType, j.Key, summary)
	return branchName
}

// GenerateBranchFromJira fetches a JIRA issue and generates a branch name
func GenerateBranchFromJira(jiraKey string, manager *Manager) (string, error) {
	issue, err := GetJiraIssue(jiraKey, manager)
	if err != nil {
		return "", err
	}

	return issue.BranchName(), nil
}

// IsJiraCliAvailable checks if the JIRA CLI is installed and available
func IsJiraCliAvailable() bool {
	_, err := exec.LookPath("jira")
	return err == nil
}

// GetJiraTicketDetails retrieves detailed JIRA ticket information using the JIRA CLI
// This method uses the more detailed JSON API compared to other jira.go functions
func GetJiraTicketDetails(jiraKey string) (*JiraTicketDetails, error) {
	// Check if jira CLI is available
	if !IsJiraCliAvailable() {
		return nil, ErrJiraCliNotFound
	}

	// Get raw JSON data using jira CLI
	cmd := exec.Command("jira", "issue", "view", jiraKey, "--raw")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get JIRA issue: %w", err)
	}

	// Parse the JSON response
	var jiraResponse jiraRawResponse
	if err := json.Unmarshal(output, &jiraResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JIRA response: %w", err)
	}

	// Build the ticket details from parsed JSON
	ticket := &JiraTicketDetails{
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
			comment := &Comment{
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

// formatJiraURL formats the JIRA URL from the self link
func formatJiraURL(selfURL, key string) string {
	// Extract base URL from self link
	// Self URL format: https://company.atlassian.net/rest/api/2/issue/12345
	if strings.Contains(selfURL, "/rest/api/") {
		baseURL := strings.Split(selfURL, "/rest/api/")[0]
		return fmt.Sprintf("%s/browse/%s", baseURL, key)
	}
	return ""
}
