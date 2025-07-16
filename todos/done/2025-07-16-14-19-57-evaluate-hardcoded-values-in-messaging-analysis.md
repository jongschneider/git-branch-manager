# Hardcoded Values in User-Facing Messaging Analysis

Based on my analysis of the codebase, I found several categories of hardcoded values in user-facing messaging throughout the worktree-manager project. Here's a comprehensive breakdown:

## 1. `gbm info` Command Hardcoded Values

### **File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go`

**Line 331**: Hardcoded candidate branches for base branch detection:
```go
candidateBranches := []string{"main", "master", "develop", "dev"}
```

**Issue**: The base branch detection logic uses hardcoded branch names instead of consulting configuration or git defaults.

**Should be replaced with**: A configurable list from `.gbm/config.toml` or dynamically detected from git configuration.

### **File**: `/Users/jschneider/code/scratch/worktree-manager/internal/info_renderer.go`

**Line 106**: Hardcoded section headers with emojis:
```go
content.WriteString("📁 WORKTREE\n")      // Line 122
content.WriteString("🎫 JIRA TICKET\n")   // Line 146  
content.WriteString("🌿 GIT STATUS\n")    // Line 216
```

**Issue**: Section headers and emojis are hardcoded instead of being configurable.

**Should be replaced with**: Configurable section headers and icons through the existing icon system.

## 2. Print Function Messages Throughout Codebase

### **Command Help Text Issues**

**File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/root.go`
- **Line 22**: References outdated `.envrc` configuration:
  ```go
  Short: "Git Branch Manager - Manage Git worktrees based on .envrc configuration"
  ```
- **Line 24**: References `.envrc` in Long description:
  ```go
  Long: `...based on environment variables defined in a .envrc file.`
  ```

**Should be replaced with**: References to `gbm.branchconfig.yaml` configuration.

### **Shell Integration Hardcoded Values**

**File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/shell-integration.go`
- **Lines 32-39**: Hardcoded shell completion commands:
  ```go
  source <(gbm completion zsh)
  source <(gbm completion bash)
  ```
- **Line 29**: Hardcoded environment variable name:
  ```go
  export GBM_SHELL_INTEGRATION=1
  ```

**Should be replaced with**: Configurable shell integration settings.

### **Completion Help Text**

**File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/completion.go`
- **Lines 21-46**: Hardcoded file paths for completion installation:
  ```go
  gbm completion bash > /etc/bash_completion.d/gbm
  gbm completion bash > /usr/local/etc/bash_completion.d/gbm
  gbm completion zsh > "${fpath[1]}/_gbm"
  gbm completion fish > ~/.config/fish/completions/gbm.fish
  ```

**Should be replaced with**: Dynamic path detection or configurable paths.

## 3. Configuration File References

### **Hardcoded Configuration Filenames**

**File**: `/Users/jschneider/code/scratch/worktree-manager/internal/config.go`
- **Line 15**: Hardcoded configuration filename:
  ```go
  DefaultBranchConfigFilename = "gbm.branchconfig.yaml"
  ```

**File**: Multiple locations throughout codebase referencing hardcoded paths:
- `.gbm/config.toml`
- `.gbm/state.toml`
- `gbm.branchconfig.yaml`

**Issue**: These filenames are hardcoded and appear in error messages.

**Should be replaced with**: Constants that can be easily changed and referenced consistently.

## 4. Default Branch and Directory Names

### **Default Values in Config**

**File**: `/Users/jschneider/code/scratch/worktree-manager/internal/config.go`
- **Line 14**: Hardcoded default directory:
  ```go
  DefaultWorktreeDirname = "worktrees"
  ```
- **Lines 87-88**: Hardcoded prefixes:
  ```go
  HotfixPrefix:    "HOTFIX",
  MergebackPrefix: "MERGE",
  ```

**Issue**: Default values appear in user messages and should be configurable.

## 5. Git Status and Icons

### **Status Messages**

**File**: `/Users/jschneider/code/scratch/worktree-manager/internal/info_renderer.go`
- **Lines 308-320**: Hardcoded git status descriptions:
  ```go
  return "🔴 Unknown"
  return fmt.Sprintf("🟡 DIRTY (%d files modified)", fileCount)
  return fmt.Sprintf("🟠 DIVERGED (↑%d ↓%d)", status.Ahead, status.Behind)
  return "🟢 CLEAN"
  ```

**Should be replaced with**: Configurable status messages and icons.

### **Priority Icons**

**File**: `/Users/jschneider/code/scratch/worktree-manager/internal/info_renderer.go`
- **Lines 336-352**: Hardcoded priority formatting:
  ```go
  case "critical", "highest":
      return "🔴 Critical"
  case "high":
      return "🟠 High"
  case "medium":
      return "🟡 Medium"
  case "low":
      return "🟢 Low"
  case "lowest":
      return "🔵 Lowest"
  ```

**Should be replaced with**: Configurable priority icons and labels.

## 6. Hardcoded Branch References in Messages

### **Command Examples in Help Text**

Multiple command help texts contain hardcoded branch names:
- `main` branch references in examples
- `master` branch references
- `develop` branch references
- `feature/*` branch patterns

**Files with hardcoded branch examples**:
- `/Users/jschneider/code/scratch/worktree-manager/cmd/add.go`
- `/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go`
- `/Users/jschneider/code/scratch/worktree-manager/cmd/mergeback.go`

## 7. Time and Duration Formatting

### **Hardcoded Time Descriptions**

**File**: `/Users/jschneider/code/scratch/worktree-manager/internal/styles.go`
- **Lines 244-258**: Hardcoded time descriptions:
  ```go
  return "just now"
  return fmt.Sprintf("%d minutes", minutes)
  return fmt.Sprintf("%d hours", hours)
  return "1 day"
  return fmt.Sprintf("%d days", days)
  ```

**Should be replaced with**: Configurable time formatting strings.

## 8. Error Messages with Hardcoded Paths

### **Error Messages**

Throughout the codebase, error messages contain hardcoded references to:
- File paths (`.gbm/config.toml`, `gbm.branchconfig.yaml`)
- Directory names (`worktrees`)
- Command names (`gbm`)

**Example from** `/Users/jschneider/code/scratch/worktree-manager/internal/config.go`:
```go
return nil, fmt.Errorf("failed to read %s file: %w", DefaultBranchConfigFilename, err)
```

## 9. Deployment Chain Message (Previously Fixed)

The hotfix command deployment chain message has been previously addressed, but it demonstrates the type of hardcoded value that was problematic:

**File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go`
- **Line 107**: Now uses dynamic chain building:
  ```go
  PrintInfo("Remember to merge back through the deployment chain: %s", deploymentChain)
  ```

## Recommendations

### **Immediate Actions**
1. **Update help text** to remove `.envrc` references and use `gbm.branchconfig.yaml`
2. **Make section headers configurable** in info renderer
3. **Replace hardcoded candidate branches** with configurable list
4. **Make file paths in completion help dynamic** or configurable

### **Configuration Enhancement**
1. **Add messaging section** to `.gbm/config.toml` for customizable messages
2. **Extend icon configuration** to include section headers and status messages
3. **Add branch detection preferences** to configuration
4. **Make time formatting configurable**

### **Code Structure Improvements**
1. **Create message constants** file for all user-facing strings
2. **Implement message templating system** for dynamic values
3. **Add localization support** for future internationalization
4. **Create consistent error message formatting** system

The scope of hardcoded values is significant and affects user experience, customization, and maintainability. A systematic approach to making these values configurable would greatly improve the tool's flexibility and user experience.