# Pull Request 命令 (pr)

## pr create - 创建 PR

```bash
# 创建 PR（自动检测当前分支作为 head）
gc pr create -R infra-test/gctest1 --title "New feature" --body "Description"

# 指定 head 分支
gc pr create -R infra-test/gctest1 --head feature-branch --title "Feature" --body "Description"

# 指定基础分支
gc pr create -R infra-test/gctest1 --base main --title "Feature" --body "Description"

# 创建草稿 PR
gc pr create -R infra-test/gctest1 --title "WIP: Feature" --draft

# 创建跨仓库 PR（从 fork 到 upstream）
gc pr create -R upstream/repo --fork myfork/repo --head feature-branch --title "Feature"

# 从最后一次提交填充标题和内容
gc pr create -R infra-test/gctest1 --fill
```

> **说明**: `--head` 参数可选，未指定时自动检测当前 Git 分支。

## pr list - 列出 PRs

```bash
# 列出所有开放的 PRs
gc pr list -R infra-test/gctest1

# 只列出已关闭的 PRs
gc pr list -R infra-test/gctest1 --state closed

# 只列出已合并的 PRs
gc pr list -R infra-test/gctest1 --state merged

# 限制数量
gc pr list -R infra-test/gctest1 --limit 10
```

## pr view - 查看 PR

```bash
# 查看 PR 详情
gc pr view 1 -R infra-test/gctest1

# 查看评论
gc pr view 1 -R infra-test/gctest1 --comments

# 在浏览器中打开
gc pr view 1 -R infra-test/gctest1 --web
```

## pr comments - 查看 PR 评论

```bash
# 查看 PR 评论列表
gc pr comments 1 -R infra-test/gctest1

# 限制评论数量
gc pr comments 1 --limit 5 -R infra-test/gctest1
```

## pr reply - 回复 PR 评论

```bash
# 回复评论讨论
gc pr reply 1 --discussion <discussion_id> --body "回复内容" -R infra-test/gctest1
```

## pr diff - 查看 PR 差异

```bash
gc pr diff 1 -R infra-test/gctest1
```

## pr checkout - 检出 PR 分支

```bash
gc pr checkout 1 -R infra-test/gctest1
```

## pr merge - 合并 PR

```bash
# 合并 PR（默认合并提交）
gc pr merge 1 -R infra-test/gctest1

# Squash 合并
gc pr merge 1 -R infra-test/gctest1 --squash

# Rebase 合并
gc pr merge 1 -R infra-test/gctest1 --rebase
```

## pr close - 关闭 PR

```bash
gc pr close 1 -R infra-test/gctest1
```

## pr reopen - 重开 PR

```bash
gc pr reopen 1 -R infra-test/gctest1
```

## pr ready - 标记就绪状态

```bash
# 标记为就绪（取消草稿）
gc pr ready 1 -R infra-test/gctest1

# 标记为草稿
gc pr ready 1 -R infra-test/gctest1 --wip
```

## pr review - 评审 PR

```bash
# 评论 PR
gc pr review 1 --comment "评审意见" -R infra-test/gctest1

# 批准 PR
gc pr review 1 --approve -R infra-test/gctest1

# 请求修改
gc pr review 1 --request -R infra-test/gctest1

# 强制通过审批（管理员权限）
gc pr review 1 --approve --force -R infra-test/gctest1
```

## pr edit - 编辑 PR

```bash
# 修改标题
gc pr edit 1 --title "新标题" -R infra-test/gctest1

# 修改描述
gc pr edit 1 --body "新描述" -R infra-test/gctest1

# 设置草稿状态
gc pr edit 1 --draft true -R infra-test/gctest1

# 取消草稿状态
gc pr edit 1 --draft false -R infra-test/gctest1

# 添加标签
gc pr edit 1 --labels bug,enhancement -R infra-test/gctest1

# 设置里程碑
gc pr edit 1 --milestone 5 -R infra-test/gctest1
```

## pr test - 触发 PR 测试

```bash
# 触发测试
gc pr test 1 -R infra-test/gctest1

# 强制通过测试（管理员权限）
gc pr test 1 --force -R infra-test/gctest1
```