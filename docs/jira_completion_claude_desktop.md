This is a great enhancement! Let me think through a few approaches for integrating JIRA tickets into your COBRA CLI with good UX.

## Approach 1: Dynamic Shell Completions (Recommended)

COBRA supports dynamic completions that can call external commands. You can have your CLI fetch JIRA issues on-demand during tab completion:

```go
// In your add command setup
addCmd.RegisterFlagCompletionFunc("branch", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    if len(args) == 1 {
        // args[0] is the worktree name (JIRA key)
        return generateBranchSuggestions(args[0]), cobra.ShellCompDirectiveNoFileComp
    }
    return nil, cobra.ShellCompDirectiveNoFileComp
})

// For the positional argument (JIRA key)
addCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
    if len(args) == 0 {
        return getJiraKeys(), cobra.ShellCompDirectiveNoFileComp
    }
    return nil, cobra.ShellCompDirectiveNoFileComp
}

func getJiraKeys() []string {
    cmd := exec.Command("jira", "issue", "list", "-a$(jira me)", "--plain")
    output, err := cmd.Output()
    if err != nil {
        return nil
    }
    
    var keys []string
    lines := strings.Split(string(output), "\n")
    for _, line := range lines[1:] { // Skip header
        if line = strings.TrimSpace(line); line != "" {
            fields := strings.Fields(line)
            if len(fields) >= 2 {
                keys = append(keys, fields[1]) // The KEY column
            }
        }
    }
    return keys
}

func generateBranchSuggestions(jiraKey string) []string {
    // Fetch the specific issue details to get type and summary
    cmd := exec.Command("jira", "issue", "view", jiraKey, "--plain")
    output, err := cmd.Output()
    if err != nil {
        return nil
    }
    
    // Parse output to extract type and summary
    // Then format as: bug/INGSVC-5544_Zoom_Clips_-_New_Integration
    issueType, summary := parseJiraIssue(output)
    branchName := fmt.Sprintf("%s/%s_%s", 
        strings.ToLower(issueType), 
        jiraKey, 
        strings.ReplaceAll(summary, " ", "_"))
    
    return []string{branchName}
}
```

## Approach 2: Smart Default Behavior

Make the command work intuitively without requiring the `-b` flag when using JIRA keys:

```go
func runAddCommand(cmd *cobra.Command, args []string) error {
    worktreeName := args[0]
    branchName, _ := cmd.Flags().GetString("branch")
    
    // Check if worktree name looks like a JIRA key
    if isJiraKey(worktreeName) && branchName == "" {
        // Auto-generate branch name from JIRA
        generatedBranch, err := generateBranchFromJira(worktreeName)
        if err != nil {
            return fmt.Errorf("failed to fetch JIRA issue: %w", err)
        }
        branchName = generatedBranch
        fmt.Printf("Auto-detected JIRA issue, using branch: %s\n", branchName)
    }
    
    return createWorktree(worktreeName, branchName)
}

func isJiraKey(s string) bool {
    // Simple regex to match JIRA key pattern
    matched, _ := regexp.MatchString(`^[A-Z]+-\d+$`, s)
    return matched
}
```

## Approach 3: Interactive Selection

Add a `--interactive` or `--jira` flag that presents a nice selection interface:

```go
func runAddWithJira(cmd *cobra.Command, args []string) error {
    issues, err := fetchJiraIssues()
    if err != nil {
        return err
    }
    
    // Use a library like survey or bubbletea for nice selection
    selected := promptForIssue(issues)
    
    worktreeName := selected.Key
    branchName := generateBranchName(selected)
    
    return createWorktree(worktreeName, branchName)
}
```

## Implementation Details

For the JIRA parsing, you'll want something like:

```go
type JiraIssue struct {
    Type    string
    Key     string
    Summary string
    Status  string
}

func parseJiraList(output string) []JiraIssue {
    var issues []JiraIssue
    lines := strings.Split(output, "\n")
    
    for _, line := range lines[1:] { // Skip header
        if line = strings.TrimSpace(line); line != "" {
            fields := strings.Fields(line)
            if len(fields) >= 4 {
                issue := JiraIssue{
                    Type:    fields[0],
                    Key:     fields[1],
                    Summary: strings.Join(fields[2:len(fields)-1], " "),
                    Status:  fields[len(fields)-1],
                }
                issues = append(issues, issue)
            }
        }
    }
    return issues
}

func (j JiraIssue) BranchName() string {
    summary := strings.ReplaceAll(j.Summary, " ", "_")
    summary = strings.ReplaceAll(summary, "-", "_")
    return fmt.Sprintf("%s/%s_%s", strings.ToLower(j.Type), j.Key, summary)
}
```

## Recommended Combination

I'd suggest implementing **Approach 1 + Approach 2**: 
- Dynamic completions for great UX during typing
- Smart defaults so `gbm add INGSVC-5544` automatically generates the branch name
- Keep the `-b` flag for overrides

This gives you the best of both worlds - great discoverability through completions, but also quick usage for power users who know their JIRA keys.

Would you like me to elaborate on any of these approaches or help with the specific parsing logic?
