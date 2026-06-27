#!/bin/bash
# Full-flow delivery loop — independent context via claude -p
# PID-locked; each tick is a fresh Claude session.
# Token tracking via stream-json output.

set -uo pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
LOCKFILE="$PROJECT_ROOT/.loop/run/full-flow.pid"
LOGDIR="$PROJECT_ROOT/.loop/history"
ISSUE_FILE="$PROJECT_ROOT/.loop/run/last-issue.txt"
PROMPT_FILE="$PROJECT_ROOT/.loop/prompts/full-flow-subprocess.md"
DELIVERIES_DIR="$PROJECT_ROOT/.loop/deliveries"

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
cd "$PROJECT_ROOT"
unset HTTP_PROXY HTTPS_PROXY http_proxy https_proxy

# Capture stream-json for token parsing
# timeout: kill hung sessions after 20min (normal delivery: 10-13min)
# nohup: survive parent shell exit
set +e
cat "$PROMPT_FILE" | nohup timeout 1200 claude -p \
  --verbose \
  --output-format stream-json \
  --permission-mode bypassPermissions \
  --add-dir "$PROJECT_ROOT" \
  > "$JSONL_FILE" 2>&1
rc=$?

# All post-processing is best-effort.
# Pipe closes when claude -p exits → python3 readlines() blocks until EOF.
python3 "$PROJECT_ROOT/.loop/scripts/process_tokens.py" "$JSONL_FILE" "$LOGFILE" "$TOKEN_FILE" "$DELIVERIES_DIR" >> "$LOGFILE" 2>&1 || true
