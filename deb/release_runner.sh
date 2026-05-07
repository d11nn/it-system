#!/usr/bin/env bash

set -euo pipefail

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <version-tag>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION_TAG="$1"

# Debian version cannot start with non-digit, so normalize a leading v/V.
DEB_VERSION="${VERSION_TAG#v}"
DEB_VERSION="${DEB_VERSION#V}"

if [[ ! "$DEB_VERSION" =~ ^[0-9][0-9A-Za-z.+:~_-]*$ ]]; then
    echo "Error: invalid version tag '$VERSION_TAG'"
    echo "Allowed input examples: v1.2.3, 1.2.3, 2.0.0-rc1"
    exit 1
fi

PKG_NAME="it-system-runner"
PKG_ARCH="amd64"
PKGROOT="pkg/${PKG_NAME}_${DEB_VERSION}_${PKG_ARCH}"
DEB_TEMPLATE_DIR="deb/runner"
DEB_OUTPUT="${PKGROOT}.deb"

echo "[+] Building runner..."
make runner

if [[ ! -f build/runner ]]; then
    echo "Error: build/runner not found"
    exit 1
fi

echo "[+] Preparing package root: ${PKGROOT}"
rm -rf "$PKGROOT"
mkdir -p \
    "$PKGROOT/DEBIAN" \
    "$PKGROOT/usr/bin" \
    "$PKGROOT/etc/it-system" \
    "$PKGROOT/lib/systemd/system" \
    "$PKGROOT/usr/share/it-system-runner"

echo "[+] Copying binary/config/service..."
install -m 755 build/runner "$PKGROOT/usr/bin/it-system-runner"
install -m 644 config-runner.yaml "$PKGROOT/etc/it-system/config-runner.yaml"
install -m 644 "$DEB_TEMPLATE_DIR/lib/systemd/system/it-system-runner.service" "$PKGROOT/lib/systemd/system/it-system-runner.service"

echo "[+] Preparing DEBIAN metadata..."
install -m 644 "$DEB_TEMPLATE_DIR/DEBIAN/control" "$PKGROOT/DEBIAN/control"
install -m 644 "$DEB_TEMPLATE_DIR/DEBIAN/conffiles" "$PKGROOT/DEBIAN/conffiles"
install -m 755 "$DEB_TEMPLATE_DIR/DEBIAN/postinst" "$PKGROOT/DEBIAN/postinst"
install -m 755 "$DEB_TEMPLATE_DIR/DEBIAN/prerm" "$PKGROOT/DEBIAN/prerm"

sed -i "s/<version>/${DEB_VERSION}/g" "$PKGROOT/DEBIAN/control"

echo "[+] Building deb package..."
dpkg-deb --build "$PKGROOT" "$DEB_OUTPUT"

echo "[✔] Done: ${DEB_OUTPUT}"
