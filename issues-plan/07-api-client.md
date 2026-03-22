# API 客户端需求

本文档详细描述 gitcode-cli API 客户端的设计需求和验收标准。

## 模块概述

API 客户端模块负责与 GitCode API 进行通信，提供 REST API 封装、认证中间件、缓存机制、重试机制和错误处理。

### 目录结构

```
api/
├── client.go           # API 客户端核心
├── http_client.go      # HTTP 客户端封装
├── queries_repo.go     # 仓库相关查询
├── queries_issue.go    # Issue 相关查询
├── queries_pr.go        # PR 相关查询
├── queries_user.go     # 用户相关查询
└── query_builder.go    # GraphQL 查询构建
```

---

## API-001: REST API 封装

### 功能描述

封装 GitCode REST API 调用，提供统一的 API 调用接口。

### 设计要求

```go
// api/client.go
type Client struct {
    httpClient *http.Client
    hostname   string
}

// REST 调用方法
func (c *Client) REST(hostname, method, path string, body io.Reader, response interface{}) error

// 便捷方法
func (c *Client) Get(hostname, path string, response interface{}) error
func (c *Client) Post(hostname, path string, body, response interface{}) error
func (c *Client) Put(hostname, path string, body, response interface{}) error
func (c *Client) Patch(hostname, path string, body, response interface{}) error
func (c *Client) Delete(hostname, path string) error
```

### 功能特性

- 支持 GET/POST/PUT/PATCH/DELETE 方法
- 自动添加认证头
- JSON 序列化/反序列化
- 统一错误处理

### 验收标准

- [ ] 支持所有 HTTP 方法
- [ ] 正确处理 JSON 序列化
- [ ] 正确处理响应状态码
- [ ] 提供清晰的错误信息

### 测试用例映射

- 参考 `gc-api-doc/test/test_*.py`

---

## API-002: 认证中间件

### 功能描述

HTTP 请求中间件，自动注入认证信息。

### 设计要求

```go
// api/http_client.go

// 认证中间件
func AddAuthTokenHeader(rt http.RoundTripper, cfg tokenGetter) http.RoundTripper {
    return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
        if req.Header.Get("Authorization") == "" {
            hostname := getHost(req)
            if token, _ := cfg.ActiveToken(hostname); token != "" {
                req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
            }
        }
        return rt.RoundTrip(req)
    })
}

// 支持的认证方式
// 1. Authorization: Bearer {token}
// 2. PRIVATE-TOKEN: {token}
// 3. access_token 查询参数
```

### Token 优先级

1. `GC_TOKEN` 环境变量
2. `GITCODE_TOKEN` 环境变量
3. Keyring 存储
4. 配置文件存储

### 验收标准

- [ ] 自动注入 Authorization 头
- [ ] 支持多种认证方式
- [ ] 支持环境变量 Token
- [ ] 支持从配置读取 Token

---

## API-003: 缓存机制

### 功能描述

缓存 API 响应以提高性能和减少请求。

### 设计要求

```go
// api/cache.go

type Cache struct {
    store map[string]cacheEntry
    mu    sync.RWMutex
}

type cacheEntry struct {
    response  []byte
    expiresAt time.Time
    etag      string
}

// 缓存中间件
func CacheMiddleware(rt http.RoundTripper, ttl time.Duration) http.RoundTripper

// 缓存键生成
func cacheKey(req *http.Request) string
```

### 缓存策略

| API 类型 | 缓存时间 | 说明 |
|----------|----------|------|
| 用户信息 | 5 分钟 | 相对稳定 |
| 仓库信息 | 1 分钟 | 可能变化 |
| Issue/PR 列表 | 30 秒 | 频繁变化 |
| Issue/PR 详情 | 1 分钟 | 可能变化 |

### 验收标准

- [ ] 正确缓存 GET 请求
- [ ] 支持 ETag 验证
- [ ] 支持缓存过期
- [ ] 支持 Cache-Control 头

---

## API-004: 重试机制

### 功能描述

自动重试失败的 API 请求。

### 设计要求

```go
// api/retry.go

type RetryConfig struct {
    MaxRetries  int           // 最大重试次数
    InitialWait time.Duration // 初始等待时间
    MaxWait     time.Duration // 最大等待时间
    Multiplier  float64       // 退避乘数
}

// 重试中间件
func RetryMiddleware(rt http.RoundTripper, cfg RetryConfig) http.RoundTripper
```

### 重试策略

| 状态码 | 行为 |
|--------|------|
| 429 (Rate Limit) | 等待 Retry-After 头指定的时间后重试 |
| 500, 502, 503 | 指数退避重试 |
| 401 | 不重试，返回认证错误 |
| 其他 | 不重试 |

### 指数退避

```
第1次重试: 等待 1 秒
第2次重试: 等待 2 秒
第3次重试: 等待 4 秒
最大等待: 30 秒
```

### 验收标准

- [ ] 支持自动重试
- [ ] 支持指数退避
- [ ] 处理 429 状态码
- [ ] 最大重试次数限制

---

## API-005: 错误处理

### 功能描述

统一处理 API 错误，提供清晰的错误信息。

### 设计要求

```go
// api/errors.go

type APIError struct {
    StatusCode  int
    Message     string
    Errors      []FieldError
    RequestID   string
    RateLimit   *RateLimitInfo
}

type FieldError struct {
    Field   string
    Message string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}
```

### HTTP 状态码处理

| 状态码 | 说明 | 处理方式 |
|--------|------|----------|
| 400 | 请求格式错误 | 显示错误详情 |
| 401 | 未认证 | 提示登录 |
| 403 | 无权限 | 显示权限信息 |
| 404 | 资源不存在 | 显示友好提示 |
| 422 | 验证错误 | 显示字段错误 |
| 429 | 请求过多 | 显示重试时间 |
| 500+ | 服务器错误 | 建议稍后重试 |

### 验收标准

- [ ] 正确解析错误响应
- [ ] 提供清晰的错误信息
- [ ] 显示字段级验证错误
- [ ] 处理 Rate Limit

### 测试用例映射

- 参考 `gc-api-doc/test/test_error_codes.py`

---

## API 端点配置

### GitCode API 端点

```yaml
# 默认配置
api:
  base_url: https://api.gitcode.com/api/v5
  timeout: 30s
  retries: 3

# 端点映射
endpoints:
  user: /user
  repos: /repos
  issues: /repos/{owner}/{repo}/issues
  pulls: /repos/{owner}/{repo}/pulls
```

### 环境变量覆盖

| 环境变量 | 说明 |
|----------|------|
| `GC_API_URL` | 覆盖 API 基础 URL |
| `GC_TIMEOUT` | 覆盖请求超时 |

---

## 请求频率限制

### GitCode API 限制

| 类型 | 限制 |
|------|------|
| 每分钟 | 50 次 |
| 每小时 | 4000 次 |

### 客户端处理

- 跟踪剩余请求次数
- 接近限制时警告
- 达到限制时等待

---

## 相关文档

- [gc-api-doc/doc/00-overview.md](../../gc-api-doc/doc/00-overview.md)
- [gc-api-doc/doc/14-error-codes.md](../../gc-api-doc/doc/14-error-codes.md)
- [gc-api-doc/test/test_error_codes.py](../../gc-api-doc/test/test_error_codes.py)

---

**最后更新**: 2026-03-22