# 测试需求

本文档详细描述 gitcode-cli 的测试策略和测试需求。

## 测试概览

### 测试金字塔

```
                ┌─────────┐
                │   E2E   │  端到端测试
                │  Tests  │  - 完整流程验证
                ├─────────┤
               ┌┴─────────┴┐
               │ Integration│  集成测试
               │   Tests    │  - 模块间交互
               ├───────────┤
          ┌────┴───────────┴────┐
          │    Unit Tests       │  单元测试
          │                     │  - 函数/方法级别
          └─────────────────────┘
```

### 测试类型

| 测试类型 | 描述 | 运行频率 | 执行时间 |
|----------|------|----------|----------|
| 单元测试 | 测试单个函数/方法 | 每次提交 | 秒级 |
| 集成测试 | 测试模块间交互 | 每次 PR | 分钟级 |
| 端到端测试 | 测试完整流程 | 发布前 | 分钟级 |
| 性能测试 | 测试性能指标 | 定期 | 分钟级 |

---

## TEST-001: 单元测试框架

### 功能描述

建立完整的单元测试框架，包括测试规范、Mock 设计和断言库。

### 测试目录结构

```
gitcode-cli/
├── pkg/
│   └── cmd/
│       ├── auth/
│       │   ├── login/
│       │   │   ├── login.go
│       │   │   └── login_test.go
│       │   └── auth_test.go
│       └── repo/
│           └── clone/
│               ├── clone.go
│               └── clone_test.go
├── internal/
│   └── config/
│       ├── config.go
│       └── config_test.go
├── api/
│   ├── client.go
│   └── client_test.go
└── test/
    ├── integration/
    │   └── auth_integration_test.go
    └── e2e/
        └── e2e_test.go
```

### 测试模板

```go
// pkg/cmd/auth/login/login_test.go
package login

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestNewCmdLogin(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"default", []string{}, false},
        {"with hostname", []string{"--hostname", "gitcode.com"}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            factory := NewMockFactory()
            cmd := NewCmdLogin(factory, nil)
            cmd.SetArgs(tt.args)

            err := cmd.Execute()
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

### Mock 设计

```go
// pkg/cmdutil/factory_mock.go
type MockFactory struct {
    Config     func() (gc.Config, error)
    HttpClient func() (*http.Client, error)
    IO         *iostreams.IOStreams
    BaseRepo   func() (gtrepo.Interface, error)
}

func NewMockFactory() *MockFactory

// internal/prompter/prompter_mock.go
type MockPrompter struct {
    ConfirmFunc  func(string, bool) (bool, error)
    InputFunc    func(string, string) (string, error)
    SelectFunc   func(string, string, []string) (int, error)
}
```

### 验收标准

- [ ] 所有模块有单元测试
- [ ] Mock 设计完整
- [ ] 测试命名规范
- [ ] 支持表格驱动测试

---

## TEST-002: 集成测试

### 功能描述

测试模块间的交互，使用真实的或模拟的 API 服务。

### 测试内容

```go
// +build integration

// test/integration/auth_integration_test.go
func TestAuthLoginIntegration(t *testing.T) {
    token := os.Getenv("GC_TEST_TOKEN")
    if token == "" {
        t.Skip("GC_TEST_TOKEN not set")
    }

    // 测试真实 API 调用
    client := api.NewClientFromHTTP(httpClient)
    user, err := api.CurrentLoginName(client, "gitcode.com")
    require.NoError(t, err)
    require.NotEmpty(t, user)
}
```

### Mock HTTP 服务器

```go
// test/server/gitcode_mock.go
type GitCodeMockServer struct {
    Server *httptest.Server
}

func NewGitCodeMockServer() *GitCodeMockServer {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/v5/user", handleUser)
    mux.HandleFunc("/api/v5/repos/", handleRepos)
    // ...
    return &GitCodeMockServer{Server: httptest.NewServer(mux)}
}
```

### 验收标准

- [ ] 支持 `+build integration` 标签
- [ ] 支持 Mock HTTP 服务器
- [ ] 支持真实 API 测试（需 Token）
- [ ] 测试数据可复用

---

## TEST-003: Mock 设计

### 功能描述

设计完整的 Mock 体系，支持单元测试隔离。

### Mock 类型

| Mock | 用途 |
|------|------|
| MockFactory | 命令工厂 |
| MockConfig | 配置管理 |
| MockPrompter | 交互提示 |
| MockHTTPServer | HTTP 服务 |
| MockGitClient | Git 操作 |

### 测试数据

```go
// test/fixtures/fixtures.go
var TestUser = map[string]interface{}{
    "id":       1,
    "username": "testuser",
    "name":     "Test User",
}

var TestRepo = map[string]interface{}{
    "id":         1,
    "name":       "test-repo",
    "full_name":  "owner/test-repo",
    "private":    false,
}
```

### 验收标准

- [ ] Mock 类型完整
- [ ] 测试数据可配置
- [ ] Mock 行为可定制

---

## TEST-004: CI/CD 集成

### 功能描述

配置 GitHub Actions 实现持续集成测试。

### GitHub Actions 配置

```yaml
# .github/workflows/test.yml
name: Test

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run unit tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run integration tests
        run: go test -v -tags=integration ./...
        env:
          GC_TEST_TOKEN: ${{ secrets.TEST_TOKEN }}

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - uses: golangci/golangci-lint-action@v3
```

### 测试矩阵

```yaml
test-matrix:
  strategy:
    matrix:
      os: [ubuntu-latest, macos-latest, windows-latest]
      go: ['1.21', '1.22']
  runs-on: ${{ matrix.os }}
```

### 验收标准

- [ ] 支持 GitHub Actions
- [ ] 支持多平台测试
- [ ] 支持多 Go 版本
- [ ] 支持覆盖率报告

---

## 覆盖率目标

| 模块 | 目标覆盖率 |
|------|-----------|
| `pkg/cmd/` | 80% |
| `internal/config/` | 90% |
| `api/` | 85% |
| `internal/authflow/` | 80% |
| `pkg/iostreams/` | 70% |

---

## 现有测试用例集成

从 https://gitcode.com/afly-infra/gc-api-doc/tree/main/test/ 目录集成以下测试：

| 测试文件 | 对应模块 |
|----------|----------|
| test_users.py | 用户 API 测试 |
| test_repositories.py | 仓库 API 测试 |
| test_issues.py | Issues API 测试 |
| test_pull_requests.py | PR API 测试 |
| test_organizations.py | 组织 API 测试 |
| test_labels.py | 标签 API 测试 |
| test_milestones.py | 里程碑 API 测试 |
| test_search.py | 搜索 API 测试 |
| test_webhooks.py | Webhook API 测试 |
| test_error_codes.py | 错误码测试 |

---

## 测试命令

```bash
# 运行所有测试
go test ./...

# 运行特定包
go test ./pkg/cmd/auth/...

# 运行特定测试
go test -run TestLogin ./pkg/cmd/auth/login/...

# 运行集成测试
go test -tags=integration ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 竞态检测
go test -race ./...

# 运行基准测试
go test -bench=. ./...
```

---

## 相关文档

- [gc-design/docs/testing/strategy.md](https://gitcode.com/afly-infra/gc-design/blob/main/docs/testing/strategy.md)
- [gc-api-doc/test/](https://gitcode.com/afly-infra/gc-api-doc/tree/main/test/)

---

**最后更新**: 2026-03-22