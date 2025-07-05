package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Settings ConfigSettings `toml:"settings"`
	State    ConfigState    `toml:"state"`
	Icons    ConfigIcons    `toml:"icons"`
}

type ConfigSettings struct {
	WorktreePrefix        string `toml:"worktree_prefix"`
	AutoFetch             bool   `toml:"auto_fetch"`
	CreateMissingBranches bool   `toml:"create_missing_branches"`
}

type ConfigState struct {
	LastSync         time.Time `toml:"last_sync"`
	TrackedVars      []string  `toml:"tracked_vars"`
	AdHocWorktrees   []string  `toml:"ad_hoc_worktrees"`
	CurrentWorktree  string    `toml:"current_worktree"`
	PreviousWorktree string    `toml:"previous_worktree"`
}

type ConfigIcons struct {
	// Status icons
	Success  string `toml:"success"`
	Warning  string `toml:"warning"`
	Error    string `toml:"error"`
	Info     string `toml:"info"`
	Orphaned string `toml:"orphaned"`
	DryRun   string `toml:"dry_run"`
	Missing  string `toml:"missing"`
	Changes  string `toml:"changes"`

	// Git status icons
	GitClean    string `toml:"git_clean"`
	GitDirty    string `toml:"git_dirty"`
	GitAhead    string `toml:"git_ahead"`
	GitBehind   string `toml:"git_behind"`
	GitDiverged string `toml:"git_diverged"`
	GitUnknown  string `toml:"git_unknown"`
}

type EnvMapping struct {
	Variables map[string]string
}

func DefaultConfig() *Config {
	return &Config{
		Settings: ConfigSettings{
			WorktreePrefix:        "worktrees",
			AutoFetch:             true,
			CreateMissingBranches: false,
		},
		State: ConfigState{
			LastSync:         time.Time{},
			TrackedVars:      []string{},
			AdHocWorktrees:   []string{},
			CurrentWorktree:  "",
			PreviousWorktree: "",
		},
		Icons: ConfigIcons{
			// Status icons
			Success:  "✅",
			Warning:  "⚠️",
			Error:    "❌",
			Info:     "💡",
			Orphaned: "🗑️",
			DryRun:   "🔍",
			Missing:  "📁",
			Changes:  "🔄",

			// Git status icons
			GitClean:    "✓",
			GitDirty:    "~",
			GitAhead:    "↑",
			GitBehind:   "↓",
			GitDiverged: "⇕",
			GitUnknown:  "?",
		},
	}
}

func LoadConfig(gbmDir string) (*Config, error) {
	configPath := filepath.Join(gbmDir, "config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}

func (c *Config) Save(gbmDir string) error {
	if err := os.MkdirAll(gbmDir, 0755); err != nil {
		return fmt.Errorf("failed to create .gbm directory: %w", err)
	}

	configPath := filepath.Join(gbmDir, "config.toml")
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

func ParseEnvrc(path string) (*EnvMapping, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open .envrc file: %w", err)
	}
	defer file.Close()

	mapping := &EnvMapping{
		Variables: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	envVarRegex := regexp.MustCompile(`^([A-Z_][A-Z0-9_]*)=(.+)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := envVarRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			varName := matches[1]
			branchName := strings.Trim(matches[2], "\"'")
			mapping.Variables[varName] = branchName
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading .envrc file: %w", err)
	}

	return mapping, nil
}
