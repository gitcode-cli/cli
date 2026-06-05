---
title: 新成员上手发布平台仓库
description: 使用 GitCode CLI 和本地检查快速生成 openLiBing 发布平台上手指南
---

# 新成员上手发布平台仓库

## 场景

新成员加入 openLiBing 发布平台团队，需要快速了解这个 Java 21 / Maven / Spring Boot 仓库如何启动、测试、访问 API 文档，以及如何围绕发布评审、Jenkins、OBS、漏洞扫描等模块贡献代码。

## 推荐 skill

- `gitcode-repo-onboarding` — 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 新加入团队的开发者
- 外部贡献者
- 售前/交付团队（了解项目结构）

## 可直接执行的 Prompt

```text
请使用 gitcode-repo-onboarding skill，帮我为 openLiBingNext/openlibing-platform-release 生成一份新成员上手指南。

目标分支：master

请全程使用 `gitcode` 命令入口；如需下载代码，使用 SSH。不要编造构建和测试命令，只基于仓库真实文件总结。

请输出 Markdown 格式的 onboarding-guide.md 内容，重点覆盖：
- 项目定位：openLiBing 发布平台；
- 技术栈：Java 21、Maven、Spring Boot、Docker Compose、Dev Container；
- 启动方式：`./start-local.sh`、`mvn spring-boot:run -Dspring-boot.run.profiles=local`；
- 测试方式：`mvn test`、`mvn verify -Pintegration-test`；
- 访问地址：http://localhost:8085、`/actuator/health`、`/swagger-ui.html`；
- 关键模块：发布评审、Jenkins 集成、OBS 信息管理、文件下载、漏洞扫描、附件管理；
- 建议第一个任务：Issue #1 或 Issue #3。
```

## 预期产出

- 一份针对发布平台仓库的 onboarding-guide.md。
- 新成员能直接知道如何启动依赖、运行应用、访问 Swagger、跑测试和选择第一个任务。

## 价值

- 将 README、pom.xml、Docker 配置、docs 和现有 issue 中分散的信息汇总成单页指南。
- 降低 Java 服务本地启动和依赖服务准备的沟通成本。
- 帮助新人从小任务 Issue #1/#3 开始，而不是直接碰发布链路核心逻辑。

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| 技术栈 | Java 21 / Maven / Spring Boot | 目标项目技术栈 |
| 构建命令 | `mvn test` / `mvn spring-boot:run` | 目标项目构建命令 |
| 推荐首个任务 | Issue #1 或 #3 | 目标仓库适合新手的 issue |

### 适用场景

- 新成员入职、外部贡献者首次接触仓库
- 售前/交付团队需要快速了解项目结构和构建方式
- 不适合：已经很熟悉的仓库（直接读 README 即可）

### 跨平台提醒

- Windows 下 clone 使用 SSH 需先配置 OpenSSH 客户端
- 构建命令因操作系统不同可能需要调整（如 `mvnw` vs `mvn`）

### 前置条件

- `gitcode auth status` 确认已登录且有仓库读权限
- SSH 已配置（`ssh -T git@gitcode.com`）
- （可选）安装 `gitcode-repo-onboarding` skill

## 本次真实执行记录

本案例使用当前版本对 `openLiBingNext/openlibing-platform-release` 执行了真实 onboarding 检查链：

- `gitcode repo view --json`：仓库为 **private**（Java 21 / Maven / Spring Boot），默认分支 `master`，SSH 地址 `git@gitcode.com:openLiBingNext/openlibing-platform-release.git`，5 个开放 Issue
- `gitcode repo stats --branch master --json`：统计接口可用，返回了提交者列表和变更统计（私有仓库限制部分字段匿名化）
- `gitcode issue list --state open --limit 5 --json`：5 个开放 Issue，适合新人起步的任务为 **Issue #1**（chore: 清理过期 TODO/JavaDoc，`type/chore` + `scope/common`）和 **Issue #3**（补充附件管理模块单元测试覆盖）
- `gitcode label list --json`：5 个标签 -- `enhancement`、`scope/common`、`scope/docs`、`type/chore`、`type/docs`

![GitCode CLI onboarding evidence](assets/openlibing-onboarding-evidence.svg)

复盘：onboarding 不只是汇总 README 信息，需要实际执行命令验证仓库可访问、构建工具可识别。本仓库缺少 CONTRIBUTING.md，新成员 guide 中需标注并建议补充。Issue #1（清理过期注释）和 Issue #3（补充测试）是典型的新人入门任务，改动范围小、验收标准明确。

## 相关案例

- 后续：[向发布平台提交高质量 Issue](./create-issue.md) — 了解仓库后开始贡献
- 关联：[多环境认证配置](./auth-setup.md) — 确保 clone 和代码下载的认证配置正确
