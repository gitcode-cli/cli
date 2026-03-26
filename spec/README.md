# 规范文档索引

本目录包含 gitcode-cli 项目的开发规范和流程文档。

## 文档结构

```
spec/
├── README.md                # 本文件：文档索引
├── development-workflow.md  # 开发工作流程（重要！）
├── coding-standards.md      # 编码规范
├── testing-guide.md         # 测试指南
├── command-template.md      # 命令开发模板
├── security.md              # 安全规范
└── workflows/               # 开发流程文档
    ├── issue-workflow.md    # Issue 流程
    ├── pr-workflow.md       # PR 流程
    ├── review-workflow.md   # 评审流程
    ├── test-workflow.md     # 测试流程
    ├── build-workflow.md    # 构建打包流程
    └── release-workflow.md  # Release 流程
```

## 快速导航

### 开发工作流程（重要！）

| 文档 | 说明 |
|------|------|
| [开发工作流程](./development-workflow.md) | **完整流程、分支规范、禁止行为、检查清单** |

### 规范文档

| 文档 | 说明 |
|------|------|
| [编码规范](./coding-standards.md) | 命名规范、文件结构、错误处理、代码风格 |
| [测试指南](./testing-guide.md) | 单元测试、实际命令测试、测试仓库限制 |
| [命令开发模板](./command-template.md) | 新命令开发模板、API 客户端用法、输出处理 |
| [安全规范](./security.md) | Token 管理、敏感信息保护、测试安全 |

### 开发流程

| 文档 | 说明 |
|------|------|
| [Issue 流程](./workflows/issue-workflow.md) | Issue 创建、标签、验证、关闭 |
| [PR 流程](./workflows/pr-workflow.md) | 分支创建、代码提交、PR 创建与合并 |
| [评审流程](./workflows/review-workflow.md) | Issue 评论、PR 审查评论、检查清单 |
| [测试流程](./workflows/test-workflow.md) | 单元测试、实际命令测试、验证结果 |
| [构建打包流程](./workflows/build-workflow.md) | 本地构建、打包命令、版本同步 |
| [Release 流程](./workflows/release-workflow.md) | 创建 Release、上传包、更新版本号 |

## 使用指南

### 开发新功能
1. 阅读 [Issue 流程](./workflows/issue-workflow.md) 创建 Issue
2. 阅读 [编码规范](./coding-standards.md) 了解命名和结构
3. 阅读 [命令开发模板](./command-template.md) 学习开发模式
4. 阅读 [测试指南](./testing-guide.md) 编写测试
5. 阅读 [PR 流程](./workflows/pr-workflow.md) 提交代码

### 修复 Bug
1. 阅读 [Issue 流程](./workflows/issue-workflow.md) 验证问题
2. 阅读 [测试指南](./testing-guide.md) 编写回归测试
3. 阅读 [PR 流程](./workflows/pr-workflow.md) 提交修复

### 发布版本
1. 阅读 [构建打包流程](./workflows/build-workflow.md) 构建包
2. 阅读 [Release 流程](./workflows/release-workflow.md) 发布

## 相关文档

| 文档 | 位置 | 说明 |
|------|------|------|
| CLAUDE.md | 根目录 | AI 辅助开发核心指南 |
| SECURITY.md | 根目录 | 安全策略 |
| RELEASE.md | 根目录 | 发布指南 |
| docs/COMMANDS.md | docs/ | 命令使用指南 |
| docs/PACKAGING.md | docs/ | 打包发布指南 |
| issues-plan/ | 根目录 | 需求管理和里程碑 |

---

**最后更新**: 2026-03-26