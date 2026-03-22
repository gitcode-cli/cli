---
name: gitcode-cmd-generator
description: |
  Generate command code templates and test files for gitcode-cli project.

  Use this skill when:
  - Creating a new command for gitcode-cli (e.g., "create a new gc release command")
  - Adding a subcommand to existing command (e.g., "add list subcommand to release")
  - Generating command template (e.g., "generate command template for workflow list")
  - Creating test files for commands (e.g., "create test for auth logout")
  - Scaffolding CLI commands following project conventions

  This skill generates production-ready Go code following gitcode-cli patterns with Cobra framework, Factory dependency injection, and IOStreams for output handling.
---

# GitCode CLI Command Generator

Generate standard command code templates and test files for the gitcode-cli project.

## Project Context

| 项目 | 值 |
|------|-----|
| 命令名 | `gc` |
| 语言 | Go 1.22+ |
| 框架 | Cobra |
| 模块路径 | `github.com/gitcode-com/gitcode-cli` |
| 配置目录 | `~/.config/gc/` |
| 环境变量前缀 | `GC_*` |

## Directory Structure

```
pkg/cmd/<category>/
├── <category>.go          # Parent command
├── <action>/
│   ├── <action>.go        # Subcommand implementation
│   └── <action>_test.go   # Subcommand tests
```

## Command Types

### 1. Parent Command Template

For a new command category (e.g., `gc release`):

```go
// Package <category> implements the <category> command
package <category>

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/<category>/<subcmd1>"
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/<category>/<subcmd2>"
)

// NewCmd<Category> creates the <category> command
func NewCmd<Category>(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "<category> <command>",
		Short: "<Short description>",
		Long: heredoc.Doc(`
			<Detailed description>.
		`),
		Example: heredoc.Doc(`
			# Example 1
			$ gc <category> <subcmd> --option value

			# Example 2
			$ gc <category> <subcmd>
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
	}

	cmd.AddCommand(<subcmd1>.NewCmd<Subcmd1>(f, nil))
	cmd.AddCommand(<subcmd2>.NewCmd<Subcmd2>(f, nil))

	return cmd
}
```

### 2. Action Command Template (create/delete/close/reopen)

```go
// Package <action> implements the <category> <action> command
package <action>

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type <Action>Options struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Number     int

	// Flags
	Repository string
}

