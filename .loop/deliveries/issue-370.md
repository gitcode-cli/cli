# Delivery: Issue #370

- **Issue**: [#370](https://gitcode.com/gitcode-cli/cli/issues/370) — docs: README 将 Docker 标记为规划中但已实现
- **PR**: [#303](https://gitcode.com/gitcode-cli/cli/-/merge_requests/303)
- **Date**: 2026-06-27
- **Type**: docs-only
- **Risk**: low

## Summary

README listed Docker as planned (`- [ ] Docker 镜像`) but Dockerfile, docker-compose.yml, and Makefile docker targets were already implemented. Moved Docker from planned to a proper installation section.

## Gates

| # | Gate | Result |
|---|------|--------|
| 1 | Implementation | ✅ Docker section added |
| 2 | Test | ⏭️ Skipped (docs-only) |
| 3 | Build | ⏭️ Skipped (docs-only) |
| 4 | Unit test | ⏭️ Skipped (docs-only) |
| 5 | Pre-commit | ✅ Passed |
| 6 | Command verify | ⏭️ Skipped (docs-only) |
| 7 | CI | ⏭️ Skipped (docs-only) |
| 8 | Risk classify | ✅ risk/low |

## Status

Merged via auto-merge (risk/low + docs-only, PR #303).
