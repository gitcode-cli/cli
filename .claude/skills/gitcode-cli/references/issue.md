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

# 按标签筛选
gc issue list -R infra-test/gctest1 --label bug

# 限制数量
gc issue list -R infra-test/gctest1 --limit 20
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
gc issue comment 1 -R infra-test/gctest1 --body "This is a comment"
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