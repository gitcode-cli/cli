# 测试流程

本文档定义测试的完整流程，与 [测试指南](../testing-guide.md) 配合使用。

## ⚠️ 重要：测试仓库限制

**必须使用以下测试仓库进行实际命令测试：**

| 仓库 | 用途 |
|------|------|
| `infra-test/gctest1` | 功能测试、集成测试（首选） |
| `infra-test` 组织下其他仓库 | 其他测试场景 |

**禁止行为：**
- ❌ 使用个人仓库测试
- ❌ 使用其他组织或用户的仓库测试
- ❌ 使用 `gitcode-cli/cli` 测试
- ❌ 在生产环境仓库测试

## 流程概览

```
开发完成 → 单元测试 → 实际命令测试 → 验证结果 → 提交 PR
```

## 1. 单元测试

### 测试时机
- 新功能开发完成后
- Bug 修复完成后
- 提交 PR 之前

### 测试步骤

```bash
# 1. 运行所有测试
go test ./...

# 2. 运行特定模块测试
go test ./pkg/cmd/issue/...

# 3. 运行特定测试用例
go test -run TestLabelCmd ./pkg/cmd/issue/label/...

# 4. 查看覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 测试要求
- 所有测试必须通过
- 新功能覆盖率 ≥ 70%
- 核心模块覆盖率 ≥ 80%

## 2. 实际命令测试

**重要：单元测试无法覆盖所有场景，必须进行实际命令测试！**

### 测试前准备

```bash
# 1. 设置 Token
export GC_TOKEN=your_token

# 2. 验证认证状态
./gc auth status

# 3. 确认测试仓库
# 只能使用：infra-test/gctest1 或 infra-test 组织下仓库
```

### 测试步骤

```bash
# 1. 构建本地版本
go build -o ./gc ./cmd/gc

# 2. 执行测试命令（示例）
./gc issue list -R infra-test/gctest1
./gc issue create --title "Test" --body "Body" -R infra-test/gctest1
./gc issue label 1 --add bug -R infra-test/gctest1

# 3. 验证结果
# 检查命令输出
# 在 Web 界面验证操作结果
```

### 测试范围

| 测试类型 | 说明 | 示例 |
|---------|------|------|
| 正常流程 | 正常使用场景 | `gc issue list` |
| 边界条件 | 极端输入 | 空值、特殊字符 |
| 错误处理 | 预期错误 | 缺少参数、无效输入 |

## 3. 验证结果

### 检查项

- [ ] 命令输出正确
- [ ] 错误信息清晰
- [ ] Web 界面显示正确
- [ ] 数据已正确保存

### 验证示例

```bash
# 创建 Issue 后验证
gc issue view <number> -R infra-test/gctest1

# 添加标签后验证
gc issue label <number> --list -R infra-test/gctest1

# 创建 PR 后验证
gc pr view <number> -R infra-test/gctest1
```

## 4. 测试检查清单

### 提交 PR 前

- [ ] 单元测试全部通过
- [ ] 实际命令测试通过
- [ ] 测试仓库使用正确
- [ ] 测试结果已记录

### PR 描述中

```markdown
## 测试结果

- [x] 单元测试通过
- [x] 实际命令测试通过
- [x] 正常流程测试
- [x] 边界条件测试
- [x] 错误处理测试
```

## 5. 测试清理

测试完成后，清理测试数据：

```bash
# 关闭测试 Issue
gc issue close <number> -R infra-test/gctest1

# 删除测试分支
git branch -D test-branch

# 清理本地构建
rm -f ./gc
```

## 完整测试流程示例

```bash
# 1. 运行单元测试
go test ./pkg/cmd/issue/label/...

# 2. 构建本地版本
go build -o ./gc ./cmd/gc

# 3. 验证认证
./gc auth status

# 4. 测试 Issue label 命令
./gc issue label 1 --add bug -R infra-test/gctest1
./gc issue label 1 --list -R infra-test/gctest1
./gc issue label 1 --remove bug -R infra-test/gctest1

# 5. 验证结果
# 检查标签是否正确添加/移除

# 6. 记录测试结果
# 在 PR 描述中记录测试通过
```

---

**最后更新**: 2026-03-26