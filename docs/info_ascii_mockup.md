$ gbm info PROJECT-123

╔═══════════════════════════════════════════════════════════════════════════════╗
║                           📋 WORKTREE INFO: PROJECT-123                       ║
╚═══════════════════════════════════════════════════════════════════════════════╝

┌─ 📁 WORKTREE ───────────────────────────────────────────────────────────────────┐
│ Name: PROJECT-123                                                               │
│ Path: /home/user/repos/myapp/worktrees/PROJECT-123                              │
│ Branch: feature/PROJECT-123_update_user_auth_flow                               │
│ Created: 2025-07-02 14:30:15 (3 days ago)                                       │
│ Status: 🟡 DIRTY (5 files modified                                              │
│ PR: https://github.com/company/myapp/pull/456 (Draft)                           │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─ 🎫 JIRA TICKET ────────────────────────────────────────────────────────────────┐
│ Key: PROJECT-123                                                                │
│ Summary: Update user auth flow                                                  │
│ Status: 🔄 In Progress → Code Review                                            │
│ Assignee: john.doe@company.com                                                  │
│ Priority: 🔴 High                    Reporter: jane.smith@company.com           │
│ Created: 2025-06-28 09:15:00         Due Date: 2025-07-10 17:00:00              │
│ Epic: AUTH-001 (User Authentication Overhaul)                                   │
│ Link: https://company.atlassian.net/browse/PROJECT-123                          │
│                                                                                 │
│ 💬 Latest Comment (2 hours ago):                                                │
│    "Please review the password validation logic before merging"                 │
│    - tech.lead@company.com                                                      │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─ 🌿 GIT STATUS ────────────────────────────────────────────────────────────────┐
│ Base Branch: develop (diverged 5 days ago at commit abc1234)                   │
│ Upstream: origin/feature/PROJECT-123_update_user_auth_flow                     │
│ Position: ↑ 4 commits ahead, ↓ 2 commits behind origin                         │
│                                                                                │
│ Last Commit: feat: implement 2FA validation (def5678) - 3 hours ago            │
│ Author: John Doe <john.doe@company.com>                                        │
│                                                                                │
│ Modified Files:                                                                │
│   M  src/auth/validator.go        (+89, -23)                                   │
│   M  src/auth/middleware.go       (+45, -12)                                   │
│   A  src/auth/twofa.go            (+156, -0)                                   │
│   M  tests/auth_test.go           (+78, -5)                                    │
│   M  docs/auth_flow.md            (+12, -3)                                    │
│                                                                                │
│ Recent Commits:                                                                │
│   def5678 feat: implement 2FA validation               (3 hours ago)           │
│   ghi9012 refactor: extract auth helpers               (1 day ago)             │
│   jkl3456 fix: handle edge case in password reset      (2 days ago)            │
│   mno7890 feat: add session timeout configuration      (3 days ago)            │
└────────────────────────────────────────────────────────────────────────────────┘

