# 文档治理规范

本文件定义 gitcode-cli 仓库的文档体系、AI 协作文档分层、唯一真相源和同步规则。

## 职责

- 定义文档分层
- 明确每类信息的唯一真相源
- 约束入口文档、用户文档、spec、skills 和 loop-kits 的边界
- 降低重复和漂移

## 文档分层

| 层级 | 位置 | 职责 |
| --- | --- | --- |
| 项目规则 | `spec/` | 内部开发、流程、质量、交付和 Loop Engineering 正式规则 |
| 用户文档 | `docs/` | 面向用户解释如何使用 `gc` |
| 命令手册 | `docs/COMMANDS.md` | 命令行为唯一真相源 |
| AI 入口 | `AGENTS.md`、`CLAUDE.md` | 不同 AI 客户端的项目入口 |
| AI 执行 skill | `gitcode-cli/skills` | 可复用 AI 执行方法 |
| Loop 标准包 | `gitcode-cli/loop-kits` | schema、policy、hooks、templates、adapters |
| 历史兼容层 | `.ai/skills/`、`.claude/skills/`、`.codex/skills/` | 旧仓内 skill 迁移参考 |
| 阶段说明 | `issues-plan/PROGRESS.md` | 背景说明，不判定单个 issue / PR 实时状态 |

## 唯一真相源

- 项目正式规则：`spec/`
- 命令行为：`docs/COMMANDS.md`
- 事实源边界：`spec/governance/source-of-truth-matrix.md`
- AI skill 真相源：`gitcode-cli/skills`
- Loop 标准包真相源：`gitcode-cli/loop-kits`
- GitCode 协作状态：GitCode issue / PR / labels / comments
- CI 执行结果：GitHub mirror Actions run

## 变更同步规则

- 命令行为变化：同步 `docs/COMMANDS.md`、`README.md`、相关外部使用说明
- 流程变化：同步 `spec/workflows/*`、`AGENTS.md`、`CLAUDE.md`
- Loop Engineering 变化：同步 `spec/loop/*`、`docs/LOOP_ENGINEERING*.md`、必要时同步 `gitcode-cli/skills` 和 `gitcode-cli/loop-kits`
- AI 执行方法变化：同步 `gitcode-cli/skills`
- schema / policy / hook / template / adapter 变化：同步 `gitcode-cli/loop-kits`
- CI 协同变化：同步 `spec/delivery/ci-workflows.md` 和 `spec/loop/mirror-ci-contract.md`

## 禁止

- 把 `.ai/skills/` 当作当前 skill 真相源
- 把 `.claude/skills/` 或 `.codex/skills/` 当作跨 AI 的唯一来源
- 在入口文档中定义与 `spec/` 冲突的项目规则
- 把 GitHub mirror CI 当作 GitCode 主仓状态
- 把一次性执行证据提交进主仓文档

## Loop Engineering 归档边界

- 一次性执行证据：GitCode issue / PR comment
- 长期工程规则：`spec/loop/`
- 用户可读说明：`docs/`
- AI 可复用执行方法：`gitcode-cli/skills`
- 标准 schema / policy / hook / template / adapter：`gitcode-cli/loop-kits`
- CI 原始事实：GitHub Actions run URL

判断标准：

- 会影响未来工程行为的，进 `spec/`
- 会指导用户理解的，进 `docs/`
- 会指导 AI 执行的，进 `gitcode-cli/skills`
- 会被机器复用的，进 `gitcode-cli/loop-kits`
- 只属于本次执行的，留在 issue / PR

## 当前阶段

当前阶段目标是 Loop Engineering Demo v1：

1. `gitcode-cli/cli` 定义规则和演示入口
2. `gitcode-cli/skills` 定义 AI 执行 skill
3. `gitcode-cli/loop-kits` 定义标准资源包

旧仓内 skill 分层只保留迁移说明。
