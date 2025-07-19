EMPTY .
EMPTY internal/testutils
PASS internal.TestParseTimestamp/valid_unix_timestamp_string (0.00s)
PASS internal.TestParseTimestamp/invalid_timestamp (0.00s)
PASS internal.TestParseTimestamp/empty_timestamp (0.00s)
PASS internal.TestParseTimestamp (0.00s)
PASS internal.TestExtractJiraTicket/standard_JIRA_ticket (0.00s)
PASS internal.TestExtractJiraTicket/JIRA_ticket_at_beginning (0.00s)
PASS internal.TestExtractJiraTicket/JIRA_ticket_in_middle (0.00s)
PASS internal.TestExtractJiraTicket/multiple_JIRA_tickets (0.00s)
PASS internal.TestExtractJiraTicket/no_JIRA_ticket (0.00s)
PASS internal.TestExtractJiraTicket/lowercase_should_not_match (0.00s)
PASS internal.TestExtractJiraTicket/numbers_only_should_not_match (0.00s)
PASS internal.TestExtractJiraTicket/single_letter_project_code (0.00s)
PASS internal.TestExtractJiraTicket/long_project_code (0.00s)
PASS internal.TestExtractJiraTicket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/hotfix_with_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/feature_with_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/bugfix_with_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/merge_with_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/hotfix_without_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/feature_without_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/branch_without_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/branch_with_underscores (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/empty_branch_name (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch/branch_with_only_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromBranch (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_JIRA_ticket (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_feat_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_fix_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_hotfix_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_merge_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_add_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_with_update_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/message_without_prefix (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/empty_message (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/single_word_prefix_only (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage/uppercase_message (0.00s)
PASS internal.TestExtractWorktreeNameFromMessage (0.00s)
PASS internal.TestExtractBranchFromRef/origin_ref_with_hotfix (0.00s)
PASS internal.TestExtractBranchFromRef/origin_ref_with_feature (0.00s)
PASS internal.TestExtractBranchFromRef/ref_without_origin (0.00s)
PASS internal.TestExtractBranchFromRef/simple_branch_name (0.00s)
PASS internal.TestExtractBranchFromRef/empty_ref (0.00s)
PASS internal.TestExtractBranchFromRef/complex_origin_ref (0.00s)
PASS internal.TestExtractBranchFromRef (0.00s)
PASS internal.TestRecentActivityStruct (0.00s)
PASS internal.TestGitMergePatternRegex/standard_merge_message (0.00s)
PASS internal.TestGitMergePatternRegex/hotfix_merge_message (0.00s)
PASS internal.TestGitMergePatternRegex/complex_branch_names (0.00s)
PASS internal.TestGitMergePatternRegex/non-merge_message (0.00s)
PASS internal.TestGitMergePatternRegex/different_merge_format (0.00s)
PASS internal.TestGitMergePatternRegex/empty_message (0.00s)
PASS internal.TestGitMergePatternRegex (0.00s)
PASS internal.TestMockRecentActivity (0.00s)
PASS internal.TestCopyFilesToWorktree_AdHocOnly (0.00s)
PASS internal.TestCopyFilesToWorktree_NoRules (0.00s)
PASS internal.TestCopyFilesToWorktree_SourceNotExists (0.00s)
PASS internal.TestAddWorktree_TrackedWorktreeNoFileCopy (0.00s)
PASS internal.TestAddWorktree_AdHocWorktreeFileCopy (0.00s)
PASS internal.TestMergeBackDetection_RealWorldDemo (0.83s)
PASS cmd.TestAddCommand_NewBranchFromRemote (0.98s)
=== RUN   TestMergeBackDetection_BasicThreeTierScenario
Git command output: 
Git command output: 
Git command output: 
Git command output: [fb69657a-a7f8-4fe6-b2e6-57110f80b73a (root-commit) 227939b] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      fb69657a-a7f8-4fe6-b2e6-57110f80b73a -> fb69657a-a7f8-4fe6-b2e6-57110f80b73a

Git command output: 
Git command output: [fb69657a-a7f8-4fe6-b2e6-57110f80b73a 5dbca0d] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch '3923342c-db87-4954-9379-985924f91c5c'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      3923342c-db87-4954-9379-985924f91c5c -> 3923342c-db87-4954-9379-985924f91c5c

Git command output: Switched to branch 'fb69657a-a7f8-4fe6-b2e6-57110f80b73a'
Your branch is ahead of 'origin/fb69657a-a7f8-4fe6-b2e6-57110f80b73a' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'bc218870-56cc-48d8-9078-7be8b1de53dd'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      bc218870-56cc-48d8-9078-7be8b1de53dd -> bc218870-56cc-48d8-9078-7be8b1de53dd

Git command output: 
Git command output: [bc218870-56cc-48d8-9078-7be8b1de53dd a77cc85] Fix critical security vulnerability
 1 file changed, 1 insertion(+)
 create mode 100644 hotfix.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
   5dbca0d..a77cc85  bc218870-56cc-48d8-9078-7be8b1de53dd -> bc218870-56cc-48d8-9078-7be8b1de53dd

Git command output: Switched to branch 'fb69657a-a7f8-4fe6-b2e6-57110f80b73a'
Your branch is ahead of 'origin/fb69657a-a7f8-4fe6-b2e6-57110f80b73a' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:69: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:69
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:62
        	Error:      	"[{prod preview [{a77cc8576e4359cc1755a86098d00218bc5e7fca Fix critical security vulnerability Test User test@example.com 2025-07-17 21:23:34 -0400 EDT false}] [{a77cc8576e4359cc1755a86098d00218bc5e7fca Fix critical security vulnerability Test User test@example.com 2025-07-17 21:23:34 -0400 EDT true}] 1 1} {preview main [{5dbca0d68c0f71c34597a5ba879ef28d798b8f80 Add gbm.branchconfig.yaml configuration Test User test@example.com 2025-07-17 21:23:34 -0400 EDT false}] [{5dbca0d68c0f71c34597a5ba879ef28d798b8f80 Add gbm.branchconfig.yaml configuration Test User test@example.com 2025-07-17 21:23:34 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_BasicThreeTierScenario
--- FAIL: TestMergeBackDetection_BasicThreeTierScenario (0.85s)
FAIL internal.TestMergeBackDetection_BasicThreeTierScenario (0.85s)
PASS cmd.TestAddCommand_NewBranch (1.05s)
=== RUN   TestMergeBackDetection_MultipleCommits
Git command output: 
Git command output: 
Git command output: 
Git command output: [b780c001-befa-483c-a053-73d75947ff53 (root-commit) bdf5035] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      b780c001-befa-483c-a053-73d75947ff53 -> b780c001-befa-483c-a053-73d75947ff53

Git command output: 
Git command output: [b780c001-befa-483c-a053-73d75947ff53 50ab741] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch '1e9cee86-3043-4d2c-b0b8-1ab6825138a8'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      1e9cee86-3043-4d2c-b0b8-1ab6825138a8 -> 1e9cee86-3043-4d2c-b0b8-1ab6825138a8

Git command output: Switched to branch 'b780c001-befa-483c-a053-73d75947ff53'
Your branch is ahead of 'origin/b780c001-befa-483c-a053-73d75947ff53' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'fa6102c7-93a2-4c3e-8633-9152ccb5c729'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      fa6102c7-93a2-4c3e-8633-9152ccb5c729 -> fa6102c7-93a2-4c3e-8633-9152ccb5c729

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 d717750] Fix database connection issue
 1 file changed, 1 insertion(+)
 create mode 100644 fix1.txt

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 3fb9173] Fix memory leak in auth module
 1 file changed, 1 insertion(+)
 create mode 100644 fix2.txt

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 74f8731] Fix race condition in cache
 1 file changed, 1 insertion(+)
 create mode 100644 fix3.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
   50ab741..74f8731  fa6102c7-93a2-4c3e-8633-9152ccb5c729 -> fa6102c7-93a2-4c3e-8633-9152ccb5c729

Git command output: Switched to branch 'b780c001-befa-483c-a053-73d75947ff53'
Your branch is ahead of 'origin/b780c001-befa-483c-a053-73d75947ff53' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:155: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:155
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:148
        	Error:      	"[{prod preview [{74f8731415d45153f1f99e05289b7cf53b2da79e Fix race condition in cache Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false} {3fb9173473bd85ecbf2b291fac05a251d8c41b34 Fix memory leak in auth module Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false} {d7177504f86140ded6d439a296e640f10a9dcaa6 Fix database connection issue Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false}] [{74f8731415d45153f1f99e05289b7cf53b2da79e Fix race condition in cache Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true} {3fb9173473bd85ecbf2b291fac05a251d8c41b34 Fix memory leak in auth module Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true} {d7177504f86140ded6d439a296e640f10a9dcaa6 Fix database connection issue Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true}] 3 3} {preview main [{50ab741cc17baa2746ee3784688227775f6de9ce Add gbm.branchconfig.yaml configuration Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false}] [{50ab741cc17baa2746ee3784688227775f6de9ce Add gbm.branchconfig.yaml configuration Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_MultipleCommits
--- FAIL: TestMergeBackDetection_MultipleCommits (0.97s)
FAIL internal.TestMergeBackDetection_MultipleCommits (0.97s)
PASS cmd.TestAddCommand_NewBranchWithBaseBranch (1.04s)
PASS internal.TestMergeBackDetection_CascadingMergebacks (1.07s)
PASS cmd.TestAddCommand_InvalidBaseBranch (1.12s)
=== RUN   TestMergeBackDetection_NoMergeBacksNeeded
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 37f008a] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main b1fb437] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch 'preview'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      preview -> preview

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'prod'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      prod -> prod

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:315: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:315
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:308
        	Error:      	"[{preview main [{b1fb4370e631b04a9c60f1512a00e120707ee275 Add gbm.branchconfig.yaml configuration Developer dev@example.com 2025-07-17 21:23:37 -0400 EDT false}] [{b1fb4370e631b04a9c60f1512a00e120707ee275 Add gbm.branchconfig.yaml configuration Developer dev@example.com 2025-07-17 21:23:37 -0400 EDT true}] 1 1}]" should have 0 item(s), but has 1
        	Test:       	TestMergeBackDetection_NoMergeBacksNeeded
    mergeback_integration_test.go:316: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:316
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:308
        	Error:      	Should be false
        	Test:       	TestMergeBackDetection_NoMergeBacksNeeded
--- FAIL: TestMergeBackDetection_NoMergeBacksNeeded (0.81s)
FAIL internal.TestMergeBackDetection_NoMergeBacksNeeded (0.81s)
PASS internal.TestMergeBackDetection_NonExistentBranches (0.43s)
PASS cmd.TestAddCommand_JIRAKeyGeneration (1.44s)
PASS cmd.TestAddCommand_GenerateBranchName/Regular_name (0.00s)
PASS cmd.TestAddCommand_GenerateBranchName/Name_with_spaces (0.00s)
PASS cmd.TestAddCommand_GenerateBranchName/Name_with_underscores (0.00s)
PASS cmd.TestAddCommand_GenerateBranchName/Already_has_prefix (0.00s)
PASS cmd.TestAddCommand_GenerateBranchName (0.00s)
=== RUN   TestMergeBackDetection_DynamicHierarchy
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 5b75469] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 0616d4b] Add gbm.branchconfig.yaml configuration
 1 file changed, 23 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch 'staging'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      staging -> staging

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'preview'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      preview -> preview

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'prod'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      prod -> prod

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'hotfix'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      hotfix -> hotfix

