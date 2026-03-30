# Issue 命令 (issue)

## issue create - 创建 Issue

```bash
# 创建 Issue
gc issue create -R infra-test/gctest1 --title "Bug: Something wrong" --body "Description here"

# 创建 Issue 并添加标签
gc issue create -R infra-test/gctest1 --title "Feature request" --body "Description" --label enhancement

# 指定受理人
gc issue create -R infra-test/gctest1 --title "Task" --body "Description" --assignee username
```

## issue list - 列出 Issues

```bash
# 列出所有开放的 Issues
gc issue list -R infra-test/gctest1

# 只列出已关闭的 Issues
gc issue list -R infra-test/gctest1 --state closed

# 列出所有状态的 Issues
gc issue list -R infra-test/gctest1 --state all

# 按标签筛选
gc issue list -R infra-test/gctest1 --label bug,enhancement

# 限制数量
gc issue list -R infra-test/gctest1 --limit 20

# 按里程碑筛选
gc issue list -R infra-test/gctest1 --milestone "v1.0"

# 按受理人筛选
gc issue list -R infra-test/gctest1 --assignee username

# 按创建者筛选
gc issue list -R infra-test/gctest1 --creator username

# 按更新时间排序
gc issue list -R infra-test/gctest1 --sort updated --direction desc

# 按创建时间筛选
gc issue list -R infra-test/gctest1 --created-after "2024-01-01"
gc issue list -R infra-test/gctest1 --created-before "2024-12-31"

# 按更新时间筛选
gc issue list -R infra-test/gctest1 --updated-after "2024-01-01"

# 关键字搜索
gc issue list -R infra-test/gctest1 --search "bug"

# 组合使用
gc issue list -R infra-test/gctest1 --state open --milestone "v1.0" --sort updated
```

## issue view - 查看 Issue

```bash
# 查看 Issue 详情
gc issue view 1 -R infra-test/gctest1

# 查看评论
gc issue view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc issue view 1 -R infra-test/gctest1 --web
```

## issue close - 关闭 Issue

```bash
gc issue close 1 -R infra-test/gctest1
```

## issue edit - 编辑 Issue

```bash
# 修改标题
gc issue edit 1 --title "New title" -R infra-test/gctest1

# 修改描述
gc issue edit 1 --body "New description" -R infra-test/gctest1

# 修改状态（close/reopen）
gc issue edit 1 --state close -R infra-test/gctest1
gc issue edit 1 --state reopen -R infra-test/gctest1

# 指派负责人
gc issue edit 1 --assignee username -R infra-test/gctest1
gc issue edit 1 --assignee user1 --assignee user2 -R infra-test/gctest1

# 设置标签
gc issue edit 1 --label bug,enhancement -R infra-test/gctest1

# 设置里程碑
gc issue edit 1 --milestone 5 -R infra-test/gctest1

# 设置为私有 Issue
gc issue edit 1 --security-hole -R infra-test/gctest1

# 组合使用
gc issue edit 1 --title "Bug fix" --assignee username --label bug --milestone 1 -R infra-test/gctest1
```

## issue reopen - 重开 Issue

```bash
gc issue reopen 1 -R infra-test/gctest1
```

## issue comment - 添加评论

```bash
# 添加评论
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment"

# 从文件读取评论内容
gc issue comment 1 -R infra-test/gctest1 --body-file comment.txt

# 从 stdin 读取评论内容
echo "Comment from stdin" | gc issue comment 1 -R infra-test/gctest1 --body-file -
```

## issue comment edit - 编辑 Issue 评论

```bash
# 按参数编辑评论
gc issue comment edit 166061383 -R infra-test/gctest1 --body "Updated comment"

# 按 --id 编辑评论
gc issue comment edit --id 166061383 -R infra-test/gctest1 --body "Updated comment"

# 从文件读取新内容
gc issue comment edit 166061383 -R infra-test/gctest1 --body-file comment.md
```

## issue comments - 列出 Issue 评论

```bash
# 列出评论
gc issue comments 1 -R infra-test/gctest1
gc issue comments 1

# 限制返回数量
gc issue comments 1 -R infra-test/gctest1 --limit 10

# 倒序排列
gc issue comments 1 -R infra-test/gctest1 --order desc

# 按更新时间筛选
gc issue comments 1 -R infra-test/gctest1 --since "2024-01-01T00:00:00+08:00"

# JSON 输出
gc issue comments 1 -R infra-test/gctest1 --json
```

## issue label - 管理 Issue 标签

```bash
# 添加标签
gc issue label 1 --add bug,enhancement -R infra-test/gctest1

# 移除标签
gc issue label 1 --remove bug -R infra-test/gctest1

# 列出标签
gc issue label 1 --list -R infra-test/gctest1
```

## issue prs - 查看 Issue 关联的 PRs

```bash
# 查看 Issue 关联的 Pull Requests
gc issue prs 123 -R infra-test/gctest1

# 获取增强信息（包含可合并状态）
gc issue prs 123 --mode 1 -R infra-test/gctest1
```
