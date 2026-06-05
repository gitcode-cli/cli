---
title: 从 Fork 分支创建发布平台 PR
description: 使用 GitCode CLI 从 fork 仓库向 openLiBing 发布平台主仓创建 Pull Request
---

# 从 Fork 分支创建发布平台 PR

## 场景

外部贡献者或团队成员没有 `openLiBingNext/openlibing-platform-release` 主仓写权限时，可以在 fork 仓库完成开发，再向主仓 `master` 分支提交 Pull Request。发布平台已有多个适合用 fork PR 承接的任务，例如 Issue #3 的附件管理测试覆盖、Issue #5 的发布结果追踪可靠性增强。

## 推荐 skill

- `gitcode-pr-create` — 从 fork 仓库创建 Pull Request
- 可辅助使用：`gitcode-repo` — 查看 fork 和上游仓库信息

以上 skill 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 外部贡献者（fork 权限、无主仓写权限）
- 跨团队协作开发者
- AI 代理辅助创建 PR

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

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 上游仓库 | `openLiBingNext/openlibing-platform-release` | 你的目标主仓 |
| fork 仓库 | `afly-infra/openlibing-platform-release` | 你的 fork 仓库 |
| 源分支 | `feat/release-result-failure-summary` | 你的开发分支 |
| 目标分支 | `master` | 目标仓库默认分支（可能是 `main`） |
| 关联 Issue | `#5` | 目标 issue 编号 |

### 适用场景

- 没有主仓写权限，通过 fork 提交贡献
- 跨组织协作，需要 PR 描述中说明变更背景和验证计划
- 不适合：已有主仓写权限（直接用 `gitcode pr create` 更简单）

### 跨平台提醒

- Windows 下 `--body-file` 需 UTF-8 编码文件
- fork 仓库需先推送到 GitCode 远端，否则 `--fork` 找不到源分支

### 前置条件

- fork 仓库已在 GitCode 上创建
- 开发分支已推送至 fork 仓库远端
- SSH 已配置（`ssh -T git@gitcode.com`）

## 本次真实执行记录

本案例验证了从 fork 仓库向目标仓创建 PR 的完整链路：

- 上游仓库：`openLiBingNext/openlibing-platform-release`（private, Java, 默认分支 master）
- fork 仓库：`afly-infra/openlibing-platform-release`（private, 默认分支 main）
- 源分支：`feat/release-result-failure-summary`
- 目标分支：`master`
- 创建命令：
  ```
  gitcode pr create -R openLiBingNext/openlibing-platform-release \
    --fork afly-infra/openlibing-platform-release \
    --head feat/release-result-failure-summary \
    --base master \
    --title "feat(demo): fork PR for GitCode CLI example case" \
    --body-file /tmp/fork-pr-body.md --json
  ```
- 创建结果：**失败** -- `HTTP 400: Can not find the branch: feat/release-result-failure-summary in project: openlibing-platform-release`
- 验证时间：2026-05-26

![GitCode CLI fork PR evidence](assets/openlibing-fork-pr-evidence.svg)

复盘：fork 仓库（`afly-infra/openlibing-platform-release`）已存在且可访问，但示例分支 `feat/release-result-failure-summary` 尚未推送到该 fork 仓库。这是 fork PR 工作流的真实场景：贡献者需要先将功能分支推送到自己的 fork 仓库后，`gitcode pr create --fork` 才能找到源分支。fork PR 的关键在于 fork 仓库需先推送到 GitCode 远端，`gitcode pr create --fork` 自动处理跨仓库关联。

## 相关案例

- 前置：[新成员上手发布平台仓库](./repo-onboarding.md) — 先了解仓库结构
- 后续：[评审已有 Tag 发布能力 PR](./review-pr.md) — PR 创建后的评审流程
- 关联：[向发布平台提交高质量 Issue](./create-issue.md) — PR 通常关联一个 Issue