Git command output: Already on 'hotfix'

Git command output: 
Git command output: [hotfix de7c02d] Emergency security patch
 1 file changed, 1 insertion(+)
 create mode 100644 emergency.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
   0616d4b..de7c02d  hotfix -> hotfix

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:425: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:425
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:418
        	Error:      	"[{hotfix prod [{de7c02df83f5ebcb48eb75f7bd066102bf7ffeeb Emergency security patch DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT false}] [{de7c02df83f5ebcb48eb75f7bd066102bf7ffeeb Emergency security patch DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT true}] 1 1} {preview main [{0616d4b893f920b333bc1d2a086faea8024e2732 Add gbm.branchconfig.yaml configuration DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT false}] [{0616d4b893f920b333bc1d2a086faea8024e2732 Add gbm.branchconfig.yaml configuration DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_DynamicHierarchy
--- FAIL: TestMergeBackDetection_DynamicHierarchy (1.19s)
FAIL internal.TestMergeBackDetection_DynamicHierarchy (1.19s)
PASS cmd.TestAddCommand_MissingBranchName (0.95s)
PASS internal.TestMergeBackAlertFormatting_RealScenario (0.81s)
PASS internal.TestCommitInfo/isUserCommit_with_email_match (0.00s)
PASS internal.TestCommitInfo/isUserCommit_with_name_match (0.00s)
PASS internal.TestCommitInfo/isUserCommit_with_no_match (0.00s)
PASS internal.TestCommitInfo (0.00s)
PASS internal.TestFormatMergeBackAlert/no_merge-backs_needed (0.00s)
PASS internal.TestFormatMergeBackAlert/nil_status (0.00s)
PASS internal.TestFormatMergeBackAlert/single_merge-back_with_user_commits (0.00s)
PASS internal.TestFormatMergeBackAlert/multiple_merge-backs_with_mixed_user_commits (0.00s)
PASS internal.TestFormatMergeBackAlert (0.00s)
PASS internal.TestFormatRelativeTime/just_now (0.00s)
PASS internal.TestFormatRelativeTime/30_minutes_ago (0.00s)
PASS internal.TestFormatRelativeTime/2_hours_ago (0.00s)
PASS internal.TestFormatRelativeTime/1_day_ago (0.00s)
PASS internal.TestFormatRelativeTime/3_days_ago (0.00s)
PASS internal.TestFormatRelativeTime (0.00s)
PASS internal.TestCheckMergeBackStatusIntegration/missing_gbm.branchconfig.yaml_file (0.03s)
=== RUN   TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file
    mergeback_test.go:206: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_test.go:206
        	Error:      	Expected value not to be nil.
        	Test:       	TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file
--- FAIL: TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file (0.03s)
FAIL internal.TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file (0.03s)
=== RUN   TestCheckMergeBackStatusIntegration
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ee6b7ff] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCheckMergeBackStatusIntegration1663872725/001/remote.git
 * [new branch]      main -> main

--- FAIL: TestCheckMergeBackStatusIntegration (0.28s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x1010d3558]

goroutine 512 [running]:
testing.tRunner.func1.2({0x101203720, 0x1014dbc70})
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1734 +0x1ac
testing.tRunner.func1()
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1737 +0x334
panic({0x101203720?, 0x1014dbc70?})
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/runtime/panic.go:792 +0x124
gbm/internal.TestCheckMergeBackStatusIntegration.func2(0x140001a9c00)
	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_test.go:207 +0xf8
testing.tRunner(0x140001a9c00, 0x101278688)
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1792 +0xe4
created by testing.(*T).Run in goroutine 506
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1851 +0x374
FAIL internal.TestCheckMergeBackStatusIntegration (0.28s)
FAIL internal
PASS cmd.TestAddCommand_NewBranchWithoutFlag (0.96s)
PASS cmd.TestAddCommand_AutoGenerateBranchWithFlag (0.95s)
PASS cmd.TestAddCommand_DuplicateWorktreeName (1.08s)
PASS cmd.TestAddCommand_ValidArgsFunction/First_argument_-_JIRA_keys (0.73s)
PASS cmd.TestAddCommand_ValidArgsFunction/Second_argument_-_branch_name (0.30s)
PASS cmd.TestAddCommand_ValidArgsFunction/Third_argument (0.00s)
PASS cmd.TestAddCommand_ValidArgsFunction (1.03s)
PASS cmd.TestAddCommand_FromWorktreeDirectory (1.11s)
=== RUN   TestCloneCommand_Basic
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 4a4d485] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 1494e64] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth 2a5692e] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 7d2c548] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at 4a4d485 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:63: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:63
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140002ed8c0)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_Basic
--- FAIL: TestCloneCommand_Basic (0.91s)
FAIL cmd.TestCloneCommand_Basic (0.91s)
=== RUN   TestCloneCommand_WithExistingGBMConfig
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 58e6749] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 6410c89] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth c8dfb85] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 555f68b] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: 
Git command output: [main e92a908] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
   58e6749..e92a908  main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at e92a908 Add gbm.branchconfig.yaml configuration
💡 Setting up gbm.branchconfig.yaml configuration...
💡 Found gbm.branchconfig.yaml in main worktree, copying to root...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:114: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:114
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"dev":internal.WorktreeConfig{Branch:"develop", MergeInto:"main", Description:"Dev branch"}, "feat":internal.WorktreeConfig{Branch:"feature/auth", MergeInto:"dev", Description:"Feat branch"}, "main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"dev":internal.WorktreeConfig{Branch:"develop", MergeInto:"main", Description:"Dev branch"}, "feat":internal.WorktreeConfig{Branch:"feature/auth", MergeInto:"dev", Description:"Feat branch"}, "main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main branch"}}, Tree:(*internal.WorktreeManager)(0x140004de120)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -18,3 +18,140 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=3) {
        	            	+   (string) (len=3) "dev": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=3) "dev",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) (len=4) "main",
        	            	+     Description: (string) (len=10) "Dev branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)({
        	            	+     Name: (string) (len=4) "main",
        	            	+     Config: (internal.WorktreeConfig) {
        	            	+      Branch: (string) (len=4) "main",
        	            	+      MergeInto: (string) "",
        	            	+      Description: (string) (len=11) "Main branch"
        	            	+     },
        	            	+     Parent: (*internal.WorktreeNode)(<nil>),
        	            	+     Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+      (*internal.WorktreeNode)(<already shown>)
        	            	+     }
        	            	+    }),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=4) "feat",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=12) "feature/auth",
        	            	+       MergeInto: (string) (len=3) "dev",
        	            	+       Description: (string) (len=11) "Feat branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) {
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   }),
        	            	+   (string) (len=4) "feat": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "feat",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=12) "feature/auth",
        	            	+     MergeInto: (string) (len=3) "dev",
        	            	+     Description: (string) (len=11) "Feat branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)({
        	            	+     Name: (string) (len=3) "dev",
        	            	+     Config: (internal.WorktreeConfig) {
        	            	+      Branch: (string) (len=7) "develop",
        	            	+      MergeInto: (string) (len=4) "main",
        	            	+      Description: (string) (len=10) "Dev branch"
        	            	+     },
        	            	+     Parent: (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=4) "main",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=4) "main",
        	            	+       MergeInto: (string) "",
        	            	+       Description: (string) (len=11) "Main branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<nil>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)(<already shown>)
        	            	+      }
        	            	+     }),
        	            	+     Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+      (*internal.WorktreeNode)(<already shown>)
        	            	+     }
        	            	+    }),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   }),
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=11) "Main branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=3) "dev",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=7) "develop",
        	            	+       MergeInto: (string) (len=4) "main",
        	            	+       Description: (string) (len=10) "Dev branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)({
        	            	+        Name: (string) (len=4) "feat",
        	            	+        Config: (internal.WorktreeConfig) {
        	            	+         Branch: (string) (len=12) "feature/auth",
        	            	+         MergeInto: (string) (len=3) "dev",
        	            	+         Description: (string) (len=11) "Feat branch"
        	            	+        },
        	            	+        Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+        Children: ([]*internal.WorktreeNode) {
        	            	+        }
        	            	+       })
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=11) "Main branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=3) "dev",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=7) "develop",
        	            	+       MergeInto: (string) (len=4) "main",
        	            	+       Description: (string) (len=10) "Dev branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)({
        	            	+        Name: (string) (len=4) "feat",
        	            	+        Config: (internal.WorktreeConfig) {
        	            	+         Branch: (string) (len=12) "feature/auth",
        	            	+         MergeInto: (string) (len=3) "dev",
        	            	+         Description: (string) (len=11) "Feat branch"
        	            	+        },
        	            	+        Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+        Children: ([]*internal.WorktreeNode) {
        	            	+        }
        	            	+       })
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_WithExistingGBMConfig
--- FAIL: TestCloneCommand_WithExistingGBMConfig (1.13s)
FAIL cmd.TestCloneCommand_WithExistingGBMConfig (1.13s)
=== RUN   TestCloneCommand_WithoutGBMConfig
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ea83bb2] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithoutGBMConfig3181735622/001/remote.git
 * [new branch]      main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithoutGBMConfig3181735622/001/remote
 * [new branch]      main       -> origin/main
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at ea83bb2 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:151: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:151
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140000a2820)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_WithoutGBMConfig
--- FAIL: TestCloneCommand_WithoutGBMConfig (0.41s)
FAIL cmd.TestCloneCommand_WithoutGBMConfig (0.41s)
=== RUN   TestCloneCommand_DifferentDefaultBranches/master_branch
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) ea83bb2] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesmaster_branch1243095124/001/remote.git
 * [new branch]      master -> master

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesmaster_branch1243095124/001/remote
 * [new branch]      master     -> origin/master
