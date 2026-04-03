# Testing/Regression Agent 任务书

## 目标

审视测试覆盖、回归矩阵和高风险空洞，只产出可直接提交 issue 的测试与回归缺口。

## 本 agent 只看什么

- 高风险命令或共享层是否缺少单元测试
- 文档或流程里列为关键门禁的路径，是否没有回归覆盖
- 是否存在“只能靠人工发现、无法被测试兜住”的缺口
- `infra-test/*` 真实验证矩阵是否覆盖关键能力

## 必看目录

- `pkg/`
- `api/`
- `scripts/`
- `docs/`
- `spec/foundations/`
- `spec/workflows/`

## 必看文件

- 全部 `*_test.go`
- `scripts/regression-core.sh`
- `docs/REGRESSION.md`
- `spec/foundations/testing-guide.md`
- `spec/workflows/test-workflow.md`

## 建议重点包

- `pkg/cmd/root`
- `pkg/cmd/release/edit`
- `pkg/cmd/release/upload`
- `pkg/cmd/repo/sync`
- `pkg/cmd/pr/merge`
- `pkg/cmd/commit/`
- `pkg/browser/`
- `internal/config/`

## 必答问题

1. 哪些高风险路径几乎没有测试？
2. 哪些近期暴露出的真实问题，文档已经记录，但仍未被测试覆盖？
3. 哪些命令族的行为依赖真实 API，但 regression 核验矩阵没有覆盖？
4. 是否存在共享层变更会影响多条命令，却没有集中测试的情况？
5. 当前测试缺口里，哪些已经足以单独提交 issue？

## 明确算 issue 的问题类型

- 高风险路径完全无测试或无回归
- 文档门禁要求与实际测试矩阵脱节
- 已知历史缺陷没有回归测试守护
- 共享层或关键命令族缺少最小可执行验证

## 不要提交 issue 的内容

- 一般性的“测试再多一点更好”
- 为了覆盖率数字而补的低价值建议
- 没有具体风险面的测试抱怨

## 输出格式

```markdown
## Finding N

- 标题建议:
- 严重度: high / medium / low
- 类型: testing
- 现象:
- 根因:
- 影响:
- 涉及目录或文件:
- 复现命令或核验方式:
- 是否疑似已有远端 issue:
```

## 通过标准

- 每条 finding 都必须说明“为什么这是测试空洞而不是实现缺陷”
- 如果建议补 regression，必须指明应该落到哪类命令验证
- 如果建议补单测，必须说明目标包或共享层

---

**最后更新**: 2026-04-03
