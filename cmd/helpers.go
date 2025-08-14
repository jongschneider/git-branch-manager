package cmd

import (
	"fmt"
	"strings"

	"gbm/internal"
)

// createBranchNameGenerator creates a function that generates branch names with the specified prefix
func createBranchNameGenerator(prefix string) func(worktreeName, jiraTicket, targetSuffix string, manager *internal.Manager) (string, error) {
	return func(worktreeName, jiraTicket, targetSuffix string, manager *internal.Manager) (string, error) {
		var branchName string

		if jiraTicket != "" && internal.IsJiraKey(jiraTicket) {
			// Generate branch name from JIRA ticket
			if manager != nil {
				jiraBranchName, err := internal.GenerateBranchFromJira(jiraTicket, manager)
				if err != nil {
					PrintVerbose("Failed to generate branch name from JIRA issue %s: %v", jiraTicket, err)
					// Fallback to simple format
					if targetSuffix != "" {
						branchName = fmt.Sprintf("%s/%s_%s", prefix, strings.ToUpper(jiraTicket), targetSuffix)
					} else {
						branchName = fmt.Sprintf("%s/%s", prefix, strings.ToUpper(jiraTicket))
					}
				} else {
					// Replace any prefix with the specified prefix
					parts := strings.Split(jiraBranchName, "/")
					if len(parts) > 1 {
						parts[0] = prefix
						baseName := strings.Join(parts, "/")
						if targetSuffix != "" {
							branchName = fmt.Sprintf("%s_%s", baseName, targetSuffix)
						} else {
							branchName = baseName
						}
					} else {
						if targetSuffix != "" {
							branchName = fmt.Sprintf("%s/%s_%s", prefix, jiraBranchName, targetSuffix)
						} else {
							branchName = fmt.Sprintf("%s/%s", prefix, jiraBranchName)
						}
					}
				}
			} else {
				// No manager available, use simple format
				if targetSuffix != "" {
					branchName = fmt.Sprintf("%s/%s_%s", prefix, strings.ToUpper(jiraTicket), targetSuffix)
				} else {
					branchName = fmt.Sprintf("%s/%s", prefix, strings.ToUpper(jiraTicket))
				}
			}
		} else {
			// Generate from worktree name
			cleanName := strings.ReplaceAll(worktreeName, " ", "-")
			cleanName = strings.ReplaceAll(cleanName, "_", "-")
			cleanName = strings.ToLower(cleanName)
			if targetSuffix != "" {
				branchName = fmt.Sprintf("%s/%s_%s", prefix, cleanName, targetSuffix)
			} else {
				branchName = fmt.Sprintf("%s/%s", prefix, cleanName)
			}
		}

		return branchName, nil
	}
}

// createBranchNameGeneratorWithCreator creates a function that generates branch names with the specified prefix using the hotfixCreator interface
func createBranchNameGeneratorWithCreator(prefix string, creator hotfixCreator) func(worktreeName, jiraTicket, targetSuffix string) (string, error) {
	return func(worktreeName, jiraTicket, targetSuffix string) (string, error) {
		// Guard clause: handle non-JIRA case first and return early
		if jiraTicket == "" || !internal.IsJiraKey(jiraTicket) {
			// Generate from worktree name
			cleanName := strings.ReplaceAll(worktreeName, " ", "-")
			cleanName = strings.ReplaceAll(cleanName, "_", "-")
			cleanName = strings.ToLower(cleanName)
			if targetSuffix != "" {
				return fmt.Sprintf("%s/%s_%s", prefix, cleanName, targetSuffix), nil
			}
			return fmt.Sprintf("%s/%s", prefix, cleanName), nil
		}

		// Generate branch name from JIRA ticket
		jiraBranchName, err := creator.GenerateBranchFromJira(jiraTicket)
		if err != nil {
			PrintVerbose("Failed to generate branch name from JIRA issue %s: %v", jiraTicket, err)
			// Fallback to simple format
			if targetSuffix != "" {
				return fmt.Sprintf("%s/%s_%s", prefix, strings.ToUpper(jiraTicket), targetSuffix), nil
			}
			return fmt.Sprintf("%s/%s", prefix, strings.ToUpper(jiraTicket)), nil
		}

		// Replace any prefix with the specified prefix
		parts := strings.Split(jiraBranchName, "/")
		if len(parts) > 1 {
			parts[0] = prefix
			baseName := strings.Join(parts, "/")
			if targetSuffix != "" {
				return fmt.Sprintf("%s_%s", baseName, targetSuffix), nil
			}
			return baseName, nil
		}

		if targetSuffix != "" {
			return fmt.Sprintf("%s/%s_%s", prefix, jiraBranchName, targetSuffix), nil
		}
		return fmt.Sprintf("%s/%s", prefix, jiraBranchName), nil
	}
}
