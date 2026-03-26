# CLAUDE.md - AI 辅助开发指南

> 项目概述和功能介绍请参阅 [README.md](./README.md)

本文档为 Claude Code 提供 gitcode-cli 项目开发指导。

## 核心信息

| 项目 | 值 |
|------|-----|
| 命令名 | `gc` |
| 语言 | Go 1.21+ |
| 框架 | Cobra |
| 配置目录 | `~/.config/gc/` |
| 环境变量前缀 | `GC_*` |

## 目录结构

```
gitcode-cli/
├── cmd/gc/           # 程序入口
├── pkg/cmd/          # 命令实现
├── pkg/cmdutil/      # Factory、工具函数
├── pkg/iostreams/    # IO 流管理
├── internal/         # 内部包（config、authflow、keyring、prompter）
├── api/              # API 客户端
└── git/              # Git 操作
```

## 关键要求（必须遵守）

1. **命名规范**: 命令名 `gc`，禁止使用 `gt`；环境变量 `GC_*`
2. **测试仓库**: 只能用 `infra-test/gctest1` 和 `gitcode-cli/cli`
3. **开发流程**: Issue → 分支 → 开发 → 测试 → PR → 合并
4. **安全**: Token 必须使用环境变量，禁止硬编码
5. **提交限制**: 单次提交不超过 800 行

## 开发工作流程（重要！）

**严格遵守以下流程，违反将导致代码管理混乱！**

```
提交 Issue → 打标签 → 创建分支 → 分支开发 → 编写测试 → 实际命令测试 → 提交 PR → Issue 评论 → PR 审查评论 → 关闭 Issue → 合并 PR
```

详细流程参见 [开发工作流程](./spec/development-workflow.md)，包含：
- 完整流程步骤
- 分支命名规范
- 标签使用规范
- 禁止行为清单
- 检查清单

## 文档索引

### 规范文档 (spec/)

| 文档 | 说明 |
|------|------|
| [编码规范](./spec/coding-standards.md) | 命名、文件结构、错误处理 |
| [测试指南](./spec/testing-guide.md) | 单元测试、实际命令测试 |
| [命令开发模板](./spec/command-template.md) | 新命令开发模板和示例 |
| [安全规范](./spec/security.md) | Token 管理、敏感信息保护 |

### 开发流程 (spec/workflows/)

| 文档 | 说明 |
|------|------|
| [Issue 流程](./spec/workflows/issue-workflow.md) | Issue 创建、标签、验证、关闭 |
| [PR 流程](./spec/workflows/pr-workflow.md) | PR 创建、关联、合并 |
| [评审流程](./spec/workflows/review-workflow.md) | Issue 评论、PR 审查评论 |
| [测试流程](./spec/workflows/test-workflow.md) | 单元测试、实际命令测试流程 |
| [构建打包流程](./spec/workflows/build-workflow.md) | 本地构建、打包命令 |
| [Release 流程](./spec/workflows/release-workflow.md) | 创建 Release、上传包 |

### 其他文档

| 文档 | 说明 |
|------|------|
| [安全策略](./SECURITY.md) | 安全策略 |
| [命令使用指南](./docs/COMMANDS.md) | 所有命令的使用示例 |
| [打包发布指南](./docs/PACKAGING.md) | DEB/RPM/PyPI 打包 |
| [需求管理](./issues-plan/) | 需求清单和里程碑 |

## 快速命令

```bash
# 构建
go build -o ./gc ./cmd/gc

# 测试
go test ./...

# 实际命令测试
./gc issue list -R infra-test/gctest1
```

## 参考资源

- [API 文档](https://gitcode.com/afly-infra/gc-api-doc)
- [GitHub CLI 源码](https://github.com/cli/cli)

---

**最后更新**: 2026-03-26