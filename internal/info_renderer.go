package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type InfoRenderer struct {
	headerStyle    lipgloss.Style
	sectionStyle   lipgloss.Style
	keyStyle       lipgloss.Style
	valueStyle     lipgloss.Style
	borderStyle    lipgloss.Style
	separatorStyle lipgloss.Style
	titleStyle     lipgloss.Style
	subtitleStyle  lipgloss.Style
	statusStyle    lipgloss.Style
	commitStyle    lipgloss.Style
	fileStyle      lipgloss.Style
	jiraStyle      lipgloss.Style
}

// getTerminalWidth returns the terminal width, with multiple fallbacks
func getTerminalWidth() int {
	// Try getting terminal size directly
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil && width > 0 {
		return width
	}

	// Fallback to COLUMNS environment variable (works in tmux)
	if columns := os.Getenv("COLUMNS"); columns != "" {
		if w, err := strconv.Atoi(columns); err == nil && w > 0 {
			return w
		}
	}

	// Try tput cols command as another fallback
	if cmd := exec.Command("tput", "cols"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			if w, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil && w > 0 {
				return w
			}
		}
	}

	// Final fallback
	return 80
}

func NewInfoRenderer() *InfoRenderer {
	// Define adaptive colors for better light/dark theme support
	primaryColor := lipgloss.AdaptiveColor{Light: "#7D56F4", Dark: "#A78BFA"}
	secondaryColor := lipgloss.AdaptiveColor{Light: "#3B82F6", Dark: "#60A5FA"}
	textColor := lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#F9FAFB"}
	subtleColor := lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	successColor := lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"}
	warningColor := lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#FCD34D"}

	termWidth := getTerminalWidth()
	contentWidth := termWidth - 10 // Account for borders and padding

	return &InfoRenderer{
		headerStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Align(lipgloss.Center).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2).
			Width(contentWidth),

		sectionStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtleColor).
			Padding(1, 2).
			Width(contentWidth),

		keyStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Width(15),

		valueStyle: lipgloss.NewStyle().
			Foreground(textColor),

		borderStyle: lipgloss.NewStyle().
			Foreground(subtleColor),

		separatorStyle: lipgloss.NewStyle().
			Foreground(subtleColor),

		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor),

		subtitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(subtleColor),

		statusStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor),

		commitStyle: lipgloss.NewStyle().
			Foreground(subtleColor),

		fileStyle: lipgloss.NewStyle().
			Foreground(warningColor),

		jiraStyle: lipgloss.NewStyle().
			Foreground(secondaryColor),
	}
}

func (r *InfoRenderer) RenderWorktreeInfo(data *WorktreeInfoData) string {
	var sections []string

	// Header
	header := r.headerStyle.Render(fmt.Sprintf("📋 WORKTREE INFO: %s", data.Name))
	sections = append(sections, header)

	// Worktree Section
	worktreeSection := r.renderWorktreeSection(data)
	sections = append(sections, worktreeSection)

	// JIRA Section (if available)
	if data.JiraTicket != nil {
		jiraSection := r.renderJiraSection(data.JiraTicket)
		sections = append(sections, jiraSection)
	}

	// Git Status Section
	gitSection := r.renderGitSection(data)
	sections = append(sections, gitSection)

	// Use lipgloss.JoinVertical with compact spacing
	result := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return result + "\n"
}

func (r *InfoRenderer) renderWorktreeSection(data *WorktreeInfoData) string {
	var content strings.Builder

	content.WriteString("📁 WORKTREE\n")
	content.WriteString(r.renderKeyValue("Name", data.Name))
	content.WriteString(r.renderKeyValue("Path", data.Path))
	content.WriteString(r.renderKeyValue("Branch", data.Branch))

	if !data.CreatedAt.IsZero() {
		timeAgo := time.Since(data.CreatedAt)
		timeStr := fmt.Sprintf("%s (%s ago)",
			data.CreatedAt.Format("2006-01-02 15:04:05"),
			r.formatDuration(timeAgo))
		content.WriteString(r.renderKeyValue("Created", timeStr))
	}

	if data.GitStatus != nil {
		status := r.formatGitStatus(data.GitStatus)
		content.WriteString(r.renderKeyValue("Status", status))
	}

	return r.sectionStyle.Render(content.String())
}

