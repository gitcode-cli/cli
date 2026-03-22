# CLAUDE.md - AI 辅助开发指南

本文档为 Claude Code 提供项目开发指导，帮助 AI 更好地理解和参与 gitcode-cli 项目开发。

## 项目概述

gitcode-cli (命令名: `gc`) 是一款面向 GitCode 平台的命令行工具，深度参考 GitHub CLI (gh) 的架构设计。

### 核心信息

- **语言**: Go 1.21+
- **命令框架**: Cobra
- **配置格式**: YAML
- **目标平台**: Linux, macOS, Windows
- **命令名**: `gc`

### 目录结构

```
gitcode-cli/
├── cmd/gc/           # 程序入口
├── pkg/cmd/          # 命令实现
├── pkg/cmdutil/      # 命令工具
├── pkg/iostreams/    # IO 流管理
├── internal/         # 内部模块
│   ├── config/       # 配置管理
│   ├── authflow/     # 认证流程
│   ├── gtrepo/       # 仓库模型
│   ├── prompter/     # 交互提示
│   ├── browser/      # 浏览器操作
│   ├── keyring/      # 安全存储
│   └── tableprinter/ # 表格输出
├── api/              # API 客户端
├── git/              # Git 操作
├── context/          # 上下文管理
└── docs/             # 文档
```

## 架构设计

### 分层架构

```
CLI 入口 → 命令层 → 服务层 → API 层 → 存储层
```

详细架构设计请参阅 `docs/architecture/overview.md`。

### 核心设计模式

1. **工厂模式** - 依赖注入 (`pkg/cmdutil/factory.go`)
2. **命令模式** - Cobra 子命令
3. **中间件模式** - HTTP 请求处理
4. **策略模式** - 多种认证方式

## 编码规范

### 命名规范

```go
// 包名：小写，简短
package config

// 导出函数：大驼峰
func NewConfig() (*Config, error)

// 内部函数：小驼峰
func parseConfig(data []byte) (*Config, error)

// 接口：动词+er 或名词
type Config interface {}
type TokenGetter interface {}

// 常量：大驼峰或全大写
const DefaultHost = "gitcode.com"
const MAX_RETRIES = 3
```

### 文件组织

```go
// 标准文件结构
package xxx

import (
    // 标准库
    "context"
    "fmt"

    // 第三方库
    "github.com/spf13/cobra"

    // 内部包（按层级排序）
    "github.com/gitcode-cli/cli/internal/config"
    "github.com/gitcode-cli/cli/pkg/cmdutil"
)

// 常量定义
const ()

// 类型定义
type ()

// 接口定义
type interface {}

// 构造函数
func New() *Xxx {}

// 公开方法
func (x *Xxx) Public() {}

// 私有方法
func (x *Xxx) private() {}
```

### 错误处理

```go
// 使用 errors.New 创建简单错误
var ErrNotFound = errors.New("not found")

// 使用 fmt.Errorf 包装错误
return fmt.Errorf("failed to get config: %w", err)

// 自定义错误类型
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}
```

## 命令开发

### 标准命令模板

```go
// pkg/cmd/xxx/xxx.go
package xxx

import (
    "github.com/spf13/cobra"
    cmdutil "github.com/gitcode-cli/cli/pkg/cmdutil"
)

type XxxOptions struct {
    // 注入依赖
    IO         *iostreams.IOStreams
    HttpClient func() (*http.Client, error)
    Config     func() (gc.Config, error)

    // 命令选项
    Option string
}

func NewCmdXxx(f *cmdutil.Factory, runF func(*XxxOptions) error) *cobra.Command {
    opts := &XxxOptions{
        IO:         f.IOStreams,
        HttpClient: f.HttpClient,
        Config:     f.Config,
    }

    cmd := &cobra.Command{
        Use:   "xxx",
        Short: "Short description",
        RunE: func(cmd *cobra.Command, args []string) error {
            if runF != nil {
                return runF(opts)
            }
            return xxxRun(opts)
        },
    }

    cmd.Flags().StringVarP(&opts.Option, "option", "o", "", "Description")

    return cmd
}

func xxxRun(opts *XxxOptions) error {
    // 1. 验证参数
    // 2. 获取依赖
    // 3. 执行业务逻辑
    // 4. 格式化输出
    return nil
}
```

### 测试模板

```go
// pkg/cmd/xxx/xxx_test.go
package xxx

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestNewCmdXxx(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {"default", []string{}, false},
        {"with option", []string{"--option", "value"}, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            factory := cmdutil.NewMockFactory()
            cmd := NewCmdXxx(factory, nil)
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

## API 客户端使用

```go
// 创建客户端
httpClient, _ := f.HttpClient()
client := api.NewClientFromHTTP(httpClient)

// REST API 调用
var user User
err := client.REST(hostname, "GET", "/api/v4/user", nil, &user)

