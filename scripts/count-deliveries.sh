#!/usr/bin/env bash
# Update statistics in .loop/deliveries/README.md from its summary table.
set -euo pipefail

README=".loop/deliveries/README.md"

# Parse only the summary table (stop at ## 统计)
TABLE=$(sed '/^## 统计/q' "$README")

total=$(echo "$TABLE" | grep -cE '^\| \[#' || true)
merged=$(echo "$TABLE" | grep -cE '^\| \[#[0-9]+\].*\| merged \|' || true)
closed=$(echo "$TABLE" | grep -cE '^\| \[#[0-9]+\].*\| closed \|' || true)
code=$(echo "$TABLE" | grep -cE '^\| \[#[0-9]+\].*\| bug \| merged \|' || true)
docs=$(echo "$TABLE" | grep -cE '^\| \[#[0-9]+\].*\| docs \|' || true)
high=$(echo "$TABLE" | grep -cE 'high \|' || true)

scores=$(echo "$TABLE" | grep -oE '[0-9]+/8' | sed 's|/8||' | tr '\n' ' ')
if [ -n "$scores" ]; then
  sum=0; count=0
  for s in $scores; do sum=$((sum + s)); count=$((count + 1)); done
  avg=$(echo "scale=1; $sum / $count" | bc)
else
  avg="N/A"
fi

# Replace everything after "## 统计" with new stats
tmp=$(mktemp)
awk -v stats="\
| 维度 | 数据 |
|------|------|
| 总 issue | $total |
| 已合并 | $merged |
| 已关闭 | $closed |
| 含代码改动 | $code |
| docs-only | $docs |
| risk/high | $high |
| 平均门禁 | ${avg}/8 |" '
/^## 统计/ { print; print ""; print stats; exit }
{ print }
' "$README" > "$tmp" && mv "$tmp" "$README"

echo "Stats written to $README"
