package internal

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// JiraIssue represents a JIRA issue with its key, type, summary, and status
type JiraIssue struct {
	Type    string
	Key     string
	Summary string
	Status  string
}

// IsJiraKey checks if a string matches the JIRA key pattern (PROJECT-NUMBER)
func IsJiraKey(s string) bool {
	matched, _ := regexp.MatchString(`^[A-Z]+-\d+$`, s)
	return matched
}

// GetJiraKeys fetches all JIRA issue keys for the current user
func GetJiraKeys() ([]string, error) {
	// Get current user first
	meCmd := exec.Command("jira", "me")
	userOutput, err := meCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current JIRA user: %w", err)
	}
	user := strings.TrimSpace(string(userOutput))

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
func GetJiraIssues() ([]JiraIssue, error) {
	// Get current user first
	meCmd := exec.Command("jira", "me")
	userOutput, err := meCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current JIRA user: %w", err)
	}
	user := strings.TrimSpace(string(userOutput))

	// Now list issues for the user
	cmd := exec.Command("jira", "issue", "list", "-a"+user, "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JIRA issues: %w", err)
	}

	return ParseJiraList(string(output)), nil
}

// GetJiraIssue fetches detailed information for a specific JIRA issue
func GetJiraIssue(key string) (*JiraIssue, error) {
	// Get current user first
	meCmd := exec.Command("jira", "me")
	userOutput, err := meCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current JIRA user: %w", err)
	}
	user := strings.TrimSpace(string(userOutput))

	// Get the issue from the list command
	cmd := exec.Command("jira", "issue", "list", "-a"+user, "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JIRA issues: %w", err)
	}

	// Parse the list output to find our specific issue
	lines := strings.Split(string(output), "\n")
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

				if issueKey == key {
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

					return &JiraIssue{
						Type:    issueType,
						Key:     issueKey,
						Summary: summary,
						Status:  status,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("JIRA issue %s not found in your assigned issues", key)
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

	return fmt.Sprintf("%s/%s_%s", issueType, j.Key, summary)
}

// GenerateBranchFromJira fetches a JIRA issue and generates a branch name
func GenerateBranchFromJira(jiraKey string) (string, error) {
	issue, err := GetJiraIssue(jiraKey)
	if err != nil {
		return "", err
	}

	return issue.BranchName(), nil
}