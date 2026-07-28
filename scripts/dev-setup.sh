#!/usr/bin/env bash
# Install and verify local development dependencies on the host.
#
# This is the container-free path: it installs the core verification baseline
# into the host so AI agents and contributors can work without pulling an
# image. Pass --with-packaging to add the container's packaging toolchain.
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
GITLEAKS_VERSION="v8.30.1"
GORELEASER_VERSION="v2.17.1"
PRE_COMMIT_VERSION="4.6.1"

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

run_as_root() {
    if [[ "$(id -u)" -eq 0 ]]; then
        "$@"
    elif have sudo; then
        sudo "$@"
    else
        gap "sudo (required to install system packages)"
        return 1
    fi
}

install_system_dependencies() {
    local missing=0
    for tool in go git make bash python3; do
        have "$tool" || missing=1
    done
    if have go && ! go_version_ok; then
        missing=1
    fi
    (( missing )) || return 0
    (( CHECK_ONLY )) && return 0

    section "System dependencies"
    if [[ "$(uname -s)" == "Darwin" ]]; then
        if ! have brew; then
            gap "Homebrew (install from https://brew.sh, then rerun this script)"
            return 0
        fi
        info "installing Go, Git, Make, Bash, Python, and pipx with Homebrew"
        brew install go git make bash python pipx
        make_gnubin="$(brew --prefix make)/libexec/gnubin"
        [[ -d "$make_gnubin" ]] && persist_path_dir "$make_gnubin"
        return 0
    fi

    if have apt-get; then
        info "installing base tools with apt-get"
        run_as_root apt-get update
        run_as_root apt-get install -y golang-go git make bash python3 python3-pip curl ca-certificates
        run_as_root apt-get install -y pipx >/dev/null 2>&1 || true
    elif have dnf; then
        info "installing base tools with dnf"
        run_as_root dnf install -y golang git make bash python3 python3-pip curl ca-certificates
        run_as_root dnf install -y pipx >/dev/null 2>&1 || true
    elif have yum; then
        info "installing base tools with yum"
        run_as_root yum install -y golang git make bash python3 python3-pip curl ca-certificates
    elif have pacman; then
        info "installing base tools with pacman"
        run_as_root pacman -Sy --needed --noconfirm go git make bash python python-pip curl ca-certificates
        run_as_root pacman -S --needed --noconfirm python-pipx >/dev/null 2>&1 || true
    elif have zypper; then
        info "installing base tools with zypper"
        run_as_root zypper --non-interactive install go git make bash python3 python3-pip curl ca-certificates
    else
        gap "supported package manager (apt, dnf, yum, pacman, zypper, or brew)"
    fi
}

persist_path_dir() {
    local dir="$1" profile="$HOME/.profile" line
    [[ "$(uname -s)" == "Darwin" ]] && profile="$HOME/.zprofile"
    line="export PATH=\"$dir:\$PATH\""
    touch "$profile"
    grep -Fqx "$line" "$profile" || printf '\n%s\n' "$line" >>"$profile"
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$dir:$PATH" ;;
    esac
}

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
    if (( CHECK_ONLY )); then
        case ":$PATH:" in
            *":$dir:"*) ;;
            *) export PATH="$dir:$PATH" ;;
        esac
    else
        persist_path_dir "$dir"
    fi
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

install_system_dependencies

section "Go toolchain"
if ! have go; then
    gap "go (install Go ${GO_MIN_MAJOR}.${GO_MIN_MINOR}+ from https://go.dev/dl/ or your package manager)"
elif ! go_version_ok; then
    gap "go $(go env GOVERSION) is older than ${GO_MIN_MAJOR}.${GO_MIN_MINOR}; upgrade from https://go.dev/dl/ or your package manager"
else
    ok "go $(go env GOVERSION) ($(command -v go))"
    ensure_go_bin_on_path
fi

if have go; then
    section "Module dependencies"
    if (( CHECK_ONLY )); then
        if GOPROXY=off go list -mod=readonly -deps ./... >/dev/null 2>&1; then
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
    install_go_tool gitleaks "github.com/zricethezav/gitleaks/v8@${GITLEAKS_VERSION}"

    if (( WITH_PACKAGING )); then
        section "Packaging tooling"
        install_go_tool nfpm "github.com/goreleaser/nfpm/v2/cmd/nfpm@${NFPM_VERSION}"
        install_go_tool goreleaser "github.com/goreleaser/goreleaser/v2@${GORELEASER_VERSION}"
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
        if have pipx && pipx install --force "pre-commit==${PRE_COMMIT_VERSION}" >/dev/null 2>&1; then
            have pre-commit || persist_path_dir "$HOME/.local/bin"
            ok "pre-commit"
        elif "$PY" -m pip install --user --upgrade "pre-commit==${PRE_COMMIT_VERSION}" >/dev/null 2>&1; then
            user_bin="$($PY -c 'import site,os; print(os.path.join(site.USER_BASE, "bin"))')"
            persist_path_dir "$user_bin"
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
    smoke_bin="$(mktemp "${TMPDIR:-/tmp}/gc-dev-setup.XXXXXX")"
    trap 'rm -f "$smoke_bin"' EXIT
    info "go build -o $smoke_bin ./cmd/gc"
    if (( CHECK_ONLY )); then
        build_cmd=(env GOPROXY=off go build -mod=readonly -o "$smoke_bin" ./cmd/gc)
    else
        build_cmd=(go build -o "$smoke_bin" ./cmd/gc)
    fi
    if "${build_cmd[@]}"; then
        ok "build"
        if "$smoke_bin" version >/dev/null 2>&1; then
            ok "gc version"
        else
            gap "gc version failed to run"
        fi
        rm -f "$smoke_bin"
        trap - EXIT
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
