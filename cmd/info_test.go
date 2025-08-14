package cmd

import (
	"errors"
	"testing"
	"time"

	"gbm/internal"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// UNIT TESTS (Using mocks - these are fast and don't require real git operations)
// ============================================================================
// These tests use the worktreeInfoProvider interface with mocks to test business logic
// without requiring real git operations. They run in milliseconds.

func TestGetWorktreeInfo(t *testing.T) {
	// Sample real data based on INGSVC-5739 from JIRA
	sampleWorktree := &internal.WorktreeInfo{
		Name:   "INGSVC-5739",
		Path:   "/Users/test/worktrees/INGSVC-5739",
		Branch: "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API",
	}

	sampleGitStatus := &internal.GitStatus{
		IsDirty:   true,
		Modified:  1,
		Untracked: 0,
		Staged:    0,
	}

	sampleCommits := []internal.CommitInfo{
		{
			Hash:      "ec2162d",
			Message:   "docs(refinitivlseg): add CFS API spec, OpenAPI, and usage notes; generate initial client via oapi-codegen",
			Author:    "jonathan-schneider-tl",
			Timestamp: time.Now().Add(-18 * time.Hour),
		},
		{
			Hash:      "8b0227a",
			Message:   "WIP - initital code gen",
			Author:    "jonathan-schneider-tl",
			Timestamp: time.Now().Add(-19 * time.Hour),
		},
	}

	sampleFileChanges := []internal.FileChange{
		{
			Path:      "pkg/refinitivlseg/.docs/DDD_plan.md",
			Status:    "?",
			Additions: 45,
			Deletions: 0,
		},
	}

	// Sample JIRA ticket data based on real API response
	sampleJiraTicket := &internal.JiraTicketDetails{
		Key:      "INGSVC-5739",
		Summary:  "New Integration - Refinitiv LSEG Messenger API",
		Status:   "In Dev.",
		Assignee: "Jonathan Schneider (jonathan.schneider@thetalake.com)",
		Priority: "High",
		Reporter: "Kannan Appachi (kannan@thetalake.com)",
		Epic:     "EPIC-2540",
		URL:      "https://thetalake.atlassian.net/browse/INGSVC-5739",
		Created:  time.Date(2025, 8, 5, 4, 28, 27, 0, time.UTC),
		LatestComment: &internal.Comment{
			Author:    "Jonathan Schneider",
			Content:   "Cursor Prompt (GPT-5):",
			Timestamp: time.Date(2025, 8, 13, 8, 49, 35, 0, time.UTC),
		},
	}

	tests := []struct {
		name         string
		mockSetup    func() *worktreeInfoProviderMock
		expectErr    func(t *testing.T, err error)
		expectData   func(t *testing.T, data *internal.WorktreeInfoData)
		worktreeName string
	}{
		{
			name:         "success - get complete worktree info with JIRA",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return []*internal.WorktreeInfo{sampleWorktree}, nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						assert.Equal(t, sampleWorktree.Path, worktreePath)
						return sampleGitStatus, nil
					},
					GetWorktreeCommitHistoryFunc: func(worktreePath string, limit int) ([]internal.CommitInfo, error) {
						assert.Equal(t, sampleWorktree.Path, worktreePath)
						assert.Equal(t, 5, limit)
						return sampleCommits, nil
					},
					GetWorktreeFileChangesFunc: func(worktreePath string) ([]internal.FileChange, error) {
						assert.Equal(t, sampleWorktree.Path, worktreePath)
						return sampleFileChanges, nil
					},
					GetJiraTicketDetailsFunc: func(jiraKey string) (*internal.JiraTicketDetails, error) {
						assert.Equal(t, "INGSVC-5739", jiraKey)
						return sampleJiraTicket, nil
					},
					// Add missing methods for getBaseBranchInfo
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 1, 0, nil
					},
					GetStateFunc: func() *internal.State {
						state := &internal.State{}
						state.WorktreeBaseBranch = map[string]string{
							"INGSVC-5739": "master",
						}
						return state
					},
					GetConfigFunc: func() *internal.Config {
						return &internal.Config{
							Settings: internal.ConfigSettings{
								CandidateBranches: []string{"main", "master", "develop"},
							},
						}
					},
					VerifyWorktreeRefFunc: func(ref string, worktreePath string) (bool, error) {
						// Mock that stored base branches exist
						return true, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.NotNil(t, data)
				assert.Equal(t, "INGSVC-5739", data.Name)
				assert.Equal(t, sampleWorktree.Path, data.Path)
				assert.Equal(t, sampleWorktree.Branch, data.Branch)
				assert.Equal(t, sampleGitStatus, data.GitStatus)
				assert.Equal(t, sampleCommits, data.Commits)
				assert.Equal(t, sampleFileChanges, data.ModifiedFiles)
				assert.Equal(t, sampleJiraTicket, data.JiraTicket)
			},
		},
		{
			name:         "success - get worktree info without JIRA (no JIRA key in name)",
			worktreeName: "feature-branch",
			mockSetup: func() *worktreeInfoProviderMock {
				noJiraWorktree := &internal.WorktreeInfo{
					Name:   "feature-branch",
					Path:   "/Users/test/worktrees/feature-branch",
					Branch: "feature/some-feature",
				}
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return []*internal.WorktreeInfo{noJiraWorktree}, nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						return sampleGitStatus, nil
					},
					GetWorktreeCommitHistoryFunc: func(worktreePath string, limit int) ([]internal.CommitInfo, error) {
						return sampleCommits, nil
					},
					GetWorktreeFileChangesFunc: func(worktreePath string) ([]internal.FileChange, error) {
						return sampleFileChanges, nil
					},
					// Add missing methods for getBaseBranchInfo
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "feature/some-feature", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/feature/some-feature", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 0, 0, nil
					},
					GetStateFunc: func() *internal.State {
						return &internal.State{}
					},
					GetConfigFunc: func() *internal.Config {
						return &internal.Config{
							Settings: internal.ConfigSettings{
								CandidateBranches: []string{"main", "master", "develop"},
							},
						}
					},
					VerifyWorktreeRefFunc: func(ref string, worktreePath string) (bool, error) {
						// Mock that stored base branches exist
						return true, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.NotNil(t, data)
				assert.Equal(t, "feature-branch", data.Name)
				assert.Nil(t, data.JiraTicket) // Should be nil since no JIRA key in name
			},
		},
		{
			name:         "success - JIRA CLI not available",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return []*internal.WorktreeInfo{sampleWorktree}, nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						return sampleGitStatus, nil
					},
					GetWorktreeCommitHistoryFunc: func(worktreePath string, limit int) ([]internal.CommitInfo, error) {
						return sampleCommits, nil
					},
					GetWorktreeFileChangesFunc: func(worktreePath string) ([]internal.FileChange, error) {
						return sampleFileChanges, nil
					},
					GetJiraTicketDetailsFunc: func(jiraKey string) (*internal.JiraTicketDetails, error) {
						assert.Equal(t, "INGSVC-5739", jiraKey)
						return nil, internal.ErrJiraCliNotFound // JIRA CLI not available
					},
					// Add missing methods for getBaseBranchInfo
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 1, 0, nil
					},
					GetStateFunc: func() *internal.State {
						state := &internal.State{}
						state.WorktreeBaseBranch = map[string]string{
							"INGSVC-5739": "master",
						}
						return state
					},
					GetConfigFunc: func() *internal.Config {
						return &internal.Config{
							Settings: internal.ConfigSettings{
								CandidateBranches: []string{"main", "master", "develop"},
							},
						}
					},
					VerifyWorktreeRefFunc: func(ref string, worktreePath string) (bool, error) {
						// Mock that stored base branches exist
						return true, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.NotNil(t, data)
				assert.Equal(t, "INGSVC-5739", data.Name)
				assert.Nil(t, data.JiraTicket) // Should be nil since JIRA CLI not available
			},
		},
		{
			name:         "error - worktree not found",
			worktreeName: "nonexistent-worktree",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return []*internal.WorktreeInfo{sampleWorktree}, nil // Different worktree
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "worktree 'nonexistent-worktree' not found")
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.Nil(t, data)
			},
		},
		{
			name:         "error - get worktrees fails",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return nil, errors.New("git worktree list failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get worktrees")
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.Nil(t, data)
			},
		},
		{
			name:         "partial success - git status fails but continues",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreesFunc: func() ([]*internal.WorktreeInfo, error) {
						return []*internal.WorktreeInfo{sampleWorktree}, nil
					},
					GetWorktreeStatusFunc: func(worktreePath string) (*internal.GitStatus, error) {
						return nil, errors.New("git status failed")
					},
					GetWorktreeCommitHistoryFunc: func(worktreePath string, limit int) ([]internal.CommitInfo, error) {
						return sampleCommits, nil
					},
					GetWorktreeFileChangesFunc: func(worktreePath string) ([]internal.FileChange, error) {
						return sampleFileChanges, nil
					},
					GetJiraTicketDetailsFunc: func(jiraKey string) (*internal.JiraTicketDetails, error) {
						assert.Equal(t, "INGSVC-5739", jiraKey)
						return nil, internal.ErrJiraCliNotFound // Simulate JIRA CLI not available
					},
					// Add missing methods for getBaseBranchInfo
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 1, 0, nil
					},
					GetStateFunc: func() *internal.State {
						state := &internal.State{}
						state.WorktreeBaseBranch = map[string]string{
							"INGSVC-5739": "master",
						}
						return state
					},
					GetConfigFunc: func() *internal.Config {
						return &internal.Config{
							Settings: internal.ConfigSettings{
								CandidateBranches: []string{"main", "master", "develop"},
							},
						}
					},
					VerifyWorktreeRefFunc: func(ref string, worktreePath string) (bool, error) {
						// Mock that stored base branches exist
						return true, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err) // Should not fail, just continue with nil git status
			},
			expectData: func(t *testing.T, data *internal.WorktreeInfoData) {
				assert.NotNil(t, data)
				assert.Equal(t, "INGSVC-5739", data.Name)
				assert.Nil(t, data.GitStatus)  // Should be nil due to error
				assert.NotNil(t, data.Commits) // Other data should still be present
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := tt.mockSetup()

			data, err := getWorktreeInfo(provider, tt.worktreeName)

			tt.expectErr(t, err)
			tt.expectData(t, data)
		})
	}
}

