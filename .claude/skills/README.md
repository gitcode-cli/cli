# Claude skill 适配层

本目录是 gitcode-cli 的 Claude skill 适配层。

## 角色

- 为 Claude 提供可直接使用的项目内 skill
- 保留 Claude 侧需要的描述、入口和引用方式
- 与共享 skill 真相源保持同步

## 权威边界

- 共享 skill 真相源：`../.ai/skills/`
- Claude 适配层：本目录
- 项目正式规则：`spec/`

本目录不是共享 skill 的唯一来源。

## 当前映射

| Claude 目录 | 共享 skill |
|-------------|------------|
| `gitcode-cli/` | `.ai/skills/gitcode-cli/` |
| `pr-reviewer/` | `.ai/skills/pr-reviewer/` |
| `issue-reviewer/` | `.ai/skills/issue-reviewer/` |
| `gc-dev-setup/` | `.ai/skills/gc-dev-setup/` |
| `gitcode-cmd-generator/` | `.ai/skills/gitcode-cmd-generator/` |

## 当前阶段说明

本阶段只完成适配层定位和映射关系。

现有 Claude skill 内容继续保留，后续再按共享源逐步校准。
