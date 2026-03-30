# 规范文档索引

`spec/` 是 gitcode-cli 的正式规范目录，定义项目开发、测试、安全、评审和后续交付规则。

如果你要修改代码、文档、流程或 AI 协作规则，先从本目录开始。

## 阅读顺序

建议按以下顺序阅读：

1. [文档治理规范](./docs-governance.md)
2. [开发工作流程](./development-workflow.md)
3. [本地构建与打包规范](./build-and-package.md)
4. [发布流程规范](./release-process.md)
5. [代码质量门禁规范](./code-quality-gates.md)
6. [编码规范](./coding-standards.md)
7. [测试指南](./testing-guide.md)
8. [安全规范](./security.md)
9. `workflows/` 中对应操作流程

## 当前结构

```
spec/
├── README.md                # 本文件：规范入口
├── docs-governance.md       # 文档治理规范
├── development-workflow.md  # 开发工作流程
├── build-and-package.md     # 本地构建与打包规范
├── release-process.md       # 发布流程规范
├── code-quality-gates.md    # 代码质量门禁规范
├── coding-standards.md      # 编码规范
├── testing-guide.md         # 测试指南
├── command-template.md      # 命令开发模板
├── security.md              # 安全规范
└── workflows/               # 独立操作流程
    ├── issue-workflow.md    # Issue 流程
    ├── pr-workflow.md       # PR 流程
    ├── review-workflow.md   # 评审流程
    └── test-workflow.md     # 测试流程
```

## 权威边界

本目录定义的规则具有以下边界：

- `spec/` 是项目正式规范唯一来源
- `docs/COMMANDS.md` 定义命令行为
- `README.md` 是项目总入口，不是规范源
- `AGENTS.md` 和 `CLAUDE.md` 是 AI 协作入口，不是项目规则源

如果其他文档与本目录冲突，以 `spec/` 为准。

## 核心规范

| 文档 | 说明 |
|------|------|
| [文档治理规范](./docs-governance.md) | 文档分层、唯一真相源、AI 协作分层、分阶段实施方案 |
| [开发工作流程](./development-workflow.md) | 完整流程、分支规范、禁止行为、检查清单 |
| [本地构建与打包规范](./build-and-package.md) | 标准构建命令、打包方式、产物边界和验证要求 |
| [发布流程规范](./release-process.md) | 版本规则、发布步骤、release notes 和发布后验证 |
| [代码质量门禁规范](./code-quality-gates.md) | 本地门禁、PR 门禁、合并门禁和 blocker 判定 |
| [编码规范](./coding-standards.md) | 命名规范、文件结构、错误处理、代码风格 |
| [测试指南](./testing-guide.md) | 单元测试、实际命令测试、测试仓库限制 |
| [命令开发模板](./command-template.md) | 新命令开发模板、API 客户端用法、输出处理 |
| [安全规范](./security.md) | Token 管理、敏感信息保护、安全审查 |

## 操作流程

| 文档 | 说明 |
|------|------|
| [Issue 流程](./workflows/issue-workflow.md) | Issue 创建、标签、验证、关闭 |
| [PR 流程](./workflows/pr-workflow.md) | 分支创建、代码提交、PR 创建与合并 |
| [评审流程](./workflows/review-workflow.md) | Issue 评论、PR 审查评论 |
| [测试流程](./workflows/test-workflow.md) | 单元测试、实际命令测试 |

## 计划补齐的规范

以下规范仍在治理计划中，但尚未落地：

- `spec/ci-workflows.md`

其中 `spec/ci-workflows.md` 放在最后阶段实施，因为当前 GitCode CI 条件尚未具备。

## 相关文档

| 文档 | 位置 | 说明 |
|------|------|------|
| README.md | 根目录 | 项目总入口 |
| AGENTS.md | 根目录 | Codex 项目级入口 |
| CLAUDE.md | 根目录 | Claude 项目级入口 |
| docs/COMMANDS.md | docs/ | 命令行为说明 |
| docs/PACKAGING.md | docs/ | 打包发布使用说明 |
| issues-plan/PROGRESS.md | issues-plan/ | 当前项目状态 |

**最后更新**: 2026-03-30
