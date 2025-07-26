package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

const (
	// Directory and file names
	DefaultWorktreeDirname      = "worktrees"
	DefaultBranchConfigFilename = "gbm.branchconfig.yaml"
	DefaultConfigDirname        = ".gbm"
	DefaultConfigFilename       = "config.toml"
	DefaultStateFilename        = "state.toml"
)

type Config struct {
	Settings ConfigSettings `toml:"settings"`
	Icons    ConfigIcons    `toml:"icons"`
	Jira     ConfigJira     `toml:"jira"`
	FileCopy ConfigFileCopy `toml:"file_copy"`
}

type ConfigSettings struct {
	WorktreePrefix              string        `toml:"worktree_prefix"`
	AutoFetch                   bool          `toml:"auto_fetch"`
	CreateMissingBranches       bool          `toml:"create_missing_branches"`
	MergeBackAlerts             bool          `toml:"merge_back_alerts"`
	HotfixPrefix                string        `toml:"hotfix_prefix"`
	MergebackPrefix             string        `toml:"mergeback_prefix"`
	MergeBackCheckInterval      time.Duration `toml:"merge_back_check_interval"`
	MergeBackUserCommitInterval time.Duration `toml:"merge_back_user_commit_interval"`
	CandidateBranches           []string      `toml:"candidate_branches"`
}

type FileCopyRule struct {
	SourceWorktree string   `toml:"source_worktree"`
	Files          []string `toml:"files"`
}

type ConfigFileCopy struct {
	Rules []FileCopyRule `toml:"rules"`
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

	// Section header icons
	WorktreeHeader string `toml:"worktree_header"`
	JiraHeader     string `toml:"jira_header"`
	GitHeader      string `toml:"git_header"`
}

type ConfigJira struct {
	Me string `toml:"me"`
}

// YAML-based configuration structures
type GBMConfig struct {
	Worktrees map[string]WorktreeConfig `yaml:"worktrees"`
	Tree      *WorktreeManager          `yaml:"-"`
}

type WorktreeConfig struct {
	Branch      string `yaml:"branch"`
	MergeInto   string `yaml:"merge_into,omitempty"`
	Description string `yaml:"description,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Settings: ConfigSettings{
			WorktreePrefix:              DefaultWorktreeDirname,
			AutoFetch:                   true,
			CreateMissingBranches:       false,
			MergeBackAlerts:             true,                                         // Enabled by default
			HotfixPrefix:                "HOTFIX",                                     // Default hotfix prefix
			MergebackPrefix:             "MERGE",                                      // Default mergeback prefix
			MergeBackCheckInterval:      3 * time.Hour,                                // Check every 3 hours by default
			MergeBackUserCommitInterval: 30 * time.Minute,                             // Alert every 30 minutes when user has commits
			CandidateBranches:           []string{"main", "master", "develop", "dev"}, // Default candidate branches
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

			// Section header icons
			WorktreeHeader: "📁",
			JiraHeader:     "🎫",
			GitHeader:      "🌿",
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
	configPath := filepath.Join(gbmDir, DefaultConfigFilename)

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

	configPath := filepath.Join(gbmDir, DefaultConfigFilename)
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer func() { _ = file.Close() }()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	return nil
}

// ParseGBMConfig parses the YAML-based branch config file
func ParseGBMConfig(path string) (*GBMConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s file: %w", DefaultBranchConfigFilename, err)
	}

	var config GBMConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	// Initialize the tree structure
	tree, err := NewWorktreeManager(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to build worktree tree: %w", err)
	}
	config.Tree = tree

	return &config, nil
}
