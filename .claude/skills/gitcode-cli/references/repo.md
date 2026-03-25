# 仓库命令 (repo)

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