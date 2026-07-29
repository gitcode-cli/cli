#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/gc-dev-setup-test.XXXXXX")"
trap 'rm -rf "$TMP_ROOT"' EXIT

STUB_BIN="$TMP_ROOT/stubs"
TOOLS_DIR="$TMP_ROOT/tools"
CALL_LOG="$TMP_ROOT/calls.log"
TOKEN_LEAK="$TMP_ROOT/token-leak"
CACHE_LOG="$TMP_ROOT/cache.log"
mkdir -p "$STUB_BIN"

cat >"$STUB_BIN/tool-stub" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail
name="$(basename "$0")"
if [[ -n "${GC_TOKEN:-}" || -n "${GITCODE_TOKEN:-}" ]]; then
    : >"${TOKEN_LEAK:?}"
    exit 90
fi
printf '%s %s\n' "$name" "$*" >>"${CALL_LOG:?}"

case "$name" in
    go)
        case "${1:-} ${2:-}" in
            "env GOVERSION") echo "go1.22.9" ;;
            "version -m")
                module="$(sed -n 's/^# module=//p' "${3:-}" | head -n 1)"
                version="$(sed -n 's/^# version=//p' "${3:-}" | head -n 1)"
                [[ -n "$module" && -n "$version" ]] && printf 'mod\t%s\t%s\n' "$module" "$version"
                ;;
            "mod download") ;;
            "list -mod=readonly")
                printf '%s|%s\n' "${GOCACHE:-}" "${GOTMPDIR:-}" >>"${CACHE_LOG:?}"
                ;;
            "build -mod=readonly"|"build -o")
                printf '%s|%s\n' "${GOCACHE:-}" "${GOTMPDIR:-}" >>"${CACHE_LOG:?}"
                out=""
                while (($#)); do
                    if [[ "$1" == "-o" ]]; then out="$2"; break; fi
                    shift
                done
                printf '#!/usr/bin/env bash\nexit 0\n' >"$out"
                chmod +x "$out"
                ;;
            "test -mod=readonly"|"test -race")
                printf '%s|%s\n' "${GOCACHE:-}" "${GOTMPDIR:-}" >>"${CACHE_LOG:?}"
                ;;
            "install "*)
                spec="${2:-}"
                package="${spec%@*}"
                version="${spec##*@}"
                binary="${package##*/}"
                [[ "$package" == github.com/zricethezav/gitleaks/* ]] && binary="gitleaks"
                module="$package"
                [[ "$package" == github.com/golangci/golangci-lint/v2/* ]] && module="github.com/golangci/golangci-lint/v2"
                [[ "$package" == github.com/zricethezav/gitleaks/* ]] && module="github.com/zricethezav/gitleaks/v8"
                mkdir -p "${GOBIN:?}"
                printf '#!/usr/bin/env bash\n# module=%s\n# version=%s\nexit 0\n' "$module" "$version" >"$GOBIN/$binary"
                chmod +x "$GOBIN/$binary"
                ;;
            *) ;;
        esac
        ;;
    python3)
        if [[ "${1:-} ${2:-}" == "-m venv" ]]; then
            venv="$3"
            mkdir -p "$venv/bin"
            cat >"$venv/bin/python" <<'PYTHON'
#!/usr/bin/env bash
set -euo pipefail
if [[ -n "${GC_TOKEN:-}" || -n "${GITCODE_TOKEN:-}" ]]; then
    : >"${TOKEN_LEAK:?}"
    exit 90
fi
if [[ "${1:-} ${2:-}" == "-m pip" ]]; then
    target="$(cd "$(dirname "$0")" && pwd)/pre-commit"
    printf '#!/usr/bin/env bash\nprintf "pre-commit 4.6.1\\n"\n' >"$target"
    chmod +x "$target"
fi
PYTHON
            chmod +x "$venv/bin/python"
        fi
        ;;
    *) ;;
esac
STUB
chmod +x "$STUB_BIN/tool-stub"
for name in go git make python3 cc; do
    ln -s "$STUB_BIN/tool-stub" "$STUB_BIN/$name"
done

run_setup() {
    env \
        PATH="$STUB_BIN:/usr/bin:/bin" \
        GC_DEV_TOOLS_DIR="$TOOLS_DIR" \
        CALL_LOG="$CALL_LOG" TOKEN_LEAK="$TOKEN_LEAK" CACHE_LOG="$CACHE_LOG" \
        GC_TOKEN="token-sentinel" GITCODE_TOKEN="token-sentinel-2" \
        GOCACHE="$TMP_ROOT/outside-cache" GOTMPDIR="$TMP_ROOT/outside-tmp" \
        bash "$ROOT_DIR/scripts/dev-setup.sh" "$@"
}

set +e
first_output="$(run_setup 2>&1)"
first_status=$?
set -e
[[ $first_status -eq 0 ]]
[[ "$first_output" == *"dev environment ready."* ]]
[[ ! -e "$TOKEN_LEAK" ]]
[[ -x "$TOOLS_DIR/bin/golangci-lint" ]]
[[ -x "$TOOLS_DIR/bin/gitleaks" ]]
[[ -L "$TOOLS_DIR/bin/pre-commit" ]]

: >"$CALL_LOG"
set +e
second_output="$(run_setup 2>&1)"
second_status=$?
set -e
[[ $second_status -eq 0 ]]
[[ "$second_output" != *"installing golangci-lint"* ]]
[[ "$second_output" != *"installing gitleaks"* ]]
[[ "$second_output" != *"installing pre-commit"* ]]
[[ ! -e "$TOKEN_LEAK" ]]

: >"$CACHE_LOG"
set +e
check_output="$(run_setup --check 2>&1)"
check_status=$?
set -e
[[ $check_status -eq 0 ]]
[[ "$check_output" == *"dev environment ready."* ]]
[[ ! -e "$TOKEN_LEAK" ]]
[[ -s "$CACHE_LOG" ]]
! grep -F "$TMP_ROOT/outside-cache" "$CACHE_LOG"
! grep -F "$TMP_ROOT/outside-tmp" "$CACHE_LOG"

rm "$TOOLS_DIR/bin/gitleaks"
set +e
missing_output="$(run_setup --check 2>&1)"
missing_status=$?
set -e
[[ $missing_status -eq 1 ]]
[[ "$missing_output" == *"MISS gitleaks v8.30.1"* ]]

printf '#!/usr/bin/env bash\nexit 0\n' >"$TOOLS_DIR/bin/gitleaks"
chmod +x "$TOOLS_DIR/bin/gitleaks"
before_hash="$(sha256sum "$TOOLS_DIR/bin/gitleaks" | awk '{print $1}')"
set +e
unmanaged_output="$(run_setup 2>&1)"
unmanaged_status=$?
set -e
after_hash="$(sha256sum "$TOOLS_DIR/bin/gitleaks" | awk '{print $1}')"
[[ $unmanaged_status -eq 1 ]]
[[ "$unmanaged_output" == *"refusing to overwrite unmanaged tool"* ]]
[[ "$before_hash" == "$after_hash" ]]

printf 'Unix dev setup install, idempotency, isolation, and guard tests passed.\n'
