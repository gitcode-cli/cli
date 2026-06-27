#!/usr/bin/env bash
# Auto-generate statistics for .loop/deliveries/README.md from its summary table.
set -euo pipefail

README=".loop/deliveries/README.md"
[ -f "$README" ] || { echo "ERROR: $README not found"; exit 1; }

python3 << 'PYEOF'
import re, os

readme = os.environ.get('README', '.loop/deliveries/README.md')
with open(readme) as f:
    content = f.read()

rows = re.findall(r'^\| \[\#.*\|$', content, re.MULTILINE)

# Counts
total = len(rows)
merged = sum(1 for r in rows if '| merged |' in r)
closed = sum(1 for r in rows if '| closed |' in r)
has_change = sum(1 for r in rows if re.search(r'\+[0-9]+/-[0-9]+', r))
docs = sum(1 for r in rows if '| docs |' in r)
high = sum(1 for r in rows if 'high |' in r)
from datetime import datetime
today_str = datetime.now().strftime('%Y-%m-%d')
today = sum(1 for r in rows if today_str in r)

# Gate average
scores = [int(m) for r in rows for m in re.findall(r'(\d+)/8', r)]
avg = f"{sum(scores)/len(scores):.1f}" if scores else "N/A"

# Token total — parse "4.7M(4.6M cache)", "260k", "—", etc.
token_raw = 0
for r in rows:
    cols = [c.strip() for c in r.split('|') if c.strip()]
    if len(cols) < 9:
        continue
    tok = cols[8]  # Tokens column (0-indexed: 8 of 11)
    if tok == '—' or not tok:
        continue
    # Extract the main number before any parenthetical: "4.7M(4.6M cache)" → "4.7M"
    main = tok.split('(')[0].strip()
    m = re.match(r'([\d.]+)\s*(k|M)', main)
    if m:
        val = float(m.group(1))
        unit = m.group(2)
        token_raw += int(val * 1_000_000) if unit == 'M' else int(val * 1_000)
if token_raw >= 1_000_000:
    token_total = f"{token_raw/1_000_000:.1f}M"
elif token_raw >= 1_000:
    token_total = f"{token_raw/1_000:.0f}k"
else:
    token_total = str(token_raw)

# Cost total
costs = [float(m) for r in rows for m in re.findall(r'¥([\d.]+)', r)]
cost_total = f"¥{sum(costs):.2f}" if costs else "¥0"

stats = f"""| 维度 | 数据 |
|------|------|
| 总 issue | {total} |
| 已合并 | {merged} |
| 已关闭 | {closed} |
| 含代码改动 | {has_change} |
| docs-only | {docs} |
| risk/high | {high} |
| 平均门禁 | {avg}/8 |
| Token 总消耗 | {token_total} |
| 总成本 | {cost_total} |
| 今日交付 | {today} |"""

idx = content.find('## 统计')
if idx >= 0:
    content = content[:idx] + '## 统计\n\n' + stats + '\n'

with open(readme, 'w') as f:
    f.write(content)
print(stats)
PYEOF
