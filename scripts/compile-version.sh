#!/usr/bin/env bash

# Copyright © 2019 The Homeport Team
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="$(cd "${BASEDIR}" && git describe --tags)"
TARGET_PATH="${BASEDIR}/binaries"

function build-binary() {
  OS="$1"
  ARCH="$2"

  echo -e "Compiling \\033[1mdyff version ${VERSION}\\033[0m for OS \\033[1m${OS}\\033[0m and architecture \\033[1m${ARCH}\\033[0m"
  TARGET_FILE="${TARGET_PATH}/dyff-${OS}-${ARCH}"
  if [[ ${OS} == "windows" ]]; then
    TARGET_FILE="${TARGET_FILE}.exe"
  fi

  GO111MODULE=on CGO_ENABLED=0 GOOS="$OS" GOARCH="$ARCH" go build \
    -tags netgo \
    -ldflags="-s -w -extldflags '-static' -X github.com/homeport/dyff/internal/cmd.version=${VERSION}" \
    -o "$TARGET_FILE" \
    "${BASEDIR}/cmd/dyff/main.go"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --local)
      build-binary "$(uname | tr '[:upper:]' '[:lower:]')" "$(uname -m | sed 's/x86_64/amd64/')"
      ;;

    --all)
      build-binary darwin amd64
      build-binary linux amd64
      ;;

    *)
      echo "unknown argument $1"
      exit 1
      ;;
  esac
  shift
done

if hash file >/dev/null 2>&1; then
  echo -e '\n\033[1mFile details of compiled binaries:\033[0m'
  file binaries/*
fi

if hash shasum >/dev/null 2>&1; then
  echo -e '\n\033[1mSHA sum of compiled binaries:\033[0m'
  shasum --algorithm 256 binaries/*

elif hash sha1sum >/dev/null 2>&1; then
  echo -e '\n\033[1mSHA sum of compiled binaries:\033[0m'
  sha1sum binaries/*
  echo
fi
