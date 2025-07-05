> how about a different approach. can we use carapace to create the dynamic completions instead? see
  docs/libs/carapace_docs/ for documentation if you need it.

> no. we are off. I want the prompt to be updated in realtime. like other shell completions. not just output
  to stdout.

> looks like you need to give the <worktree-name> when you use the -j flag. the user might not know what to
  call the worktree before checking the jira tickets. could we make it just `gbm add -j` when using that
  flag? and then dynamically create the worktree names after the jira ticket number (KEY)?

> read docs/spec.md docs/todo.md and docs/libs/jira-cli_readme.md. i think
  what we want is
  ```sh
  jira add -j
  ```
  that or (jira add --jira) should execute `jira issue list -a$(jira me) --plain` and that create a search with a select with the options being the output from the `jira issue list -a$(jira me) --plain`. using the search should filter rows out.




> how about a different approach. can we use carapace to create the dynamic completions instead? see
  docs/libs/carapace_docs/ for documentation if you need it.

> no. we are off. I want the prompt to be updated in realtime. like other shell completions. not just output
  to stdout.

> looks like you need to give the <worktree-name> when you use the -j flag. the user might not know what to
  call the worktree before checking the jira tickets. could we make it just `gbm add -j` when using that
  flag? and then dynamically create the worktree names after the jira ticket number (KEY)?

> read docs/spec.md docs/todo.md, docs/libs/charmbracelet_huh_readme.md and docs/libs/jira-cli_readme.md. i think
  what we want is
  ```sh
  jira add -j
  ```
  that or (jira add --jira) should execute `jira issue list -a$(jira me) --plain` and that create a search with a
  select with the options being the output from the `jira issue list -a$(jira me) --plain`. using the search should
  filter rows out.

