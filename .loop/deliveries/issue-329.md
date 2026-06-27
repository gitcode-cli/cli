# Issue #329 — refactor: getBody 函数在 issue/create 和 pr/create 中完全重复

## 状态流转

| 阶段 | 时间 | 状态标签 |
|------|------|---------|
| Triage | 2026-06-26 | `status/triage` `type/refactor` `scope/cmdutil` `risk/medium` |
| Verified | 2026-06-27 10:31 | `status/verified` |
| In Progress | 2026-06-27 10:31 | `status/in-progress` |
| PR Created | 2026-06-27 10:37 | PR #296 `status/draft` |
| Self Check | 2026-06-27 10:38 | 自检完成 |
| Multi-Role Review | 2026-06-27 10:39 | 4 角色评审全部通过 |
| Merged | 2026-06-27 10:40 | `status/merged` |

## 8 Gate 完成表

| # | 门禁 | 状态 | 证据 |
|---|------|:--:|------|
| 1 | 开发实现 | ✅ | `cmdutil.ReadBody` 替换两处重复 `getBody` |
| 2 | 测试 | ✅ | 新增 6 个 ReadBody 单元测试（空输入、body设置、冲突、stdin、文件、错误路径） |
| 3 | 本地构建 | ✅ | `go build -o ./gc ./cmd/gc` 成功 |
| 4 | 单元测试 | ✅ | 1250 测试全部通过 |
| 5 | Pre-commit | ✅ | 10 hooks 全部通过 |
| 6 | 实际命令验证 | ✅ | infra-test/gctest1: `--body`, `--body-file`, 冲突检测均正确 |
| 7 | 远端 CI | ✅ | https://github.com/gitcode-cli/cli/actions/runs/28276023093 — Test/Lint/Build(ubuntu) ✅ |
| 8 | 风险分级 | ✅ | risk/medium — `classify-change-risk.py` 确认 |

## 多角色评审

| 角色 | 判定 | 备注 |
|------|:--:|------|
| 代码审查 | ✅ Approve | P2: fmt.Errorf → errors.New 优化建议 |
| 安全审查 | ✅ Approve | 无新增安全风险 |
| 测试审查 | ✅ Approve | 6 测试覆盖全部分支，预存缺口已记录 |
| 文档审查 | ✅ Approve | 纯内部重构，无需文档更新 |

## PR / CI 证据

- **PR**: [#296](https://gitcode.com/gitcode-cli/cli/pulls/296)
- **CI**: https://github.com/gitcode-cli/cli/actions/runs/28276023093
- **验证记录**: https://gitcode.com/gitcode-cli/cli/issues/329#comment-177415965
- **自检记录**: https://gitcode.com/gitcode-cli/cli/issues/329#comment-177417601
- **评审汇总**: https://gitcode.com/gitcode-cli/cli/pulls/296#comment-3b96013f
- **合入确认**: https://gitcode.com/gitcode-cli/cli/issues/329#comment-177418221

## 变更摘要

| 文件 | 变更 |
|------|------|
| `pkg/cmdutil/text_file.go` | +32 行: 新增 `ReadBody()` 函数 |
| `pkg/cmdutil/text_file_test.go` | +62 行: 新增 6 个单元测试 |
| `pkg/cmd/issue/create/create.go` | -26/+2 行: getBody → cmdutil.ReadBody |
| `pkg/cmd/pr/create/create.go` | -26/+2 行: getBody → cmdutil.ReadBody |
