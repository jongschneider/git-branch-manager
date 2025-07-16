# Analysis of `gbm hotfix` Command Issues

## Problem Summary
The `gbm hotfix` command has two main issues:
1. **Second argument autocompletion** shows JIRA keys instead of JIRA summaries
2. **Branch naming** includes incorrect prefix (`hotfix/bug/` instead of `hotfix/`)

## Code Analysis

### Current Hotfix Command Implementation
**File:** `/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go`

**Autocompletion Logic (Lines 89-111):**
```go
cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    if len(args) == 0 || len(args) == 1 {
        // Both first and second arguments: JIRA keys with summaries
        jiraIssues, err := internal.GetJiraIssues(manager)
        if err != nil {
            return nil, cobra.ShellCompDirectiveNoFileComp
        }

        var completions []string
        for _, issue := range jiraIssues {
            completion := fmt.Sprintf("%s\t%s", issue.Key, issue.Summary)
            completions = append(completions, completion)
        }
        return completions, cobra.ShellCompDirectiveNoFileComp
    }
    return nil, cobra.ShellCompDirectiveNoFileComp
}
```

### Correct Add Command Implementation
**File:** `/Users/jschneider/code/scratch/worktree-manager/cmd/add.go`

**Autocompletion Logic (Lines 111-155):**
```go
cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    if len(args) == 0 {
        // First argument: JIRA keys with summaries
        jiraIssues, err := internal.GetJiraIssues(manager)
        if err != nil {
            return nil, cobra.ShellCompDirectiveNoFileComp
        }

        var completions []string
        for _, issue := range jiraIssues {
            completion := fmt.Sprintf("%s\t%s", issue.Key, issue.Summary)
            completions = append(completions, completion)
        }
        return completions, cobra.ShellCompDirectiveNoFileComp
    } else if len(args) == 1 {
        // Second argument: branch name suggestions based on JIRA key
        worktreeName := args[0]
        if internal.IsJiraKey(worktreeName) {
            branchName, err := internal.GenerateBranchFromJira(worktreeName, manager)
            if err != nil {
                branchName = fmt.Sprintf("feature/%s", strings.ToLower(worktreeName))
            }
            return []string{branchName}, cobra.ShellCompDirectiveNoFileComp
        }
        return nil, cobra.ShellCompDirectiveNoFileComp
    }
    return nil, cobra.ShellCompDirectiveNoFileComp
}
```

### Branch Name Generation Issue
**File:** `/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go`

**Function:** `generateHotfixBranchName()` (Lines 183-217)

**Current Logic:**
1. Calls `internal.GenerateBranchFromJira(jiraTicket, manager)`
2. **Problem:** Replaces `feature/` or `bugfix/` prefixes with `hotfix/`
3. **Issue:** JIRA generates `bugfix/INGSVC-5638_EMAIL_Invalid_Date_7_16_2025_13_00_00`
4. **Result:** Becomes `hotfix/bug/INGSVC-5638_EMAIL_Invalid_Date_7_16_2025_13_00_00`

**JIRA Branch Generation:**
**File:** `/Users/jschneider/code/scratch/worktree-manager/internal/jira.go`

**Function:** `BranchName()` (Lines 220-234)
```go
func (j *JiraIssue) BranchName() string {
    summary := strings.ReplaceAll(j.Summary, " ", "_")
    summary = strings.ReplaceAll(summary, "-", "_")
    summary = regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(summary, "_")
    summary = regexp.MustCompile(`_+`).ReplaceAllString(summary, "_")
    summary = strings.Trim(summary, "_")

    issueType := strings.ToLower(j.Type)
    if issueType == "story" || issueType == "improvement" {
        issueType = "feature"
    }

    return fmt.Sprintf("%s/%s_%s", issueType, j.Key, summary)
}
```

## Issues Identified

### Issue 1: Second Argument Autocompletion
**Problem:** Shows JIRA keys instead of JIRA summaries
**Current Behavior:** `gbm hf INGSVC-5638 <TAB>` shows JIRA keys again
**Expected Behavior:** Should show JIRA summary like `gbm add` does

### Issue 2: Branch Naming Prefix Replacement
**Problem:** Incorrect prefix replacement in `generateHotfixBranchName()`
**Current Logic:** 
```go
branchName = strings.ReplaceAll(branchName, "feature/", "hotfix/")
branchName = strings.ReplaceAll(branchName, "bugfix/", "hotfix/")
```
**Issue:** When JIRA generates `bugfix/INGSVC-5638_...`, replacement becomes `hotfix/bug/INGSVC-5638_...`
**Expected:** Should be `hotfix/INGSVC-5638_...`

## Command Usage Pattern Analysis

### Current hotfix command usage:
```bash
gbm hotfix <worktree-name> [jira-ticket]
```

### Expected behavior based on todo:
```bash
gbm hf INGSVC-5638 INGSVC-5638  # Second arg should show JIRA summary, not key
```

**Note:** The usage pattern suggests first argument is JIRA key, second argument should provide JIRA summary context, similar to `gbm add`.