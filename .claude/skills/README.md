# Claude skill 历史适配层

本目录曾经作为 gitcode-cli 仓库内 Claude skill 适配层。

当前正式 skill 真相源已经迁移到独立仓库：

- https://gitcode.com/gitcode-cli/skills

## 当前边界

- 本目录只作为历史兼容层和迁移参考。
- 不在本目录新增正式 skill。
- Claude 使用的最新 GitCode CLI skill 应从 `gitcode-cli/skills` 安装或同步。
- 项目正式规则仍以 `spec/` 为准。

## Loop Engineering

Loop Engineering 相关执行 skill 应落在 `gitcode-cli/skills`：

- `gitcode-loop-engineering`
- `gitcode-loop-ci`
- `gitcode-loop-archive`

本目录不得定义低于 `spec/` 或 `gitcode-cli/skills` 的规则。
