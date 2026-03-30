# AI 协作规范

## 职责

定义 Codex、Claude、共享 skill 源和客户端适配层之间的关系。

## 适用场景

- 调整 `AGENTS.md` 或 `CLAUDE.md`
- 调整 `.ai/skills/`、`.claude/skills/`、`.codex/skills/`
- 设计多 AI 协作边界

## 必须

- 以 `spec/` 作为项目正式规则源
- 以 `AGENTS.md` 和 `CLAUDE.md` 作为各自客户端入口
- 以 `.ai/skills/` 作为共享 skill 真相源
- 让 `.claude/skills/` 与 `.codex/skills/` 只承担适配职责

## 禁止

- 在 AI 入口文档中定义与 `spec/` 冲突的项目规则
- 把任一客户端的适配层当作跨 AI 的唯一来源
- 让外部可分发 skill 依赖仓库私有 `spec/` 或 `docs/` 路径

## 同步要求

- AI 入口变化时同步 `AGENTS.md`、`CLAUDE.md`
- 共享源变化时同步 `.ai/skills/`、相关适配层和分发说明
- 外部可分发 skill 变化时同步 `.ai/distribution/`

## 不负责什么

- 命令实现细节
- 一般编码规范
- 本地构建与发布流程

## 权威关系

按以下顺序理解：

1. `spec/` 定义项目正式规则
2. `AGENTS.md` / `CLAUDE.md` 定义不同 AI 客户端如何进入规则体系
3. `.ai/skills/` 定义共享场景技能
4. `.claude/skills/` / `.codex/skills/` 定义客户端适配
5. `.ai/distribution/` 定义可分发的通用 skill 包

## 当前结构

- Codex 入口：[../../AGENTS.md](../../AGENTS.md)
- Claude 入口：[../../CLAUDE.md](../../CLAUDE.md)
- 共享源：[../../.ai/README.md](../../.ai/README.md)
- 文档治理：[docs-governance.md](./docs-governance.md)

## 下一步去看哪里

- 如果你在改技能分发，继续看 [../../.ai/distribution/gc-core/README.md](../../.ai/distribution/gc-core/README.md)
- 如果你在改同步边界，继续看 [docs-governance.md](./docs-governance.md)
