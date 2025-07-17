Perfect! I found the issue and can see that the code in `cmd/root.go` is correctly using `fmt.Fprintf` with `%s\n` to add a newline, but the problematic line in `cmd/list.go` is using `fmt.Fprint` without a newline.

## Summary

I found the print formatting issue in the worktree-manager codebase. Here are the exact locations and details:

### **Main Issue Location:**
- **File**: `/Users/jschneider/code/scratch/worktree-manager/cmd/list.go`
- **Line**: 97
- **Code**: `_, _ = fmt.Fprint(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))`

### **Problem:**
The issue is that `fmt.Fprint` is being used instead of `fmt.Fprintln` or `fmt.Fprintf` with a newline character. This causes the output to not end with a newline, and in shells like zsh, this results in a `%` character being displayed to indicate that the line doesn't end with a newline.

### **Root Cause:**
The `internal.FormatInfo()` function returns a formatted string with an icon and text, but it doesn't include a trailing newline. When this formatted string is printed using `fmt.Fprint`, no newline is added, causing the shell prompt to appear immediately after the text.

### **Correct Usage Example:**
In `/Users/jschneider/code/scratch/worktree-manager/cmd/root.go` at lines 84 and 107, the code correctly uses:
```go
fmt.Fprintf(os.Stderr, "%s\n", internal.FormatInfo(msg))
fmt.Fprintf(os.Stderr, "%s\n", internal.FormatError("ERROR: "+msg))
```

### **Other Locations:**
The same issue exists in the worktree copy at `/Users/jschneider/code/scratch/worktree-manager/worktrees/HOTFIX_INGSVC-5638/cmd/list.go` line 97.

### **Solution:**
Change line 97 in `cmd/list.go` from:
```go
_, _ = fmt.Fprint(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))
```

To either:
```go
_, _ = fmt.Fprintln(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))
```
or:
```go
_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", internal.FormatInfo("Run 'gbm sync' to synchronize changes"))
```

This will ensure that the output ends with a proper newline character, preventing the `%` character from appearing in the shell.

Based on my comprehensive search of the codebase, I found several instances of printing functions that don't include newlines. Here's a detailed breakdown:

## Summary of Findings

### 1. `fmt.Fprint` instances (4 total):

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/list.go`**
- **Line 90**: `_, _ = fmt.Fprint(cmd.OutOrStdout(), table.String())`
- **Line 97**: `_, _ = fmt.Fprint(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))`

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/root.go`**
- **Line 219**: `fmt.Fprint(os.Stderr, alert)`

**File: `/Users/jschneider/code/scratch/worktree-manager/worktrees/HOTFIX_INGSVC-5638/cmd/list.go`**
- **Line 90**: `_, _ = fmt.Fprint(cmd.OutOrStdout(), table.String())`
- **Line 97**: `_, _ = fmt.Fprint(cmd.OutOrStdout(), internal.FormatInfo("Run 'gbm sync' to synchronize changes"))`

**File: `/Users/jschneider/code/scratch/worktree-manager/worktrees/HOTFIX_INGSVC-5638/cmd/root.go`**
- **Line 219**: `fmt.Fprint(os.Stderr, alert)`

### 2. `fmt.Print` instances (10 total):

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/info.go`**
- **Line 155**: `fmt.Print(output)` - This outputs rendered worktree info

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/shell-integration.go`**
- **Line 85**: `fmt.Print(shellCode)` - This outputs shell integration script

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/add.go`**
- **Line 180**: `fmt.Print(internal.FormatPrompt("Select a branch: "))` - Interactive prompt
- **Line 191**: `fmt.Print(internal.FormatPrompt("Enter new branch name: "))` - Interactive prompt

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/switch.go`**
- **Line 107**: `fmt.Print(targetPath)` - Outputs path for shell integration

**File: `/Users/jschneider/code/scratch/worktree-manager/cmd/sync.go`**
- **Line 84**: `fmt.Print(message + " [y/N]: ")` - Interactive confirmation prompt

Plus corresponding duplicates in the worktrees folder.

### 3. Other Write methods:
- Only found in documentation files (not actual code)
- Test files using `w.Write([]byte(input + "\n"))` which properly includes newlines

## Analysis and Recommendations

### Issues that need fixing:

1. **`/Users/jschneider/code/scratch/worktree-manager/cmd/list.go` line 97**:
   - Same issue as we already identified - info message without newline
   - Should be `fmt.Fprintln` or add `\n` to the message

2. **`/Users/jschneider/code/scratch/worktree-manager/cmd/root.go` line 219**:
   - Alert message without newline - could cause formatting issues
   - Should likely be `fmt.Fprintln(os.Stderr, alert)`

### Instances that are intentionally correct:

1. **Interactive prompts** (lines in add.go, sync.go):
   - These are meant to stay on the same line for user input
   - Correct behavior - no newline needed

2. **Shell integration output** (shell-integration.go):
   - Script output that should not have extra newlines
   - Correct behavior

3. **Path output** (switch.go):
   - Used for shell integration, should not have newlines
   - Correct behavior

4. **Table output** (list.go line 90):
   - Table already includes proper formatting
   - Has a subsequent `fmt.Fprintln` to add spacing
   - Correct behavior

5. **Info output** (info.go):
   - Rendered output that likely already includes proper formatting
   - May be correct, but should be verified

The most critical issues are:
- **Line 97 in list.go** (both main and worktree copies)
- **Line 219 in root.go** (both main and worktree copies)

These should be changed to include newlines to prevent formatting issues in terminal output.