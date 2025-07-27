# Replace duplicate branch status calls in cmd/info.go with new utilities
**Status:** Done
**Agent PID:** 52831

## Original Todo
**Replace duplicate branch status calls in cmd/info.go with new utilities**
- Update lines 291, 300, 309 to use new GitManager methods
- Test: `gbm info` command still works correctly

## Description
Review and verify that branch status calls in cmd/info.go have been properly converted to use GitManager utilities, and identify any remaining direct git commands that could benefit from additional utilities. The main target lines (291, 300, 309) appear to already use GetCurrentBranchInPath, GetUpstreamBranch, and GetAheadBehindCount utilities.

## Implementation Plan
- [x] Verify that lines 291, 300, 309 in cmd/info.go are already using GitManager utilities (GetCurrentBranchInPath, GetUpstreamBranch, GetAheadBehindCount)
- [x] Check for any remaining direct git commands in cmd/info.go that could be converted to utilities
- [x] Identify any inconsistencies where some calls use utilities and others use direct git commands
- [x] Test that `gbm info` command works correctly with current utility usage
- [x] Run validation commands: `just lint`, `just format`, `just test-changed`

## Notes

### Findings

✅ **Target lines already converted:**
- Line 291: `manager.GetGitManager().GetCurrentBranchInPath(worktreePath)` ✅
- Line 298: `manager.GetGitManager().GetUpstreamBranch(worktreePath)` ✅  
- Line 304: `manager.GetGitManager().GetAheadBehindCount(worktreePath)` ✅

❌ **Remaining direct git commands in cmd/info.go:**
- Line 161: `git log` for recent commits
- Line 197: `git diff --numstat` for unstaged changes  
- Line 238: `git diff --cached --numstat` for staged changes
- Line 327: `git rev-parse --verify` for branch verification
- Line 331: `git merge-base --is-ancestor` for ancestry checking

⚠️ **Inconsistencies found:**
- internal/git.go:685,719 - PushWorktree/PullWorktree still use direct `ExecGitCommand` for getting current branch instead of `GetCurrentBranchInPath`

### Conclusion
The original todo task (lines 291, 300, 309) is **already complete**. The utilities are working correctly.