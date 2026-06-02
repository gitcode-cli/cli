# CI 工作流规范

本文件定义 gitcode-cli 的远端 CI 工作流规范，包括 CI 定位、触发方式、Job 映射、AI 编排流程和证据纳入规则。

## 职责

- 定义远端 CI 在项目质量门禁体系中的定位
- 定义 AI 如何通过 `gh` CLI 触发和监控 CI
- 定义 CI Job 与 `spec/foundations/code-quality-gates.md` 的映射关系
- 定义 CI 结果如何纳入 PR 自检证据

## 适用场景

- AI 协作者在本地验证完成后、进入自检前，触发远端 CI 验证
- 跨平台（Linux/macOS/Windows）构建与测试验证
- PR 自检中引用 CI 结果作为自动化证据

## 必须

- CI 由 AI 通过 `gh` CLI 触发，不作为人工手动步骤
- CI 结果必须纳入 PR 自检证据
- CI 失败时不得进入 `status/self-checked`
- CI Job 定义变更时同步本文件

## 禁止

- 把 CI 通过当作跳过本地验证的理由
- 把 CI 通过当作独立执行主体评审的替代
- 在 CI 运行中或失败时宣称"开发完成"
- 修改 CI 定义绕过本文件规定的门禁映射

## 同步要求

- CI Job 变化时同步本文件和 `spec/foundations/code-quality-gates.md`
- CI 触发流程变化时同步 `spec/workflows/ai-local-development-workflow.md`

## 不负责什么

- 本地构建与单元测试（由本地门禁负责）
- 独立执行主体语义审查（由评审流程负责）
- 真实命令验证（由 `infra-test/*` 负责）

---

## 1. CI 定位

### 1.1 在质量门禁体系中的层级

远端 CI 位于本地门禁和 PR 门禁之间，作为**自动化证据补充层**：

```
本地开发门禁（必须，不可跳过）
  → 远端 CI 验证（AI 触发，自动化证据采集）
  → PR 门禁（自检证据 + CI 结果）
  → 合并门禁（独立评审 + 人工确认）
```

CI 不替代任何现有门禁层，只自动化采集原本需要人工手动执行的跨平台验证证据。

### 1.2 运行平台

CI 运行在 **GitHub Actions**（GitHub 镜像仓 `github.com/gitcode-cli/cli`）。

GitCode 主仓（`gitcode.com/gitcode-cli/cli`）当前不作为 CI 运行平台。

### 1.3 工具链

| 操作 | 工具 | 说明 |
|------|------|------|
| 触发 CI | `gh workflow run` | 触发 GitHub Actions workflow |
| 监控运行 | `gh run watch` | 实时等待 CI 完成 |
| 查看日志 | `gh run view --log` | 失败时获取详细日志 |
| 列出运行 | `gh run list` | 查看历史 CI 运行记录 |

`gh` 是 GitHub CLI，操作对象是 GitHub 镜像仓。GitCode 平台操作仍使用 `gc`。

---

## 2. CI Job 定义

### 2.1 Job 概览

当前 `.github/workflows/ci.yml` 包含以下 Job：

| Job | 运行环境 | 内容 | 对应质量门禁 |
|-----|---------|------|-------------|
| `lint` | ubuntu-latest | golangci-lint | 代码规范检查（`coding-standards.md`） |
| `test` | ubuntu / macos / windows | `go test -v -race -coverprofile` | 单元测试 + 竞态检测 + 覆盖率（`testing-guide.md`） |
| `build` | ubuntu / macos / windows | `go build` + `gc version` | 跨平台构建验证（`build-and-package.md`） |
| `docker` | ubuntu-latest | Docker 构建 + shell 补全生成 | 容器化构建验证 |

### 2.2 Job 依赖关系

```
lint ──┐
       ├──→ build (needs: test)
test ──┘
           docker (needs: test)
```

- `lint` 和 `test` 并行启动
- `build` 和 `docker` 等待 `test` 通过后执行
- 任何 Job 失败即整体 CI 失败

### 2.3 与质量门禁的映射

| 质量门禁要求 | CI 覆盖 |
|-------------|---------|
| `go test ./...` | `test` Job（3 OS，`-race`） |
| `go build` | `build` Job（3 OS） |
| 格式/规范检查 | `lint` Job（golangci-lint） |
| 跨平台兼容 | `test` + `build` 覆盖 ubuntu/macos/windows |

