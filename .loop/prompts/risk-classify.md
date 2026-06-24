# /goal: 风险分级

## Prompt

```
/goal until python3 scripts/classify-change-risk.py --base origin/main
outputs a risk level AND the result is recorded in self-check
完成后更新 .loop/deliveries/issue-<N>.md 记录风险等级
```

## 注意事项

- 风险脚本可能因累积 diff 误判为 high
- 如果判断不合理，在自检中说明原因

## .loop/ 更新

```markdown
| risk_classify | <ts> | risk/low (or medium/high) |
```
