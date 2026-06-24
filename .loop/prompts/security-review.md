# /goal: 安全审查

## Prompt

```
/goal until security review passes:
  - git diff origin/main 中无硬编码 token/password/secret
  - 文档、测试中未误写真实凭证
  - 涉及认证/配置/权限路径的改动已对照 spec/foundations/security.md 检查
  - 审查结论已写入自检记录
```

## 评估器检查点

- 对话中是否出现 "无硬编码凭证" 或等效文本
- 安全审查结论是否已写入自检

## 注意事项

- docs-only 改动仍需检查文档中是否误写凭证
- 高风险改动（auth/release/delete）需额外关注
