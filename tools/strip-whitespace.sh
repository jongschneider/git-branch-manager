#!/bin/bash

# Strip trailing whitespace from a file
# Usage: strip-whitespace.sh <file>

set -euo pipefail

if [ $# -ne 1 ]; then
    echo "Usage: $0 <file>" >&2
    exit 1
fi

file="$1"

if [ ! -f "$file" ]; then
    echo "Error: File '$file' not found" >&2
    exit 1
fi

# Strip trailing whitespace and tabs from each line
# Handle macOS vs Linux sed differences
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS requires empty string after -i
    sed -i '' 's/[[:space:]]*$//' "$file"
else
    # Linux sed
    sed -i 's/[[:space:]]*$//' "$file"
fi

echo "Stripped trailing whitespace from: $file"
