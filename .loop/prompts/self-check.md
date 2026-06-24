# /goal: 作者自检模板

## Prompt

```
/goal until PR self-check is complete for issue #<ISSUE_NUMBER>:
  - 根因或实现理由、修改范围、测试结果、CI 证据
  - 安全审查结果、文档同步说明、风险分级、未覆盖项
  AND gc issue comment <ISSUE_NUMBER> with self-check record -R gitcode-cli/cli
  AND gc pr comment <PR_NUMBER> with gate compliance summary -R gitcode-cli/cli
  AND gc issue label <ISSUE_NUMBER> --add status/self-checked -R gitcode-cli/cli
完成后更新 .loop/deliveries/issue-<ISSUE_NUMBER>.md 和 .loop/deliveries/README.md 汇总表
```

## .loop/ 更新

```markdown
| status/in-progress | status/self-checked | <ts> | self-check comment <url> |
```
