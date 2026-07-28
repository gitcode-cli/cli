#!/usr/bin/env bash
# Install and verify the core host toolchain without a dev container.
#
# Usage:
#   bash scripts/dev-setup.sh          # install missing tools, then verify
#   bash scripts/dev-setup.sh --check  # offline verification; no persistent changes
#
# Packaging and release tools remain in .devcontainer/ and the documented
# release workflow. This script never needs GitCode credentials.

set -euo pipefail

unset GC_TOKEN GITCODE_TOKEN

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

GO_MIN_MAJOR=1
GO_MIN_MINOR=22
GOLANGCI_VERSION="v2.12.2"
GITLEAKS_VERSION="v8.30.1"
PRE_COMMIT_VERSION="4.6.1"
TOOLS_ROOT="${GC_DEV_TOOLS_DIR:-$HOME/.local/share/gitcode-cli/dev-tools}"
MANAGED_BIN="$TOOLS_ROOT/bin"
PRE_COMMIT_VENV="$TOOLS_ROOT/pre-commit"
CHECK_ONLY=0
MISSING=()

# Managed tools are available to this process only. No profile is modified.
export PATH="$MANAGED_BIN:$PATH"

case "${1:-}" in
    "") ;;
    --check) CHECK_ONLY=1 ;;
    -h|--help) sed -n '2,9p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
esac
[[ $# -le 1 ]] || { echo "too many arguments" >&2; exit 2; }

info() { printf '  %s\n' "$1"; }
ok() { printf 'ok   %s\n' "$1"; }
gap() { printf 'MISS %s\n' "$1"; MISSING+=("$1"); }
section() { printf '\n[%s]\n' "$1"; }
have() { command -v "$1" >/dev/null 2>&1; }

managed_or_path() {
    local name="$1"
    if [[ -x "$MANAGED_BIN/$name" ]]; then
        printf '%s\n' "$MANAGED_BIN/$name"
    else
        command -v "$name" 2>/dev/null || true
    fi
}

run_as_root() {
    if [[ "$(id -u)" -eq 0 ]]; then
        "$@"
    elif have sudo; then
        sudo -- "$@"
    else
        gap "sudo (required to install system packages)"
        return 1
    fi
}

go_version_ok() {
    local raw major minor
    raw="$(go env GOVERSION 2>/dev/null || true)"
    raw="${raw#go}"
    major="${raw%%.*}"
    minor="${raw#*.}"; minor="${minor%%.*}"
    [[ "$major" =~ ^[0-9]+$ && "$minor" =~ ^[0-9]+$ ]] || return 1
    (( major > GO_MIN_MAJOR || (major == GO_MIN_MAJOR && minor >= GO_MIN_MINOR) ))
}

compiler_available() {
    have cc || have gcc || have clang
}

install_system_dependencies() {
    (( CHECK_ONLY )) && return 0

    local need=0
    for tool in go git make bash python3; do have "$tool" || need=1; done
    compiler_available || need=1
    (( need )) || return 0

    section "System dependencies"
    case "$(uname -s)" in
        Darwin)
            if ! have brew; then
                gap "Homebrew (install from https://brew.sh)"
                return
            fi
            info "installing Go, Git, GNU Make, Bash, and Python with Homebrew"
            brew install go git make bash python
            if ! have make && have gmake; then
                mkdir -p "$MANAGED_BIN"
                ln -sfn "$(command -v gmake)" "$MANAGED_BIN/make"
            fi
            compiler_available || gap "C compiler (run: xcode-select --install)"
            ;;
        *)
            if have apt-get; then
                run_as_root apt-get update
                run_as_root apt-get install -y golang-go git make bash python3 python3-venv build-essential ca-certificates
            elif have dnf; then
                run_as_root dnf install -y golang git make bash python3 gcc gcc-c++ ca-certificates
            elif have yum; then
                run_as_root yum install -y golang git make bash python3 gcc gcc-c++ ca-certificates
            elif have pacman; then
                # Do not use pacman -Sy: partial upgrades are unsupported.
                run_as_root pacman -S --needed --noconfirm go git make bash python base-devel ca-certificates
            elif have zypper; then
                run_as_root zypper --non-interactive install go git make bash python3 gcc gcc-c++ ca-certificates
            else
                gap "supported package manager (apt, dnf, yum, pacman, zypper, or brew)"
            fi
            ;;
    esac
}

go_tool_version() {
    local binary="$1" module="$2"
    go version -m "$binary" 2>/dev/null | awk -v module="$module" '$1 == "mod" && $2 == module { print $3; exit }'
}

ensure_go_tool() {
    local name="$1" module="$2" version="$3" binary current
    binary="$(managed_or_path "$name")"
    current=""
    [[ -n "$binary" ]] && current="$(go_tool_version "$binary" "$module")"
    if [[ "$current" == "$version" ]]; then
        ok "$name $version ($binary)"
        return
    fi
    if (( CHECK_ONLY )); then
        gap "$name $version${current:+ (found $current)}"
        return
    fi
    mkdir -p "$MANAGED_BIN"
    info "installing $name $version into $MANAGED_BIN"
    if GOBIN="$MANAGED_BIN" go install "$module@$version"; then
        ok "$name $version ($MANAGED_BIN/$name)"
    else
        gap "$name $version (go install failed)"
    fi
}

