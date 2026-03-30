# gc-core 通用 skill 包

`gc-core` 是 gitcode-cli 仓库导出的第一版通用 skill 包。

它的目标是：

- 让其他项目在不依赖本仓库内部 `spec/`、`docs/` 路径的前提下使用 `gc`
- 提供一组可直接复制、安装或再分发的基础 skill
- 将“仓库内协作 skill”和“可分发通用 skill”明确分开

## 包含内容

第一版固定包含 7 个 skill：

- `gc-auth`
- `gc-repo`
- `gc-issue`
- `gc-pr`
- `gc-release`
- `gc-review`
- `gc-regression`

## 设计边界

本目录下的 skill：

- 不依赖本仓库内部 `docs/COMMANDS.md`
- 不依赖本仓库内部 `spec/`
- 不依赖本仓库的目录结构

这些 skill 面向的是“在其他项目里使用 `gc`”，不是“继续开发 gitcode-cli 仓库本身”。

## 使用方式

典型使用方式：

1. 将需要的 skill 目录复制到目标环境
2. 将 skill 安装到目标 AI 客户端的 skill 目录
3. 根据目标项目的仓库、分支、权限和工作流进行最小适配

详细安装与分发步骤见：

- [INSTALL.md](./INSTALL.md)

## 注意事项

- `gc` 必须已安装并可执行
- 需要访问私有仓库时，应先完成 `gc` 认证
- GitCode 平台能力以真实 API 和 CLI 当前行为为准
- 这些通用 skill 不会替代目标项目自己的开发流程规范
