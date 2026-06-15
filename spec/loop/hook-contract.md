# Loop Hook 契约

## 定位

Loop Hooks 不是 Git hooks。它们是 Loop Engineering 阶段门禁，可由 AI skill、未来的 `gc loop` 命令或外部 runner 调用。

## 标准 hook

| Hook | 目标 |
| --- | --- |
| `pre-loop` | 检查 auth、仓库、policy、issue 状态 |
| `pre-change` | 检查非 main、issue verified、范围清楚 |
| `post-change` | 收集 diff、测试、命令验证、文档同步 |
| `pre-pr` | 检查 self-check、安全项、关联 issue |
| `post-pr` | 等待 mirror SHA、查询 Actions、生成 CI evidence |
| `pre-merge` | 检查 review、CI、风险和人工批准 |
| `post-merge` | 确认 `origin/main` 包含代码 |
| `archive` | 判断知识资产归档目标 |

## 输出状态

- `pass`
- `fail`
- `blocked`
- `needs_human`

## Phase 1-3 边界

Phase 1-3 只定义 hook 契约和占位，不实现复杂自动化。具体可复用 hook 资产落在 `gitcode-cli/loop-kits`。
