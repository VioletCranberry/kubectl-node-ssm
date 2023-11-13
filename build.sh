#!/bin/bash
set -euxo pipefail

snapshot='false'
while [[ $# -gt 0 ]]; do
    case "$1" in
        -s|--snapshot)
            snapshot='true'
            shift
            ;;
        *)
            echo "usage: $0 [-s | --snapshot]" >&2
            exit 1
            ;;
    esac
done

checksum_cmd=""
# common on macOS
if command -v shasum >/dev/null 2>&1; then
    checksum_cmd="shasum -a 256"
# common on Linux
elif command -v sha256sum >/dev/null 2>&1; then
    checksum_cmd="sha256sum"
fi

goreleaser_base_cmd='goreleaser release --timeout=10m --fail-fast --clean --skip=publish'
if [ "$snapshot" = "true" ]; then
    goreleaser_cmd="$goreleaser_base_cmd --snapshot"
else
    goreleaser_cmd="$goreleaser_base_cmd --auto-snapshot"
fi

$goreleaser_cmd
cd dist || exit 1

for archive in *.tar.gz *.zip; do
  if [ -f "$archive" ]; then
    $checksum_cmd "$archive" > "$archive.sha256"
  fi
done
