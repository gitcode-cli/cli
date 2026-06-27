#!/usr/bin/env bash
# Delivery integrity check — detect mismatches between git merges,
# delivery files, and README rows. Also flag stubs and residue.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" 2>/dev/null && pwd)"
# Fallback: when invoked via stdin, $0 is "bash" — use git root instead
if [ ! -d "$SCRIPT_DIR/../../.loop" ]; then
    SCRIPT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)/.loop/scripts"
fi
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DELIVERIES="$PROJECT_ROOT/.loop/deliveries"
README="$DELIVERIES/README.md"

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

errors=0
warnings=0

# Helper: increment without tripping set -e (post-increment returns 0 → exit)
incr() { local -n _v=$1; _v=$((_v + $2)); }

section() { echo -e "\n${GREEN}── $1${NC}"; }

# ── 1. Residue ──────────────────────────────────────────────
section "1. Residue files"
found=0
for f in "$DELIVERIES"/*.bak; do
    [ -f "$f" ] || continue
    echo -e "  ${RED}✗${NC} residue: $(basename "$f") ($(wc -c < "$f") bytes)"
    incr errors 1
    found=1
done
[ $found -eq 0 ] && echo "  none"

# ── 2. Stub delivery files (<500 bytes) ─────────────────────
section "2. Stub delivery files (<500 bytes, missing gate evidence)"
stubs=()
while IFS= read -r f; do
    stubs+=("$(basename "$f" .md)")
done < <(find "$DELIVERIES" -maxdepth 1 -name 'issue-*.md' -size -500c | sort -V)
if [ ${#stubs[@]} -gt 0 ]; then
    echo -e "  ${YELLOW}⚠${NC} ${#stubs[@]} stubs: ${stubs[*]}"
    incr warnings "${#stubs[@]}"
else
    echo "  none"
fi

# ── 3. Delivery files ↔ README cross-check ──────────────────
section "3. Delivery files ↔ README cross-check"

mapfile -t file_issues < <(ls "$DELIVERIES"/issue-*.md 2>/dev/null | sed 's/.*issue-//;s/\.md//' | sort -n)
mapfile -t readme_issues < <(/usr/bin/grep -oP 'issue-\d+' "$README" 2>/dev/null | sed 's/issue-//' | sort -n | uniq)

# Files without README row
for n in "${file_issues[@]}"; do
    if ! printf '%s\n' "${readme_issues[@]}" | /usr/bin/grep -qx "$n"; then
        echo -e "  ${RED}✗${NC} issue-$n.md exists but no README row"
        incr errors 1
    fi
done

# README rows without file
for n in "${readme_issues[@]}"; do
    if ! printf '%s\n' "${file_issues[@]}" | /usr/bin/grep -qx "$n"; then
        echo -e "  ${RED}✗${NC} README row for #$n but issue-$n.md missing"
        incr errors 1
    fi
done

[ $errors -eq 0 ] && echo "  all matched"

# ── 4. Git merges missing delivery files ────────────────────
section "4. Merged issues without delivery record (last 90 days)"

# Temp file to avoid process-substitution pipefail triggering set -e
tmp_merges=$(mktemp)
git -C "$PROJECT_ROOT" log origin/main --merges --format='%h %s' --since='90 days ago' 2>/dev/null | /usr/bin/grep -iE 'issue|#' > "$tmp_merges" || true

while IFS= read -r line; do
    [ -z "$line" ] && continue
    n=$(echo "$line" | /usr/bin/grep -oP 'issue[-\s#]?\K\d+' | head -1 || true)
    [ -z "$n" ] && continue
    if [ ! -f "$DELIVERIES/issue-$n.md" ]; then
        echo -e "  ${YELLOW}⚠${NC} merged #$n has no delivery file ($(echo "$line" | cut -c1-80))"
        incr warnings 1
    fi
done < "$tmp_merges"
rm -f "$tmp_merges"

[ $warnings -eq 0 ] && echo "  none"

# ── 5. README rows missing token/cost data ──────────────────
section "5. README rows missing token/cost data"

while IFS= read -r line; do
    n=$(echo "$line" | /usr/bin/grep -oP 'issue-\d+' | head -1)
    tok=$(echo "$line" | cut -d'|' -f10 | xargs)
    if [ "$tok" = "—" ]; then
        echo -e "  ${YELLOW}⚠${NC} $n: token column empty"
        incr warnings 1
    fi
done < <(/usr/bin/grep '^| \[#' "$README" || true)

[ $warnings -eq 0 ] && echo "  none"

# ── Summary ─────────────────────────────────────────────────
echo -e "\n${GREEN}── Summary${NC}"
echo "  Delivery files: ${#file_issues[@]}"
echo "  README rows:    ${#readme_issues[@]}"
echo "  Stubs:          ${#stubs[@]}"
if [ $errors -gt 0 ] || [ $warnings -gt 0 ]; then
    echo -e "  ${RED}Errors:   $errors${NC}"
    echo -e "  ${YELLOW}Warnings: $warnings${NC}"
    exit 1
else
    echo -e "  ${GREEN}All checks passed${NC}"
fi
