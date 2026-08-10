#!/usr/bin/env bash
# Keep every package-manager metadata file aligned to the repository VERSION.

set -euo pipefail

RAW_VERSION="${1:?usage: sync-package-version.sh VERSION [--check]}"
MODE="${2:-}"
VERSION="$(bash scripts/validate-release-version.sh "${RAW_VERSION}")"

check_versions() {
    test "$(tr -d '\r\n' < VERSION)" = "${VERSION}"
    grep -Fq "version: \"${VERSION}\"" nfpm-amd64.yaml
    grep -Fq "version: \"${VERSION}\"" nfpm-arm64.yaml
    grep -Fq "version = \"${VERSION}\"" pyproject.toml
    grep -Fq "__version__ = \"${VERSION}\"" gc_cli/__init__.py
    grep -Fq "\"version\": \"${VERSION}\"" npm/package.json
}

if [[ "${MODE}" == "--check" ]]; then
    check_versions
    printf 'package versions are aligned at %s\n' "${VERSION}"
    exit 0
fi
if [[ -n "${MODE}" ]]; then
    printf 'unknown argument: %s\n' "${MODE}" >&2
    exit 2
fi

printf '%s\n' "${VERSION}" > VERSION
sed -i "s/version: \".*\"/version: \"${VERSION}\"/" nfpm-amd64.yaml nfpm-arm64.yaml
sed -i "s/version = \".*\"/version = \"${VERSION}\"/" pyproject.toml
sed -i "s/__version__ = \".*\"/__version__ = \"${VERSION}\"/" gc_cli/__init__.py
sed -i "0,/\"version\": \".*\"/s//\"version\": \"${VERSION}\"/" npm/package.json
check_versions
printf 'updated package versions to %s\n' "${VERSION}"
