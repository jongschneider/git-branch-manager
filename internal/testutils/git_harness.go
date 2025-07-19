package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type GitTestRepo struct {
	RemoteDir string
	LocalDir  string
	TempDir   string
	RepoName  string
	Config    RepoConfig
	t         *testing.T
}

type RepoConfig struct {
	DefaultBranch string
	UserName      string
	UserEmail     string
	RemoteName    string
}

// YAML-based configuration structures for GBM config generation
type gbmConfig struct {
	Worktrees map[string]WorktreeConfig `yaml:"worktrees"`
}

type WorktreeConfig struct {
	Branch      string `yaml:"branch"`
	MergeInto   string `yaml:"merge_into,omitempty"`
	Description string `yaml:"description,omitempty"`
}

var defaultConfig = RepoConfig{
	DefaultBranch: "main",
	UserName:      "Test User",
	UserEmail:     "test@example.com",
	RemoteName:    "origin",
}

func NewGitTestRepo(t *testing.T, opts ...RepoOption) *GitTestRepo {
	tempDir := t.TempDir()
	remoteDir := filepath.Join(tempDir, "remote.git")

	repo := &GitTestRepo{
		TempDir:   tempDir,
		RemoteDir: remoteDir,
		LocalDir:  filepath.Join(tempDir, "local"),
		RepoName:  extractRepoName(remoteDir),
		Config:    defaultConfig,
		t:         t,
	}

	for _, opt := range opts {
		opt(repo)
	}

	if err := repo.setupBareRemote(); err != nil {
		t.Fatalf("Failed to setup bare remote: %v", err)
	}

	if err := repo.setupLocalRepo(); err != nil {
		t.Fatalf("Failed to setup local repo: %v", err)
	}

	if err := repo.configureGitUser(); err != nil {
		t.Fatalf("Failed to configure git user: %v", err)
	}

	if err := repo.createInitialCommit(); err != nil {
		t.Fatalf("Failed to create initial commit: %v", err)
	}

	return repo
}

func (r *GitTestRepo) setupBareRemote() error {
	if err := os.MkdirAll(r.RemoteDir, 0o755); err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	cmd := exec.Command("git", "init", "--bare")
	cmd.Dir = r.RemoteDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize bare repository: %w", err)
	}

	// Set the default branch for the remote repository
	cmd = exec.Command("git", "symbolic-ref", "HEAD", "refs/heads/"+r.Config.DefaultBranch)
	cmd.Dir = r.RemoteDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set default branch: %w", err)
	}

	return nil
}

func (r *GitTestRepo) setupLocalRepo() error {
	cmd := exec.Command("git", "clone", r.RemoteDir, r.LocalDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func (r *GitTestRepo) configureGitUser() error {
	if err := r.runGitCommand("config", "user.name", r.Config.UserName); err != nil {
		return fmt.Errorf("failed to configure git user name: %w", err)
	}

	if err := r.runGitCommand("config", "user.email", r.Config.UserEmail); err != nil {
		return fmt.Errorf("failed to configure git user email: %w", err)
	}

	return nil
}

func (r *GitTestRepo) createInitialCommit() error {
	readmePath := filepath.Join(r.LocalDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Repository\n"), 0o644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	if err := r.runGitCommand("add", "README.md"); err != nil {
		return fmt.Errorf("failed to add README.md: %w", err)
	}

	if err := r.runGitCommand("commit", "-m", "Initial commit"); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	if err := r.runGitCommand("branch", "-M", r.Config.DefaultBranch); err != nil {
		return fmt.Errorf("failed to rename branch to %s: %w", r.Config.DefaultBranch, err)
	}

	if err := r.runGitCommand("push", r.Config.RemoteName, r.Config.DefaultBranch); err != nil {
		return fmt.Errorf("failed to push initial commit: %w", err)
	}

	return nil
}

func (r *GitTestRepo) runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = r.LocalDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git command failed: %s (output: %s)", err, string(output))
	} else {
		fmt.Println("Git command output:", string(output))
	}
	return nil
}

func (r *GitTestRepo) runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = r.LocalDir
	return cmd.CombinedOutput()
}

func (r *GitTestRepo) Cleanup() {
}

func (r *GitTestRepo) GetLocalPath() string {
	return r.LocalDir
}

func (r *GitTestRepo) GetRemotePath() string {
	return r.RemoteDir
}

func (r *GitTestRepo) GetRepoName() string {
	return r.RepoName
}

