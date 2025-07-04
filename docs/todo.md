# TODO
- [x] option to log to file
- [x] binary can be called anywhere inside the git repo
- [x] add clone verb
    - [x] clones a bare repo
    - [x] creates the MAIN worktree using the HEAD branch
    - [x] detects if the worktree has a .envrc
        * YES - use that as the .envrc for the repo
        * NO - tell user to create a .envrc file. suggest to user to generate one based on the initial worktree we created.
- [x] suport `add` verb to add worktree
    - [x] normal worktree on new branch
    - [x] normal worktree on existing branch
    - [x] worktree based on JIRA ticket (use `jira-cli` to select ticket)
        - [x] configure your own rules for parsing jira ticket for branch name
    - [x] https://github.com/ankitpokhrel/jira-cli/discussions/356
- [x] enhance `switch` command (future features)
    - [x] auto-completion for worktree names
    - [x] automatic directory switching with shell integration
    - [x] fuzzy matching (e.g., `gbm switch prod` matches `PROD`)
- [x] add `push` verb
    - [x] `gbm push` pushes worktree if in a worktree - otherwise error
    - [x] `gbm push <worktree_name>` pushes named worktree - no matter what directoy you are in
    - [x] `gbm push --all` pushes all managed worktrees
- [x] add `pull` verb
    - [x] `gbm pull` pulls worktree if in a worktree - otherwise error
    - [x] `gbm pull <worktree_name>` pulls named worktree - no matter what directoy you are in
    - [x] `gbm pull --all` pulls all managed worktrees
- [x] use lipgloss for tables and styling
- [x] fix `gbm pull`
```sh
󰀵 jschneider  ~/code/scratch/integrator/worktrees/MAIN   master  󰟓 v1.24.4
  gbm pull
Pulling current worktree 'MAIN'...
There is no tracking information for the current branch.
Please specify which branch you want to merge with.
See git-pull(1) for details.

    git pull <remote> <branch>

If you wish to set tracking information for this branch you can do so with:

    git branch --set-upstream-to=origin/<branch> master

Error: exit status 1
Usage:
  gbm pull [worktree-name] [flags]

Flags:
      --all    Pull all worktrees
  -h, --help   help for pull

Global Flags:
  -c, --config string         specify custom .envrc path
  -d, --debug                 enable debug logging to ./gbm.log
  -w, --worktree-dir string   override worktree directory location

ERROR: Error: exit status 1
```
- [x] add configuration for controlling the icons for git status, repo validations, etc
    - [x] ./.gbm/config.toml
- [ ] track merge backs... not sure how yet and how to prompt the user.
    - [ ] helper to create a merge worktree
- [ ] configuration for copying files into new worktrees (.env, anything not tracked by git)
- [x] add completion support (cobra built-in bash/zsh/fish/powershell)
- [ ] add carapace completion
