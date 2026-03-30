# CLAUDE.md

本文件是 Claude 在 gitcode-cli 仓库中的项目级入口文档。

如果任务涉及代码、文档、流程、评审或发布，请先阅读：

1. [spec/README.md](./spec/README.md)
2. [docs/README.md](./docs/README.md)
3. [README.md](./README.md)

## 1. 入口职责

`CLAUDE.md` 的职责是：

- 告诉 Claude 应先看哪些正式规范
- 说明仓库中的 AI 协作入口和场景化 skill
- 提醒 Claude 不得绕过 `spec/` 中定义的正式规则

`CLAUDE.md` 不是项目规则源。

项目规则以 [spec/README.md](./spec/README.md) 和 `spec/` 目录中的正式规范为准。

## 2. 必读文档

先读：

1. [spec/README.md](./spec/README.md)

再根据任务进入对应规范，不要机械顺序通读全部文档。

常用任务入口：

- 改命令行为：`spec/workflows/development-workflow.md`、`spec/governance/docs-governance.md`、`spec/foundations/code-quality-gates.md`
- 改 API / auth / config：`spec/foundations/coding-standards.md`、`spec/foundations/security.md`、`spec/foundations/testing-guide.md`
- 补测试或做真实命令验证：`spec/foundations/testing-guide.md`、`spec/workflows/test-workflow.md`
- 提交 PR / 做 review：`spec/workflows/pr-workflow.md`、`spec/workflows/review-workflow.md`
- 改构建 / 打包 / 发布：`spec/delivery/build-and-package.md`、`spec/delivery/release-process.md`

如果任务是具体操作流程，再进入：

- [Issue 流程](./spec/workflows/issue-workflow.md)
- [PR 流程](./spec/workflows/pr-workflow.md)
- [评审流程](./spec/workflows/review-workflow.md)
- [测试流程](./spec/workflows/test-workflow.md)

## 3. 核心执行规则

Claude 在本仓库中必须遵守：

- 命令名固定为 `gc`
- 项目正式规范以 `spec/` 为准
- 命令行为以 [docs/COMMANDS.md](./docs/COMMANDS.md) 为准
- 当前状态以 [issues-plan/PROGRESS.md](./issues-plan/PROGRESS.md) 为准
- 代码变更后必须同步检查相关文档
- 实际命令测试只能使用 `infra-test/*`
- 不得在 `main` 直接开发
- 不得在文档中写入真实 token 或凭证

## 4. AI 协作入口

Claude 侧的项目内 skill 目录位于：

- [.claude/skills](./.claude/skills)

使用原则：

- `.ai/skills/` 是共享 skill 真相源
- `.claude/skills/` 是 Claude 适配层
- 场景技能不能定义与 `spec/` 冲突的项目规则
- 行为变化后，必须同步检查相关 skill 文档

Codex 侧适配层位于 `.codex/skills/`。

## 5. 常用入口

- 用户命令手册：[docs/COMMANDS.md](./docs/COMMANDS.md)
- 认证说明：[docs/AUTH.md](./docs/AUTH.md)
- 回归说明：[docs/REGRESSION.md](./docs/REGRESSION.md)
- 打包说明：[docs/PACKAGING.md](./docs/PACKAGING.md)
- 发布说明：[RELEASE.md](./RELEASE.md)
- 贡献说明：[CONTRIBUTING.md](./CONTRIBUTING.md)

## 6. 当前阶段说明

当前治理已完成：

- `spec/governance/docs-governance.md`
- `spec/delivery/build-and-package.md`
- `spec/delivery/release-process.md`
- `spec/foundations/code-quality-gates.md`
- `.ai/skills/` 共享 skill 真相源
- `.codex/skills/` 适配层结构

当前阶段目标是继续校准共享 skill 与客户端适配层，而不是扩展新的规则源。
