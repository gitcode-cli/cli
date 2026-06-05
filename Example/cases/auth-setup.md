---
title: 多环境 GitCode CLI 认证配置
description: 在 CI、本地开发和容器三种环境下配置 GitCode CLI 认证，排查 token 来源优先级问题
---

# 多环境 GitCode CLI 认证配置

## 场景

团队在不同环境下使用 GitCode CLI 时，认证方式的差异经常导致"本地能跑、CI 报 401"。常见问题：环境变量 `GC_TOKEN` 和 `gitcode auth login` 的优先级、Windows 凭据存储差异、CI 无 TTY 环境下的非交互式登录。这个案例展示如何在三种典型环境下正确配置认证并排查问题。

## 推荐 skill

- `gitcode-auth` — 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 开发者（本地环境配置）
- DevOps（CI/CD 流水线配置）
- 平台团队（容器化环境配置）

## 可直接执行的 Prompt

```text
请使用 gitcode-auth skill，帮我排查当前环境的 GitCode CLI 认证配置。

请全程使用 `gitcode` 命令入口。请按顺序：

1. 先确认 CLI 版本和认证状态：
   gitcode version
   gitcode auth status --json

2. 检查 token 来源（环境变量 vs 存储凭据）：
   - 是否设置了 GC_TOKEN 或 GITCODE_TOKEN 环境变量？
   - 如果两者都设置了，GC_TOKEN 优先级更高，这可能不是期望的

3. 根据环境类型给出配置建议：
   - 本地开发：gitcode auth login --with-token --git-protocol ssh
   - CI 环境：export GC_TOKEN="<token>" （不存储在配置文件中）
   - 容器环境：通过 secret 注入 GC_TOKEN 环境变量

4. 验证认证是否生效：
   gitcode repo view openLiBingNext/openlibing-platform-release --json

请输出当前认证状态、token 来源、潜在问题和修复建议。不要输出完整 token 值。
```

## 预期产出

- 当前环境的认证状态报告
- token 来源分析（环境变量 vs 存储凭据）
- 环境适配建议（本地/CI/容器）
- 认证修复操作指南

## 价值

- 避免"本地能跑 CI 报错"的认证配置不一致
- 明确 GC_TOKEN 优先级规则，防止环境变量污染
- 为团队提供标准化的认证配置模板

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 测试仓库 | `openLiBingNext/openlibing-platform-release` | 你的私有仓库 |
| SSH 偏好 | `--git-protocol ssh` | 可改为 `--git-protocol https` |

### 适用场景

- 首次配置 GitCode CLI
- CI/CD 流水线中集成 GitCode CLI
- 认证失败排查（401、403 错误）
- 不适合：OAuth/SSO 相关的认证问题（需平台侧配置）

### 跨平台提醒

- **Windows**：`gc` 是 `Get-Content` 别名，统一使用 `gitcode`
- **Linux/macOS**：`GC_TOKEN` 环境变量优先级高于 `GITCODE_TOKEN`
- **CI 环境**：使用 `echo "$TOKEN" | gitcode auth login --with-token` 需确保无 TTY

### 前置条件

- GitCode CLI 已安装
- 已获取 GitCode 个人访问 token
- 对测试仓库有读权限

## 相关案例

- 后续：[对发布平台仓库做 CLI 冒烟验证](./regression-after-install.md) — 认证成功后验证 CLI 能力
- 关联：[新成员上手发布平台仓库](./repo-onboarding.md) — 认证是 clone 和 onboarding 的前提

## 本次真实执行记录

本案例使用当前版本验证了 GitCode CLI 的认证配置：

- CLI 版本：gitcode version 0.5.8 (commit: 1bea9ac, built: 2026-05-26)
- 登录用户：aflyingto
- 认证方式：GC_TOKEN（环境变量）
- Git 协议：https
- 验证时间：2026-05-26
- 验证仓库：`gitcode repo view openLiBingNext/openlibing-platform-release --json` 返回正常

关键发现：当前环境使用 GC_TOKEN 环境变量认证。`GC_TOKEN` 优先级高于 `GITCODE_TOKEN`，高于 `gitcode auth login` 存储的凭据。在同时设置多个认证源时，`auth status --json` 的 token 来源字段（显示为 `GC_TOKEN`）是排查关键。注意 `auth status --hostname gitcode.com --json` 返回 `logged_in: false`，说明带 `--hostname` 参数时仅检查存储凭据，不会回退到环境变量检测。

![GitCode CLI auth evidence](assets/openlibing-auth-evidence.svg)

复盘：认证问题是"能跑"和"不能跑"之间的第一道门槛。`gitcode auth status --json` 应作为所有自动化脚本的第一步健康检查。
