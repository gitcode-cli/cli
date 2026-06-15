# GitHub Mirror CI 契约

## 背景

GitCode 主仓不运行 CI。gitcode-cli 的 CI 执行源是 GitHub mirror Actions。

## 事实源划分

- GitCode issue / PR：协作事实
- GitCode `origin/main`：主干完成事实
- GitHub Actions：CI 执行事实
- commit SHA：GitCode 与 GitHub CI 的绑定键

## 流程

```text
GitCode PR
-> 获取 PR head commit SHA
-> 等待 GitHub mirror 同步该 SHA
-> 查询 GitHub Actions run
-> 生成 CI evidence
-> 回写 GitCode PR comment
```

## 规则

- CI evidence 必须包含 commit SHA、workflow/run URL、状态和获取时间
- mirror 未同步时状态是 `ci_waiting`
- CI 不可访问时状态是 `ci_unavailable`，不得编造通过或失败
- CI 失败时进入 fix loop，不得直接 archive
- GitHub mirror 不能替代 GitCode merged PR 或 `origin/main` 判定
