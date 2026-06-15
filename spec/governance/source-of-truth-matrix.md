# 真相源矩阵

本文档定义 gitcode-cli 仓库中不同信息类型的真相源、适用边界和判定优先级。

## 职责

- 说明哪些来源可以直接作为事实依据
- 避免把阶段说明、聊天记录、旧 skill 目录或镜像仓状态误当作事实源
- 为 Loop Engineering 的状态推进、证据记录和归档提供判定基础

## 真相源矩阵

| 信息类型 | 真相源 | 可否直接判事实 | 说明 |
| --- | --- | --- | --- |
| 项目正式规则 | `spec/` | 是 | 最高规则源；入口文档、skills、docs 不得覆盖 |
| 命令行为 | `docs/COMMANDS.md` | 是 | `gc` 命令行为唯一真相源 |
| 用户使用说明 | `docs/` | 有条件地是 | 只对用户文档范围生效，不定义内部流程 |
| GitCode 协作状态 | GitCode issue / PR / labels / comments | 是 | issue、PR、评审讨论和状态推进事实源 |
| 主干完成事实 | GitCode merged PR + `origin/main` | 是 | 判断功能是否合入主干必须同时看远端主干事实 |
| CI 执行结果 | GitHub mirror Actions run | 是 | 只作为 CI 执行事实源，必须绑定 commit SHA |
| CI 与 GitCode 关联 | commit SHA | 是 | GitCode PR head SHA 与 GitHub Actions run 的绑定键 |
| AI skill | `gitcode-cli/skills` | 有条件地是 | 定义 AI 执行方法，不定义低于 `spec/` 的项目规则 |
| Loop 标准件 | `gitcode-cli/loop-kits` | 有条件地是 | 定义 schema、policy、hooks、templates、adapters 的通用契约 |
| 仓内旧 skill 目录 | `.ai/skills/`、`.claude/skills/`、`.codex/skills/` | 否 | 仅作为历史兼容/迁移参考 |
| 阶段说明 | `issues-plan/PROGRESS.md` | 否 | 可作为背景，不作为单个 issue / PR 实时状态事实源 |
| 本地 runtime | `.loop/runtime/` | 否 | 只作为本地临时缓存，不提交、不判事实 |
| 会话记忆 | `MEMORY.md` 或聊天上下文 | 否 | 只作恢复辅助；与 `spec/` 或远端事实冲突时以后者为准 |

## 判定优先级

1. `spec/` 定义项目规则
2. GitCode issue / PR / labels / comments 定义协作状态
3. GitCode merged PR + `origin/main` 定义主干完成事实
4. GitHub mirror Actions run 定义 CI 执行事实
5. `docs/COMMANDS.md` 定义命令行为
6. `gitcode-cli/skills` 定义 AI 执行方法
7. `gitcode-cli/loop-kits` 定义可复用标准件
8. 本地文件、阶段说明、会话记忆只作辅助

## Loop Engineering 特别规则

- 状态长期事实必须写回 GitCode issue / PR。
- CI 证据必须引用 GitHub Actions run URL 和 commit SHA。
- GitHub mirror 不能替代 GitCode 主仓，也不能作为合并事实源。
- `loop-kits` 的 schema 和 templates 不保存项目运行状态。
- 可复用规则进入 `spec/`，可复用 AI 方法进入 `gitcode-cli/skills`，可复用标准件进入 `gitcode-cli/loop-kits`。

## 下一步去哪里

- 改 AI 协作边界：继续看 [ai-collaboration.md](./ai-collaboration.md)
- 改文档分层：继续看 [docs-governance.md](./docs-governance.md)
- 改 Loop Engineering：继续看 [../loop/README.md](../loop/README.md)