💡 Discovering default branch...
💡 Default branch: master
💡 Creating main worktree...
Preparing worktree (checking out 'master')
HEAD is now at ea83bb2 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in master worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"master":internal.WorktreeConfig{Branch:"master", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"master":internal.WorktreeConfig{Branch:"master", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x1400052c700)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=6) "master": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=6) "master",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=6) "master",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=6) "master",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=6) "master",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/master_branch
--- FAIL: TestCloneCommand_DifferentDefaultBranches/master_branch (0.49s)
FAIL cmd.TestCloneCommand_DifferentDefaultBranches/master_branch (0.49s)
=== RUN   TestCloneCommand_DifferentDefaultBranches/develop_branch
Git command output: 
Git command output: 
Git command output: 
Git command output: [develop (root-commit) b131a02] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesdevelop_branch1705929057/001/remote.git
 * [new branch]      develop -> develop

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesdevelop_branch1705929057/001/remote
 * [new branch]      develop    -> origin/develop
💡 Discovering default branch...
💡 Default branch: develop
💡 Creating main worktree...
Preparing worktree (checking out 'develop')
HEAD is now at b131a02 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in develop worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"develop":internal.WorktreeConfig{Branch:"develop", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"develop":internal.WorktreeConfig{Branch:"develop", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x14000228900)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=7) "develop": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=7) "develop",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=7) "develop",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/develop_branch
--- FAIL: TestCloneCommand_DifferentDefaultBranches/develop_branch (0.43s)
FAIL cmd.TestCloneCommand_DifferentDefaultBranches/develop_branch (0.43s)
=== RUN   TestCloneCommand_DifferentDefaultBranches/custom_branch
Git command output: 
Git command output: 
Git command output: 
Git command output: [custom-main (root-commit) b131a02] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchescustom_branch2505295580/001/remote.git
 * [new branch]      custom-main -> custom-main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchescustom_branch2505295580/001/remote
 * [new branch]      custom-main -> origin/custom-main
💡 Discovering default branch...
💡 Default branch: custom-main
💡 Creating main worktree...
Preparing worktree (checking out 'custom-main')
HEAD is now at b131a02 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in custom-main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"custom-main":internal.WorktreeConfig{Branch:"custom-main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"custom-main":internal.WorktreeConfig{Branch:"custom-main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x14000228020)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=11) "custom-main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=11) "custom-main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=11) "custom-main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=11) "custom-main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=11) "custom-main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/custom_branch
--- FAIL: TestCloneCommand_DifferentDefaultBranches/custom_branch (0.44s)
FAIL cmd.TestCloneCommand_DifferentDefaultBranches/custom_branch (0.44s)
=== RUN   TestCloneCommand_DifferentDefaultBranches
--- FAIL: TestCloneCommand_DifferentDefaultBranches (1.36s)
FAIL cmd.TestCloneCommand_DifferentDefaultBranches (1.36s)
PASS cmd.TestCloneCommand_DirectoryStructure (0.90s)
PASS cmd.TestCloneCommand_InvalidRepository (0.02s)
PASS cmd.TestExtractRepoName/GitHub_HTTPS (0.00s)
PASS cmd.TestExtractRepoName/GitHub_SSH (0.00s)
PASS cmd.TestExtractRepoName/Without_.git (0.00s)
PASS cmd.TestExtractRepoName/Local_path (0.00s)
PASS cmd.TestExtractRepoName/Local_path_with_.git (0.00s)
PASS cmd.TestExtractRepoName/Empty_string (0.00s)
PASS cmd.TestExtractRepoName/Single_path (0.00s)
PASS cmd.TestExtractRepoName (0.00s)
=== RUN   TestCreateDefaultGBMConfig
    clone_test.go:300: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:300
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140002286a0)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCreateDefaultGBMConfig
--- FAIL: TestCreateDefaultGBMConfig (0.00s)
FAIL cmd.TestCreateDefaultGBMConfig (0.00s)
PASS cmd.TestBuildDeploymentChain/production_branch_shows_full_chain (0.00s)
PASS cmd.TestBuildDeploymentChain/preview_branch_shows_partial_chain (0.00s)
PASS cmd.TestBuildDeploymentChain/main_branch_shows_only_itself (0.00s)
PASS cmd.TestBuildDeploymentChain/unknown_branch_shows_only_itself (0.00s)
PASS cmd.TestBuildDeploymentChain (0.00s)
PASS cmd.TestFindMergeIntoTarget/production_merges_into_preview (0.00s)
PASS cmd.TestFindMergeIntoTarget/preview_merges_into_main (0.00s)
PASS cmd.TestFindMergeIntoTarget/main_has_no_merge_target (0.00s)
PASS cmd.TestFindMergeIntoTarget/unknown_branch_has_no_merge_target (0.00s)
PASS cmd.TestFindMergeIntoTarget (0.00s)
PASS cmd.TestGenerateHotfixBranchName/simple_worktree_name (0.00s)
PASS cmd.TestGenerateHotfixBranchName/worktree_name_with_spaces (0.00s)
PASS cmd.TestGenerateHotfixBranchName/worktree_name_with_underscores (0.00s)
PASS cmd.TestGenerateHotfixBranchName/JIRA_ticket_as_worktree_name (0.00s)
PASS cmd.TestGenerateHotfixBranchName/mixed_case_worktree_name (0.00s)
PASS cmd.TestGenerateHotfixBranchName (0.00s)
PASS cmd.TestHotfixWorktreeNaming/default_prefix (0.00s)
PASS cmd.TestHotfixWorktreeNaming/custom_prefix (0.00s)
PASS cmd.TestHotfixWorktreeNaming/empty_prefix (0.00s)
PASS cmd.TestHotfixWorktreeNaming/single_char_prefix (0.00s)
PASS cmd.TestHotfixWorktreeNaming (0.00s)
PASS cmd.TestIsProductionBranchName/production_branch (0.00s)
PASS cmd.TestIsProductionBranchName/prod_branch (0.00s)
PASS cmd.TestIsProductionBranchName/main_branch (0.00s)
PASS cmd.TestIsProductionBranchName/master_branch (0.00s)
PASS cmd.TestIsProductionBranchName/release_branch (0.00s)
PASS cmd.TestIsProductionBranchName/production_with_version (0.00s)
PASS cmd.TestIsProductionBranchName/feature_branch (0.00s)
PASS cmd.TestIsProductionBranchName/development_branch (0.00s)
PASS cmd.TestIsProductionBranchName/staging_branch (0.00s)
PASS cmd.TestIsProductionBranchName (0.00s)
=== RUN   TestListCommand_EmptyRepository
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ec9498c] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_EmptyRepository2573104520/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 215210a] Add empty gbm.branchconfig.yaml
 1 file changed, 1 insertion(+)
 create mode 100644 gbm.branchconfig.yaml

Error: failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
    list_test.go:85: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/list_test.go:85
        	Error:      	Received unexpected error:
        	            	failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
        	Test:       	TestListCommand_EmptyRepository
--- FAIL: TestListCommand_EmptyRepository (0.32s)
FAIL cmd.TestListCommand_EmptyRepository (0.32s)
PASS cmd.TestListCommand_WithGBMConfigWorktrees (1.27s)
PASS cmd.TestListCommand_UntrackedWorktrees (1.13s)
PASS cmd.TestListCommand_OrphanedWorktrees (1.70s)
PASS cmd.TestListCommand_GitStatus (1.05s)
PASS cmd.TestListCommand_ExpectedBranchDisplay (1.62s)
=== RUN   TestListCommand_SortedOutput
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ba964d4] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 3593f35] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth 896ae67] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 9b9980b] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: 
Git command output: [main 5a1a165] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
   ba964d4..5a1a165  main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at 5a1a165 Add gbm.branchconfig.yaml configuration
💡 Setting up gbm.branchconfig.yaml configuration...
💡 Found gbm.branchconfig.yaml in main worktree, copying to root...
💡 Initializing worktree management...
💡 Repository cloned successfully!
⚠️  Merge-back required in tracked branches:

