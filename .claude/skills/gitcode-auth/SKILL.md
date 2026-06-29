---
name: gitcode-auth
description: Use GitCode CLI authentication commands to log in, inspect auth state, troubleshoot token source and protocol settings, or log out. Trigger for GitCode auth, token, login, logout, permission, or credential problems.
---

# gitcode-auth

Use GitCode CLI authentication safely across Windows and Linux.

## Command Entry

Use `gitcode` in reusable instructions. On Linux/macOS, `gc` is equivalent. On Windows PowerShell, `gc` is `Get-Content`, so prefer `gitcode`.

Verify the installed CLI first:

```bash
gitcode version
gitcode help --json
```

## Workflow

1. Inspect the current state before changing credentials.
2. Prefer environment variables for CI and temporary sessions.
3. Use `--with-token` for non-interactive login.
4. Never print or paste full tokens unless the user explicitly asks and the output location is safe.

## Commands

```bash
# Status and token source
gitcode auth status
gitcode auth status --json
gitcode auth status --hostname gitcode.com

# Token for scripts
gitcode auth token
gitcode auth token --json

# Environment variable auth
export GC_TOKEN="your_gitcode_token"
export GITCODE_TOKEN="your_gitcode_token"

# Non-interactive login
echo "YOUR_TOKEN" | gitcode auth login --with-token

# Select Git protocol during login
echo "YOUR_TOKEN" | gitcode auth login --with-token --git-protocol ssh

# Logout
gitcode auth logout
gitcode auth logout --yes
```

PowerShell examples:

```powershell
$env:GC_TOKEN = "your_gitcode_token"
gitcode auth status
Get-Content token.txt | gitcode auth login --with-token
```

## Rules

- `GC_TOKEN` has priority over `GITCODE_TOKEN`; environment variables override stored login.
- `auth logout` removes stored credentials but cannot unset environment variables.
- For code download and sync, prefer SSH and verify `ssh -T git@gitcode.com`.
- For automation, prefer JSON output and avoid commands that require an interactive TTY.
- Redact tokens in summaries, logs, PR comments, and issue comments.
