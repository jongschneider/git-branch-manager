# TODOs

## Git Command Deduplication Project

**Gameplan Reference**: See `docs/git_command_deduplication_gameplan.md` for full implementation strategy.

### Phase 1: High Priority - Branch Status Utilities






### Phase 2: Medium Priority - Repository Introspection


- [ ] **Add GetCommitHash utility function to GitManager**
  - Location: `internal/git.go`
  - Replace pattern: `git rev-parse <ref>`
  - Files affected: `internal/git_add.go:48,54`

- [ ] **Replace duplicate rev-parse calls in cmd/info.go with new utilities**
  - Update line 338 to use VerifyRef
  - Test: `gbm info` command still works correctly

- [ ] **Replace duplicate rev-parse calls in internal/git_add.go with new utilities**
  - Update lines 48, 54, 89 to use new utilities
  - Test: `gbm add` command still works correctly

### Phase 3: Low Priority - Info Command Extraction

- [ ] **Add GetCommitHistory utility function to GitManager**
  - Location: `internal/git.go`
  - Replace pattern: `git log --oneline --format=...`
  - Files affected: `cmd/info.go:161`

- [ ] **Add GetFileChanges utility function to GitManager**
  - Location: `internal/git.go`
  - Replace pattern: `git diff --numstat`, `git diff --cached --numstat`
  - Files affected: `cmd/info.go:197,238`

- [ ] **Extract git log calls from cmd/info.go to use GetCommitHistory**
  - Update line 161 to use new utility
  - Test: `gbm info` shows correct commit history

- [ ] **Extract git diff calls from cmd/info.go to use GetFileChanges**
  - Update lines 197, 238 to use new utility
  - Test: `gbm info` shows correct file changes

### Validation

- [ ] **Run full test suite to validate all Git command abstractions**
  - Command: `just test` or `go test ./...`
  - Ensure: All tests pass after all refactoring
  - Verify: No direct `exec.Command("git", ...)` calls outside utilities

## Notes

- Work through each todo sequentially to minimize risk
- Run tests after each change to catch regressions early
- Reference the gameplan document for implementation details and guidelines
- Each utility function should use `enhanceGitError()` for consistent error handling