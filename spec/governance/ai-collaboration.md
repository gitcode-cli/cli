# AI 协作规范

## 职责

定义 Codex、Claude、GitCode CLI skills、Loop Engineering 标准包和项目规则之间的关系。

## 适用场景

- 调整 `AGENTS.md` 或 `CLAUDE.md`
- 调整仓内历史 skill 兼容层
- 调整 `gitcode-cli/skills`
- 调整 `gitcode-cli/loop-kits`
- 设计多 AI 协作、Loop Engineering 或动态工作流边界

## 必须

- 以 `spec/` 作为项目正式规则源
- 以 `AGENTS.md` 和 `CLAUDE.md` 作为不同 AI 客户端进入规则体系的入口
- 以 `gitcode-cli/skills` 作为 AI skill 真相源
- 以 `gitcode-cli/loop-kits` 作为 Loop Engineering 标准包真相源
- 让仓内 `.ai/skills/`、`.claude/skills/`、`.codex/skills/` 只承担历史兼容/迁移参考职责
- 让作者主体与独立评审主体保持分离
- 让 GitCode issue / PR 保存长期状态和证据

## 禁止

- 在 AI 入口文档或 skills 中定义与 `spec/` 冲突的项目规则
- 把仓内旧 skill 目录当作当前真相源
- 让外部可复用 skill 依赖 gitcode-cli 仓库私有路径
- 让同一执行主体同时扮演作者和独立评审
- 让 GitHub mirror 状态替代 GitCode 主仓状态
- 把本地 `.loop/runtime/` 或聊天上下文当作长期事实源

## 同步要求

- AI 入口变化时同步 `AGENTS.md`、`CLAUDE.md`
- 项目规则变化时同步 `spec/`
- AI 执行方法变化时同步 `gitcode-cli/skills`
- Loop schema、policy、hook、template、adapter 变化时同步 `gitcode-cli/loop-kits`
- 外部用户说明变化时同步 `docs/`

## 权威关系

按以下顺序理解：

1. `spec/` 定义项目正式规则
2. `spec/governance/source-of-truth-matrix.md` 定义事实源边界
3. `AGENTS.md` / `CLAUDE.md` 定义不同 AI 客户端如何进入规则体系
4. `gitcode-cli/skills` 定义 AI 如何执行某类任务
5. `gitcode-cli/loop-kits` 定义可复用标准件
6. `.ai/skills/` / `.claude/skills/` / `.codex/skills/` 只保留历史兼容说明

## Dynamic Workflows 边界

Claude Dynamic Workflows 或其他 agent team 能力可用于大规模审计、迁移、CI 批量排障、跨文件分析和多角色评审。

这些并行执行能力不能越过：

- issue / PR 的 GitCode 事实状态
- `spec/` 定义的人工确认点
- GitHub mirror CI 的 commit SHA 证据绑定
- 作者与独立评审主体分离要求

## 下一步去哪里

- 改技能分发：查看 `gitcode-cli/skills`
- 改 Loop 标准包：查看 `gitcode-cli/loop-kits`
- 执行仓库内 AI 本地开发：继续看 [../workflows/ai-local-development-workflow.md](../workflows/ai-local-development-workflow.md)
