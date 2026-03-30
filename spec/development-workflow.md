# 开发工作流程

本文档定义完整的开发工作流程，**严格遵守以下流程，违反将导致代码管理混乱！**

## 职责

定义从 issue 到 merge 的总开发流程，以及在本仓库中开发时必须遵守的顺序和禁止行为。

## 适用场景

- 修复 issue
- 开发新功能
- 提交 PR
- 按仓库规范完成评论、评审和合并

## 必须

- 修复前先验证问题是否仍存在
- 不在 `main` 直接开发
- 补测试、做本地验证和实际命令验证
- 在提 PR 前完成安全审查
- 按流程补 issue comment 和 PR review comment

## 禁止

- 看到 issue 就直接写代码
- 跳过验证、测试或安全检查
- 在主分支直接提交

## 同步要求

- 改动流程时同步 `spec/*`
- 改动命令行为时同步 `docs/COMMANDS.md`
- 改动 AI 协作或入口规则时同步 `AGENTS.md`、`CLAUDE.md` 和相关 skills

## 不负责什么

- 具体编码风格
- 测试覆盖细节
- 文档分层设计
- 代码质量 blocker 判定

## 修复 Issue 前的验证（重要）

**在开始修复任何 Issue 之前，必须先验证问题是否仍存在：**

### 验证步骤

1. **用当前版本验证**：执行复现步骤，确认问题是否存在
2. **检查时间线**：Issue 创建时间、相关代码提交时间、最新发布版本时间
3. **检查是否已修复**：对比 Issue 创建时和当前代码版本

### 错误示例

```bash
# ❌ 看到 Issue 就立即创建分支、写代码
gc issue view 123
git checkout -b bugfix/issue-123
# 直接开始修改...
```

### 正确示例

```bash
# ✅ 先验证问题是否存在
./gc issue view 123
# 查看 Issue 描述的复现步骤

# 执行复现命令，确认问题是否存在
./gc xxx command --params

# 如果问题不存在，在 Issue 中评论说明并关闭
gc issue comment 123 --body "经测试，当前版本该命令正常工作，问题已不存在。" -R gitcode-cli/cli
gc issue close 123 -R gitcode-cli/cli

# 如果问题确实存在，再开始修复流程
```

### 常见误判场景

| 场景 | 说明 |
|------|------|
| 临时 API 故障 | GitCode API 临时返回异常数据，被误认为代码 Bug |
| 特定参数触发 | 使用了特定参数组合才触发的错误 |
| 环境差异 | 本地环境与 Issue 报告者环境不同 |
| 已被其他 PR 修复 | 相关代码已在其他 PR 中修复 |

## 流程概览

```
提交 Issue → 打标签 → 创建分支 → 分支开发 → 编写测试 → 本地测试 → 实际命令测试 → 安全审查 → 提交 PR → Issue 评论 → PR 审查评论 → 关闭 Issue → 合并 PR
```

## 完整流程步骤

### 1. 提交 Issue

发现 BUG 或需要新特性后，首先在项目中创建 Issue。

```bash
gc issue create --title "Bug: 描述问题" --body "问题描述" -R gitcode-cli/cli
```

### 2. 打标签

Issue 创建后立即打上合适的标签。

```bash
gc issue label <number> --add bug -R gitcode-cli/cli
```

### 3. 创建开发分支

从 main 分支创建对应类型的分支，**绝对禁止在 main 分支直接修改**。

```bash
# 确保在 main 分支并更新
git checkout main
git pull

# 创建开发分支
git checkout -b feature/issue-<number>
# 或
git checkout -b bugfix/issue-<number>
```

### 4. 在分支开发

在新建的分支上进行开发，不在 main 分支修改。

### 5. 编写测试用例

为新功能或修复编写单元测试。

```bash
# 运行测试
go test ./pkg/cmd/xxx/...
```

### 6. 本地测试

运行单元测试确保功能正常。

```bash
go test ./...
```

### 7. 实际命令测试

**单元测试无法覆盖所有场景，必须进行实际命令测试！**

```bash
# 构建本地版本
go build -o ./gc ./cmd/gc

# 在测试仓库验证
export GC_TOKEN=your_token
./gc issue list -R infra-test/gctest1
```

### 8. 安全审查

**提交代码前必须进行安全审查！**

检查以下项目：

- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] `.gitignore` 已忽略敏感文件
- [ ] 测试代码不包含真实 Token
- [ ] 文档中不包含真实凭证

```bash
# 检查即将提交的内容
git diff --cached

# 检查是否有敏感信息
git diff --cached | grep -iE "token|password|secret|api_key"

# 检查是否有敏感文件被追踪
git ls-files | grep -iE "\.pem|\.key|\.env|credentials|secret"
```

详细安全规范参见 [安全规范](./security.md)。

### 9. 提交代码

提交代码，commit 信息关联 Issue。

```bash
git add .
git commit -m "feat: add xxx command

Closes #<number>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
```

### 10. 推送分支

```bash
git push -u origin feature/issue-<number>
```

### 11. 创建 PR

创建 PR 合并到 main，描述中关联 Issue。

```bash
gc pr create --title "feat: add xxx command" --body "Closes #<number>" --base main -R gitcode-cli/cli
```

### 12. Issue 评论

在 Issue 中添加完成说明。

