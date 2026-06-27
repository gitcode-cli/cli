#!/bin/bash
# Full-flow delivery loop — independent context via claude -p
# PID-locked; each tick is a fresh Claude session.
# Token tracking via stream-json output.

set -uo pipefail

LOCKFILE="/home/wpf/claude-code/vibe-coding/cli/.loop/run/full-flow.pid"
LOGDIR="/home/wpf/claude-code/vibe-coding/cli/.loop/history"
ISSUE_FILE="/home/wpf/claude-code/vibe-coding/cli/.loop/run/last-issue.txt"
PROMPT_FILE="/home/wpf/claude-code/vibe-coding/cli/.loop/prompts/full-flow-subprocess.md"
DELIVERIES_DIR="/home/wpf/claude-code/vibe-coding/cli/.loop/deliveries"

mkdir -p "$(dirname "$LOCKFILE")" "$LOGDIR" "$DELIVERIES_DIR"

# --- PID lock ---
if [ -f "$LOCKFILE" ]; then
    old_pid=$(cat "$LOCKFILE")
    if kill -0 "$old_pid" 2>/dev/null; then
        echo "[$(date -Iseconds)] SKIP: pid=$old_pid still running" | tee -a "$LOGDIR/run.log"
        exit 0
    fi
    rm -f "$LOCKFILE"
fi
echo $$ > "$LOCKFILE"
cleanup() { rm -f "$LOCKFILE"; }
trap cleanup EXIT INT TERM

# --- Logging ---
TS=$(date +%Y-%m-%d-%H%M%S)
LOGFILE="$LOGDIR/$TS-full-flow.log"
JSONL_FILE="$LOGDIR/$TS-full-flow.jsonl"
TOKEN_FILE="$LOGDIR/$TS-full-flow.tokens.json"

echo "[$(date -Iseconds)] START pid=$$" | tee "$LOGFILE"

# --- Run ---
cd /home/wpf/claude-code/vibe-coding/cli
unset HTTP_PROXY HTTPS_PROXY http_proxy https_proxy

# Capture stream-json for token parsing
# timeout: kill hung sessions after 20min (normal delivery: 10-13min)
# nohup: survive parent shell exit
set +e
cat "$PROMPT_FILE" | nohup timeout 1200 claude -p \
  --verbose \
  --output-format stream-json \
  --permission-mode bypassPermissions \
  --add-dir /home/wpf/claude-code/vibe-coding/cli \
  > "$JSONL_FILE" 2>&1
rc=$?

# Brief pause: ensure JSONL fully flushed before parsing
sleep 2

# All post-processing is best-effort
set +e
python3 -c "
import json, sys

with open('$JSONL_FILE') as f:
    lines = f.readlines()

for line in lines:
    line = line.strip()
    if not line or not line.startswith('{'):
        continue
    try:
        d = json.loads(line)
    except json.JSONDecodeError:
        continue
    if d.get('type') == 'assistant':
        msg = d.get('message', {})
        for c in msg.get('content', []):
            if c.get('type') == 'text':
                print(c['text'])
" >> "$LOGFILE" 2>/dev/null || true

