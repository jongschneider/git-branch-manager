# Fix hardcoded merge-back chain in hotfix command message
**Status:** Done
**Agent PID:** 50911

## Original Todo
- Fix hardcoded merge-back chain in hotfix command message
    * Currently shows: "Remember to merge back through the deployment chain: master → preview → main"
    * Should dynamically determine the correct deployment chain from gbm.branchconfig.yaml
    * The actual chain should be based on the merge_into configuration in the branch config

## Description
The deployment chain message in the hotfix command is showing the merge direction backwards. Currently it shows "main → preview → production" but should show "production → preview → main" to indicate the proper merge-back flow direction.

## Implementation Plan
- [x] Fix the `buildMergeChain` function in `cmd/hotfix.go` to build the chain in the correct direction (forward merge flow)
- [x] Update the logic to follow `merge_into` relationships forward instead of backward
- [x] Create test to verify the correct deployment chain direction
- [x] Test with sample configuration to ensure "production → preview → main" flow

## Notes
[Implementation notes]