# 安全策略

## 敏感信息保护

本项目严格遵守以下安全规则：

### Token 管理

1. **禁止硬编码 Token**：Token 绝对不能写入源代码或任何持久化存储
2. **使用环境变量**：Token 仅通过环境变量传递
   ```bash
   export GC_TOKEN="your_token"
   # 或
   export GITCODE_TOKEN="your_token"
   ```
3. **内存存储**：Token 仅在内存中保存，程序结束后自动清除

### 禁止提交的内容

以下文件类型已被 `.gitignore` 忽略，禁止提交：

| 文件类型 | 说明 |
|----------|------|
| `*.pem`, `*.key`, `*.p12`, `*.pfx` | 私钥/证书文件 |
| `id_rsa*`, `id_ed*` | SSH 密钥 |
| `.env`, `.env.*` | 环境变量文件 |
| `*.secret` | 密钥文件 |
| `credentials.json` | 凭证文件 |
| `token.txt`, `*.token` | Token 文件 |
| `secrets.yaml`, `secrets.yml` | 密钥配置 |

### 代码审查检查项

提交前必须确认：

- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] `.gitignore` 已忽略敏感文件
- [ ] 测试代码不包含真实 Token
- [ ] 文档中不包含真实凭证

### 测试规范

- **测试组织**: `infra-test`
- **测试仓库**: `infra-test/gctest1`
- **Token 来源**: 环境变量或运行时输入

### CI/CD 安全

在 GitHub Actions 中使用 Secrets 存储敏感信息：

```yaml
env:
  GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  DOCKER_USERNAME: ${{ secrets.DOCKER_USERNAME }}
  DOCKER_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}
```

### 安全检查命令

```bash
# 检查即将提交的内容
git diff --cached

# 检查历史中的敏感信息
git log -p | grep -iE "token|password|secret|api_key"

# 检查是否有敏感文件被追踪
git ls-files | grep -iE "\.pem|\.key|\.env|credentials|secret"
```

## 报告安全问题

如果您发现安全问题，请通过以下方式报告：

1. 不要在公开 Issue 中报告安全问题
2. 发送邮件到项目维护者
3. 提供详细的问题描述和复现步骤

---

**最后更新**: 2026-03-23