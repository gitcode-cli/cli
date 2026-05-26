---
title: 向发布平台提交高质量 Issue
description: 使用 GitCode CLI 和 gitcode-issue-create skill 为 openLiBing 发布平台提交可执行 Issue
---

# 向发布平台提交高质量 Issue

## 场景

openLiBing 发布平台涉及发布评审、Jenkins 任务、OBS 制品下载、附件管理、漏洞扫描和发布结果追踪。用户经常会先发现一个“现象”，例如发布单卡在发布中、附件下载失败、已有 Tag 无法复用、某个制品缺少失败原因。这个案例展示如何把零散现象整理成维护者能直接分析和分派的 GitCode Issue。

## 推荐 skill

- `gitcode-issue-create`

## 适用人群

- 产品经理提交需求
- 测试人员提交缺陷
- 开源用户反馈问题
- AI 代理帮助用户整理 issue

## 可直接执行的 Prompt

```text
请使用 gitcode-issue-create skill，帮我向 openLiBingNext/openlibing-platform-release 提交一个高质量 Issue。

请全程使用 `gitcode` 命令入口；如果信息不足，先问我。

我的原始描述：
发布平台的发布决策当前已经有制品级 release_result 表，但用户在评审单详情页只能看到最终成功/失败，无法快速定位哪个软件包、哪个 Jenkins 阶段失败。希望在发布结果追踪中补充“制品级失败原因聚合视图”：

- 背景：Issue #5 已经讨论了异步发布可靠性和 release_result 预插入，本需求是在它基础上的用户可视化增强。
- 目标：按 reviewId 聚合 release_result，展示每个制品的状态、失败阶段、失败摘要、最后更新时间。
- 期望：后端提供查询接口，前端或调用方可以直接展示失败列表。
- 价值：发布负责人不用翻 Jenkins 日志和数据库就能判断失败原因。
- 验收：支持按 reviewId 查询；失败原因为空时给出默认文案；不影响现有发布流程；补充单元测试。

请先给出 issue 预览，等我确认后再创建。
```

## 预期产出

- 一个面向 `openLiBingNext/openlibing-platform-release` 的 Issue 草稿。
- 自动查重后识别是否应关联已有 Issue #5，而不是重复创建。
- 建议使用现有 `enhancement` 标签，必要时提示仓库缺少更细的 `scope/release`、`type/feature` 标签。

## 价值

- 把“发布失败看不懂”这类用户反馈变成可分派的后端增强任务。
- 让需求自然关联到已有异步发布可靠性 Issue #5，保留上下文。
- 让维护者直接看到影响范围、验收标准和测试要求，减少来回追问。

## 复用方式

复用到其他仓库时，保留“背景、目标、期望、价值、验收”结构，将仓库名、已有关联 issue 和业务模块替换为目标项目即可。
