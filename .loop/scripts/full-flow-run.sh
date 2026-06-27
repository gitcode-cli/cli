#!/bin/bash
# Full-flow delivery loop — independent context via claude -p
# PID-locked; each tick is a fresh Claude session.

set -uo pipefail

LOCKFILE="/home/wpf/claude-code/vibe-coding/cli/.loop/run/full-flow.pid"
LOGDIR="/home/wpf/claude-code/vibe-coding/cli/.loop/history"
PROMPT_FILE="/home/wpf/claude-code/vibe-coding/cli/.loop/prompts/full-flow-subprocess.md"

mkdir -p "$(dirname "$LOCKFILE")" "$LOGDIR"

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
echo "[$(date -Iseconds)] START pid=$$" | tee "$LOGFILE"

# --- Run ---
cd /home/wpf/claude-code/vibe-coding/cli
unset HTTP_PROXY HTTPS_PROXY http_proxy https_proxy

# claude -p: prompt via stdin pipe, may exit non-zero on tool failures
set +e
cat "$PROMPT_FILE" | claude -p \
  --permission-mode bypassPermissions \
  --add-dir /home/wpf/claude-code/vibe-coding/cli \
  >> "$LOGFILE" 2>&1
rc=$?
set -e

echo "[$(date -Iseconds)] DONE rc=$rc" | tee -a "$LOGFILE"
