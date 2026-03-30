# 命令开发模板

本文档提供 gitcode-cli 新命令开发的模板和示例。

## 命令结构

### 文件位置
```
pkg/cmd/<command>/<command>.go
pkg/cmd/<command>/<command>_test.go
pkg/cmd/<command>/subcommand.go  # 子命令（如有）
```

## 基本模板

```go
// pkg/cmd/xxx/xxx.go
package xxx

import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/MakeNowJust/heredoc/v2"
    "github.com/spf13/cobra"

    "gitcode.com/gitcode-cli/cli/api"
    cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
    "gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

type XxxOptions struct {
    IO         *iostreams.IOStreams
    HttpClient func() (*http.Client, error)

    // Arguments
    Repository string

    // Flags
    Option string
    JSON   bool
}

func NewCmdXxx(f *cmdutil.Factory, runF func(*XxxOptions) error) *cobra.Command {
    opts := &XxxOptions{
        IO:         f.IOStreams,
        HttpClient: f.HttpClient,
    }

    cmd := &cobra.Command{
        Use:   "xxx",
        Short: "Short description",
        Long: heredoc.Doc(`
            Longer description of the command.
        `),
        Example: heredoc.Doc(`
            # Example 1
            $ gc xxx --option value

            # Example 2
            $ gc xxx -R owner/repo
        `),
        Args: cobra.NoArgs,  // 或 cobra.ExactArgs(1)
        RunE: func(cmd *cobra.Command, args []string) error {
            if runF != nil {
                return runF(opts)
            }
            return xxxRun(opts)
        },
    }

    // 添加 flags
    cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
    cmd.Flags().StringVarP(&opts.Option, "option", "o", "", "Option description")
    cmd.Flags().BoolVar(&opts.JSON, "json", false, "Output as JSON")

    return cmd
}

func xxxRun(opts *XxxOptions) error {
    cs := opts.IO.ColorScheme()

    // 1. 获取 HTTP 客户端
    httpClient, err := opts.HttpClient()
    if err != nil {
        return fmt.Errorf("failed to create HTTP client: %w", err)
    }

    // 2. 创建 API 客户端
    client := api.NewClientFromHTTP(httpClient)
    token := getEnvToken()
    if token == "" {
        return fmt.Errorf("not authenticated. Run: gc auth login")
    }
    client.SetToken(token, "environment")

    // 3. 解析仓库
    owner, repo, err := parseRepo(opts.Repository)
    if err != nil {
        return err
    }

    // 4. 执行操作
    // ... API 调用 ...

    // 5. 输出结果
    if opts.JSON {
        if err := json.NewEncoder(opts.IO.Out).Encode(map[string]string{"status": "ok"}); err != nil {
            return fmt.Errorf("failed to encode JSON output: %w", err)
        }
        return nil
    }

    fmt.Fprintf(opts.IO.Out, "%s Success\n", cs.Green("✓"))

    return nil
}

func parseRepo(repo string) (string, string, error) {
    if repo == "" {
        return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
    }
    parts := strings.Split(repo, "/")
    if len(parts) != 2 {
        return "", "", fmt.Errorf("invalid repository format: %s", repo)
    }
    return parts[0], parts[1], nil
}

func getEnvToken() string {
    if token := os.Getenv("GC_TOKEN"); token != "" {
        return token
    }
    return os.Getenv("GITCODE_TOKEN")
}
```

## API 客户端用法

### 创建客户端
```go
httpClient, _ := f.HttpClient()
client := api.NewClientFromHTTP(httpClient)
client.SetToken(token, "environment")
```

### GET 请求
```go
var user User
err := client.Get("/repos/"+owner+"/"+repo+"/user", &user)
```

### POST 请求
```go
var result Result
err := client.Post("/repos/"+owner+"/"+repo+"/issues", &CreateOptions{
    Title: "Title",
    Body:  "Body",
}, &result)
```

### PATCH 请求
```go
err := client.Patch("/repos/"+owner+"/"+repo+"/issues/"+number, &UpdateOptions{
    Title: "New Title",
}, &result)
```

### DELETE 请求
```go
err := client.Delete("/repos/"+owner+"/"+repo+"/issues/"+number)
```

## 配置访问

```go
cfg, _ := f.Config()

// 获取 Git 协议
protocol := cfg.GitProtocol(hostname).Value

// 获取认证配置
authCfg := cfg.Authentication()
token, source := authCfg.ActiveToken(hostname)
```

## 输出处理

### 颜色输出
```go
cs := opts.IO.ColorScheme()

// 成功消息
fmt.Fprintf(opts.IO.Out, "%s Success\n", cs.Green("✓"))

// 错误消息
fmt.Fprintf(opts.IO.ErrOut, "%s Error\n", cs.Red("!"))

// 粗体
fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold(title))
```

### 表格输出
```go
tp := tableprinter.New(opts.IO, tableprinter.WithHeader("ID", "TITLE", "STATE"))

for _, item := range items {
    tp.AddField(item.ID)
    tp.AddField(item.Title)
    tp.AddField(item.State)
    tp.EndRow()
}

tp.Render()
```

## 交互提示

```go
err := cmdutil.ConfirmOrAbort(cmdutil.ConfirmOptions{
    IO:       opts.IO,
    Yes:      opts.Yes,
    Expected: "owner/repo",
    Prompt:   "Type the repository name to confirm: ",
})
```

## 测试模板

```go
// xxx_test.go
package xxx

import (
    "testing"

    cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdXxx(t *testing.T) {
    tests := []struct {
        name    string
        args    []string
        wantErr bool
    }{
        {
            name:    "normal case",
            args:    []string{"--option", "value"},
            wantErr: false,
        },
        {
            name:    "error case",
            args:    []string{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f := cmdutil.TestFactory()
            cmd := NewCmdXxx(f, func(opts *XxxOptions) error {
                return nil
            })
            cmd.SetArgs(tt.args)

            err := cmd.Execute()
            if (err != nil) != tt.wantErr {
                t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

**最后更新**: 2026-03-26
