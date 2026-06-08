#!/usr/bin/env bash
# Test release workflow version input validation.

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
validator="${script_dir}/validate-release-version.sh"
failures=0

fail() {
    echo "FAIL: $*" >&2
    failures=$((failures + 1))
}

run_valid() {
    local input="$1"
    local expected="$2"
    local output

    if ! output="$("${validator}" "${input}" 2>&1)"; then
        fail "expected valid version '${input}', got: ${output}"
        return
    fi

    if [[ "${output}" != "${expected}" ]]; then
        fail "expected '${input}' to normalize to '${expected}', got '${output}'"
    fi
}

run_invalid() {
    local input="$1"
    local output

    if output="$("${validator}" "${input}" 2>&1)"; then
        fail "expected invalid version '${input}', got normalized output: ${output}"
    fi
}

if output="$("${validator}" 2>&1)"; then
    fail "expected missing version argument to fail, got: ${output}"
fi

run_valid "v1.2.3" "1.2.3"
run_valid "1.2.3" "1.2.3"
run_valid "v1.2.3-beta.1" "1.2.3-beta.1"
run_valid "1.2.3-rc.1" "1.2.3-rc.1"
run_valid "v10.20.30-alpha.1" "10.20.30-alpha.1"
run_valid "v1.2.3-alpha.0" "1.2.3-alpha.0"

run_invalid ""
run_invalid "v01.2.3"
run_invalid "v1.02.3"
run_invalid "v1.2.03"
run_invalid "v1"
run_invalid "v1.2"
run_invalid "v1.2.3+build"
run_invalid "v1.2.3-"
run_invalid "v1.2.3-alpha"
run_invalid "v1.2.3-beta.01"
run_invalid "v1.2.3-hotfix.1"
run_invalid "v1.2.3-alpha..1"
run_invalid "v1.2.3-alpha.-1"
run_invalid "v1.2.3-alpha."
run_invalid "v1.2.3;echo bad"
run_invalid $'v1.2.3\necho bad'
run_invalid "v1.2.3/../../x"

if ((failures > 0)); then
    exit 1
fi

echo "release version validation tests passed"
