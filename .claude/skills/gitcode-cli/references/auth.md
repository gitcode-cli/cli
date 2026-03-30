# 认证命令 (auth)

> 本文档是 Claude 参考层，不是命令行为真相源。
> 认证行为以 `docs/COMMANDS.md`、`docs/AUTH.md` 和 `spec/` 为准。

## auth login - 登录

```bash
# 交互式登录
gc auth login

# 使用 Token 登录
gc auth login --token YOUR_TOKEN
```

## auth status - 查看认证状态

```bash
gc auth status
```

## auth token - 显示 Token

```bash
gc auth token
```

## auth logout - 登出

```bash
gc auth logout
```

> 当前认证优先级：
> `GC_TOKEN` -> `GITCODE_TOKEN` -> 本地登录配置。