func (r *GitTestRepo) CreateBranch(name, content string) error {
	if err := r.runGitCommand("checkout", "-b", name); err != nil {
		return fmt.Errorf("failed to create branch %s: %w", name, err)
	}

	contentPath := filepath.Join(r.LocalDir, "content.txt")
	if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write content file: %w", err)
	}

	if err := r.runGitCommand("add", "content.txt"); err != nil {
		return fmt.Errorf("failed to add content file: %w", err)
	}

	if err := r.runGitCommand("commit", "-m", fmt.Sprintf("Add content for %s", name)); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	if err := r.runGitCommand("push", r.Config.RemoteName, name); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", name, err)
	}

	if err := r.runGitCommand("checkout", r.Config.DefaultBranch); err != nil {
		return fmt.Errorf("failed to return to default branch: %w", err)
	}

	return nil
}

func (r *GitTestRepo) CreateBranchFrom(name, baseBranch, content string) error {
	if err := r.runGitCommand("checkout", baseBranch); err != nil {
		return fmt.Errorf("failed to checkout base branch %s: %w", baseBranch, err)
	}

	if err := r.runGitCommand("checkout", "-b", name); err != nil {
		return fmt.Errorf("failed to create branch %s from %s: %w", name, baseBranch, err)
	}

	contentPath := filepath.Join(r.LocalDir, "content.txt")
	if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write content file: %w", err)
	}

	if err := r.runGitCommand("add", "content.txt"); err != nil {
		return fmt.Errorf("failed to add content file: %w", err)
	}

	if err := r.runGitCommand("commit", "-m", fmt.Sprintf("Add content for %s", name)); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	if err := r.runGitCommand("push", r.Config.RemoteName, name); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", name, err)
	}

	if err := r.runGitCommand("checkout", r.Config.DefaultBranch); err != nil {
		return fmt.Errorf("failed to return to default branch: %w", err)
	}

	return nil
}

func (r *GitTestRepo) SwitchToBranch(name string) error {
	if err := r.runGitCommand("checkout", name); err != nil {
		return fmt.Errorf("failed to switch to branch %s: %w", name, err)
	}
	return nil
}

func (r *GitTestRepo) WriteFile(path, content string) error {
	fullPath := filepath.Join(r.LocalDir, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", path, err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

func (r *GitTestRepo) CreateGBMConfig(worktrees map[string]WorktreeConfig) error {
	config := &gbmConfig{
		Worktrees: worktrees,
	}

	// Marshal the config to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal GBM config to YAML: %w", err)
	}

	// Add a header comment to the YAML
	content := "# Git Branch Manager Configuration\n\n# Worktree definitions - key is the worktree name, value defines the branch and merge strategy\n" + string(yamlData)

	fmt.Println("🚨", content)
	return r.WriteFile("gbm.branchconfig.yaml", content)
}

func (r *GitTestRepo) CommitChanges(message string) error {
	if err := r.runGitCommand("add", "."); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	if err := r.runGitCommand("commit", "-m", message); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

func (r *GitTestRepo) CommitChangesWithForceAdd(message string) error {
	if err := r.runGitCommand("add", "-f", "."); err != nil {
		return fmt.Errorf("failed to add changes: %w", err)
	}

	if err := r.runGitCommand("commit", "-m", message); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	return nil
}

func (r *GitTestRepo) PushBranch(branch string) error {
	if err := r.runGitCommand("push", r.Config.RemoteName, branch); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", branch, err)
	}
	return nil
}

func (r *GitTestRepo) ConvertToBare() string {
	return r.RemoteDir
}

func (r *GitTestRepo) ListBranches() ([]string, error) {
	output, err := r.runCommand("git", "branch", "-r")
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var branches []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "HEAD") {
			branch := strings.TrimPrefix(line, r.Config.RemoteName+"/")
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

func (r *GitTestRepo) WithWorkingDir(dir string) func() {
	originalDir, _ := os.Getwd()
	os.Chdir(dir)
	return func() {
		os.Chdir(originalDir)
	}
}

func (r *GitTestRepo) InLocalRepo(fn func() error) error {
	restore := r.WithWorkingDir(r.LocalDir)
	defer restore()
	return fn()
}

func (r *GitTestRepo) CreateSynchronizedBranch(name string) error {
	if err := r.runGitCommand("checkout", "-b", name); err != nil {
		return fmt.Errorf("failed to create branch %s: %w", name, err)
	}

	if err := r.runGitCommand("push", r.Config.RemoteName, name); err != nil {
		return fmt.Errorf("failed to push branch %s: %w", name, err)
	}

	return nil
}

// extractRepoName extracts the repository name from a URL or path
func extractRepoName(repoUrl string) string {
	// Remove .git suffix if present
	url := strings.TrimSuffix(repoUrl, ".git")
	// Extract the last part of the URL (repository name)
	parts := strings.Split(url, "/")
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}
	return "repository"
}
