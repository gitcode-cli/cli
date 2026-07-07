## 作者自检

- 作者主体标识: AI 实现子代理 (glm-5.2 via opencode, 分支 bugfix/issue-395)
- 根因或实现理由: RemoteURL 直接传 name 给 git remote get-url，无 -- 分隔无 ValidateRef。恶意 remote 名（--upload-pack=）会被 git 解释为 option。与同文件 SafeFetch/SafeCheckout 不一致。
- 主要修改: 2 文件 — git.go RemoteURL 加 ValidateRef + --；git_test.go 新增 TestRemoteURLRejectsOptionInjection
  - git/git.go: RemoteURL 加 ValidateRef(name) 校验 + `--` 分隔符
  - git/git_test.go: 新增 TestRemoteURLRejectsOptionInjection（option 注入/dash 前缀/空/shell metacharacter 4 用例）
- 影响范围: git.RemoteURL（查询 remote URL）；不改变用户可见命令行为（正常 remote 名通过，恶意被拒）
- 单元测试: ✅ TestRemoteURLRejectsOptionInjection 全 PASS，-race 通过
- 构建: ✅ go build ./... 全包通过；go vet ./git/... 无问题
- 实际命令验证: ⏩ 豁免 — 内部 git 封装，UT 充分覆盖（4 恶意 name 用例 + ValidateRef 已有测试）
- 安全审查: ✅ 防 git option 注入，不涉及 token/凭证；与同文件 SafeFetch/SafeCheckout 风格一致
- 文档同步: ✅ 内部行为，不改命令行为，docs/COMMANDS.md 无需改；设计文档已补
- 风险: medium（classify-change-risk → risk=medium，runtime 路径）
- 未覆盖项: CI 验证待 PR 推送触发
- 自检结论: 可进入 ready-for-review（medium 风险，需独立 AI 评审）

---

**关联 Issue**: #395
**Closes #395**
