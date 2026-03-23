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
| PR | pull request | pull request |
| API 版本 | v3/graphql | v5 |
| 端点 | /repos/owner/repo | /projects/owner%2Frepo |

## 开发优先级

1. **P0**: `auth login`, `repo clone`, `issue create/list/view`, `pr create/list/view`, `pr review`
2. **P1**: `pr checkout/merge`, `repo create/fork`, `config`
3. **P2**: `api`, `extension`, 自动补全

## 重要注意事项

1. **命名**: 命令名 `gc`，禁止使用 `gt`；环境变量 `GC_*`
2. **安全**: Token 必须使用 keyring 存储
3. **错误**: 提供清晰的错误信息和修复建议
4. **测试**: 新功能必须有单元测试

## 敏感信息保护（重要！）

**严格遵守以下规则，违反将导致严重安全问题！**

### Token 处理要求

1. **禁止写入文件**: Token 绝对不能写入配置文件、源代码或任何持久化存储
2. **内存存储**: Token 仅在内存中保存，程序结束后自动清除
3. **禁止提交**: Token 绝对不能提交到版本控制系统
4. **测试时传递**: 测试时通过环境变量 `GC_TOKEN` 或命令行参数传递

### 测试配置

- 测试组织: `infra-test`
- API 基础 URL: `https://api.gitcode.com/api/v5`
- Token 来源: 环境变量或运行时输入

### 测试仓库限制（重要！）

**严格遵守**: 只能使用以下指定的测试仓库，禁止随意使用其他仓库进行测试。

**允许使用的测试仓库**:
- `infra-test/gctest1`
- `gitcode-cli/cli`

**禁止行为**:
- 不要使用个人仓库测试
- 不要使用其他组织或用户的仓库测试
- 测试前确认仓库在允许列表中

### 代码审查检查项

提交前必须确认：
- [ ] 没有硬编码的 Token 或密钥
- [ ] 配置文件中不包含敏感信息
- [ ] .gitignore 已忽略敏感文件
- [ ] 测试代码不包含真实 Token

## 进度跟踪

### 状态定义

| 状态 | 图标 | 说明 |
|------|------|------|
| 待开发 | 📋 | 需求已定义，等待开发 |
| 开发中 | 🚧 | 正在开发中 |
| 已完成 | ✅ | 功能已实现并通过验收 |
| 暂停 | ⏸️ | 开发暂停 |
| 取消 | ❌ | 需求已取消 |

### 状态刷新要求

**重要**: 每次完成以下操作后，必须更新 `issues-plan/PROGRESS.md`：

1. **开始开发任务时**: 将任务状态从 📋 改为 🚧
2. **完成任务时**: 将任务状态改为 ✅，更新完成日期
3. **提交代码后**: 更新提交记录表
4. **里程碑完成时**: 更新里程碑总览表

### 进度文件

- 总进度表: `issues-plan/PROGRESS.md`
- 需求清单: `issues-plan/01-requirements-overview.md`
- 里程碑详情: `issues-plan/milestones/`

## 提交规范

### 提交信息

使用 Conventional Commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`

### 提交要求

- **单次提交限制**: 每次代码提交不超过 **800 行**
- **及时提交**: 完成一个功能点或修复后立即提交，避免大量代码积压
- **原子提交**: 每个提交应是一个独立的、完整的功能或修复
- **立即推送**: 每次提交后立即推送到远端，确保代码同步

## 开发工作流程（重要！）

**严格遵守以下流程，违反将导致代码管理混乱！**

### 流程步骤

1. **提交 Issue**: 发现 BUG 或需要新特性后，首先在项目中创建 Issue
2. **打标签**: Issue 创建后立即打上合适的标签（如 `bug`、`enhancement`）
3. **创建开发分支**: 从 main 分支创建对应类型的分支
4. **在分支开发**: 不直接在 main 分支修改
5. **编写测试用例**: 为新功能或修复编写单元测试
6. **本地测试**: 运行单元测试确保功能正常
7. **实际命令测试**: 使用 `gc` 命令在测试仓库进行实际功能验证
8. **提交 PR**: 开发完成后，创建 PR 合并到 main 分支
9. **关联 Issue**: PR 描述中必须关联对应的 Issue（如 `Fixes #123` 或 `Closes #123`）
10. **审查 PR 提交**: 检查 PR 的代码变更，确保质量
11. **Issue 评论**: 在 Issue 中添加完成说明（如何解决的、做了哪些改动）
12. **PR 审查评论**: 在 PR 中提交审查评论，说明改动内容和测试结果
13. **关闭 Issue**: 审查通过后关闭关联的 Issue
14. **合并 PR**: 确认所有测试通过后合并 PR

### 分支命名规范

| 类型 | 命名格式 | 示例 |
|------|----------|------|
| BUG 修复 | `bugfix/issue-<number>` | `bugfix/issue-3` |
| 新特性 | `feature/issue-<number>` | `feature/issue-5` |

### 标签使用规范

| 标签 | 使用场景 |
|------|----------|
| `bug` | 错误修复 |
| `enhancement` | 功能增强/新特性 |
| `documentation` | 文档更新 |
| `help wanted` | 需要帮助 |
| `question` | 需要讨论 |

