# issues-plan 需求管理目录

本目录用于管理 gitcode-cli 项目的全量开发需求和交付、验收进展。

## 目录结构

```
issues-plan/
├── README.md                      # 本文件，需求管理说明
├── 00-overview.md                 # 项目总览
├── 01-requirements-overview.md    # 需求总清单
├── 02-architecture.md             # 架构设计说明
├── 03-module-auth.md              # 认证模块需求
├── 04-module-repo.md              # 仓库模块需求
├── 05-module-issue.md             # Issue模块需求
├── 06-module-pr.md                # PR模块需求
├── 07-api-client.md               # API客户端需求
├── 08-config.md                   # 配置管理需求
├── 09-testing.md                  # 测试需求
├── 10-deployment.md               # 部署发布需求
├── 11-next-development-plan.md    # 下一步开发计划
├── agent-cli-improvement/         # Agent-Friendly CLI 专项改进方案
├── PROGRESS.md                    # 项目交付进度跟踪
└── milestones/                    # 里程碑追踪
    ├── m1-foundation.md           # 里程碑1：基础架构
    ├── m2-auth.md                 # 里程碑2：认证功能
    ├── m3-repo.md                 # 里程碑3：仓库功能
    ├── m4-issue.md                # 里程碑4：Issue功能
    ├── m5-pr.md                   # 里程碑5：PR功能
    └── m6-release.md              # 里程碑6：Release功能
```

## 文档索引

| 文档 | 说明 |
|------|------|
| [00-overview.md](./00-overview.md) | 项目总览、设计目标、技术栈 |
| [01-requirements-overview.md](./01-requirements-overview.md) | 全量需求清单、优先级、状态追踪 |
| [02-architecture.md](./02-architecture.md) | 架构设计、分层结构、设计模式 |
| [03-module-auth.md](./03-module-auth.md) | 认证模块详细需求 |
| [04-module-repo.md](./04-module-repo.md) | 仓库模块详细需求 |
| [05-module-issue.md](./05-module-issue.md) | Issue模块详细需求 |
| [06-module-pr.md](./06-module-pr.md) | PR模块详细需求（重点：代码检视） |
| [07-api-client.md](./07-api-client.md) | API客户端设计需求 |
| [08-config.md](./08-config.md) | 配置管理设计需求 |
| [09-testing.md](./09-testing.md) | 测试策略和测试需求 |
| [10-deployment.md](./10-deployment.md) | 部署发布流程需求 |
| [11-next-development-plan.md](./11-next-development-plan.md) | 收口与一致性阶段开发计划 |
| [agent-cli-improvement/README.md](./agent-cli-improvement/README.md) | 参考 agent-cli-guide 的 CLI 改进专项方案 |

## 如何使用需求文档

### 需求状态定义

| 状态 | 说明 |
|------|------|
| 📋 待开发 | 需求已定义，等待开发 |
| 🚧 开发中 | 正在开发中 |
| ✅ 已完成 | 功能已实现并通过验收 |
| ⏸️ 暂停 | 开发暂停 |
| ❌ 取消 | 需求已取消 |

### 优先级定义

| 优先级 | 说明 |
|--------|------|
| P0 | 核心功能，第一版本必须实现 |
| P1 | 重要功能，第二版本实现 |
| P2 | 增强功能，后续版本实现 |
| P3 | 可选功能，根据反馈决定 |

### 需求更新流程

1. 在对应的模块文档中添加/修改需求
2. 更新 01-requirements-overview.md 中的需求清单
3. 在对应的里程碑文档中更新任务分解
4. 提交 PR 并进行评审

## 相关文档

- [CLAUDE.md](../CLAUDE.md) - AI 辅助开发指南
- [COMMANDS.md](../docs/COMMANDS.md) - 命令使用指南
- [PACKAGING.md](../docs/PACKAGING.md) - 打包发布指南
- [SECURITY.md](../SECURITY.md) - 安全策略
- [gc-design](https://gitcode.com/gitcode-cli/gc-design) - 详细设计文档
- [gc-api-doc](https://gitcode.com/gitcode-cli/gc-api-doc) - GitCode API 文档

---

**最后更新**: 2026-03-30
