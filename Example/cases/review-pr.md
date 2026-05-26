---
title: 评审已有 Tag 发布能力 PR
description: 使用 GitCode CLI 评审 openLiBing 发布平台 PR #4 的发布流程变更风险
---

# 评审已有 Tag 发布能力 PR

## 场景

`openLiBingNext/openlibing-platform-release` 的 PR #4 为发布流程增加“使用已有 Tag 发布”的能力，涉及 `ExecuteTagModel`、`OpenEulerJenkinsOperation`、`executeTag` 分支和 Jenkins 参数传递。此类变更容易引入权限、命令执行、tag/commit 校验和向后兼容风险，适合用结构化 PR 评审案例展示 GitCode CLI 的价值。

## 推荐 skill

- `gitcode-pr-review`
- 可辅助使用：`gitcode-review`

## 可直接执行的 Prompt

```text
请使用 gitcode-pr-review skill，对 openLiBingNext/openlibing-platform-release 的 PR #4 做一次工程评审。

请全程使用 `gitcode` 命令入口。重点看：
- `tagName` 新字段是否保持向后兼容；
- 既有自动创建 Tag 流程是否被影响；
- `git ls-remote` / `ProcessBuilder` 调用是否有命令注入、超时和错误处理风险；
- tag commitId 与产品包 commit.txt 的匹配校验是否覆盖异常路径；
- Jenkins `source=existingTag` 参数是否需要文档或配置同步；
- 测试是否覆盖 tag 匹配、不匹配、不存在、commit.txt 获取失败和 tagName 为空路径。

请先输出评审报告，不要直接 approve；如果需要在 PR 上发表评论，先给我预览。
```

## 预期产出

- 一份针对 PR #4 的结构化评审报告。
- 明确列出 blocker、建议项、测试缺口和是否建议合并。
- 可选生成可直接通过 `gitcode pr review --comment-file` 发布的评论正文。

## 价值

- 把发布链路变更从“看代码风格”提升到“验证发布安全与兼容性”。
- 对 Jenkins、Tag、commitId 校验这类高风险路径建立固定审查清单。
- 让维护者可以复用同一 prompt 评审后续发布流程 PR。

## 复用方式

复用到其他 PR 时，替换仓库和 PR 编号，并把重点检查项改为该 PR 涉及的业务风险即可。
