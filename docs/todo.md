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
- [x] add completion support (cobra built-in bash/zsh/fish/powershell)
- [x] add configuration for controlling the icons for git status, repo validations, etc
    - [x] ./.gbm/config.toml
- [x] track all worktrees created with `gbm` in the `list` and `status` commands
- [x]  add `remove` verb
- [x] support `gbm switch -` to go to previous worktree
- [x] sort branches by .envrc first, then worktree createdAt DESC
- [x] add info verb (see info_prd.md and info_ascii_mockup.md)
- [x] jira-cli support
- [x] combine list and status. they do the same thing.
    - [x] use `list` and remove `status`
    - [x] columns should be WORKTREE | BRANCH | GIT STATUS | SYNC STATUS | PATH (if not enough room in terminal, omit PATH)
- [x] make output adaptive layout.
    - [x] responsive table design for gbm list - hides PATH column when terminal is narrow (< 100 chars)
- [x] remove the `clean` verb
- [x] make .gbm.config.yaml tests actually use yaml for validation instead of string contains.
- [x] what's the point of TestCloneCommand_EmptyRepository?
- [x] configuration for copying files into new worktrees (.env, anything not tracked by git)
- [ ] track merge backs... not sure how yet and how to prompt the user.
    - [ ] helper to create a merge worktree
- [ ] add carapace completion
- [ ] control branch name JIRA through config file (use current as defaults)
- [x] add jira me to config.toml
- [ ] replace confirmation with bubbltea confirmation (lipgloss?)
- [ ] add `theme` verb with default themes
- [ ] support `JIRA_CONFIG_FILE` env var (see https://github.com/ankitpokhrel/jira-cli?tab=readme-ov-file#multiple-projects)
- [ ] add hooks (for automating tasked before/after a command is run)
    * example: after `gbm add` copy `.env` file from MAIN
- [ ] review and improve fuzzy matching logic in switch command
    * current behavior: input is normalized to uppercase, then fuzzy matched against actual config keys
    * this creates inconsistent behavior between direct matches (uppercase) and fuzzy matches (config case)
    * consider making the matching logic more intuitive and consistent
- [x] fix .gbm.config.yaml creation (clone.go) and make default branch not be MAIN, but what the default branch is called
    * don't enforce CAPS for worktree name
- [x] add confirmation to the `gbm sync --force` because it is destructive.
    * give list of what will be destroyed
- [x] make sure `gbm add xxxxxx` when performed in a worktree will create the new worktree in ./worktree/xxxx and not in the current worktree dir (like native wt add)
- [x] add a base branch argument to `gbm add` so it acts like `git worktree add` command
    * if no base branch is supplied, use the default branch.
- [ ] add `gbm hotfix` or `gbm hf`
- [ ] add `gbm mergeback` or `gbm mb`
- [ ] ai plugin for:
    * merge conflict resolution
    * commit messages
- [ ] consider a timestamp based approach for checking drift
    * .gbm/drift_tracker_timestamp
- [ ] make --fetch the default for `gbm  sync`
  Summary:
  - gbm sync = "Sync worktrees with the branch versions I currently have locally"
  - gbm sync --fetch = "First get the latest from remote, then sync worktrees with those updated versions"
  The --fetch flag ensures you're always working with the most current version of your branches, which is
  especially important in collaborative environments where branches are frequently updated.