// GraphQL 查询
query := `
    query($login: String!) {
        user(login: $login) {
            name
            email
        }
    }
`
variables := map[string]interface{}{"login": "username"}
err := client.GraphQL(hostname, query, variables, &result)
```

## 配置访问

```go
// 获取配置
cfg, _ := f.Config()

// 获取配置项
protocol := cfg.GitProtocol(hostname).Value
editor := cfg.Editor(hostname).Value

// 认证配置
authCfg := cfg.Authentication()
token, source := authCfg.ActiveToken(hostname)
user, _ := authCfg.ActiveUser(hostname)

// 设置配置
cfg.Set(hostname, "editor", "vim")
cfg.Write()
```

## 输出处理

```go
// 颜色输出
cs := opts.IO.ColorScheme()
fmt.Fprintf(opts.IO.Out, "%s Success\n", cs.Green("✓"))

// 表格输出
tp := tableprinter.New(opts.IO, tableprinter.WithHeader("ID", "TITLE"))
tp.AddField("123")
tp.AddField("Title")
tp.EndRow()
tp.Render()

// JSON 输出
if len(opts.JSONFields) > 0 {
    return jsonFields(data, opts.JSONFields, opts.IO.Out)
}
```

## 交互提示

```go
// 确认
confirmed, err := opts.Prompter.Confirm("Are you sure?", true)

// 输入
name, err := opts.Prompter.Input("Name:", "")

// 选择
index, err := opts.Prompter.Select("Choose:", "option1", []string{"option1", "option2"})

// Markdown 编辑器
body, err := opts.Prompter.MarkdownEditor("Body:", "", true)
```

## 开发优先级

### 第一阶段：核心功能

1. `gc auth login` - 认证登录
2. `gc repo clone` - 克隆仓库
3. `gc issue create/list/view` - Issue 管理
4. `gc mr create/list/view` - MR 管理
5. `gc mr review` - 代码检视

### 第二阶段：增强功能

1. `gc mr checkout/merge` - MR 检出/合并
2. `gc repo create/fork` - 仓库创建/Fork
3. `gc config` - 配置管理
4. `gc alias` - 别名管理

### 第三阶段：扩展功能

1. `gc api` - API 调用
2. `gc extension` - 扩展系统
3. 自动补全优化

## GitCode API 差异

GitCode API 基于 GitLab API，与 GitHub API 有以下主要差异：

| 功能 | GitHub | GitCode/GitLab |
|------|--------|----------------|
| Issue 编号 | #123 | #123 (iid) |
| MR/PR | pull request | merge request |
| API 版本 | v3/graphql | v4/v5 |
| 认证 | OAuth/Token | OAuth/Token |
| 端点格式 | /repos/owner/repo | /projects/owner%2Frepo |

## 常见任务

### 添加新命令

1. 在 `pkg/cmd/` 下创建目录
2. 实现命令逻辑
3. 在父命令中注册
4. 编写测试
5. 更新文档

### 添加 API 查询

1. 在 `api/queries_*.go` 添加函数
2. 定义请求/响应结构
3. 实现 API 调用
4. 添加测试

### 添加配置项

1. 在 `internal/config/config.go` 添加方法
2. 更新默认配置模板
3. 更新文档

## 测试运行

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test -run TestLogin ./pkg/cmd/auth/login/...

# 运行带覆盖率
go test -coverprofile=coverage.out ./...

# 运行集成测试
go test -tags=integration ./...
```

## 代码质量

```bash
# 格式化
go fmt ./...

# Lint
golangci-lint run ./...

# 静态检查
go vet ./...
```

## 参考资源

- [GitHub CLI 源码](https://github.com/cli/cli)
- [Cobra 文档](https://github.com/spf13/cobra)
- [GitLab API 文档](https://docs.gitlab.com/ee/api/)
- [设计文档](./issues-plan/)

## 注意事项

1. **命名规范（重要）**:
   - 命令名称统一使用 `gc`，**禁止使用 `gt`**
   - 配置目录统一使用 `~/.config/gc/`
   - 环境变量统一使用 `GC_*` 前缀（如 `GC_TOKEN`, `GC_HOST`）
   - 二进制文件命名为 `gc`
   - 如果发现文档或代码中误使用 `gt`，应立即修正为 `gc`
2. **安全性**: Token 必须使用 keyring 存储，避免明文
3. **错误处理**: 提供清晰的错误信息，指导用户如何修复
4. **兼容性**: 支持 Linux/macOS/Windows 三大平台
5. **性能**: API 请求使用缓存，避免重复请求
6. **测试**: 新功能必须有单元测试覆盖

## 提交信息规范

使用 Conventional Commits:

```
feat: add mr review command
fix: token validation error
docs: update installation guide
test: add unit tests for login
refactor: simplify config loading
```

---

**最后更新**: 2026-03-22

**维护者**: GitCode CLI Team