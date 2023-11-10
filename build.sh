#!/bin/bash
set -euxo pipefail

checksum_command=""
# common on macOS
if command -v shasum >/dev/null 2>&1; then
    checksum_command="shasum -a 256"
# common on Linux
elif command -v sha256sum >/dev/null 2>&1; then
    checksum_command="sha256sum"
fi

goreleaser release --timeout=1m --auto-snapshot --fail-fast --clean --skip=publish
cd dist || exit 1

for archive in *.tar.gz *.zip; do
  if [ -f "$archive" ]; then
    $checksum_command "$archive" > "$archive.sha256"
  fi
done
