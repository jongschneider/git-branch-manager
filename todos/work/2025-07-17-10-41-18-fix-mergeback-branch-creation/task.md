# Fix mergeback branch creation
**Status:** InProgress
**Agent PID:** 35219

## Original Todo
- fix mergeback branch creation
    *  Name:           MERGE_INGSVC-5638_production
        - merge is production -> master so the worktree name should be "MERGE_INGSVC-5638_master". append the target to the name, not the source.
    * Branch:         merge/INGSVC-5638_EMAIL_Invalid_Date_7_16_2025_13_00_00_production 
        - merge is production -> master so the branch name should be "merge/INGSVC-5638_EMAIL_Invalid_Date_7_16_2025_13_00_00_master". append the target to the branch name, not the source. 
    *  Base Branch:    production-2025-07-1
        - the workflow for mergeback is:
            1. create a merge branch based on the target
            2. merge source into target
            3. create pull request of merge branch into target
    * what we really want to do is run
    ```sh
     git log origin/<targetBranch>..origin/<sourceBranch> --oneline
```
        - whatever branch is associated with the most recent (first line) should be what we use for the branch naming convention
        - if there are no lines, that means <targetBranch> is up to date with <sourceBranch>

## Description
Fix the mergeback branch and worktree naming logic to correctly append the target branch name instead of the source branch name. Currently, when merging production → master, it creates "MERGE_INGSVC-5638_production" but should create "MERGE_INGSVC-5638_master".

The issue is in the `findNextBranchAndWorktreeInChain` function which finds the branch that merges INTO production (e.g., preview → production) rather than what production merges into (production → master).

## Implementation Plan
- [x] Fix the `findNextBranchAndWorktreeInChain` function in cmd/mergeback.go:217-229 to find what production merges into instead of what merges into production
- [x] **REWRITE APPROACH**: Replace the flawed merge logic with proper tree structure parsing using docs/worktree_nodes.md
- [x] Implement WorktreeManager tree parsing to identify the deepest leaf nodes (production branches)
- [x] Use `git log origin/<targetBranch>..origin/<sourceBranch> --oneline` to check if mergebacks are needed
- [x] Traverse from deepest leaf toward root until finding a branch that needs mergeback
- [x] Add tests for worktree_tree.go
- [x] Add automated tests to verify production → master creates names with "_master" suffix
- [ ] User test: Run mergeback command and verify branch/worktree names use target branch

## Notes
- Created `internal/worktree_tree.go` with full tree structure implementation from docs/worktree_nodes.md
- Embedded WorktreeManager in GBMConfig and auto-initialize in ParseGBMConfig
- Completely rewrote mergeback logic to use proper tree traversal:
  1. Find deepest leaf nodes (production branches)
  2. Use `git log origin/<target>..origin/<source> --oneline` to check for commits
  3. Traverse from deepest leaves toward root until finding a branch that needs mergeback
- This should now correctly identify production → master mergebacks and name them with "_master" suffix
- Added comprehensive tests in `cmd/mergeback_tree_test.go` 
- Key test `TestMergebackNamingProductionToMaster` passes, confirming the fix works for the main issue
- The test verifies: MERGE_INGSVC-5638_master (correct) vs MERGE_INGSVC-5638_production (incorrect)
- **FIXED CRITICAL ISSUE**: Updated WorktreeManager to allow multiple root nodes 
- This fixes failing sync tests that were encountering "multiple root nodes found" errors
- Real repositories often have multiple independent deployment chains (e.g., main + develop trees)