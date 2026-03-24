---
name: pr-reviewer
description: |
  Review Pull Requests for GitCode CLI project. Check code quality, security issues,
  and compliance with coding standards.

  TRIGGER when: user asks to review PR, check PR quality, audit PR security, analyze
  PR changes, or phrases like "评审PR", "审查PR", "检查PR", "代码审查", "安全审查".
---

# PR Reviewer

自动化评审 GitCode CLI 项目的 Pull Request，检查代码质量、安全问题和规范合规性。

## 评审流程

### 1. 获取 PR 信息

```bash
# 查看 PR 详情
gc pr view <number> -R gitcode-cli/cli

# 查看 PR 代码变更
gc pr diff <number> -R gitcode-cli/cli
```

### 2. 检查清单

#### 测试检查

| 检查项 | 说明 |
|--------|------|
| 单元测试 | 新增代码是否有测试覆盖 |
| 测试通过 | `go test ./...` 是否通过 |
| 边界条件 | 是否测试边界情况 |

#### 安全检查

| 检查项 | 说明 |
|--------|------|
| 敏感信息 | 无硬编码 Token、密码、密钥 |
| 输入验证 | 用户输入是否验证 |
| SQL 注入 | 数据库查询是否安全 |
| 命令注入 | 是否安全执行外部命令 |
| XSS | 输出是否转义 |

#### 代码规范

| 检查项 | 说明 |
|--------|------|
| 命名规范 | 变量、函数、文件命名符合规范 |
| 注释完整 | 复杂逻辑是否有注释 |
| 无冗余代码 | 无死代码、重复代码 |
| 错误处理 | 错误是否正确处理 |

#### 变更分析

| 检查项 | 说明 |
|--------|------|
| 变更范围 | 修改文件数量和影响范围 |
| 破坏性变更 | 是否影响向后兼容 |
| 文档同步 | 是否需要更新文档 |

## 安全敏感模式

检查以下敏感信息模式：

```
# Token/密钥
(GC_TOKEN|TOKEN|API_KEY|SECRET|PASSWORD)\s*[=:]\s*["\']?[^\s"\']+
# 私钥
-----BEGIN (RSA |EC |DSA )?PRIVATE KEY-----
# 连接字符串
(mysql|postgres|mongodb)://[^\s]+:[^\s]+@
```

## 执行命令

```bash
# 添加评论
gc pr review <number> --comment "评审意见" -R gitcode-cli/cli

# 批准 PR
gc pr review <number> --approve -R gitcode-cli/cli

# 请求修改
gc pr review <number> --request -R gitcode-cli/cli

# 运行测试
go test ./...
go test -coverprofile=coverage.out ./...
```

## 评审结论模板

### 通过

```markdown
## 评审结论: ✅ 通过

### 检查结果
- [x] 测试检查通过
- [x] 安全检查通过
- [x] 代码规范通过

### 代码质量
[简述代码质量评价]

### 建议
[可选的改进建议]
```

### 需要修改

```markdown
## 评审结论: ⚠️ 需要修改

### 问题列表
1. **[安全问题]** [问题描述]
2. **[代码规范]** [问题描述]

### 修改建议
[具体的修改建议]
```

## 注意事项

1. 重点关注安全问题和敏感信息
2. 检查新增代码是否有测试覆盖
3. 注意破坏性变更
4. 评论保持专业和建设性