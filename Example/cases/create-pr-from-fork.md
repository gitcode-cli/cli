---
title: 从 Fork 分支创建发布平台 PR
description: 使用 GitCode CLI 从 fork 仓库向 openLiBing 发布平台主仓创建 Pull Request
---

# 从 Fork 分支创建发布平台 PR

## 场景

外部贡献者或团队成员没有 `openLiBingNext/openlibing-platform-release` 主仓写权限时，可以在 fork 仓库完成开发，再向主仓 `master` 分支提交 Pull Request。发布平台已有多个适合用 fork PR 承接的任务，例如 Issue #3 的附件管理测试覆盖、Issue #5 的发布结果追踪可靠性增强。

## 推荐 skill

- `gitcode-pr-create`
- 可辅助使用：`gitcode-repo`

## 可直接执行的 Prompt

```text
请使用 gitcode-pr-create skill，帮我从 fork 仓库向 openLiBing 发布平台主仓创建 Pull Request。

上下文：
- 上游仓库：openLiBingNext/openlibing-platform-release
- fork 仓库：afly-infra/openlibing-platform-release
- 工作分支：feat/release-result-failure-summary
- 目标分支：master

请全程使用 `gitcode` 命令入口；涉及代码传输默认使用 SSH。

我的变更说明：
- 关联 Issue：计划解决“发布结果中缺少制品级失败原因聚合视图”，并关联 Issue #5 的 release_result 可靠性增强上下文。
- 主要改动：新增按 reviewId 聚合 release_result 的查询能力；补充失败原因默认文案；增加 service 层单元测试。
- 验证计划：运行相关 service 测试和 `mvn test -Dtest=ReleaseDecisionServiceImplTest`；如改动 API，再补充接口文档。
- 风险：release_result 状态枚举和错误信息字段需要与现有发布流程保持兼容。

请先给出 PR 标题和正文预览，等我确认后再创建。
```

## 预期产出

- 从 `afly-infra/openlibing-platform-release` 的功能分支向 `openLiBingNext/openlibing-platform-release:master` 创建 PR。
- PR 描述中自然包含关联 Issue、发布平台业务价值、验证计划和兼容性风险。

## 价值

- 外部贡献者只需要 fork 权限和 SSH 配置，就能参与发布平台建设。
- PR 内容围绕发布平台真实模块展开，Reviewer 可以快速判断影响范围。
- 统一 `gitcode pr create --fork` 路径，避免手写 Web 表单遗漏测试和风险说明。

## 复用方式

复用时替换上游仓库、fork 仓库、分支名、关联 Issue 和验证计划即可。若目标项目默认分支不是 `master`，同步替换目标分支。
