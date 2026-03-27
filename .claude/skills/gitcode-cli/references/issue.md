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