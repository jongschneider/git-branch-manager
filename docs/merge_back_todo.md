# Merge-back Detection and Alerting System - Implementation Todo List

## 1. Foundation Setup

### 1.1 Data Structures
- [x] Create `MergeBackStatus` struct with `MergeBacksNeeded` and `HasUserCommits` fields
- [x] Create `MergeBackInfo` struct with branch information and commit details
- [x] Create `CommitInfo` struct with hash, message, author, timestamp, and user flag

### 1.2 Core Package Structure
- [x] Create internal package for merge-back detection logic
- [x] Set up error handling types for graceful degradation
- [x] Create configuration constants for git commands and timeouts

## 2. Environment Configuration Parsing

### 2.1 .envrc File Processing
- [x] Implement `.envrc` file parser that preserves line order
- [x] Handle comments and empty lines correctly
- [x] Extract environment variable to branch mappings
- [x] Build hierarchy chain from file order (bottom entries merge into entries above)

### 2.2 Branch Hierarchy Logic
- [x] Create function to determine merge-back hierarchy from `.envrc` order
- [x] Handle dynamic hierarchy (branches may be added/removed)
- [x] Validate that mapped branches exist in repository
- [x] Skip non-existent branches without failing the process

## 3. Git Operations Integration

### 3.1 Commit Analysis
- [x] Implement `git log TARGET..SOURCE --format=%H|%s|%an|%ae|%ct` command execution
- [x] Parse git log output to extract commit information
- [x] Handle empty results gracefully
- [x] Extract commit hash, message, author, email, and timestamp

### 3.2 User Identification
- [x] Get current user email from `git config user.email`
- [x] Implement fallback to `git config user.name` if email unavailable
- [x] Handle missing git configuration gracefully
- [x] Implement case-sensitive exact matching for user identification

### 3.3 Git Command Error Handling
- [x] Graceful degradation on git command failures
- [x] Silent skip when not in git repository
- [x] Brief warning for git command failures
- [x] Brief warning for branch access issues

## 4. Detection Logic Implementation

### 4.1 Core Detection Algorithm
- [x] Check if in git repository with `.envrc` file
- [x] Parse `.envrc` to get ordered environment variables
- [x] Verify each mapped branch exists
- [x] For each adjacent pair in hierarchy, execute commit analysis
- [x] Identify user-authored commits vs. other commits
- [x] Build MergeBackStatus with complete information

### 4.2 Performance Optimization
- [x] Ensure merge-back check completes within 500ms for typical repositories
- [x] Implement efficient git command execution
- [x] Handle large repositories appropriately

## 5. Alert Display System

### 5.1 Alert Formatting
- [x] Design alert format with warning icon and clear messaging
- [x] Include commit details for user-authored commits
- [x] Format output for terminal readability
- [x] Show commit count and user-specific count

### 5.2 Alert Content
- [x] Display branch pairs requiring merge-back (e.g., "PROD → PREVIEW")
- [x] Show total commit count and user commit count
- [x] Include commit hash, message, and timestamp for user commits
- [x] Use relative time formatting (e.g., "2 days ago")

### 5.3 Alert Integration
- [x] Integrate with root command's `PersistentPreRun` function
- [x] Display alerts before primary command output
- [x] Output alerts to stderr to avoid interfering with command output
- [x] Ensure alerts don't prevent primary gbm command execution

## 6. Error Handling and Edge Cases

### 6.1 File System Errors
- [x] Silent skip when `.envrc` doesn't exist
- [x] Handle corrupted or unreadable `.envrc` files
- [x] Handle permission issues with git repository access

### 6.2 Git Repository Errors
- [x] Silent skip when not in git repository
- [x] Handle corrupted git repositories
- [x] Handle missing git configuration scenarios

### 6.3 Edge Cases
- [x] Handle empty `.envrc` files
- [x] Handle single branch configurations
- [x] Handle non-existent branches in hierarchy
- [x] Handle repositories with no commits

## 7. Documentation and Examples

### 7.1 Alert Output Examples
- [x] Document clean state (no alerts) behavior
- [x] Document single user commit alert format
- [x] Document multiple branches with mixed authors format
- [x] Document complex environment setup examples

### 7.2 Configuration Examples
- [x] Document standard three-tier `.envrc` setup
- [x] Document dynamic two-tier setup
- [x] Document complex environment setup with multiple tiers

## 8. Performance and Compatibility

### 8.1 Performance Requirements
- [x] Ensure 500ms execution time target is met
- [x] Profile git command execution performance
- [x] Optimize for repositories with many branches

### 8.2 Integration Compatibility
- [x] Ensure compatibility with existing gbm functionality
- [x] Test integration with all existing gbm commands
- [x] Verify no breaking changes to current workflows

## 9. Final Integration and Validation

### 9.1 Command Integration
- [x] Integrate detection logic with all gbm commands
- [x] Test with `gbm sync`, `gbm status`, `gbm list`, `gbm validate`
- [x] Ensure consistent behavior across all commands

### 9.2 User Experience Validation
- [x] Test with realistic merge-back scenarios
- [x] Validate alert clarity and actionability
- [x] Test with different terminal environments
- [x] Verify performance in real-world repositories
