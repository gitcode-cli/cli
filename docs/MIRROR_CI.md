# Mirror CI

GitCode 主仓不运行 CI。gitcode-cli 使用 GitHub mirror Actions 作为 CI 执行源。

## 规则

- GitCode issue / PR 是协作事实源。
- GitCode `origin/main` 是主干完成事实源。
- GitHub Actions 是 CI 执行事实源。
- commit SHA 是 GitCode PR 与 GitHub Actions run 的绑定键。
- GitHub mirror 不能替代 GitCode 主仓状态。

## 标准流程

```text
GitCode PR
-> 获取 head commit SHA
-> 等待 GitHub mirror 同步该 SHA
-> 查询 GitHub Actions
-> 生成 CI evidence
-> 回写 GitCode PR comment
```

## 状态

- `ci_waiting`：mirror 尚未同步或 run 尚未完成
- `ci_passed`：指定 SHA 的 required jobs 通过
- `ci_failed`：指定 SHA 的 required jobs 失败
- `ci_unavailable`：CI 事实源不可访问

不得在 `ci_waiting` 或 `ci_unavailable` 时宣称 CI 已通过。
