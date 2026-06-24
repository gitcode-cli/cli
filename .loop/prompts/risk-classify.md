# /goal: 风险分级

## Prompt

```
/goal until python3 scripts/classify-change-risk.py --base origin/main
outputs a risk level AND the result is recorded in self-check
```

## 评估器检查点

- 对话中出现 `risk=low`、`risk=medium` 或 `risk=high`
- 结果已写入自检记录

## 注意事项

- 风险脚本可能因累积 diff 误判为 high
- 如果判断不合理，在自检中说明原因
- risk/low → 自动推进，risk/high → 暂停确认
