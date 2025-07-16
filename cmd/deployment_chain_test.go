package cmd

import (
	"testing"

	"gbm/internal"
	"github.com/stretchr/testify/assert"
)

func TestBuildDeploymentChain(t *testing.T) {
	// Create a test configuration
	config := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				MergeInto:   "", // Final branch
				Description: "Main branch",
			},
			"preview": {
				Branch:      "preview",
				MergeInto:   "main",
				Description: "Preview branch",
			},
			"production": {
				Branch:      "production",
				MergeInto:   "preview",
				Description: "Production branch",
			},
		},
	}

	tests := []struct {
		name       string
		baseBranch string
		expected   []string
	}{
		{
			name:       "production branch shows full chain",
			baseBranch: "production",
			expected:   []string{"production", "preview", "main"},
		},
		{
			name:       "preview branch shows partial chain",
			baseBranch: "preview",
			expected:   []string{"preview", "main"},
		},
		{
			name:       "main branch shows only itself",
			baseBranch: "main",
			expected:   []string{"main"},
		},
		{
			name:       "unknown branch shows only itself",
			baseBranch: "unknown",
			expected:   []string{"unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildMergeChain(tt.baseBranch, config)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindMergeIntoTarget(t *testing.T) {
	config := &internal.GBMConfig{
		Worktrees: map[string]internal.WorktreeConfig{
			"main": {
				Branch:      "main",
				MergeInto:   "",
				Description: "Main branch",
			},
			"preview": {
				Branch:      "preview",
				MergeInto:   "main",
				Description: "Preview branch",
			},
			"production": {
				Branch:      "production",
				MergeInto:   "preview",
				Description: "Production branch",
			},
		},
	}

	tests := []struct {
		name         string
		sourceBranch string
		expected     string
	}{
		{
			name:         "production merges into preview",
			sourceBranch: "production",
			expected:     "preview",
		},
		{
			name:         "preview merges into main",
			sourceBranch: "preview",
			expected:     "main",
		},
		{
			name:         "main has no merge target",
			sourceBranch: "main",
			expected:     "",
		},
		{
			name:         "unknown branch has no merge target",
			sourceBranch: "unknown",
			expected:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMergeIntoTarget(tt.sourceBranch, config)
			assert.Equal(t, tt.expected, result)
		})
	}
}