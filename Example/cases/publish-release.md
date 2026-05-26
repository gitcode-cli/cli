---
title: 发布 openLiBing 发布平台版本
description: 使用 GitCode CLI 为 openLiBing 发布平台生成发布说明、创建 release 并上传资产
---

# 发布 openLiBing 发布平台版本

## 场景

`openLiBingNext/openlibing-platform-release` 当前尚未创建 release。维护者可以在合并一组发布平台增强后，基于 `master` 分支发布第一个内部版本，归档 Java 构建产物、发布说明和验证记录。

## 推荐 skill

- `gitcode-release-helper`
- 可辅助使用：`gitcode-release`

## 可直接执行的 Prompt

```text
请使用 gitcode-release-helper skill，帮我为 openLiBingNext/openlibing-platform-release 发布版本 v1.0.0-platform-release。

请全程使用 `gitcode` 命令入口；发布前先给我 release notes、资产清单和验证计划预览，等我确认后再创建 release。

输入：
- previous_tag: 无历史 release，请基于 master 当前提交生成首个版本说明
- version: v1.0.0-platform-release
- target_branch: master
- asset_files: target/openlibing-platform-release-1.0.jar, docs/API接口文档.md, docs/敏感信息配置说明.md
- 本版本重点变化：
  - 发布平台 Java 21 / Spring Boot 应用基础能力
  - 发布评审、Jenkins 集成、OBS 信息管理、漏洞扫描、文件下载能力
  - 本地 Docker Compose / Dev Container 开发方式
  - 当前 open issues 中的后续规划：附件管理测试覆盖、已有 Tag 发布、发布结果追踪可靠性
```

## 预期产出

- 一个面向 `openLiBingNext/openlibing-platform-release` 的 GitCode Release 预案。
- 清晰的首版 release notes，包括已具备能力、验证方式、已知后续事项。
- 确认后可创建 release 并上传 jar 与文档资产。

## 价值

- 将“代码仓已经可运行”转化为可对外交付的版本说明。
- 让发布负责人在创建 release 前先看到资产清单、验证计划和安全检查点。
- 可作为发布平台自身版本发布的示范，也能复用于其他 Java 服务。

## 复用方式

复用时替换版本号、目标分支、构建产物和本版本重点变化。若仓库已有 release，应将 `previous_tag` 改为上一个真实 tag。