func TestGetBaseBranchInfo(t *testing.T) {
	sampleConfig := &internal.Config{
		Settings: internal.ConfigSettings{
			CandidateBranches: []string{"main", "master", "develop"},
		},
	}

	sampleState := &internal.State{}
	// Set up stored base branch
	sampleState.WorktreeBaseBranch = map[string]string{
		"INGSVC-5739": "master",
	}

	tests := []struct {
		name         string
		worktreePath string
		worktreeName string
		mockSetup    func() *worktreeInfoProviderMock
		expectErr    func(t *testing.T, err error)
		expectData   func(t *testing.T, data *internal.BranchInfo)
	}{
		{
			name:         "success - get branch info with stored base branch",
			worktreePath: "/Users/test/worktrees/INGSVC-5739",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 1, 0, nil
					},
					GetStateFunc: func() *internal.State {
						return sampleState
					},
					GetConfigFunc: func() *internal.Config {
						return sampleConfig
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expectData: func(t *testing.T, data *internal.BranchInfo) {
				assert.NotNil(t, data)
				assert.Equal(t, "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", data.Upstream)
				assert.Equal(t, 1, data.AheadBy)
				assert.Equal(t, 0, data.BehindBy)
			},
		},
		{
			name:         "success - fallback to candidate branches when no stored base",
			worktreePath: "/Users/test/worktrees/feature-branch",
			worktreeName: "feature-branch",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "feature/some-feature", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/feature/some-feature", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 2, 1, nil
					},
					GetStateFunc: func() *internal.State {
						return &internal.State{} // No stored base branches
					},
					GetConfigFunc: func() *internal.Config {
						return sampleConfig
					},
					VerifyWorktreeRefFunc: func(ref string, worktreePath string) (bool, error) {
						// Simulate that "main" exists and is a valid base
						if ref == "main" {
							return true, nil
						}
						return false, nil
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
			expectData: func(t *testing.T, data *internal.BranchInfo) {
				assert.NotNil(t, data)
				assert.Equal(t, "origin/feature/some-feature", data.Upstream)
				assert.Equal(t, 2, data.AheadBy)
				assert.Equal(t, 1, data.BehindBy)
			},
		},
		{
			name:         "error - get current branch fails",
			worktreePath: "/Users/test/worktrees/INGSVC-5739",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "", errors.New("git branch failed")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "git branch failed")
			},
			expectData: func(t *testing.T, data *internal.BranchInfo) {
				assert.Nil(t, data)
			},
		},
		{
			name:         "error - get upstream branch fails",
			worktreePath: "/Users/test/worktrees/INGSVC-5739",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "", errors.New("no upstream branch")
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to get upstream branch")
			},
			expectData: func(t *testing.T, data *internal.BranchInfo) {
				assert.Nil(t, data)
			},
		},
		{
			name:         "partial success - ahead/behind count fails but continues",
			worktreePath: "/Users/test/worktrees/INGSVC-5739",
			worktreeName: "INGSVC-5739",
			mockSetup: func() *worktreeInfoProviderMock {
				return &worktreeInfoProviderMock{
					GetWorktreeCurrentBranchFunc: func(worktreePath string) (string, error) {
						return "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeUpstreamBranchFunc: func(worktreePath string) (string, error) {
						return "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", nil
					},
					GetWorktreeAheadBehindCountFunc: func(worktreePath string) (int, int, error) {
						return 0, 0, errors.New("ahead/behind count failed")
					},
					GetStateFunc: func() *internal.State {
						return sampleState
					},
					GetConfigFunc: func() *internal.Config {
						return sampleConfig
					},
				}
			},
			expectErr: func(t *testing.T, err error) {
				assert.NoError(t, err) // Should not fail, just use 0,0 for ahead/behind
			},
			expectData: func(t *testing.T, data *internal.BranchInfo) {
				assert.NotNil(t, data)
				assert.Equal(t, "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API", data.Upstream)
				assert.Equal(t, 0, data.AheadBy)  // Should be 0 due to error
				assert.Equal(t, 0, data.BehindBy) // Should be 0 due to error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := tt.mockSetup()

			data, err := getBaseBranchInfo(tt.worktreePath, tt.worktreeName, provider)

			tt.expectErr(t, err)
			tt.expectData(t, data)
		})
	}
}

