package testutils

import (
	"fmt"
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

func NewEnvrcRepo(t *testing.T, mapping map[string]string) *GitTestRepo {
	repo := NewMultiBranchRepo(t)

	if err := repo.CreateEnvrc(mapping); err != nil {
		t.Fatalf("Failed to create .envrc: %v", err)
	}

	if err := repo.CommitChangesWithForceAdd("Add .envrc configuration"); err != nil {
		t.Fatalf("Failed to commit .envrc: %v", err)
	}

	if err := repo.PushBranch("main"); err != nil {
		t.Fatalf("Failed to push main branch: %v", err)
	}

	return repo
}

func NewStandardEnvrcRepo(t *testing.T) *GitTestRepo {
	mapping := map[string]string{
		"MAIN": "main",
		"DEV":  "develop",
		"FEAT": "feature/auth",
		"PROD": "production/v1.0",
	}

	return NewEnvrcRepo(t, mapping)
}

func NewRepoWithConflictingBranches(t *testing.T) *GitTestRepo {
	repo := NewBasicRepo(t)

	if err := repo.CreateBranch("feature/conflict", "Original content"); err != nil {
		t.Fatalf("Failed to create conflict branch: %v", err)
	}

	if err := repo.SwitchToBranch("main"); err != nil {
		t.Fatalf("Failed to switch to main: %v", err)
	}

	if err := repo.WriteFile("conflict.txt", "Main branch content"); err != nil {
		t.Fatalf("Failed to write conflict file: %v", err)
	}

	if err := repo.CommitChanges("Add conflict file on main"); err != nil {
		t.Fatalf("Failed to commit on main: %v", err)
	}

	if err := repo.SwitchToBranch("feature/conflict"); err != nil {
		t.Fatalf("Failed to switch to conflict branch: %v", err)
	}

	if err := repo.WriteFile("conflict.txt", "Feature branch content"); err != nil {
		t.Fatalf("Failed to write conflict file: %v", err)
	}

	if err := repo.CommitChanges("Add conflict file on feature"); err != nil {
		t.Fatalf("Failed to commit on feature: %v", err)
	}

	if err := repo.SwitchToBranch("main"); err != nil {
		t.Fatalf("Failed to switch back to main: %v", err)
	}

	return repo
}

func NewLargeHistoryRepo(t *testing.T) *GitTestRepo {
	repo := NewBasicRepo(t)

	for i := 1; i <= 10; i++ {
		if err := repo.WriteFile("history.txt", fmt.Sprintf("Version %d content", i)); err != nil {
			t.Fatalf("Failed to write history file %d: %v", i, err)
		}

		if err := repo.CommitChanges(fmt.Sprintf("Version %d", i)); err != nil {
			t.Fatalf("Failed to commit version %d: %v", i, err)
		}
	}

	if err := repo.PushBranch("main"); err != nil {
		t.Fatalf("Failed to push main branch: %v", err)
	}

	return repo
}

func NewEmptyRepo(t *testing.T) *GitTestRepo {
	repo := NewGitTestRepo(t,
		WithDefaultBranch("main"),
		WithUser("Test User", "test@example.com"),
	)

	return repo
}