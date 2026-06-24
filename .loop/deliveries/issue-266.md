# Delivery Record: Issue #266
- **Title**: regression-core should refuse non infra-test repositories
- **Type**: bug
- **Status**: merged
- **PR**: [#273](https://gitcode.com/gitcode-cli/cli/pulls/273)
- **Date**: 2026-06-24

## Gate Compliance
| # | Gate | Result | Evidence |
|---|------|--------|----------|
| 1 | 验证 | ✅ | GC_REGRESSION_REPO=owner/repo → rejected with exit 1 |
| 2 | 开发 | ✅ | scripts/regression-core.sh +14 |
| 3 | 构建 | skipped | bash script only |
| 4 | UT | ✅ | 手动验证: 拒绝+接受 |
| 5 | Pre-commit | ✅ | |
| 6-7 | — | skipped | |
| 8 | 风险 | ✅ | risk/medium |
| + | 合并 | ✅ | merged |
