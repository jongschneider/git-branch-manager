# Analysis: Fix hardcoded merge-back chain in hotfix command message

## Task Agents Analysis

### Agent 1: Find hardcoded merge-back chain

Based on my investigation, here's what I found about the hardcoded merge-back chain message in the hotfix command:

#### 1. Location of the Message

The message is located in `/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go` at **line 107**:

```go
PrintInfo("Remember to merge back through the deployment chain: %s", deploymentChain)
```

#### 2. How the Current Chain is Defined

The deployment chain is **not actually hardcoded** as "master → preview → main". Instead, it's dynamically built by the `buildDeploymentChain` function (lines 266-279) which:

1. Loads the `gbm.branchconfig.yaml` configuration file
2. Finds the production branch (the branch with no `merge_into` target)
3. Traverses the merge configuration to build the complete chain
4. Joins the chain with "→" arrows

The function `buildMergeChain` (lines 302-319) builds the chain by:
- Starting from the base branch (production branch)
- Finding branches that merge into the current branch
- Following the chain until no more merge targets exist

#### 3. How the gbm.branchconfig.yaml merge_into Configuration Works

The `gbm.branchconfig.yaml` file uses a YAML structure defined in `/Users/jschneider/code/scratch/worktree-manager/internal/config.go` (lines 70-78):

```go
type GBMConfig struct {
    Worktrees map[string]WorktreeConfig `yaml:"worktrees"`
}

type WorktreeConfig struct {
    Branch      string `yaml:"branch"`
    MergeInto   string `yaml:"merge_into,omitempty"`
    Description string `yaml:"description,omitempty"`
}
```

The `merge_into` field defines where a branch should be merged. The system builds a chain by:
- Finding the production branch (one with no `merge_into` field)
- Following the chain of branches that merge into each other

#### 4. Structure of the Branch Config

The branch configuration follows this pattern:

```yaml
worktrees:
  main:
    branch: main
    description: "Main production branch"
    # No merge_into - this is the root/production branch
  preview:
    branch: preview
    merge_into: main
    description: "Preview/staging branch"
  master:
    branch: master
    merge_into: preview
    description: "Master production branch"
```

In this example:
- `master` is the actual production branch (bottom of chain)
- `master` merges into `preview`
- `preview` merges into `main`
- `main` has no `merge_into` target (top of chain)

#### Key Issue Found

The current implementation has a logical issue in the `buildMergeChain` function. It builds the chain from production branch upward, but the message suggests it should show the merge flow from production → staging → main. The current code shows the reverse order because it finds branches that merge INTO the current branch, when it should show the flow FROM the current branch.

The `findBranchThatMergesInto` function (lines 321-329) finds branches that merge INTO a target, creating a chain that goes from production to final destination, which is the correct direction for showing merge-back flow.

### Agent 2: Analyze gbm.branchconfig.yaml structure

Based on my analysis of the codebase, I can now provide you with a comprehensive understanding of the gbm.branchconfig.yaml file structure and the deployment chain functionality.

## gbm.branchconfig.yaml File Structure

The `gbm.branchconfig.yaml` file uses the following structure:

```yaml
# Git Branch Manager Configuration

# Worktree definitions - key is the worktree name, value defines the branch and merge strategy
worktrees:
  worktree_name:
    branch: actual_branch_name
    merge_into: target_worktree_name  # Optional - defines merge hierarchy
    description: "Human readable description"  # Optional
```

### Key Components:

1. **Worktree Name**: The key used to identify the worktree (e.g., `main`, `preview`, `production`)
2. **Branch**: The actual git branch name associated with the worktree
3. **MergeInto**: Defines the merge hierarchy by specifying which worktree this one merges into
4. **Description**: Optional human-readable description

## Examples from the Codebase

### Example 1: Three-tier deployment chain
```yaml
worktrees:
  main:
    branch: main
    merge_into: ""  # Root - no merge target
    description: "Main branch"
  preview:
    branch: preview
    merge_into: "main"
    description: "Preview environment"
  production:
    branch: production
    merge_into: "preview"
    description: "Production environment"
```

### Example 2: Simple two-tier chain
```yaml
worktrees:
  main:
    branch: main
    merge_into: ""
    description: "Main branch"
  production:
    branch: production
    merge_into: "main"
    description: "Production environment"
```

## How merge_into is configured

The `merge_into` field establishes the deployment chain:
- **Empty string or omitted**: Indicates this is the root/final branch in the chain
- **Target worktree name**: Points to the worktree that this branch should merge into

## Deployment Chain Logic

The deployment chain is built using two key functions:

### 1. `buildDeploymentChain(baseBranch, manager)`
- Takes a base branch and builds the complete deployment chain
- Returns a string representation like "production → preview → main"

### 2. `buildMergeChain(baseBranch, config)`
- Traverses the merge configuration to build the complete chain
- Starts from the base branch and follows the `merge_into` relationships
- Returns an array of branch names in deployment order

## Test Examples

From the test files, here are actual configurations being used:

### Standard GBM Config from tests:
```yaml
worktrees:
  main:
    branch: main
    description: "Main branch"
  preview:
    branch: preview
    merge_into: "main"
    description: "Preview branch"
  staging:
    branch: staging
    merge_into: "preview"
    description: "Staging branch"
  dev:
    branch: dev
    merge_into: "staging"
    description: "Dev branch"
  feat:
    branch: feat
    merge_into: "dev"
    description: "Feat branch"
  prod:
    branch: prod
    merge_into: "feat"
    description: "Prod branch"
```

### Production Branch Finding Logic
The code includes logic to find the "production branch" (the branch at the bottom of the chain):
- Looks for branches with no `merge_into` target (root branches)
- Uses heuristics to identify production branches by name (prod, production, master, main, release)

## Key Files and Locations

- **/Users/jschneider/code/scratch/worktree-manager/internal/config.go**: Contains the `GBMConfig` and `WorktreeConfig` structs
- **/Users/jschneider/code/scratch/worktree-manager/cmd/hotfix.go**: Contains the `buildDeploymentChain` and `buildMergeChain` functions
- **/Users/jschneider/code/scratch/worktree-manager/internal/testutils/git_harness.go**: Contains the `CreateGBMConfig` function that shows how the YAML is generated

## Current Issue

The current task in the todos indicates there's a hardcoded deployment chain message that needs to be fixed. The problem is that the hotfix command shows a hardcoded message like "Remember to merge back through the deployment chain: master → preview → main" instead of dynamically determining the correct chain from the gbm.branchconfig.yaml file.

This is exactly what the `buildDeploymentChain` function is designed to solve - it dynamically reads the configuration and builds the proper deployment chain based on the `merge_into` relationships.