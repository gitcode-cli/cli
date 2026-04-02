# .ai 目录说明

`.ai/` 是 gitcode-cli 的跨 AI 协作目录。

本目录的目标是：

- 为 Claude 和 Codex 提供共享的项目级 skill 真相源
- 避免把某一客户端的 skill 目录误当作全项目唯一来源
- 让 skill 可以版本化、评审、迁移和跨项目复用
- 提供可脱离本仓库使用的通用 `gc` skill 包

## 权威边界

- `.ai/skills/` 是共享 skill 真相源
- `.ai/distribution/` 是可分发的通用 skill 包
- `.claude/skills/` 是 Claude 适配层
- `.codex/skills/` 是 Codex 适配层

共享源定义场景和边界，适配层负责不同客户端如何落地这些场景。

`.ai/` 不是项目正式规则源。

项目正式规则仍以 [spec/README.md](../spec/README.md) 和 `spec/` 目录为准。
不同信息类型的事实边界见 [spec/governance/source-of-truth-matrix.md](../spec/governance/source-of-truth-matrix.md)。

## 仓库内 skill 与通用 skill 的区别

- `.ai/skills/` 面向 gitcode-cli 仓库内部协作
- `.ai/distribution/` 面向外部项目复用 `gc`

通用 skill 包不得依赖本仓库内部 `spec/` 或 `docs/COMMANDS.md` 路径。

当前可分发包的安装与分发说明见：

- `.ai/distribution/gc-core/INSTALL.md`

## 当前阶段说明

当前阶段先完成结构化治理：

- 新增 `.ai/skills/`
- 新增 `.codex/skills/`
- 为现有 skill 建立共享源与适配层映射

本阶段不做：

- 自动同步脚本
- 全量重写所有现有 skill 内容
- 把客户端差异抹平成完全相同的文本

当前已提供：

- `scripts/sync-ai-skills.sh`

该脚本用于基于共享源更新 Codex 适配层，并为缺失的 Claude 适配目录生成占位入口；不会覆盖现有 Claude skill 正文。
