package internal

func (gm *GitManager) RemoveWorktree(worktreePath string) error {
	if err := execGitCommandRun(gm.repoPath, "worktree", "remove", worktreePath, "--force"); err != nil {
		return enhanceGitError(err, "worktree remove")
	}

	return nil
}
