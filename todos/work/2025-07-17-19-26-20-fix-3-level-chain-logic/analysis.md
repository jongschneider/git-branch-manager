Perfect! Now I found the root cause. Here's the detailed analysis:

## Root Cause Analysis

### The Problem
The 3-level chain logic is failing because there's a **parameter order mismatch** in the `hasCommitsBetweenBranches` function call.

### Expected 3-Level Chain
From the test setup:
- **production** (production-2025-05-1) → **preview** (production-2025-07-1) → **master** (master)

The test expects:
1. Production has new commits that need to merge to preview
2. So the merge target should be preview (branch: "production-2025-07-1", worktree: "preview")

### The Bug

In line 159 of `cmd/mergeback.go`:
```go
hasCommits, err := hasCommitsBetweenBranches(gitManager, leaf.Parent.Config.Branch, leaf.Config.Branch)
```

This calls:
```go
hasCommitsBetweenBranches(gitManager, "production-2025-07-1", "production-2025-05-1")
```

But the function signature and implementation expects:
```go
hasCommitsBetweenBranches(gitManager, targetBranch, sourceBranch string)
```

The git command then becomes:
```bash
git log production-2025-07-1..production-2025-05-1 --oneline
```

This checks if `production-2025-05-1` has commits that `production-2025-07-1` doesn't have - **which is backwards!**

### What Should Happen

The function should be called as:
```go
hasCommitsBetweenBranches(gitManager, leaf.Config.Branch, leaf.Parent.Config.Branch)
```

Which would become:
```bash
git log production-2025-07-1..production-2025-05-1 --oneline  
```

Wait, that's not right either. Let me check the git log syntax...

Actually, `git log A..B` shows commits in B that are not in A. So if we want to check if source branch has commits that target branch doesn't have, we want `git log target..source`.

### The Correct Fix

The issue is in line 159. It should be:
```go
hasCommits, err := hasCommitsBetweenBranches(gitManager, leaf.Config.Branch, leaf.Parent.Config.Branch)
```

Because:
- `leaf.Config.Branch` is the source (production-2025-05-1) 
- `leaf.Parent.Config.Branch` is the target (production-2025-07-1)
- We want to check if source has commits that target doesn't have

### Why Tests Are Failing

1. **TestFindMergeTargetWithTreeStructure**: 
   - The function returns "master" instead of "production-2025-07-1" (preview)
   - This happens because `hasCommitsBetweenBranches` returns `false` when it should return `true`
   - So the algorithm skips to `findNextMergeTargetInChain` which likely also has the same bug

2. **TestMergebackNamingWithTreeStructure**:
   - Since the wrong target is found ("master" instead of "preview")
   - The branch name becomes "MERGE_INGSVC-5638_master" instead of "MERGE_INGSVC-5638_preview"

### Secondary Bug in findNextMergeTargetInChain

Looking at line 206, there's the same bug:
```go
hasCommits, err := hasCommitsBetweenBranches(gitManager, leaf.Parent.Parent.Config.Branch, leaf.Parent.Config.Branch)
```

This should also be reversed to check if the source (leaf.Parent) has commits that the target (leaf.Parent.Parent) doesn't have.

### Summary

The core issue is that **both calls to `hasCommitsBetweenBranches` have their parameters reversed**, causing the mergeback detection logic to fail and skip up the chain incorrectly, ultimately returning "master" instead of the correct intermediate step "preview".