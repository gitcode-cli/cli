# gc-issue

使用 `gc` 完成 GitCode issue 相关操作。

## 触发场景

- 创建 issue
- 查看或列出 issue
- 编辑 issue
- 关闭 / 重开 issue
- 添加或查看 issue 评论

## 常用命令

```bash
# 创建 issue
gc issue create -R owner/repo --title "Bug report" --body "Description"

# 列出 issue
gc issue list -R owner/repo --state open

# 查看 issue
gc issue view 123 -R owner/repo

# 编辑 issue
gc issue edit 123 -R owner/repo --title "New title"

# 关闭 / 重开 issue
gc issue close 123 -R owner/repo
gc issue reopen 123 -R owner/repo

# 评论
gc issue comment 123 -R owner/repo --body "Comment text"
gc issue comments 123 -R owner/repo --limit 10
gc issue comment edit 166061383 -R owner/repo --body "Updated comment"
```

## 使用约束

- 未传 `-R` 时，是否支持从当前仓库自动推断，取决于具体命令和当前目录
- 修改 issue 前，建议先查看当前内容，避免覆盖已有信息
- 批量筛选和复杂查询时，优先先跑 `gc issue list` 验证过滤条件