```bash
gc issue comment <number> --body "## 修复完成

### 解决方案
描述如何解决的

### 测试结果
- [x] 单元测试通过
- [x] 实际命令测试通过" -R gitcode-cli/cli
```

### 13. PR 审查评论

在 PR 中提交审查评论。

```bash
gc pr review <pr_number> --comment "## 审查结果

### 改动内容
- 新增 xxx 命令

### 测试结果
- [x] 单元测试通过
- [x] 实际命令测试通过" -R gitcode-cli/cli
```

### 14. 关闭 Issue

审查通过后关闭关联的 Issue。

```bash
gc issue close <number> -R gitcode-cli/cli
```

### 15. 合并 PR

确认所有测试通过后合并 PR。

```bash
gc pr merge <pr_number> -R gitcode-cli/cli
```

## 下一步去看哪里

- 如果你在写代码，继续看 [编码规范](./coding-standards.md)
- 如果你在补测试，继续看 [测试指南](./testing-guide.md)
- 如果你准备提交，继续看 [代码质量门禁规范](./code-quality-gates.md)
- 如果你改了命令或协作规则，继续看 [文档治理规范](./docs-governance.md)

### 16. 拉取最新代码

```bash
git checkout main && git pull
```

## 分支命名规范

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| BUG 修复 | `bugfix/issue-<number>` | `bugfix/issue-33` |
| 新特性 | `feature/issue-<number>` | `feature/issue-23` |
| 文档更新 | `docs/issue-<number>` | `docs/issue-5` |

## 标签使用规范

| 标签 | 使用场景 |
|------|----------|
| `bug` | 错误修复 |
| `enhancement` | 功能增强/新特性 |
| `documentation` | 文档更新 |
| `help wanted` | 需要帮助 |
| `question` | 需要讨论 |

## 提交规范

### 提交信息

使用 Conventional Commits：

- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `test:` 测试相关
- `refactor:` 重构

### 提交要求

- **单次提交限制**: 每次代码提交不超过 **800 行**
- **及时提交**: 完成一个功能点或修复后立即提交
- **原子提交**: 每个提交应是一个独立的、完整的功能或修复
- **立即推送**: 每次提交后立即推送到远端
- **关联 Issue**: commit 信息和 PR 描述中都要关联 Issue

## 完整工作流检查清单

开发完成后必须确认：

### 功能检查
- [ ] 单元测试全部通过 (`go test ./...`)
- [ ] 在测试仓库进行实际命令测试
- [ ] Issue 已打标签
- [ ] PR 已创建并关联 Issue
- [ ] Issue 已添加完成评论
- [ ] PR 已提交审查评论
- [ ] Issue 已关闭
- [ ] PR 已合并

### 安全检查
- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] `.gitignore` 已忽略敏感文件
- [ ] 测试代码不包含真实 Token
- [ ] 文档中不包含真实凭证

## 禁止行为

### 流程相关
- ❌ 直接在 main 分支开发
- ❌ 不创建 Issue 直接开发
- ❌ Issue 创建后不打标签
- ❌ PR 不关联 Issue
- ❌ 未编写测试用例就提交 PR
- ❌ 单元测试未通过就提交 PR
- ❌ 未进行实际命令测试就合并 PR
- ❌ 未添加 Issue 评论就关闭 Issue
- ❌ 未提交 PR 审查评论就合并 PR

### 安全相关
- ❌ 硬编码 Token 或密钥到代码中
- ❌ 将敏感文件提交到版本控制
- ❌ 在测试代码中使用真实 Token
- ❌ 在文档中记录真实凭证
- ❌ 使用非授权的测试仓库

## 完整流程示例

```bash
# 1. 确保在 main 分支并更新
git checkout main
git pull

# 2. 创建开发分支
git checkout -b feature/issue-23

# 3. 开发代码
# ... 编写代码 ...

# 4. 编写测试用例
# ... 创建 xxx_test.go ...

# 5. 运行单元测试
go test ./pkg/cmd/issue/label/...

# 6. 实际命令测试
./gc issue label 1 --add bug -R infra-test/gctest1

# 7. 安全审查
git diff | grep -iE "token|password|secret|api_key"
# 确认无敏感信息

# 8. 提交代码
git add .
git commit -m "feat(issue): add label command

Closes #23

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"

# 9. 推送分支
git push -u origin feature/issue-23

# 10. 创建 PR
gc pr create --title "feat: add issue label command" --body "Closes #23" --base main -R gitcode-cli/cli

# 11. 给 Issue 打标签
gc issue label 23 --add enhancement -R gitcode-cli/cli

# 12. 在 Issue 中添加完成评论
gc issue comment 23 --body "## 实现完成

- 新增 gc issue label 命令
- 支持添加、移除、列出标签
- 单元测试通过
- 实际命令测试通过" -R gitcode-cli/cli

# 13. 在 PR 中提交审查评论
gc pr review <pr_number> --comment "## 审查结果

### 改动内容
- 新增 issue label 命令

### 安全检查
- [x] 无硬编码 Token
- [x] 无敏感信息泄露

### 测试结果
- [x] 单元测试通过
- [x] 实际命令测试通过" -R gitcode-cli/cli

# 14. 关闭 Issue
gc issue close 23 -R gitcode-cli/cli

# 15. 合并 PR
gc pr merge <pr_number> -R gitcode-cli/cli

# 16. 拉取最新代码
git checkout main && git pull
```

---

**最后更新**: 2026-03-26
