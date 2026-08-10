#!/usr/bin/env bash
# Assemble the npm package exclusively from the verified GoReleaser binaries.

set -euo pipefail

VERSION="${1:?usage: prepare-npm-package.sh VERSION DIST_DIR OUTPUT_DIR}"
DIST_DIR="${2:?usage: prepare-npm-package.sh VERSION DIST_DIR OUTPUT_DIR}"
OUTPUT_DIR="${3:?usage: prepare-npm-package.sh VERSION DIST_DIR OUTPUT_DIR}"
PLATFORMS_DIR="npm/bin/platforms"
mkdir -p "${OUTPUT_DIR}"
OUTPUT_DIR="$(cd "${OUTPUT_DIR}" && pwd)"

extract_tar_binary() {
    local archive="$1" target="$2" temp_dir
    temp_dir="$(mktemp -d)"
    tar -xzf "${archive}" -C "${temp_dir}"
    install -m 0755 "${temp_dir}/gc" "${target}"
    rm -rf "${temp_dir}"
}

(
    cd npm
    npm version "${VERSION}" --no-git-tag-version --allow-same-version --ignore-scripts
)
mkdir -p "${PLATFORMS_DIR}"
install -m 0755 "${DIST_DIR}/gc_linux_amd64" "${PLATFORMS_DIR}/gc-linux-amd64"
install -m 0755 "${DIST_DIR}/gc_linux_arm64" "${PLATFORMS_DIR}/gc-linux-arm64"
extract_tar_binary "${DIST_DIR}/gc_${VERSION}_darwin_amd64.tar.gz" "${PLATFORMS_DIR}/gc-darwin-amd64"
extract_tar_binary "${DIST_DIR}/gc_${VERSION}_darwin_arm64.tar.gz" "${PLATFORMS_DIR}/gc-darwin-arm64"
unzip -p "${DIST_DIR}/gc_${VERSION}_windows_amd64.zip" gc.exe > "${PLATFORMS_DIR}/gc-windows-amd64.exe"

(
    cd npm
    npm test
    npm pack --pack-destination "${OUTPUT_DIR}"
)
