# Loop Engineering 总规范

## 目标

gitcode-cli 的 Loop Engineering Demo v1 用于展示：

- AI 不依赖聊天记忆推进流程
- 状态推进必须有证据
- 事实源和执行源分离
- 作者、自检、独立评审和合并确认分层
- 可复用知识有明确归档位置

## 组成

| 层 | 位置 | 职责 |
| --- | --- | --- |
| 规则层 | `gitcode-cli/cli/spec/loop` | 状态机、证据、归档、人工确认点 |
| 执行层 | `gitcode-cli/skills` | AI 如何执行某类 loop |
| 标准层 | `gitcode-cli/loop-kits` | schema、policy、hooks、templates、adapters |
| 协作事实 | GitCode issue / PR | 长期状态、评论、评审、证据 |
| CI 事实 | GitHub mirror Actions | commit SHA 绑定的 CI 执行结果 |

## 必须

- 每个长期任务必须能从 GitCode issue / PR 恢复上下文
- 每次状态推进必须有可追溯证据
- CI 证据必须绑定 commit SHA
- 独立评审不能由作者自检替代
- 合并完成必须以 GitCode merged PR 和 `origin/main` 为准
- 归档前必须判断知识资产落点

## 禁止

- 把本地 `.loop/runtime` 当作长期事实源
- 把 GitHub mirror 当作开发主仓或合并事实源
- 把一次性执行证据提交进主仓文档
- 让 skill 定义低于 `spec/` 的项目硬规则
- 在 Phase 1-3 中实现自动驾驶式 `gc loop run`
