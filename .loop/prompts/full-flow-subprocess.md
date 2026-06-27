从 status/triage 取一个 issue，推进到 merged。每次只处理一个。

## 前置
- 独立 git worktree（`.claude/worktrees/issue-N-<ts>`），用后即删
- 禁止在 main 开发、跳过验证、作者自检当独立评审

## 流程
1. 取 issue（`gc issue list -R gitcode-cli/cli --state opened --label status/triage --limit 5`），选最小 scope；若 triage 空→孤儿 PR 检查
2. 判定 docs-only 还是代码改动，走对应分支
3. 状态机: triage→verified→in-progress→draft→self-checked→ready→approved→merged
   risk/low 自动合，risk/high 暂停

## 门禁
| # | 门禁 | docs-only | 代码改动 |
|---|------|:--:|------|
| 1 | 实现 | — | 修复 |
| 2 | 测试 | 跳过 | go test ./... 全通过 |
| 3 | 构建 | 跳过 | go build 成功 |
| 4 | UT | 跳过 | 全通过 |
| 5 | Pre-commit | 必须 | 必须 |
| 6 | 命令验证 | 跳过 | infra-test/* 至少一条 |
| 7 | CI | 跳过 | gh CLI 触发 GitHub Actions，等待全绿，PR 附 run URL |
| 8 | 风险分级 | 必须 | classify-change-risk.py |

## 证据
- Issue: 验证记录 + 自检 9 项（含 CI URL）
- PR: 评审结论 + CI URL + gate 表
- docs-only 跳过评审；其余路径必须多角色独立评审
- CI 未跑写 ✅ 算违规

## 交付
创建 `.loop/deliveries/issue-N.md`，更新 README。末尾输出 `ISSUE_NUM=<N>`。

## 孤儿 PR（仅 triage 为空时）
`gc pr list --state open --json`，找本人非 draft PR→完整读评论→对照 spec/workflows/development-workflow.md §5.3 补缺失→合并。
