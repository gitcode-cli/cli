# API 命令 (api)

> 本文档是 Claude 参考层，不是命令行为真相源。
> API 命令行为以 `docs/COMMANDS.md` 和 `spec/` 为准。

## api - 调用 GitCode API

```bash
# 读取仓库 API 原始响应
gc api repos/owner/repo

# 读取 PR 文件列表
gc api repos/owner/repo/pulls/1/files

# 带查询参数的 API，包含 & 时请整体加引号
gc api 'repos/owner/repo/commits?path=README.md&sha=main'

# 指定 HTTP 方法和请求体文件
gc api repos/owner/repo/pulls/1 --method PATCH --input body.json
```

## 使用约束

- endpoint 可写成 `repos/owner/repo` 或 `/api/v5/repos/owner/repo`
- 输出为远端原始响应 body，不额外包装 JSON
- 写操作前必须确认目标资源、权限和请求体内容

---

**最后更新**: 2026-05-27