// Test helper function to create a sample WorktreeInfoData for testing display logic
func createSampleWorktreeInfoData() *internal.WorktreeInfoData {
	return &internal.WorktreeInfoData{
		Name:      "INGSVC-5739",
		Path:      "/Users/test/worktrees/INGSVC-5739",
		Branch:    "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API",
		CreatedAt: time.Date(2025, 8, 13, 17, 11, 33, 0, time.UTC),
		GitStatus: &internal.GitStatus{
			IsDirty:  true,
			Modified: 1,
		},
		BaseInfo: &internal.BranchInfo{
			Name:     "bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API",
			Upstream: "origin/bug/INGSVC-5739_New_Integration_Refinitiv_LSEG_Messenger_API",
			AheadBy:  1,
			BehindBy: 0,
		},
		Commits: []internal.CommitInfo{
			{
				Hash:      "ec2162d",
				Message:   "docs(refinitivlseg): add CFS API spec, OpenAPI, and usage notes; generate initial client via oapi-codegen",
				Author:    "jonathan-schneider-tl",
				Timestamp: time.Now().Add(-18 * time.Hour),
			},
		},
		ModifiedFiles: []internal.FileChange{
			{
				Path:      "pkg/refinitivlseg/.docs/DDD_plan.md",
				Status:    "?",
				Additions: 45,
				Deletions: 0,
			},
		},
		JiraTicket: &internal.JiraTicketDetails{
			Key:      "INGSVC-5739",
			Summary:  "New Integration - Refinitiv LSEG Messenger API",
			Status:   "In Dev.",
			Priority: "High",
			URL:      "https://thetalake.atlassian.net/browse/INGSVC-5739",
		},
	}
}

func TestDisplayWorktreeInfo(t *testing.T) {
	// This function mainly delegates to InfoRenderer, so we just test that it doesn't panic
	// and handles nil config gracefully

	data := createSampleWorktreeInfoData()

	t.Run("success - display with valid config", func(t *testing.T) {
		config := internal.DefaultConfig()

		// Should not panic
		assert.NotPanics(t, func() {
			displayWorktreeInfo(data, config)
		})
	})

	t.Run("success - display with nil config", func(t *testing.T) {
		// Should not panic and should handle nil config by using default
		assert.NotPanics(t, func() {
			displayWorktreeInfo(data, nil)
		})
	})
}
