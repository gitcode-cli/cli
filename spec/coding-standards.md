# 编码规范

本文档定义 gitcode-cli 项目的编码规范。

## 命名规范

### 包名
- 小写、简短、有意义
- 不使用下划线或驼峰

```go
package config    // ✅ 正确
package configUtils  // ❌ 错误
package config_utils  // ❌ 错误
```

### 导出名称
- 导出的函数、类型、常量使用大驼峰（PascalCase）

```go
// ✅ 正确
func NewConfig() (*Config, error)
type UserManager struct {}
const DefaultHost = "gitcode.com"

// ❌ 错误
func newConfig() (*Config, error)  // 内部函数
func Get_config()  // 下划线
```

### 内部名称
- 内部（未导出）的函数、变量使用小驼峰（camelCase）

```go
// ✅ 正确
func parseConfig(data []byte) (*Config, error)
func getUserInfo(id int) *User

// ❌ 错误
func ParseConfig(data []byte) (*Config, error)  // 应该导出
func get_user_info(id int) *User  // 下划线
```

### 常量
- 导出常量使用大驼峰
- 未导出常量使用小驼峰

```go
const (
    DefaultHost    = "gitcode.com"  // 导出
    defaultTimeout = 30             // 内部
)
```

### 接口
- 单方法接口以 "-er" 结尾

```go
// ✅ 正确
type Reader interface { Read() }
type Writer interface { Write() }

// ❌ 错误
type Read interface { Read() }
```

## 文件结构

### import 顺序
按以下顺序分组，组间空行分隔：
1. 标准库
2. 第三方库
3. 内部包

```go
package xxx

import (
    "context"
    "fmt"
    "net/http"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"

    "gitcode.com/gitcode-cli/cli/internal/config"
    "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)
```

### 文件内容顺序
```go
package xxx

import (...)

const (...)    // 常量

type (...)     // 类型定义

func New()     // 构造函数（放前面）

func (x *Xxx) Public() {}   // 公开方法

func (x *Xxx) private() {}  // 私有方法

func helperFunc() {}        // 辅助函数
```

## 错误处理

### 简单错误
```go
var ErrNotFound = errors.New("not found")
```

### 包装错误
```go
if err != nil {
    return fmt.Errorf("failed to get config: %w", err)
}
```

### 错误信息
- 小写开头，不以句号结尾
- 描述发生了什么，不是为什么

```go
// ✅ 正确
return fmt.Errorf("failed to connect to server")
return fmt.Errorf("config file not found")

// ❌ 错误
return fmt.Errorf("Failed to connect to server.")  // 大写、句号
return fmt.Errorf("the server is down")  // 不是描述发生了什么
```

## 代码风格

### 行长度
- 最大 120 字符
- 超长行适当换行

### 函数长度
- 单个函数不超过 50 行
- 超过时应拆分为多个函数

### 注释
- 导出的函数、类型必须有注释
- 注释以函数名开头

```go
// NewConfig creates a new configuration instance.
func NewConfig() (*Config, error) {
    // ...
}

// GetUser returns the user with the given ID.
func (s *UserService) GetUser(id int) (*User, error) {
    // ...
}
```

## 代码组织

### 目录结构
```
pkg/cmd/xxx/          # 命令实现
├── xxx.go            # 主命令
├── xxx_test.go       # 单元测试
└── subcommand.go     # 子命令（如有）

internal/             # 内部包
├── config/           # 配置管理
├── authflow/         # 认证流程
└── prompter/         # 交互提示

api/                  # API 客户端
├── client.go         # 客户端
├── queries_issue.go  # Issue 相关
└── queries_pr.go     # PR 相关
```

### 避免的问题
- ❌ 循环导入
- ❌ 过深的嵌套（最多 3 层）
- ❌ 过多的函数参数（最多 5 个）
- ❌ 全局变量

## 参考

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)

---

**最后更新**: 2026-03-26