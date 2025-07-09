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

// Global icon manager instance
var globalIconManager *IconManager

// IconManager manages configurable icons
type IconManager struct {
	config *Config
}

// NewIconManager creates a new icon manager with configuration
func NewIconManager(config *Config) *IconManager {
	return &IconManager{config: config}
}

// SetGlobalIconManager sets the global icon manager instance
func SetGlobalIconManager(manager *IconManager) {
	globalIconManager = manager
}

// GetGlobalIconManager returns the global icon manager instance
func GetGlobalIconManager() *IconManager {
	if globalIconManager == nil {
		// Return a default manager if none is set
		return NewIconManager(DefaultConfig())
	}
	return globalIconManager
}

// Icon getter methods
func (im *IconManager) Success() string     { return im.config.Icons.Success }
func (im *IconManager) Warning() string     { return im.config.Icons.Warning }
func (im *IconManager) Error() string       { return im.config.Icons.Error }
func (im *IconManager) Info() string        { return im.config.Icons.Info }
func (im *IconManager) Orphaned() string    { return im.config.Icons.Orphaned }
func (im *IconManager) DryRun() string      { return im.config.Icons.DryRun }
func (im *IconManager) Missing() string     { return im.config.Icons.Missing }
func (im *IconManager) Changes() string     { return im.config.Icons.Changes }
func (im *IconManager) GitClean() string    { return im.config.Icons.GitClean }
func (im *IconManager) GitDirty() string    { return im.config.Icons.GitDirty }
func (im *IconManager) GitAhead() string    { return im.config.Icons.GitAhead }
func (im *IconManager) GitBehind() string   { return im.config.Icons.GitBehind }
func (im *IconManager) GitDiverged() string { return im.config.Icons.GitDiverged }
func (im *IconManager) GitUnknown() string  { return im.config.Icons.GitUnknown }

var (
	// Colors
	primaryColor = lipgloss.Color("#7D56F4")
	successColor = lipgloss.Color("#04B575")
	warningColor = lipgloss.Color("#F59E0B")
	errorColor   = lipgloss.Color("#EF4444")
	infoColor    = lipgloss.Color("#3B82F6")
	subtleColor  = lipgloss.Color("#6B7280")

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
	iconManager := GetGlobalIconManager()

	switch icon {
	case iconManager.Success():
		return StatusOKStyle.Render(icon) + " " + text
	case iconManager.Warning():
		return StatusWarningStyle.Render(icon) + " " + text
	case iconManager.Error():
		return StatusErrorStyle.Render(icon) + " " + text
	case iconManager.Orphaned():
		return StatusErrorStyle.Render(icon) + " " + text
	case iconManager.Info():
		return StatusInfoStyle.Render(icon) + " " + text
	case iconManager.DryRun():
		return StatusInfoStyle.Render(icon) + " " + text
	case iconManager.Missing():
		return StatusWarningStyle.Render(icon) + " " + text
	case iconManager.Changes():
		return StatusWarningStyle.Render(icon) + " " + text
	default:
		return icon + " " + text
	}
}

// Helper functions for common icons
func FormatSuccess(text string) string {
	iconManager := GetGlobalIconManager()
	return FormatStatusIcon(iconManager.Success(), text)
}

func FormatWarning(text string) string {
	iconManager := GetGlobalIconManager()
	return FormatStatusIcon(iconManager.Warning(), text)
}

func FormatError(text string) string {
	iconManager := GetGlobalIconManager()
	return FormatStatusIcon(iconManager.Error(), text)
}

func FormatInfo(text string) string {
	iconManager := GetGlobalIconManager()
	return FormatStatusIcon(iconManager.Info(), text)
}

// Git status formatting
func FormatGitStatus(status *GitStatus) string {
	iconManager := GetGlobalIconManager()

	if status == nil {
		return StatusInfoStyle.Render(iconManager.GitUnknown())
	}

	if status.IsDirty {
		return StatusWarningStyle.Render(iconManager.GitDirty())
	}

	if status.Ahead > 0 && status.Behind > 0 {
		return StatusErrorStyle.Render(iconManager.GitDiverged())
	}

	if status.Ahead > 0 {
		return StatusInfoStyle.Render(iconManager.GitAhead())
	}

	if status.Behind > 0 {
		return StatusWarningStyle.Render(iconManager.GitBehind())
	}

	return StatusOKStyle.Render(iconManager.GitClean())
}

// Time formatting utilities
func FormatDuration(d time.Duration) string {
	if d < time.Hour {
		minutes := int(d.Minutes())
		if minutes < 1 {
			return "just now"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}

	if d < 24*time.Hour {
		hours := int(d.Hours())
		return fmt.Sprintf("%d hours", hours)
	}

	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func FormatRelativeTime(t time.Time) string {
	duration := time.Since(t)
	formatted := FormatDuration(duration)

	if formatted == "just now" {
		return formatted
	}
	return formatted + " ago"
}

// Terminal utilities
// GetTerminalWidth returns the terminal width, with multiple fallbacks
func GetTerminalWidth() int {
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
