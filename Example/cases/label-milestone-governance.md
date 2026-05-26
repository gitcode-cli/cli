---
title: 标签体系与里程碑治理
description: 使用 GitCode CLI 为发布平台仓库建立 type/scope/priority/status 标签体系和版本里程碑
---

# 标签体系与里程碑治理

## 场景

openLiBing 发布平台仓库当前标签较少，缺少 `scope/*`（模块范围）、`type/*`（工作类型）、`priority/*`（优先级）分类，Issue 和 PR 难以按模块筛选。同时没有设置里程碑，版本规划无锚点。这个案例展示如何用 GitCode CLI 建立一套轻量但完备的标签体系和里程碑。

## 推荐 skill

- `gitcode-label-milestone` — 标签和里程碑的创建、编辑、删除、关联

## 适用人群

- 维护者、项目负责人
- Scrum Master、项目经理
- 开源社区经理

## 可直接执行的 Prompt

```text
请使用 gitcode-label-milestone skill，帮我为 openLiBingNext/openlibing-platform-release 建立标签体系和版本里程碑。

请全程使用 `gitcode` 命令入口。所有创建操作先 dry-run 预览，等我确认后再执行。

步骤：
1. 先查看当前有哪些标签和里程碑：
   gitcode label list -R openLiBingNext/openlibing-platform-release --json
   gitcode milestone list -R openLiBingNext/openlibing-platform-release --json

2. 建议补充以下标签（按模块和类型）：
   - scope: scope/release, scope/attachment, scope/file-transfer, scope/jenkins, scope/obs
   - type: type/feature, type/bug, type/test, type/docs, type/chore, type/security
   - priority: priority/critical, priority/high, priority/medium, priority/low

3. 建议创建里程碑：
   - "v1.0 发布平台基础能力"（当前已合并的功能）
   - "v1.1 发布结果追踪增强"（Issue #5 及相关）
   - "v1.2 附件管理与测试覆盖"（Issue #3 及相关）

4. 将现有 open issue 关联到合适的里程碑。

请先输出标签和里程碑的规划方案，不要直接创建。
```

## 预期产出

- 当前标签体系分析报告
- 推荐的标签创建清单（含颜色和描述）
- 版本里程碑规划
- Issue 到里程碑的关联建议

## 价值

- 让 Issue/PR 可按模块、类型、优先级筛选
- 为版本规划和发布说明提供结构化元数据
- 防止标签膨胀（只补必要的，不创建用不到的）

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| 模块范围 | release, attachment, file-transfer, jenkins, obs | 目标项目的模块列表 |
| 里程碑名 | v1.0, v1.1, v1.2 | 你的版本号 |

### 适用场景

- 新仓库建立标签体系
- 版本规划前设立里程碑
- 标签混乱需要治理
- 不适合：已有成熟标签体系和 CI 自动标签的仓库（做增量补充即可）

### 跨平台提醒

- 标签颜色使用 hex 格式 `#ff0000`，不带 alpha 通道
- 里程碑日期使用 `YYYY-MM-DD` 格式

### 前置条件

- 对目标仓库有管理权限
- 了解目标项目的模块划分和版本规划
- （可选）安装 `gitcode-label-milestone` skill

## 相关案例

- 前置：[整理发布平台 Issue 队列](./triage-issues.md) — 了解当前 Issue 状态后再建标签
- 后续：[发布 openLiBing 发布平台版本](./publish-release.md) — 里程碑关联的版本发布
- 关联：[Issue 生命周期管理](./issue-lifecycle.md) — 标签在 Issue 流转中的应用

## 本次真实执行记录

本案例分析了 `openLiBingNext/openlibing-platform-release` 的标签和里程碑现状：

- 分析时间：2026-05-26
- 当前标签：5 个 — `enhancement`、`scope/common`、`scope/docs`、`type/chore`、`type/docs`
- 当前里程碑：0 个（无任何里程碑）
- 缺失的标签维度：
  - 无 `priority/*` 标签（无法标记紧急程度）
  - 无 `status/*` 标签（无法标记工作流状态）
  - scope 标签覆盖不全（仅有 common 和 docs，缺 release、attachment、file-transfer、jenkins、obs）
  - type 标签覆盖不全（仅有 chore 和 docs，缺 feature、bug、test、security）
- 建议新增标签数量：约 12 个（scope 5 + type 4 + priority 4，仅补必要项，status 标签另议）
- 建议新增里程碑：3 个（v1.0 基础能力、v1.1 发布追踪、v1.2 附件测试）

重要提示：创建标签前先 `--dry-run`，确认颜色和名称不与已有标签冲突。使用 `gitcode label create` 逐个创建，不要批量操作。

![GitCode CLI label milestone evidence](assets/openlibing-label-milestone-evidence.svg)

复盘：标签治理的原则是"只补必要的"。标签太少无法分类，太多变成噪音。以仓库当前 5 个标签为起点，补 12 个分到 4 个维度（type/scope/priority/status），总数控制在 15-20 之间是合理的。