# --- Post-process: extract token data ---
RESULT=$(python3 -c "
import json
with open('$JSONL_FILE') as f:
    for line in f:
        line = line.strip()
        if not line or not line.startswith('{'):
            continue
        try:
            d = json.loads(line)
        except json.JSONDecodeError:
            continue
        if d.get('type') == 'result':
            last_result = line
print(last_result)
" 2>/dev/null)

if [ -n "$RESULT" ]; then
    echo "$RESULT" | python3 -c "
import json, sys
d = json.load(sys.stdin)
u = d.get('usage', {})
mu = d.get('modelUsage', {})
total_in = u.get('input_tokens', 0)
total_out = u.get('output_tokens', 0)
total = total_in + total_out
cache_read = u.get("cache_read_input_tokens", 0)
cache_create = u.get("cache_creation_input_tokens", 0)
cost = d.get('total_cost_usd', 0)
dur_ms = d.get('duration_ms', 0)
turns = d.get('num_turns', 0)

# Find per-model breakdown
models = []
for m, v in mu.items():
    models.append(f\"{m}: in={v.get('inputTokens',0)} out={v.get('outputTokens',0)} cost=\${v.get('costUSD',0):.4f}\")

print(f'---')
print(f'TOKENS:  in={total_in}  out={total_out}  total={total}  cache_read={cache_read}  cache_create={cache_create}  cost=\${cost:.4f}  duration={dur_ms}ms  turns={turns}')
for m in models:
    print(f'  {m}')

# Write token JSON for delivery integration
	# Calculate DeepSeek cost (RMB per million tokens)
	deepseek_cost_rmb = 0.0
	for m, v in mu.items():
	    model_in = v.get("inputTokens", 0)
	    model_cache = v.get("cacheReadInputTokens", 0)
	    model_out = v.get("outputTokens", 0)
	    deepseek_cost_rmb += (model_in / 1_000_000) * 3.0
	    deepseek_cost_rmb += (model_cache / 1_000_000) * 0.025
	    deepseek_cost_rmb += (model_out / 1_000_000) * 6.0
	deepseek_cost_usd = round(deepseek_cost_rmb / 7.2, 4)
	deepseek_cost_rmb = round(deepseek_cost_rmb, 4)

	with open('$TOKEN_FILE', 'w') as f:
	    json.dump({
	        "input_tokens": total_in,
	        "cache_miss_tokens": total_in,
	        "cache_read_input_tokens": cache_read,
	        "output_tokens": total_out,
	        "total_tokens": total,
	        "cost_rmb": deepseek_cost_rmb,
	        "cost_usd": deepseek_cost_usd,
	        "duration_ms": dur_ms,
	        "num_turns": turns,
	        "pricing": "DeepSeek: ¥3/M cache-miss, ¥0.025/M cache-hit, ¥6/M output",
	    }, f, indent=2)
" >> "$LOGFILE"
else
    echo "---" >> "$LOGFILE"
    echo "WARNING: no result object found in stream-json output" >> "$LOGFILE"
fi

echo "[$(date -Iseconds)] DONE rc=$rc" | tee -a "$LOGFILE"

# Auto-refresh stats
bash scripts/count-deliveries.sh >> "$LOGFILE" 2>&1 || true

# --- Post-process: inject token data into delivery files ---
# Search JSONL (raw stream) first, fall back to log
ISSUE_NUM=$(python3 -c "
import re
# Try JSONL first (more reliable)
try:
    with open('$JSONL_FILE') as f:
        content = f.read()
    m = re.search(r'ISSUE_NUM=(\d+)', content)
    if m: print(m.group(1)); exit()
except: pass
# Fall back to log
try:
    with open('$LOGFILE') as f:
        content = f.read()
    m = re.search(r'ISSUE_NUM=(\d+)', content)
    if m: print(m.group(1))
except: pass
" 2>/dev/null)
if [ -n "$ISSUE_NUM" ] && [ -f "$TOKEN_FILE" ]; then
    DELIVERY_FILE="$DELIVERIES_DIR/issue-$ISSUE_NUM.md"
    README_FILE="$DELIVERIES_DIR/README.md"

    # Build token summary block
    TOKEN_BLOCK=$(python3 -c "
import json
with open('$TOKEN_FILE') as f:
    t = json.load(f)
in_k = t['input_tokens']/1000
out_k = t['output_tokens']/1000
total_k = t['total_tokens']/1000
dur = f\"{t['duration_ms']/1000:.0f}s\"
cost_rmb = t.get("cost_rmb", 0)
cost_usd = t.get("cost_usd", 0)
print(f\"\"\"
## Token 消耗

| 指标 | 值 |
|------|-----|
| 输入 tokens | {t['input_tokens']:,} ({in_k:.0f}k) |
| 输出 tokens | {t['output_tokens']:,} ({out_k:.0f}k) |
| 缓存命中 | {cache_read:,} ({cache_read/1000:.0f}k) |
| 缓存写入 | {cache_create:,} ({cache_create/1000:.0f}k) |
| 总计 tokens | {t['total_tokens']:,} ({total_k:.0f}k) |
| 成本 (DeepSeek) | ¥{cost_rmb} (~\${cost_usd}) |
| 耗时 | {dur} |
| 轮次 | {t['num_turns']} |
\"\"\")
")

    # Append to delivery file
    if [ -f "$DELIVERY_FILE" ]; then
        echo "$TOKEN_BLOCK" >> "$DELIVERY_FILE"
        echo "[$(date -Iseconds)] token data appended to $DELIVERY_FILE" | tee -a "$LOGFILE"
    fi

    # Update README row: inject token count into a new "tokens" column
    TOKEN_STR=$(python3 -c "
import json
t = json.load(open('$TOKEN_FILE'))
tk = t['total_tokens']
if tk >= 1_000_000:
    print(f'{tk/1_000_000:.1f}M')
elif tk >= 1_000:
    print(f'{tk/1000:.0f}k')
else:
    print(str(tk))
")
    README_FILE="$DELIVERIES_DIR/README.md"
    if [ -f "$README_FILE" ]; then
        python3 -c "
import sys

with open('$README_FILE') as f:
    lines = f.readlines()

target = '| [#$ISSUE_NUM](issue-$ISSUE_NUM.md) |'
new_lines = []
for line in lines:
    if target in line:
        # Strip trailing newline, append token column, re-add newline
        line = line.rstrip('\n')
        if line.endswith('|'):
            line = line + ' $TOKEN_STR |'
        else:
            line = line + ' $TOKEN_STR |'
        line = line + '\n'
    new_lines.append(line)

with open('$README_FILE', 'w') as f:
    f.writelines(new_lines)
"
        echo "[$(date -Iseconds)] README row updated with token=$TOKEN_STR" | tee -a "$LOGFILE"
    fi
else
    echo "[$(date -Iseconds)] no issue num or token data, skipping delivery update" | tee -a "$LOGFILE"
fi
