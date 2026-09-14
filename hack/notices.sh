#!/usr/bin/env bash
# Generates THIRD_PARTY_NOTICES from the licences of every module linked into
# the nosleep binary. The Go standard library is covered by Go's own licence and
# is not a third party; test-only dependencies are not distributed and so are
# not included.
#
# Run from the repository root:
#
#	hack/notices.sh > THIRD_PARTY_NOTICES
#
# CI checks that the committed file is up to date.
set -euo pipefail

cd "$(dirname "$0")/.."

cat <<'HEADER'
THIRD-PARTY NOTICES

The nosleep binary links the following modules, reproduced here with their
licences as those licences require.

HEADER

go list -deps -f '{{if .Module}}{{.Module.Path}}{{"\t"}}{{.Module.Version}}{{"\t"}}{{.Module.Dir}}{{end}}' ./cmd/nosleep \
  | grep -v '^github.com/elliotwms/nosleep' \
  | sort -u \
  | while IFS=$'\t' read -r path version dir; do
      license=""
      for candidate in LICENSE LICENSE.txt LICENSE.md LICENCE COPYING; do
        if [ -f "$dir/$candidate" ]; then
          license="$dir/$candidate"
          break
        fi
      done

      if [ -z "$license" ]; then
        echo "no licence file found for $path@$version in $dir" >&2
        exit 1
      fi

      echo "--------------------------------------------------------------------------------"
      echo "$path $version"
      echo "--------------------------------------------------------------------------------"
      echo
      cat "$license"
      echo
    done
