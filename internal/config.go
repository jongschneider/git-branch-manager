package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Settings ConfigSettings `toml:"settings"`
	State    ConfigState    `toml:"state"`
	Icons    ConfigIcons    `toml:"icons"`
	Jira     ConfigJira     `toml:"jira"`
	FileCopy ConfigFileCopy `toml:"file_copy"`
}

type ConfigSettings struct {
	WorktreePrefix        string `toml:"worktree_prefix"`
	AutoFetch             bool   `toml:"auto_fetch"`
	CreateMissingBranches bool   `toml:"create_missing_branches"`
	MergeBackAlerts       bool   `toml:"merge_back_alerts"`
}

type FileCopyRule struct {
	SourceWorktree string   `toml:"source_worktree"`
	Files          []string `toml:"files"`
}

type ConfigFileCopy struct {
	Rules []FileCopyRule `toml:"rules"`
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

type ConfigJira struct {
	Me string `toml:"me"`
}

// YAML-based configuration structures
type GBMConfig struct {
	Worktrees map[string]WorktreeConfig `yaml:"worktrees"`
}

type WorktreeConfig struct {
	Branch      string `yaml:"branch"`
	MergeInto   string `yaml:"merge_into,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Settings: ConfigSettings{
			WorktreePrefix:        "worktrees",
			AutoFetch:             true,
			CreateMissingBranches: false,
			MergeBackAlerts:       false, // Disabled by default
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
		Jira: ConfigJira{
			Me: "", // Will be populated when first used
		},
		FileCopy: ConfigFileCopy{
			Rules: []FileCopyRule{},
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

// GetGBMDir returns the path to the .gbm directory for the given repository root
func GetGBMDir(repoRoot string) string {
	return filepath.Join(repoRoot, ".gbm")
}

func (c *Config) Save(gbmDir string) error {
	if err := os.MkdirAll(gbmDir, 0o755); err != nil {
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

// ParseGBMConfig parses the YAML-based .gbm.config.yaml file
func ParseGBMConfig(path string) (*GBMConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read .gbm.config.yaml file: %w", err)
	}

	var config GBMConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
}
