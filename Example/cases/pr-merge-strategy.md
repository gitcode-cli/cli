---
title: PR 合并策略与清理
description: 使用 GitCode CLI 管理 PR 的 checkout 验证、合并策略（merge/squash/rebase）和过期 PR 清理
---

# PR 合并策略与清理

## 场景

openLiBing 发布平台当前有多个 PR，部分处于 open 状态较久、部分已合并但分支未删除。维护者需要：验证 PR 是否仍可合并、选择合适的合并策略、关闭过期 PR、清理已合并分支。这个案例展示 PR 生命周期后半段的管理操作。

## 推荐 skill

- `gitcode-pr` — 来自 [gitcode-cli/skills](https://gitcode.com/gitcode-cli/skills) 项目（`git@gitcode.com:gitcode-cli/skills.git`），可独立安装使用

## 适用人群

- 维护者（合并 PR、清理过期 PR）
- 发布负责人（选择合并策略、管理发布分支）
- DevOps（自动化 PR 合并流程）

## 可直接执行的 Prompt

```text
请使用 gitcode-pr skill，帮我管理 openLiBingNext/openlibing-platform-release 的 PR 队列。

请全程使用 `gitcode` 命令入口。写到远端前先 dry-run 或给我预览。

操作计划：
1. 列出所有 open PR，检查每个的状态：
   gitcode pr list -R openLiBingNext/openlibing-platform-release --state open --json

2. 对每个 open PR 检查：
   - 是否仍有 mergeable 标记？
   - 是否有冲突需要解决？
   - 最后更新时间多久了？

3. 对已完成的 PR，选择合适的合并策略：
   - 单次提交 → `--method squash`（保持 main 历史整洁）
   - 多个独立提交 → `--method merge`（保留提交历史）
   - 线性历史 → `--method rebase`

4. 对过期/不再需要的 PR，建议关闭：
   gitcode pr close <PR编号> -R openLiBingNext/openlibing-platform-release --yes --json

5. 合并后删除源分支：
   gitcode pr merge <PR编号> -R openLiBingNext/openlibing-platform-release --method squash --delete-branch --yes --json

请先输出分析报告和操作建议，等我确认后再执行远端写操作。
```

## 预期产出

- 当前 PR 队列状态报告
- 每个 PR 的推荐操作（合并/关闭/等更新）
- 合并策略选择建议
- 过期 PR 清理清单

## 价值

- 防止过期 PR 堆积，保持 PR 队列清晰
- 根据变更性质选择合适的合并策略，保持提交历史可读
- 减少"PR 已合并但分支忘了删"的仓库垃圾

## 复用方式

### 替换清单

| 占位符 | 案例值 | 替换为 |
|---|---|---|
| 仓库 | `openLiBingNext/openlibing-platform-release` | 目标仓库 |
| PR 编号 | 实际 open PR 编号 | 待处理的 PR |

### 适用场景

- 版本发布前整理 PR 队列
- 定期清理过期 PR
- 团队建立 PR 合并规范
- 不适合：受保护分支需要多人审批的场景（需先在 Web 端审批）

### 跨平台提醒

- `--yes` 跳过交互确认，CI 环境必须使用
- `pr merge` 可能因分支保护规则失败，需确认仓库设置

### 前置条件

- 对目标仓库有 PR 合并权限
- 了解仓库的分支保护规则
- （可选）安装 `gitcode-pr` skill

## 相关案例

- 前置：[从 Fork 分支创建发布平台 PR](./create-pr-from-fork.md) — PR 创建
- 前置：[评审已有 Tag 发布能力 PR](./review-pr.md) — 评审通过后才合并
- 关联：[发布 openLiBing 发布平台版本](./publish-release.md) — 合并后发布新版本

## 本次真实执行记录

本案例分析了 `openLiBingNext/openlibing-platform-release` 的 PR 队列：

- 分析时间：2026-05-26
- Open PR 数量：5 个（PR #1 到 #5，0 个已合并）
- 最近合并 PR：0 个（仓库尚无已合并的 PR）
- PR 状态快照：
  - PR #5 "docs: sync GitCode CLI example cases" — open, mergeable=true, 创建于 2026-05-26
  - PR #4 "feat(release): support existing tag in executeTag flow" — open, mergeable=true, 创建于 2026-05-19
  - PR #3 "test(attachment): 补充附件管理模块单元测试覆盖" — open, mergeable=true, 创建于 2026-05-18 (作者: kerer-sk)
  - PR #2 "refactor: optimize file transfer subsystem..." — open, **mergeable=false**, 创建于 2026-05-18
  - PR #1 "chore: clean stale TODO and outdated JavaDoc in HwCloudClient" — open, mergeable=true, 创建于 2026-05-18
- 合并策略建议：
  - PR #1（单次提交的文档清理）→ 推荐 `--method squash`
  - PR #2（多模块重构，但 mergeable=false 需先解决冲突）→ 解决冲突后视提交历史选择 `--method merge` 或 `squash`
  - PR #3（测试补充，836 行新增）→ 推荐 `--method squash`
  - PR #4（feature 实现，提交历史清晰）→ 可用 `--method merge`
  - PR #5（文档同步）→ 推荐 `--method squash`
- 过期 PR 识别：PR #1、#2、#3、#4 均已创建超过一周且无更新，需要确认作者是否仍计划推进；PR #2 存在合并冲突（mergeable=false），优先级最高需要处理

关键操作注意：合并前确认 `mergeable=true`，合并后使用 `--delete-branch` 清理源分支。`pr merge` 受分支保护规则限制，提前确认仓库设置。

![GitCode CLI PR merge evidence](assets/openlibing-pr-merge-evidence.svg)

复盘：PR 队列管理的痛点不是"合并不了"，而是"不知道哪些该合并、哪些该关闭"。定期（每两周或每个 milestone 结束前）用 `gitcode pr list --json` 生成队列报告，可以防止 PR 腐烂。