feat → dev: 1 commits need merge-back (0 by you)

dev → main: 1 commits need merge-back (0 by you)


💡 ✅ Successfully synchronized worktrees
💡 Adding worktree 'adhoc' on branch 'production/v1.0'
💡 Using default base branch: main
Error: failed to add worktree: branch 'production/v1.0' exists but is not based on 'main'. Please delete the branch and try again, or use a different branch name
Usage:
  gbm add <worktree-name> [branch-name] [base-branch] [flags]

Flags:
  -h, --help          help for add
  -i, --interactive   Interactive mode to select branch
  -b, --new-branch    Create a new branch for the worktree

Global Flags:
      --debug                 enable debug logging to ./gbm.log
      --worktree-dir string   override worktree directory location

    list_test.go:317: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/list_test.go:317
        	Error:      	Received unexpected error:
        	            	failed to add worktree: branch 'production/v1.0' exists but is not based on 'main'. Please delete the branch and try again, or use a different branch name
        	Test:       	TestListCommand_SortedOutput
--- FAIL: TestListCommand_SortedOutput (1.87s)
FAIL cmd.TestListCommand_SortedOutput (1.87s)
PASS cmd.TestListCommand_NoGitRepository (0.02s)
PASS cmd.TestListCommand_NoGBMConfigFile (0.30s)
PASS cmd.TestGetSmartMergebackCompletions/function_exists_and_handles_no_manager_gracefully (0.86s)
PASS cmd.TestGetSmartMergebackCompletions (0.86s)
PASS cmd.TestGetJiraCompletions/function_exists_and_handles_no_JIRA_gracefully (0.71s)
PASS cmd.TestGetJiraCompletions (0.71s)
PASS cmd.TestMergebackValidArgsFunction/returns_completions_for_first_argument (1.49s)
PASS cmd.TestMergebackValidArgsFunction (1.49s)
PASS cmd.TestCompletionFormatting (0.00s)
PASS cmd.TestCompletionPrioritization (0.00s)
PASS cmd.TestCompletionFallback/falls_back_to_JIRA_when_no_activities (0.00s)
PASS cmd.TestCompletionFallback (0.00s)
PASS cmd.TestCompletionSeparator (0.00s)
PASS cmd.TestCompletionEdgeCases/activity_with_empty_worktree_name (0.00s)
PASS cmd.TestCompletionEdgeCases/activity_with_empty_branch_name (0.00s)
PASS cmd.TestCompletionEdgeCases/activity_with_long_branch_name (0.00s)
PASS cmd.TestCompletionEdgeCases (0.00s)
PASS cmd.TestCompletionIntegration/smart_completions_with_activity (0.84s)
PASS cmd.TestCompletionIntegration/completion_formatting (0.79s)
PASS cmd.TestCompletionIntegration (2.35s)
PASS cmd.TestMergebackCommand/no_arguments_-_should_attempt_auto-detection (0.76s)
PASS cmd.TestMergebackCommand/with_worktree_name (0.73s)
PASS cmd.TestMergebackCommand/with_worktree_name_and_jira_ticket (1.08s)
PASS cmd.TestMergebackCommand/too_many_arguments (0.56s)
PASS cmd.TestMergebackCommand (3.14s)
PASS cmd.TestGenerateMergebackBranchName/simple_worktree_name_with_target (0.00s)
PASS cmd.TestGenerateMergebackBranchName/JIRA_ticket_with_target (0.00s)
PASS cmd.TestGenerateMergebackBranchName/worktree_with_spaces_and_underscores (0.00s)
PASS cmd.TestGenerateMergebackBranchName/uppercase_target_worktree (0.00s)
PASS cmd.TestGenerateMergebackBranchName (0.00s)
PASS cmd.TestFilterAndValidateActivities/filters_out_feature_branches (0.00s)
PASS cmd.TestFilterAndValidateActivities (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/hotfix_branch_with_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/feature_branch_with_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/bugfix_branch_with_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/hotfix_branch_without_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/branch_without_prefix (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch/empty_branch_name (0.00s)
PASS cmd.TestExtractWorktreeNameFromBranch (0.00s)
PASS cmd.TestExtractJiraTicket/commit_message_with_JIRA_ticket (0.00s)
PASS cmd.TestExtractJiraTicket/commit_message_with_multiple_JIRA_tickets (0.00s)
PASS cmd.TestExtractJiraTicket/commit_message_without_JIRA_ticket (0.00s)
PASS cmd.TestExtractJiraTicket/JIRA_ticket_at_end_of_message (0.00s)
PASS cmd.TestExtractJiraTicket/lowercase_jira_pattern_should_not_match (0.00s)
PASS cmd.TestExtractJiraTicket (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/message_with_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/message_without_JIRA_ticket (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/message_with_feat_prefix (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/message_with_update_prefix (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/empty_message (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage/single_word_message (0.00s)
PASS cmd.TestExtractWorktreeNameFromMessage (0.00s)
PASS cmd.TestFindPotentialMergeTargets/production_branch_merges_into_preview (0.00s)
PASS cmd.TestFindPotentialMergeTargets/preview_branch_merges_into_main (0.00s)
PASS cmd.TestFindPotentialMergeTargets/main_branch_has_no_merge_target (0.00s)
PASS cmd.TestFindPotentialMergeTargets/hotfix_branch_with_no_explicit_config (0.00s)
PASS cmd.TestFindPotentialMergeTargets/unknown_branch (0.00s)
PASS cmd.TestFindPotentialMergeTargets (0.00s)
PASS cmd.TestFindMergeTargetBranchAndWorktree/handles_missing_config_gracefully (0.32s)
PASS cmd.TestFindMergeTargetBranchAndWorktree (0.32s)
PASS cmd.TestMergebackIntegration/manual_mergeback_creation (0.48s)
PASS cmd.TestMergebackIntegration/verify_branch_naming (0.01s)
PASS cmd.TestMergebackIntegration (1.08s)
PASS cmd.TestMergebackWorktreeNaming/with_mergeback_prefix (0.00s)
=== RUN   TestMergebackWorktreeNaming/without_mergeback_prefix
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 2471dd5] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackWorktreeNamingwithout_mergeback_prefix2184777373/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 409e35a] Add gbm.branchconfig.yaml
 1 file changed, 7 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Error: failed to determine merge target branch: no mergeback targets found
Usage:
  gbm mergeback [worktree-name] [jira-ticket] [flags]

Aliases:
  mergeback, mb

Flags:
  -h, --help   help for mergeback

Global Flags:
      --debug                 enable debug logging to ./gbm.log
      --worktree-dir string   override worktree directory location

    mergeback_test.go:544: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:544
        	Error:      	Received unexpected error:
        	            	failed to determine merge target branch: no mergeback targets found
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
    mergeback_test.go:547: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:547
        	Error:      	unable to find file "worktrees/fix-auth_main"
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
    mergeback_test.go:555: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:555
        	Error:      	"" does not contain "merge/SHOP-456_main"
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
        	Messages:   	Branch should include target suffix
--- FAIL: TestMergebackWorktreeNaming/without_mergeback_prefix (0.36s)
FAIL cmd.TestMergebackWorktreeNaming/without_mergeback_prefix (0.36s)
=== RUN   TestMergebackWorktreeNaming
--- FAIL: TestMergebackWorktreeNaming (0.36s)
FAIL cmd.TestMergebackWorktreeNaming (0.36s)
PASS cmd.TestFindMergeTargetWithTreeStructure/Production_to_Master_mergeback_needed (0.69s)
=== RUN   TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 639b387] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 27be594] Add content for production-2025-05-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      production-2025-05-1 -> production-2025-05-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to a new branch 'production-2025-07-1'

Git command output: 
Git command output: [production-2025-07-1 1f4ee7f] Add content for production-2025-07-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      production-2025-07-1 -> production-2025-07-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-2025-07-1'

Git command output: 
Git command output: [production-2025-07-1 7b465db] Add preview change
 2 files changed, 16 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 preview-change.txt

    mergeback_tree_test.go:132: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:132
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : "production-2025-07-1"
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+production-2025-07-1
        	Test:       	TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date
        	Messages:   	Should return correct target branch
    mergeback_tree_test.go:133: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:133
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : "preview"
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+preview
        	Test:       	TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date
        	Messages:   	Should return correct target worktree
--- FAIL: TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date (0.67s)
FAIL cmd.TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date (0.67s)
=== RUN   TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 639b387] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructureSimple_two-level_chain52591763/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-branch'

Git command output: 
Git command output: [production-branch 1426055] Add content for production-branch
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructureSimple_two-level_chain52591763/001/remote.git
 * [new branch]      production-branch -> production-branch

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-branch'

Git command output: 
Git command output: [production-branch 8988250] Add production change
 2 files changed, 12 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 prod-change.txt

    mergeback_tree_test.go:131: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:131
        	Error:      	Received unexpected error:
        	            	no mergeback targets found
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
    mergeback_tree_test.go:132: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:132
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : ""
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
        	Messages:   	Should return correct target branch
    mergeback_tree_test.go:133: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:133
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : ""
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
        	Messages:   	Should return correct target worktree
--- FAIL: TestFindMergeTargetWithTreeStructure/Simple_two-level_chain (0.47s)
FAIL cmd.TestFindMergeTargetWithTreeStructure/Simple_two-level_chain (0.47s)
=== RUN   TestFindMergeTargetWithTreeStructure
--- FAIL: TestFindMergeTargetWithTreeStructure (1.83s)
FAIL cmd.TestFindMergeTargetWithTreeStructure (1.83s)
PASS cmd.TestMergebackNamingWithTreeStructure (1.09s)
=== RUN   TestMergebackNamingProductionToMaster
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 1b76672] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackNamingProductionToMaster1650401132/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 7ee6d5b] Add content for production-2025-05-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackNamingProductionToMaster1650401132/001/remote.git
 * [new branch]      production-2025-05-1 -> production-2025-05-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 3e53007] Add production change for master
 2 files changed, 12 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 production-change.txt

    mergeback_tree_test.go:246: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:246
        	Error:      	Received unexpected error:
        	            	no mergeback targets found
        	Test:       	TestMergebackNamingProductionToMaster
