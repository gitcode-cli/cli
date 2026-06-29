# Delivery: Issue #376

- **Issue**: [#376](https://gitcode.com/gitcode-cli/cli/issues/376) — docs: Docker 文档变更未按 docs-governance.md 同步到相关文档
- **PR**: [#311](https://gitcode.com/gitcode-cli/cli/pulls/311) — docs: add Docker references to COMMANDS.md, AGENTS.md, CLAUDE.md
- **Type**: docs-only
- **Risk**: low
- **Branch**: worktree-issue-376-docker-docs
- **Date**: 2026-06-29

## Changes

Per docs-governance.md §6.1, added Docker references to 3 entry documents:

| File | Change |
|------|--------|
| `docs/COMMANDS.md` | New Docker usage subsection under prerequisites |
| `AGENTS.md` | Docker install and usage entry in §5 common entry points |
| `CLAUDE.md` | Docker install and usage entry in §5 common entry points |

## Gates

| # | Gate | Status |
|---|------|--------|
| 1 | Implementation | ✅ |
| 2 | Tests | N/A (docs-only) |
| 3 | Build | N/A (docs-only) |
| 4 | UT | N/A (docs-only) |
| 5 | Pre-commit | ✅ 10/10 pass |
| 6 | Command verification | N/A (docs-only) |
| 7 | CI | N/A (docs-only) |
| 8 | Risk classification | risk/low |

## Verification

- Pre-commit: 10/10 pass
- `grep -i docker` confirms references in all 3 files
- Classifier: `risk=low` for all 3 files (documentation assets)
