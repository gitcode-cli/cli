## 开发计划: Issue #395

| # | 文件 | 操作 | 说明 |
|---|------|------|------|
| 1 | git/git.go | 修改 (+3/-1) | RemoteURL 加 ValidateRef + `--` |
| 2 | git/git_test.go | 修改 (+22/-0) | 新增 TestRemoteURLRejectsOptionInjection（4 用例） |

### 测试矩阵
| 类型 | 覆盖 | 状态 |
|------|------|------|
| UT | option 注入（--upload-pack=）被拒 | ✅ |
| UT | dash 前缀（-bogus）被拒 | ✅ |
| UT | 空 remote 被拒 | ✅ |
| UT | shell metacharacter 被拒 | ✅ |
