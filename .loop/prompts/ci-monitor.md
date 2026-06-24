# /loop: CI 监控 + 修复

## Prompt

```
/loop 5m 监控 gh run list --workflow=ci.yml --branch <BRANCH>:
  - CI 全绿 → 记录 run ID 到自检证据，停止
  - CI 失败 → 获取日志，诊断根因
    - 本次改动导致的 → 修复，commit，push，回到监控
    - 环境/平台偶发 → 在自检中记录，继续推进
    - 预存 bug → 记录 issue，继续推进
```

## 替换参数

- `<BRANCH>`: 当前开发分支名
- 间隔可按需调整（推荐 5m）

## 注意事项

- 使用 `gh` CLI 而非 `gc`
- CI 仅在 GitHub 镜像仓有效
- 代理问题：`unset HTTP_PROXY HTTPS_PROXY` 后再执行 `gh`
