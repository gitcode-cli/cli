# Delivery: Issue #371

- **Issue**: [#371](https://gitcode.com/gitcode-cli/cli/issues/371) — docs: issues-plan 目录残留 19 处 internal/config 旧路径引用
- **PR**: [#302](https://gitcode.com/gitcode-cli/cli/merge_requests/302)
- **Merge Commit**: `352aa39`
- **Date**: 2026-06-27

## Summary
Replaced 12 `internal/config/` → `pkg/config/` references and annotated 7 `internal/authflow/`/`internal/prompter/` references as historical design across 8 files in `issues-plan/`.

## Type
docs-only

## Risk
risk/low

## Gates
| # | Gate | Status |
|---|------|--------|
| 1 | 实现 | N/A |
| 2 | 测试 | N/A |
| 3 | 构建 | N/A |
| 4 | UT | N/A |
| 5 | Pre-commit | ✅ |
| 6 | 命令验证 | N/A |
| 7 | CI | N/A |
| 8 | 风险分级 | risk/low |
