# 开发计划模板

基于选定方案，AI 输出开发计划（文件清单、修改顺序、验证策略），写入 Issue comment。

## 输出格式

```markdown
## 开发计划

### 文件清单

| # | 文件 | 操作 | 预计行数 | 说明 |
|---|------|------|----------|------|
| 1 | `path/to/file.go` | 修改 | +N/-M | <说明> |
| 2 | `path/to/file_test.go` | 修改 | +N/-M | <说明> |

### 执行顺序
1. <步骤1> → 验证: <检查方式>
2. <步骤2> → 验证: <检查方式>
3. ...

### 测试策略

| 层级 | 方式 | 覆盖目标 |
|------|------|----------|
| 编译 | `go build ./...` | 语法/类型正确 |
| 单元 | `go test ./pkg/xxx/...` | 修改路径的 happy path + error path |
| 全量 | `go test ./...` | 无回归 |
| 命令 | `./gc xxx -R infra-test/gctest1` | 实际命令行为不变 |
| CI | GitHub Actions | 全平台通过 |
```

## 示例

```markdown
## 开发计划

### 文件清单

| # | 文件 | 操作 | 预计行数 | 说明 |
|---|------|------|----------|------|
| 1 | `pkg/cmd/repo/sync/sync.go` | 修改 | +1/-3 | 移除 var gitRun，clone 调用改为 opts.GitRun |
| 2 | `pkg/cmd/repo/sync/sync_test.go` | 修改 | +23/-12 | mock 注入方式改为 opts 字段 |
| 3 | `pkg/cmd/pr/sync/sync.go` | 修改 | +13/-13 | SyncOptions 新增 GitRun/GitRunInDir；syncCommits 签名变更 |
| 4 | `pkg/cmd/pr/sync/sync_test.go` | 修改 | +9/-17 | mock 注入方式改为 opts 字段 / 局部函数 |

### 执行顺序
1. repo/sync/sync.go → 验证: go build + grep 确认 var gitRun 已移除
2. repo/sync/sync_test.go → 验证: go test ./pkg/cmd/repo/sync/ -count=1
3. pr/sync/sync.go → 验证: go build + grep 确认 var gitRunWithEnv/gitRunInDirWithEnv 已移除
4. pr/sync/sync_test.go → 验证: go test ./pkg/cmd/pr/sync/ -count=1
5. 全量回归 → 验证: go test ./... -count=1
6. pre-commit → 验证: pre-commit run --all-files
7. 风险分级 → 验证: python3 scripts/classify-change-risk.py --base origin/main

### 测试策略

| 层级 | 方式 | 覆盖目标 |
|------|------|----------|
| 编译 | `go build ./cmd/gc` | 类型正确，无未定义引用 |
| 单元 | `go test ./pkg/cmd/repo/sync/ ./pkg/cmd/pr/sync/` | 29 个测试全部通过 |
| 全量 | `go test ./...` | 1268 个测试无回归 |
| 命令 | 不适用 (纯重构，行为不变) | — |
| CI | GitHub Actions | linux/macOS test + build + docker |
```
