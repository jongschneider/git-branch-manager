# Refactor Workflow

This document outlines the systematic approach for detecting duplicate code and redundant functionality in the codebase.

## Overview

The goal is to identify functions and methods that perform similar or identical operations to enable refactoring and code consolidation.

## Process

### 1. Job Organization

- Process files in directory-based batches ("jobs")
- Each job targets a specific directory or logical grouping of files
- Jobs are numbered sequentially (001, 002, etc.)

### 2. Analysis Steps

For each job:

1. **Scan Directory**: Identify all `.go` files in the target directory
2. **Extract Functions**: Parse each file to find:
   - Function declarations (`func functionName()`)
   - Method declarations (`func (receiver) methodName()`)
   - Struct methods and interface implementations
3. **Document Details**: For each function/method record:
   - Function/method name
   - File path (relative to repo root)
   - Line number where declared
   - Brief description of functionality
   - Input parameters and return types
   - Usage locations (where it's called from)
4. **Identify Patterns**: Look for:
   - Similar function names
   - Similar parameter signatures
   - Similar functionality descriptions
   - Potential duplicates or redundant code

### 3. Output Format

Results are appended to `docs/refactor_report.md` with the following structure:

```markdown
# <job_number>_<job_name>

## Directory: <target_directory>

### Functions Found

#### <function_name>
- **File**: `<file_path>:<line_number>`
- **Signature**: `<function_signature>`
- **Description**: <what_it_does>
- **Usage**: <where_it_is_called>

### Potential Duplicates/Redundancy

- List of functions that appear to have similar functionality
- Suggested consolidation opportunities

---
```

### 4. Job Checklist

Jobs to process (tracked in `docs/refactor_workflow_state.md`):

- [ ] **001_cmd_directory** - Analyze all files in `cmd/` directory
- [ ] **002_internal_core** - Analyze core internal packages (`config.go`, `git.go`, `manager.go`, `jira.go`)
- [ ] **003_internal_utils** - Analyze utility packages (`styles.go`, `table.go`, `info_renderer.go`)
- [ ] **004_internal_testutils** - Analyze test utilities in `internal/testutils/`
- [ ] **005_internal_mergeback** - Analyze mergeback functionality (`mergeback.go` and related tests)
- [ ] **006_main_entry** - Analyze main entry point (`main.go`)
- [ ] **007_cross_package** - Cross-package analysis for duplicate patterns

### 5. State Management Rules

**CRITICAL**: Follow these rules when updating `docs/refactor_workflow_state.md`:

1. **Never Remove Completed Jobs**: Once a job is marked as completed, it must remain in the "Completed Jobs" section
2. **Move Jobs Properly**: When completing a job:
   - Move it from "Planned Jobs" to "Completed Jobs" section
   - Add completion date: `(completed YYYY-MM-DD)`
   - Mark with `[x]` checkbox
3. **Update Current Job Section**: After completing a job:
   - Reset to "None" / "Ready for next job"
   - Clear started/completed dates
4. **Preserve Job History**: The completed jobs section serves as a historical record of progress

### 6. Analysis Tools

Use these approaches for comprehensive analysis:

1. **AST Parsing**: Parse Go files to extract function declarations
2. **Grep/Search**: Find function calls and usage patterns
3. **Signature Comparison**: Compare function signatures for similarities
4. **Code Similarity**: Identify functions with similar logic flow