--- FAIL: TestMergebackNamingProductionToMaster (0.51s)
FAIL cmd.TestMergebackNamingProductionToMaster (0.51s)
PASS cmd.TestPullCommand_CurrentWorktree (2.19s)
PASS cmd.TestPullCommand_NamedWorktree (2.56s)
PASS cmd.TestPullCommand_AllWorktrees (3.17s)
PASS cmd.TestPullCommand_NotInWorktree (2.23s)
PASS cmd.TestPullCommand_NonexistentWorktree (2.47s)
PASS cmd.TestPullCommand_NotInGitRepo (0.02s)
PASS cmd.TestPullCommand_FastForward (2.16s)
PASS cmd.TestPullCommand_UpToDate (1.98s)
PASS cmd.TestPullCommand_WithLocalChanges (2.10s)
PASS cmd.TestPushCommand_CurrentWorktree (2.40s)
PASS cmd.TestPushCommand_NamedWorktree (2.69s)
PASS cmd.TestPushCommand_AllWorktrees (3.73s)
PASS cmd.TestPushCommand_WithExistingUpstream (2.16s)
PASS cmd.TestPushCommand_WithoutUpstream (2.01s)
PASS cmd.TestPushCommand_NotInWorktree (2.18s)
PASS cmd.TestPushCommand_NonexistentWorktree (2.35s)
PASS cmd.TestPushCommand_NotInGitRepo (0.02s)
PASS cmd.TestPushCommand_WithLocalCommits (2.68s)
PASS cmd.TestPushCommand_UpToDate (2.10s)
PASS cmd.TestPushCommand_EmptyWorktreeList (0.62s)
PASS cmd.TestRemoveCommand_SuccessfulRemoval (2.43s)
PASS cmd.TestRemoveCommand_NonexistentWorktree (2.17s)
PASS cmd.TestRemoveCommand_NotInGitRepo (0.02s)
PASS cmd.TestRemoveCommand_UncommittedChangesWithoutForce (2.25s)
PASS cmd.TestRemoveCommand_ForceWithUncommittedChanges (2.16s)
PASS cmd.TestRemoveCommand_ForceBypassesConfirmation (2.20s)
PASS cmd.TestRemoveCommand_UserAcceptsConfirmation (2.26s)
PASS cmd.TestRemoveCommand_UserAcceptsConfirmationWithYes (2.20s)
PASS cmd.TestRemoveCommand_UserDeclinesConfirmation (2.13s)
PASS cmd.TestRemoveCommand_UserDeclinesWithEmptyInput (2.21s)
PASS cmd.TestRemoveCommand_RemovalFromWorktreeDirectory (1.93s)
PASS cmd.TestRemoveCommand_UpdatesWorktreeList (2.36s)
PASS cmd.TestRemoveCommand_CleanupFilesystem (2.24s)
PASS cmd.TestShouldShowMergeBackAlerts_DisabledByConfig (0.03s)
PASS cmd.TestShouldShowMergeBackAlerts_EnabledWithNoState (0.02s)
PASS cmd.TestShouldShowMergeBackAlerts_TimestampLogic (0.06s)
PASS cmd.TestShouldShowMergeBackAlerts_UserCommitInterval (0.03s)
PASS cmd.TestUpdateLastMergebackCheck (0.29s)
PASS cmd.TestSwitchCommand_BasicWorktreeSwitching/switch_to_main_worktree (0.30s)
PASS cmd.TestSwitchCommand_BasicWorktreeSwitching/switch_to_dev_worktree (0.33s)
PASS cmd.TestSwitchCommand_BasicWorktreeSwitching/switch_to_feat_worktree (0.31s)
PASS cmd.TestSwitchCommand_BasicWorktreeSwitching/switch_to_prod_worktree (0.31s)
PASS cmd.TestSwitchCommand_BasicWorktreeSwitching (3.05s)
PASS cmd.TestSwitchCommand_PrintPathFlag/print_path_for_main (0.32s)
PASS cmd.TestSwitchCommand_PrintPathFlag/print_path_for_dev (0.31s)
PASS cmd.TestSwitchCommand_PrintPathFlag (2.41s)
PASS cmd.TestSwitchCommand_FuzzyMatching/case_insensitive_match_-_dev (0.33s)
PASS cmd.TestSwitchCommand_FuzzyMatching/case_insensitive_match_-_main (0.30s)
PASS cmd.TestSwitchCommand_FuzzyMatching/substring_match_-_fea (0.42s)
PASS cmd.TestSwitchCommand_FuzzyMatching/prefix_match_preference_-_ma (0.41s)
PASS cmd.TestSwitchCommand_FuzzyMatching/nonexistent_worktree (0.43s)
PASS cmd.TestSwitchCommand_FuzzyMatching (3.81s)
PASS cmd.TestSwitchCommand_ListWorktrees (2.32s)
PASS cmd.TestSwitchCommand_PreviousWorktree (2.81s)
PASS cmd.TestSwitchCommand_NoPreviousWorktree (2.18s)
PASS cmd.TestSwitchCommand_ShellIntegration (2.23s)
PASS cmd.TestSwitchCommand_ErrorConditions/not_in_git_repository (0.02s)
PASS cmd.TestSwitchCommand_ErrorConditions/worktree_does_not_exist (2.28s)
PASS cmd.TestSwitchCommand_ErrorConditions (2.30s)
PASS cmd.TestSyncCommand_BasicOperations/sync_with_existing_gbm_config_creates_all_worktrees (1.86s)
PASS cmd.TestSyncCommand_BasicOperations/sync_with_minimal_gbm_config (0.64s)
PASS cmd.TestSyncCommand_BasicOperations/sync_with_already_synced_repo_is_idempotent (2.26s)
PASS cmd.TestSyncCommand_BasicOperations (4.76s)
PASS cmd.TestSyncCommand_Flags/dry-run_flag_shows_changes_without_applying (1.38s)
PASS cmd.TestSyncCommand_Flags/force_flag_removes_orphaned_worktrees_with_confirmation (1.34s)
PASS cmd.TestSyncCommand_Flags (2.72s)
PASS cmd.TestSyncCommand_SyncScenarios/branch_reference_changed (1.24s)
PASS cmd.TestSyncCommand_SyncScenarios/new_worktree_added (1.24s)
PASS cmd.TestSyncCommand_SyncScenarios/worktree_removed (1.23s)
PASS cmd.TestSyncCommand_SyncScenarios/no_changes_needed (1.24s)
PASS cmd.TestSyncCommand_SyncScenarios (4.95s)
PASS cmd.TestSyncCommand_UntrackedWorktrees/untracked_worktree_preserved_by_default (1.24s)
PASS cmd.TestSyncCommand_UntrackedWorktrees/untracked_worktree_removed_with_--force (1.29s)
PASS cmd.TestSyncCommand_UntrackedWorktrees/dry-run_shows_untracked_worktree_would_be_removed (1.18s)
PASS cmd.TestSyncCommand_UntrackedWorktrees/tracked_worktrees_updated,_untracked_preserved_without_force (1.34s)
PASS cmd.TestSyncCommand_UntrackedWorktrees (5.06s)
PASS cmd.TestSyncCommand_ErrorHandling/not_a_git_repository (0.02s)
PASS cmd.TestSyncCommand_ErrorHandling/missing_gbm_config_file (0.35s)
PASS cmd.TestSyncCommand_ErrorHandling/invalid_branch_reference (0.61s)
PASS cmd.TestSyncCommand_ErrorHandling (0.98s)
PASS cmd.TestSyncCommand_Integration/complete_sync_workflow (1.41s)
PASS cmd.TestSyncCommand_Integration/sync_after_manual_worktree_changes (1.95s)
PASS cmd.TestSyncCommand_Integration (3.35s)
PASS cmd.TestSyncCommand_ForceConfirmationDirectManagerTest (1.26s)
PASS cmd.TestSyncCommand_ForceConfirmation/user_confirms_deletion_with_'y' (1.30s)
PASS cmd.TestSyncCommand_ForceConfirmation/user_confirms_deletion_with_'yes' (1.61s)
PASS cmd.TestSyncCommand_ForceConfirmation/user_cancels_deletion_with_'n' (1.32s)
PASS cmd.TestSyncCommand_ForceConfirmation/user_cancels_deletion_with_empty_response (1.36s)
PASS cmd.TestSyncCommand_ForceConfirmation/user_cancels_deletion_with_'no' (1.30s)
PASS cmd.TestSyncCommand_ForceConfirmation (6.88s)
PASS cmd.TestSyncCommand_WorktreePromotion (2.30s)
PASS cmd.TestValidateCommand_AllBranchesValid (1.50s)
PASS cmd.TestValidateCommand_SomeBranchesInvalid (1.39s)
PASS cmd.TestValidateCommand_BranchExistence/local_branches_only (1.45s)
PASS cmd.TestValidateCommand_BranchExistence/remote_branches_only (1.08s)
PASS cmd.TestValidateCommand_BranchExistence/both_local_and_remote_branches (1.45s)
PASS cmd.TestValidateCommand_BranchExistence/non-existent_branches (0.44s)
PASS cmd.TestValidateCommand_BranchExistence/branches_with_special_characters/slashes (1.22s)
PASS cmd.TestValidateCommand_BranchExistence (5.63s)
PASS cmd.TestValidateCommand_MissingGBMConfig (0.32s)
PASS cmd.TestValidateCommand_InvalidGBMConfigSyntax (0.36s)
=== RUN   TestValidateCommand_EmptyGBMConfig
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 02d5d7f] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestValidateCommand_EmptyGBMConfig918491144/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 99f2746] Add empty gbm.branchconfig.yaml
 1 file changed, 0 insertions(+), 0 deletions(-)
 create mode 100644 gbm.branchconfig.yaml

    validate_test.go:299: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/validate_test.go:299
        	Error:      	Received unexpected error:
        	            	failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
        	Test:       	TestValidateCommand_EmptyGBMConfig
        	Messages:   	Validate command should succeed with empty gbm.branchconfig.yaml
