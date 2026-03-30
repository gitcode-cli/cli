# 仓库命令 (repo)

> 本文档是 Claude 参考层，不是命令行为真相源。
> 仓库命令行为以 `docs/COMMANDS.md` 和 `spec/` 为准。

## repo view - 查看仓库

```bash
# 查看仓库详情
gc repo view infra-test/gctest1

# 在浏览器中打开
gc repo view infra-test/gctest1 --web
```

## repo list - 列出仓库

```bash
# 列出自己的仓库
gc repo list

# 列出指定组织的仓库
gc repo list --owner infra-test

# 限制数量
gc repo list --limit 10

# 只列出公开仓库
gc repo list --visibility public
```

## repo create - 创建仓库

```bash
# 创建公开仓库
gc repo create my-repo --public

# 创建私有仓库
gc repo create my-repo --private

# 创建带描述的仓库
gc repo create my-repo --public --description "My project"
```

## repo fork - Fork 仓库

```bash
# Fork 仓库到自己的账户
gc repo fork owner/repo

# Fork 并克隆到本地
gc repo fork owner/repo --clone
```

## repo delete - 删除仓库

```bash
# 删除仓库（危险操作，需确认）
gc repo delete owner/repo
```

## repo stats - 代码贡献统计

```bash
# 获取 main 分支代码贡献统计
gc repo stats --branch main -R owner/repo

# 按作者筛选
gc repo stats --branch main --author username -R owner/repo

# 仅显示个人统计
gc repo stats --branch main --only-self -R owner/repo

# 指定日期范围
gc repo stats --branch main --since 2024-01-01 --until 2024-12-31 -R owner/repo
```

---

**最后更新**: 2026-03-26
