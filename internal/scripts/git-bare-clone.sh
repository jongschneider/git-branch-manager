#!/bin/bash
# Make directory for your git repository
repo=$(basename $1 .git)
mkdir "$repo"
pushd "$repo"
echo "Cloning bare repository to .git..."
git clone --bare $1 .git
pushd '.git' > /dev/null
echo "Adjusting origin fetch locations..."
# Explicitly sets the remote origin fetch so we can fetch remote branches
git config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"
# Fetch all branches from remote
echo "Fetching all branches from remote..."
git fetch origin
popd > /dev/null
echo "Success."

