# Security/Auth Agent 任务书

## 目标

审视认证、凭证处理、危险写操作和确认语义，只产出可独立提交 issue 的安全或认证问题。

## 本 agent 只看什么

- token 来源与优先级是否统一
- 本地认证状态是否被业务命令绕过
- 凭证是否被放进 URL、日志、错误信息或本地持久化路径
- 删除、合并、覆盖等高风险写操作是否统一使用确认语义
- 安全相关接口语义和实际实现是否一致

## 必看目录

- `internal/config/`
- `pkg/cmdutil/`
- `pkg/cmd/`
- `api/`

## 必看文件

- `internal/config/auth_config.go`
- `internal/config/config.go`
- `pkg/cmdutil/auth.go`
- `pkg/cmdutil/confirm.go`
- `api/client.go`

## 建议重点命令族

- `pkg/cmd/auth/`
- `pkg/cmd/pr/merge/`
- `pkg/cmd/repo/sync/`
- `pkg/cmd/release/`
- `pkg/cmd/issue/comment/`
- 任何带 `delete`、`merge`、`sync`、`upload`、`edit` 的命令

## 必答问题

1. `gc auth login` 写入的认证状态，是否被业务命令真正消费？
2. 是否存在把 token 放进 URL 或命令参数的路径？
3. 危险写操作是否统一经过 `cmdutil.ConfirmOrAbort()` 或等价确认？
4. 接口暴露的安全开关、参数、存储策略，是否存在“名义支持、实际上无效”？
5. 哪些问题属于 CLI 自身缺陷，哪些属于上游 API 限制？

## 明确算 issue 的问题类型

- 本地已登录 token 不被业务命令使用
- token 被拼接到 URL 或其他高暴露面
- 高风险状态变更缺失统一确认
- 安全接口或参数名义存在、实际被忽略

## 不要提交 issue 的内容

- 没有实际风险放大的理论担忧
- 已经被当前实现充分防护的问题
- 只涉及一般代码整洁度的问题

## 输出格式

```markdown
## Finding N

- 标题建议:
- 严重度: high / medium / low
- 类型: security
- 现象:
- 根因:
- 影响:
- 涉及目录或文件:
- 复现命令或核验方式:
- 是否疑似已有远端 issue:
```

## 通过标准

- 每条 finding 都必须明确写出风险暴露面
- 如果问题依赖用户本地环境，必须说明成立条件
- 如果问题本质是上游 API 限制，不得错误归类为 CLI 安全漏洞

---

**最后更新**: 2026-04-03