--- FAIL: TestValidateCommand_EmptyGBMConfig (0.39s)
FAIL cmd.TestValidateCommand_EmptyGBMConfig (0.39s)
PASS cmd.TestValidateCommand_NotInGitRepository (0.02s)
PASS cmd.TestValidateCommand_CorruptGitRepository (0.32s)
PASS cmd.TestValidateCommand_DuplicateWorktrees (0.36s)
PASS cmd.TestValidateCommand_VeryLongBranchNames (0.79s)
FAIL cmd

=== Failed
=== FAIL: cmd TestCloneCommand_Basic (0.91s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 4a4d485] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 1494e64] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth 2a5692e] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 7d2c548] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_Basic1638806173/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at 4a4d485 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:63: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:63
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140002ed8c0)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_Basic

=== FAIL: cmd TestCloneCommand_WithExistingGBMConfig (1.13s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 58e6749] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 6410c89] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth c8dfb85] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 555f68b] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: 
Git command output: [main e92a908] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote.git
   58e6749..e92a908  main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithExistingGBMConfig2809337861/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at e92a908 Add gbm.branchconfig.yaml configuration
💡 Setting up gbm.branchconfig.yaml configuration...
💡 Found gbm.branchconfig.yaml in main worktree, copying to root...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:114: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:114
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"dev":internal.WorktreeConfig{Branch:"develop", MergeInto:"main", Description:"Dev branch"}, "feat":internal.WorktreeConfig{Branch:"feature/auth", MergeInto:"dev", Description:"Feat branch"}, "main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"dev":internal.WorktreeConfig{Branch:"develop", MergeInto:"main", Description:"Dev branch"}, "feat":internal.WorktreeConfig{Branch:"feature/auth", MergeInto:"dev", Description:"Feat branch"}, "main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main branch"}}, Tree:(*internal.WorktreeManager)(0x140004de120)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -18,3 +18,140 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=3) {
        	            	+   (string) (len=3) "dev": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=3) "dev",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) (len=4) "main",
        	            	+     Description: (string) (len=10) "Dev branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)({
        	            	+     Name: (string) (len=4) "main",
        	            	+     Config: (internal.WorktreeConfig) {
        	            	+      Branch: (string) (len=4) "main",
        	            	+      MergeInto: (string) "",
        	            	+      Description: (string) (len=11) "Main branch"
        	            	+     },
        	            	+     Parent: (*internal.WorktreeNode)(<nil>),
        	            	+     Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+      (*internal.WorktreeNode)(<already shown>)
        	            	+     }
        	            	+    }),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=4) "feat",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=12) "feature/auth",
        	            	+       MergeInto: (string) (len=3) "dev",
        	            	+       Description: (string) (len=11) "Feat branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) {
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   }),
        	            	+   (string) (len=4) "feat": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "feat",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=12) "feature/auth",
        	            	+     MergeInto: (string) (len=3) "dev",
        	            	+     Description: (string) (len=11) "Feat branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)({
        	            	+     Name: (string) (len=3) "dev",
        	            	+     Config: (internal.WorktreeConfig) {
        	            	+      Branch: (string) (len=7) "develop",
        	            	+      MergeInto: (string) (len=4) "main",
        	            	+      Description: (string) (len=10) "Dev branch"
        	            	+     },
        	            	+     Parent: (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=4) "main",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=4) "main",
        	            	+       MergeInto: (string) "",
        	            	+       Description: (string) (len=11) "Main branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<nil>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)(<already shown>)
        	            	+      }
        	            	+     }),
        	            	+     Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+      (*internal.WorktreeNode)(<already shown>)
        	            	+     }
        	            	+    }),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   }),
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=11) "Main branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=3) "dev",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=7) "develop",
        	            	+       MergeInto: (string) (len=4) "main",
        	            	+       Description: (string) (len=10) "Dev branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)({
        	            	+        Name: (string) (len=4) "feat",
        	            	+        Config: (internal.WorktreeConfig) {
        	            	+         Branch: (string) (len=12) "feature/auth",
        	            	+         MergeInto: (string) (len=3) "dev",
        	            	+         Description: (string) (len=11) "Feat branch"
        	            	+        },
        	            	+        Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+        Children: ([]*internal.WorktreeNode) {
        	            	+        }
        	            	+       })
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=11) "Main branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+     (*internal.WorktreeNode)({
        	            	+      Name: (string) (len=3) "dev",
        	            	+      Config: (internal.WorktreeConfig) {
        	            	+       Branch: (string) (len=7) "develop",
        	            	+       MergeInto: (string) (len=4) "main",
        	            	+       Description: (string) (len=10) "Dev branch"
        	            	+      },
        	            	+      Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+      Children: ([]*internal.WorktreeNode) (len=1) {
        	            	+       (*internal.WorktreeNode)({
        	            	+        Name: (string) (len=4) "feat",
        	            	+        Config: (internal.WorktreeConfig) {
        	            	+         Branch: (string) (len=12) "feature/auth",
        	            	+         MergeInto: (string) (len=3) "dev",
        	            	+         Description: (string) (len=11) "Feat branch"
        	            	+        },
        	            	+        Parent: (*internal.WorktreeNode)(<already shown>),
        	            	+        Children: ([]*internal.WorktreeNode) {
        	            	+        }
        	            	+       })
        	            	+      }
        	            	+     })
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_WithExistingGBMConfig

=== FAIL: cmd TestCloneCommand_WithoutGBMConfig (0.41s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ea83bb2] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithoutGBMConfig3181735622/001/remote.git
 * [new branch]      main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_WithoutGBMConfig3181735622/001/remote
 * [new branch]      main       -> origin/main
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at ea83bb2 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:151: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:151
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140000a2820)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_WithoutGBMConfig

=== FAIL: cmd TestCloneCommand_DifferentDefaultBranches/master_branch (0.49s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) ea83bb2] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesmaster_branch1243095124/001/remote.git
 * [new branch]      master -> master

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesmaster_branch1243095124/001/remote
 * [new branch]      master     -> origin/master
💡 Discovering default branch...
💡 Default branch: master
💡 Creating main worktree...
Preparing worktree (checking out 'master')
HEAD is now at ea83bb2 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in master worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"master":internal.WorktreeConfig{Branch:"master", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"master":internal.WorktreeConfig{Branch:"master", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x1400052c700)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=6) "master": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=6) "master",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=6) "master",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=6) "master",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=6) "master",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/master_branch

=== FAIL: cmd TestCloneCommand_DifferentDefaultBranches/develop_branch (0.43s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [develop (root-commit) b131a02] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesdevelop_branch1705929057/001/remote.git
 * [new branch]      develop -> develop

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchesdevelop_branch1705929057/001/remote
 * [new branch]      develop    -> origin/develop
💡 Discovering default branch...
💡 Default branch: develop
💡 Creating main worktree...
Preparing worktree (checking out 'develop')
HEAD is now at b131a02 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in develop worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"develop":internal.WorktreeConfig{Branch:"develop", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"develop":internal.WorktreeConfig{Branch:"develop", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x14000228900)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=7) "develop": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=7) "develop",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=7) "develop",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=7) "develop",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/develop_branch

=== FAIL: cmd TestCloneCommand_DifferentDefaultBranches/custom_branch (0.44s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [custom-main (root-commit) b131a02] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchescustom_branch2505295580/001/remote.git
 * [new branch]      custom-main -> custom-main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCloneCommand_DifferentDefaultBranchescustom_branch2505295580/001/remote
 * [new branch]      custom-main -> origin/custom-main
💡 Discovering default branch...
💡 Default branch: custom-main
💡 Creating main worktree...
Preparing worktree (checking out 'custom-main')
HEAD is now at b131a02 Initial commit
💡 Setting up gbm.branchconfig.yaml configuration...
💡 No gbm.branchconfig.yaml found in custom-main worktree, creating new one...
💡 Initializing worktree management...
💡 Repository cloned successfully!
    clone_test.go:197: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:197
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"custom-main":internal.WorktreeConfig{Branch:"custom-main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"custom-main":internal.WorktreeConfig{Branch:"custom-main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x14000228020)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=11) "custom-main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=11) "custom-main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=11) "custom-main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=11) "custom-main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=11) "custom-main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCloneCommand_DifferentDefaultBranches/custom_branch

=== FAIL: cmd TestCloneCommand_DifferentDefaultBranches (1.36s)

