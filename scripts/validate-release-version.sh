#!/usr/bin/env bash
# Validate a release workflow version input and print the package version.

set -euo pipefail

version="${1:-}"

if [[ -z "${version}" ]]; then
    echo "Invalid release version: value is required" >&2
    exit 1
fi

if [[ ! "${version}" =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z][0-9A-Za-z.-]*)?$ ]]; then
    echo "Invalid release version: ${version}" >&2
    echo "Expected format: vMAJOR.MINOR.PATCH, MAJOR.MINOR.PATCH, or a prerelease variant" >&2
    exit 1
fi

printf '%s\n' "${version#v}"
