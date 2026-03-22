# CLAUDE.md - AI 辅助开发指南

> 项目概述和功能介绍请参阅 [README.md](./README.md)

本文档为 Claude Code 提供 gitcode-cli 项目开发指导。

## 核心信息

| 项目 | 值 |
|------|-----|
| 命令名 | `gc` |
| 语言 | Go 1.21+ |
| 框架 | Cobra |
| 配置目录 | `~/.config/gc/` |
| 环境变量前缀 | `GC_*` |

## 目录结构

```
gitcode-cli/
├── cmd/gc/           # 程序入口
├── pkg/cmd/          # 命令实现
├── pkg/cmdutil/      # Factory、工具函数
├── pkg/iostreams/    # IO 流管理
├── internal/
│   ├── config/       # 配置管理
│   ├── authflow/     # 认证流程
│   ├── keyring/      # 安全存储
│   └── prompter/     # 交互提示
├── api/              # API 客户端
└── git/              # Git 操作
```

## 编码规范

### 命名

```go
// 包名：小写简短
package config

// 导出：大驼峰
func NewConfig() (*Config, error)

// 内部：小驼峰
func parseConfig(data []byte) (*Config, error)

// 常量：大驼峰
const DefaultHost = "gitcode.com"
```

### 文件结构

```go
package xxx

import (
    "context"           // 1. 标准库
    "github.com/spf13/cobra"  // 2. 第三方库
    "github.com/gitcode-cli/cli/internal/config"  // 3. 内部包
)

const ()    // 常量
type ()     // 类型
func New()  // 构造函数
func (x *Xxx) Public() {}  // 公开方法
func (x *Xxx) private() {} // 私有方法
```

### 错误处理

```go
// 简单错误
var ErrNotFound = errors.New("not found")

// 包装错误
return fmt.Errorf("failed to get config: %w", err)
```

## 命令开发模板

```go
// pkg/cmd/xxx/xxx.go
package xxx

import (
    "github.com/spf13/cobra"
    cmdutil "github.com/gitcode-cli/cli/pkg/cmdutil"
)

type XxxOptions struct {
    IO         *iostreams.IOStreams
    HttpClient func() (*http.Client, error)
    Config     func() (gc.Config, error)
    Option     string
}

func NewCmdXxx(f *cmdutil.Factory, runF func(*XxxOptions) error) *cobra.Command {
    opts := &XxxOptions{
        IO: f.IOStreams,
        HttpClient: f.HttpClient,
        Config: f.Config,
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
    // 1. 验证参数 2. 获取依赖 3. 执行逻辑 4. 格式化输出
    return nil
}
```

## API 客户端

```go
httpClient, _ := f.HttpClient()
client := api.NewClientFromHTTP(httpClient)

// REST 调用
var user User
err := client.REST(hostname, "GET", "/api/v5/user", nil, &user)
```

## 配置访问

```go
cfg, _ := f.Config()
protocol := cfg.GitProtocol(hostname).Value
authCfg := cfg.Authentication()
token, source := authCfg.ActiveToken(hostname)
```

## 输出处理

```go
// 颜色
cs := opts.IO.ColorScheme()
fmt.Fprintf(opts.IO.Out, "%s Success\n", cs.Green("✓"))

// 表格
tp := tableprinter.New(opts.IO, tableprinter.WithHeader("ID", "TITLE"))
tp.AddField("123")
tp.EndRow()
tp.Render()
```

## 交互提示

```go
confirmed, _ := opts.Prompter.Confirm("Are you sure?", true)
name, _ := opts.Prompter.Input("Name:", "")
index, _ := opts.Prompter.Select("Choose:", "opt1", []string{"opt1", "opt2"})
```

## 测试

```bash
go test ./...                          # 所有测试
go test -run TestLogin ./pkg/cmd/auth/...  # 特定测试
go test -coverprofile=coverage.out ./...   # 覆盖率
go test -tags=integration ./...            # 集成测试
```

## GitCode API 差异

| 功能 | GitHub | GitCode |
|------|--------|---------|
| MR/PR | pull request | merge request |
| API 版本 | v3/graphql | v5 |
| 端点 | /repos/owner/repo | /projects/owner%2Frepo |

## 开发优先级

1. **P0**: `auth login`, `repo clone`, `issue create/list/view`, `mr create/list/view`, `mr review`
2. **P1**: `mr checkout/merge`, `repo create/fork`, `config`
3. **P2**: `api`, `extension`, 自动补全

## 重要注意事项

1. **命名**: 命令名 `gc`，禁止使用 `gt`；环境变量 `GC_*`
2. **安全**: Token 必须使用 keyring 存储
3. **错误**: 提供清晰的错误信息和修复建议
4. **测试**: 新功能必须有单元测试

## 提交规范

### 提交信息

使用 Conventional Commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`

### 提交要求

- **单次提交限制**: 每次代码提交不超过 **800 行**
- **及时提交**: 完成一个功能点或修复后立即提交，避免大量代码积压
- **原子提交**: 每个提交应是一个独立的、完整的功能或修复

## 参考文档

- [需求文档](./issues-plan/) - 完整需求和里程碑
- [架构设计](./issues-plan/02-architecture.md)
- [API 客户端](./issues-plan/07-api-client.md)
- [GitHub CLI 源码](https://github.com/cli/cli)

---

**最后更新**: 2026-03-22