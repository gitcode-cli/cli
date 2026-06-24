#!/usr/bin/env bash
# Generate statistics from .loop/deliveries/README.md summary table.
set -euo pipefail

README=".loop/deliveries/README.md"

total=$(grep -cE '^\| \[#' "$README" || true)
merged=$(grep -cE '^\| \[#[0-9]+\].*\| merged \|' "$README" || true)
closed=$(grep -cE '^\| \[#[0-9]+\].*\| closed \|' "$README" || true)
code=$(grep -cE '^\| \[#[0-9]+\].*\| bug \| merged \|' "$README" || true)
docs=$(grep -cE '^\| \[#[0-9]+\].*\| docs \|' "$README" || true)
high=$(grep -cE 'high.+' "$README" || true)

# Average gate score: extract last column, filter numbers with /8
scores=$(grep -oE '[0-9]+/8' "$README" | sed 's|/8||' | tr '\n' ' ')
if [ -n "$scores" ]; then
  sum=0; count=0
  for s in $scores; do sum=$((sum + s)); count=$((count + 1)); done
  avg=$(echo "scale=1; $sum / $count" | bc)
else
  avg="N/A"
fi

cat <<EOF
| 维度 | 数据 |
|------|------|
| 总 issue | $total |
| 已合并 | $merged |
| 已关闭 | $closed |
| 含代码改动 | $code |
| docs-only | $docs |
| risk/high | $high |
| 平均门禁 | ${avg}/8 |
EOF
