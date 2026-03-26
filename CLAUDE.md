# CLAUDE.md - AI 辅助开发指南

> 项目概述和功能介绍请参阅 [README.md](./README.md)

## 核心信息

| 项目 | 值 |
|------|-----|
| 命令名 | `gc` |
| 语言 | Go 1.21+ |
| 框架 | Cobra |
| 配置目录 | `~/.config/gc/` |
| 环境变量前缀 | `GC_*` |

## 关键要求（必须遵守）

1. **命名规范**: 命令名 `gc`，禁止使用 `gt`；环境变量 `GC_*`
2. **测试仓库**: 只能用 `infra-test/gctest1` 和 `gitcode-cli/cli`
3. **开发流程**: Issue → 分支 → 开发 → 测试 → PR → 合并
4. **安全**: Token 必须使用环境变量，禁止硬编码
5. **提交限制**: 单次提交不超过 800 行

---

## 场景入口

根据任务类型，优先调用对应的 Skill。

### 优先调用 Skill

| 场景 | Skill | 说明 |
|------|-------|------|
| 评审 Issue | `/issue-reviewer` | 自动分析 Issue、添加评论和标签 |
| 评审 PR | `/pr-reviewer` | 检查代码质量、安全问题、规范合规 |
| 开发新命令 | `/gitcode-cmd-generator` | 生成命令代码模板和测试文件 |
| 初始化环境 | `/gc-dev-setup` | 检查并初始化本地开发环境 |

### 独立操作（无 Skill）

| 场景 | 文档 | 说明 |
|------|------|------|
| CI 自动发布 | [RELEASE.md](./RELEASE.md) | GitHub Actions 自动发布流程 |
| 本地打包发布 | [docs/PACKAGING.md](./docs/PACKAGING.md) | 本地构建 DEB/RPM/PyPI 包 |
| 安全策略 | [SECURITY.md](./SECURITY.md) | Token 管理、敏感信息保护 |

### 完整开发流程

新功能开发或 Bug 修复，遵循完整流程：

```
提交 Issue → 打标签 → 创建分支 → 分支开发 → 编写测试 → 实际命令测试 → 安全审查 → 提交 PR → Issue 评论 → PR 审查评论 → 关闭 Issue → 合并 PR
```

详细流程参见 [开发工作流程](./spec/development-workflow.md)。

---

## 规范文档索引

完整规范文档参见 [spec/README.md](./spec/README.md)。

| 文档 | 说明 |
|------|------|
| [开发工作流程](./spec/development-workflow.md) | 完整流程、分支规范、禁止行为 |
| [编码规范](./spec/coding-standards.md) | 命名、文件结构、错误处理 |
| [测试指南](./spec/testing-guide.md) | 单元测试、实际命令测试 |
| [命令开发模板](./spec/command-template.md) | 新命令开发模板和示例 |
| [安全规范](./spec/security.md) | Token 管理、敏感信息保护 |

---

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