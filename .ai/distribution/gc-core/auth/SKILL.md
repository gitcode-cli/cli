# gc-auth

使用 `gc` 完成 GitCode 认证相关操作。

## 触发场景

- 需要配置 `gc` 认证
- 需要检查当前认证状态
- 需要查看当前生效的 token 来源
- 需要清理本地登录状态

## 常用命令

```bash
# 环境变量认证
export GC_TOKEN="your_gitcode_token"

# 交互式登录
gc auth login

# 使用 token 登录
gc auth login --token YOUR_TOKEN

# 查看认证状态
gc auth status

# 显示当前生效 token
gc auth token

# 清理本地登录状态
gc auth logout
```

## 使用约束

- 环境变量优先于本地登录配置
- `gc auth logout` 只清理本地配置，不会自动取消环境变量
- 在共享机器或 CI 环境中，优先使用环境变量

## 常见提醒

- `GC_TOKEN` 是首选环境变量
- 如使用 `GITCODE_TOKEN`，其优先级低于 `GC_TOKEN`
- 在排查权限问题前，先运行 `gc auth status`
