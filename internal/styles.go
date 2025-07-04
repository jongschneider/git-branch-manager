package internal

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	successColor   = lipgloss.Color("#04B575")
	warningColor   = lipgloss.Color("#F59E0B")
	errorColor     = lipgloss.Color("#EF4444")
	infoColor      = lipgloss.Color("#3B82F6")
	subtleColor    = lipgloss.Color("#6B7280")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
		MarginLeft(0).
		MarginRight(0)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		MarginBottom(1)

	SubHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(subtleColor)

	// Status styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true)

	InfoStyle = lipgloss.NewStyle().
		Foreground(infoColor)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Align(lipgloss.Left)

	TableBorderStyle = lipgloss.NewStyle().
		Foreground(subtleColor)

	// Status icon styles
	StatusOKStyle = lipgloss.NewStyle().
		Foreground(successColor)

	StatusWarningStyle = lipgloss.NewStyle().
		Foreground(warningColor)

	StatusErrorStyle = lipgloss.NewStyle().
		Foreground(errorColor)

	StatusInfoStyle = lipgloss.NewStyle().
		Foreground(infoColor)

	// Utility styles
	BoldStyle = lipgloss.NewStyle().
		Bold(true)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(subtleColor)

	// Message styles
	VerboseStyle = lipgloss.NewStyle().
		Foreground(subtleColor).
		Italic(true)

	PromptStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true)
)

// Status formatting functions
func FormatSuccess(text string) string {
	return SuccessStyle.Render(text)
}

func FormatWarning(text string) string {
	return WarningStyle.Render(text)
}

func FormatError(text string) string {
	return ErrorStyle.Render(text)
}

func FormatInfo(text string) string {
	return InfoStyle.Render(text)
}

func FormatVerbose(text string) string {
	return VerboseStyle.Render(text)
}

func FormatHeader(text string) string {
	return HeaderStyle.Render(text)
}

func FormatSubHeader(text string) string {
	return SubHeaderStyle.Render(text)
}

func FormatBold(text string) string {
	return BoldStyle.Render(text)
}

func FormatSubtle(text string) string {
	return SubtleStyle.Render(text)
}

func FormatPrompt(text string) string {
	return PromptStyle.Render(text)
}

// Status icon formatting with consistent styling
func FormatStatusIcon(icon, text string) string {
	switch icon {
	case "✅":
		return StatusOKStyle.Render(icon) + " " + text
	case "⚠️":
		return StatusWarningStyle.Render(icon) + " " + text
	case "❌":
		return StatusErrorStyle.Render(icon) + " " + text
	case "🗑️":
		return StatusErrorStyle.Render(icon) + " " + text
	case "💡":
		return StatusInfoStyle.Render(icon) + " " + text
	default:
		return icon + " " + text
	}
}

// Git status formatting
func FormatGitStatus(status *GitStatus) string {
	if status == nil {
		return StatusInfoStyle.Render("?")
	}

	if status.IsDirty {
		return StatusWarningStyle.Render("~")
	}

	if status.Ahead > 0 && status.Behind > 0 {
		return StatusErrorStyle.Render("⇕")
	}

	if status.Ahead > 0 {
		return StatusInfoStyle.Render("↑")
	}

	if status.Behind > 0 {
		return StatusWarningStyle.Render("↓")
	}

	return StatusOKStyle.Render("✓")
}