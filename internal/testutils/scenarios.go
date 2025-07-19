package testutils

import (
	"strings"
	"testing"
)

func NewBasicRepo(t *testing.T) *GitTestRepo {
	return NewGitTestRepo(t,
		WithDefaultBranch("main"),
		WithUser("Test User", "test@example.com"),
	)
}

func NewMultiBranchRepo(t *testing.T) *GitTestRepo {
	repo := NewBasicRepo(t)

	if err := repo.CreateBranch("develop", "Development content"); err != nil {
		t.Fatalf("Failed to create develop branch: %v", err)
	}

	if err := repo.CreateBranch("feature/auth", "Authentication feature"); err != nil {
		t.Fatalf("Failed to create feature/auth branch: %v", err)
	}

	if err := repo.CreateBranch("production/v1.0", "Production release"); err != nil {
		t.Fatalf("Failed to create production/v1.0 branch: %v", err)
	}

	return repo
}


// stringMappingToWorktreeConfigs converts the old string mapping format to the new explicit format
// This preserves the old behavior for backward compatibility in scenarios
func stringMappingToWorktreeConfigs(mapping map[string]string) map[string]WorktreeConfig {
	worktrees := make(map[string]WorktreeConfig)

	// Standard order that tests expect - matches the original .envrc order
	orderedKeys := []string{"main", "preview", "staging", "dev", "feat", "prod", "hotfix"}

	// Find keys that exist in the mapping
	var existingKeys []string
	for _, key := range orderedKeys {
		if _, exists := mapping[key]; exists {
			existingKeys = append(existingKeys, key)
		}
	}

	// Add any remaining keys alphabetically
	var remainingKeys []string
	for key := range mapping {
		found := false
		for _, standardKey := range orderedKeys {
			if key == standardKey {
				found = true
				break
			}
		}
		if !found {
			remainingKeys = append(remainingKeys, key)
		}
	}

	// Sort remaining keys for determinism
	for i := 0; i < len(remainingKeys); i++ {
		for j := i + 1; j < len(remainingKeys); j++ {
			if remainingKeys[i] > remainingKeys[j] {
				remainingKeys[i], remainingKeys[j] = remainingKeys[j], remainingKeys[i]
			}
		}
	}

	existingKeys = append(existingKeys, remainingKeys...)

	// Build hierarchy - each key merges to the previous one (except the first)
	for i, key := range existingKeys {
		config := WorktreeConfig{
			Branch:      mapping[key],
			Description: strings.ToUpper(key[:1]) + key[1:] + " branch",
		}

		// Add merge_into relationship for all except the first (root)
		if i > 0 {
			config.MergeInto = existingKeys[i-1]
		}

		worktrees[key] = config
	}

	return worktrees
}

func NewGBMConfigRepo(t *testing.T, mapping map[string]string) *GitTestRepo {
	repo := NewMultiBranchRepo(t)

	// Convert string mapping to WorktreeConfig mapping
	worktrees := stringMappingToWorktreeConfigs(mapping)

	if err := repo.CreateGBMConfig(worktrees); err != nil {
		t.Fatalf("Failed to create gbm.branchconfig.yaml: %v", err)
	}

	if err := repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml configuration"); err != nil {
		t.Fatalf("Failed to commit gbm.branchconfig.yaml: %v", err)
	}

	if err := repo.PushBranch("main"); err != nil {
		t.Fatalf("Failed to push main branch: %v", err)
	}

	return repo
}


func NewStandardGBMConfigRepo(t *testing.T) *GitTestRepo {
	mapping := map[string]string{
		"main": "main",
		"dev":  "develop",
		"feat": "feature/auth",
		"prod": "production/v1.0",
	}

	return NewGBMConfigRepo(t, mapping)
}