### 测试要求

#### 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定模块测试
go test ./pkg/cmd/issue/...

# 运行特定测试用例
go test -run TestLabelCmd ./pkg/cmd/issue/label/...

# 查看测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### 测试用例规范
- 每个新命令必须有对应的测试文件（如 `label_test.go`）
- 测试用例应覆盖：正常流程、边界条件、错误处理
- 测试文件放在与源文件相同的目录

#### 实际命令测试（重要！）

**单元测试无法覆盖所有场景，必须进行实际命令测试！**

测试步骤：
1. 在测试仓库（`infra-test/gctest1`）进行实际命令测试
2. 验证命令的输入输出是否符合预期
3. 检查 API 调用是否正确

```bash
# 设置 Token
export GC_TOKEN=your_token

# 示例：测试 PR edit 命令
gc pr edit 2 --title "Test title" -R infra-test/gctest1
gc pr view 2 -R infra-test/gctest1

# 示例：测试 PR review 命令
gc pr review 2 --approve -R infra-test/gctest1

# 示例：测试 issue label 命令
gc issue label 5 --add enhancement -R gitcode-cli/cli
```

#### 测试文件模板
```go
// xxx_test.go
package xxx

import (
	"testing"
)

func TestXxxCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "normal case",
			args:    []string{"--flag", "value"},
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
			// 测试逻辑
		})
	}
}
```

### 示例流程

```bash
# 1. 确保在 main 分支并更新
git checkout main
git pull

# 2. 创建开发分支（BUG修复用 bugfix，新特性用 feature）
git checkout -b feature/issue-5

# 3. 开发代码
# ... 编写代码 ...

# 4. 编写测试用例
# ... 创建 xxx_test.go ...

# 5. 运行单元测试
go test ./pkg/cmd/issue/label/...
# 确保测试通过
# ok  gitcode.com/gitcode-cli/cli/pkg/cmd/issue/label  0.123s

# 6. 提交代码
git add .
git commit -m "feat(issue): add label command"

# 7. 推送分支
git push -u origin feature/issue-5

# 8. 创建 PR（关联 Issue）
gc pr create --title "feat: add issue label command" --body "Closes #5" --base main -R gitcode-cli/cli

# 9. 给 Issue 打标签
gc issue label 5 --add enhancement -R gitcode-cli/cli

# 10. 实际命令测试（在测试仓库验证功能）
export GC_TOKEN=your_token
gc issue label 1 --add bug -R infra-test/gctest1
gc issue label 1 --list -R infra-test/gctest1

# 11. 审查 PR 提交
gc pr view <pr_number> -R gitcode-cli/cli
gc pr diff <pr_number> -R gitcode-cli/cli

# 12. 在 Issue 中添加完成评论
gc issue comment 5 --body "已完成实现：\n- 新增 gc issue label 命令\n- 支持添加、移除、列出标签\n- 单元测试通过\n- 实际命令测试通过" -R gitcode-cli/cli

# 13. 在 PR 中提交审查评论
gc pr review <pr_number> --comment "## 审查结果\n\n### 改动内容\n- 新增 issue label 命令\n\n### 测试结果\n- [x] 单元测试通过\n- [x] 实际命令测试通过" -R gitcode-cli/cli

# 14. 关闭 Issue
gc issue close 5 -R gitcode-cli/cli

# 15. 合并 PR
gc pr merge <pr_number> -R gitcode-cli/cli

# 16. 拉取最新代码
git checkout main && git pull
```

### 完整工作流检查清单

开发完成后必须确认：

- [ ] 单元测试全部通过 (`go test ./...`)
- [ ] 在测试仓库进行实际命令测试
- [ ] Issue 已打标签
- [ ] PR 已创建并关联 Issue
- [ ] PR 提交已审查
- [ ] Issue 已添加完成评论
- [ ] PR 已提交审查评论
- [ ] Issue 已关闭
- [ ] PR 已合并

### 禁止行为

- ❌ 直接在 main 分支开发
- ❌ 不创建 Issue 直接开发
- ❌ Issue 创建后不打标签
- ❌ PR 不关联 Issue
- ❌ 未编写测试用例就提交 PR
- ❌ 单元测试未通过就提交 PR
- ❌ 未进行实际命令测试就合并 PR
- ❌ 未审查 PR 提交就合并
- ❌ 未添加 Issue 评论就关闭 Issue
- ❌ 未提交 PR 审查评论就合并 PR

## 参考文档

- [需求文档](./issues-plan/) - 完整需求和里程碑
- [架构设计](./issues-plan/02-architecture.md)
- [API 客户端](./issues-plan/07-api-client.md)
- [GitHub CLI 源码](https://github.com/cli/cli)

### API 开发参考（重要）

开发 API 相关功能时，必须参考以下验证过的文档：

- **https://gitcode.com/afly-infra/gc-api-doc/tree/main/test/** - 所有 API 测试用例，已验证通过
- **https://gitcode.com/afly-infra/gc-api-doc/tree/main/doc/** - API 说明文档

这些文档包含正确的 API 端点、请求格式、响应格式和认证方式。

---

**最后更新**: 2026-03-23