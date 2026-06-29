---
name: gitcode-commit
description: Use GitCode CLI commit commands to view commits, retrieve diff or patch output, and manage commit comments. Trigger for GitCode commit inspection, commit review, or commit comment workflows.
---

# gitcode-commit

Inspect commits and manage commit comments.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Inspect Commits

```bash
gitcode commit view <sha> -R owner/repo
gitcode commit view <sha> -R owner/repo --show-diff
gitcode commit view <sha> -R owner/repo --json
gitcode commit view <sha> -R owner/repo --web

gitcode commit diff <sha> -R owner/repo
gitcode commit patch <sha> -R owner/repo
```

## Comments

```bash
gitcode commit comments create <sha> -R owner/repo --body "Comment text"
gitcode commit comments list -R owner/repo --json
gitcode commit comments list -R owner/repo --page 1 --per-page 50 --json
gitcode commit comments list-by-sha <sha> -R owner/repo --json
gitcode commit comments view <id> -R owner/repo --json
gitcode commit comments edit <id> -R owner/repo --body "Updated comment"
```

## Rules

- Prefer `--json` for list and view commands when parsing output.
- Use exact commit SHA values; avoid ambiguous short SHAs in automation.
- When reviewing a commit, inspect both metadata and diff/patch.
- Do not paste sensitive data into commit comments.
