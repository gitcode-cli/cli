# 需求总清单

本文档列出 gitcode-cli 项目的全量需求，按模块分类，包含优先级和状态追踪。

## 需求汇总

| 模块 | 需求数量 | P0 | P1 | P2 | 完成状态 |
|------|----------|----|----|----|----------|
| 认证 (auth) | 6 | 6 | 0 | 0 | 6/6 ✅ |
| 仓库 (repo) | 6 | 5 | 1 | 0 | 6/6 ✅ |
| Issue | 7 | 5 | 2 | 0 | 7/7 ✅ |
| PR | 10 | 8 | 2 | 0 | 10/10 ✅ |
| Release | 6 | 4 | 2 | 0 | 6/6 ✅ |
| API客户端 | 5 | 5 | 0 | 0 | 5/5 ✅ |
| 配置管理 | 4 | 3 | 1 | 0 | 4/4 ✅ |
| 测试 | 4 | 4 | 0 | 0 | 4/4 ✅ |
| 部署 | 3 | 3 | 0 | 0 | 3/3 ✅ |
| 文档 | 5 | 5 | 0 | 0 | 5/5 ✅ |
| **总计** | **56** | **48** | **8** | **0** | **56/56** |

---

## 认证模块 (auth)

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| AUTH-001 | auth login - OAuth Device Flow 认证 | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |
| AUTH-002 | auth login --with-token Token 认证 | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |
| AUTH-003 | auth logout 登出账户 | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |
| AUTH-004 | auth status 查看认证状态 | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |
| AUTH-005 | auth token 打印认证 Token | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |
| AUTH-006 | auth switch 切换账户 | P0 | ✅ 已完成 | [03-module-auth.md](./03-module-auth.md) |

---

## 仓库模块 (repo)

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| REPO-001 | repo clone 克隆仓库 | P0 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |
| REPO-002 | repo create 创建仓库 | P0 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |
| REPO-003 | repo fork Fork 仓库 | P0 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |
| REPO-004 | repo view 查看仓库 | P0 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |
| REPO-005 | repo list 列出仓库 | P0 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |
| REPO-006 | repo delete 删除仓库 | P1 | ✅ 已完成 | [04-module-repo.md](./04-module-repo.md) |

---

## Issue 模块

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| ISSUE-001 | issue create 创建 Issue | P0 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-002 | issue list 列出 Issues | P0 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-003 | issue view 查看 Issue | P0 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-004 | issue close 关闭 Issue | P0 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-005 | issue reopen 重开 Issue | P0 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-006 | issue comment 添加评论 | P1 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |
| ISSUE-007 | issue edit 编辑 Issue | P1 | ✅ 已完成 | [05-module-issue.md](./05-module-issue.md) |

---

## PR 模块

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| PR-001 | pr create 创建 PR | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-002 | pr list 列出 PRs | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-003 | pr view 查看 PR | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-004 | pr checkout 检出 PR 分支 | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-005 | pr merge 合并 PR | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-006 | pr close 关闭 PR | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-007 | pr reopen 重开 PR | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-008 | pr review 代码检视（重点功能） | P0 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-009 | pr diff 查看 PR 差异 | P1 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |
| PR-010 | pr ready 标记为就绪/WIP | P1 | ✅ 已完成 | [06-module-pr.md](./06-module-pr.md) |

---

## Release 模块

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| REL-001 | release create 创建 Release | P0 | ✅ 已完成 | [PACKAGING.md](../docs/PACKAGING.md) |
| REL-002 | release list 列出 Releases | P0 | ✅ 已完成 | [COMMANDS.md](../docs/COMMANDS.md) |
| REL-003 | release view 查看 Release | P0 | ✅ 已完成 | [COMMANDS.md](../docs/COMMANDS.md) |
| REL-004 | release delete 删除 Release | P0 | ✅ 已完成 | [COMMANDS.md](../docs/COMMANDS.md) |
| REL-005 | release edit 编辑 Release | P1 | ✅ 已完成 | [COMMANDS.md](../docs/COMMANDS.md) |
| REL-006 | release upload/download 资产管理 | P1 | ✅ 已完成 | [PACKAGING.md](../docs/PACKAGING.md) |

---

## API 客户端

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| API-001 | REST API 封装 | P0 | ✅ 已完成 | [07-api-client.md](./07-api-client.md) |
| API-002 | 认证中间件 | P0 | ✅ 已完成 | [07-api-client.md](./07-api-client.md) |
| API-003 | 缓存机制 | P0 | ✅ 已完成 | [07-api-client.md](./07-api-client.md) |
| API-004 | 重试机制 | P0 | ✅ 已完成 | [07-api-client.md](./07-api-client.md) |
| API-005 | 错误处理 | P0 | ✅ 已完成 | [07-api-client.md](./07-api-client.md) |

---

## 配置管理

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| CFG-001 | YAML 配置格式 | P0 | ✅ 已完成 | [08-config.md](./08-config.md) |
| CFG-002 | 配置存储位置 | P0 | ✅ 已完成 | [08-config.md](./08-config.md) |
| CFG-003 | Keyring 安全存储 | P0 | ✅ 已完成 | [08-config.md](./08-config.md) |
| CFG-004 | 环境变量支持 | P1 | ✅ 已完成 | [08-config.md](./08-config.md) |

---

## 测试

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| TEST-001 | 单元测试框架 | P0 | ✅ 已完成 | [09-testing.md](./09-testing.md) |
| TEST-002 | 集成测试 | P0 | ✅ 已完成 | [09-testing.md](./09-testing.md) |
| TEST-003 | Mock 设计 | P0 | ✅ 已完成 | [09-testing.md](./09-testing.md) |
| TEST-004 | CI/CD 集成 | P0 | ✅ 已完成 | [09-testing.md](./09-testing.md) |

---

## 部署

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| DEPLOY-001 | 构建流程 | P0 | ✅ 已完成 | [10-deployment.md](./10-deployment.md) |
| DEPLOY-002 | 多平台发布 | P0 | ✅ 已完成 | [10-deployment.md](./10-deployment.md) |
| DEPLOY-003 | 版本管理 | P0 | ✅ 已完成 | [10-deployment.md](./10-deployment.md) |

---

## 文档

| ID | 需求 | 优先级 | 状态 | 详细文档 |
|----|------|--------|------|----------|
| DOC-001 | COMMANDS.md 命令使用指南 | P0 | ✅ 已完成 | [COMMANDS.md](../docs/COMMANDS.md) |
| DOC-002 | PACKAGING.md 打包发布指南 | P0 | ✅ 已完成 | [PACKAGING.md](../docs/PACKAGING.md) |
| DOC-003 | SECURITY.md 安全策略 | P0 | ✅ 已完成 | [SECURITY.md](../SECURITY.md) |
| DOC-004 | LICENSE MIT许可证 | P0 | ✅ 已完成 | [LICENSE](../LICENSE) |
| DOC-005 | 仓库迁移到 gitcode-cli/cli | P0 | ✅ 已完成 | - |

---

## 状态图例

- 📋 待开发：需求已定义，等待开发
- 🚧 开发中：正在开发中
- ✅ 已完成：功能已实现并通过验收
- ⏸️ 暂停：开发暂停
- ❌ 取消：需求已取消

---

**最后更新**: 2026-03-23