package testutils

import (
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


func NewGBMConfigRepo(t *testing.T, mapping map[string]string) *GitTestRepo {
	repo := NewMultiBranchRepo(t)

	if err := repo.CreateGBMConfig(mapping); err != nil {
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