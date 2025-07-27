# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Current State

This repository is in its initial state with minimal structure:
- `docs/spec.md` - Empty specification file

## Architecture

The repository appears to be intended for a worktree manager tool based on the name, but no implementation exists yet.

## Development Setup

No build system, dependencies, or development tooling has been configured yet. When development begins, standard practices should be established for the chosen technology stack.

## Testing

### Installing and Testing gbm

To test the latest version of gbm (git branch manager):

1. **Install latest version**: Run `goi` (which is `go install ./...`) in this repo
2. **Testing environment**: Use tmux session named "test" for testing gbm functionality
3. **Availability**: After installation with `goi`, gbm will be available in the tmux session for testing

This allows for quick iteration and testing of changes to the gbm tool.
