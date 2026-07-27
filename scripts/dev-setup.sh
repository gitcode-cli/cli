#!/usr/bin/env bash
# Install and verify local development dependencies on the host.
#
# This is the container-free path: it targets the same toolchain baseline as
# .devcontainer/ but installs into the host so AI agents and contributors can
# run local verification without pulling an image.
#
# Usage:
#   ./scripts/dev-setup.sh            # install missing dependencies, then verify
#   ./scripts/dev-setup.sh --check    # verify only, install nothing (exit 1 if gaps)
#   ./scripts/dev-setup.sh --with-packaging   # also install nfpm / goreleaser / python build
#
# Scope: this script installs tooling only. It never touches credentials and
# never reads or prints GC_TOKEN / GITCODE_TOKEN.

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

GO_MIN_MAJOR=1
GO_MIN_MINOR=22
NFPM_VERSION="v2.40.0"
GOLANGCI_VERSION="v2.12.2"

CHECK_ONLY=0
WITH_PACKAGING=0

for arg in "$@"; do
    case "$arg" in
        --check) CHECK_ONLY=1 ;;
        --with-packaging) WITH_PACKAGING=1 ;;
        -h|--help)
            sed -n '2,15p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'
            exit 0
            ;;
        *)
            echo "unknown argument: $arg" >&2
            exit 2
            ;;
    esac
done

MISSING=()

info() { printf '  %s\n' "$1"; }
ok() { printf 'ok   %s\n' "$1"; }
gap() { printf 'MISS %s\n' "$1"; MISSING+=("$1"); }
section() { printf '\n[%s]\n' "$1"; }

have() { command -v "$1" >/dev/null 2>&1; }

go_bin_dir() {
    local gobin
    gobin="$(go env GOBIN 2>/dev/null || true)"
    if [[ -n "$gobin" ]]; then
        printf '%s\n' "$gobin"
    else
        printf '%s/bin\n' "$(go env GOPATH)"
    fi
}

# Go tools land in GOPATH/bin, which is often absent from PATH on fresh hosts.
ensure_go_bin_on_path() {
    local dir
    dir="$(go_bin_dir)"
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$dir:$PATH" ;;
    esac
}

go_version_ok() {
    local raw major minor
    raw="$(go env GOVERSION 2>/dev/null || true)"
    raw="${raw#go}"
    major="${raw%%.*}"
    minor="${raw#*.}"
    minor="${minor%%.*}"
    [[ -n "$major" && -n "$minor" ]] || return 1
    (( major > GO_MIN_MAJOR )) && return 0
    (( major == GO_MIN_MAJOR && minor >= GO_MIN_MINOR ))
}

install_go_tool() {
    local name="$1" module="$2"
    if have "$name"; then
        ok "$name ($(command -v "$name"))"
        return 0
    fi
    if (( CHECK_ONLY )); then
        gap "$name"
        return 0
    fi
    info "installing $name from $module"
    if go install "$module"; then
        ensure_go_bin_on_path
        if have "$name"; then
            ok "$name ($(command -v "$name"))"
        else
            gap "$name (installed but not on PATH; add $(go_bin_dir) to PATH)"
        fi
    else
        gap "$name (go install failed)"
    fi
}

section "Go toolchain"
if ! have go; then
    gap "go (install Go ${GO_MIN_MAJOR}.${GO_MIN_MINOR}+ from https://go.dev/dl/ or your package manager)"
elif ! go_version_ok; then
    gap "go $(go env GOVERSION) is older than ${GO_MIN_MAJOR}.${GO_MIN_MINOR}"
else
    ok "go $(go env GOVERSION) ($(command -v go))"
    ensure_go_bin_on_path
fi

if have go; then
    section "Module dependencies"
    if (( CHECK_ONLY )); then
        if go list -deps ./... >/dev/null 2>&1; then
            ok "module dependencies resolved"
        else
            gap "module dependencies (run: go mod download)"
        fi
    else
        info "go mod download"
        go mod download
        ok "module dependencies resolved"
    fi

    section "Lint tooling"
    install_go_tool golangci-lint \
        "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_VERSION}"

    section "Secret scanning"
    install_go_tool gitleaks "github.com/zricethezav/gitleaks/v8@latest"

    if (( WITH_PACKAGING )); then
        section "Packaging tooling"
        install_go_tool nfpm "github.com/goreleaser/nfpm/v2/cmd/nfpm@${NFPM_VERSION}"
        install_go_tool goreleaser "github.com/goreleaser/goreleaser/v2@latest"
    fi
fi

section "Shell and build utilities"
for tool in git make bash; do
    if have "$tool"; then
        ok "$tool ($(command -v "$tool"))"
    else
        gap "$tool (install via your package manager)"
    fi
done

section "Python tooling"
PY=""
if have python3; then
    PY="python3"
elif have python; then
    PY="python"
fi

if [[ -z "$PY" ]]; then
    gap "python3 (required by scripts/validate-ai-record.py and friends)"
else
    ok "$PY ($("$PY" --version 2>&1))"
    if "$PY" -m pre_commit --version >/dev/null 2>&1 || have pre-commit; then
        ok "pre-commit"
    elif (( CHECK_ONLY )); then
        gap "pre-commit"
    else
        info "installing pre-commit"
        # Debian/Ubuntu mark the system interpreter as externally managed;
        # --user keeps us out of the way of the distro package manager.
        if "$PY" -m pip install --user --upgrade pre-commit >/dev/null 2>&1; then
            ok "pre-commit"
        else
            gap "pre-commit (pip install failed; try pipx install pre-commit)"
        fi
    fi

    if (( WITH_PACKAGING )); then
        if "$PY" -m build --version >/dev/null 2>&1; then
            ok "python build"
        elif (( CHECK_ONLY )); then
            gap "python build"
        else
            info "installing build/wheel/setuptools"
            if "$PY" -m pip install --user --upgrade build wheel setuptools >/dev/null 2>&1; then
                ok "python build"
            else
                gap "python build (pip install failed)"
            fi
        fi
    fi
fi

section "Optional tooling"
check_optional() {
    local tool="$1" why="$2"
    if have "$tool"; then
        ok "$tool ($(command -v "$tool"))"
    else
        info "optional: $tool not found ($why)"
    fi
}
check_optional gh "needed to inspect CI runs"
check_optional docker "needed for make docker-* targets"
check_optional rpmbuild "needed for local RPM packaging"

section "Verification"
if have go && go_version_ok; then
    info "go build -o ./gc ./cmd/gc"
    if go build -o ./gc ./cmd/gc; then
        ok "build"
        if ./gc version >/dev/null 2>&1; then
            ok "gc version"
        else
            gap "gc version failed to run"
        fi
        rm -f ./gc
    else
        gap "go build failed"
    fi
else
    info "skipped (Go unavailable)"
fi

printf '\n'
if (( ${#MISSING[@]} > 0 )); then
    printf 'incomplete: %d dependency gap(s)\n' "${#MISSING[@]}"
    for item in "${MISSING[@]}"; do
        printf '  - %s\n' "$item"
    done
    printf '\nGo tools install into %s; make sure it is on PATH.\n' \
        "$(have go && go_bin_dir || echo '$(go env GOPATH)/bin')"
    exit 1
fi

cat <<'EOF'
dev environment ready.

Next steps:
  go test ./...
  go build -o ./gc ./cmd/gc && ./gc version
  ./scripts/regression-core.sh

Real command verification needs GC_TOKEN (or GITCODE_TOKEN) exported by hand
and must target infra-test/* repositories only.
EOF
