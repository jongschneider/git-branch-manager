# fix print formatting
**Status:** Done
**Agent PID:** 98877

## Original Todo
- fix print formatting
    *   gbm list
┌────────────────────┬──────────────────────────────────────────────────────────┬────────────┬─────────────┬────────────────────────────────────────────────────────────────────────────┐
│      WORKTREE      │                          BRANCH                          │ GIT STATUS │ SYNC STATUS │                                    PATH                                    │
├────────────────────┼──────────────────────────────────────────────────────────┼────────────┼─────────────┼────────────────────────────────────────────────────────────────────────────┤
│ master             │ master                                                   │ ✓          │ ✅ IN_SYNC  │ /Users/jschneider/code/scratch/email_ingester/worktrees/master             │
│ production         │ production-2025-07-1                                     │ ✓          │ ✅ IN_SYNC  │ /Users/jschneider/code/scratch/email_ingester/worktrees/production         │
│ HOTFIX_INGSVC-5638 │ hotfix/INGSVC-5638_EMAIL_Invalid_Date_7_16_2025_13_00_00 │ ✓          │ UNTRACKED   │ /Users/jschneider/code/scratch/email_ingester/worktrees/HOTFIX_INGSVC-5638 │
└────────────────────┴──────────────────────────────────────────────────────────┴────────────┴─────────────┴────────────────────────────────────────────────────────────────────────────┘

💡 Run 'gbm sync' to synchronize changes%
    * note the "%" printed at the end of "changes"
    * research where else this is happening and fix it.

## Description
Fix print formatting issues where missing newlines cause shell prompts to display incorrectly. The main issue is in `gbm list` command output where a "%" character appears after "changes" due to missing newlines in fmt.Fprint statements. This happens because shells like zsh show "%" to indicate lines that don't end with newlines.

## Implementation Plan
- [x] Fix main formatting issue in cmd/list.go:97 - change fmt.Fprint to fmt.Fprintln for sync message
- [x] Fix alert formatting issue in cmd/root.go:219 - change fmt.Fprint to fmt.Fprintln for alert messages
- [x] Fix corresponding issues in worktree copies (worktrees/HOTFIX_INGSVC-5638/cmd/list.go:97 and root.go:219)
- [x] Verify info.go:155 output formatting is correct (may need newline)
- [x] Remove "\n" from RenderWorktreeInfo and change fmt.Print to fmt.Println in info.go for consistency
- [x] Apply same changes to worktree copy (info_renderer.go and info.go)
- [x] Automated test: Run `go test ./...` to ensure changes don't break existing functionality
- [x] User test: Build and run `gbm list` command to verify no "%" appears after messages

## Notes
- Applied consistent strategy: Use fmt.Fprintln everywhere, remove manual "\n" from strings
- Fixed 6 files total: 2 main files + 4 worktree copies
- Build succeeded and gbm list command no longer shows "%" character after messages
- All formatting now uses fmt.Fprintln/fmt.Println consistently