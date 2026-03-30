# Codex skill 适配层

本目录是 gitcode-cli 的 Codex skill 适配层。

## 角色

- 为 Codex 提供仓库内可版本化的 skill 适配入口
- 与共享 skill 真相源保持同步
- 避免仅依赖用户本地 `~/.codex/skills/` 作为唯一来源

## 权威边界

- 共享 skill 真相源：`../.ai/skills/`
- Codex 适配层：本目录
- 项目正式规则：`spec/`

本目录不是项目规则源，只负责 Codex 侧适配。

## 当前映射

| Codex 目录 | 共享 skill |
|------------|------------|
| `gitcode-cli/` | `.ai/skills/gitcode-cli/` |
| `pr-reviewer/` | `.ai/skills/pr-reviewer/` |
| `issue-reviewer/` | `.ai/skills/issue-reviewer/` |
| `gc-dev-setup/` | `.ai/skills/gc-dev-setup/` |
| `gitcode-cmd-generator/` | `.ai/skills/gitcode-cmd-generator/` |

## 当前阶段说明

本阶段先补齐目录结构和适配说明。

后续阶段再根据共享 skill 真相源逐步补充更细的 Codex 侧技能文档。
