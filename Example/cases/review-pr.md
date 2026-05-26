---
title: 对 Pull Request 做工程评审
description: 使用 GitCode CLI 查看 PR、diff、评论并输出结构化工程评审结论
---

# 对 Pull Request 做工程评审

## 场景

维护者或 Reviewer 需要对一个 GitCode Pull Request 做独立工程评审，重点关注缺陷、回归、安全风险、测试缺口和合并阻塞点。

## 推荐 skill

- `gitcode-pr-review`
- 可辅助使用：`gitcode-review`

## 可直接执行的 Prompt

```text
请使用 GitCode CLI 对 <owner/repo> 的 PR #<number> 做一次工程评审。

要求：
1. 全程使用 `gitcode` 命令，不使用 `gc`。
2. 优先使用 `gitcode-pr-review` skill；如果未安装该 skill，请按同等流程执行。
3. 先读取 PR 元数据、评论和 diff：
   - gitcode pr view <number> -R <owner/repo> --comments --json
   - gitcode pr comments <number> -R <owner/repo> --json
   - gitcode pr diff <number> -R <owner/repo>
4. 如需本地验证，使用：
   - gitcode pr checkout <number> -R <owner/repo>
5. 评审输出必须优先列出问题，按严重程度排序：
   - Critical / High / Medium / Low / Info
6. 重点检查：
   - 行为回归
   - 安全风险
   - 缺失测试
   - 文档同步
   - API/CLI 向后兼容
   - 是否存在敏感信息
7. 如需要发表评论，先生成 review-report.md，再使用：
   - gitcode pr review <number> -R <owner/repo> --comment-file review-report.md
8. 只有在没有 blocker 且验证充分时，才建议 approve；不要默认批准。

输出：
- Findings
- Open Questions
- Verification
- Residual Risk
- 是否建议合并
```

## 预期产出

- 一份结构化 PR 评审报告。
- 可选的 GitCode PR review 评论。

## 价值

- 让评审从“泛泛看过”变成可审计的工程判断。
- 帮助团队统一 blocker 和非 blocker 的分级。
- 减少漏测、漏文档、漏安全风险。

## 复用方式

替换 `<owner/repo>` 和 `<number>` 即可用于任意 GitCode PR。