func (r *InfoRenderer) renderJiraSection(jira *JiraTicketDetails) string {
	var content strings.Builder

	content.WriteString("🎫 JIRA TICKET\n")
	content.WriteString(r.renderKeyValue("Key", jira.Key))

	if jira.Summary != "" {
		termWidth := getTerminalWidth()
		summaryWidth := termWidth - 25 // Account for borders, padding, and key label
		if summaryWidth < 30 {
			summaryWidth = 30 // minimum width
		}
		wrappedSummary := r.wrapText(jira.Summary, summaryWidth)
		content.WriteString(r.renderKeyValue("Summary", wrappedSummary))
	}

	if jira.Status != "" {
		content.WriteString(r.renderKeyValue("Status", jira.Status))
	}

	if jira.Assignee != "" {
		content.WriteString(r.renderKeyValue("Assignee", jira.Assignee))
	}

	if jira.Priority != "" {
		priority := r.formatPriority(jira.Priority)
		content.WriteString(r.renderKeyValue("Priority", priority))
	}

	if jira.Reporter != "" {
		content.WriteString(r.renderKeyValue("Reporter", jira.Reporter))
	}

	if !jira.Created.IsZero() {
		content.WriteString(r.renderKeyValue("Created", jira.Created.Format("2006-01-02 15:04:05")))
	}

	if jira.DueDate != nil {
		content.WriteString(r.renderKeyValue("Due Date", jira.DueDate.Format("2006-01-02 15:04:05")))
	}

	if jira.Epic != "" {
		content.WriteString(r.renderKeyValue("Epic", jira.Epic))
	}

	if jira.URL != "" {
		content.WriteString(r.renderKeyValue("Link", jira.URL))
	}

	if jira.LatestComment != nil {
		timeAgo := time.Since(jira.LatestComment.Timestamp)
		commentHeader := fmt.Sprintf("💬 Latest Comment (%s ago):", r.formatDuration(timeAgo))
		content.WriteString(commentHeader + "\n")

		// Wrap the comment text to fit within borders
		termWidth := getTerminalWidth()
		commentWidth := termWidth - 30 // Account for borders, padding, and indentation
		if commentWidth < 40 {
			commentWidth = 40 // minimum width
		}
		wrappedComment := r.wrapText(jira.LatestComment.Content, commentWidth)
		for _, line := range strings.Split(wrappedComment, "\n") {
			content.WriteString(fmt.Sprintf("    %s\n", line))
		}
		content.WriteString(fmt.Sprintf("    - %s", jira.LatestComment.Author))
	}

	return r.sectionStyle.Render(content.String())
}

