# 历史 skill 兼容层

本目录曾经存放 gitcode-cli 的共享 skill 源。

当前正式 skill 真相源已经迁移到独立仓库：

- https://gitcode.com/gitcode-cli/skills

## 当前边界

- 本目录只作为历史兼容层和迁移参考。
- 不在本目录新增正式 skill。
- 不把本目录作为跨 AI 的 skill 真相源。
- 项目规则仍以 `spec/` 为准。

## 新增或修改 skill

请在 `gitcode-cli/skills` 仓库中新增或修改：

- GitCode CLI 通用命令 skill
- GitCode 工作流 skill
- Loop Engineering 执行 skill

本仓库内 `.claude/skills/` 与 `.codex/skills/` 也只作为历史适配层，不再承担真相源职责。
