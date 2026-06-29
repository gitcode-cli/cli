---
name: gitcode-pr-review
description: Conduct an engineering review of a GitCode PR using GitCode CLI diff, comments, review, approval, and local verification. Trigger when users request PR review, independent review, approval readiness, or merge risk analysis.
---

# gitcode-pr-review

Review a GitCode PR as an engineering reviewer.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Review Stance

Lead with bugs, regressions, security risks, missing tests, and release blockers. Keep style-only comments secondary.

## Workflow

1. Read PR metadata and comments.
2. Inspect diff.
3. Checkout locally if deeper validation is needed.
4. Run relevant tests or explain what could not be run.
5. Post findings or approval.

## Commands

```bash
gitcode pr view 123 -R owner/repo --comments --json
gitcode pr comments 123 -R owner/repo --json
gitcode pr diff 123 -R owner/repo
gitcode pr checkout 123 -R owner/repo
```

Post findings:

```bash
gitcode pr comment 123 -R owner/repo --body-file review-findings.md --json
gitcode pr comment 123 -R owner/repo --path path/to/file.go --position 12 --body-file inline.md --json
gitcode pr review 123 -R owner/repo --comment-file review-report.md
```

Approve:

```bash
gitcode pr review 123 -R owner/repo --approve --comment-file approval.md
```

## Report Format

```markdown
## Findings

- [High] path/file:line - issue, impact, recommended fix.

## Open Questions

- ...

## Verification

- ...

## Summary

- ...
```

## Rules

- If no issues are found, say so clearly and mention residual test risk.
- Do not approve unless permission and confidence are both present.
- Avoid leaking secrets from diffs into comments.
- If a requested action is not supported by the GitCode API/CLI, state the limitation.
