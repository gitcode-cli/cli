# 需求分析模板

AI 拿到 issue 后、写任何代码前，按此模板输出需求分析，写入 Issue comment。

## 输出格式

```markdown
## 需求分析

### 问题定义
<一句话描述核心问题，引用 issue 中的关键信息>

### 影响范围
| 层级 | 受影响 |
|------|--------|
| 模块 | <pkg/xxx> |
| 文件 | <具体文件路径> |
| 函数/类型 | <函数名或类型名> |

### 根因推断
<为什么会出这个问题：设计缺陷/实现遗漏/边界未覆盖/外部变更/...>

### 成功标准
- [ ] <可验证的条件1>
- [ ] <可验证的条件2>

### 约束条件
- 不得破坏: <现有行为/API/接口>
- 兼容性: <向后兼容/需迁移指南>
- 其他: <性能/安全/文档要求>
```

## 示例

```markdown
## 需求分析

### 问题定义
`pkg/cmd/repo/sync` 和 `pkg/cmd/pr/sync` 使用包级 `var` 声明 git mock 函数，阻塞并行测试 (`go test -parallel`)。

### 影响范围
| 层级 | 受影响 |
|------|--------|
| 模块 | pkg/cmd/repo/sync, pkg/cmd/pr/sync |
| 文件 | sync.go x2, sync_test.go x2 |
| 函数/类型 | SyncOptions, syncRun, syncCommits |

### 根因推断
为了测试时替换 git 函数，在包级声明了 `var gitRun = gitpkg.RunWithEnv` 等变量。这是 Go 中常见的 mock 模式，但引入包级可变全局状态，与 `t.Parallel()` 不兼容。

### 成功标准
- [ ] 移除所有包级可变 git 函数变量
- [ ] 测试可通过 Options struct 字段注入 mock
- [ ] go test -count=1 ./pkg/cmd/repo/sync/ ./pkg/cmd/pr/sync/ 全部通过
- [ ] 命令行为不变

### 约束条件
- 不得破坏: repo sync / pr sync 现有命令行为
- 兼容性: 向后兼容，NewCmdSync 签名不变
- 其他: 遵循 command-template.md 中行为函数注入模式
```