func (r *InfoRenderer) renderGitSection(data *WorktreeInfoData) string {
	var content strings.Builder

	content.WriteString("🌿 GIT STATUS\n")

	if data.BaseInfo != nil {
		if data.BaseInfo.Name != "" {
			content.WriteString(r.renderKeyValue("Base Branch", data.BaseInfo.Name))
		}
		if data.BaseInfo.Upstream != "" {
			content.WriteString(r.renderKeyValue("Upstream", data.BaseInfo.Upstream))
		}
		if data.BaseInfo.AheadBy > 0 || data.BaseInfo.BehindBy > 0 {
			position := fmt.Sprintf("↑ %d commits ahead, ↓ %d commits behind",
				data.BaseInfo.AheadBy, data.BaseInfo.BehindBy)
			content.WriteString(r.renderKeyValue("Position", position))
		}
	}

	// Recent commits
	if len(data.Commits) > 0 {
		latest := data.Commits[0]
		timeAgo := time.Since(latest.Timestamp)
		termWidth := getTerminalWidth()
		commitWidth := termWidth - 25 // Account for borders, padding, and key label
		if commitWidth < 30 {
			commitWidth = 30
		}
		wrappedMessage := r.wrapText(latest.Message, commitWidth-20) // Reserve space for hash and time
		lastCommit := fmt.Sprintf("%s (%s) - %s ago",
			wrappedMessage,
			latest.Hash[:7],
			r.formatDuration(timeAgo))
		content.WriteString(r.renderKeyValue("Last Commit", lastCommit))
		content.WriteString(r.renderKeyValue("Author", latest.Author))
	}

	// Modified files
	if len(data.ModifiedFiles) > 0 {
		content.WriteString("Modified Files:\n")
		termWidth := getTerminalWidth()
		filePathWidth := termWidth - 40 // Account for borders, status, and changes
		if filePathWidth < 20 {
			filePathWidth = 20
		}

		for _, file := range data.ModifiedFiles {
			statusIcon := r.getStatusIcon(file.Status)
			changes := fmt.Sprintf("(+%d, -%d)", file.Additions, file.Deletions)

			// Truncate long file paths intelligently
			displayPath := file.Path
			if len(displayPath) > filePathWidth {
				displayPath = "..." + displayPath[len(displayPath)-filePathWidth+3:]
			}

			// Use lipgloss for better formatting
			statusCol := lipgloss.NewStyle().Width(3).Render(statusIcon)
			pathCol := lipgloss.NewStyle().Width(filePathWidth).Render(displayPath)
			changesCol := r.fileStyle.Render(changes)

			line := lipgloss.JoinHorizontal(lipgloss.Top, "  ", statusCol, "  ", pathCol, " ", changesCol)
			content.WriteString(line + "\n")
		}
	}

	// Recent commits list
	if len(data.Commits) > 1 {
		content.WriteString("Recent Commits:\n")
		for i, commit := range data.Commits {
			if i == 0 {
				continue // Skip the first one as it's shown in "Last Commit"
			}
			timeAgo := time.Since(commit.Timestamp)
			line := fmt.Sprintf("  %s %-40s (%s ago)\n",
				commit.Hash[:7],
				commit.Message,
				r.formatDuration(timeAgo))
			content.WriteString(line)
		}
	}

	return r.sectionStyle.Render(content.String())
}

func (r *InfoRenderer) renderKeyValue(key, value string) string {
	keyColumn := r.keyStyle.Render(key + ":")
	valueColumn := r.valueStyle.Render(value)

	// Use lipgloss.JoinHorizontal for proper alignment
	return lipgloss.JoinHorizontal(lipgloss.Top, keyColumn, " ", valueColumn) + "\n"
}

func (r *InfoRenderer) formatGitStatus(status *GitStatus) string {
	if status == nil {
		return "🔴 Unknown"
	}

	if status.IsDirty {
		fileCount := status.Modified + status.Staged + status.Untracked
		return fmt.Sprintf("🟡 DIRTY (%d files modified)", fileCount)
	}

	if status.Ahead > 0 || status.Behind > 0 {
		return fmt.Sprintf("🟠 DIVERGED (↑%d ↓%d)", status.Ahead, status.Behind)
	}

	return "🟢 CLEAN"
}

func (r *InfoRenderer) formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%d minutes", int(d.Minutes()))
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	} else {
		return fmt.Sprintf("%d days", int(d.Hours()/24))
	}
}

func (r *InfoRenderer) getStatusIcon(status string) string {
	switch status {
	case "A":
		return "A"
	case "M":
		return "M"
	case "D":
		return "D"
	default:
		return "?"
	}
}

func (r *InfoRenderer) formatPriority(priority string) string {
	lowercasePriority := strings.ToLower(priority)
	switch lowercasePriority {
	case "critical", "highest":
		return "🔴 Critical"
	case "high":
		return "🟠 High"
	case "medium":
		return "🟡 Medium"
	case "low":
		return "🟢 Low"
	case "lowest":
		return "🔵 Lowest"
	default:
		return priority
	}
}

func (r *InfoRenderer) wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	if len(text) <= width {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		// If a single word is longer than width, just add it on its own line
		if len(word) > width {
			if currentLine != "" {
				result.WriteString(currentLine + "\n")
			}
			result.WriteString(word + "\n")
			currentLine = ""
			continue
		}

		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine + "\n")
			currentLine = word
		}
	}

	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}
