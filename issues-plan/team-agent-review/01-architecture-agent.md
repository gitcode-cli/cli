# Architecture Agent 任务书

## 目标

审视当前主干代码里的架构和分层问题，只产出可以独立提交 issue 的明确问题。

## 本 agent 只看什么

- 命令层是否大量重复初始化 glue code
- `cmdutil` / `Factory` 是否真正成为共享入口
- API client、repo 解析、auth 解析是否在多个命令族中漂移
- 当前仓库上下文推断是否只在局部落地

## 必看目录

- `pkg/cmd/`
- `pkg/cmdutil/`
- `api/`
- `internal/config/`

## 必看文件

- `pkg/cmd/root/root.go`
- `pkg/cmdutil/factory.go`
- `pkg/cmdutil/repo.go`
- `pkg/cmdutil/auth.go`
- `api/client.go`
- `internal/config/auth_config.go`

## 建议抽样命令族

- `pkg/cmd/pr/`
- `pkg/cmd/issue/`
- `pkg/cmd/release/`
- `pkg/cmd/commit/`
- `pkg/cmd/repo/`

## 必答问题

1. 当前命令层是否存在大面积重复的 `HTTP client -> API client -> token -> repo` 启动序列？
2. `cmdutil.Factory` 是否已经成为命令初始化的唯一入口？如果没有，分裂点在哪里？
3. 统一的 repo 解析能力是否在所有资源域一致落地？
4. 统一的认证模型是否被命令层绕过？
5. 是否存在系统级默认行为缺陷，而不是单个命令 bug？

## 明确算 issue 的问题类型

- 默认工厂或共享层存在系统级错误配置
- 同一类命令在不同资源域出现明显不一致的初始化语义
- 共享抽象已经存在，但业务命令大面积绕开，导致行为分裂
- 当前仓库上下文推断能力只在部分命令生效，形成不可预测行为

## 不要提交 issue 的内容

- 纯代码风格建议
- “未来可以加一层 service” 这类空泛建议
- 没有明确影响面的抽象设计偏好

## 输出格式

```markdown
## Finding N

- 标题建议:
- 严重度: high / medium / low
- 类型: architecture
- 现象:
- 根因:
- 影响:
- 涉及目录或文件:
- 复现命令或核验方式:
- 是否疑似已有远端 issue:
```

## 通过标准

- 每条 finding 都能直接转成一个 issue
- 每条 finding 都至少附 2 个以上代码位置，证明不是偶发点
- 对“重复初始化”类问题，必须给出至少一个跨命令族的证据

---

**最后更新**: 2026-04-03
