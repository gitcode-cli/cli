#385 refactor: full-flow-run.sh 存在冗余 README_FILE 赋值 | 2026-06-27
PR: N/A (pre-existing fix) | risk: low | code-change
Gates: G5(✅) G8(✅) | CI: skipped (fix already on main)
Commit: dfb8b20 (refactor(loop): extract token processing to standalone Python script)

## 处置说明

Issue 描述的冗余 `README_FILE` 赋值（`full-flow-run.sh` 第 179 行和第 226 行完全相同）已在 commit `dfb8b20` 中修复。

该重构将约 200 行内联 Python 代码从 `full-flow-run.sh` 提取到独立的 `process_tokens.py`，重复的 `README_FILE` 赋值被一并消除：
- `full-flow-run.sh`: 60 行，不再包含 `README_FILE`
- `process_tokens.py`: 仅在 114 行定义一次 `readme_file`
