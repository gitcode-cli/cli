# 流程状态更新 Checklist

本文档提供 Issue/PR 状态标签更新的简明 checklist，确保流程推进时状态标签同步更新。

## Issue 状态更新时机

| 步骤 | 操作命令 | 必须添加 | 必须移除 |
|------|----------|----------|----------|
| 创建 Issue 后 | `gc issue label <n> --add type/bug,status/triage,scope/xxx -R owner/repo` | `type/*`, `status/triage`, `scope/*` | - |
| 验证完成后 | `gc issue label <n> --add status/verified -R owner/repo` | `status/verified` | `status/triage` |
| 开始开发时 | `gc issue label <n> --add status/in-progress -R owner/repo` | `status/in-progress` | - |
| PR 合入后 | `gc issue label <n> --add status/merged -R owner/repo && gc issue close <n> -R owner/repo` | `status/merged` | `status/in-progress`, `status/triage` |

## PR 状态更新时机

| 步骤 | 操作命令 | 必须添加 | 必须移除 |
|------|----------|----------|----------|
| 创建 PR 后 | `gc pr edit <n> -R owner/repo --labels status/draft,type/*,risk/*,scope/*` | `status/draft`, `type/*`, `risk/*`, `scope/*` | - |
| 自检完成后 | `gc pr edit <n> -R owner/repo --labels status/self-checked` | `status/self-checked` | `status/draft` |
| 进入评审前 | `gc pr edit <n> -R owner/repo --labels status/ready-for-review` | `status/ready-for-review` | `status/self-checked` |
| 评审通过后 | `gc pr edit <n> -R owner/repo --labels status/approved` | `status/approved` | `status/ready-for-review` |
| 合入主干后 | `gc pr edit <n> -R owner/repo --labels status/merged` | `status/merged` | `status/approved` |

## Label 维度要求

### Issue 必须包含

- **类型**: `type/bug`, `type/feature`, `type/docs`, `type/refactor`
- **状态**: `status/triage`, `status/verified`, `status/in-progress`, `status/merged`
- **范围**: `scope/auth`, `scope/repo`, `scope/issue`, `scope/pr`, `scope/release`, `scope/docs`, `scope/testing`
- **风险**: `risk/low`, `risk/medium`, `risk/high` (可选，验证后添加)

### PR 必须包含

- **类型**: `type/bug`, `type/feature`, `type/docs`, `type/refactor`
- **状态**: `status/draft`, `status/self-checked`, `status/ready-for-review`, `status/approved`, `status/merged`
- **范围**: `scope/*` (与 Issue 对应)
- **风险**: `risk/low`, `risk/medium`, `risk/high` (自检时添加)

## 快速参考命令

```bash
# Issue 状态推进
gc issue label 123 --add status/verified --remove status/triage -R owner/repo
gc issue label 123 --add status/in-progress -R owner/repo
gc issue label 123 --add status/merged --remove status/in-progress,status/triage -R owner/repo

# PR 状态推进
gc pr edit 456 -R owner/repo --labels status/draft,type/feature,risk/low,scope/pr
gc pr edit 456 -R owner/repo --labels status/self-checked
gc pr edit 456 -R owner/repo --labels status/ready-for-review
gc pr edit 456 -R owner/repo --labels status/approved
gc pr edit 456 -R owner/repo --labels status/merged
```

---

**重要**: 状态标签更新是流程关键步骤，不得跳过或延迟执行。
