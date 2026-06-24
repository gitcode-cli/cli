# /goal: 安全审查

## Prompt

```
/goal until security review passes:
  - git diff origin/main 中无硬编码 token/password/secret
  - 文档、测试中未误写真实凭证
  - 涉及认证/配置/权限路径的改动已对照 spec/foundations/security.md 检查
  - 审查结论已写入自检记录
完成后更新 .loop/deliveries/issue-<N>.md 和 .loop/deliveries/README.md 汇总表
```

## 注意事项

- docs-only 改动仍需检查文档中是否误写凭证

## .loop/ 更新

```markdown
| security_review | <ts> | clean / findings |
```
