---
name: gitcode-review
description: Review GitCode pull requests and commits using GitCode CLI comments, review approval, replies, PR diffs, and commit comments. Trigger when the user asks to review, approve, comment on, or inspect a GitCode PR or commit.
---

# gitcode-review

Use GitCode CLI for PR and commit review.

## Command Entry

Use `gitcode` for cross-platform instructions. Linux/macOS may use `gc`; Windows PowerShell should use `gitcode`.

## Review Workflow

1. Read PR metadata and comments.
2. Inspect diff and changed files.
3. Run local checks if the repository is available.
4. Leave precise comments for actionable findings.
5. Approve only when the change is safe and policy allows approval.

## Inspect PR

```bash
gitcode pr view 123 -R owner/repo --json
gitcode pr view 123 -R owner/repo --comments --json
gitcode pr comments 123 -R owner/repo --json
gitcode pr diff 123 -R owner/repo
```

## Comment and Reply

```bash
gitcode pr comment 123 -R owner/repo --body "Short comment" --json
gitcode pr comment 123 -R owner/repo --body-file review-note.md --json
echo "Long comment" | gitcode pr comment 123 -R owner/repo --body-file - --json

gitcode pr reply 123 -R owner/repo --discussion <discussion_id> --body "Reply"
```

Inline comment:

```bash
gitcode pr diff 123 -R owner/repo
gitcode pr comment 123 -R owner/repo --path path/to/file.go --position 12 --body-file inline.md --json
```

## Submit Review or Approval

```bash
gitcode pr review 123 -R owner/repo --comment-file review-report.md
gitcode pr review 123 -R owner/repo --approve --comment "LGTM"
gitcode pr review 123 -R owner/repo --approve --comment-file approval.md
```

`--approve` requires platform approval permission, which is separate from merge permission.

## Commit Review

```bash
gitcode commit view <sha> -R owner/repo --show-diff --json
gitcode commit diff <sha> -R owner/repo
gitcode commit patch <sha> -R owner/repo

gitcode commit comments list -R owner/repo --json
gitcode commit comments list-by-sha <sha> -R owner/repo --json
gitcode commit comments create <sha> -R owner/repo --body "Comment"
gitcode commit comments view <id> -R owner/repo --json
gitcode commit comments edit <id> -R owner/repo --body "Updated"
```

## Review Output

Lead with findings, ordered by severity:

```markdown
## Findings

- [High] path/file.go:42 - Problem and impact.
- [Medium] path/file.go:87 - Problem and impact.

## Open Questions

- ...

## Summary

Short change summary and verification notes.
```

## Rules

- Do not approve a PR you have not inspected.
- Distinguish blockers from suggestions.
- Avoid posting secrets or full tokens in review comments.
- If CLI lacks a platform capability, state the limitation instead of inventing a workflow.