=== FAIL: cmd TestCreateDefaultGBMConfig (0.00s)
    clone_test.go:300: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/clone_test.go:300
        	Error:      	Not equal: 
        	            	expected: &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(nil)}
        	            	actual  : &internal.GBMConfig{Worktrees:map[string]internal.WorktreeConfig{"main":internal.WorktreeConfig{Branch:"main", MergeInto:"", Description:"Main production branch"}}, Tree:(*internal.WorktreeManager)(0x140002286a0)}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -8,3 +8,30 @@
        	            	  },
        	            	- Tree: (*internal.WorktreeManager)(<nil>)
        	            	+ Tree: (*internal.WorktreeManager)({
        	            	+  nodes: (map[string]*internal.WorktreeNode) (len=1) {
        	            	+   (string) (len=4) "main": (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  },
        	            	+  roots: ([]*internal.WorktreeNode) (len=1) {
        	            	+   (*internal.WorktreeNode)({
        	            	+    Name: (string) (len=4) "main",
        	            	+    Config: (internal.WorktreeConfig) {
        	            	+     Branch: (string) (len=4) "main",
        	            	+     MergeInto: (string) "",
        	            	+     Description: (string) (len=22) "Main production branch"
        	            	+    },
        	            	+    Parent: (*internal.WorktreeNode)(<nil>),
        	            	+    Children: ([]*internal.WorktreeNode) {
        	            	+    }
        	            	+   })
        	            	+  }
        	            	+ })
        	            	 })
        	Test:       	TestCreateDefaultGBMConfig

=== FAIL: cmd TestListCommand_EmptyRepository (0.32s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ec9498c] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_EmptyRepository2573104520/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 215210a] Add empty gbm.branchconfig.yaml
 1 file changed, 1 insertion(+)
 create mode 100644 gbm.branchconfig.yaml

Error: failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
    list_test.go:85: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/list_test.go:85
        	Error:      	Received unexpected error:
        	            	failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
        	Test:       	TestListCommand_EmptyRepository

=== FAIL: cmd TestListCommand_SortedOutput (1.87s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ba964d4] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      main -> main

Git command output: Switched to a new branch 'develop'

Git command output: 
Git command output: [develop 3593f35] Add content for develop
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      develop -> develop

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'feature/auth'

Git command output: 
Git command output: [feature/auth 896ae67] Add content for feature/auth
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      feature/auth -> feature/auth

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: Switched to a new branch 'production/v1.0'

Git command output: 
Git command output: [production/v1.0 9b9980b] Add content for production/v1.0
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
 * [new branch]      production/v1.0 -> production/v1.0

Git command output: Switched to branch 'main'
Your branch is up to date with 'origin/main'.

Git command output: 
Git command output: [main 5a1a165] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote.git
   ba964d4..5a1a165  main -> main

💡 Cloning repository using git-bare-clone.sh...
💡 Cloning bare repository to .git...
Cloning into bare repository '.git'...
done.
💡 Adjusting origin fetch locations...
💡 Fetching all branches from remote...
From /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestListCommand_SortedOutput833795518/001/remote
 * [new branch]      develop         -> origin/develop
 * [new branch]      feature/auth    -> origin/feature/auth
 * [new branch]      main            -> origin/main
 * [new branch]      production/v1.0 -> origin/production/v1.0
💡 Discovering default branch...
💡 Default branch: main
💡 Creating main worktree...
Preparing worktree (checking out 'main')
HEAD is now at 5a1a165 Add gbm.branchconfig.yaml configuration
💡 Setting up gbm.branchconfig.yaml configuration...
💡 Found gbm.branchconfig.yaml in main worktree, copying to root...
💡 Initializing worktree management...
💡 Repository cloned successfully!
⚠️  Merge-back required in tracked branches:

feat → dev: 1 commits need merge-back (0 by you)

dev → main: 1 commits need merge-back (0 by you)


💡 ✅ Successfully synchronized worktrees
💡 Adding worktree 'adhoc' on branch 'production/v1.0'
💡 Using default base branch: main
Error: failed to add worktree: branch 'production/v1.0' exists but is not based on 'main'. Please delete the branch and try again, or use a different branch name
Usage:
  gbm add <worktree-name> [branch-name] [base-branch] [flags]

Flags:
  -h, --help          help for add
  -i, --interactive   Interactive mode to select branch
  -b, --new-branch    Create a new branch for the worktree

Global Flags:
      --debug                 enable debug logging to ./gbm.log
      --worktree-dir string   override worktree directory location

    list_test.go:317: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/list_test.go:317
        	Error:      	Received unexpected error:
        	            	failed to add worktree: branch 'production/v1.0' exists but is not based on 'main'. Please delete the branch and try again, or use a different branch name
        	Test:       	TestListCommand_SortedOutput

=== FAIL: cmd TestMergebackWorktreeNaming/without_mergeback_prefix (0.36s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 2471dd5] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackWorktreeNamingwithout_mergeback_prefix2184777373/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 409e35a] Add gbm.branchconfig.yaml
 1 file changed, 7 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Error: failed to determine merge target branch: no mergeback targets found
Usage:
  gbm mergeback [worktree-name] [jira-ticket] [flags]

Aliases:
  mergeback, mb

Flags:
  -h, --help   help for mergeback

Global Flags:
      --debug                 enable debug logging to ./gbm.log
      --worktree-dir string   override worktree directory location

    mergeback_test.go:544: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:544
        	Error:      	Received unexpected error:
        	            	failed to determine merge target branch: no mergeback targets found
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
    mergeback_test.go:547: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:547
        	Error:      	unable to find file "worktrees/fix-auth_main"
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
    mergeback_test.go:555: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_test.go:555
        	Error:      	"" does not contain "merge/SHOP-456_main"
        	Test:       	TestMergebackWorktreeNaming/without_mergeback_prefix
        	Messages:   	Branch should include target suffix

=== FAIL: cmd TestMergebackWorktreeNaming (0.36s)

=== FAIL: cmd TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date (0.67s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 639b387] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 27be594] Add content for production-2025-05-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      production-2025-05-1 -> production-2025-05-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to a new branch 'production-2025-07-1'

Git command output: 
Git command output: [production-2025-07-1 1f4ee7f] Add content for production-2025-07-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructurePreview_to_Master_mergeback_needed_when_production_is_up_to_date2318494873/001/remote.git
 * [new branch]      production-2025-07-1 -> production-2025-07-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-2025-07-1'

Git command output: 
Git command output: [production-2025-07-1 7b465db] Add preview change
 2 files changed, 16 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 preview-change.txt

    mergeback_tree_test.go:132: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:132
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : "production-2025-07-1"
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+production-2025-07-1
        	Test:       	TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date
        	Messages:   	Should return correct target branch
    mergeback_tree_test.go:133: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:133
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : "preview"
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+preview
        	Test:       	TestFindMergeTargetWithTreeStructure/Preview_to_Master_mergeback_needed_when_production_is_up_to_date
        	Messages:   	Should return correct target worktree

=== FAIL: cmd TestFindMergeTargetWithTreeStructure/Simple_two-level_chain (0.47s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 639b387] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructureSimple_two-level_chain52591763/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-branch'

Git command output: 
Git command output: [production-branch 1426055] Add content for production-branch
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestFindMergeTargetWithTreeStructureSimple_two-level_chain52591763/001/remote.git
 * [new branch]      production-branch -> production-branch

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-branch'

Git command output: 
Git command output: [production-branch 8988250] Add production change
 2 files changed, 12 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 prod-change.txt

    mergeback_tree_test.go:131: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:131
        	Error:      	Received unexpected error:
        	            	no mergeback targets found
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
    mergeback_tree_test.go:132: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:132
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : ""
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
        	Messages:   	Should return correct target branch
    mergeback_tree_test.go:133: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:133
        	Error:      	Not equal: 
        	            	expected: "master"
        	            	actual  : ""
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1 +1 @@
        	            	-master
        	            	+
        	Test:       	TestFindMergeTargetWithTreeStructure/Simple_two-level_chain
        	Messages:   	Should return correct target worktree

=== FAIL: cmd TestFindMergeTargetWithTreeStructure (1.83s)

=== FAIL: cmd TestMergebackNamingProductionToMaster (0.51s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [master (root-commit) 1b76672] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackNamingProductionToMaster1650401132/001/remote.git
 * [new branch]      master -> master

Git command output: Switched to a new branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 7ee6d5b] Add content for production-2025-05-1
 1 file changed, 1 insertion(+)
 create mode 100644 content.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergebackNamingProductionToMaster1650401132/001/remote.git
 * [new branch]      production-2025-05-1 -> production-2025-05-1

Git command output: Switched to branch 'master'
Your branch is up to date with 'origin/master'.

Git command output: Switched to branch 'production-2025-05-1'

Git command output: 
Git command output: [production-2025-05-1 3e53007] Add production change for master
 2 files changed, 12 insertions(+)
 create mode 100644 gbm.branchconfig.yaml
 create mode 100644 production-change.txt

    mergeback_tree_test.go:246: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback_tree_test.go:246
        	Error:      	Received unexpected error:
        	            	no mergeback targets found
        	Test:       	TestMergebackNamingProductionToMaster

=== FAIL: cmd TestValidateCommand_EmptyGBMConfig (0.39s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 02d5d7f] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestValidateCommand_EmptyGBMConfig918491144/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 99f2746] Add empty gbm.branchconfig.yaml
 1 file changed, 0 insertions(+), 0 deletions(-)
 create mode 100644 gbm.branchconfig.yaml

    validate_test.go:299: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/cmd/validate_test.go:299
        	Error:      	Received unexpected error:
        	            	failed to load gbm.branchconfig.yaml: failed to build worktree tree: no root nodes found (all nodes have merge_into)
        	Test:       	TestValidateCommand_EmptyGBMConfig
        	Messages:   	Validate command should succeed with empty gbm.branchconfig.yaml

