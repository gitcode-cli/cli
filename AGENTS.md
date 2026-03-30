# AGENTS.md

本文件是 Codex 和其他通用代理在 gitcode-cli 仓库中的项目级入口文档。

如果任务涉及代码、文档、流程、评审或发布，请先阅读：

1. [spec/README.md](./spec/README.md)
2. [docs/README.md](./docs/README.md)
3. [README.md](./README.md)

## 1. 入口职责

`AGENTS.md` 的职责是：

- 为 Codex 提供仓库级入口
- 指向正式规范、用户文档和后续 skill 分层
- 约束代理不得绕过项目正式规则

`AGENTS.md` 不是项目规则源。

项目正式规则以 [spec/README.md](./spec/README.md) 和 `spec/` 目录中的规范文档为准。

## 2. 必读文档

按以下顺序阅读：

1. [文档治理规范](./spec/docs-governance.md)
2. [开发工作流程](./spec/development-workflow.md)
3. [本地构建与打包规范](./spec/build-and-package.md)
4. [发布流程规范](./spec/release-process.md)
5. [代码质量门禁规范](./spec/code-quality-gates.md)
6. [编码规范](./spec/coding-standards.md)
7. [测试指南](./spec/testing-guide.md)
8. [安全规范](./spec/security.md)

具体流程任务再进入：

- [Issue 流程](./spec/workflows/issue-workflow.md)
- [PR 流程](./spec/workflows/pr-workflow.md)
- [评审流程](./spec/workflows/review-workflow.md)
- [测试流程](./spec/workflows/test-workflow.md)

## 3. 核心执行规则

代理在本仓库中必须遵守：

- 项目命令固定为 `gc`
- 项目正式规范以 `spec/` 为准
- 命令行为以 [docs/COMMANDS.md](./docs/COMMANDS.md) 为准
- 当前状态以 [issues-plan/PROGRESS.md](./issues-plan/PROGRESS.md) 为准
- 代码或流程变化后必须同步检查相关文档
- 实际命令测试只能使用 `infra-test/*`
- 不得在 `main` 直接开发
- 不得提交构建产物、评估输出或真实凭证

## 4. Codex 入口边界

当前仓库内的 Codex 项目级入口是：

- `AGENTS.md`

当前仓库已引入：

- `.ai/skills/` 共享 skill 真相源
- `.codex/skills/` Codex 适配层

Codex 仍应先以 `spec/` 和本文件为主要入口，再进入共享源或适配层。

## 5. 常用入口

- 用户文档入口：[docs/README.md](./docs/README.md)
- 命令手册：[docs/COMMANDS.md](./docs/COMMANDS.md)
- 认证说明：[docs/AUTH.md](./docs/AUTH.md)
- 回归说明：[docs/REGRESSION.md](./docs/REGRESSION.md)
- 打包说明：[docs/PACKAGING.md](./docs/PACKAGING.md)
- 发布说明：[RELEASE.md](./RELEASE.md)
- 贡献说明：[CONTRIBUTING.md](./CONTRIBUTING.md)
- Claude 入口：[CLAUDE.md](./CLAUDE.md)

## 6. 当前阶段说明

当前治理已完成：

- `spec/docs-governance.md`
- `spec/build-and-package.md`
- `spec/release-process.md`
- `spec/code-quality-gates.md`
- `.ai/skills/` 共享 skill 真相源
- `.codex/skills/` 适配层结构

当前阶段目标是继续校准共享源与客户端适配层，而不是新建额外规则源。
