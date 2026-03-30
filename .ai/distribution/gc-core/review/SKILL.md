# gc-review

使用 `gc` 完成 GitCode Pull Request 审查相关操作。

## 触发场景

- 对 PR 留评论
- 批准 PR
- 查看 PR 评论
- 回复 PR 评论

## 常用命令

```bash
# 查看 PR
gc pr view 123 -R owner/repo

# 查看评论
gc pr comments 123 -R owner/repo

# 回复评论
gc pr reply 123 -R owner/repo --discussion <discussion_id> --body "Reply"

# 评论审查
gc pr review 123 -R owner/repo --comment "Review comment"

# 批准
gc pr review 123 -R owner/repo --approve
gc pr review 123 -R owner/repo --approve --comment "LGTM"
```

## 已知限制

- `gc pr review --request` 当前会明确提示 GitCode API 暂不支持该动作
- PR 评论的 resolved / unresolved 状态仍需在 Web UI 中处理

## 使用建议

- 先查看 PR 和 diff，再发表评论
- 审查结论应明确区分 blocker 和建议项
- 平台限制应在评论中如实说明，不要假设 CLI 支持未实现能力
