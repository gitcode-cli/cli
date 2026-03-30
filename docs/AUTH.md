# 认证模型

本文档定义当前版本 GitCode CLI 的认证来源优先级和持久化语义。

## 当前模型

当前版本支持：

1. 环境变量认证
2. 本地配置持久化认证
3. 单主机单活动用户模型

当前版本不支持：

- 多账户切换的完整产品化体验
- 操作系统 keyring 集成

## 认证来源优先级

同一主机下，活动 token 的解析顺序为：

1. `GC_TOKEN`
2. `GITCODE_TOKEN`
3. 本地配置文件中的已登录 token

也就是说：

- 只要设置了环境变量，环境变量始终覆盖本地存储
- `auth logout` 只会清理本地存储，不会替你取消环境变量

## 持久化边界

`auth login` 成功后会将以下信息写入本地配置目录：

- 主机名
- 用户名
- token
- Git protocol

默认位置：

```bash
~/.config/gc/auth.json
```

可通过 `GC_CONFIG_DIR` 覆盖配置目录。

## 各命令统一语义

### auth login

- 校验 token 有效性
- 成功后写入本地配置
- 若同时存在环境变量，后续命令仍优先使用环境变量

### auth status

- 按统一优先级解析当前活动 token
- 显示 token 来源：`GC_TOKEN`、`GITCODE_TOKEN` 或 `config`

### auth token

- 输出当前活动 token
- 与 `auth status` 看到的来源保持一致

### auth logout

- 清理本地配置中的认证信息
- 如果当前活动 token 来自环境变量，会明确提示用户手动 `unset`
