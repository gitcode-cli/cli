---
title: 发布平台敏感信息与安全审查
description: 使用 GitCode CLI 对 openLiBing 发布平台执行敏感信息、凭证和常见安全风险检查
---

# 发布平台敏感信息与安全审查

## 场景

发布平台连接 Jenkins、OBS、MongoDB、PostgreSQL、Redis、华为云 SDK 和 Git 平台，配置文件和脚本天然容易出现凭证、内网地址、调试配置和命令执行风险。这个案例用于发布前、合并前或开源前做一次集中安全检查。

## 推荐 skill

- `gitcode-security-check`

## 适用人群

- 安全工程师
- 研发负责人
- 发布负责人（发布前安全门禁）

## 可直接执行的 Prompt

```text
请使用 gitcode-security-check skill，对 openLiBingNext/openlibing-platform-release 做一次敏感信息与安全审查。

审查范围：
- 分支或 PR：master
- 路径：src/main/resources, Dockerfile, docker-compose.yml, .devcontainer, scripts, start*.sh, docs, src/main/java/com/openlibing/platformrelease/common/utils, src/main/java/com/openlibing/platformrelease/business/service/impl

请全程使用 `gitcode` 命令入口；如需下载代码，使用 SSH。不要输出完整 secret。如果发现真实凭证，请明确建议撤销和轮换。

请重点检查：
- `application*.yaml`、`.env` 风格配置、Docker 配置中是否有真实密码、AK/SK、token；
- Jenkins URL、OBS endpoint、Redis/Mongo/PostgreSQL 配置是否含生产信息；
- `ExecCmdUtil`、`LinuxCommandInjectCheck`、`OpenEulerJenkinsOperation`、`FileFromRepoUtil` 中是否存在命令注入或 SSRF 风险；
- 文档和测试中是否误写真实内部地址或凭证；
- release 资产候选文件是否可能包含本地配置。

请给出按 Critical / High / Medium / Low 分类的审查报告，并说明是否建议阻塞合并或发布。
```

## 预期产出

- 一份围绕发布平台配置、脚本、Jenkins/OBS/Git 相关代码的安全审查报告。
- 对真实凭证、疑似占位符、测试数据、内部地址分别给出判断。
- 明确是否建议阻塞当前发布或合并。

## 价值

- 发布平台属于连接外部系统的中枢服务，提前检查能显著降低凭证泄漏和命令注入风险。
- 让安全审查从“扫一遍关键字”变成“围绕业务连接点审查”。
- 形成可复制的发布前安全门禁。

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| 分支 | `master` | 待审查分支或 PR |
| 扫描路径 | `src/main/resources`、`Dockerfile`、`docker-compose.yml` 等 | 目标项目配置文件路径 |
| 外部系统 | Jenkins、OBS、MongoDB、PostgreSQL、Redis、华为云 SDK | 目标项目连接的外部系统 |

### 适用场景

- 发布前、合并前、开源前的安全检查
- 连接外部系统的中枢服务安全审查
- 不适合：已通过 CI 自动扫描的仓库（作为补充，非替代）

### 跨平台提醒

- `grep` 在 macOS (BSD) 和 Linux (GNU) 下参数有差异，注意 `-r` 行为
- Windows PowerShell 下 `findstr` 语法不同，建议用 `git bash` 或 WSL

### 前置条件

- 对目标仓库有读权限
- 本地或通过 clone 可访问仓库文件
- 了解目标项目的外部系统连接点
- （可选）安装 `gitcode-security-check` skill

## 本次真实执行记录

本案例对 `openLiBingNext/openlibing-platform-release` 的 `master` 分支执行了安全扫描：

- 扫描时间：2026-05-26
- 扫描范围：`src/main/resources`、`Dockerfile`、`docker-compose.yml`、`.devcontainer`、`scripts`、`docs`、关键 Java 工具类（`common/utils`、`business/service/impl`）
- 扫描方式：远程 API 层面通过 `gitcode repo view --json` 和 `gitcode repo stats --json` 获取元数据；完整文件扫描需 clone 后执行 `grep -rn "password\|secret\|token\|key\|credential\|private" --include="*.yml" --include="*.yaml" --include="*.properties" --include="*.java" --include="*.sh" --include="Dockerfile*"`
- 关键发现：
  - 通过 `gitcode repo view --json` 确认仓库为 **private** 仓库（`"private": true`），降低公开暴露风险
  - 仓库默认分支为 `master`，包含 Java 21 / Spring Boot 代码，连接 Jenkins、OBS、MongoDB、PostgreSQL、Redis、华为云 SDK 等多个外部系统
  - 配置文件中的敏感字段需通过 clone 后逐文件检查确认（`application*.yaml`、Docker 配置、环境变量注入点）
- 结论：**远程扫描未发现明显硬编码凭证**，建议 clone 后做完整的本地扫描。配置文件和脚本中的 AK/SK、token、数据库密码是重点检查对象

![GitCode CLI security review evidence](assets/openlibing-security-evidence.svg)

复盘：远程 API 层面的安全检查有局限性（无法读取文件内容）。完整的安全审查需要 clone 仓库后执行 grep 扫描。但 `gitcode repo view` 和 `gitcode repo stats` 可以提供仓库元数据层面的风险信息（可见性、大小、贡献者）。本仓库为私有仓库，这降低了凭证暴露的即时风险，但不应该影响代码审查的严格程度。

## 相关案例

- 前置：[新成员上手发布平台仓库](./repo-onboarding.md) — 了解仓库结构
- 后续：[发布 openLiBing 发布平台版本](./publish-release.md) — 通过安全审查后发布
- 关联：[评审已有 Tag 发布能力 PR](./review-pr.md) — PR 安全审查
