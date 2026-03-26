# 安全规范

本文档定义 gitcode-cli 项目的安全要求。

## Token 管理

### 基本原则

1. **禁止硬编码** - Token 绝对不能写入源代码或任何持久化存储
2. **内存存储** - Token 仅在内存中保存，程序结束后自动清除
3. **禁止提交** - Token 绝对不能提交到版本控制系统

### Token 获取方式

```bash
# 方式一：环境变量（推荐）
export GC_TOKEN="your_gitcode_token"
# 或
export GITCODE_TOKEN="your_gitcode_token"

# 方式二：交互式登录
gc auth login --token YOUR_TOKEN
```

### 获取 Token
1. 登录 [GitCode](https://gitcode.com)
2. 进入 设置 → 私人令牌
3. 生成新令牌并复制

## 敏感信息保护

### 禁止提交的内容

| 文件类型 | 说明 |
|----------|------|
| `*.pem`, `*.key`, `*.p12`, `*.pfx` | 私钥/证书文件 |
| `id_rsa*`, `id_ed*` | SSH 密钥 |
| `.env`, `.env.*` | 环境变量文件 |
| `*.secret` | 密钥文件 |
| `credentials.json` | 凭证文件 |
| `token.txt`, `*.token` | Token 文件 |
| `secrets.yaml`, `secrets.yml` | 密钥配置 |

### .gitignore 配置

确保以下内容已在 `.gitignore` 中：
```gitignore
# Secrets
.env
.env.*
*.pem
*.key
credentials.json
token.txt
secrets.yaml
```

## 测试安全

### 测试仓库限制

**只能使用以下测试仓库：**

| 仓库 | 用途 |
|------|------|
| `infra-test/gctest1` | 主要测试仓库（首选） |
| `infra-test` 组织下其他仓库 | 其他测试场景 |

**禁止行为：**
- ❌ 使用个人仓库测试
- ❌ 使用其他组织或用户的仓库测试
- ❌ 使用 `gitcode-cli/cli` 测试

### 测试 Token 来源

- 通过环境变量 `GC_TOKEN` 传递
- 或运行时手动输入

### 测试配置

```
测试组织: infra-test
API 基础 URL: https://api.gitcode.com/api/v5
Token 来源: 环境变量或运行时输入
```

## 代码审查检查项

提交代码前必须确认：

- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] `.gitignore` 已忽略敏感文件
- [ ] 测试代码不包含真实 Token
- [ ] 文档中不包含真实凭证

## CI/CD 安全

### GitHub Actions

使用 Secrets 存储敏感信息：

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  GC_TOKEN: ${{ secrets.GC_TOKEN }}
```

### PyPI 发布

使用 Trusted Publishing（OIDC），无需配置 API Token。

## 安全检查命令

```bash
# 检查即将提交的内容
git diff --cached

# 检查历史中的敏感信息
git log -p | grep -iE "token|password|secret|api_key"

# 检查是否有敏感文件被追踪
git ls-files | grep -iE "\.pem|\.key|\.env|credentials|secret"
```

## 报告安全问题

如果您发现安全问题，请：

1. 不要在公开 Issue 中报告安全问题
2. 发送邮件到项目维护者
3. 提供详细的问题描述和复现步骤

---

**最后更新**: 2026-03-26