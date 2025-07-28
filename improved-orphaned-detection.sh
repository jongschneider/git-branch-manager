#!/bin/bash
# Improved orphaned task detection script
# Handles directories with special characters, command substitution, etc.

mkdir -p todos/work todos/done

orphaned_count=0
orphaned_tasks=()

# Use find to get all task.md files, avoiding glob issues
while IFS= read -r -d $'\0' task_file; do
    # Skip if file doesn't exist (shouldn't happen with find, but safety check)
    [ -f "$task_file" ] || continue
    
    # Extract PID with better error handling
    pid=$(grep "^**Agent PID:" "$task_file" 2>/dev/null | cut -d' ' -f3 | tr -d ' ')
    
    # Skip if PID is empty or process is still running
    if [ -n "$pid" ] && ps -p "$pid" >/dev/null 2>&1; then
        continue
    fi
    
    # This is an orphaned task
    orphaned_count=$((orphaned_count + 1))
    
    # Get directory name and task title safely
    task_dir=$(dirname "$task_file")
    task_name=$(basename "$task_dir")
    task_title=$(head -1 "$task_file" 2>/dev/null | sed 's/^# //' | tr -d '\n\r')
    
    # Store for output
    orphaned_tasks+=("$orphaned_count. $task_name: $task_title")
    
done < <(find todos/work -name 'task.md' -type f -print0)

# Output results
for task in "${orphaned_tasks[@]}"; do
    echo "$task"
done