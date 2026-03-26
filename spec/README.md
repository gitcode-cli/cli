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
└── workflows/               # 独立操作流程
    ├── issue-workflow.md    # Issue 操作流程
    ├── pr-workflow.md       # PR 操作流程
    ├── review-workflow.md   # 评审流程
    └── test-workflow.md     # 测试流程
```

---

## 核心文档

| 文档 | 说明 |
|------|------|
| [开发工作流程](./development-workflow.md) | **完整流程、分支规范、禁止行为、检查清单** |

---

## 规范文档

| 文档 | 说明 |
|------|------|
| [编码规范](./coding-standards.md) | 命名规范、文件结构、错误处理、代码风格 |
| [测试指南](./testing-guide.md) | 单元测试、实际命令测试、测试仓库限制 |
| [命令开发模板](./command-template.md) | 新命令开发模板、API 客户端用法、输出处理 |
| [安全规范](./security.md) | Token 管理、敏感信息保护、安全审查 |

---

## 独立操作流程

| 文档 | 说明 | 对应 Skill |
|------|------|------------|
| [Issue 流程](./workflows/issue-workflow.md) | Issue 创建、标签、验证、关闭 | `/issue-reviewer` |
| [PR 流程](./workflows/pr-workflow.md) | 分支创建、代码提交、PR 创建与合并 | `/pr-reviewer` |
| [评审流程](./workflows/review-workflow.md) | Issue 评论、PR 审查评论 | `/issue-reviewer`, `/pr-reviewer` |
| [测试流程](./workflows/test-workflow.md) | 单元测试、实际命令测试 | - |

---

## 相关文档

| 文档 | 位置 | 说明 |
|------|------|------|
| CLAUDE.md | 根目录 | AI 辅助开发入口（场景 + Skill 索引） |
| SECURITY.md | 根目录 | 安全策略 |
| RELEASE.md | 根目录 | 发布指南（GitHub Actions） |
| docs/COMMANDS.md | docs/ | 命令使用指南 |
| docs/PACKAGING.md | docs/ | 打包发布指南 |
| issues-plan/ | 根目录 | 需求管理和里程碑 |

---

**最后更新**: 2026-03-26