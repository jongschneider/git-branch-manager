- Fix hardcoded merge-back chain in hotfix command message
    * Currently shows: "Remember to merge back through the deployment chain: master ’ preview ’ main"
    * Should dynamically determine the correct deployment chain from gbm.branchconfig.yaml
    * The actual chain should be based on the merge_into configuration in the branch config