=== FAIL: internal TestMergeBackDetection_BasicThreeTierScenario (0.85s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [fb69657a-a7f8-4fe6-b2e6-57110f80b73a (root-commit) 227939b] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      fb69657a-a7f8-4fe6-b2e6-57110f80b73a -> fb69657a-a7f8-4fe6-b2e6-57110f80b73a

Git command output: 
Git command output: [fb69657a-a7f8-4fe6-b2e6-57110f80b73a 5dbca0d] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch '3923342c-db87-4954-9379-985924f91c5c'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      3923342c-db87-4954-9379-985924f91c5c -> 3923342c-db87-4954-9379-985924f91c5c

Git command output: Switched to branch 'fb69657a-a7f8-4fe6-b2e6-57110f80b73a'
Your branch is ahead of 'origin/fb69657a-a7f8-4fe6-b2e6-57110f80b73a' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'bc218870-56cc-48d8-9078-7be8b1de53dd'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
 * [new branch]      bc218870-56cc-48d8-9078-7be8b1de53dd -> bc218870-56cc-48d8-9078-7be8b1de53dd

Git command output: 
Git command output: [bc218870-56cc-48d8-9078-7be8b1de53dd a77cc85] Fix critical security vulnerability
 1 file changed, 1 insertion(+)
 create mode 100644 hotfix.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_BasicThreeTierScenario328134025/001/remote.git
   5dbca0d..a77cc85  bc218870-56cc-48d8-9078-7be8b1de53dd -> bc218870-56cc-48d8-9078-7be8b1de53dd

Git command output: Switched to branch 'fb69657a-a7f8-4fe6-b2e6-57110f80b73a'
Your branch is ahead of 'origin/fb69657a-a7f8-4fe6-b2e6-57110f80b73a' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:69: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:69
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:62
        	Error:      	"[{prod preview [{a77cc8576e4359cc1755a86098d00218bc5e7fca Fix critical security vulnerability Test User test@example.com 2025-07-17 21:23:34 -0400 EDT false}] [{a77cc8576e4359cc1755a86098d00218bc5e7fca Fix critical security vulnerability Test User test@example.com 2025-07-17 21:23:34 -0400 EDT true}] 1 1} {preview main [{5dbca0d68c0f71c34597a5ba879ef28d798b8f80 Add gbm.branchconfig.yaml configuration Test User test@example.com 2025-07-17 21:23:34 -0400 EDT false}] [{5dbca0d68c0f71c34597a5ba879ef28d798b8f80 Add gbm.branchconfig.yaml configuration Test User test@example.com 2025-07-17 21:23:34 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_BasicThreeTierScenario

=== FAIL: internal TestMergeBackDetection_MultipleCommits (0.97s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [b780c001-befa-483c-a053-73d75947ff53 (root-commit) bdf5035] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      b780c001-befa-483c-a053-73d75947ff53 -> b780c001-befa-483c-a053-73d75947ff53

Git command output: 
Git command output: [b780c001-befa-483c-a053-73d75947ff53 50ab741] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch '1e9cee86-3043-4d2c-b0b8-1ab6825138a8'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      1e9cee86-3043-4d2c-b0b8-1ab6825138a8 -> 1e9cee86-3043-4d2c-b0b8-1ab6825138a8

Git command output: Switched to branch 'b780c001-befa-483c-a053-73d75947ff53'
Your branch is ahead of 'origin/b780c001-befa-483c-a053-73d75947ff53' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'fa6102c7-93a2-4c3e-8633-9152ccb5c729'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
 * [new branch]      fa6102c7-93a2-4c3e-8633-9152ccb5c729 -> fa6102c7-93a2-4c3e-8633-9152ccb5c729

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 d717750] Fix database connection issue
 1 file changed, 1 insertion(+)
 create mode 100644 fix1.txt

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 3fb9173] Fix memory leak in auth module
 1 file changed, 1 insertion(+)
 create mode 100644 fix2.txt

Git command output: 
Git command output: [fa6102c7-93a2-4c3e-8633-9152ccb5c729 74f8731] Fix race condition in cache
 1 file changed, 1 insertion(+)
 create mode 100644 fix3.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_MultipleCommits3072156424/001/remote.git
   50ab741..74f8731  fa6102c7-93a2-4c3e-8633-9152ccb5c729 -> fa6102c7-93a2-4c3e-8633-9152ccb5c729

Git command output: Switched to branch 'b780c001-befa-483c-a053-73d75947ff53'
Your branch is ahead of 'origin/b780c001-befa-483c-a053-73d75947ff53' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:155: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:155
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:148
        	Error:      	"[{prod preview [{74f8731415d45153f1f99e05289b7cf53b2da79e Fix race condition in cache Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false} {3fb9173473bd85ecbf2b291fac05a251d8c41b34 Fix memory leak in auth module Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false} {d7177504f86140ded6d439a296e640f10a9dcaa6 Fix database connection issue Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false}] [{74f8731415d45153f1f99e05289b7cf53b2da79e Fix race condition in cache Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true} {3fb9173473bd85ecbf2b291fac05a251d8c41b34 Fix memory leak in auth module Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true} {d7177504f86140ded6d439a296e640f10a9dcaa6 Fix database connection issue Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true}] 3 3} {preview main [{50ab741cc17baa2746ee3784688227775f6de9ce Add gbm.branchconfig.yaml configuration Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT false}] [{50ab741cc17baa2746ee3784688227775f6de9ce Add gbm.branchconfig.yaml configuration Alice alice@example.com 2025-07-17 21:23:35 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_MultipleCommits

=== FAIL: internal TestMergeBackDetection_NoMergeBacksNeeded (0.81s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 37f008a] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main b1fb437] Add gbm.branchconfig.yaml configuration
 1 file changed, 15 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch 'preview'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      preview -> preview

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'prod'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_NoMergeBacksNeeded3729196665/001/remote.git
 * [new branch]      prod -> prod

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:315: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:315
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:308
        	Error:      	"[{preview main [{b1fb4370e631b04a9c60f1512a00e120707ee275 Add gbm.branchconfig.yaml configuration Developer dev@example.com 2025-07-17 21:23:37 -0400 EDT false}] [{b1fb4370e631b04a9c60f1512a00e120707ee275 Add gbm.branchconfig.yaml configuration Developer dev@example.com 2025-07-17 21:23:37 -0400 EDT true}] 1 1}]" should have 0 item(s), but has 1
        	Test:       	TestMergeBackDetection_NoMergeBacksNeeded
    mergeback_integration_test.go:316: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:316
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:308
        	Error:      	Should be false
        	Test:       	TestMergeBackDetection_NoMergeBacksNeeded

=== FAIL: internal TestMergeBackDetection_DynamicHierarchy (1.19s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) 5b75469] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      main -> main

Git command output: 
Git command output: [main 0616d4b] Add gbm.branchconfig.yaml configuration
 1 file changed, 23 insertions(+)
 create mode 100644 gbm.branchconfig.yaml

Git command output: Switched to a new branch 'staging'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      staging -> staging

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'preview'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      preview -> preview

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'prod'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      prod -> prod

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Git command output: Switched to a new branch 'hotfix'

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
 * [new branch]      hotfix -> hotfix

Git command output: Already on 'hotfix'

Git command output: 
Git command output: [hotfix de7c02d] Emergency security patch
 1 file changed, 1 insertion(+)
 create mode 100644 emergency.txt

Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestMergeBackDetection_DynamicHierarchy3850042465/001/remote.git
   0616d4b..de7c02d  hotfix -> hotfix

Git command output: Switched to branch 'main'
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

    mergeback_integration_test.go:425: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:425
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go:444
        	            				/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_integration_test.go:418
        	Error:      	"[{hotfix prod [{de7c02df83f5ebcb48eb75f7bd066102bf7ffeeb Emergency security patch DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT false}] [{de7c02df83f5ebcb48eb75f7bd066102bf7ffeeb Emergency security patch DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT true}] 1 1} {preview main [{0616d4b893f920b333bc1d2a086faea8024e2732 Add gbm.branchconfig.yaml configuration DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT false}] [{0616d4b893f920b333bc1d2a086faea8024e2732 Add gbm.branchconfig.yaml configuration DevOps devops@example.com 2025-07-17 21:23:38 -0400 EDT true}] 1 1}]" should have 1 item(s), but has 2
        	Test:       	TestMergeBackDetection_DynamicHierarchy

=== FAIL: internal TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file (0.03s)
    mergeback_test.go:206: 
        	Error Trace:	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_test.go:206
        	Error:      	Expected value not to be nil.
        	Test:       	TestCheckMergeBackStatusIntegration/empty_gbm.branchconfig.yaml_file

=== FAIL: internal TestCheckMergeBackStatusIntegration (0.28s)
Git command output: 
Git command output: 
Git command output: 
Git command output: [main (root-commit) ee6b7ff] Initial commit
 1 file changed, 1 insertion(+)
 create mode 100644 README.md

Git command output: 
Git command output: To /var/folders/hl/ppqjmzrx72q_x8gf0087xt7r0000gn/T/TestCheckMergeBackStatusIntegration1663872725/001/remote.git
 * [new branch]      main -> main

panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x0 pc=0x1010d3558]

goroutine 512 [running]:
testing.tRunner.func1.2({0x101203720, 0x1014dbc70})
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1734 +0x1ac
testing.tRunner.func1()
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1737 +0x334
panic({0x101203720?, 0x1014dbc70?})
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/runtime/panic.go:792 +0x124
gbm/internal.TestCheckMergeBackStatusIntegration.func2(0x140001a9c00)
	/Users/jschneider/code/scratch/worktree-manager/internal/mergeback_test.go:207 +0xf8
testing.tRunner(0x140001a9c00, 0x101278688)
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1792 +0xe4
created by testing.(*T).Run in goroutine 506
	/nix/store/rq7irijkj3nhapmjcv9d96xgkisj55x2-go-1.24.4/share/go/src/testing/testing.go:1851 +0x374

DONE 337 tests, 23 failures in 173.892s