// NewCmd<Action> creates the <action> command
func NewCmd<Action>(f *cmdutil.Factory, runF func(*<Action>Options) error) *cobra.Command {
	opts := &<Action>Options{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "<action> [<number>]",
		Short: "<Short description>",
		Long: heredoc.Doc(`
			<Detailed description>.
		`),
		Example: heredoc.Doc(`
			# Example 1
			$ gc <category> <action> 123

			# Example 2
			$ gc <category> <action> 123 -R owner/repo
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Parse number from args[0]
			}

			if runF != nil {
				return runF(opts)
			}
			return <action>Run(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")

	return cmd
}

func <action>Run(opts *<Action>Options) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// TODO: Implement API call

	fmt.Fprintf(opts.IO.Out, "%s <Message>\n", cs.Green("✓"))
	return nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
```

### 3. List Command Template

```go
// Package list implements the <category> list command
package list

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Flags
	Repository string
	State      string
	Limit      int
}

// NewCmdList creates the list command
func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List <items>",
		Long: heredoc.Doc(`
			List <items> in a GitCode repository.
		`),
		Example: heredoc.Doc(`
			# List <items>
			$ gc <category> list

			# List <items> in a specific repository
			$ gc <category> list -R owner/repo

			# List with filters
			$ gc <category> list --state open --limit 20
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().StringVarP(&opts.State, "state", "s", "open", "Filter by state (open/closed/all)")
	cmd.Flags().IntVarP(&opts.Limit, "limit", "L", 30, "Maximum number of items to list")

	return cmd
}

func listRun(opts *ListOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// TODO: Implement API call and table output
	return nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
```

### 4. View Command Template

```go
// Package view implements the <category> view command
package view

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/gitcode-com/gitcode-cli/api"
	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
	"github.com/gitcode-com/gitcode-cli/pkg/iostreams"
)

type ViewOptions struct {
	IO         *iostreams.IOStreams
	HttpClient func() (*http.Client, error)

	// Arguments
	Number     int

	// Flags
	Repository string
	Web        bool
}

// NewCmdView creates the view command
func NewCmdView(f *cmdutil.Factory, runF func(*ViewOptions) error) *cobra.Command {
	opts := &ViewOptions{
		IO:         f.IOStreams,
		HttpClient: f.HttpClient,
	}

	cmd := &cobra.Command{
		Use:   "view [<number>]",
		Short: "View <item> details",
		Long: heredoc.Doc(`
			View details of a <item>.
		`),
		Example: heredoc.Doc(`
			# View <item>
			$ gc <category> view 123

			# View in browser
			$ gc <category> view 123 --web
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				// Parse number from args[0]
			}

			if runF != nil {
				return runF(opts)
			}
			return viewRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository (owner/repo)")
	cmd.Flags().BoolVarP(&opts.Web, "web", "w", false, "Open in browser")

	return cmd
}

func viewRun(opts *ViewOptions) error {
	cs := opts.IO.ColorScheme()

	httpClient, err := opts.HttpClient()
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	client := api.NewClientFromHTTP(httpClient)
	token := getEnvToken()
	if token == "" {
		return fmt.Errorf("not authenticated. Run: gc auth login")
	}
	client.SetToken(token, "environment")

	// TODO: Implement API call

	fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold("<Title>"))
	return nil
}

func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
```

## Test File Template

```go
package <action>

import (
	"testing"

	cmdutil "github.com/gitcode-com/gitcode-cli/pkg/cmdutil"
)

func TestNewCmd<Action>(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "basic test",
			args:    []string{"<required-args>"},
			wantErr: false,
		},
		{
			name:    "with flags",
			args:    []string{"<required-args>", "--flag", "value"},
			wantErr: false,
		},
		{
			name:    "missing required",
			args:    []string{},
			wantErr: false, // Command runs, error in run function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmd<Action>(f, func(opts *<Action>Options) error {
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

## Code Conventions

### Import Order
1. Standard library (e.g., `fmt`, `net/http`, `os`)
2. Third-party libraries (e.g., `github.com/spf13/cobra`)
3. Internal packages (e.g., `github.com/gitcode-com/gitcode-cli/api`)

### Naming Conventions

| 元素 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写，简短 | `create`, `list`, `view` |
| Options结构体 | 大驼峰 | `CreateOptions`, `ListOptions` |
| 构造函数 | `NewCmd<Action>` | `NewCmdCreate`, `NewCmdList` |
| Run函数 | `<action>Run` | `createRun`, `listRun` |
| 常量 | 大驼峰 | `DefaultHost` |

### Common Flags

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--repo` | `-R` | Repository (owner/repo) |
| `--title` | `-t` | Title |
| `--body` | `-b` | Body/description |
| `--label` | `-l` | Labels (comma-separated) |
| `--state` | `-s` | State filter |
| `--limit` | `-L` | Result limit |
| `--web` | `-w` | Open in browser |

### Color Output

```go
cs := opts.IO.ColorScheme()
fmt.Fprintf(opts.IO.Out, "%s Success message\n", cs.Green("✓"))
fmt.Fprintf(opts.IO.ErrOut, "%s Error message\n", cs.Red("✗"))
fmt.Fprintf(opts.IO.Out, "%s\n", cs.Bold("Title"))
```

## Registration

After creating a new command, register it in the root command:

```go
// pkg/cmd/root/root.go
import (
	"github.com/gitcode-com/gitcode-cli/pkg/cmd/<category>"
)

// In NewRootCmd:
cmd.AddCommand(<category>.NewCmd<Category>(f))
```

## Commit Requirements

- 单次提交不超过 **800 行**
- 及时提交，避免大量代码积压
- 使用 Conventional Commits: `feat:`, `fix:`, `docs:`, `test:`

## Reference Files

For detailed API patterns, see `references/api-patterns.md`.