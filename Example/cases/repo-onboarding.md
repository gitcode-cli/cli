---
title: 新成员上手发布平台仓库
description: 使用 GitCode CLI 和本地检查快速生成 openLiBing 发布平台上手指南
---

# 新成员上手发布平台仓库

## 场景

新成员加入 openLiBing 发布平台团队，需要快速了解这个 Java 21 / Maven / Spring Boot 仓库如何启动、测试、访问 API 文档，以及如何围绕发布评审、Jenkins、OBS、漏洞扫描等模块贡献代码。

## 推荐 skill

- `gitcode-repo-onboarding`

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

复用到其他仓库时，替换仓库名、默认分支、启动命令、测试命令和关键模块即可。
