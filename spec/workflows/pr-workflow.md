# PR 流程

本文档定义 Pull Request 的生命周期管理、自检要求和合并规则。

## 职责

- 定义 PR 的标准状态流转
- 约束作者自检与独立评审的边界
- 明确合并前必须具备的证据

## 流程概览

```
创建分支
→ 开发代码
→ 编写测试
→ 提交代码
→ 创建 PR draft
→ 作者自检
→ ready-for-review
→ 独立评审
→ approved
→ merged
```

## 1. 创建分支

### 分支命名规范

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| BUG 修复 | `bugfix/issue-<number>` | `bugfix/issue-33` |
| 新特性 | `feature/issue-<number>` | `feature/issue-23` |
| 文档更新 | `docs/issue-<number>` | `docs/issue-5` |
| 重构 | `refactor/issue-<number>` | `refactor/issue-19` |

### 创建步骤

```bash
git checkout main
git pull
git checkout -b feature/issue-23
```

## 2. PR 标签规范

PR 至少应补齐以下维度：

- 类型：`type/bug`、`type/feature`、`type/docs`、`type/refactor`
- 状态：`status/draft`、`status/self-checked`、`status/ready-for-review`、`status/changes-requested`、`status/approved`、`status/merged`
- 风险：`risk/low`、`risk/medium`、`risk/high`
- 范围：`scope/auth`、`scope/repo`、`scope/issue`、`scope/pr`、`scope/release`、`scope/docs`

## 3. 开发与测试

### 开发规范

- 遵循 [编码规范](../foundations/coding-standards.md)
- 使用 [命令开发模板](../foundations/command-template.md)
- 命令行为变化后同步相关文档

### 最低测试要求

```bash
go test ./...
go build -o ./gc ./cmd/gc
```

如果改动影响具体命令行为，还必须补：

```bash
go test ./pkg/cmd/xxx/...
./gc xxx -R infra-test/gctest1
```

## 4. 提交代码

### 提交规范

使用 Conventional Commits：

- `feat:` 新功能
- `fix:` Bug 修复
- `docs:` 文档更新
- `test:` 测试相关
- `refactor:` 重构

### 提交步骤

```bash
git add <files>
git commit -m "feat: add issue label command

Refs #23

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>"
git push -u origin feature/issue-23
```

### 提交要求

- 单次提交应保持范围清晰
- 不把无关改动混入同一个 PR
- 未经过自检的代码不应直接宣称可合并

## 5. 创建 PR

### 创建命令

```bash
gc pr create --title "feat: add issue label command" \
  --body "## 变更内容

- 新增 gc issue label 命令

## 关联

Refs #23" \
  --base main \
  -R gitcode-cli/cli
```

创建后先补 `status/draft`，不要一开始就当作 ready for review。

## 6. 作者自检

作者自检是必需步骤，但不等同于独立评审。

### 作者自检模板

```markdown
## 作者自检

- 根因或实现理由:
- 主要修改:
- 单元测试:
- 构建:
- 实际命令验证:
- 安全审查:
- 文档同步:
- 风险:
- 未覆盖项:
```

### 自检要求

- 自检记录应写在 PR 中
- 自检记录必须包含安全审查结果；若本次改动不涉及认证、凭证、权限或危险写路径，应明确写出无相关影响
- 自检完成后，PR 才能进入 `status/self-checked`
- 只有自检，不代表可以直接合并

示例：

```bash
gc pr review <number> --comment "## 作者自检

- 根因或实现理由: ...
- 主要修改: ...
- 单元测试: go test ./...
- 构建: go build -o ./gc ./cmd/gc
- 实际命令验证: ./gc issue label 1 --add bug -R infra-test/gctest1
- 安全审查: 已检查无硬编码凭证，本次改动不涉及认证和危险写路径
- 文档同步: 已更新 docs/COMMANDS.md 与 spec
- 风险: ...
- 未覆盖项: ..." -R gitcode-cli/cli
```

## 7. Ready For Review

进入 `status/ready-for-review` 前必须满足：

- PR 已有作者自检记录
- 关联 issue 已进入 `status/ready-for-review`
- 所有本地门禁已完成
- 安全审查已完成
- 文档同步已完成或明确说明无需更新

## 8. 独立评审与合并

### 评审前提

- 评审者不能把作者自检当作评审完成
- 作者本人不得把自己的评论当作独立通过结论

### 合并前检查

- [ ] 单元测试通过
- [ ] 构建通过
- [ ] 实际命令测试通过，或已说明未执行原因
- [ ] 安全审查已完成，或已明确说明本次改动无安全敏感影响
- [ ] issue 已有验证记录和开发进度记录
- [ ] PR 已有作者自检记录
- [ ] PR 已有独立评审结论
- [ ] 文档已同步

### 合并命令

```bash
gc pr merge <number> -R gitcode-cli/cli
```

### 合并后操作

```bash
git checkout main
git pull
```

合并后更新状态：

- PR：`status/merged`
- Issue：`status/merged` 或关闭

## 9. 检查清单

- [ ] 分支已创建
- [ ] 代码已开发
- [ ] 测试已编写并执行
- [ ] PR 已以 `status/draft` 创建
- [ ] 作者自检已完成
- [ ] PR 已进入 `status/ready-for-review`
- [ ] PR 已获得独立评审结论
- [ ] PR 已合并

---

**最后更新**: 2026-04-02