CI **不覆盖**的质量门禁（仍需本地或人工执行）：

- 真实命令验证（`infra-test/*`）
- 安全审查（凭证扫描、敏感信息检查）
- 文档同步检查
- 工作区卫生检查
- 独立执行主体语义审查

---

## 3. AI 触发与监控流程

### 3.1 触发时机

AI 在以下条件下触发 CI：

1. 本地 `go test ./...` 和 `go build` 已通过
2. 代码已提交并推送到 GitHub 镜像仓
3. 改动涉及代码路径（docs-only 可跳过 CI）

### 3.2 标准触发命令

```bash
# 触发 CI（在项目根目录执行）
gh workflow run ci.yml

# 获取最新运行的 run ID
gh run list --workflow=ci.yml --limit 1 --json databaseId --jq '.[0].databaseId'

# 实时等待 CI 完成
gh run watch <run-id>

# CI 完成后查看结论
gh run view <run-id> --json conclusion --jq '.conclusion'
```

### 3.3 推荐的一体化触发与等待脚本

```bash
# 触发 CI 并等待完成（一步到位）
gh workflow run ci.yml
sleep 5  # 等待 GitHub 创建运行记录
RUN_ID=$(gh run list --workflow=ci.yml --limit 1 --json databaseId --jq '.[0].databaseId')
gh run watch $RUN_ID
CONCLUSION=$(gh run view $RUN_ID --json conclusion --jq '.conclusion')
echo "CI conclusion: $CONCLUSION"
```

### 3.4 失败处理

CI 失败时，AI 必须：

1. 获取失败 Job 的详细日志：`gh run view <run-id> --log --job=<job-id>`
2. 分析根因（代码问题 vs 环境问题 vs 偶发问题）
3. 修复后重新推送并重新触发 CI
4. 在 PR 自检中记录 CI 失败与修复过程

如果 CI 失败原因是环境/平台偶发问题（非代码问题），可在自检中明确说明，仍可继续推进流程。

---

## 4. CI 证据纳入自检

### 4.1 最小 CI 证据

PR 作者自检中至少包含：

- CI run ID 或 run URL
- CI 结论（success / failure）
- 各 Job 状态摘要
- 如 CI 失败，失败原因和修复记录

### 4.2 自检模板中的 CI 条目

```markdown
## CI 验证

- Run URL: https://github.com/gitcode-cli/cli/actions/runs/<run-id>
- 结论: success
- Job 状态:
  - lint: ✅
  - test (ubuntu): ✅
  - test (macos): ✅
  - test (windows): ✅
  - build (ubuntu): ✅
  - build (macos): ✅
  - build (windows): ✅
  - docker: ✅
```

### 4.3 CI 未执行的处理

如果因以下原因未执行 CI，必须在自检中明确说明：

- docs-only 改动（写明"不涉及代码路径，已跳过 CI"）
- GitHub 镜像仓不可达（写明具体错误）
- 其他合理原因（需明确记录）

---

## 5. 约束与边界

### 5.1 CI 通过 ≠ 可以合并

CI 通过只表示自动化检查无问题。以下事项仍需独立完成：

- 真实命令验证（`infra-test/*`）
- 安全审查
- 文档同步
- 独立执行主体评审
- 高风险改动的人工最终确认

### 5.2 CI 不定义新门禁

CI 是现有质量门禁的自动化实现，不得引入高于 `spec/foundations/code-quality-gates.md` 的额外要求。

### 5.3 CI 配置变更

修改 `.github/workflows/ci.yml` 的行为等同于修改构建/测试门禁，必须：

- 同步更新本文件中的 Job 描述和映射表
- 在 PR 中说明变更理由
- 变更后至少成功运行一次 CI 作为自验证

---

## 6. Release CI

`.github/workflows/release.yml` 用于版本发布构建，不属于日常开发 CI。

触发方式：AI 在发布流程中通过 `gh workflow run release.yml -f version=vX.Y.Z` 触发。

Release CI 规范详见 `spec/delivery/release-process.md`。

---

## 下一步去看哪里

- CI 不通过的修复流程：看 [测试流程](../workflows/test-workflow.md)
- CI 结果如何影响合并：看 [代码质量门禁规范](../foundations/code-quality-gates.md)
- Release CI 详情：看 [发布流程规范](./release-process.md)
- AI 如何编排 CI：看 [AI 本地开发流程](../workflows/ai-local-development-workflow.md)

---

**最后更新**: 2026-06-02