pre_commit_version() {
    local binary="$1"
    "$binary" --version 2>/dev/null | awk '{print $2}'
}

ensure_pre_commit() {
    local binary current
    binary="$(managed_or_path pre-commit)"
    current=""
    [[ -n "$binary" ]] && current="$(pre_commit_version "$binary")"
    if [[ "$current" == "$PRE_COMMIT_VERSION" ]]; then
        ok "pre-commit $PRE_COMMIT_VERSION ($binary)"
        return
    fi
    if (( CHECK_ONLY )); then
        gap "pre-commit $PRE_COMMIT_VERSION${current:+ (found $current)}"
        return
    fi
    mkdir -p "$TOOLS_ROOT" "$MANAGED_BIN"
    info "installing pre-commit $PRE_COMMIT_VERSION into $PRE_COMMIT_VENV"
    if python3 -m venv "$PRE_COMMIT_VENV" &&
        "$PRE_COMMIT_VENV/bin/python" -m pip install --disable-pip-version-check "pre-commit==$PRE_COMMIT_VERSION" >/dev/null; then
        ln -sfn "$PRE_COMMIT_VENV/bin/pre-commit" "$MANAGED_BIN/pre-commit"
        ok "pre-commit $PRE_COMMIT_VERSION ($MANAGED_BIN/pre-commit)"
    else
        gap "pre-commit $PRE_COMMIT_VERSION (venv install failed)"
    fi
}

install_system_dependencies

section "Go toolchain"
if ! have go; then
    gap "go ${GO_MIN_MAJOR}.${GO_MIN_MINOR}+ (install from https://go.dev/dl/)"
elif ! go_version_ok; then
    if (( ! CHECK_ONLY )) && [[ "$(uname -s)" == "Darwin" ]] && have brew; then
        brew upgrade go || true
    fi
    go_version_ok || gap "go $(go env GOVERSION) is older than ${GO_MIN_MAJOR}.${GO_MIN_MINOR}; upgrade from https://go.dev/dl/"
else
    ok "go $(go env GOVERSION) ($(command -v go))"
fi

GO_READY=0
if have go && go_version_ok; then GO_READY=1; fi

if (( GO_READY )); then
    section "Module dependencies"
    if (( CHECK_ONLY )); then
        GOPROXY=off go list -mod=readonly -deps ./... >/dev/null 2>&1 \
            && ok "module dependencies resolved" \
            || gap "module dependencies (run without --check to download)"
    else
        go mod download && ok "module dependencies resolved"
    fi

    section "Verification tools"
    ensure_go_tool golangci-lint github.com/golangci/golangci-lint/v2 "$GOLANGCI_VERSION"
    ensure_go_tool gitleaks github.com/zricethezav/gitleaks/v8 "$GITLEAKS_VERSION"
fi

section "System tools"
for tool in git make bash python3; do
    have "$tool" && ok "$tool ($(command -v "$tool"))" || gap "$tool"
done
compiler_available && ok "C compiler" || gap "C compiler (required by go test -race)"

if have python3; then
    ensure_pre_commit
else
    gap "pre-commit (Python 3 unavailable)"
fi

section "Verification"
if (( GO_READY )); then
    smoke_dir="$(mktemp -d "${TMPDIR:-/tmp}/gc-dev-setup.XXXXXX")"
    trap 'rm -rf "$smoke_dir"' EXIT
    if (( CHECK_ONLY )); then
        build_cmd=(env GOPROXY=off go build -mod=readonly -o "$smoke_dir/gc" ./cmd/gc)
        race_cmd=(env GOPROXY=off go test -mod=readonly -race ./pkg/config)
    else
        build_cmd=(go build -o "$smoke_dir/gc" ./cmd/gc)
        race_cmd=(go test -race ./pkg/config)
    fi
    "${build_cmd[@]}" && "$smoke_dir/gc" version >/dev/null && ok "build and gc version" || gap "build or gc version"
    compiler_available && "${race_cmd[@]}" >/dev/null && ok "race-enabled test" || gap "race-enabled test"
    rm -rf "$smoke_dir"; trap - EXIT
fi

printf '\n'
if (( ${#MISSING[@]} )); then
    printf 'incomplete: %d dependency gap(s)\n' "${#MISSING[@]}"
    printf '  - %s\n' "${MISSING[@]}"
    printf '\nManaged tools directory: %s\n' "$MANAGED_BIN"
    printf 'Add it to PATH manually if you want these tools in future shells.\n'
    exit 1
fi

printf 'dev environment ready.\nManaged tools: %s\n' "$MANAGED_BIN"
printf 'This script made no shell profile changes. Add that directory to PATH manually if desired.\n'
