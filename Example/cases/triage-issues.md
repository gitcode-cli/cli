---
title: 批量整理 Issue 队列
description: 使用 GitCode CLI 对 Issue 队列进行分类、优先级排序、补标签和重复项识别
---

# 批量整理 Issue 队列

## 场景

项目维护者需要整理长期积压的 Issue，识别重复问题、缺信息问题、高优先级问题，并补齐标签和评论。

## 推荐 skill

- `gitcode-issue-triage`

## 可直接执行的 Prompt

```text
请使用 GitCode CLI 帮我整理 <owner/repo> 的 Issue 队列。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-issue-triage` skill；如果未安装该 skill，请按同等流程执行。
3. 先读取仓库标签和里程碑：
   - gitcode label list -R <owner/repo> --json
   - gitcode milestone list -R <owner/repo> --json
4. 拉取待整理 Issue：
   - gitcode issue list -R <owner/repo> --state open --limit <limit> --json
5. 对每个 issue 判断：
   - 类型：bug / feature / enhancement / docs / question
   - 优先级：critical / high / medium / low
   - 状态：needs-info / ready / duplicate / blocked
   - 是否缺少复现、验收标准或版本信息
6. 不要批量修改前直接执行远端写操作。先输出计划表让我确认。
7. 我确认后，再逐条执行：
   - gitcode issue label <number> -R <owner/repo> --add <labels>
   - gitcode issue comment <number> -R <owner/repo> --body-file <comment-file> --json
8. 如需关闭重复 issue，必须先给出重复依据并等待确认。

输出：
- Issue triage 表格
- 推荐标签变更
- 推荐评论
- 需要用户确认的批量操作清单
```

## 预期产出

- 一份 issue 队列整理计划。
- 经确认后的标签和评论变更。
- 重复项和缺信息项列表。

## 价值

- 让维护者快速掌握 backlog 状态。
- 降低 issue 队列长期失控的风险。
- 为版本规划和人力分配提供依据。

## 复用方式

替换仓库和 `limit` 即可用于不同项目的定期 issue 运营。
