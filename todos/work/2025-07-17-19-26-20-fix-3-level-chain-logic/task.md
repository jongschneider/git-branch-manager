# Fix 3-level chain logic in mergeback functionality
**Status:** AwaitingCommit
**Agent PID:** 35219

## Original Todo
Fix 3-level chain logic in mergeback functionality - The `TestMergebackNamingWithTreeStructure` test fails because 3-level chains (production → preview → master) incorrectly return "master" instead of "preview" for the target branch

## Description
Fix parameter order bug in mergeback 3-level chain logic. The `hasCommitsBetweenBranches` function calls in `cmd/mergeback.go` have reversed parameters, causing the algorithm to incorrectly detect merge requirements and skip intermediate steps in the chain (production → preview → master), resulting in "master" being returned instead of "preview" as the merge target.

## Implementation Plan
- [x] Fix parameter order in line 159 of `cmd/mergeback.go`: change `hasCommitsBetweenBranches(gitManager, leaf.Parent.Config.Branch, leaf.Config.Branch)` to `hasCommitsBetweenBranches(gitManager, leaf.Config.Branch, leaf.Parent.Config.Branch)` 
- [x] Fix parameter order in line 206 of `cmd/mergeback.go`: change `hasCommitsBetweenBranches(gitManager, leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch)` to `hasCommitsBetweenBranches(gitManager, leaf.Parent.Config.Branch, leaf.Parent.Parent.Config.Branch)`
- [x] Investigate and fix integration test failures caused by parameter order changes - the mergeback detection is now returning multiple results when expecting single results, and detecting mergebacks where none should exist
- [x] Automated test: Run `go test -run TestMergebackNamingWithTreeStructure` and `go test -run TestFindMergeTargetWithTreeStructure` to verify the fixes work
- [x] User test: Manually test mergeback command with 3-level chain to verify correct branch/worktree naming

## Notes
Additional fix: Resolved nil pointer dereference in `CheckMergeBackStatus` when config file is empty by ensuring function always returns a valid status object instead of nil.