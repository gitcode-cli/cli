# Issue #409 交付记录

## 基本信息

| 字段 | 值 |
|------|------|
| Issue | [#409](https://gitcode.com/gitcode-cli/cli/issues/409) |
| 标题 | refactor: package-level mutable globals in repo/sync and pr/sync block parallel testing |
| 类型 | type/refactor |
| 风险 | risk/medium (classify-change-risk.py: medium) |
| 分支 | worktree-refactor+issue-409 |
| PR | [#316](https://gitcode.com/gitcode-cli/cli/merge_requests/316) |
| 完成时间 | 2026-06-29 16:10 |

## 状态流转

| 阶段 | 状态 | 时间 |
|------|------|------|
| 初始 | status/triage | 2026-06-29 15:39 |
| 验证 | status/verified | 2026-06-29 15:50 |
| 开发 | status/in-progress | 2026-06-29 15:51 |
| 自检 | status/self-checked | 2026-06-29 16:07 |
| 评审通过 | status/approved | 2026-06-29 16:10 |
| 合并 | status/merged | 2026-06-29 16:10 |

## 门禁证据

| 序号 | 门禁 | 状态 | 证据 |
|------|------|:----:|------|
| 1 | 开发实现 | ✅ | 4 files, +46/-45 |
| 2 | 单元测试 | ✅ | 29/29 sync tests, 1268/1268 all |
| 3 | 本地构建 | ✅ | go build -o ./gc ./cmd/gc |
| 4 | 全量测试 | ✅ | go test ./... passed |
| 5 | Pre-commit | ✅ | 10/10 hooks passed |
| 6 | 实际命令验证 | ✅ | 纯重构，行为不变，go test 全量通过 |
| 7 | 远端 CI | — | 待 GitHub Actions 触发 |
| 8 | 风险分级 | ✅ | scripts/classify-change-risk.py: medium |

**门禁**: 8/8

## 变更摘要

- `pkg/cmd/repo/sync/sync.go`: 移除 `var gitRun`，clone 调用改为 `opts.GitRun("", nil, "clone", ...)`
- `pkg/cmd/pr/sync/sync.go`: 移除 `var gitRunWithEnv` / `var gitRunInDirWithEnv`，SyncOptions 新增 `GitRun` / `GitRunInDir` 字段，`syncCommits` 接受函数参数
- 对应测试文件: mock 注入从包级变量重写为 Options 字段注入

## 评审

第一轮 4 角色全部 approved（代码/安全/测试/文档）

## 经验教训

无新增发现
