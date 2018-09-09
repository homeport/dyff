#!/usr/bin/env bash

# Copyright © 2018 Matthias Diester
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
VERSION="$(cd "$BASEDIR" && git describe --tags)"
VERFILE="$BASEDIR/internal/cmd/version.go"

SKIP_FULL_BUILD=0
SKIP_LOCAL_BUILD=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --only-local)
      SKIP_FULL_BUILD=1
      ;;

    --no-local)
      SKIP_LOCAL_BUILD=1
      ;;

    *)
      echo "unknown argument $1"
      exit 1
      ;;
  esac
  shift
done

on_exit() {
  if [[ -f "${VERFILE}.bak" ]]; then
    mv "${VERFILE}.bak" "${VERFILE}"
  fi
}

# Run on exit function to clean-up
trap on_exit EXIT

# Backup current version of the version subcommand and set current tag as version
cp -p "${VERFILE}" "${VERFILE}.bak"
perl -pi -e "s/const version = \"\\(development\\)\"/const version = \"${VERSION}\"/g" "${VERFILE}"

# Compile a local version into GOPATH bin if it exists
if [[ ${SKIP_LOCAL_BUILD} == 0 ]]; then
  if [[ -n ${GOPATH+x} ]]; then
    if [[ -d "${GOPATH}/bin" ]]; then
      echo -e "Compiling \\033[1mdyff version ${VERSION}\\033[0m for local machine to \\033[1m${GOPATH}/bin\\033[0m"
      (cd "${BASEDIR}/cmd/dyff/" && go install)
    fi
  fi
fi

# Stop here if skip full build flag was used
if [[ $SKIP_FULL_BUILD == 1 ]]; then
  exit 0
fi

# Compile all possible operating systems and architectures into the binaries directory (to be used for distribution)
TARGET_PATH="${BASEDIR}/binaries"
mkdir -p "$TARGET_PATH"
while read -r OS ARCH; do
  echo -e "Compiling \\033[1mdyff version ${VERSION}\\033[0m for OS \\033[1m${OS}\\033[0m and architecture \\033[1m${ARCH}\\033[0m"
  TARGET_FILE="${TARGET_PATH}/dyff-${OS}-${ARCH}"
  if [[ ${OS} == "windows" ]]; then
    TARGET_FILE="${TARGET_FILE}.exe"
  fi

  (cd "$BASEDIR" && GOOS="$OS" GOARCH="$ARCH" go build -ldflags='-s -w -extldflags "-static"' -o "$TARGET_FILE" cmd/dyff/main.go)

done <<EOL
darwin	386
darwin	amd64
freebsd	386
freebsd	amd64
linux	386
linux	amd64
linux	ppc64
linux	ppc64le
linux	s390x
windows	386
windows	amd64
EOL
