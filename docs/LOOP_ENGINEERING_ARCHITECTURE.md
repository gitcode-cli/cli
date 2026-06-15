# Loop Engineering Architecture

## 分层

```text
Loop Engineering
|
|-- spec / policy
|   定义规则、状态机、人工确认点
|
|-- skills
|   定义 AI 该如何执行某类 loop
|
|-- loop-kits
|   定义 hooks、schema、templates、adapters
|
|-- GitCode
|   保存 issue、PR、评审、状态和证据
|
`-- GitHub mirror Actions
    提供 commit SHA 绑定的 CI 执行事实
```

## 调用关系

```text
User / Automation
-> AI skill
-> read spec / policy
-> use loop-kits templates / hooks / adapters
-> operate GitCode issue / PR with gitcode
-> query GitHub mirror CI
-> write evidence back to GitCode
```

## 长期记忆

- 状态在 GitCode issue / PR
- 证据在 GitCode comment 和 GitHub Actions URL
- 规则在 `spec/`
- AI 方法在 `gitcode-cli/skills`
- 标准件在 `gitcode-cli/loop-kits`
- 本地 `.loop/runtime` 只是缓存
