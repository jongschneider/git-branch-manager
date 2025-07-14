package internal

import (
	"path/filepath"
	"strings"
	"testing"

	"gbm/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestMergeBackDetection_RealWorldDemo(t *testing.T) {
	// This test demonstrates the real-world usage of merge-back detection
	// and alert formatting with an actual git repository

	repo := testutils.NewGitTestRepo(t,
		testutils.WithDefaultBranch("main"),
		testutils.WithUser("Jane Developer", "jane@company.com"),
	)

	// Set up a typical production environment hierarchy
	err := repo.CreateGBMConfig(map[string]string{
		"main":    "main",
		"staging": "staging",
		"prod":    "production",
	})
	require.NoError(t, err)

	err = repo.CommitChangesWithForceAdd("Add gbm.branchconfig.yaml configuration")
	require.NoError(t, err)

	// Create synchronized production branches
	err = repo.CreateSynchronizedBranch("staging")
	require.NoError(t, err)

	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	err = repo.CreateSynchronizedBranch("production")
	require.NoError(t, err)

	// Simulate a production hotfix scenario
	err = repo.WriteFile("security_patch.py", `# Critical security patch
def validate_user_input(input_data):
    # Fix SQL injection vulnerability
    return sanitize_sql(input_data)`)
	require.NoError(t, err)

	err = repo.CommitChanges("CVE-2024-1234: Fix SQL injection in user authentication")
	require.NoError(t, err)

	err = repo.PushBranch("production")
	require.NoError(t, err)

	// Switch back to main and test the detection
	err = repo.SwitchToBranch("main")
	require.NoError(t, err)

	// Demonstrate the merge-back detection in action
	err = repo.InLocalRepo(func() error {
		configPath := filepath.Join(repo.GetLocalPath(), DefaultBranchConfigFilename)
		status, err := CheckMergeBackStatus(configPath)
		require.NoError(t, err)
		require.NotNil(t, status)

		// Generate the alert that would be shown to the user
		alert := FormatMergeBackAlert(status)

		// Print the actual alert that would be displayed
		separator := strings.Repeat("=", 60)
		t.Logf("\n%s", separator)
		t.Logf("MERGE-BACK ALERT DEMO")
		t.Logf("%s", separator)
		if alert != "" {
			t.Logf("%s", alert)
		} else {
			t.Logf("✅ No merge-backs required - all branches are synchronized!")
		}
		t.Logf("%s", separator)

		return nil
	})
	require.NoError(t, err)
}
